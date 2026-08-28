package gate

// --- maturação: o gate informativo que já está limpo (QUALITY §7) ---
//
// Um gate nasce informativo num projeto que ainda não cumpre o limiar dele. O problema é
// o que acontece DEPOIS: quando o projeto passa a cumprir, ninguém volta ao anchors.yaml
// para promovê-lo — e o gate fica informativo para sempre, medindo sem defender.
//
// A maturação é uma decisão humana, e continua sendo: o Anchors não promove sozinho. Mas
// ele pode LEMBRAR, e o lembrete só é útil se aparecer onde a pessoa já está olhando —
// no `check`, no `status`, no `next`. Um aviso que exige rodar um comando específico para
// ser visto é um aviso que ninguém vê.

// Promovivel é um gate informativo cujo veredito está limpo: ele mede algo que o projeto
// já cumpre, e promovê-lo a bloqueante passaria a DEFENDER isso sem custo nenhum hoje.
type Promovivel struct {
	// Gate é o nome, como declarado no anchors.yaml.
	Gate string
	// Passou é quantos nós ele aprovou. Zero significa que ele não teve o que medir —
	// e aí não há o que promover: ver `GatesPromoviveis`.
	Passou int
}

// GatesPromoviveis devolve os gates informativos que estão limpos no perfil dado.
//
// Três condições, e cada uma existe para não sugerir promoção enganosa:
//
//   - informativo — o bloqueante já defende, não há o que sugerir;
//   - zero reprovações — promover um gate que reprova barraria o trabalho na hora;
//   - ao menos uma aprovação — um gate que nunca teve o que medir não está "limpo", está
//     sem dado. Promovê-lo daria a impressão de defesa que não existe, e é exatamente o
//     silêncio que o Anchors combate em outros lugares.
func GatesPromoviveis(p Profile) []Promovivel {
	var out []Promovivel
	for _, nome := range p.GateNames() {
		s := p.ByGate[nome]
		if s.Blocking || s.Fail > 0 || s.Pass == 0 {
			continue
		}
		out = append(out, Promovivel{Gate: nome, Passou: s.Pass})
	}
	return out
}
