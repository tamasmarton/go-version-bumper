# Version Bumper

**Version Bumper** is a lightweight Go CLI tool that automates version management for Node.js projects.
It updates the `package.json` version (`major`, `minor`, or `patch`), stages and commits the change in Git,
and allows you to customize the commit message via config or command-line flags.

---

## Features

- bumps the version in `package.json` (or automatically `packages.json`) (`major` / `minor` / `patch`),
- stages and commits the file in Git,
- uses a configurable commit message template from a YAML file placed next to the executable,
- default commit message: just the version (`"%s"`).

## Installation

```
go build -o version-bumper .
```

Put the `version-bumper` binary on your PATH, or run it from the repository root.

## Usage

By default it performs a `patch` bump and automatically detects `package.json` / `packages.json`:

```
./version-bumper
```

Pár példa:

```
./version-bumper minor
./version-bumper major
./version-bumper -file path/to/package.json patch
./version-bumper -dry-run minor
./version-bumper -commit-template "chore(release): %s"
```

Flags:

- `-file` (default: `package.json`) – path to the JSON file; if this does not exist and `packages.json` does, that will be used.
- `-level` (default: `patch`) – `patch | minor | major`
- `-no-git` – skip `git add` and `git commit`
- `-commit-template` – override the commit message template via flag (fmt.Sprintf pattern, `%s` is the version)
- `-dry-run` – only prints what would happen, without making changes

## Configuration

The commit message template can also be provided via a `version-bump.yml` (or `version-bump.yaml`) file placed next to the executable.

Example:

```yaml
commit_message: 'chore(release): %s'
```

If no config file is found, the default is `"%s"`, meaning only the version number is used as the commit message.
Note: the template uses `fmt.Sprintf`, so `%s` will be replaced by the version.

## Notes

- The JSON is rewritten with pretty formatting (2-space indent). Key order and formatting may change on the first run.
- Pre-release / build-suffixed versions are accepted (e.g. `1.2.3-beta.1`), but after a bump the version is reset to the clean `MAJOR.MINOR.PATCH` format.

## License

MIT
