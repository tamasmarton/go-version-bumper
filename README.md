# Version Bumper

<p align="left">
  <a href="https://github.com/tamasmarton/go-version-bumper/actions/workflows/release.yaml">
    <img src="https://github.com/tamasmarton/go-version-bumper/actions/workflows/release.yaml/badge.svg" alt="Build">
  </a>
  <a href="https://github.com/tamasmarton/go-version-bumper/releases">
    <img src="https://img.shields.io/github/v/release/tamasmarton/go-version-bumper?sort=semver" alt="Latest Release">
  </a>
  <a href="https://github.com/tamasmarton/go-version-bumper/releases">
    <img src="https://img.shields.io/github/downloads/tamasmarton/go-version-bumper/total" alt="Downloads">
  </a>
  <a href="https://pkg.go.dev/github.com/tamasmarton/go-version-bumper">
    <img src="https://img.shields.io/github/go-mod/go-version/tamasmarton/go-version-bumper" alt="Go Version">
  </a>
  <a href="https://github.com/tamasmarton/go-version-bumper/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/tamasmarton/go-version-bumper" alt="License">
  </a>
  <a href="https://goreportcard.com/report/github.com/tamasmarton/go-version-bumper">
    <img src="https://goreportcard.com/badge/github.com/tamasmarton/go-version-bumper" alt="Go Report Card">
  </a>
</p>

Version Bumper is a lightweight Go CLI that automates version management for Node.js (and similar) projects by bumping the `version` field inside a `package.json` (or a fallback `packages.json`), committing the change, and letting you customize the commit message.

---

## Table of Contents

- [Features](#features)
- [Quick Install](#quick-install)
- [Usage](#usage)
- [Flags](#flags)
- [Configuration](#configuration)
- [Examples](#examples)
- [Behavior & Notes](#behavior--notes)
- [Planned / Ideas](#planned--ideas)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- Bumps the version in `package.json` (or automatically falls back to `packages.json`) using `patch`, `minor`, or `major`.
- Stages and commits the modified file with a customizable commit message.
- Commit message template can come from:
  - Command‑line flag (`-commit-template`)
  - YAML config file (`version-bump.yml` / `version-bump.yaml`) placed next to the executable
  - Default template falls back to `"%s"` (just the new version).
- Dry‑run mode to preview changes without writing or committing.
- Pretty reformatting (2‑space indentation) of the JSON.
- Accepts pre-release inputs (e.g. `1.2.3-beta.1`) and normalizes output to plain `MAJOR.MINOR.PATCH` after the bump.

---

## Quick Install

### Go (direct)

```
go install github.com/tamasmarton/go-version-bumper@latest
```

Places the `version-bumper` binary in `$GOBIN` (or `$GOPATH/bin`).

### Homebrew (after the tap formula is published)

```
brew tap tamasmarton/homebrew-tap
brew install version-bumper
```

### Manual build

```
git clone https://github.com/tamasmarton/go-version-bumper.git
cd go-version-bumper
go build -o version-bumper .
```

Add the resulting binary to your PATH if desired.

---

## Usage

Default (patch bump, auto-detects `package.json` first, else `packages.json`):

```
./version-bumper
```

Specify the level explicitly:

```
./version-bumper minor
./version-bumper major
```

Provide a custom file:

```
./version-bumper -file path/to/package.json patch
```

Dry run (no write, no commit):

```
./version-bumper -dry-run minor
```

Custom commit template at runtime:

```
./version-bumper -commit-template "chore(release): %s"
```

Skip git operations (only modify file):

```
./version-bumper -no-git patch
```

---

## Flags

| Flag               | Default                                                   | Description                                                                                       |
| ------------------ | --------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `-file`            | `package.json` (fallback to `packages.json` if not found) | Path to the JSON file to update.                                                                  |
| `-level`           | `patch`                                                   | One of `patch`, `minor`, `major`. (You can also supply the level as a positional argument.)       |
| `-no-git`          | `false`                                                   | If set, skips `git add` and `git commit`.                                                         |
| `-commit-template` | (none)                                                    | Overrides commit message template (`fmt.Sprintf` pattern, `%s` is replaced with the new version). |
| `-dry-run`         | `false`                                                   | Prints planned changes without modifying files or running git.                                    |

Order of precedence for determining the bump level:

1. Positional argument (e.g. `./version-bumper minor`)
2. `-level` flag
3. Default: `patch`

---

## Configuration

You can define a YAML file named `version-bump.yml` or `version-bump.yaml` next to the executable to provide a default commit template:

```yaml
commit_message: 'chore(release): %s'
```

If not present, the default template is simply:

```
"%s"
```

Template rules:

- Uses Go `fmt.Sprintf`, so `%s` is replaced by the new version string.
- No other automatic placeholders exist (keep it simple).

---

## Examples

Patch bump (default):

```
./version-bumper
```

Minor bump with dry run and custom commit template (no changes written):

```
./version-bumper -dry-run -commit-template "release: %s" minor
```

Major bump, custom file, skip git:

```
./version-bumper -file ./frontend/package.json -no-git major
```

Commit message via config file (place `version-bump.yml` beside binary):

```
commit_message: "chore(release): %s"
```

Run:

```
./version-bumper patch
```

---

## Behavior & Notes

- JSON rewrite: The tool rewrites the entire file with consistent 2‑space indentation; key ordering may change the first time.
- Pre-releases: Input versions like `1.2.3-beta.1` are parsed; the bump logic strips pre-release metadata and increments the base semver part.
- Git operations: Unless `-no-git` or `-dry-run` is set, the modified file is staged and committed in the current repository.
- Safety: If parsing fails or the `version` field is missing/invalid, the tool exits with a non-zero status.
- Idempotence: Running the same bump level repeatedly will increment each time (there is no caching).

---

## Planned / Ideas

Potential future improvements (PRs welcome):

- Support custom version field names.
- Pre-release bump flags (e.g. `--pre beta`).
- Optional tag creation (`git tag vX.Y.Z`) and push automation.
- JSON path override (e.g. nested version fields).

Open an issue to discuss before larger changes.

---

## Contributing

1. Fork the repository.
2. Create a feature branch:
   ```
   git checkout -b feat/your-idea
   ```
3. Build & test locally.
4. Open a PR with a clear description and rationale.

Please keep the tool lean—avoid adding large dependency chains or unrelated features.

---

## License

MIT — see [LICENSE](./LICENSE) for details.
