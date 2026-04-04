"use client";

import { Bot, ExternalLink } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";
import { useActorName } from "@/features/workspace";
import type { ChannelMessage } from "@/shared/types";

interface MessageItemProps {
  message: ChannelMessage;
}

export function MessageItem({ message }: MessageItemProps) {
  const { getActorName, getActorInitials, getActorAvatarUrl } = useActorName();

  const isAgent = message.author_type === "agent";
  const isSystem = message.type === "system" || message.type === "issue_created";
  const name = getActorName(message.author_type, message.author_id);
  const initials = getActorInitials(message.author_type, message.author_id);
  const avatarUrl = getActorAvatarUrl(message.author_type, message.author_id);

  if (isSystem) {
    return (
      <div className="flex items-center gap-2 py-1.5 px-2">
        <div className="h-px flex-1 bg-border" />
        <span className="text-xs text-muted-foreground flex items-center gap-1">
          {message.type === "issue_created" && <ExternalLink className="h-3 w-3" />}
          {message.content}
        </span>
        <div className="h-px flex-1 bg-border" />
      </div>
    );
  }

  const time = new Date(message.created_at).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <div className="group flex gap-3 rounded-md px-2 py-2 hover:bg-muted/50">
      <Avatar className={cn("h-8 w-8 shrink-0", isAgent && "ring-2 ring-purple-500/20")}>
        {avatarUrl && <AvatarImage src={avatarUrl} alt={name} />}
        <AvatarFallback className={cn("text-xs", isAgent && "bg-purple-100 text-purple-700")}>
          {isAgent ? <Bot className="h-4 w-4" /> : initials}
        </AvatarFallback>
      </Avatar>
      <div className="min-w-0 flex-1">
        <div className="flex items-baseline gap-2">
          <span className={cn("text-sm font-semibold", isAgent && "text-purple-600")}>
            {name}
          </span>
          <span className="text-xs text-muted-foreground">{time}</span>
        </div>
        <p className="text-sm whitespace-pre-wrap break-words">{message.content}</p>
      </div>
    </div>
  );
}
