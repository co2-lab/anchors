package gate

import "strings"

// --- a REGRA dentro do gate ---
//
// Um gate não verifica uma coisa: `spec-completa` cobra duas (não haver placeholder, e
// haver ao menos uma regra catalogada), `header-conforme` cobra várias. Até aqui o
// veredito dizia apenas QUAL GATE reprovou, e isso tem duas consequências:
//
//   - quem dispensa um caso deliberado (a spec que nasce antes da feature) só podia
//     dispensar o gate INTEIRO, e junto ia toda verificação que ele fazia bem;
//   - dois defeitos diferentes no mesmo gate são indistinguíveis num relatório — o
//     leitor precisa da mensagem para saber qual dos dois aconteceu, e mensagem muda.
//
// A regra tem ID: `<gate>/<regra>` — `spec-completa/sem-placeholder`. É o mesmo
// princípio que a doutrina já aplica ao artefato (`WRKSP-B01`): um identificador estável
// é o que permite falar de uma decisão sem descrevê-la de novo.

// RegraID é o identificador estável de uma verificação dentro de um gate.
//
// O formato é `<gate>/<regra>`, e o `<gate>` faz parte de propósito: nomes de regra
// curtos ("sem-placeholder", "tem-codigo") se repetiriam entre gates, e um ID que colide
// não identifica nada.
type RegraID string

// Gate devolve a parte de gate do ID — o que vem antes da barra.
func (r RegraID) Gate() string {
	g, _, _ := strings.Cut(string(r), "/")
	return g
}

// Regra devolve a parte de regra do ID. Vazio quando o gate não declarou regras (o
// veredito é do gate como um todo), e aí o ID é só o nome do gate.
func (r RegraID) Regra() string {
	_, reg, _ := strings.Cut(string(r), "/")
	return reg
}

// NovaRegraID monta o ID. Um `regra` vazio devolve só o gate: é o caso de um gate que
// faz uma verificação só, e para o qual dividir em regras não acrescentaria nada.
func NovaRegraID(gate, regra string) RegraID {
	if regra == "" {
		return RegraID(gate)
	}
	return RegraID(gate + "/" + regra)
}

// Dispensa é o conjunto de regras que o usuário pediu para não confrontar nesta execução,
// com o motivo de cada dispensa.
//
// O motivo é parte do dado, e não um comentário à parte: uma dispensa sem justificativa
// escrita é indistinguível de alguém fugindo de um gate que achou defeito. Com o motivo,
// o relatório diz o que foi pulado E por quê — e quem ler depois não precisa adivinhar.
type Dispensa struct {
	// PorRegra mapeia o ID (`spec-completa/sem-placeholder`) ou o nome do gate inteiro
	// (`trinca-completa`) para o motivo.
	PorRegra map[string]string
	// Alvos restringe a dispensa a caminhos específicos, por ID de regra.
	//
	// Sem isso, dispensar `trinca-completa` para commitar 4 specs novas apagava o gate
	// para o REPOSITÓRIO INTEIRO — e uma trinca que quebrou por descuido noutro lugar
	// passava junto, sem que nada acusasse. O mascaramento que a dispensa por regra
	// existe para evitar, um nível acima.
	//
	// A forma é `regra@alvo=motivo`, onde o alvo é o CÓDIGO do artefato (`WRKSP`) ou o
	// caminho dele. Um ID sem `@` continua valendo para tudo: é a saída grossa, e há
	// casos legítimos (um gate recém-declarado num projeto que ainda não o cumpre em
	// lugar nenhum).
	Alvos map[string][]string
}

// Dispensou diz se esta regra foi dispensada, e devolve o motivo.
//
// Aceita as duas granularidades: dispensar `spec-completa` cobre todas as regras dele,
// e dispensar `spec-completa/sem-placeholder` cobre só aquela. A primeira é a saída
// grossa para quem não conhece as regras; a segunda é a que preserva o resto do gate.
func (d Dispensa) Dispensou(id RegraID) (string, bool) {
	motivo, ok := d.dispensouRegra(id)
	if !ok {
		return "", false
	}
	// Uma dispensa COM alvos não vale para o gate inteiro: ela precisa ser confrontada
	// alvo a alvo, e quem pergunta sem dizer o alvo recebe "não dispensado". É o que
	// impede o filtro de gates de remover o gate da lista e apagá-lo para todo o
	// repositório.
	if len(d.alvosDe(id)) > 0 {
		return "", false
	}
	return motivo, true
}

