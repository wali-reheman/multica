"use client";

<<<<<<< HEAD
import { useState } from "react";
import { toast } from "sonner";
=======
import { useState, useCallback } from "react";
import { Loader2 } from "lucide-react";
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
<<<<<<< HEAD
=======
  DialogDescription,
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { api } from "@/shared/api";
<<<<<<< HEAD
import { useProjectStore } from "../store";

interface AddProjectDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function AddProjectDialog({ open, onOpenChange }: AddProjectDialogProps) {
  const [localPath, setLocalPath] = useState("");
  const [name, setName] = useState("");
  const [initGit, setInitGit] = useState(true);
  const [loading, setLoading] = useState(false);
  const addProject = useProjectStore((s) => s.addProject);
  const setActiveProject = useProjectStore((s) => s.setActiveProject);

  const handleSubmit = async () => {
    if (!localPath.trim()) return;
    setLoading(true);
    try {
      const project = await api.createProject({
        local_path: localPath.trim(),
        name: name.trim() || undefined,
        init_git: initGit,
      });
      addProject(project);
      setActiveProject(project.id);
      toast.success(`Project "${project.name}" added`);
      setLocalPath("");
      setName("");
      onOpenChange(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add project");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Project</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="local-path" className="text-xs">
              Local Path
            </Label>
            <Input
              id="local-path"
              placeholder="/path/to/your/project"
=======
import { toast } from "sonner";
import { useProjectStore } from "../store";

export function AddProjectDialog({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const addProject = useProjectStore((s) => s.addProject);
  const setActiveProject = useProjectStore((s) => s.setActiveProject);
  const [name, setName] = useState("");
  const [localPath, setLocalPath] = useState("");
  const [initGit, setInitGit] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!localPath.trim()) return;

      setSubmitting(true);
      try {
        const project = await api.createProject({
          name: name.trim() || undefined,
          local_path: localPath.trim(),
          init_git: initGit,
        });
        addProject(project);
        setActiveProject(project.id);
        toast.success("Project added");
        setName("");
        setLocalPath("");
        setInitGit(false);
        onOpenChange(false);
      } catch {
        toast.error("Failed to add project");
      } finally {
        setSubmitting(false);
      }
    },
    [name, localPath, initGit, addProject, setActiveProject, onOpenChange],
  );

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Project</DialogTitle>
          <DialogDescription>
            Register a local directory as a project
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="project-path">Local path</Label>
            <Input
              id="project-path"
              placeholder="/path/to/project"
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
              value={localPath}
              onChange={(e) => setLocalPath(e.target.value)}
              autoFocus
            />
          </div>
<<<<<<< HEAD
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="project-name" className="text-xs">
              Name (optional)
            </Label>
            <Input
              id="project-name"
              placeholder="Auto-detected from folder name"
=======
          <div className="space-y-2">
            <Label htmlFor="project-name">Name (optional)</Label>
            <Input
              id="project-name"
              placeholder="Derived from directory name if empty"
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
<<<<<<< HEAD
          <div className="flex items-center gap-2">
            <Checkbox
              id="init-git"
              checked={initGit}
              onCheckedChange={(checked) => setInitGit(!!checked)}
            />
            <Label htmlFor="init-git" className="text-xs">
              Initialize git if not a repository
            </Label>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={!localPath.trim() || loading}>
            {loading ? "Adding..." : "Add Project"}
          </Button>
        </DialogFooter>
=======
          <label className="flex items-center gap-2 cursor-pointer">
            <Checkbox
              checked={initGit}
              onCheckedChange={() => setInitGit(!initGit)}
            />
            <span className="text-sm">Initialize git repository</span>
          </label>
          <DialogFooter>
            <Button
              type="submit"
              disabled={!localPath.trim() || submitting}
            >
              {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
              Add Project
            </Button>
          </DialogFooter>
        </form>
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
      </DialogContent>
    </Dialog>
  );
}
