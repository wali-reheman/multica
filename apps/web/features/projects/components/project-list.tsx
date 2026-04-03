"use client";

<<<<<<< HEAD
import { FolderGit2, Folder } from "lucide-react";
import type { Project } from "@/shared/types";
import { cn } from "@/lib/utils";

interface ProjectListProps {
  projects: Project[];
  activeId: string | null;
  onSelect: (id: string) => void;
}

export function ProjectList({ projects, activeId, onSelect }: ProjectListProps) {
  if (projects.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 p-6 text-muted-foreground">
        <Folder className="size-8 opacity-40" />
        <p className="text-xs">No projects yet</p>
        <p className="text-xs">Add a local folder to get started</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col py-1">
      {projects.map((project) => (
        <button
          key={project.id}
          type="button"
          onClick={() => onSelect(project.id)}
          className={cn(
            "flex items-start gap-3 px-4 py-2.5 text-left hover:bg-accent/50 transition-colors",
            activeId === project.id && "bg-accent"
          )}
        >
          <div className="mt-0.5">
            {project.is_git_repo ? (
              <FolderGit2 className="size-4 text-muted-foreground" />
            ) : (
              <Folder className="size-4 text-muted-foreground" />
            )}
          </div>
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-medium">{project.name}</p>
            <p className="truncate text-xs text-muted-foreground">
              {project.local_path}
            </p>
            {project.language && (
              <span className="mt-1 inline-block rounded bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">
                {project.language}
              </span>
            )}
          </div>
        </button>
      ))}
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
    </div>
  );
}
