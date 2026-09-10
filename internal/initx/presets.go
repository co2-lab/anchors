package initx

import "github.com/co2-lab/anchors/internal/config"

// Presets de estrutura por stack. Destilados de um estudo das estruturas de projeto
// consagradas por linguagem/framework (docs oficiais + convenção de comunidade). O
// init OFERECE um preset conforme a stack detectada/escolhida; as perguntas seguem
// SEMPRE sendo feitas (o preset só pré-preenche os defaults — como a inferência).
//
// Cada preset traz as LAYERS de código no formato do anchors.yaml. As layers de
// ARTEFATO (spec/feature/test/guide/plan) NÃO entram aqui — vêm da escolha de
// artefatos (ApplyArtifactChoice), que é ortogonal à stack. A layer de teste do
// preset é a exceção: a convenção de nome de teste é específica da stack.

// Preset é uma estrutura de projeto consagrada, pronta para virar Layers.
type Preset struct {
	Name        string // id curto: "go", "node-ts", "python-web", ...
	Title       string // rótulo humano para o menu do init
	Pattern     string // padrão arquitetural: "layered", "clean", "feature-sliced", ...
	Modular     bool   // organiza por feature/módulo vertical?
	ModuleGlob  string // se modular: glob do diretório de módulo (ex: "src/features/*")
	TestPattern string // convenção de teste (ex: "**/*_test.go", "**/*.test.ts")
	Colocated   bool   // teste/spec ao lado do código?
	Layers      []PresetLayer
	// CoverageHint: como este stack emite os artefatos que `anchors ingest` consome —
	// o comando de teste com cobertura + os formatos (JUnit para execução, lcov para
	// linha). Reduz o atrito de adoção: o init mostra isto, o projeto não descobre só.
	CoverageHint string
}

// PresetLayer é uma camada de código de um preset (espelha config.Layer, sem os
// campos que o init preenche depois).
type PresetLayer struct {
	Name       string
	Pattern    string
	Kind       string // quase sempre "code"; a camada de teste é "test"
	Tags       []string
	CodePrefix string // Camada 1 do código de identidade, quando a camada é um módulo
}

// ToLayers converte as PresetLayers em config.Layer prontas para o anchors.yaml.
func (p Preset) ToLayers() map[string]config.Layer {
	out := map[string]config.Layer{}
	for _, l := range p.Layers {
		kind := l.Kind
		if kind == "" {
			kind = "code"
		}
		out[l.Name] = config.Layer{
			Pattern:    l.Pattern,
			Kind:       kind,
			Tags:       l.Tags,
			CodePrefix: l.CodePrefix,
		}
	}
	return out
}

