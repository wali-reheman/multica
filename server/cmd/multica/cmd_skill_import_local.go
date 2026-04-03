package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/multica-ai/multica/server/internal/cli"
)

// localItemType identifies the kind of Claude Code artifact being imported.
type localItemType string

const (
	localTypeSkill       localItemType = "skill"
	localTypeAgent       localItemType = "agent"
	localTypeCommand     localItemType = "command"
	localTypeHook        localItemType = "hook"
	localTypeConventions localItemType = "conventions"
)

// localImportItem represents a single item discovered on disk, ready to upload.
type localImportItem struct {
	Name        string
	Description string
	Content     string // SKILL.md body or main file content
	Files       []localImportFile
	Type        localItemType
	SourcePath  string // absolute path on disk
}

type localImportFile struct {
	Path    string // relative path within the skill
	Content string
}

var skillImportLocalCmd = &cobra.Command{
	Use:   "import-local",
	Short: "Import skills, agents, commands, and hooks from a local Claude Code directory",
	Long: `Scan a local Claude Code directory (defaults to ~/.claude/) and import
artifacts into the current Multica workspace as skills.

Supported types:
  skill       SKILL.md-based skills from skills/ directory
  agent       Agent definitions from agents/ directory
  command     Slash commands from commands/ directory
  hook        Hook scripts from hooks/ directory
  conventions CLAUDE.md coding conventions
  all         All of the above (default)

Each imported item is stored as a Multica skill with source metadata in
its config field, so you can identify locally-imported items in the UI.

Examples:
  multica skill import-local                        # Import everything from ~/.claude/
  multica skill import-local --path ~/my-project/.claude/
  multica skill import-local --type skill            # Only skills
  multica skill import-local --type agent,command    # Agents and commands
  multica skill import-local --dry-run               # Preview without importing`,
	RunE: runSkillImportLocal,
}

func init() {
	skillCmd.AddCommand(skillImportLocalCmd)

	skillImportLocalCmd.Flags().String("path", "", "Root directory to scan (default: ~/.claude/)")
	skillImportLocalCmd.Flags().String("type", "all", "What to import: skill, agent, command, hook, conventions, all (comma-separated)")
	skillImportLocalCmd.Flags().Bool("dry-run", false, "Preview what would be imported without creating anything")
	skillImportLocalCmd.Flags().String("output", "table", "Output format: table or json")
}

func runSkillImportLocal(cmd *cobra.Command, _ []string) error {
	rootPath, _ := cmd.Flags().GetString("path")
	if rootPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		rootPath = filepath.Join(home, ".claude")
	}

	// Expand ~ if present
	if strings.HasPrefix(rootPath, "~") {
		home, _ := os.UserHomeDir()
		rootPath = filepath.Join(home, rootPath[1:])
	}

	info, err := os.Stat(rootPath)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("directory not found: %s", rootPath)
	}

	typeFlag, _ := cmd.Flags().GetString("type")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	output, _ := cmd.Flags().GetString("output")

	types := parseTypeFlags(typeFlag)

	// Discover items
	var items []localImportItem

	if types["skill"] {
		found, err := discoverSkills(rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error scanning skills: %v\n", err)
		}
		items = append(items, found...)
	}

	if types["agent"] {
		found, err := discoverAgents(rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error scanning agents: %v\n", err)
		}
		items = append(items, found...)
	}

	if types["command"] {
		found, err := discoverCommands(rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error scanning commands: %v\n", err)
		}
		items = append(items, found...)
	}

	if types["hook"] {
		found, err := discoverHooks(rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error scanning hooks: %v\n", err)
		}
		items = append(items, found...)
	}

	if types["conventions"] {
		found, err := discoverConventions(rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error scanning conventions: %v\n", err)
		}
		items = append(items, found...)
	}

	if len(items) == 0 {
		fmt.Println("No items found to import.")
		return nil
	}

	// Dry run — just list
	if dryRun {
		return printDryRun(items, output)
	}

	// Actual import
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}

	return importItems(client, items, output)
}

// parseTypeFlags parses a comma-separated type flag into a lookup map.
func parseTypeFlags(flag string) map[string]bool {
	m := map[string]bool{}
	flag = strings.TrimSpace(strings.ToLower(flag))
	if flag == "all" || flag == "" {
		m["skill"] = true
		m["agent"] = true
		m["command"] = true
		m["hook"] = true
		m["conventions"] = true
		return m
	}
	for _, t := range strings.Split(flag, ",") {
		m[strings.TrimSpace(t)] = true
	}
	return m
}

