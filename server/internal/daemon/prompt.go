package daemon

import (
	"fmt"
	"strings"
)

// BuildPrompt constructs the task prompt for an agent CLI.
// Keep this minimal — detailed instructions live in CLAUDE.md / AGENTS.md
// injected by execenv.InjectRuntimeConfig.
func BuildPrompt(task Task) string {
	if task.ChannelID != "" {
		return buildChannelPrompt(task)
	}
	return buildIssuePrompt(task)
}

func buildIssuePrompt(task Task) string {
	var b strings.Builder
	b.WriteString("You are running as a local coding agent for a Multica workspace.\n\n")
	fmt.Fprintf(&b, "Your assigned issue ID is: %s\n\n", task.IssueID)
	fmt.Fprintf(&b, "Start by running `multica issue get %s --output json` to understand your task, then complete it.\n", task.IssueID)
	return b.String()
}

func buildChannelPrompt(task Task) string {
	var b strings.Builder
	b.WriteString("You are running as a local coding agent for a Multica workspace.\n\n")
	b.WriteString("You have been @mentioned in a channel conversation. Your job is to read the recent messages, understand the context, and reply helpfully.\n\n")
	fmt.Fprintf(&b, "Channel ID: %s\n", task.ChannelID)
	if task.ChannelMessageID != "" {
		fmt.Fprintf(&b, "Triggering message ID: %s\n", task.ChannelMessageID)
	}
	b.WriteString("\nSteps:\n")
	fmt.Fprintf(&b, "1. Run `multica channel messages %s --output json` to read recent messages\n", task.ChannelID)
	b.WriteString("2. Understand what is being discussed and what is being asked of you\n")
	fmt.Fprintf(&b, "3. Reply with `multica channel reply %s --content \"your response\"`\n\n", task.ChannelID)
	b.WriteString("If the conversation suggests work that should be tracked as an issue:\n")
	fmt.Fprintf(&b, "- Suggest a task: `multica channel suggest %s --title \"...\" --description \"...\"`\n", task.ChannelID)
	return b.String()
}
