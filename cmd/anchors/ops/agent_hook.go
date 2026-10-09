// @anchors
//   code: AGHKG
//   ref: INHKN

package ops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// agentHookCommand is what the agent's hook runs: it records the agent's long commands for
// `anchors monitor` (DESIGN-process-monitor.md).
const agentHookCommand = "anchors monitor hook"

// agentSettingsPath is where the hook goes: the project's Claude Code settings, shared with
// the team, or — with user — the user's own, which no project versions. A user-level hook
// runs in every folder, and records nothing outside a project with an anchors.yaml.
func agentSettingsPath(root string, user bool) (string, error) {
	if user {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("no home folder to write the user's settings in: %w", err)
		}
		return filepath.Join(home, ".claude", "settings.json"), nil
	}
	return filepath.Join(root, ".claude", "settings.json"), nil
}

// installAgentHook adds, to the Claude Code settings at path, the hook that runs before and
// after each command of the agent and records the long ones. Everything else in the settings
// is kept as it was, and a hook already there is not added twice. It says what it did.
func installAgentHook(path string) (string, error) {
	settings := map[string]any{}
	if b, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(b, &settings); err != nil {
			return "", fmt.Errorf("%s is not valid JSON (%v) — fix it, then install again", path, err)
		}
	}
	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
	}
	added := 0
	for _, event := range []string{"PreToolUse", "PostToolUse"} {
		list, _ := hooks[event].([]any)
		if hasAgentHook(list) {
			continue
		}
		list = append(list, map[string]any{
			"matcher": "Bash",
			"hooks":   []any{map[string]any{"type": "command", "command": agentHookCommand, "timeout": 10}},
		})
		hooks[event] = list
		added++
	}
	if added == 0 {
		return fmt.Sprintf("✓ agent hook already in %s", filepath.ToSlash(path)), nil
	}
	settings["hooks"] = hooks
	b, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("✓ agent hook written to %s — the agent's long commands are recorded for `anchors monitor`", filepath.ToSlash(path)), nil
}

// hasAgentHook says whether a list of hook entries already runs the agent hook.
func hasAgentHook(list []any) bool {
	for _, e := range list {
		entry, _ := e.(map[string]any)
		inner, _ := entry["hooks"].([]any)
		for _, h := range inner {
			hm, _ := h.(map[string]any)
			if c, _ := hm["command"].(string); strings.Contains(c, agentHookCommand) {
				return true
			}
		}
	}
	return false
}
