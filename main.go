package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// CommitMessage is a fmt.Sprintf template with one %s placeholder for the version.
	// Example: "chore(release): %s"
	CommitMessage string `yaml:"commit_message"`
}

// loadConfig tries to load a YAML config placed next to the executable.
// It looks for "version-bump.yml" or "version-bump.yaml". If not found, returns defaults.
func loadConfig() (Config, error) {
	cfg := Config{
		CommitMessage: "%s",
	}
	exePath, err := os.Executable()
	if err != nil {
		return cfg, err
	}
	exeDir := filepath.Dir(exePath)

	candidates := []string{
		filepath.Join(exeDir, "version-bump.yml"),
		filepath.Join(exeDir, "version-bump.yaml"),
	}

	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil {
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return cfg, fmt.Errorf("config parse error in %s: %w", p, err)
			}
			return cfg, nil
		}
	}

	return cfg, nil
}

type PackageJSON struct {
	raw map[string]any
}

func readPackageJSON(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var m map[string]any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("parsing JSON %s: %w", path, err)
	}
	return &PackageJSON{raw: m}, nil
}

func (p *PackageJSON) GetVersion() (string, error) {
	v, ok := p.raw["version"]
	if !ok {
		return "", errors.New(`"version" field not found in JSON file`)
	}
	vs, ok := v.(string)
	if !ok {
		return "", errors.New(`"version" field is not a string`)
	}
	return vs, nil
}

func (p *PackageJSON) SetVersion(v string) {
	p.raw["version"] = v
}

func (p *PackageJSON) Write(path string) error {
	// Pretty print with 2-space indent (common in JS projects).
	data, err := json.MarshalIndent(p.raw, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing JSON: %w", err)
	}
	// Ensure trailing newline
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

var semverRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:[-+].*)?$`)

func bumpVersion(current string, level string) (string, error) {
	m := semverRe.FindStringSubmatch(current)
	if m == nil {
		return "", fmt.Errorf("version %q is not a supported semver (expected MAJOR.MINOR.PATCH optionally with pre-release/build)", current)
	}
	maj := atoi(m[1])
	min := atoi(m[2])
	patch := atoi(m[3])

	switch level {
	case "patch":
		patch++
	case "minor":
		min++
		patch = 0
	case "major":
		maj++
		min = 0
		patch = 0
	default:
		return "", fmt.Errorf("unknown bump level %q (use: patch, minor, major)", level)
	}

	return fmt.Sprintf("%d.%d.%d", maj, min, patch), nil
}

func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func main() {
	var filePath string
	var level string
	var noGit bool
	var commitTemplateOverride string
	var dryRun bool

	flag.StringVar(&filePath, "file", "package.json", "Path to package.json (or packages.json if that is your file)")
	flag.StringVar(&level, "level", "patch", "Bump level: patch | minor | major")
	flag.BoolVar(&noGit, "no-git", false, "Do not run git add/commit")
	flag.StringVar(&commitTemplateOverride, "commit-template", "", "Override commit message template (fmt.Sprintf pattern, use %s for version)")
	flag.BoolVar(&dryRun, "dry-run", false, "Show what would change without writing or committing")
	flag.Parse()

	// Backwards compatible: allow positional arg as level if provided.
	if flag.NArg() >= 1 {
		level = strings.ToLower(flag.Arg(0))
	}

	// Auto-detect packages.json if package.json doesn't exist and user didn't override -file.
	if filePath == "package.json" && !fileExists(filePath) && fileExists("packages.json") {
		filePath = "packages.json"
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
	}
	commitTpl := cfg.CommitMessage
	if commitTpl == "" {
		commitTpl = "%s"
	}
	if commitTemplateOverride != "" {
		commitTpl = commitTemplateOverride
	}

	// Read and bump version
	pkg, err := readPackageJSON(filePath)
	if err != nil {
		exitErr(err)
	}
	curVer, err := pkg.GetVersion()
	if err != nil {
		exitErr(err)
	}
	newVer, err := bumpVersion(curVer, strings.ToLower(level))
	if err != nil {
		exitErr(err)
	}

	if dryRun {
		fmt.Printf("Current version: %s\n", curVer)
		fmt.Printf("Bump level: %s\n", level)
		fmt.Printf("New version: %s\n", newVer)
		fmt.Printf("Commit template: %q -> message: %q\n", commitTpl, fmt.Sprintf(commitTpl, newVer))
		fmt.Println("Dry run: no files written, no git commands executed.")
		return
	}

	// Write updated version
	pkg.SetVersion(newVer)
	if err := pkg.Write(filePath); err != nil {
		exitErr(err)
	}
	fmt.Printf("Updated %s: %s -> %s\n", filePath, curVer, newVer)

	if noGit {
		return
	}

	// Stage and commit
	if err := runGit("add", filePath); err != nil {
		exitErr(fmt.Errorf("git add failed: %w", err))
	}
	commitMsg := fmt.Sprintf(commitTpl, newVer)
	if err := runGit("commit", "-m", commitMsg); err != nil {
		exitErr(fmt.Errorf("git commit failed: %w", err))
	}
	fmt.Printf("Committed with message: %q\n", commitMsg)
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}