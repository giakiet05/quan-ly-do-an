// src/services/project-round-service.ts
import { apiFetch } from "../utils/api-fetch";
import type {
    ProjectRoundResponse,
    CreateProjectRoundRequest,
    UpdateProjectRoundRequest
} from "../dtos/project-dto";
import type { ProjectRound } from "../types/project-round";
import { projectRoundMapper } from "../mappers/project-round-mapper";

export async function getProjectRounds(classroomId: string): Promise<ProjectRound[]> {
    const response = await apiFetch<ProjectRoundResponse[]>(
        `/api/projects/classrooms/${classroomId}/rounds`
    );

    if (!response) return [];

    return response.map(dto => projectRoundMapper.toEntity(dto));
}

/**
 * Lấy chi tiết một vòng dự án
 */
export async function getProjectRoundById(roundId: string): Promise<ProjectRound | null> {
    const response = await apiFetch<ProjectRoundResponse>(`/api/projects/classrooms/:classroom_id/rounds/${roundId}`);
    if (!response) return null;
    return projectRoundMapper.toEntity(response);
}

/**
 * Tạo mới vòng dự án
 */
export async function createProjectRound(data: CreateProjectRoundRequest): Promise<ProjectRound> {
    const response = await apiFetch<ProjectRoundResponse>("/api/projects/rounds", {
        method: "POST",
        body: JSON.stringify(data)
    });

    return projectRoundMapper.toEntity(response);
}

/**
 * Cập nhật vòng dự án
 */
export async function updateProjectRound(
    data: UpdateProjectRoundRequest
): Promise<ProjectRound> {
    const response = await apiFetch<ProjectRoundResponse>(`/api/projects/rounds/`, {
        method: "PUT",
        body: JSON.stringify(data)
    });

    return projectRoundMapper.toEntity(response);
}

/**
 * Xóa vòng dự án (Soft delete)
 */
export async function deleteProjectRound(roundId: string, classroomId: string): Promise<void> {
    await apiFetch(`/api/projects/classrooms/${classroomId}/rounds/${roundId}`, {
        method: "DELETE"
    });
}