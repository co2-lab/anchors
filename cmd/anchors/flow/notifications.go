package flow

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// --- notifications.md: a message to every agent, delivered by `anchors next` ---
//
// Whoever runs the project had no way to TELL the agents something. A comment on a card
// reaches only the agent holding that card, and only if it rereads the card. Measured in
// blue-eyes, 2026-09-24: five agents on another machine kept opening bugs as decisions
// after `anchors escalate --bug` existed, and the only channel to say "update the binary"
// was a comment on each escalation after the fact.
//
// `anchors next` is the one command every agent runs between two pieces of work, so it is
// where the message goes: the content of `notifications.md` is printed on top, before the
// card, every time. Removing the content stops it.
//
// In github mode the file is read from the INTEGRATION BRANCH on the platform, not from the
// agent's checkout: a worktree can be days behind, and a message must reach every agent as
// soon as it is merged. In local mode the file at the project root is the message.

// NotificationsFile is the file at the repository root whose content `anchors next` prints.
const NotificationsFile = "notifications.md"

// ghRaw runs `gh` for the notifications read; a variable so the tests replace it.
var ghRaw = func(args ...string) ([]byte, error) {
	return exec.Command("gh", args...).CombinedOutput()
}

var htmlCommentRE = regexp.MustCompile(`(?s)<!--.*?-->`)

// notificationText is what is worth printing: the content without HTML comments, trimmed.
// The comments let the file explain itself ("write here what every agent must read") and
// still count as empty — an explanation printed on every `next` would teach agents to skip
// the block.
func notificationText(raw string) string {
	return strings.TrimSpace(htmlCommentRE.ReplaceAllString(raw, ""))
}

// notificationsFromBranch reads `notifications.md` from `branch` of `repo` on the platform.
// A missing file is not an error: most of the time there is nothing to say.
func notificationsFromBranch(repo, branch string) (string, error) {
	path := "repos/" + repo + "/contents/" + NotificationsFile
	if branch != "" {
		path += "?ref=" + url.QueryEscape(branch)
	}
	out, err := ghRaw("api", "-H", "Accept: application/vnd.github.raw", path)
	if err != nil {
		if strings.Contains(string(out), "404") || strings.Contains(string(out), "Not Found") {
			return "", nil
		}
		return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// notificationsLocal reads `notifications.md` at the project root (local mode).
func notificationsLocal(root string) string {
	b, err := os.ReadFile(filepath.Join(root, NotificationsFile))
	if err != nil {
		return ""
	}
	return string(b)
}

// printNotifications prints the message on top of `next`. A read failure is ONE line and
// never stops the command: the message is an aid, and an agent without work because the
// notice could not be read would be worse than the notice missing.
func printNotifications(raw, source string, readErr error) {
	if readErr != nil {
		fmt.Printf("· could not read %s (%s): %v\n\n", NotificationsFile, source, readErr)
		return
	}
	text := notificationText(raw)
	if text == "" {
		return
	}
	fmt.Printf("📣 NOTIFICATIONS — %s (%s). Read it, and act on it before your card:\n\n", NotificationsFile, source)
	for _, l := range strings.Split(text, "\n") {
		fmt.Println("  " + l)
	}
	fmt.Printf("\n%s\n\n", strings.Repeat("─", 72))
}
