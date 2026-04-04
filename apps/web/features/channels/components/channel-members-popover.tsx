"use client";

import { useState } from "react";
import { Users, Bot, UserPlus } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";
import { api } from "@/shared/api";
import { useActorName } from "@/features/workspace";
import { useWorkspaceStore } from "@/features/workspace";
import { useChannelStore } from "../store";
import type { ChannelMember } from "@/shared/types";

interface ChannelMembersPopoverProps {
  channelId: string;
  members: ChannelMember[];
}

export function ChannelMembersPopover({ channelId, members }: ChannelMembersPopoverProps) {
  const { getActorName, getActorInitials } = useActorName();
  const agents = useWorkspaceStore((s) => s.agents).filter((a) => !a.archived_at);
  const wsMembers = useWorkspaceStore((s) => s.members);
  const [addType, setAddType] = useState<"agent" | "member">("agent");
  const [addId, setAddId] = useState("");

  const handleAdd = async () => {
    if (!addId) return;
    try {
      await api.addChannelMember(channelId, addType, addId);
      useChannelStore.getState().fetch();
      setAddId("");
      toast.success("Member added");
    } catch {
      toast.error("Failed to add member");
    }
  };

  // Filter out already-added entities.
  const memberIds = new Set(members.map((m) => m.member_id));
  const availableAgents = agents.filter((a) => !memberIds.has(a.id));
  const availableMembers = wsMembers.filter((m) => !memberIds.has(m.user_id));

  return (
    <Popover>
      <PopoverTrigger
        render={
          <Button variant="ghost" size="sm" className="gap-1.5">
            <Users className="h-4 w-4" />
            <span>{members.length}</span>
          </Button>
        }
      />
      <PopoverContent align="end" className="w-72">
        <div className="space-y-3">
          <h4 className="text-sm font-medium">Members</h4>
          <div className="space-y-1.5">
            {members.map((m) => (
              <div key={`${m.member_type}:${m.member_id}`} className="flex items-center gap-2 text-sm">
                <Avatar className="h-6 w-6">
                  <AvatarFallback className={cn("text-[10px]", m.member_type === "agent" && "bg-purple-100 text-purple-700")}>
                    {m.member_type === "agent" ? <Bot className="h-3 w-3" /> : getActorInitials(m.member_type, m.member_id)}
                  </AvatarFallback>
                </Avatar>
                <span className="truncate">{getActorName(m.member_type, m.member_id)}</span>
                {m.role === "owner" && (
                  <span className="ml-auto text-xs text-muted-foreground">owner</span>
                )}
              </div>
            ))}
          </div>

          {(availableAgents.length > 0 || availableMembers.length > 0) && (
            <div className="border-t pt-3 space-y-2">
              <h4 className="text-xs font-medium text-muted-foreground">Add member</h4>
              <div className="flex gap-2">
                <Select value={addId} onValueChange={(v) => {
                  if (!v) return;
                  const parts = v.split(":");
                  setAddType(parts[0] as "agent" | "member");
                  setAddId(parts[1] ?? "");
                }}>
                  <SelectTrigger className="h-8 text-xs">
                    <SelectValue placeholder="Select..." />
                  </SelectTrigger>
                  <SelectContent>
                    {availableAgents.map((a) => (
                      <SelectItem key={a.id} value={`agent:${a.id}`}>
                        {a.name} (Agent)
                      </SelectItem>
                    ))}
                    {availableMembers.map((m) => (
                      <SelectItem key={m.user_id} value={`member:${m.user_id}`}>
                        {m.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <Button size="sm" className="h-8" onClick={handleAdd} disabled={!addId}>
                  <UserPlus className="h-3.5 w-3.5" />
                </Button>
              </div>
            </div>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
