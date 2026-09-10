//go:build windows

package daemon

import "os"

// alive: o processo existe? No Windows, FindProcess já abre um handle real
// (ao contrário do Unix, onde sempre "acha" o PID) — a existência do handle
// já é o suficiente.
func alive(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}

// terminate mata o processo diretamente — Windows não tem SIGTERM; Process.Kill
// chama TerminateProcess. O watcher não recebe chance de encerrar o loop de
// forma graciosa nesta plataforma (ver watch.go).
func terminate(proc *os.Process) error {
	return proc.Kill()
}
