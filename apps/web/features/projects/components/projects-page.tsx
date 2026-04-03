"use client";

<<<<<<< HEAD
import { useEffect, useState } from "react";
import { FolderGit2, Plus, GitBranch } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
import { useProjectStore } from "../store";
import { ProjectList } from "./project-list";
import { ProjectDetail } from "./project-detail";
import { AddProjectDialog } from "./add-project-dialog";

export function ProjectsPage() {
<<<<<<< HEAD
=======
  const isLoading = useAuthStore((s) => s.isLoading);
  const workspace = useWorkspaceStore((s) => s.workspace);
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
  const projects = useProjectStore((s) => s.projects);
  const loading = useProjectStore((s) => s.loading);
  const activeProjectId = useProjectStore((s) => s.activeProjectId);
  const fetch = useProjectStore((s) => s.fetch);
<<<<<<< HEAD

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
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
        </div>
      </div>
    );
  }

  return (
<<<<<<< HEAD
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
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
  );
}
