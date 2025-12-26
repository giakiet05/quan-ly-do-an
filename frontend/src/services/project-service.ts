import { apiFetch } from "./api";
import type { Project, ProjectRound } from "../models";

// ==================== PROJECT SERVICES ====================

/**
 * Get all projects for a classroom
 */
export async function getClassroomProjects(
  classroomId: string
): Promise<Project[]> {
  return apiFetch<Project[]>(`/api/classrooms/${classroomId}/projects`);
}

/**
 * Get a single project by ID
 */
export async function getProject(projectId: string): Promise<Project> {
  return apiFetch<Project>(`/api/projects/${projectId}`);
}

/**
 * Create a new project
 */
export async function createProject(
  classroomId: string,
  projectData: Partial<Project>
): Promise<Project> {
  return apiFetch<Project>(`/api/classrooms/${classroomId}/projects`, {
    method: "POST",
    body: JSON.stringify(projectData),
  });
}

/**
 * Update a project
 */
export async function updateProject(
  projectId: string,
  updates: Partial<Project>
): Promise<Project> {
  return apiFetch<Project>(`/api/projects/${projectId}`, {
    method: "PUT",
    body: JSON.stringify(updates),
  });
}

/**
 * Delete a project
 */
export async function deleteProject(projectId: string): Promise<void> {
  return apiFetch<void>(`/api/projects/${projectId}`, {
    method: "DELETE",
  });
}

// ==================== PROJECT ROUND SERVICES ====================

/**
 * Get all project rounds for a classroom
 */
export async function getProjectRounds(
  classroomId: string
): Promise<ProjectRound[]> {
  return apiFetch<ProjectRound[]>(`/api/classrooms/${classroomId}/rounds`);
}

/**
 * Get a single project round
 */
export async function getProjectRound(
  roundId: string
): Promise<ProjectRound> {
  return apiFetch<ProjectRound>(`/api/rounds/${roundId}`);
}

/**
 * Create a new project round
 */
export async function createProjectRound(
  classroomId: string,
  roundData: Partial<ProjectRound>
): Promise<ProjectRound> {
  return apiFetch<ProjectRound>(`/api/classrooms/${classroomId}/rounds`, {
    method: "POST",
    body: JSON.stringify(roundData),
  });
}

/**
 * Update a project round
 */
export async function updateProjectRound(
  roundId: string,
  updates: Partial<ProjectRound>
): Promise<ProjectRound> {
  return apiFetch<ProjectRound>(`/api/rounds/${roundId}`, {
    method: "PUT",
    body: JSON.stringify(updates),
  });
}

/**
 * Delete a project round
 */
export async function deleteProjectRound(roundId: string): Promise<void> {
  return apiFetch<void>(`/api/rounds/${roundId}`, {
    method: "DELETE",
  });
}