// ---------------------------------------------------------------------------
// Discovery functions
// ---------------------------------------------------------------------------

// discoverSkills scans <root>/skills/*/SKILL.md
func discoverSkills(root string) ([]localImportItem, error) {
	skillsDir := filepath.Join(root, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var items []localImportItem
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		skillDir := filepath.Join(skillsDir, e.Name())
		skillMdPath := filepath.Join(skillDir, "SKILL.md")
		content, err := os.ReadFile(skillMdPath)
		if err != nil {
			continue // no SKILL.md — skip
		}

		name, description := parseLocalSkillFrontmatter(string(content))
		if name == "" {
			name = e.Name()
		}

		item := localImportItem{
			Name:        name,
			Description: description,
			Content:     string(content),
			Type:        localTypeSkill,
			SourcePath:  skillDir,
		}

		// Collect supporting files
		err = filepath.Walk(skillDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			relPath, _ := filepath.Rel(skillDir, path)
			if strings.EqualFold(relPath, "SKILL.md") {
				return nil // already captured as content
			}
			// Skip very large files (>1MB)
			if info.Size() > 1<<20 {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			item.Files = append(item.Files, localImportFile{
				Path:    relPath,
				Content: string(data),
			})
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error walking skill directory %s: %v\n", skillDir, err)
		}

		items = append(items, item)
	}
	return items, nil
}

// discoverAgents scans <root>/agents/*.md
func discoverAgents(root string) ([]localImportItem, error) {
	agentsDir := filepath.Join(root, "agents")
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var items []localImportItem
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(agentsDir, e.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		// Skip very large files (>1MB)
		if len(content) > 1<<20 {
			fmt.Fprintf(os.Stderr, "Warning: skipping large agent file %s (>1MB)\n", e.Name())
			continue
		}

		baseName := strings.TrimSuffix(e.Name(), ".md")
		name := "agent:" + baseName
		description := fmt.Sprintf("Agent definition: %s", baseName)

		items = append(items, localImportItem{
			Name:        name,
			Description: description,
			Content:     string(content),
			Type:        localTypeAgent,
			SourcePath:  path,
		})
	}
	return items, nil
}

// discoverCommands scans <root>/commands/**/*.md recursively
func discoverCommands(root string) ([]localImportItem, error) {
	commandsDir := filepath.Join(root, "commands")
	if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
		return nil, nil
	}

	var items []localImportItem
	err := filepath.Walk(commandsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		if info.Size() > 1<<20 {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(commandsDir, path)
		// Build command name from path: gsd/execute-phase.md -> cmd:gsd:execute-phase
		cmdName := strings.TrimSuffix(relPath, ".md")
		cmdName = strings.ReplaceAll(cmdName, string(filepath.Separator), ":")
		name := "cmd:" + cmdName

		// Use the first non-empty line as description
		description := firstLine(string(content))
		if len(description) > 120 {
			description = description[:117] + "..."
		}

		items = append(items, localImportItem{
			Name:        name,
			Description: description,
			Content:     string(content),
			Type:        localTypeCommand,
			SourcePath:  path,
		})
		return nil
	})
	return items, err
}

