---
title: "Installing the Anchors CLI"
description: "How to install the Anchors CLI on macOS, Linux, Windows, CI/CD pipelines, and Docker."
---

The **Anchors CLI** (`anchors`) is distributed as a single static binary with zero external runtime dependencies. It runs seamlessly on local developer machines, inside AI pair programming environments (such as Claude Code, Cursor, Windsurf, or terminal agents), and within continuous integration (CI/CD) pipelines.

---

## 1. Quick Install (Recommended)

### macOS & Linux via Homebrew
The simplest and most automated way to install and maintain Anchors on macOS and Linux is through Homebrew:

```bash
# Add the official CO2 Lab tap
brew tap co2-lab/tap

# Install Anchors
brew install anchors

# Verify installation
anchors version
```

### Automated Shell Installer (macOS & Linux)
If you do not use Homebrew, you can install the latest release directly via our secure installer script:

```bash
curl -fsSL https://raw.githubusercontent.com/co2-lab/anchors/main/install.sh | bash
```

The installer detects your operating system and architecture (`amd64` or `arm64`), downloads the latest verified GitHub release binary, verifies checksums, and places it in `/usr/local/bin` (or `~/.local/bin`).

---

## 2. Install via Go Toolchain

If you have Go (1.23+) installed on your machine:

```bash
go install github.com/co2-lab/anchors/cmd/anchors@latest
```

Ensure that `$GOPATH/bin` (typically `~/go/bin`) is in your system's `PATH`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

---

## 3. Direct Binary Download (GitHub Releases)

You can download pre-compiled, signed standalone binaries directly from [GitHub Releases](https://github.com/co2-lab/anchors/releases):

| Operating System | Architecture | Binary Package |
| :--- | :--- | :--- |
| **macOS** | Apple Silicon (`M1/M2/M3/M4`) | `anchors_darwin_arm64.tar.gz` |
| **macOS** | Intel (`x86_64`) | `anchors_darwin_amd64.tar.gz` |
| **Linux** | 64-bit x86 (`amd64`) | `anchors_linux_amd64.tar.gz` |
| **Linux** | ARM64 (`aarch64`) | `anchors_linux_arm64.tar.gz` |
| **Windows** | 64-bit (`x86_64`) | `anchors_windows_amd64.zip` |

### Manual Installation Steps:
```bash
# Example for Linux AMD64
curl -LO https://github.com/co2-lab/anchors/releases/latest/download/anchors_linux_amd64.tar.gz
tar -xzf anchors_linux_amd64.tar.gz
sudo mv anchors /usr/local/bin/
sudo chmod +x /usr/local/bin/anchors
```

---

## 4. Windows Installation

### Via PowerShell (Direct Download):
Open PowerShell as Administrator:
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

## 5. Docker & CI/CD Environments

### Docker Image
Anchors is available on GitHub Packages as a minimal, scratch-based container:

```bash
docker pull ghcr.io/co2-lab/anchors:latest

# Run against current directory
docker run --rm -v $(pwd):/workspace -w /workspace ghcr.io/co2-lab/anchors:latest check
```

### GitHub Actions
To run Anchors in your GitHub Actions CI workflows, add the official setup action:

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
      - name: Run Anchors Pipeline
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

## 6. Shell Autocompletions

Enable tab autocompletions for all commands and flags in your shell:

### Zsh (macOS / Linux)
```bash
anchors completion zsh > "${fpath[1]}/_anchors"
# Or append to your .zshrc:
echo 'source <(anchors completion zsh)' >> ~/.zshrc
```

### Bash
```bash
anchors completion bash > /etc/bash_completion.d/anchors
# Or local profile:
echo 'source <(anchors completion bash)' >> ~/.bashrc
```

### Fish
```bash
anchors completion fish > ~/.config/fish/completions/anchors.fish
```

---

## 7. Verifying Your Installation

Once installed, run:

```bash
# Print installed version and build metadata
anchors version

# Run ecosystem diagnostic
anchors doctor
```

If you see the output with no errors, you're ready to proceed to [The Workflow](/docs/workflow/) or inspect the [Commands Reference](/docs/cli/commands/)!