// DispensouAlvo diz se ESTA regra está dispensada para ESTE alvo.
//
// O alvo é identificado por CÓDIGO (`WRKSP`) ou por CAMINHO
// (`packages/shared/Workspace.spec.md`), e as duas formas valem.
//
// O código é a forma preferida, e a razão é a estabilidade: ele é a identidade do
// artefato e sobrevive a mover ou renomear o arquivo. Uma dispensa presa ao caminho
// deixa de valer silenciosamente quando alguém reorganiza pastas — e o commit seguinte
// reprova sem que nada explique o que mudou.
//
// O caminho continua aceito porque nem todo alvo tem código: um `package.json`, um
// arquivo de configuração, um teste sem identidade declarada.
func (d Dispensa) DispensouAlvo(id RegraID, alvo, codigo string) (string, bool) {
	motivo, ok := d.dispensouRegra(id)
	if !ok {
		return "", false
	}
	alvos := d.alvosDe(id)
	if len(alvos) == 0 {
		return motivo, true
	}
	for _, a := range alvos {
		if a == alvo || (codigo != "" && a == codigo) {
			return motivo, true
		}
	}
	return "", false
}

// dispensouRegra procura o motivo pelo ID exato e, depois, pelo gate inteiro.
func (d Dispensa) dispensouRegra(id RegraID) (string, bool) {
	if len(d.PorRegra) == 0 {
		return "", false
	}
	if motivo, ok := d.PorRegra[string(id)]; ok {
		return motivo, true
	}
	if motivo, ok := d.PorRegra[id.Gate()]; ok {
		return motivo, true
	}
	return "", false
}

// alvosDe devolve os caminhos a que a dispensa desta regra está restrita.
func (d Dispensa) alvosDe(id RegraID) []string {
	if len(d.Alvos) == 0 {
		return nil
	}
	if a, ok := d.Alvos[string(id)]; ok {
		return a
	}
	return d.Alvos[id.Gate()]
}

// ParseDispensa lê a forma textual `regra=motivo,regra=motivo` — o que chega por flag ou
// por variável de ambiente.
//
// Uma entrada sem `=` é recusada, e não aceita em silêncio com motivo vazio: o motivo é
// a única coisa que separa dispensa deliberada de gate ignorado, e aceitá-lo ausente
// esvaziaria a garantia.
func ParseDispensa(bruto string) (Dispensa, []string) {
	d := Dispensa{PorRegra: map[string]string{}, Alvos: map[string][]string{}}
	var erros []string
	for _, parte := range strings.Split(bruto, ",") {
		parte = strings.TrimSpace(parte)
		if parte == "" {
			continue
		}
		regra, motivo, ok := strings.Cut(parte, "=")
		regra, motivo = strings.TrimSpace(regra), strings.TrimSpace(motivo)
		// `regra@caminho` restringe a dispensa àquele alvo. O `@` vem antes do `=`, então
		// o corte acima já separou o motivo — aqui sobra separar a regra do caminho.
		var alvo string
		if r, a, temAlvo := strings.Cut(regra, "@"); temAlvo {
			regra, alvo = strings.TrimSpace(r), strings.TrimSpace(a)
			if alvo == "" {
				erros = append(erros, "`"+parte+"` — falta o caminho depois do `@`")
				continue
			}
		}
		if !ok || motivo == "" {
			erros = append(erros, "`"+parte+"` — falta o motivo (use `regra=por quê`)")
			continue
		}
		if regra == "" {
			erros = append(erros, "`"+parte+"` — falta a regra")
			continue
		}
		d.PorRegra[regra] = motivo
		if alvo != "" {
			d.Alvos[regra] = append(d.Alvos[regra], alvo)
		}
	}
	return d, erros
}
