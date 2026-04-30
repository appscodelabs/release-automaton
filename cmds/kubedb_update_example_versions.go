/*
Copyright AppsCode Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type kubedbActiveVersions map[string][]string

type kubedbDatabaseFile struct {
	Example kubedbDatabaseExample `json:"example"`
}

type kubedbDatabaseExample struct {
	Kind string             `json:"kind"`
	Spec kubedbDatabaseSpec `json:"spec"`
}

type kubedbDatabaseSpec struct {
	Version string `json:"version"`
}

type kubedbVersionUpdateResult struct {
	Filename string
	Kind     string
	Old      string
	New      string
	Changed  bool
	Reason   string
}

func NewCmdKubeDBUpdateExampleVersions() *cobra.Command {
	var workspace string
	var dbDir string
	var dryRun bool

	cmd := &cobra.Command{
		Use:               "update-example-versions <active-versions-file>",
		Short:             "Update example.spec.version in KubeDB database JSON files",
		DisableAutoGenTag: true,
		Args:              cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resolvedDir := dbDir
			if !filepath.IsAbs(resolvedDir) {
				resolvedDir = filepath.Join(workspace, resolvedDir)
			}
			return runKubeDBExampleVersionUpdate(args[0], resolvedDir, dryRun)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", ".", "Path to directory containing static-assets repository")
	cmd.Flags().StringVar(&dbDir, "dir", "data/products/kubedb/databases", "Directory containing KubeDB database JSON files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print updates without writing files")

	return cmd
}

func runKubeDBExampleVersionUpdate(versionsPath, dbDir string, dryRun bool) error {
	versions, err := loadKubeDBActiveVersions(versionsPath)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return fmt.Errorf("failed to read database directory %q: %w", dbDir, err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var updated int
	var skipped int
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filename := filepath.Join(dbDir, entry.Name())
		result, err := updateKubeDBVersionInFile(filename, versions, dryRun)
		if err != nil {
			return err
		}

		if result.Changed {
			updated++
			fmt.Printf("updated %s (%s): %s -> %s\n", result.Filename, result.Kind, result.Old, result.New)
			continue
		}

		skipped++
		fmt.Printf("skipped %s: %s\n", result.Filename, result.Reason)
	}

	fmt.Printf("done. updated=%d skipped=%d\n", updated, skipped)
	return nil
}

func loadKubeDBActiveVersions(path string) (kubedbActiveVersions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read active versions file %q: %w", path, err)
	}

	var versions kubedbActiveVersions
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("failed to parse active versions file %q: %w", path, err)
	}
	return versions, nil
}

func updateKubeDBVersionInFile(filename string, versions kubedbActiveVersions, dryRun bool) (kubedbVersionUpdateResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return kubedbVersionUpdateResult{}, fmt.Errorf("failed to read %q: %w", filename, err)
	}

	var doc kubedbDatabaseFile
	if err := json.Unmarshal(data, &doc); err != nil {
		return kubedbVersionUpdateResult{}, fmt.Errorf("failed to parse %q: %w", filename, err)
	}

	kind := doc.Example.Kind
	if kind == "" {
		return kubedbVersionUpdateResult{Filename: filename, Reason: "missing example.kind"}, nil
	}

	currentVersion := doc.Example.Spec.Version
	if currentVersion == "" {
		return kubedbVersionUpdateResult{Filename: filename, Kind: kind, Reason: "missing example.spec.version"}, nil
	}

	candidates, ok := versions[kind]
	if !ok || len(candidates) == 0 {
		return kubedbVersionUpdateResult{Filename: filename, Kind: kind, Old: currentVersion, Reason: "no active version entry for kind"}, nil
	}

	targetVersion, err := chooseKubeDBVersion(candidates, currentVersion)
	if err != nil {
		return kubedbVersionUpdateResult{}, fmt.Errorf("failed to choose version for %q (%s): %w", filename, kind, err)
	}

	if targetVersion == currentVersion {
		return kubedbVersionUpdateResult{Filename: filename, Kind: kind, Old: currentVersion, New: targetVersion, Reason: "already up to date"}, nil
	}

	oldNeedle := []byte(fmt.Sprintf(`"version": %q`, currentVersion))
	newNeedle := []byte(fmt.Sprintf(`"version": %q`, targetVersion))

	if bytes.Count(data, oldNeedle) != 1 {
		return kubedbVersionUpdateResult{}, fmt.Errorf("expected exactly one occurrence of %q in %q", oldNeedle, filename)
	}

	updated := bytes.Replace(data, oldNeedle, newNeedle, 1)
	if !dryRun {
		if err := os.WriteFile(filename, updated, 0o644); err != nil {
			return kubedbVersionUpdateResult{}, fmt.Errorf("failed to write %q: %w", filename, err)
		}
	}

	return kubedbVersionUpdateResult{Filename: filename, Kind: kind, Old: currentVersion, New: targetVersion, Changed: true}, nil
}

func chooseKubeDBVersion(candidates []string, current string) (string, error) {
	if len(candidates) == 0 {
		return "", errors.New("empty candidate list")
	}

	if idx := strings.Index(current, "-"); idx > 0 {
		prefix := current[:idx+1]
		if startsKubeDBVersionWithLetter(current) {
			if v := firstKubeDBVersionMatch(candidates, func(candidate string) bool {
				return strings.HasPrefix(candidate, prefix)
			}); v != "" {
				return v, nil
			}
		}

		if startsKubeDBVersionWithDigit(current) {
			suffix := current[idx:]
			if v := firstKubeDBVersionMatch(candidates, func(candidate string) bool {
				return strings.HasSuffix(candidate, suffix)
			}); v != "" {
				return v, nil
			}
		}
	}

	if !strings.Contains(current, "-") && startsKubeDBVersionWithDigit(current) {
		if v := firstKubeDBVersionMatch(candidates, func(candidate string) bool {
			return !strings.Contains(candidate, "-")
		}); v != "" {
			return v, nil
		}
	}

	return candidates[0], nil
}

func firstKubeDBVersionMatch(values []string, fn func(string) bool) string {
	for _, value := range values {
		if fn(value) {
			return value
		}
	}
	return ""
}

func startsKubeDBVersionWithLetter(s string) bool {
	if s == "" {
		return false
	}
	b := s[0]
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func startsKubeDBVersionWithDigit(s string) bool {
	if s == "" {
		return false
	}
	b := s[0]
	return b >= '0' && b <= '9'
}
