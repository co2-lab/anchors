// Package daemon gerencia o watcher em background: daemonização (re-exec
// desanexado), PID file, log, e o flag de pausa. Mantém o terminal livre — o
// `start` retorna o controle imediatamente.
package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// StateDir é onde o estado do daemon vive (dentro do projeto).
const StateDir = ".anchors"

type Paths struct {
	Dir    string
	PID    string
	Log    string
	Paused string
	Meta   string
}

func PathsFor(root string) Paths {
	d := filepath.Join(root, StateDir)
	return Paths{
		Dir:    d,
		PID:    filepath.Join(d, "watch.pid"),
		Log:    filepath.Join(d, "watch.log"),
		Paused: filepath.Join(d, "watch.paused"),
		Meta:   filepath.Join(d, "watch.meta"),
	}
}

// Running devolve o PID do daemon vivo, ou 0 se não há nenhum. Limpa PID stale
// (processo morto que deixou o arquivo para trás).
func Running(p Paths) int {
	data, err := os.ReadFile(p.PID)
	if err != nil {
		return 0
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		return 0
	}
	if !alive(pid) {
		_ = os.Remove(p.PID) // stale
		return 0
	}
	return pid
}

// WritePID grava o PID do daemon e a meta (início).
func WritePID(p Paths, pid int) error {
	if err := os.MkdirAll(p.Dir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(p.PID, []byte(strconv.Itoa(pid)), 0o644); err != nil {
		return err
	}
	return nil
}

// WriteMeta registra desde quando (o daemon não pode usar time.Now via args do
// script de workflow, mas aqui é processo normal — pode).
func WriteMeta(p Paths, started time.Time, root string) error {
	line := fmt.Sprintf("started=%s\nroot=%s\n", started.Format(time.RFC3339), root)
	return os.WriteFile(p.Meta, []byte(line), 0o644)
}

func ReadMeta(p Paths) string {
	data, _ := os.ReadFile(p.Meta)
	return string(data)
}

// Stop encerra o daemon e limpa o PID (o mecanismo de encerramento é
// específico de plataforma — ver daemon_unix.go / daemon_windows.go).
func Stop(p Paths) error {
	pid := Running(p)
	if pid == 0 {
		return fmt.Errorf("watcher não está rodando")
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := terminate(proc); err != nil {
		return err
	}
	_ = os.Remove(p.PID)
	return nil
}

// Pause/Resume/IsPaused via arquivo-flag (simples, inspecionável). O loop do
// watcher checa IsPaused a cada evento.
func Pause(p Paths) error  { return os.WriteFile(p.Paused, []byte("1"), 0o644) }
func Resume(p Paths) error { return os.Remove(p.Paused) }
func IsPaused(p Paths) bool {
	_, err := os.Stat(p.Paused)
	return err == nil
}

// Cleanup remove os artefatos de estado (ao encerrar o loop).
func Cleanup(p Paths) {
	_ = os.Remove(p.PID)
	_ = os.Remove(p.Paused)
}