// discoverHooks scans <root>/hooks/*.js (and *.sh, *.py)
func discoverHooks(root string) ([]localImportItem, error) {
	hooksDir := filepath.Join(root, "hooks")
	entries, err := os.ReadDir(hooksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	hookExts := map[string]bool{".js": true, ".sh": true, ".py": true, ".ts": true}

	var items []localImportItem
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if !hookExts[ext] {
			continue
		}
		path := filepath.Join(hooksDir, e.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if len(content) > 1<<20 {
			continue
		}

		baseName := strings.TrimSuffix(e.Name(), ext)
		name := "hook:" + baseName
		description := fmt.Sprintf("Hook script: %s", e.Name())

		items = append(items, localImportItem{
			Name:        name,
			Description: description,
			Content:     string(content),
			Type:        localTypeHook,
			SourcePath:  path,
		})
	}
	return items, nil
}

// discoverConventions reads <root>/CLAUDE.md
func discoverConventions(root string) ([]localImportItem, error) {
	claudePath := filepath.Join(root, "CLAUDE.md")
	content, err := os.ReadFile(claudePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return []localImportItem{{
		Name:        "claude-md",
		Description: "Global CLAUDE.md coding conventions",
		Content:     string(content),
		Type:        localTypeConventions,
		SourcePath:  claudePath,
	}}, nil
}

// ---------------------------------------------------------------------------
// Import logic
// ---------------------------------------------------------------------------

func importItems(client *cli.APIClient, items []localImportItem, output string) error {
	type resultEntry struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Status string `json:"status"`
		ID     string `json:"id,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	var results []resultEntry
	created := 0
	skipped := 0
	failed := 0

	for _, item := range items {
		// Build request body
		files := make([]map[string]string, 0, len(item.Files))
		for _, f := range item.Files {
			files = append(files, map[string]string{
				"path":    f.Path,
				"content": f.Content,
			})
		}

		config := map[string]any{
			"source":     "local",
			"local_type": string(item.Type),
			"local_path": item.SourcePath,
		}

		body := map[string]any{
			"name":        item.Name,
			"description": item.Description,
			"content":     item.Content,
			"config":      config,
		}
		if len(files) > 0 {
			body["files"] = files
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var result map[string]any
		err := client.PostJSON(ctx, "/api/skills", body, &result)
		cancel()

		entry := resultEntry{
			Name: item.Name,
			Type: string(item.Type),
		}

		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "409") || strings.Contains(errMsg, "already exists") {
				entry.Status = "skipped"
				entry.Error = "already exists"
				skipped++
			} else {
				entry.Status = "failed"
				entry.Error = errMsg
				failed++
			}
		} else {
			entry.Status = "created"
			entry.ID = strVal(result, "id")
			created++
		}

		results = append(results, entry)

		// Print progress for table output
		if output != "json" {
			icon := "+"
			switch entry.Status {
			case "skipped":
				icon = "~"
			case "failed":
				icon = "!"
			}
			suffix := ""
			if entry.Error != "" {
				suffix = " (" + entry.Error + ")"
			}
			fmt.Printf("  %s %-14s %s%s\n", icon, "["+entry.Type+"]", entry.Name, suffix)
		}
	}

	if output == "json" {
		return cli.PrintJSON(os.Stdout, map[string]any{
			"created": created,
			"skipped": skipped,
			"failed":  failed,
			"items":   results,
		})
	}

	fmt.Printf("\nDone: %d created, %d skipped, %d failed (total: %d)\n", created, skipped, failed, len(items))
	return nil
}

func printDryRun(items []localImportItem, output string) error {
	if output == "json" {
		type previewItem struct {
			Name       string `json:"name"`
			Type       string `json:"type"`
			SourcePath string `json:"source_path"`
			FileCount  int    `json:"file_count"`
		}
		preview := make([]previewItem, len(items))
		for i, it := range items {
			preview[i] = previewItem{
				Name:       it.Name,
				Type:       string(it.Type),
				SourcePath: it.SourcePath,
				FileCount:  len(it.Files),
			}
		}
		return cli.PrintJSON(os.Stdout, map[string]any{
			"total": len(items),
			"items": preview,
		})
	}

	headers := []string{"TYPE", "NAME", "FILES", "SOURCE"}
	rows := make([][]string, 0, len(items))
	for _, it := range items {
		rows = append(rows, []string{
			string(it.Type),
			it.Name,
			fmt.Sprintf("%d", len(it.Files)),
			it.SourcePath,
		})
	}
	cli.PrintTable(os.Stdout, headers, rows)
	fmt.Printf("\n%d items would be imported. Run without --dry-run to proceed.\n", len(items))
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseLocalSkillFrontmatter extracts name and description from YAML frontmatter.
func parseLocalSkillFrontmatter(content string) (name, description string) {
	if !strings.HasPrefix(content, "---") {
		return "", ""
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return "", ""
	}
	frontmatter := content[3 : 3+end]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
			name = strings.Trim(name, "\"'")
		} else if strings.HasPrefix(line, "description:") {
			description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			description = strings.Trim(description, "\"'")
		}
	}
	return name, description
}

// firstLine returns the first non-empty, non-frontmatter line from content.
func firstLine(content string) string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			continue
		}
		// Skip markdown headings prefix
		cleaned := strings.TrimLeft(trimmed, "# ")
		if cleaned != "" {
			return cleaned
		}
	}
	return ""
}
