"use client";

import { Folder, FolderGit2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import type { Project } from "@/shared/types";
import { cn } from "@/lib/utils";

function ProjectListItem({
  project,
  isSelected,
  onClick,
}: {
  project: Project;
  isSelected: boolean;
  onClick: () => void;
}) {
  const Icon = project.is_git_repo ? FolderGit2 : Folder;

  return (
    <button
      onClick={onClick}
      className={cn(
        "flex w-full items-center gap-3 px-4 py-3 text-left transition-colors",
        isSelected ? "bg-accent" : "hover:bg-accent/50",
      )}
    >
      <div
        className={cn(
          "flex h-8 w-8 shrink-0 items-center justify-center rounded-lg",
          project.is_git_repo ? "bg-primary/10" : "bg-muted",
        )}
      >
        <Icon className="h-4 w-4 text-muted-foreground" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="truncate text-sm font-medium">{project.name}</div>
        <div className="mt-0.5 truncate text-xs text-muted-foreground">
          {project.language ?? "Unknown"} &middot;{" "}
          {project.file_count.toLocaleString()} files
        </div>
      </div>
    </button>
  );
}

export function ProjectList({
  projects,
  selectedId,
  onSelect,
  onAddClick,
}: {
  projects: Project[];
  selectedId: string | null;
  onSelect: (id: string) => void;
  onAddClick: () => void;
}) {
  return (
    <div className="overflow-y-auto h-full border-r">
      <div className="flex h-12 items-center justify-between border-b px-4">
        <h1 className="text-sm font-semibold">Projects</h1>
        <Button variant="ghost" size="icon-xs" onClick={onAddClick}>
          <Plus className="h-4 w-4" />
        </Button>
      </div>
      {projects.length === 0 ? (
        <div className="flex flex-col items-center justify-center px-4 py-12">
          <FolderGit2 className="h-8 w-8 text-muted-foreground/40" />
          <p className="mt-3 text-sm text-muted-foreground">No projects yet</p>
          <Button
            variant="outline"
            size="sm"
            className="mt-3"
            onClick={onAddClick}
          >
            Add Project
          </Button>
        </div>
      ) : (
        <div className="divide-y">
          {projects.map((project) => (
            <ProjectListItem
              key={project.id}
              project={project}
              isSelected={project.id === selectedId}
              onClick={() => onSelect(project.id)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
