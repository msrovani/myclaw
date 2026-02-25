package skills

import (
	"bufio"
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// AgentSkill represents a dynamically loaded skill from the .agent markdown files.
type AgentSkill struct {
	id      string
	desc    string
	content string // The markdown body
}

func (a AgentSkill) ID() string {
	return a.id
}

func (a AgentSkill) Description() string {
	return a.desc
}

// Execute for an AgentSkill simply returns its markdown content for now.
// A full implementation would pipe this into an LLM or execute bash scripts.
func (a AgentSkill) Execute(ctx context.Context, req Request) (Response, error) {
	return Response{
		Result: map[string]any{
			"content": a.content,
		},
	}, nil
}

// LoadAgentDir scans the provided root directory (e.g., ".agent") for markdown files,
// parses their YAML frontmatter for 'name' and 'description', and registers them into the runtime.
func LoadAgentDir(root string, registry *Registry) error {
	slog.Info("skills: scanning .agent mount", "path", root)

	if _, err := os.Stat(root); os.IsNotExist(err) {
		slog.Warn("skills: .agent mount not found, skipping integration bridge", "path", root)
		return nil
	}

	extracted := 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Very rudimentary YAML frontmatter parser
		scanner := bufio.NewScanner(file)
		var name, desc string
		var content strings.Builder
		inFrontmatter := false
		lineCount := 0

		for scanner.Scan() {
			line := scanner.Text()
			lineCount++

			if lineCount == 1 && strings.TrimSpace(line) == "---" {
				inFrontmatter = true
				continue
			}

			if inFrontmatter {
				if strings.TrimSpace(line) == "---" {
					inFrontmatter = false
					continue
				}

				// naive key: value parse
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					if key == "name" {
						name = val
					} else if key == "description" {
						desc = val
					}
				}
			} else {
				content.WriteString(line)
				content.WriteString("\n")
			}
		}

		// Fallbacks if missing frontmatter
		if name == "" {
			name = strings.TrimSuffix(filepath.Base(path), ".md")
		}
		if desc == "" {
			desc = "Dynamically mounted from " + path
		}

		skill := AgentSkill{
			id:      name,
			desc:    desc,
			content: content.String(),
		}

		if err := registry.Register(skill); err != nil {
			if errors.Is(err, ErrSkillAlreadyRegistered) {
				slog.Debug("skills: duplicate agent skill skipped", "id", name)
			} else {
				slog.Warn("skills: failed to mount agent skill", "id", name, "error", err)
			}
		} else {
			extracted++
		}

		return nil
	})

	slog.Info("skills: integration bridge complete", "mounted", extracted)
	return err
}
