import { defineCollection } from 'astro:content';
import { docsLoader } from '@astrojs/starlight/loaders';
import { docsSchema } from '@astrojs/starlight/schema';

// Toda a documentação é servida sob `/docs`, com i18n: português é o locale
// RAIZ (sem prefixo), inglês é secundário (`/en/docs/...`).
//
// Layout em disco: `src/content/docs/<locale>/<caminho>.md`.
// O Starlight extrai o locale do PRIMEIRO segmento do id gerado, então:
//   pt/conceito.md   -> id `docs/conceito`        (raiz, sem locale)
//   en/concept.md    -> id `en/docs/concept`      (locale primeiro)
const ROOT_LOCALE = 'pt';

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
				const rest = parts.join('/');
				const docs = rest ? `docs/${rest}` : 'docs';
				return locale === ROOT_LOCALE ? docs : `${locale}/${docs}`;
			},
		}),
		schema: docsSchema(),
	}),
};
