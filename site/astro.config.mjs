// @ts-check
import { defineConfig } from 'astro/config';
import cliCommands from './src/generated/cli-commands.json' with { type: 'json' };
import gateGroups from './src/generated/gates.json' with { type: 'json' };
import starlight from '@astrojs/starlight';

// Anchors — landing (custom Astro in `src/pages`) + docs (Starlight under `/docs`).
// i18n: English is DEFAULT (root, without prefix); Portuguese is secondary (`/pt/docs/...`).
// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Anchors',
			favicon: '/favicon.svg',
			logo: {
				src: './public/anchors-mark-red.svg',
				alt: 'Anchors',
			},
			customCss: ['./src/styles/starlight-theme.css'],
			components: {
				Header: './src/components/overrides/Header.astro',
			},
			defaultLocale: 'root',
			locales: {
				root: { label: 'English', lang: 'en' },
				pt: { label: 'Português', lang: 'pt-BR' },
			},
			head: [
				{
					tag: 'script',
					content: "if (typeof localStorage !== 'undefined' && !localStorage.getItem('starlight-theme')) { localStorage.setItem('starlight-theme', 'light'); document.documentElement.dataset.theme = 'light'; }",
				},
			],
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/co2-lab/anchors' },
			],
			sidebar: [
				{
					label: 'The CLI',
					translations: { 'pt-BR': 'O CLI', pt: 'O CLI' },
					items: [
						{ label: 'Overview', translations: { 'pt-BR': 'Visão Geral', pt: 'Visão Geral' }, slug: 'docs/cli' },
						{ label: 'Installation', translations: { 'pt-BR': 'Instalação', pt: 'Instalação' }, slug: 'docs/cli/installation' },
						{ label: 'Commands Reference', translations: { 'pt-BR': 'Referência de Comandos', pt: 'Referência de Comandos' }, slug: 'docs/cli/commands' },
						{
							label: 'Command Guides',
							translations: { 'pt-BR': 'Guias de Comandos', pt: 'Guias de Comandos' },
							// One page per command, generated from the CLI (`go run ./tools/sitecli`).
							items: cliCommands,
						},
					],
				},
				{
					label: 'Core Concepts',
					translations: { 'pt-BR': 'Conceitos Fundamentais', pt: 'Conceitos Fundamentais' },
					items: [
						{ label: 'Overview', translations: { 'pt-BR': 'Visão Geral', pt: 'Visão Geral' }, slug: 'docs/concept' },
						{ label: 'What is an Anchor', translations: { 'pt-BR': 'O que é uma Âncora', pt: 'O que é uma Âncora' }, slug: 'docs/concepts/anchor' },
						{ label: 'The Unit', translations: { 'pt-BR': 'A Unidade', pt: 'A Unidade' }, slug: 'docs/concepts/unit' },
						{ label: 'The Graph & Map', translations: { 'pt-BR': 'O Grafo e o Mapa', pt: 'O Grafo e o Mapa' }, slug: 'docs/concepts/graph-and-map' },
						{ label: 'Layers & Regimes', translations: { 'pt-BR': 'Camadas e Regimes', pt: 'Camadas e Regimes' }, slug: 'docs/concepts/layers-and-regimes' },
						{ label: 'Gates & Verdicts', translations: { 'pt-BR': 'Gates e Vereditos', pt: 'Gates e Vereditos' }, slug: 'docs/concepts/gates-and-verdicts' },
						{ label: 'Traceability & Codes', translations: { 'pt-BR': 'Rastreabilidade e Códigos', pt: 'Rastreabilidade e Códigos' }, slug: 'docs/concepts/traceability-and-codes' },
						{ label: 'Propagation & Impact', translations: { 'pt-BR': 'Propagação e Impacto', pt: 'Propagação e Impacto' }, slug: 'docs/concepts/propagation-and-impact' },
						{ label: 'Product Doctrine', translations: { 'pt-BR': 'Doutrina de Produto', pt: 'Doutrina de Produto' }, slug: 'docs/concepts/doctrine' },
						{ label: 'Feature Flags', translations: { 'pt-BR': 'Feature Flags', pt: 'Feature Flags' }, slug: 'docs/concepts/feature-flags' },
						{ label: 'AI Judgment', translations: { 'pt-BR': 'Julgamento por IA', pt: 'Julgamento por IA' }, slug: 'docs/concepts/ai-judgment' },
						{ label: 'Maturity & Health', translations: { 'pt-BR': 'Maturidade e Saúde', pt: 'Maturidade e Saúde' }, slug: 'docs/concepts/maturity-and-health' },
					],
				},
				{
					label: 'Project Layers',
					translations: { 'pt-BR': 'Camadas do Projeto', pt: 'Camadas do Projeto' },
					items: [
						{ label: 'Layers Guide', translations: { 'pt-BR': 'Guia de Camadas', pt: 'Guia de Camadas' }, slug: 'docs/layers' },
					],
				},
				{
					label: 'The Pillars',
					translations: { 'pt-BR': 'Os pilares', pt: 'Os pilares' },
					items: [
						{ label: 'Project Structure', translations: { 'pt-BR': 'Estrutura de Projeto', pt: 'Estrutura de Projeto' }, slug: 'docs/structure' },
						{ label: 'Planning', translations: { 'pt-BR': 'Planejamento', pt: 'Planejamento' }, slug: 'docs/planning' },
						{ label: 'Spec', translations: { 'pt-BR': 'Spec', pt: 'Spec' }, slug: 'docs/spec' },
						{ label: 'Spec Types', translations: { 'pt-BR': 'Tipos de Spec', pt: 'Tipos de Spec' }, slug: 'docs/spec-types' },
						{ label: 'Traceability', translations: { 'pt-BR': 'Rastreabilidade', pt: 'Rastreabilidade' }, slug: 'docs/traceability' },
						{ label: 'Propagation', translations: { 'pt-BR': 'Propagação', pt: 'Propagação' }, slug: 'docs/propagation' },
						{ label: 'Quality', translations: { 'pt-BR': 'Qualidade', pt: 'Qualidade' }, slug: 'docs/quality' },
					],
				},
				{
					label: 'Using Anchors',
					translations: { 'pt-BR': 'Usar o Anchors', pt: 'Usar o Anchors' },
					items: [
						{ label: 'The Workflow', translations: { 'pt-BR': 'O fluxo de trabalho', pt: 'O fluxo de trabalho' }, slug: 'docs/workflow' },
						{ label: 'The anchors.yaml', translations: { 'pt-BR': 'O anchors.yaml', pt: 'O anchors.yaml' }, slug: 'docs/anchors-yaml' },
						{ label: 'Freezing the Project', translations: { 'pt-BR': 'Congelar o projeto', pt: 'Congelar o projeto' }, slug: 'docs/freeze' },
					],
				},
				{
					label: 'Gates Catalog',
					translations: { 'pt-BR': 'Catálogo de Gates', pt: 'Catálogo de Gates' },
					items: [
						{ label: 'Gates Index', translations: { 'pt-BR': 'Índice Geral de Gates', pt: 'Índice Geral de Gates' }, slug: 'docs/gates' },
						// One group per section of the catalog, generated from the code (`go run ./tools/sitecli`).
						...gateGroups,
					],
				},
			],
		}),
	],
});
