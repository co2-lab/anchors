//go:build windows

package flow

import (
	"os/exec"
	"syscall"
	"testing"
)

// Not linked to a rule: it runs only on Windows, and no Windows run has proven it yet.
//
// On Windows there is no session: the child is created in a new process group, so the
// console's CTRL_C_EVENT sent to the parent's group does not reach it.
func TestDetach_childGetsANewProcessGroup(t *testing.T) {
	c := exec.Command("cmd", "/c", "exit")
	detach(c)
	if c.SysProcAttr == nil || c.SysProcAttr.CreationFlags&syscall.CREATE_NEW_PROCESS_GROUP == 0 {
		t.Errorf("the child must be created with CREATE_NEW_PROCESS_GROUP: %+v", c.SysProcAttr)
	}
}
