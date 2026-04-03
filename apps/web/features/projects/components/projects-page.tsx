"use client";

import { useEffect, useCallback, useState } from "react";
import { FolderGit2 } from "lucide-react";
import { useDefaultLayout } from "react-resizable-panels";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/features/auth";
import { useWorkspaceStore } from "@/features/workspace";
import { useWSEvent } from "@/features/realtime";
import type { Project } from "@/shared/types";
import { useProjectStore } from "../store";
import { ProjectList } from "./project-list";
import { ProjectDetail } from "./project-detail";
import { AddProjectDialog } from "./add-project-dialog";

export function ProjectsPage() {
  const isLoading = useAuthStore((s) => s.isLoading);
  const workspace = useWorkspaceStore((s) => s.workspace);
  const projects = useProjectStore((s) => s.projects);
  const loading = useProjectStore((s) => s.loading);
  const activeProjectId = useProjectStore((s) => s.activeProjectId);
  const fetch = useProjectStore((s) => s.fetch);
  const setActiveProject = useProjectStore((s) => s.setActiveProject);
  const addProject = useProjectStore((s) => s.addProject);
  const updateProject = useProjectStore((s) => s.updateProject);
  const removeProject = useProjectStore((s) => s.removeProject);
  const [addOpen, setAddOpen] = useState(false);

  const { defaultLayout, onLayoutChanged } = useDefaultLayout({
    id: "multica_projects_layout",
  });

  useEffect(() => {
    if (workspace) fetch();
  }, [workspace, fetch]);

  const handleProjectCreated = useCallback(
    (payload: unknown) => {
      const data = payload as { project: Project };
      addProject(data.project);
    },
    [addProject],
  );

  const handleProjectUpdated = useCallback(
    (payload: unknown) => {
      const data = payload as { project: Project };
      updateProject(data.project.id, data.project);
    },
    [updateProject],
  );

  const handleProjectDeleted = useCallback(
    (payload: unknown) => {
      const data = payload as { project_id: string };
      removeProject(data.project_id);
      if (activeProjectId === data.project_id) {
        setActiveProject(null);
      }
    },
    [removeProject, activeProjectId, setActiveProject],
  );

  useWSEvent("project:created", handleProjectCreated);
  useWSEvent("project:updated", handleProjectUpdated);
  useWSEvent("project:deleted", handleProjectDeleted);

  const activeProject = projects.find((p) => p.id === activeProjectId) ?? null;

  if (isLoading || loading) {
    return (
      <div className="flex flex-1 min-h-0">
        <div className="w-72 border-r">
          <div className="flex h-12 items-center justify-between border-b px-4">
            <Skeleton className="h-4 w-20" />
          </div>
          <div className="divide-y">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3 px-4 py-3">
                <Skeleton className="h-5 w-5 rounded" />
                <div className="flex-1 space-y-1.5">
                  <Skeleton className="h-4 w-28" />
                  <Skeleton className="h-3 w-20" />
                </div>
              </div>
            ))}
          </div>
        </div>
        <div className="flex-1 p-6 space-y-6">
          <Skeleton className="h-5 w-32" />
          <div className="space-y-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-16 w-full rounded-lg" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <>
      <ResizablePanelGroup
        orientation="horizontal"
        className="flex-1 min-h-0"
        defaultLayout={defaultLayout}
        onLayoutChanged={onLayoutChanged}
      >
        <ResizablePanel
          id="list"
          defaultSize={280}
          minSize={240}
          maxSize={400}
          groupResizeBehavior="preserve-pixel-size"
        >
          <ProjectList
            projects={projects}
            selectedId={activeProjectId}
            onSelect={setActiveProject}
            onAddClick={() => setAddOpen(true)}
          />
        </ResizablePanel>

        <ResizableHandle />

        <ResizablePanel id="detail" minSize="50%">
          {activeProject ? (
            <ProjectDetail key={activeProject.id} project={activeProject} />
          ) : (
            <div className="flex h-full flex-col items-center justify-center text-muted-foreground">
              <FolderGit2 className="h-10 w-10 text-muted-foreground/30" />
              <p className="mt-3 text-sm">Select a project to view details</p>
              <Button
                variant="outline"
                size="sm"
                className="mt-4"
                onClick={() => setAddOpen(true)}
              >
                Add Project
              </Button>
            </div>
          )}
        </ResizablePanel>
      </ResizablePanelGroup>

      <AddProjectDialog open={addOpen} onOpenChange={setAddOpen} />
    </>
  );
}
