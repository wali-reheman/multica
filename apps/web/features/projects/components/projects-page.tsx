"use client";

import { useEffect, useState } from "react";
import { FolderGit2, Plus, GitBranch } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { useProjectStore } from "../store";
import { ProjectList } from "./project-list";
import { ProjectDetail } from "./project-detail";
import { AddProjectDialog } from "./add-project-dialog";

export function ProjectsPage() {
  const projects = useProjectStore((s) => s.projects);
  const loading = useProjectStore((s) => s.loading);
  const activeProjectId = useProjectStore((s) => s.activeProjectId);
  const fetch = useProjectStore((s) => s.fetch);

  const [addOpen, setAddOpen] = useState(false);

  useEffect(() => {
    fetch();
  }, [fetch]);

  const activeProject = projects.find((p) => p.id === activeProjectId);

  if (loading) {
    return (
      <div className="flex h-full flex-col gap-4 p-6">
        <Skeleton className="h-8 w-48" />
        <div className="flex gap-4">
          <Skeleton className="h-64 w-72" />
          <Skeleton className="h-64 flex-1" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      {/* Header */}
      <div className="flex items-center justify-between border-b px-6 py-3">
        <div className="flex items-center gap-2">
          <FolderGit2 className="size-4 text-muted-foreground" />
          <h1 className="text-sm font-medium">Projects</h1>
          <span className="text-xs text-muted-foreground">
            {projects.length} project{projects.length !== 1 ? "s" : ""}
          </span>
        </div>
        <Button size="sm" variant="outline" onClick={() => setAddOpen(true)}>
          <Plus className="mr-1 size-3.5" />
          Add Project
        </Button>
      </div>

      {/* Body */}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar project list */}
        <div className="w-72 shrink-0 overflow-y-auto border-r">
          <ProjectList
            projects={projects}
            activeId={activeProjectId}
            onSelect={(id) => useProjectStore.getState().setActiveProject(id)}
          />
        </div>

        {/* Main content */}
        <div className="flex-1 overflow-y-auto">
          {activeProject ? (
            <ProjectDetail project={activeProject} />
          ) : (
            <div className="flex h-full flex-col items-center justify-center gap-3 text-muted-foreground">
              <GitBranch className="size-10 opacity-40" />
              <p className="text-sm">Select a project to view version history</p>
            </div>
          )}
        </div>
      </div>

      <AddProjectDialog open={addOpen} onOpenChange={setAddOpen} />
    </div>
  );
}
