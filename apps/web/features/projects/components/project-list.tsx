"use client";

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
    </div>
  );
}
