// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// Anchors — landing (Astro custom em `src/pages`) + docs (Starlight sob `/docs`).
// i18n: português é o DEFAULT (raiz, sem prefixo); inglês é secundário (`/en`).
// https://astro.build/config
export default defineConfig({
	// NÃO declarar `i18n` aqui: o Starlight deriva a config i18n do Astro a
	// partir do seu próprio `locales` (0.41.x proíbe as duas ao mesmo tempo).
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
				root: { label: 'Português', lang: 'pt-BR' },
				en: { label: 'English', lang: 'en' },
			},
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/co2-lab/anchors' },
			],
			// A documentação é servida sob `/docs` (via `generateId` em
			// `src/content.config.ts`), deixando a raiz `/` livre para a landing
			// de marketing em `src/pages/index.astro`.
			sidebar: [
				{
					label: 'Conceito',
					translations: { en: 'Concept' },
					items: [{ slug: 'docs/conceito' }],
				},
				{
					label: 'Os pilares',
					translations: { en: 'The pillars' },
					items: [
						{ label: 'Estrutura de Projeto', translations: { en: 'Project Structure' }, slug: 'docs/estrutura' },
						{ label: 'Planejamento', translations: { en: 'Planning' }, slug: 'docs/planejamento' },
						{ label: 'Spec', translations: { en: 'Spec' }, slug: 'docs/spec' },
						{ label: 'Tipos de Spec', translations: { en: 'Spec Types' }, slug: 'docs/tipos-de-spec' },
						{ label: 'Rastreabilidade', translations: { en: 'Traceability' }, slug: 'docs/rastreabilidade' },
						{ label: 'Propagação', translations: { en: 'Propagation' }, slug: 'docs/propagacao' },
						{ label: 'Qualidade', translations: { en: 'Quality' }, slug: 'docs/qualidade' },
					],
				},
				{
					label: 'O CLI',
					translations: { en: 'The CLI' },
					items: [{ label: 'Visão geral', translations: { en: 'Overview' }, slug: 'docs/cli' }],
				},
			],
		}),
	],
});
