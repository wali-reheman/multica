"use client";

import { create } from "zustand";
import type { Project, CommitInfo, BranchInfo, GitStatus } from "@/shared/types";
import { toast } from "sonner";
import { api } from "@/shared/api";
import { createLogger } from "@/shared/logger";

const logger = createLogger("project-store");

interface ProjectState {
  projects: Project[];
  loading: boolean;
  activeProjectId: string | null;
  commits: CommitInfo[];
  commitsLoading: boolean;
  branches: BranchInfo[];
  branchesLoading: boolean;
  status: GitStatus | null;
  statusLoading: boolean;

  fetch: () => Promise<void>;
  setProjects: (projects: Project[]) => void;
  addProject: (project: Project) => void;
  updateProject: (id: string, updates: Partial<Project>) => void;
  removeProject: (id: string) => void;
  setActiveProject: (id: string | null) => void;

  fetchCommits: (projectId: string) => Promise<void>;
  fetchBranches: (projectId: string) => Promise<void>;
  fetchStatus: (projectId: string) => Promise<void>;
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projects: [],
  loading: true,
  activeProjectId: null,
  commits: [],
  commitsLoading: false,
  branches: [],
  branchesLoading: false,
  status: null,
  statusLoading: false,

  fetch: async () => {
    logger.debug("fetch start");
    const isInitialLoad = get().projects.length === 0;
    if (isInitialLoad) set({ loading: true });
    try {
      const res = await api.listProjects({ limit: 100 });
      logger.info("fetched", res.projects.length, "projects");
      set({ projects: res.projects, loading: false });
    } catch (err) {
      logger.error("fetch failed", err);
      toast.error("Failed to load projects");
      if (isInitialLoad) set({ loading: false });
    }
  },

  setProjects: (projects) => set({ projects }),
  addProject: (project) =>
    set((s) => ({
      projects: s.projects.some((p) => p.id === project.id)
        ? s.projects
        : [...s.projects, project],
    })),
  updateProject: (id, updates) =>
    set((s) => ({
      projects: s.projects.map((p) => (p.id === id ? { ...p, ...updates } : p)),
    })),
  removeProject: (id) =>
    set((s) => ({ projects: s.projects.filter((p) => p.id !== id) })),
  setActiveProject: (id) => set({ activeProjectId: id }),

  fetchCommits: async (projectId: string) => {
    set({ commitsLoading: true });
    try {
      const res = await api.getProjectCommits(projectId, { limit: 50 });
      set({ commits: res.commits, commitsLoading: false });
    } catch (err) {
      logger.error("fetchCommits failed", err);
      toast.error("Failed to load commits");
      set({ commitsLoading: false });
    }
  },

  fetchBranches: async (projectId: string) => {
    set({ branchesLoading: true });
    try {
      const res = await api.getProjectBranches(projectId);
      set({ branches: res.branches, branchesLoading: false });
    } catch (err) {
      logger.error("fetchBranches failed", err);
      set({ branchesLoading: false });
    }
  },

  fetchStatus: async (projectId: string) => {
    set({ statusLoading: true });
    try {
      const status = await api.getProjectStatus(projectId);
      set({ status, statusLoading: false });
    } catch (err) {
      logger.error("fetchStatus failed", err);
      set({ statusLoading: false });
    }
  },
}));
