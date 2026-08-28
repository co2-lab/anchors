// Dicionário de strings da landing (página Astro custom, fora do Starlight).
// Português é o idioma DEFAULT (raiz `/`); os demais são secundários (`/en/`, `/es/`, `/ja/`, `/zh/`).
// A landing lê deste dicionário via `t(lang, 'chave')` — uma página, cinco idiomas.

export const languages = {
	pt: 'Português',
	en: 'English',
	es: 'Español',
	ja: '日本語',
	zh: '简体中文',
} as const;

export const defaultLang = 'pt';

export type Lang = keyof typeof languages;

// Estrutura espelhada em todos os idiomas. Chaves em inglês; valores traduzidos.
export const ui = {
	pt: {
		'meta.title': 'Anchors — continuidade para desenvolvimento assistido por IA',
		'meta.description':
			'Anchors é um framework de continuidade: âncoras que guiam o trabalho na ida e o confrontam na volta, e não podem mentir sobre o código.',

		'nav.problem': 'O problema',
		'nav.solution': 'A solução',
		'nav.philosophy': 'Os pilares',
		'nav.numbers': 'O CLI',
		'nav.docs': 'Docs',

		'hero.badge': 'Em desenvolvimento ativo',
		'hero.title.1': 'A sessão termina. O contexto some.',
		'hero.title.2': 'O projeto continua.',
		'hero.sub':
			'Anchors é um framework spec-first de continuidade para desenvolvimento assistido por IA: âncoras que guiam o trabalho na ida, seguram a corda, e confrontam o código na volta — e não podem mentir.',
		'hero.cta.start': 'Ler o conceito',
		'hero.cta.github': 'Ver no GitHub',
		'hero.note.1': 'O nome vem da ',
		'hero.note.strong': 'escalada',
		'hero.note.2': ': pontos fixos que dizem para onde ir e demarcam por onde você passou.',

		'climb.label.climbing': 'Na escalada:',
		'climb.label.anchors': 'No Anchors:',
		'climb.eyebrow': 'A inspiração',
		'climb.title': 'Pegamos emprestada a lógica da escalada.',
		'climb.lead': 'Numa parede, uma âncora faz três coisas — e cada uma resolve um problema real de quem sobe sem poder confiar só na memória. O Anchors aplica a mesma lógica ao desenvolvimento.',
		'climb.1.title': 'Aponta para onde ir',
		'climb.1.climbing': 'Marca o próximo ponto na superfície sem apoio — é por saber para onde ir que você ousa dar o passo seguinte.',
		'climb.1.dev': 'O plano e a spec dizem o que construir antes da primeira linha — a IA (ou você) decide com direção, não no escuro.',
		'climb.2.title': 'Segura a corda',
		'climb.2.climbing': 'Se você escorregar, a âncora vira o safepoint que interrompe a queda — sem ela, um erro custaria a escalada inteira.',
		'climb.2.dev': 'Os gates de qualidade são esse safepoint: pegam a divergência na hora, então um erro custa uma correção — nunca o projeto inteiro.',
		'climb.3.title': 'Demarca por onde você passou',
		'climb.3.climbing': 'A rota fica registrada — quem vier depois não precisa decifrar a parede de novo, só seguir o que já foi validado.',
		'climb.3.dev': 'O grafo de dependências e o histórico de issues fazem o mesmo: a próxima sessão lê o que ficou provado e continua dali — não recomeça do zero.',

		'specfirst.eyebrow': 'A categoria',
		'specfirst.title': 'Anchors é spec-first.',
		'specfirst.lead': 'Spec-first quer dizer que a spec vem antes do código — e continua sendo a verdade depois dele. Não é documentação escrita a posteriori: é a origem de onde tudo mais deriva e contra a qual tudo é confrontado.',
		'specfirst.1.title': 'A spec nasce antes do código',
		'specfirst.1.body': 'Nada se implementa sem uma spec que diga o que construir. A IA (ou você) escreve a partir dela — nunca o contrário.',
		'specfirst.2.title': 'A spec nunca fica desatualizada em silêncio',
		'specfirst.2.body': 'Se o código diverge da spec, o Anchors detecta e confronta — a spec continua sendo a verdade, ou é atualizada explicitamente.',
		'specfirst.3.title': 'A spec é rastreável, não decorativa',
		'specfirst.3.body': 'Cada requisito da spec tem um código de identidade que atravessa feature, teste e código — a rastreabilidade sabe exatamente o que cobre o quê.',

		'problem.eyebrow': 'O problema',
		'problem.title': 'Toda sessão de IA sofre de amnésia estrutural.',
		'problem.1.title': 'A sessão acaba, o contexto some',
		'problem.1.body':
			'Outro agente, outro dia, outro humano — a próxima sessão recomeça sem saber as regras, sem saber o que já foi decidido e por quê.',
		'problem.2.title': 'Código cresce, direção não',
		'problem.2.body':
			'O projeto ganha linhas, mas perde rumo. Ferramentas de agente já resolvem a amnésia dentro de uma sessão — não entre elas.',
		'problem.3.title': 'Testes verdes, requisitos sem prova',
		'problem.3.body':
			'A suíte passa enquanto pedaços inteiros da spec seguem sem cobertura nenhuma — e ninguém percebe até doer.',
		'problem.punch.1': 'O que falta não é memória. É ',
		'problem.punch.strong': 'artefato ancorado no repositório',
		'problem.punch.2': ' — algo que a sessão futura é obrigada a respeitar.',

		'solution.eyebrow': 'A solução',
		'solution.title': 'Âncoras que apontam, seguram e demarcam.',
		'solution.lead':
			'Cada âncora cumpre uma ou mais das três funções da escalada: diz para onde ir, segura a corda para você não cair, e marca por onde você passou. Juntas formam um grafo material, versionado — e cada aresta sabe se está em dia.',
		'flow.1.name': 'Planejamento',
		'flow.1.what': 'semeia',
		'flow.2.name': 'Spec',
		'flow.2.what': 'define a verdade',
		'flow.3.name': 'Rastreabilidade',
		'flow.3.what': 'conecta',
		'flow.4.name': 'Propagação e Qualidade',
		'flow.4.what': 'confrontam',
		'solution.caption.1': 'Quando algo diverge, o Anchors ',
		'solution.caption.2':
			'detecta e apresenta — nunca arbitra sozinho. Ou o trabalho está errado, ou a âncora ficou para trás.',

		'thesis.eyebrow': 'A tese',
		'thesis.title': 'Maturidade não é um número. É o vigor dos pilares.',
		'thesis.lead.pre': 'Um projeto é maduro quando tem ',
		'thesis.lead.1': 'todos os seus pilares',
		'thesis.lead.em': 'implementados e vigorosos',
		'thesis.lead.2': ' — não quando um dashboard isolado diz que sim.',
		'pillar.1.name': 'Estrutura de Projeto',
		'pillar.1.body': 'A forma do repositório, lida de modo determinístico e confirmada com você.',
		'pillar.2.name': 'Planejamento',
		'pillar.2.body': 'O norte: o plano vivo que semeia as specs.',
		'pillar.3.name': 'Spec',
		'pillar.3.body': 'A origem da verdade: a âncora-base, o safepoint do qual tudo pende.',
		'pillar.4.name': 'Rastreabilidade',
		'pillar.4.body': 'A cola: identidade contínua de spec → feature → teste → código, e o mapa entre os arquivos.',
		'pillar.5.name': 'Propagação',
		'pillar.5.body': 'O motor: faz uma alteração percorrer o organismo, marcando o que ficou stale.',
		'pillar.6.name': 'Qualidade',
		'pillar.6.body': 'Os gates: determinísticos onde dá, julgamento por IA onde não dá.',
		'thesis.card.1.title': 'Duas vias, nunca se contradizem',
		'thesis.card.1.body':
			'O plano vivo (âncoras + grafo, a verdade atual) e o histórico (issues resolvidas, laudos, datados e imutáveis).',
		'thesis.card.2.title': 'Opt-out honesto, não silencioso',
		'thesis.card.2.body':
			'Dispensar uma exigência é explícito, registrado e datado. Dispensa a exigência, nunca o registro de que foi dispensada.',
		'thesis.card.3.title': 'Sincronia incremental',
		'thesis.card.3.body':
			'Cada aresta do grafo sabe se está em dia. Uma mudança repropaga só o que ela tocou — nunca o projeto inteiro.',
		'thesis.card.4.title': 'Issue nunca reabre',
		'thesis.card.4.body':
			'Quando um confronto falha e o Anchors não resolve sozinho, registra uma issue imutável. Recorrência é issue nova — o histórico nunca é reescrito.',

		'proof.eyebrow': 'O CLI',
		'proof.title': 'A ferramenta que uma IA opera para exercitar o ciclo.',
		'proof.note':
			'O Anchors não embute IA: é a ferramenta que a IA usa, em qualquer cliente — Claude Code, GPT, Gemini. Ela lê `anchors guide`, aprende o fluxo, e opera com os comandos.',

		'flow2.1.phase': 'Planejar',
		'flow2.1.artifact': 'guia de plano',
		'flow2.1.cmd': 'anchors guide plan',
		'flow2.2.phase': 'Especificar',
		'flow2.2.artifact': '.spec.md',
		'flow2.2.cmd': 'anchors guide spec',
		'flow2.3.phase': 'Mapear',
		'flow2.3.artifact': 'grafo de dependências',
		'flow2.3.cmd': 'anchors map build',
		'flow2.4.phase': 'Implementar',
		'flow2.4.artifact': 'código + feature',
		'flow2.4.cmd': 'anchors guide code',
		'flow2.5.phase': 'Testar',
		'flow2.5.artifact': 'testes',
		'flow2.5.cmd': 'anchors guide test',
		'flow2.6.phase': 'Confrontar',
		'flow2.6.artifact': 'issue se divergir',
		'flow2.6.cmd': 'anchors check / doctor',

		'cta.title': 'Explore a documentação.',
		'cta.body':
			'Do mecanismo comum aos seis pilares: entenda como o Anchors mantém um projeto coerente ao longo do tempo.',
		'cta.docs': 'Ler os docs',
		'cta.github': 'GitHub',

		'footer.tagline':
			'Um framework de continuidade para desenvolvimento assistido por IA — âncoras que não podem mentir sobre o código.',
		'footer.project': 'um projeto',
		'footer.opensource': 'fonte disponível (BSL)',
	},

	en: {
		'meta.title': 'Anchors — continuity for AI-assisted development',
		'meta.description':
			'Anchors is a continuity framework: anchors that guide the work on the way up and confront it on the way back, and cannot lie about the code.',

		'nav.problem': 'The problem',
		'nav.solution': 'The solution',
		'nav.philosophy': 'The pillars',
		'nav.numbers': 'The CLI',
		'nav.docs': 'Docs',

		'hero.badge': 'In active development',
		'hero.title.1': 'The session ends. The context vanishes.',
		'hero.title.2': 'The project goes on.',
		'hero.sub':
			'Anchors is a spec-first continuity framework for AI-assisted development: anchors that guide the work on the way up, hold the rope, and confront the code on the way back — and cannot lie.',
		'hero.cta.start': 'Read the concept',
		'hero.cta.github': 'View on GitHub',
		'hero.note.1': 'The name comes from ',
		'hero.note.strong': 'climbing',
		'hero.note.2': ': fixed points that say where to go and mark where you’ve been.',

		'climb.label.climbing': 'In climbing:',
		'climb.label.anchors': 'In Anchors:',
		'climb.eyebrow': 'The inspiration',
		'climb.title': 'We borrowed the logic of climbing.',
		'climb.lead': 'On a wall, an anchor does three things — and each one solves a real problem for someone who can’t rely on memory alone to get up safely. Anchors applies the same logic to development.',
		'climb.1.title': 'Points where to go',
		'climb.1.climbing': 'It marks the next point on unsupported terrain — knowing where to go is what lets you dare the next move.',
		'climb.1.dev': 'The plan and the spec say what to build before the first line — the AI (or you) decides with direction, not in the dark.',
		'climb.2.title': 'Holds the rope',
		'climb.2.climbing': 'If you slip, the anchor becomes the safepoint that stops the fall — without it, one mistake would cost the whole climb.',
		'climb.2.dev': 'The quality gates are that safepoint: they catch the divergence right away, so a mistake costs one fix — never the whole project.',
		'climb.3.title': 'Marks where you’ve been',
		'climb.3.climbing': 'The route stays recorded — whoever comes next doesn’t have to decode the wall again, just follow what’s already validated.',
		'climb.3.dev': 'The dependency graph and the issue history do the same: the next session reads what’s already proven and continues from there — it doesn’t start from zero.',

		'specfirst.eyebrow': 'The category',
		'specfirst.title': 'Anchors is spec-first.',
		'specfirst.lead': 'Spec-first means the spec comes before the code — and stays the truth after it. It’s not documentation written after the fact: it’s the origin everything else derives from, and the thing everything is confronted against.',
		'specfirst.1.title': 'The spec is born before the code',
		'specfirst.1.body': 'Nothing gets implemented without a spec that says what to build. The AI (or you) writes from it — never the other way around.',
		'specfirst.2.title': 'The spec never goes stale silently',
		'specfirst.2.body': 'If the code diverges from the spec, Anchors detects it and confronts it — the spec stays the truth, or it gets updated explicitly.',
		'specfirst.3.title': 'The spec is traceable, not decorative',
		'specfirst.3.body': 'Every requirement in the spec carries an identity code that runs through feature, test, and code — traceability knows exactly what covers what.',

		'problem.eyebrow': 'The problem',
		'problem.title': 'Every AI session suffers structural amnesia.',
		'problem.1.title': 'The session ends, the context vanishes',
		'problem.1.body':
			'Another agent, another day, another human — the next session starts over without knowing the rules, without knowing what was already decided and why.',
		'problem.2.title': 'Code grows, direction doesn’t',
		'problem.2.body':
			'The project gains lines but loses direction. Agent tooling already handles amnesia within a session — not between sessions.',
		'problem.3.title': 'Green tests, unproven requirements',
		'problem.3.body':
			'The suite passes while entire chunks of the spec stay uncovered — and nobody notices until it hurts.',
		'problem.punch.1': 'What’s missing isn’t memory. It’s ',
		'problem.punch.strong': 'artifact anchored in the repository',
		'problem.punch.2': ' — something the next session is forced to respect.',

		'solution.eyebrow': 'The solution',
		'solution.title': 'Anchors that guide, hold, and mark.',
		'solution.lead':
			'Every anchor fulfills one or more of the three climbing functions: it points where to go, it holds the rope so you don’t fall, and it marks where you’ve been. Together they form a material, versioned graph — and every edge knows if it’s in sync.',
		'flow.1.name': 'Planning',
		'flow.1.what': 'seeds',
		'flow.2.name': 'Spec',
		'flow.2.what': 'defines the truth',
		'flow.3.name': 'Traceability',
		'flow.3.what': 'connects',
		'flow.4.name': 'Propagation & Quality',
		'flow.4.what': 'confront',
		'solution.caption.1': 'When something diverges, Anchors ',
		'solution.caption.2':
			'detects and surfaces it — never arbitrates on its own. Either the work is wrong, or the anchor fell behind.',

		'thesis.eyebrow': 'The thesis',
		'thesis.title': 'Maturity isn’t a number. It’s the pillars’ vigor.',
		'thesis.lead.pre': 'A project is mature when ',
		'thesis.lead.1': 'all of its pillars',
		'thesis.lead.em': 'are implemented and vigorous',
		'thesis.lead.2': ' — not when one isolated dashboard says so.',
		'pillar.1.name': 'Project Structure',
		'pillar.1.body': 'The shape of the repository, read deterministically and confirmed with you.',
		'pillar.2.name': 'Planning',
		'pillar.2.body': 'The north: the live plan that seeds the specs.',
		'pillar.3.name': 'Spec',
		'pillar.3.body': 'The origin of truth: the base anchor, the safepoint everything hangs from.',
		'pillar.4.name': 'Traceability',
		'pillar.4.body': 'The glue: continuous identity from spec → feature → test → code, and the map between files.',
		'pillar.5.name': 'Propagation',
		'pillar.5.body': 'The engine: makes a change ripple through the organism, marking what went stale.',
		'pillar.6.name': 'Quality',
		'pillar.6.body': 'The gates: deterministic where possible, AI judgment where not.',
		'thesis.card.1.title': 'Two tracks, never contradicting',
		'thesis.card.1.body':
			'The live plan (anchors + graph, the current truth) and the history (resolved issues, dated and immutable reports).',
		'thesis.card.2.title': 'Honest opt-out, never silent',
		'thesis.card.2.body':
			'Waiving a requirement is explicit, logged, and dated. It waives the requirement, never the record that it was waived.',
		'thesis.card.3.title': 'Incremental sync',
		'thesis.card.3.body':
			'Every edge in the graph knows if it’s in sync. A change re-propagates only what it touched — never the whole project.',
		'thesis.card.4.title': 'An issue never reopens',
		'thesis.card.4.body':
			'When a confrontation fails and Anchors can’t resolve it alone, it logs an immutable issue. Recurrence is a new issue — history is never rewritten.',

		'proof.eyebrow': 'The CLI',
		'proof.title': 'The tool an AI operates to run the cycle.',
		'proof.note':
			'Anchors doesn’t embed AI: it’s the tool an AI uses, in any client — Claude Code, GPT, Gemini. It reads `anchors guide`, learns the flow, and operates the commands.',

		'flow2.1.phase': 'Plan',
		'flow2.1.artifact': 'plan guide',
		'flow2.1.cmd': 'anchors guide plan',
		'flow2.2.phase': 'Specify',
		'flow2.2.artifact': '.spec.md',
		'flow2.2.cmd': 'anchors guide spec',
		'flow2.3.phase': 'Map',
		'flow2.3.artifact': 'dependency graph',
		'flow2.3.cmd': 'anchors map build',
		'flow2.4.phase': 'Implement',
		'flow2.4.artifact': 'code + feature',
		'flow2.4.cmd': 'anchors guide code',
		'flow2.5.phase': 'Test',
		'flow2.5.artifact': 'tests',
		'flow2.5.cmd': 'anchors guide test',
		'flow2.6.phase': 'Confront',
		'flow2.6.artifact': 'issue if it diverges',
		'flow2.6.cmd': 'anchors check / doctor',

		'cta.title': 'Explore the documentation.',
		'cta.body':
			'From the common mechanism to the six pillars: understand how Anchors keeps a project coherent over time.',
		'cta.docs': 'Read the docs',
		'cta.github': 'GitHub',

		'footer.tagline':
			'A continuity framework for AI-assisted development — anchors that cannot lie about the code.',
		'footer.project': 'a project by',
		'footer.opensource': 'source-available (BSL)',
	},

	es: {
		'meta.title': 'Anchors — continuidad para desarrollo asistido por IA',
		'meta.description':
			'Anchors es un framework de continuidad: anclas que guían el trabajo de ida y lo confrontan a la vuelta, y no pueden mentir sobre el código.',

		'nav.problem': 'El problema',
		'nav.solution': 'La solución',
		'nav.philosophy': 'Los pilares',
		'nav.numbers': 'El CLI',
		'nav.docs': 'Docs',

		'hero.badge': 'En desarrollo activo',
		'hero.title.1': 'La sesión termina. El contexto desaparece.',
		'hero.title.2': 'El proyecto continúa.',
		'hero.sub':
			'Anchors es un framework spec-first de continuidad para desarrollo asistido por IA: anclas que guían el trabajo de ida, sostienen la cuerda, y confrontan el código a la vuelta — y no pueden mentir.',
		'hero.cta.start': 'Leer el concepto',
		'hero.cta.github': 'Ver en GitHub',
		'hero.note.1': 'El nombre viene de la ',
		'hero.note.strong': 'escalada',
		'hero.note.2': ': puntos fijos que dicen hacia dónde ir y marcan por dónde pasaste.',

		'climb.label.climbing': 'En la escalada:',
		'climb.label.anchors': 'En Anchors:',
		'climb.eyebrow': 'La inspiración',
		'climb.title': 'Tomamos prestada la lógica de la escalada.',
		'climb.lead': 'En una pared, un ancla hace tres cosas — y cada una resuelve un problema real de quien sube sin poder confiar solo en la memoria. Anchors aplica la misma lógica al desarrollo.',
		'climb.1.title': 'Indica hacia dónde ir',
		'climb.1.climbing': 'Marca el próximo punto en terreno sin apoyo — saber hacia dónde ir es lo que te permite arriesgar el siguiente paso.',
		'climb.1.dev': 'El plan y la spec dicen qué construir antes de la primera línea — la IA (o tú) decide con dirección, no a ciegas.',
		'climb.2.title': 'Sostiene la cuerda',
		'climb.2.climbing': 'Si resbalas, el ancla se convierte en el safepoint que detiene la caída — sin ella, un error costaría toda la escalada.',
		'climb.2.dev': 'Los gates de calidad son ese safepoint: atrapan la divergencia al instante, así un error cuesta una corrección — nunca el proyecto entero.',
		'climb.3.title': 'Marca por dónde pasaste',
		'climb.3.climbing': 'La ruta queda registrada — quien venga después no tiene que descifrar la pared de nuevo, solo seguir lo ya validado.',
		'climb.3.dev': 'El grafo de dependencias y el historial de issues hacen lo mismo: la próxima sesión lee lo ya probado y continúa desde ahí — no empieza de cero.',

		'specfirst.eyebrow': 'La categoría',
		'specfirst.title': 'Anchors es spec-first.',
		'specfirst.lead': 'Spec-first significa que la spec viene antes del código — y sigue siendo la verdad después de él. No es documentación escrita después de los hechos: es el origen del que todo lo demás deriva, y aquello contra lo que todo se confronta.',
		'specfirst.1.title': 'La spec nace antes que el código',
		'specfirst.1.body': 'Nada se implementa sin una spec que diga qué construir. La IA (o tú) escribe a partir de ella — nunca al revés.',
		'specfirst.2.title': 'La spec nunca queda desactualizada en silencio',
		'specfirst.2.body': 'Si el código diverge de la spec, Anchors lo detecta y lo confronta — la spec sigue siendo la verdad, o se actualiza explícitamente.',
		'specfirst.3.title': 'La spec es rastreable, no decorativa',
		'specfirst.3.body': 'Cada requisito de la spec lleva un código de identidad que atraviesa feature, test y código — la trazabilidad sabe exactamente qué cubre qué.',

		'problem.eyebrow': 'El problema',
		'problem.title': 'Toda sesión de IA sufre amnesia estructural.',
		'problem.1.title': 'La sesión termina, el contexto desaparece',
		'problem.1.body':
			'Otro agente, otro día, otro humano — la próxima sesión empieza de cero sin saber las reglas, sin saber qué se decidió y por qué.',
		'problem.2.title': 'El código crece, la dirección no',
		'problem.2.body':
			'El proyecto gana líneas pero pierde rumbo. Las herramientas de agentes ya resuelven la amnesia dentro de una sesión — no entre sesiones.',
		'problem.3.title': 'Tests en verde, requisitos sin prueba',
		'problem.3.body':
			'La suite pasa mientras partes enteras de la spec siguen sin cobertura — y nadie lo nota hasta que duele.',
		'problem.punch.1': 'Lo que falta no es memoria. Es ',
		'problem.punch.strong': 'artefacto anclado en el repositorio',
		'problem.punch.2': ' — algo que la próxima sesión está obligada a respetar.',

		'solution.eyebrow': 'La solución',
		'solution.title': 'Anclas que guían, sostienen y marcan.',
		'solution.lead':
			'Cada ancla cumple una o más de las tres funciones de la escalada: indica hacia dónde ir, sostiene la cuerda para que no caigas, y marca por dónde pasaste. Juntas forman un grafo material, versionado — y cada arista sabe si está al día.',
		'flow.1.name': 'Planificación',
		'flow.1.what': 'siembra',
		'flow.2.name': 'Spec',
		'flow.2.what': 'define la verdad',
		'flow.3.name': 'Trazabilidad',
		'flow.3.what': 'conecta',
		'flow.4.name': 'Propagación y Calidad',
		'flow.4.what': 'confrontan',
		'solution.caption.1': 'Cuando algo diverge, Anchors lo ',
		'solution.caption.2':
			'detecta y presenta — nunca arbitra solo. O el trabajo está mal, o el ancla quedó atrás.',

		'thesis.eyebrow': 'La tesis',
		'thesis.title': 'La madurez no es un número. Es el vigor de los pilares.',
		'thesis.lead.pre': 'Un proyecto es maduro cuando tiene ',
		'thesis.lead.1': 'todos sus pilares',
		'thesis.lead.em': 'implementados y vigorosos',
		'thesis.lead.2': ' — no cuando un dashboard aislado lo dice.',
		'pillar.1.name': 'Estructura de Proyecto',
		'pillar.1.body': 'La forma del repositorio, leída de modo determinístico y confirmada contigo.',
		'pillar.2.name': 'Planificación',
		'pillar.2.body': 'El norte: el plan vivo que siembra las specs.',
		'pillar.3.name': 'Spec',
		'pillar.3.body': 'El origen de la verdad: el ancla base, el safepoint del cual todo pende.',
		'pillar.4.name': 'Trazabilidad',
		'pillar.4.body': 'La cola: identidad continua de spec → feature → test → código, y el mapa entre los archivos.',
		'pillar.5.name': 'Propagación',
		'pillar.5.body': 'El motor: hace que un cambio recorra el organismo, marcando lo que quedó stale.',
		'pillar.6.name': 'Calidad',
		'pillar.6.body': 'Los gates: determinísticos donde se puede, juicio por IA donde no.',
		'thesis.card.1.title': 'Dos vías, nunca se contradicen',
		'thesis.card.1.body':
			'El plan vivo (anclas + grafo, la verdad actual) y el histórico (issues resueltas, informes fechados e inmutables).',
		'thesis.card.2.title': 'Opt-out honesto, nunca silencioso',
		'thesis.card.2.body':
			'Renunciar a un requisito es explícito, registrado y fechado. Exime del requisito, nunca del registro de que fue eximido.',
		'thesis.card.3.title': 'Sincronía incremental',
		'thesis.card.3.body':
			'Cada arista del grafo sabe si está al día. Un cambio repropaga solo lo que tocó — nunca el proyecto entero.',
		'thesis.card.4.title': 'Una issue nunca reabre',
		'thesis.card.4.body':
			'Cuando una confrontación falla y Anchors no puede resolverla solo, registra una issue inmutable. La recurrencia es una issue nueva — el histórico nunca se reescribe.',

		'proof.eyebrow': 'El CLI',
		'proof.title': 'La herramienta que una IA opera para ejecutar el ciclo.',
		'proof.note':
			'Anchors no integra IA: es la herramienta que una IA usa, en cualquier cliente — Claude Code, GPT, Gemini. Lee `anchors guide`, aprende el flujo, y opera los comandos.',

		'flow2.1.phase': 'Planificar',
		'flow2.1.artifact': 'guía de plan',
		'flow2.1.cmd': 'anchors guide plan',
		'flow2.2.phase': 'Especificar',
		'flow2.2.artifact': '.spec.md',
		'flow2.2.cmd': 'anchors guide spec',
		'flow2.3.phase': 'Mapear',
		'flow2.3.artifact': 'grafo de dependencias',
		'flow2.3.cmd': 'anchors map build',
		'flow2.4.phase': 'Implementar',
		'flow2.4.artifact': 'código + feature',
		'flow2.4.cmd': 'anchors guide code',
		'flow2.5.phase': 'Probar',
		'flow2.5.artifact': 'tests',
		'flow2.5.cmd': 'anchors guide test',
		'flow2.6.phase': 'Confrontar',
		'flow2.6.artifact': 'issue si diverge',
		'flow2.6.cmd': 'anchors check / doctor',

		'cta.title': 'Explora la documentación.',
		'cta.body':
			'Del mecanismo común a los seis pilares: entiende cómo Anchors mantiene un proyecto coherente a lo largo del tiempo.',
		'cta.docs': 'Leer los docs',
		'cta.github': 'GitHub',

		'footer.tagline':
			'Un framework de continuidad para desarrollo asistido por IA — anclas que no pueden mentir sobre el código.',
		'footer.project': 'un proyecto de',
		'footer.opensource': 'fuente disponible (BSL)',
	},

	ja: {
		'meta.title': 'Anchors — AI支援開発のための継続性フレームワーク',
		'meta.description':
			'Anchors は継続性フレームワークです。作業を導き、戻ってきたときにコードと照合するアンカー — それは嘘をつけません。',

		'nav.problem': '課題',
		'nav.solution': '解決策',
		'nav.philosophy': '柱',
		'nav.numbers': 'CLI',
		'nav.docs': 'ドキュメント',

		'hero.badge': '活発に開発中',
		'hero.title.1': 'セッションは終わる。コンテキストは消える。',
		'hero.title.2': 'プロジェクトは続く。',
		'hero.sub':
			'Anchors は spec-first な、AI 支援開発のための継続性フレームワークです。行きは作業を導き、綱を支え、帰りにはコードと照合します — そして嘘をつけません。',
		'hero.cta.start': 'コンセプトを読む',
		'hero.cta.github': 'GitHub で見る',
		'hero.note.1': '名前は',
		'hero.note.strong': 'クライミング',
		'hero.note.2': 'に由来します。行き先を示し、通った道を記す固定点です。',

		'climb.label.climbing': 'クライミングでは：',
		'climb.label.anchors': 'Anchors では：',
		'climb.eyebrow': '着想の源',
		'climb.title': 'クライミングの論理を借りた。',
		'climb.lead': '壁の上で、アンカーは3つのことをします — そのそれぞれが、記憶だけに頼らず安全に登る人が直面する現実の課題を解決します。Anchors は同じ論理を開発に応用します。',
		'climb.1.title': '行き先を示す',
		'climb.1.climbing': '支えのない面での次のポイントを示します — 行き先がわかるからこそ、次の一手に踏み出せます。',
		'climb.1.dev': '計画と spec が最初の1行を書く前に何を作るべきかを示します — AI（またはあなた）は闇雲にではなく、方向を持って判断します。',
		'climb.2.title': '綱を支える',
		'climb.2.climbing': '滑落したとき、アンカーが落下を止めるセーフポイントになります — なければ、一つのミスがクライミング全体を失わせます。',
		'climb.2.dev': '品質ゲートがそのセーフポイントです。逸脱をその場で捕まえるので、ミスの代償は一つの修正で済み、プロジェクト全体には及びません。',
		'climb.3.title': '通った道を記す',
		'climb.3.climbing': 'ルートは記録されます — 後から来る人は壁を一から読み解く必要はなく、すでに検証された道をたどるだけです。',
		'climb.3.dev': '依存グラフと issue の履歴も同じ役割を果たします。次のセッションはすでに証明されたことを読み、そこから続けます — ゼロからは始めません。',

		'specfirst.eyebrow': 'カテゴリー',
		'specfirst.title': 'Anchors は spec-first です。',
		'specfirst.lead': 'spec-first とは、spec がコードより先に存在し、コードが書かれた後も真実であり続けるということです。後付けのドキュメントではありません — すべてがそこから派生し、すべてがそれと照合される起点です。',
		'specfirst.1.title': 'spec はコードより先に生まれる',
		'specfirst.1.body': '何を作るべきかを示す spec なしには何も実装されません。AI（またはあなた）はそこから書きます — 逆はありません。',
		'specfirst.2.title': 'spec は黙って古くならない',
		'specfirst.2.body': 'コードが spec から乖離すれば、Anchors はそれを検出し照合します — spec が真実であり続けるか、明示的に更新されます。',
		'specfirst.3.title': 'spec は追跡可能で、飾りではない',
		'specfirst.3.body': 'spec の各要件は feature・テスト・コードを貫く識別コードを持ちます — トレーサビリティは何が何をカバーするかを正確に把握します。',

		'problem.eyebrow': '課題',
		'problem.title': 'すべての AI セッションは構造的健忘症に苦しむ。',
		'problem.1.title': 'セッションが終わるとコンテキストが消える',
		'problem.1.body':
			'別のエージェント、別の日、別の人間 — 次のセッションはルールを知らず、何がなぜ決まったかも知らずに再出発します。',
		'problem.2.title': 'コードは増えても方向は残らない',
		'problem.2.body':
			'プロジェクトは行数を増やしても方向を失います。エージェントツールはセッション内の健忘症は解決済みですが、セッション間はまだです。',
		'problem.3.title': 'テストは緑、要件は未証明',
		'problem.3.body':
			'スイートは通過するのに、仕様の丸ごとの部分がカバーされないまま — 痛い目に遭うまで誰も気づきません。',
		'problem.punch.1': '足りないのは記憶ではありません。',
		'problem.punch.strong': 'リポジトリに固定された成果物',
		'problem.punch.2': 'です — 次のセッションが従わざるを得ないもの。',

		'solution.eyebrow': '解決策',
		'solution.title': '導き、支え、記すアンカー。',
		'solution.lead':
			'各アンカーはクライミングの3つの機能のうち一つ以上を果たします：行き先を示す、落ちないように綱を支える、通った道を記す。それらが集まって、実体を持つバージョン管理されたグラフを形成します — 各エッジは同期しているかどうかを知っています。',
		'flow.1.name': '計画',
		'flow.1.what': '種をまく',
		'flow.2.name': 'Spec',
		'flow.2.what': '真実を定義する',
		'flow.3.name': 'トレーサビリティ',
		'flow.3.what': '繋ぐ',
		'flow.4.name': '伝播と品質',
		'flow.4.what': '照合する',
		'solution.caption.1': '何かが乖離したとき、Anchors は',
		'solution.caption.2':
			'検出して提示します — 自分では決して裁定しません。作業が間違っているか、アンカーが遅れているかのどちらかです。',

		'thesis.eyebrow': 'テーゼ',
		'thesis.title': '成熟度は数値ではない。柱の活力だ。',
		'thesis.lead.pre': 'プロジェクトが成熟しているのは、',
		'thesis.lead.1': 'すべての柱',
		'thesis.lead.em': 'が実装され活力を持つとき',
		'thesis.lead.2': 'です — 一つの孤立したダッシュボードがそう言うときではありません。',
		'pillar.1.name': 'プロジェクト構造',
		'pillar.1.body': 'リポジトリの形。決定論的に読み取られ、あなたと確認されます。',
		'pillar.2.name': '計画',
		'pillar.2.body': '北：spec の種をまく生きた計画。',
		'pillar.3.name': 'Spec',
		'pillar.3.body': '真実の起源：すべてがぶら下がる基点となるアンカー。',
		'pillar.4.name': 'トレーサビリティ',
		'pillar.4.body': '糊：spec → feature → test → code を貫く継続的な識別性と、ファイル間のマップ。',
		'pillar.5.name': '伝播',
		'pillar.5.body': 'エンジン：変更を組織全体に伝播させ、stale になった箇所を示す。',
		'pillar.6.name': '品質',
		'pillar.6.body': 'ゲート：可能な限り決定論的に、そうでなければ AI の判断で。',
		'thesis.card.1.title': '二つの流れ、決して矛盾しない',
		'thesis.card.1.body':
			'生きた計画（アンカー＋グラフ、現在の真実）と履歴（解決済みの issue、日付付きで不変のレポート）。',
		'thesis.card.2.title': '正直なオプトアウト、沈黙ではない',
		'thesis.card.2.body':
			'要件の免除は明示的で、記録され、日付が付きます。免除するのは要件だけで、免除された記録は残ります。',
		'thesis.card.3.title': '増分的な同期',
		'thesis.card.3.body':
			'グラフの各エッジは同期しているかを知っています。変更は触れた部分だけを再伝播します — プロジェクト全体ではありません。',
		'thesis.card.4.title': 'issue は再オープンしない',
		'thesis.card.4.body':
			'照合が失敗し Anchors が単独で解決できないとき、不変の issue を記録します。再発は新しい issue です — 履歴は書き換えられません。',

		'proof.eyebrow': 'CLI',
		'proof.title': 'AI がサイクルを動かすために操作するツール。',
		'proof.note':
			'Anchors は AI を内蔵しません。Claude Code、GPT、Gemini など、どのクライアントでも AI が使うツールです。`anchors guide` を読み、フローを学び、コマンドを操作します。',

		'flow2.1.phase': '計画',
		'flow2.1.artifact': 'プランガイド',
		'flow2.1.cmd': 'anchors guide plan',
		'flow2.2.phase': '仕様化',
		'flow2.2.artifact': '.spec.md',
		'flow2.2.cmd': 'anchors guide spec',
		'flow2.3.phase': 'マッピング',
		'flow2.3.artifact': '依存グラフ',
		'flow2.3.cmd': 'anchors map build',
		'flow2.4.phase': '実装',
		'flow2.4.artifact': 'コード＋feature',
		'flow2.4.cmd': 'anchors guide code',
		'flow2.5.phase': 'テスト',
		'flow2.5.artifact': 'テスト',
		'flow2.5.cmd': 'anchors guide test',
		'flow2.6.phase': '照合',
		'flow2.6.artifact': '乖離時は issue',
		'flow2.6.cmd': 'anchors check / doctor',

		'cta.title': 'ドキュメントを見る。',
		'cta.body':
			'共通のメカニズムから六つの柱まで — Anchors がどのようにプロジェクトを時間とともに一貫させるかを理解しましょう。',
		'cta.docs': 'ドキュメントを読む',
		'cta.github': 'GitHub',

		'footer.tagline':
			'AI 支援開発のための継続性フレームワーク — コードについて嘘をつけないアンカー。',
		'footer.project': 'プロジェクト提供：',
		'footer.opensource': 'ソース公開（BSL）',
	},

	zh: {
		'meta.title': 'Anchors — 面向 AI 辅助开发的连续性框架',
		'meta.description':
			'Anchors 是一个连续性框架：锚点在去程引导工作，在回程与代码对照 —— 它们不能说谎。',

		'nav.problem': '问题',
		'nav.solution': '解决方案',
		'nav.philosophy': '支柱',
		'nav.numbers': 'CLI',
		'nav.docs': '文档',

		'hero.badge': '积极开发中',
		'hero.title.1': '会话结束，上下文消失。',
		'hero.title.2': '项目仍在继续。',
		'hero.sub':
			'Anchors 是一个 spec-first（规范优先）的、面向 AI 辅助开发的连续性框架：锚点在去程引导工作、拉住绳索，在回程与代码对照 —— 它们不能说谎。',
		'hero.cta.start': '阅读概念',
		'hero.cta.github': '在 GitHub 查看',
		'hero.note.1': '这个名字来自',
		'hero.note.strong': '攀登',
		'hero.note.2': '：固定点指明方向，并标记你走过的路。',

		'climb.label.climbing': '在攀登中：',
		'climb.label.anchors': '在 Anchors 中：',
		'climb.eyebrow': '灵感来源',
		'climb.title': '我们借用了攀登的逻辑。',
		'climb.lead': '在岩壁上，一个锚点做三件事 —— 每一件都解决了一个不能仅凭记忆安全攀登的人所面临的真实问题。Anchors 将同样的逻辑应用到开发中。',
		'climb.1.title': '指明方向',
		'climb.1.climbing': '标记出无支撑地形上的下一个点 —— 正因为知道该往哪走，你才敢迈出下一步。',
		'climb.1.dev': '计划和 spec 在写下第一行之前就说明要构建什么 —— AI（或你）是带着方向做决定，而不是在黑暗中摸索。',
		'climb.2.title': '拉住绳索',
		'climb.2.climbing': '如果你滑落，锚点就成了阻止坠落的安全点 —— 没有它，一次失误就会葬送整次攀登。',
		'climb.2.dev': '质量关卡就是这个安全点：它当场捕获偏差，让一次失误的代价只是一次修正 —— 而不是整个项目。',
		'climb.3.title': '标记走过的路',
		'climb.3.climbing': '路线被记录下来 —— 后来者不必重新解读岩壁，只需沿着已验证的路径前进。',
		'climb.3.dev': '依赖图和 issue 历史做的是同一件事：下一次会话读取已经证实的内容，从那里继续 —— 而不是从零开始。',

		'specfirst.eyebrow': '类别定位',
		'specfirst.title': 'Anchors 是 spec-first 的。',
		'specfirst.lead': 'spec-first 意味着 spec 先于代码存在 —— 并且在代码写出之后依然是真相。它不是事后补写的文档：它是一切的起点，也是一切被对照检验的基准。',
		'specfirst.1.title': 'spec 先于代码诞生',
		'specfirst.1.body': '没有说明要构建什么的 spec，就不会有任何实现。AI（或你）依据它来编写 —— 而不是反过来。',
		'specfirst.2.title': 'spec 不会悄悄过时',
		'specfirst.2.body': '如果代码偏离了 spec，Anchors 会检测并对照 —— spec 要么继续作为真相，要么被明确更新。',
		'specfirst.3.title': 'spec 是可追溯的，不是摆设',
		'specfirst.3.body': 'spec 中的每个需求都带有贯穿 feature、测试和代码的标识码 —— 可追溯性精确掌握谁覆盖了什么。',

		'problem.eyebrow': '问题',
		'problem.title': '每一次 AI 会话都患有结构性失忆。',
		'problem.1.title': '会话结束，上下文消失',
		'problem.1.body':
			'另一个代理、另一天、另一个人 —— 下一次会话从零开始，不知道规则，不知道已经决定了什么、为什么。',
		'problem.2.title': '代码在增长，方向没有',
		'problem.2.body':
			'项目行数增加，方向却在流失。代理工具已经解决了单次会话内的失忆问题 —— 但会话之间还没有。',
		'problem.3.title': '测试全绿，需求却未被证明',
		'problem.3.body':
			'测试套件通过了，而规范中的整块内容仍未被覆盖 —— 直到出问题才有人注意到。',
		'problem.punch.1': '缺少的不是记忆，而是',
		'problem.punch.strong': '固定在仓库中的产物',
		'problem.punch.2': '—— 下一次会话必须遵守的东西。',

		'solution.eyebrow': '解决方案',
		'solution.title': '引导、支撑、标记的锚点。',
		'solution.lead':
			'每个锚点履行攀登三种功能中的一项或多项：指明方向、拉住绳索防止跌落、标记走过的路。它们共同构成一个具体的、受版本控制的图 —— 每条边都知道自己是否同步。',
		'flow.1.name': '规划',
		'flow.1.what': '播种',
		'flow.2.name': 'Spec',
		'flow.2.what': '定义真相',
		'flow.3.name': '可追溯性',
		'flow.3.what': '连接',
		'flow.4.name': '传播与质量',
		'flow.4.what': '对照',
		'solution.caption.1': '当出现偏差时，Anchors 会',
		'solution.caption.2':
			'检测并呈现 —— 从不独自裁定。要么是工作有误，要么是锚点落后了。',

		'thesis.eyebrow': '论点',
		'thesis.title': '成熟度不是一个数字，而是支柱的活力。',
		'thesis.lead.pre': '一个项目成熟，是因为它的',
		'thesis.lead.1': '所有支柱',
		'thesis.lead.em': '都已实现且充满活力',
		'thesis.lead.2': '—— 而不是因为某个孤立的仪表盘这么说。',
		'pillar.1.name': '项目结构',
		'pillar.1.body': '仓库的形态，以确定性方式读取，并与你确认。',
		'pillar.2.name': '规划',
		'pillar.2.body': '方向：播种规范的活计划。',
		'pillar.3.name': 'Spec',
		'pillar.3.body': '真相的起源：一切依附其上的基础锚点。',
		'pillar.4.name': '可追溯性',
		'pillar.4.body': '粘合剂：从 spec → feature → test → code 的连续标识，以及文件之间的映射。',
		'pillar.5.name': '传播',
		'pillar.5.body': '引擎：让一次变更传遍整个系统，标记出已过期（stale）的部分。',
		'pillar.6.name': '质量',
		'pillar.6.body': '关卡：能确定性的地方就确定性判断，不能的地方交给 AI 判断。',
		'thesis.card.1.title': '两条轨道，从不矛盾',
		'thesis.card.1.body':
			'活的计划（锚点 + 图，当前的真相）与历史（已解决的 issue，带日期且不可变的报告）。',
		'thesis.card.2.title': '诚实的豁免，从不悄悄发生',
		'thesis.card.2.body':
			'豁免一项要求是显式的、被记录的、带日期的。它豁免的是要求本身，而不是"曾被豁免"这一记录。',
		'thesis.card.3.title': '增量同步',
		'thesis.card.3.body':
			'图中的每条边都知道自己是否同步。一次变更只重新传播它触及的部分 —— 从不是整个项目。',
		'thesis.card.4.title': 'issue 永不重新打开',
		'thesis.card.4.body':
			'当一次对照失败、Anchors 无法独自解决时，它会记录一个不可变的 issue。再次发生就是新的 issue —— 历史从不被改写。',

		'proof.eyebrow': 'CLI',
		'proof.title': 'AI 用来运行这个循环的工具。',
		'proof.note':
			'Anchors 不内置 AI：它是 AI 在任何客户端中使用的工具 —— Claude Code、GPT、Gemini。它读取 `anchors guide`，学习流程，并操作命令。',

		'flow2.1.phase': '规划',
		'flow2.1.artifact': '计划指南',
		'flow2.1.cmd': 'anchors guide plan',
		'flow2.2.phase': '定义规范',
		'flow2.2.artifact': '.spec.md',
		'flow2.2.cmd': 'anchors guide spec',
		'flow2.3.phase': '映射',
		'flow2.3.artifact': '依赖图',
		'flow2.3.cmd': 'anchors map build',
		'flow2.4.phase': '实现',
		'flow2.4.artifact': '代码 + feature',
		'flow2.4.cmd': 'anchors guide code',
		'flow2.5.phase': '测试',
		'flow2.5.artifact': '测试',
		'flow2.5.cmd': 'anchors guide test',
		'flow2.6.phase': '对照',
		'flow2.6.artifact': '偏差则生成 issue',
		'flow2.6.cmd': 'anchors check / doctor',

		'cta.title': '浏览文档。',
		'cta.body':
			'从共同机制到六大支柱：理解 Anchors 如何让项目长期保持一致。',
		'cta.docs': '阅读文档',
		'cta.github': 'GitHub',

		'footer.tagline':
			'面向 AI 辅助开发的连续性框架 —— 不能对代码说谎的锚点。',
		'footer.project': '一个项目，来自',
		'footer.opensource': '源码可见（BSL）',
	},
} as const;

export function useTranslations(lang: Lang) {
	return function t(key: keyof (typeof ui)['pt']): string {
		return (ui[lang] as Record<string, string>)[key] ?? (ui.pt as Record<string, string>)[key];
	};
}
