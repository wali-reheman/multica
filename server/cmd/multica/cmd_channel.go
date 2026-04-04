package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/multica-ai/multica/server/internal/cli"
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Manage channels and chat",
}

var channelMessagesCmd = &cobra.Command{
	Use:   "messages <channel-id>",
	Short: "List recent messages in a channel",
	Args:  cobra.ExactArgs(1),
	RunE:  runChannelMessages,
}

var channelReplyCmd = &cobra.Command{
	Use:   "reply <channel-id>",
	Short: "Send a message to a channel",
	Args:  cobra.ExactArgs(1),
	RunE:  runChannelReply,
}

var channelSuggestCmd = &cobra.Command{
	Use:   "suggest <channel-id>",
	Short: "Suggest a task in a channel (pending approval)",
	Args:  cobra.ExactArgs(1),
	RunE:  runChannelSuggest,
}

func init() {
	channelCmd.AddCommand(channelMessagesCmd)
	channelCmd.AddCommand(channelReplyCmd)
	channelCmd.AddCommand(channelSuggestCmd)

	channelMessagesCmd.Flags().String("output", "json", "Output format: json")
	channelMessagesCmd.Flags().Int("limit", 50, "Number of messages to return")

	channelReplyCmd.Flags().String("content", "", "Message content (required)")

	channelSuggestCmd.Flags().String("title", "", "Task title (required)")
	channelSuggestCmd.Flags().String("description", "", "Task description")
	channelSuggestCmd.Flags().String("priority", "none", "Task priority")
	channelSuggestCmd.Flags().String("assignee", "", "Assignee agent ID")
}

func runChannelMessages(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}

	channelID := args[0]
	limit, _ := cmd.Flags().GetInt("limit")
	output, _ := cmd.Flags().GetString("output")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var messages []any
	if err := client.GetJSON(ctx, fmt.Sprintf("/api/channels/%s/messages?limit=%d", channelID, limit), &messages); err != nil {
		return fmt.Errorf("list messages: %w", err)
	}

	if output == "json" {
		return cli.PrintJSON(os.Stdout, messages)
	}
	return cli.PrintJSON(os.Stdout, messages)
}

func runChannelReply(cmd *cobra.Command, args []string) error {
	content, _ := cmd.Flags().GetString("content")
	if content == "" {
		return fmt.Errorf("--content is required")
	}

	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}

	channelID := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := map[string]any{
		"content": content,
	}

	var result any
	if err := client.PostJSON(ctx, fmt.Sprintf("/api/channels/%s/messages", channelID), body, &result); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Message sent.\n")
	return nil
}

func runChannelSuggest(cmd *cobra.Command, args []string) error {
	title, _ := cmd.Flags().GetString("title")
	if title == "" {
		return fmt.Errorf("--title is required")
	}

	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}

	channelID := args[0]
	description, _ := cmd.Flags().GetString("description")
	priority, _ := cmd.Flags().GetString("priority")
	assignee, _ := cmd.Flags().GetString("assignee")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := map[string]any{
		"title":    title,
		"priority": priority,
	}
	if description != "" {
		body["description"] = description
	}
	if assignee != "" {
		body["assignee_type"] = "agent"
		body["assignee_id"] = assignee
	}

	var result any
	if err := client.PostJSON(ctx, fmt.Sprintf("/api/channels/%s/suggestions", channelID), body, &result); err != nil {
		return fmt.Errorf("suggest task: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Task suggested (pending approval).\n")
	return nil
}
