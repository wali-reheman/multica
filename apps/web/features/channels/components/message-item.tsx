"use client";

import { Bot, ExternalLink } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Tooltip, TooltipTrigger, TooltipContent } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { useActorName } from "@/features/workspace";
import type { ChannelMessage } from "@/shared/types";

interface MessageItemProps {
  message: ChannelMessage;
  grouped?: boolean;
}

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatFullDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString([], {
    weekday: "short",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function MessageItem({ message, grouped = false }: MessageItemProps) {
  const { getActorName, getActorInitials, getActorAvatarUrl } = useActorName();

  const isAgent = message.author_type === "agent";
  const isSystem = message.type === "system" || message.type === "issue_created";
  const name = getActorName(message.author_type, message.author_id);
  const initials = getActorInitials(message.author_type, message.author_id);
  const avatarUrl = getActorAvatarUrl(message.author_type, message.author_id);

  if (isSystem) {
    return (
      <div className="flex items-center gap-3 py-1 px-2 my-1">
        <div className="h-px flex-1 bg-border" />
        <span className="text-[11px] text-muted-foreground flex items-center gap-1 shrink-0">
          {message.type === "issue_created" && <ExternalLink className="h-3 w-3" />}
          {message.content}
        </span>
        <div className="h-px flex-1 bg-border" />
      </div>
    );
  }

  // Grouped message (same author within 5 minutes) — compact, no avatar
  if (grouped) {
    return (
      <div className="group flex gap-3 rounded-md pl-[44px] pr-2 py-0.5 hover:bg-muted/40">
        <Tooltip>
          <TooltipTrigger className="invisible group-hover:visible shrink-0 text-[10px] text-muted-foreground w-0 -ml-[36px]">
            {formatTime(message.created_at)}
          </TooltipTrigger>
          <TooltipContent side="left" className="text-xs">
            {formatFullDate(message.created_at)}
          </TooltipContent>
        </Tooltip>
        <p className="text-sm whitespace-pre-wrap break-words min-w-0">{message.content}</p>
      </div>
    );
  }

  return (
    <div className="group flex gap-3 rounded-md px-2 py-1.5 hover:bg-muted/40 mt-3 first:mt-0">
      <Avatar className={cn("h-8 w-8 shrink-0 mt-0.5", isAgent && "ring-2 ring-purple-500/20")}>
        {avatarUrl && <AvatarImage src={avatarUrl} alt={name} />}
        <AvatarFallback className={cn("text-xs", isAgent && "bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-300")}>
          {isAgent ? <Bot className="h-4 w-4" /> : initials}
        </AvatarFallback>
      </Avatar>
      <div className="min-w-0 flex-1">
        <div className="flex items-baseline gap-2">
          <span className={cn("text-sm font-semibold", isAgent && "text-purple-600 dark:text-purple-400")}>
            {name}
          </span>
          <Tooltip>
            <TooltipTrigger className="text-[11px] text-muted-foreground">
              {formatTime(message.created_at)}
            </TooltipTrigger>
            <TooltipContent side="top" className="text-xs">
              {formatFullDate(message.created_at)}
            </TooltipContent>
          </Tooltip>
        </div>
        <p className="text-sm whitespace-pre-wrap break-words">{message.content}</p>
      </div>
    </div>
  );
}