// Presets é o catálogo — destilado do estudo de estruturas consagradas por stack
// (docs oficiais + convenção de comunidade, pesquisado 2026). Cada preset é a
// estrutura MAIS consagrada da stack. Quando há duas concorrentes, escolhi a
// dominante e anotei; o usuário sempre pode editar o anchors.yaml depois.
//
// `code_prefix` fica VAZIO aqui de propósito: nos presets modulares, o prefixo é
// DEDUZIDO pelo init a partir do nome de cada módulo real (via code.ModulePrefix) —
// não dá para hardcodar porque os módulos são do projeto, não da stack.
var Presets = []Preset{
	// ---------- Frontend / JS ----------
	{
		Name: "node-ts", Title: "Node.js/TypeScript backend (Nest — modular)",
		Pattern: "package-by-feature", Modular: true, ModuleGlob: "src/modules/*",
		TestPattern: "src/**/*.spec.ts", Colocated: true,
		CoverageHint: "jest --coverage --coverageReporters=lcov --reporters=default --reporters=jest-junit → coverage/lcov.info + junit.xml",
		Layers: []PresetLayer{
			{Name: "modules", Pattern: "src/modules/**/*.ts", Kind: "code", Tags: []string{"backend", "module"}},
			{Name: "core", Pattern: "src/core/**/*.ts", Kind: "code", Tags: []string{"backend"}},
			{Name: "common", Pattern: "src/common/**/*.ts", Kind: "code", Tags: []string{"backend"}},
			{Name: "unit-test", Pattern: "src/**/*.spec.ts", Kind: "test", Tags: []string{"backend"}},
			{Name: "e2e-test", Pattern: "test/**/*.e2e-spec.ts", Kind: "test", Tags: []string{"backend", "e2e"}},
		},
	},
	{
		Name: "express-ts", Title: "Node.js/TypeScript backend (Express/Fastify — layered)",
		Pattern: "layered", Modular: false,
		TestPattern: "**/*.test.ts", Colocated: false,
		Layers: []PresetLayer{
			{Name: "routes", Pattern: "src/{routes,controllers}/**/*.ts", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "services", Pattern: "src/services/**/*.ts", Kind: "code", Tags: []string{"backend"}},
			{Name: "repositories", Pattern: "src/repositories/**/*.ts", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "models", Pattern: "src/models/**/*.ts", Kind: "code", Tags: []string{"backend", "domain"}},
			{Name: "test", Pattern: "**/*.test.ts", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "nextjs", Title: "Next.js (App Router)",
		Pattern: "file-based + horizontal", Modular: false,
		TestPattern: "**/*.test.{ts,tsx}", Colocated: true,
		CoverageHint: "jest/vitest --coverage (lcov) + reporter junit; e2e Playwright tem --reporter=junit",
		Layers: []PresetLayer{
			{Name: "routes", Pattern: "{src/,}app/**/*.{tsx,ts}", Kind: "code", Tags: []string{"frontend", "web"}},
			{Name: "components", Pattern: "{src/,}components/**/*.{tsx,ts}", Kind: "code", Tags: []string{"frontend"}},
			{Name: "lib", Pattern: "{src/,}lib/**/*.ts", Kind: "code", Tags: []string{"frontend"}},
			{Name: "hooks", Pattern: "{src/,}hooks/**/*.ts", Kind: "code", Tags: []string{"frontend"}},
			{Name: "test", Pattern: "**/*.{test,spec}.{ts,tsx}", Kind: "test", Tags: []string{"frontend"}},
		},
	},
	{
		Name: "angular", Title: "Angular (feature-based)",
		Pattern: "feature-based", Modular: true, ModuleGlob: "src/app/features/*",
		TestPattern: "src/**/*.spec.ts", Colocated: true,
		Layers: []PresetLayer{
			{Name: "features", Pattern: "src/app/features/**/*.ts", Kind: "code", Tags: []string{"frontend", "module"}},
			{Name: "core", Pattern: "src/app/core/**/*.ts", Kind: "code", Tags: []string{"frontend"}},
			{Name: "shared", Pattern: "src/app/shared/**/*.ts", Kind: "code", Tags: []string{"frontend"}},
			{Name: "test", Pattern: "src/**/*.spec.ts", Kind: "test", Tags: []string{"frontend"}},
		},
	},
	{
		Name: "nuxt", Title: "Nuxt (Vue — directory convention)",
		Pattern: "directory-convention", Modular: false,
		TestPattern: "**/*.{test,spec}.ts", Colocated: false,
		Layers: []PresetLayer{
			{Name: "pages", Pattern: "{app/,}pages/**/*.vue", Kind: "code", Tags: []string{"frontend", "web"}},
			{Name: "components", Pattern: "{app/,}components/**/*.vue", Kind: "code", Tags: []string{"frontend"}},
			{Name: "composables", Pattern: "{app/,}composables/**/*.ts", Kind: "code", Tags: []string{"frontend"}},
			{Name: "layouts", Pattern: "{app/,}layouts/**/*.vue", Kind: "code", Tags: []string{"frontend"}},
			{Name: "server", Pattern: "server/**/*.ts", Kind: "code", Tags: []string{"backend"}},
			{Name: "stores", Pattern: "{app/,}stores/**/*.ts", Kind: "code", Tags: []string{"frontend"}},
			{Name: "test", Pattern: "**/*.{test,spec}.ts", Kind: "test", Tags: []string{"frontend"}},
		},
	},
	{
		Name: "expo-rn", Title: "React Native / Expo (Expo Router)",
		Pattern: "file-based + horizontal", Modular: false,
		TestPattern: "src/**/*.test.{ts,tsx}", Colocated: true,
		CoverageHint: "jest --coverage --coverageReporters=lcov + jest-junit → coverage/lcov.info + junit.xml",
		Layers: []PresetLayer{
			{Name: "screens", Pattern: "src/app/**/*.tsx", Kind: "code", Tags: []string{"mobile", "frontend"}},
			{Name: "components", Pattern: "src/components/**/*.tsx", Kind: "code", Tags: []string{"mobile", "frontend"}},
			{Name: "hooks", Pattern: "src/hooks/**/*.ts", Kind: "code", Tags: []string{"mobile"}},
			{Name: "test", Pattern: "src/**/*.test.{ts,tsx}", Kind: "test", Tags: []string{"mobile"}},
		},
	},

	// ---------- JVM / .NET / PHP ----------
	{
		Name: "spring", Title: "Java/Kotlin + Spring Boot (layered)",
		Pattern: "layered", Modular: false,
		TestPattern: "src/test/**/*.{java,kt}", Colocated: false,
		CoverageHint: "mvn test gera surefire XML (JUnit) em target/surefire-reports; jacoco:report → jacoco.xml (converta p/ lcov ou use cobertura de linha do jacoco)",
		Layers: []PresetLayer{
			{Name: "controllers", Pattern: "src/main/{java,kotlin}/**/controller/**/*.{java,kt}", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "services", Pattern: "src/main/{java,kotlin}/**/service/**/*.{java,kt}", Kind: "code", Tags: []string{"backend"}},
			{Name: "repositories", Pattern: "src/main/{java,kotlin}/**/repository/**/*.{java,kt}", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "domain", Pattern: "src/main/{java,kotlin}/**/{model,domain,entity}/**/*.{java,kt}", Kind: "code", Tags: []string{"backend", "domain"}},
			{Name: "test", Pattern: "src/test/{java,kotlin}/**/*.{java,kt}", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "dotnet-clean", Title: ".NET/C# (Clean Architecture)",
		Pattern: "clean", Modular: false,
		TestPattern: "tests/**/*.cs", Colocated: false,
		Layers: []PresetLayer{
			{Name: "api", Pattern: "src/**/*.Api/**/*.cs", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "application", Pattern: "src/**/*.Application/**/*.cs", Kind: "code", Tags: []string{"backend"}},
			{Name: "domain", Pattern: "src/**/*.Domain/**/*.cs", Kind: "code", Tags: []string{"backend", "domain"}},
			{Name: "infrastructure", Pattern: "src/**/*.Infrastructure/**/*.cs", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "test", Pattern: "tests/**/*.cs", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "laravel", Title: "PHP / Laravel (MVC)",
		Pattern: "mvc", Modular: false,
		TestPattern: "tests/{Feature,Unit}/**/*.php", Colocated: false,
		Layers: []PresetLayer{
			{Name: "controllers", Pattern: "app/Http/Controllers/**/*.php", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "models", Pattern: "app/Models/**/*.php", Kind: "code", Tags: []string{"backend", "domain", "data"}},
			{Name: "requests", Pattern: "app/Http/Requests/**/*.php", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "routes", Pattern: "routes/**/*.php", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "migrations", Pattern: "database/migrations/**/*.php", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "views", Pattern: "resources/views/**/*.blade.php", Kind: "code", Tags: []string{"frontend", "web"}},
			{Name: "test", Pattern: "tests/{Feature,Unit}/**/*.php", Kind: "test", Tags: []string{"backend"}},
		},
	},

	// ---------- Ruby / Elixir / Flutter / C++ ----------
	{
		Name: "rails", Title: "Ruby on Rails (MVC/convention)",
		Pattern: "mvc", Modular: false,
		TestPattern: "{test,spec}/**/*_{test,spec}.rb", Colocated: false,
		Layers: []PresetLayer{
			{Name: "models", Pattern: "app/models/**/*.rb", Kind: "code", Tags: []string{"backend", "domain"}},
			{Name: "controllers", Pattern: "app/controllers/**/*.rb", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "views", Pattern: "app/views/**/*.{erb,haml,slim}", Kind: "code", Tags: []string{"frontend", "web"}},
			{Name: "jobs", Pattern: "app/jobs/**/*.rb", Kind: "code", Tags: []string{"backend"}},
			{Name: "migrations", Pattern: "db/migrate/**/*.rb", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "test", Pattern: "{test,spec}/**/*_{test,spec}.rb", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "phoenix", Title: "Elixir / Phoenix (contexts)",
		Pattern: "contexts", Modular: true, ModuleGlob: "lib/*/*/",
		TestPattern: "test/**/*_test.exs", Colocated: false,
		Layers: []PresetLayer{
			{Name: "contexts", Pattern: "lib/*/[!_]*/**/*.ex", Kind: "code", Tags: []string{"backend", "domain", "module"}},
			{Name: "web-controllers", Pattern: "lib/*_web/controllers/**/*.ex", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "web-live", Pattern: "lib/*_web/live/**/*.ex", Kind: "code", Tags: []string{"frontend", "web"}},
			{Name: "migrations", Pattern: "priv/repo/migrations/**/*.exs", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "test", Pattern: "test/**/*_test.exs", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "flutter", Title: "Flutter / Dart (feature-first)",
		Pattern: "feature-first", Modular: true, ModuleGlob: "lib/**/features/*/",
		TestPattern: "test/**/*_test.dart", Colocated: false,
		Layers: []PresetLayer{
			{Name: "features", Pattern: "lib/**/features/*/**/*.dart", Kind: "code", Tags: []string{"mobile", "frontend", "module"}},
			{Name: "shared", Pattern: "lib/**/{common,shared,core,utils}/**/*.dart", Kind: "code", Tags: []string{"mobile"}},
			{Name: "test", Pattern: "test/**/*_test.dart", Kind: "test", Tags: []string{"mobile"}},
			{Name: "integration-test", Pattern: "integration_test/**/*_test.dart", Kind: "test", Tags: []string{"mobile", "e2e"}},
		},
	},
	{
		Name: "cpp", Title: "C++ (include/ + src/ split)",
		Pattern: "include-src-split", Modular: false,
		TestPattern: "{tests,test}/**/*.{cpp,cc,cxx}", Colocated: false,
		Layers: []PresetLayer{
			{Name: "public-headers", Pattern: "include/**/*.{h,hpp,hxx}", Kind: "code", Tags: []string{"api", "public"}},
			{Name: "sources", Pattern: "src/**/*.{cpp,cc,cxx}", Kind: "code", Tags: []string{"impl"}},
			{Name: "test", Pattern: "{tests,test}/**/*.{cpp,cc,cxx}", Kind: "test", Tags: []string{"test"}},
		},
	},

	// ---------- Python / Go / Rust ----------
	{
		Name: "python-lib", Title: "Python (src-layout / biblioteca — PyPA)",
		Pattern: "src-layout", Modular: false,
		TestPattern: "tests/**/*.py", Colocated: false,
		Layers: []PresetLayer{
			{Name: "src", Pattern: "src/**/*.py", Kind: "code", Tags: []string{"lib"}},
			{Name: "test", Pattern: "tests/**/*.py", Kind: "test", Tags: []string{"lib"}},
		},
	},
	{
		Name: "django", Title: "Python / Django (apps por domínio)",
		Pattern: "app-based", Modular: true, ModuleGlob: "*/",
		TestPattern: "*/tests/**/*.py", Colocated: false,
		CoverageHint: "pytest --cov --cov-report=lcov --junitxml=junit.xml (pytest-cov)",
		Layers: []PresetLayer{
			{Name: "settings", Pattern: "config/settings/*.py", Kind: "code", Tags: []string{"backend", "config"}},
			{Name: "models", Pattern: "*/models.py", Kind: "code", Tags: []string{"backend", "domain", "module"}},
			{Name: "views", Pattern: "*/views.py", Kind: "code", Tags: []string{"backend", "web", "module"}},
			{Name: "urls", Pattern: "*/urls.py", Kind: "code", Tags: []string{"backend", "web", "module"}},
			{Name: "migrations", Pattern: "*/migrations/*.py", Kind: "code", Tags: []string{"backend", "data"}},
			{Name: "templates", Pattern: "**/templates/**/*.html", Kind: "code", Tags: []string{"frontend"}},
			{Name: "test", Pattern: "*/tests/**/*.py", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "fastapi", Title: "Python / FastAPI (app + routers)",
		Pattern: "layered", Modular: false,
		TestPattern: "tests/**/*.py", Colocated: false,
		CoverageHint: "pytest --cov --cov-report=lcov --junitxml=junit.xml (pytest-cov)",
		Layers: []PresetLayer{
			{Name: "routers", Pattern: "app/routers/*.py", Kind: "code", Tags: []string{"backend", "web"}},
			{Name: "internal", Pattern: "app/internal/**/*.py", Kind: "code", Tags: []string{"backend"}},
			{Name: "entrypoint", Pattern: "app/main.py", Kind: "code", Tags: []string{"backend"}},
			{Name: "test", Pattern: "tests/**/*.py", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "go", Title: "Go (standard project layout)",
		Pattern: "standard-go-layout", Modular: true, ModuleGlob: "internal/*",
		TestPattern: "**/*_test.go", Colocated: true,
		CoverageHint: "go test ./... -coverprofile=cover.out; gcov2lcov < cover.out > lcov.info; go test -json | go-junit-report > junit.xml",
		Layers: []PresetLayer{
			{Name: "cmd", Pattern: "cmd/*/*.go", Kind: "code", Tags: []string{"cli", "entrypoint"}},
			{Name: "internal", Pattern: "internal/**/*.go", Kind: "code", Tags: []string{"backend", "module"}},
			{Name: "pkg", Pattern: "pkg/**/*.go", Kind: "code", Tags: []string{"lib", "public"}},
			{Name: "test", Pattern: "**/*_test.go", Kind: "test", Tags: []string{"backend"}},
		},
	},
	{
		Name: "rust", Title: "Rust (cargo layout)",
		Pattern: "cargo-layout", Modular: false,
		TestPattern: "tests/*.rs", Colocated: true,
		CoverageHint: "cargo test -- -Z unstable-options --format junit > junit.xml; cargo llvm-cov --lcov --output-path lcov.info",
		Layers: []PresetLayer{
			{Name: "src", Pattern: "src/**/*.rs", Kind: "code", Tags: []string{"lib"}},
			{Name: "bins", Pattern: "src/bin/*.rs", Kind: "code", Tags: []string{"bin"}},
			{Name: "integ-test", Pattern: "tests/*.rs", Kind: "test", Tags: []string{"integration"}},
		},
	},
}

func PresetByName(name string) (Preset, bool) {
	for _, p := range Presets {
		if p.Name == name {
			return p, true
		}
	}
	return Preset{}, false
}

func PresetNames() []string {
	out := make([]string, len(Presets))
	for i, p := range Presets {
		out[i] = p.Name
	}
	return out
}
