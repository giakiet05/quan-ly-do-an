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
    // apiFetch đã unwrap ApiResponse rồi → trả về trực tiếp ProjectResponse
    const project = await apiFetch<ProjectResponse>(
        `/api/projects/classrooms/${classroomId}/${projectId}`
    );

    // Optional: Kiểm tra nếu null hoặc undefined
    if (!project) {
        throw new Error("Không tìm thấy đề tài");
    }

    return project;
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
    const response = await apiFetch<ProjectResponse>(
        "/api/projects",
        {
            method: "PUT",
            body: JSON.stringify(data),
        }
    );
    return response;
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


export async function downloadProjectTemplate(): Promise<void> {
    const response = await apiFetch<Response>("/api/projects/template", {
        method: "GET",
    });

    // apiFetch trả về Response object khi content-type không phải json
    // Ta cần ép kiểu hoặc gọi .blob() từ kết quả trả về
    const blob = await response.blob();

    // Tạo link tải xuống
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "project_import_template.xlsx";
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
}

// 2. UPLOAD EXCEL
// Sử dụng apiFetch thay vì fetch thường để được tự động Refresh Token nếu hết hạn
export async function uploadProjectsExcel(
    classroomId: string,
    projectRoundId: string,
    file: File
): Promise<any> {
    const formData = new FormData();
    formData.append("file", file);

    // Không cần set Header thủ công, apiFetch sẽ tự nhận diện FormData
    // và bỏ qua việc set 'Content-Type': 'application/json'
    const response = await apiFetch(
        `/api/projects/classrooms/${classroomId}/rounds/${projectRoundId}/upload-excel`,
        {
            method: "POST",
            body: formData,
        }
    );

    return response;
}