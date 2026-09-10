package config

import "testing"

// O PADRÃO é avisar, e isso é uma decisão, não um acaso: o pipeline desatualizado ainda faz
// o trabalho antigo, e derrubá-lo troca "faz menos do que devia" por "não faz nada".
//
// O teste existe porque a inversão desse default é silenciosa — quem trocasse `false` por
// `true` aqui quebraria o CI de todo projeto que só queria o aviso, e nada acusaria.
func TestPipelineVelhoSoBarraQuandoOProjetoPede(t *testing.T) {
	var semDeclarar Workflow
	if semDeclarar.PipelineVelhoBarra() {
		t.Error("sem declarar, o padrão é AVISAR — barrar o CI de quem não pediu é regressão")
	}
	if (&Workflow{StalePipelineBlocks: true}).PipelineVelhoBarra() != true {
		t.Error("quem declarou `stale_pipeline_blocks: true` quer barrar")
	}
	// Nil não pode estourar: o `Workflow` é opcional no anchors.yaml, e o `doctor` chama
	// isto em projeto que nunca declarou a seção.
	var nulo *Workflow
	if nulo.PipelineVelhoBarra() {
		t.Error("sem a seção `workflow:`, o padrão continua sendo avisar")
	}
}
