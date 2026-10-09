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

// agentSettingsPath is the project's Claude Code settings, shared with the team.
func agentSettingsPath(root string) string {
	return filepath.Join(root, ".claude", "settings.json")
}

// installAgentHook adds, to the project's Claude Code settings, the hook that runs before and
// after each command of the agent and records the long ones. Everything else in the settings
// is kept as it was, and a hook already there is not added twice. It says what it did.
func installAgentHook(root string) (string, error) {
	path := agentSettingsPath(root)
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
		return fmt.Sprintf("✓ agent hook already in %s", filepath.ToSlash(rel(root, path))), nil
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
	return fmt.Sprintf("✓ agent hook written to %s — the agent's long commands are recorded for `anchors monitor`", filepath.ToSlash(rel(root, path))), nil
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

func rel(root, path string) string {
	if r, err := filepath.Rel(root, path); err == nil {
		return r
	}
	return path
}
