package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the configuration from the .llm-context.json file.
type Config struct {
	Dirs        []string          `json:"dir"`
	Files       []string          `json:"file"`
	Output      string            `json:"output"`
	CutComments bool              `json:"cut_comments"`
	Exceptions  map[string]string `json:"exceptions"`
}

func main() {
	cfg, err := findAndReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	var contentBuilder strings.Builder

	// Process individual files
	for _, file := range cfg.Files {
		processedContent, err := processFile(file, cfg.CutComments, cfg.Exceptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", file, err)
			continue
		}
		contentBuilder.WriteString(processedContent)
		contentBuilder.WriteString("\n\n")
	}

	// Process directories
	for _, dir := range cfg.Dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				// Ignore common nuisance files
				if d.Name() == ".DS_Store" || d.Name() == ".keep" {
					return nil
				}
				processedContent, err := processFile(path, cfg.CutComments, cfg.Exceptions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", path, err)
					return nil // Continue walking
				}
				contentBuilder.WriteString(processedContent)
				contentBuilder.WriteString("\n\n")
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", dir, err)
		}
	}

	// Write the final content to the output file
	err = os.WriteFile(cfg.Output, []byte(strings.TrimSpace(contentBuilder.String())+"\n"), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file %s: %v\n", cfg.Output, err)
		os.Exit(1)
	}

	fmt.Printf("Context successfully written to %s\n", cfg.Output)
}

// findAndReadConfig searches for .llm-context.json and unmarshals it.
func findAndReadConfig() (*Config, error) {
	configFileName := ".llm-context.json"
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", configFileName, err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", configFileName, err)
	}

	if cfg.Output == "" {
		return nil, fmt.Errorf("output file path is not specified in config")
	}

	return &cfg, nil
}

// processFile reads a file, optionally cleans it, and formats it with a header.
func processFile(filePath string, cutComments bool, exceptions map[string]string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	processedContent := string(content)
	if cutComments {
		processedContent = removeCommentsAndEmptyLines(processedContent)
	}

	fileName := filepath.Base(filePath)
	lang, ok := exceptions[fileName]
	if !ok {
		fileExt := filepath.Ext(filePath)
		lang = strings.TrimPrefix(fileExt, ".")
	}

	// Format with XML-style tags
	return fmt.Sprintf("<file name=\"%s\" lang=\"%s\">\n%s\n</file>", filePath, lang, processedContent), nil
}

// removeCommentsAndEmptyLines filters out comment lines and empty lines.
func removeCommentsAndEmptyLines(content string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			continue // Skip empty lines
		}
		if strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "#") || strings.HasPrefix(trimmedLine, "/*") || strings.HasPrefix(trimmedLine, "*") {
			continue // Skip comment lines
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
