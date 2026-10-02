import { defineCollection } from 'astro:content';
import { docsLoader } from '@astrojs/starlight/loaders';
import { docsSchema } from '@astrojs/starlight/schema';

// Documentation is served under `/docs` (English is the root/default locale, without prefix).
// Portuguese is secondary, served under `/pt/docs/...`.
//
// Disk layout: `src/content/docs/<locale>/<path>.md`.
// Starlight extracts the locale from the first segment of the generated ID:
//   en/concept.md    -> id `docs/concept`         (root default, no locale prefix)
//   pt/conceito.md   -> id `pt/docs/concept`      (secondary locale with prefix)
const ROOT_LOCALE = 'en';

const ptToCanonical: Record<string, string> = {
	conceito: 'concept',
	congelar: 'freeze',
	estrutura: 'structure',
	'fluxo-de-trabalho': 'workflow',
	planejamento: 'planning',
	propagacao: 'propagation',
	qualidade: 'quality',
	rastreabilidade: 'traceability',
	'tipos-de-spec': 'spec-types',
	camadas: 'layers',
	'conceitos/ancora': 'concepts/anchor',
	'conceitos/unidade': 'concepts/unit',
	'conceitos/grafo-e-mapa': 'concepts/graph-and-map',
	'conceitos/camadas-e-regimes': 'concepts/layers-and-regimes',
	'conceitos/gates-e-vereditos': 'concepts/gates-and-verdicts',
	'conceitos/rastreabilidade-e-codigos': 'concepts/traceability-and-codes',
	'conceitos/propagacao-e-impacto': 'concepts/propagation-and-impact',
	'conceitos/doutrina-de-produto': 'concepts/doctrine',
	'conceitos/feature-flags': 'concepts/feature-flags',
	'conceitos/julgamento-ia': 'concepts/ai-judgment',
	'conceitos/maturidade-e-saude': 'concepts/maturity-and-health',
	'cli/instalacao': 'cli/installation',
	'cli/comandos': 'cli/commands',
	'cli/comandos/check': 'cli/commands/check',
	'cli/comandos/doctor': 'cli/commands/doctor',
	'cli/comandos/map': 'cli/commands/map',
	'cli/comandos/init': 'cli/commands/init',
	'cli/comandos/flow': 'cli/commands/flow',
	'cli/comandos/freeze': 'cli/commands/freeze',
};

export const collections = {
	docs: defineCollection({
		loader: docsLoader({
			generateId: ({ entry }) => {
				const slug = entry
					.replace(/\.(md|mdx)$/, '')
					.replace(/\/?index$/, '')
					.replace(/^\/+/, '');
				const parts = slug.split('/').filter(Boolean);
				const locale = parts.shift() ?? ROOT_LOCALE;
				let rest = parts.join('/');
				if (locale === 'pt' && ptToCanonical[rest]) {
					rest = ptToCanonical[rest];
				}
				const docs = rest ? `docs/${rest}` : 'docs';
				return locale === ROOT_LOCALE ? docs : `${locale}/${docs}`;
			},
		}),
		schema: docsSchema(),
	}),
};
