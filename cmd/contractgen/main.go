package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	root, err := moduleRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "contractgen: %v\n", err)
		os.Exit(1)
	}

	contractsDir := filepath.Join(root, "contracts")
	testsDir := filepath.Join(root, "contracttests")

	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "contractgen: mkdir contracttests: %v\n", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir(contractsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "contractgen: readdir contracts: %v\n", err)
		os.Exit(1)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}

		yamlPath := filepath.Join(contractsDir, e.Name())
		data, err := os.ReadFile(yamlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "contractgen: read %s: %v\n", e.Name(), err)
			os.Exit(1)
		}

		var c Contract
		if err := yaml.Unmarshal(data, &c); err != nil {
			fmt.Fprintf(os.Stderr, "contractgen: parse %s: %v\n", e.Name(), err)
			os.Exit(1)
		}

		base := strings.TrimSuffix(e.Name(), ".yaml")

		// Generate contracts/*_gen.go
		{
			var buf bytes.Buffer
			if err := contractsTmpl.Execute(&buf, c); err != nil {
				fmt.Fprintf(os.Stderr, "contractgen: contracts template %s: %v\n", e.Name(), err)
				os.Exit(1)
			}
			out := filepath.Join(contractsDir, base+"_gen.go")
			if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "contractgen: write %s: %v\n", out, err)
				os.Exit(1)
			}
			fmt.Printf("contractgen: wrote %s\n", out)
		}

		// Generate contracttests/*_contract_gen.go
		{
			var buf bytes.Buffer
			if err := testsTmpl.Execute(&buf, c); err != nil {
				fmt.Fprintf(os.Stderr, "contractgen: tests template %s: %v\n", e.Name(), err)
				os.Exit(1)
			}
			out := filepath.Join(testsDir, base+"_contract_gen.go")
			if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "contractgen: write %s: %v\n", out, err)
				os.Exit(1)
			}
			fmt.Printf("contractgen: wrote %s\n", out)
		}
	}
}

// moduleRoot walks up from the current working directory to find go.mod.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
