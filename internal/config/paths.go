// @anchors
//   code: PTHCD
//   ref: PTHSA

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// Aliases are the short names a path of anchors.yaml may open with, each with the path it
// stands for: `"@guides": docs/engineering/guides` lets `governs[].from` read
// `@guides/SCREEN_GUIDE.md`. A long path written in a dozen fields drifted the day one of
// them was moved; written once, it moves in one line.
//
// The fields that hold a path say so in a `path:` struct tag, and ExpandPaths, called by
// Load, is the one place that reads it: no consumer of a field translates an alias, so none
// can forget to.
type Aliases map[string]string

// PathKind is what a field's path is: a file the Anchors reads, a glob over the project's
// files, or a file the project or the Anchors writes.
type PathKind string

const (
	PathRead  PathKind = "read"
	PathGlob  PathKind = "glob"
	PathWrite PathKind = "write"
)

var aliasNameRE = regexp.MustCompile(`^@[A-Za-z0-9_.-]+$`)

// validate refuses an alias whose name does not open with `@`, whose path is empty, a URL
// or itself opens with an alias.
func (a Aliases) validate() error {
	for _, name := range a.names() {
		value := a[name]
		switch {
		case !aliasNameRE.MatchString(name):
			return fmt.Errorf("`aliases.%s`: an alias is named `@` and letters, digits, `_`, `.` or `-` (`@guides`)", name)
		case strings.TrimSpace(value) == "":
			return fmt.Errorf("`aliases.%s` names no path", name)
		case isURL(value):
			return fmt.Errorf("`aliases.%s` is a URL (%s): a path of anchors.yaml names a file of the project", name, value)
		case strings.HasPrefix(value, "@"):
			return fmt.Errorf("`aliases.%s` opens with another alias (%s): write the path it stands for", name, value)
		}
	}
	return nil
}

func (a Aliases) names() []string {
	out := make([]string, 0, len(a))
	for name := range a {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Expand replaces the alias that opens p — its first segment — with the alias's path; a
// p that opens with none is returned as written. ok is false when the first segment opens
// with `@`, names no alias and no entry of the root has that name: a typo, which would
// otherwise make a glob match nothing without a word.
func (a Aliases) Expand(root, p string) (string, bool) {
	if !strings.HasPrefix(p, "@") {
		return p, true
	}
	head, rest, _ := strings.Cut(p, "/")
	value, declared := a[head]
	if !declared {
		_, err := os.Stat(filepath.Join(root, head))
		return p, err == nil
	}
	value = strings.TrimSuffix(value, "/")
	if rest == "" {
		return value, true
	}
	return value + "/" + rest, true
}

func isURL(p string) bool {
	return strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://")
}

// ExpandPaths translates the aliases in every field of c that holds a path — the ones with
// a `path:` tag —, naming the field of a path that cites an undeclared alias.
func (c *Config) ExpandPaths(root string) error {
	if err := c.Aliases.validate(); err != nil {
		return err
	}
	return walkPaths(reflect.ValueOf(c).Elem(), "", func(field string, _ PathKind, p string) (string, error) {
		out, ok := c.Aliases.Expand(root, p)
		if !ok {
			head, _, _ := strings.Cut(p, "/")
			return "", fmt.Errorf("`%s` cites the alias `%s`, which `aliases:` does not declare", field, head)
		}
		return out, nil
	})
}

// walkPaths calls f on every path held by a field with a `path:` tag under v, and keeps
// what it returns. The tag is the kind, and `,keys` when the paths are the keys of a map
// (`no_signal:`, a glob → its reason). It descends into structs, pointers, slices and maps.
func walkPaths(v reflect.Value, field string, f func(field string, kind PathKind, p string) (string, error)) error {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return walkPaths(v.Elem(), field, f)
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if !sf.IsExported() {
				continue
			}
			name := yamlName(sf)
			if name == "" {
				continue
			}
			sub := join(field, name)
			if tag, ok := sf.Tag.Lookup("path"); ok {
				kind, opt, _ := strings.Cut(tag, ",")
				if err := applyPaths(v.Field(i), sub, PathKind(kind), opt == "keys", f); err != nil {
					return err
				}
				continue
			}
			if err := walkPaths(v.Field(i), sub, f); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := walkPaths(v.Index(i), fmt.Sprintf("%s[%d]", field, i), f); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range sortedKeys(v) {
			elem := reflect.New(v.Type().Elem()).Elem()
			elem.Set(v.MapIndex(k))
			if err := walkPaths(elem, join(field, k.String()), f); err != nil {
				return err
			}
			v.SetMapIndex(k, elem)
		}
	}
	return nil
}

// applyPaths applies f to the paths of a tagged field: a string, a list of them, the values
// of a map of either or, with keys, the keys of a map.
func applyPaths(v reflect.Value, field string, kind PathKind, keys bool, f func(string, PathKind, string) (string, error)) error {
	switch {
	case v.Kind() == reflect.String:
		if v.String() == "" {
			return nil
		}
		out, err := f(field, kind, v.String())
		if err != nil {
			return err
		}
		v.SetString(out)
	case v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String:
		for i := 0; i < v.Len(); i++ {
			if err := applyPaths(v.Index(i), fmt.Sprintf("%s[%d]", field, i), kind, false, f); err != nil {
				return err
			}
		}
	case v.Kind() == reflect.Map && keys:
		for _, k := range sortedKeys(v) {
			out, err := f(join(field, k.String()), kind, k.String())
			if err != nil {
				return err
			}
			if out != k.String() {
				elem := v.MapIndex(k)
				v.SetMapIndex(k, reflect.Value{})
				v.SetMapIndex(reflect.ValueOf(out).Convert(v.Type().Key()), elem)
			}
		}
	case v.Kind() == reflect.Map:
		for _, k := range sortedKeys(v) {
			elem := reflect.New(v.Type().Elem()).Elem()
			elem.Set(v.MapIndex(k))
			if err := applyPaths(elem, join(field, k.String()), kind, false, f); err != nil {
				return err
			}
			v.SetMapIndex(k, elem)
		}
	default:
		return fmt.Errorf("`%s`: a `path:` tag on a %s, which holds no path", field, v.Type())
	}
	return nil
}

func yamlName(sf reflect.StructField) string {
	name, _, _ := strings.Cut(sf.Tag.Get("yaml"), ",")
	if name == "-" {
		return ""
	}
	if name == "" {
		return strings.ToLower(sf.Name)
	}
	return name
}

func join(field, name string) string {
	if field == "" {
		return name
	}
	return field + "." + name
}

func sortedKeys(v reflect.Value) []reflect.Value {
	keys := v.MapKeys()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	return keys
}
