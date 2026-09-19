package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/cmd/anchors/governance"
	"github.com/spf13/cobra"
)

func TestGuiaSoCitaComandoQueExiste(t *testing.T) {
	registrados := map[string]bool{}
	var colhe func(*cobra.Command)
	colhe = func(c *cobra.Command) {
		registrados[c.Name()] = true
		for _, f := range c.Commands() {
			colhe(f)
		}
	}
	colhe(newRootCmd())

	re := regexp.MustCompile(`anchors ([a-z][a-z-]*)`)
	vistos := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(governance.WorkGuide, -1) {
		nome := m[1]
		if vistos[nome] {
			continue
		}
		vistos[nome] = true
		if nome == "guide" || registrados[nome] {
			continue
		}
		t.Errorf("o guia manda rodar `anchors %s`, e esse comando não existe — quem o "+
			"seguir recebe `unknown command` e para", nome)
	}
	if len(vistos) == 0 {
		t.Fatal("nenhum `anchors <comando>` no guia — o regex quebrou e o teste passaria vazio")
	}
}

func TestPipelineSoEnsinaComandoQueExiste(t *testing.T) {
	registrados := map[string]bool{}
	var colhe func(*cobra.Command)
	colhe = func(c *cobra.Command) {
		registrados[c.Name()] = true
		for _, f := range c.Commands() {
			colhe(f)
		}
	}
	colhe(newRootCmd())

	dir := filepath.Join("..", "..", "internal", "initx", "workflows")
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("nao consegui ler os workflows: %v", err)
	}

	re := regexp.MustCompile(`anchors ([a-z][a-z-]*)`)
	total := 0
	for _, e := range entradas {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("%s: %v", e.Name(), err)
		}
		vistos := map[string]bool{}
		for _, m := range re.FindAllStringSubmatch(string(b), -1) {
			nome := m[1]
			if vistos[nome] {
				continue
			}
			vistos[nome] = true
			total++
			if registrados[nome] {
				continue
			}
			t.Errorf("%s cita `anchors %s`, e esse comando nao existe -- quem seguir "+
				"a instrucao recebe `unknown command`", e.Name(), nome)
		}
	}
	if total == 0 {
		t.Fatal("nenhum `anchors <comando>` nos workflows -- o regex quebrou e o teste passaria vazio")
	}
}
