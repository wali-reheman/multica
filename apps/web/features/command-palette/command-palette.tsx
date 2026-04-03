"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import {
  Inbox,
  ListTodo,
  Bot,
  Monitor,
  Settings,
  BookOpenText,
  SquarePen,
  CircleUser,
  Search,
  Sun,
  Moon,
  Laptop,
} from "lucide-react";
import {
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandShortcut,
  CommandSeparator,
} from "@/components/ui/command";
import { useIssueStore } from "@/features/issues/store";
import { useWorkspaceStore } from "@/features/workspace";
import { useModalStore } from "@/features/modals";
import { StatusIcon } from "@/features/issues/components/status-icon";
import { PriorityIcon } from "@/features/issues/components/priority-icon";
import { useTheme } from "next-themes";

export function CommandPalette() {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const router = useRouter();
  const issues = useIssueStore((s) => s.issues);
  const workspace = useWorkspaceStore((s) => s.workspace);
  const { setTheme } = useTheme();

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((o) => !o);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const runCommand = useCallback(
    (command: () => void) => {
      setOpen(false);
      setSearch("");
      command();
    },
    [],
  );

  const filteredIssues = useMemo(() => {
    if (!search) return issues.slice(0, 10);
    const q = search.toLowerCase();
    return issues
      .filter(
        (i) =>
          i.title.toLowerCase().includes(q) ||
          i.identifier.toLowerCase().includes(q),
      )
      .slice(0, 10);
  }, [issues, search]);

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput
        placeholder="Search issues, navigate, or run a command..."
        value={search}
        onValueChange={setSearch}
      />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>

        {/* Issues */}
        {filteredIssues.length > 0 && (
          <CommandGroup heading="Issues">
            {filteredIssues.map((issue) => (
              <CommandItem
                key={issue.id}
                value={`${issue.identifier} ${issue.title}`}
                onSelect={() =>
                  runCommand(() => router.push(`/issues/${issue.id}`))
                }
              >
                <StatusIcon status={issue.status} className="shrink-0" />
                <span className="text-muted-foreground text-xs shrink-0">
                  {issue.identifier}
                </span>
                <span className="truncate">{issue.title}</span>
                <PriorityIcon
                  priority={issue.priority}
                  className="ml-auto shrink-0"
                />
              </CommandItem>
            ))}
          </CommandGroup>
        )}

        <CommandSeparator />

        {/* Navigation */}
        <CommandGroup heading="Navigation">
          <CommandItem
            value="Inbox"
            onSelect={() => runCommand(() => router.push("/inbox"))}
          >
            <Inbox className="text-muted-foreground" />
            <span>Inbox</span>
          </CommandItem>
          <CommandItem
            value="My Issues"
            onSelect={() => runCommand(() => router.push("/my-issues"))}
          >
            <CircleUser className="text-muted-foreground" />
            <span>My Issues</span>
          </CommandItem>
          <CommandItem
            value="Issues"
            onSelect={() => runCommand(() => router.push("/issues"))}
          >
            <ListTodo className="text-muted-foreground" />
            <span>Issues</span>
          </CommandItem>
          <CommandItem
            value="Agents"
            onSelect={() => runCommand(() => router.push("/agents"))}
          >
            <Bot className="text-muted-foreground" />
            <span>Agents</span>
          </CommandItem>
          <CommandItem
            value="Runtimes"
            onSelect={() => runCommand(() => router.push("/runtimes"))}
          >
            <Monitor className="text-muted-foreground" />
            <span>Runtimes</span>
          </CommandItem>
          <CommandItem
            value="Skills"
            onSelect={() => runCommand(() => router.push("/skills"))}
          >
            <BookOpenText className="text-muted-foreground" />
            <span>Skills</span>
          </CommandItem>
          <CommandItem
            value="Settings"
            onSelect={() => runCommand(() => router.push("/settings"))}
          >
            <Settings className="text-muted-foreground" />
            <span>Settings</span>
          </CommandItem>
        </CommandGroup>

        <CommandSeparator />

        {/* Actions */}
        <CommandGroup heading="Actions">
          <CommandItem
            value="Create new issue"
            onSelect={() =>
              runCommand(() => useModalStore.getState().open("create-issue"))
            }
          >
            <SquarePen className="text-muted-foreground" />
            <span>Create new issue</span>
            <CommandShortcut>C</CommandShortcut>
          </CommandItem>
          <CommandItem
            value="Search issues"
            onSelect={() => runCommand(() => router.push("/issues"))}
          >
            <Search className="text-muted-foreground" />
            <span>Search issues</span>
          </CommandItem>
        </CommandGroup>

        <CommandSeparator />

        {/* Theme */}
        <CommandGroup heading="Theme">
          <CommandItem
            value="Light theme"
            onSelect={() => runCommand(() => setTheme("light"))}
          >
            <Sun className="text-muted-foreground" />
            <span>Light</span>
          </CommandItem>
          <CommandItem
            value="Dark theme"
            onSelect={() => runCommand(() => setTheme("dark"))}
          >
            <Moon className="text-muted-foreground" />
            <span>Dark</span>
          </CommandItem>
          <CommandItem
            value="System theme"
            onSelect={() => runCommand(() => setTheme("system"))}
          >
            <Laptop className="text-muted-foreground" />
            <span>System</span>
          </CommandItem>
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
