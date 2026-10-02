---
title: "Instalação do CLI Anchors"
description: "Como instalar o CLI do Anchors no macOS, Linux, Windows, ambientes de CI/CD e Docker."
---

O **CLI do Anchors** (`anchors`) é distribuído como um binário estático único sem nenhuma dependência externa de runtime. Ele roda perfeitamente na sua máquina local, dentro de sessões de IA (como Claude Code, Cursor, Windsurf ou agentes de terminal) e nas esteiras de integração contínua (CI/CD).

---

## 1. Instalação Rápida (Recomendada)

### macOS & Linux via Homebrew
A forma mais simples e automatizada de instalar e manter o Anchors atualizado no macOS e Linux é pelo Homebrew:

```bash
# Adiciona o tap oficial do CO2 Lab
brew tap co2-lab/tap

# Instala o Anchors
brew install anchors

# Confere a instalação
anchors version
```

### Script Automatizado de Instalação (macOS & Linux)
Se você não usa Homebrew, pode instalar diretamente o binário mais recente com nosso script seguro:

```bash
curl -fsSL https://raw.githubusercontent.com/co2-lab/anchors/main/install.sh | bash
```

O script detecta automaticamente o sistema operacional e a arquitetura (`amd64` ou `arm64`), baixa a release verificada do GitHub, confere os checksums e move o binário para `/usr/local/bin` (ou `~/.local/bin`).

---

## 2. Instalação via Toolchain do Go

Se você já possui o Go (1.23+) instalado:

```bash
go install github.com/co2-lab/anchors/cmd/anchors@latest
```

Certifique-se de que o diretório `$GOPATH/bin` (geralmente `~/go/bin`) esteja no seu `PATH`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

---

## 3. Download Direto do Binário (GitHub Releases)

Você pode baixar os binários pré-compilados diretamente da página de [Releases no GitHub](https://github.com/co2-lab/anchors/releases):

| Sistema Operacional | Arquitetura | Pacote do Binário |
| :--- | :--- | :--- |
| **macOS** | Apple Silicon (`M1/M2/M3/M4`) | `anchors_darwin_arm64.tar.gz` |
| **macOS** | Intel (`x86_64`) | `anchors_darwin_amd64.tar.gz` |
| **Linux** | 64-bit x86 (`amd64`) | `anchors_linux_amd64.tar.gz` |
| **Linux** | ARM64 (`aarch64`) | `anchors_linux_arm64.tar.gz` |
| **Windows** | 64-bit (`x86_64`) | `anchors_windows_amd64.zip` |

### Passos de Instalação Manual:
```bash
# Exemplo para Linux AMD64
curl -LO https://github.com/co2-lab/anchors/releases/latest/download/anchors_linux_amd64.tar.gz
tar -xzf anchors_linux_amd64.tar.gz
sudo mv anchors /usr/local/bin/
sudo chmod +x /usr/local/bin/anchors
```

---

## 4. Instalação no Windows

### Via PowerShell (Download Direto):
Abra o PowerShell como Administrador:
```powershell
Invoke-WebRequest -Uri "https://github.com/co2-lab/anchors/releases/latest/download/anchors_windows_amd64.zip" -OutFile "$env:TEMPnchors.zip"
Expand-Archive -Path "$env:TEMPnchors.zip" -DestinationPath "$env:ProgramFiles\Anchors" -Force
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:ProgramFiles\Anchors", "Machine")
```

### Via Scoop:
```powershell
scoop bucket add co2-lab https://github.com/co2-lab/scoop-bucket.git
scoop install anchors
```

---

## 5. Ambientes Docker e CI/CD

### Imagem Docker
O Anchors está disponível no GitHub Packages como um container minimalista:

```bash
docker pull ghcr.io/co2-lab/anchors:latest

# Executar na pasta atual do projeto
docker run --rm -v $(pwd):/workspace -w /workspace ghcr.io/co2-lab/anchors:latest check
```

### GitHub Actions
Para rodar os gates do Anchors nas suas PRs do GitHub Actions:

```yaml
name: Anchors Quality Gates

on: [push, pull_request]

jobs:
  anchors-gate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: co2-lab/setup-anchors@v1
        with:
          version: 'latest'
      - name: Rodar Pipeline do Anchors
        run: anchors check
```

### GitLab CI
```yaml
anchors_check:
  image: ghcr.io/co2-lab/anchors:latest
  stage: test
  script:
    - anchors check
```

---

## 6. Autocompletar no Shell

Ative o preenchimento automático de comandos e flags via tecla `Tab`:

### Zsh (macOS / Linux)
```bash
echo 'source <(anchors completion zsh)' >> ~/.zshrc
```

### Bash
```bash
echo 'source <(anchors completion bash)' >> ~/.bashrc
```

### Fish
```bash
anchors completion fish > ~/.config/fish/completions/anchors.fish
```

---

## 7. Verificando a Instalação

Após instalar, execute:

```bash
# Imprime versão instalada e dados de compilação
anchors version

# Roda diagnóstico do ecossistema do projeto
anchors doctor
```

Tudo pronto! Siga para [O fluxo de trabalho](/pt/docs/workflow/) ou explore a [Referência de Comandos](/pt/docs/cli/commands/)!
