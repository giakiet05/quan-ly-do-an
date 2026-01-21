import { apiFetch } from "../utils/api-fetch";
import type {
    CreateProjectRequest,
    CreateProjectsRequest,
    UpdateProjectRequest,
} from "../dtos/project-dto";
import type { ApiResponse, PaginatedResponse } from "../dtos/api-response-dto";
import type { ProjectResponse } from "../dtos/project-dto";

// ==================== PROJECT MANAGEMENT ====================

/**
 * Get all projects in a classroom and round
 */
export async function getProjects(
    classroomId: string,
    projectRoundId: string
): Promise<ProjectResponse[]> {
    const response = await apiFetch<ProjectResponse[]>(
        `/api/projects/classrooms/${classroomId}/rounds/${projectRoundId}/projects`
    );
    return response || [];
}

/**
 * Get a specific project by ID
 */
export async function getProject(
    classroomId: string,
    projectId: string
): Promise<ProjectResponse> {
    const response = await apiFetch<ApiResponse<ProjectResponse>>(
        `/api/projects/classrooms/${classroomId}/${projectId}`
    );
    return response.data!;
}

/**
 * Create a new project
 */
export async function createProject(
    data: CreateProjectRequest
): Promise<ProjectResponse> {
    const response = await apiFetch<ProjectResponse>(
        "/api/projects",
        {
            method: "POST",
            body: JSON.stringify(data),
        }
    );
    console.log("API response:", response);
    return response;
}

/**
 * Create multiple projects
 */
export async function createProjects(
    data: CreateProjectsRequest
): Promise<{ created_count: number }> {
    const response = await apiFetch<ApiResponse<{ created_count: number }>>(
        "/api/projects/bulk",
        {
            method: "POST",
            body: JSON.stringify(data),
        }
    );
    return response.data!;
}

/**
 * Update a project
 */
export async function updateProject(
    data: UpdateProjectRequest
): Promise<ProjectResponse> {
    const response = await apiFetch<ApiResponse<ProjectResponse>>(
        "/api/projects",
        {
            method: "PUT",
            body: JSON.stringify(data),
        }
    );
    return response.data!;
}

/**
 * Delete a project
 */
export async function deleteProject(
    classroomId: string,
    projectId: string
): Promise<void> {
    await apiFetch(`/api/projects/classrooms/${classroomId}/${projectId}`, {
        method: "DELETE",
    });
}
