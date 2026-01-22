import { apiFetch } from "../utils/api-fetch";
import type {
    CreatePeriodRequest,
    CreatePeriodsRequest,
    UpdatePeriodRequest,
} from "../dtos/period-dto";
import type { ApiResponse } from "../dtos/api-response-dto";
import type { PeriodResponse } from "../dtos/period-dto";

// ==================== PERIOD MANAGEMENT ====================

/**
 * Get all periods in a classroom and round
 */
export async function getPeriods(
    classroomId: string,
    roundId: string
): Promise<PeriodResponse[]> {
    const response = await apiFetch<PeriodResponse[]>(
        `/api/projects/classrooms/${classroomId}/rounds/${roundId}/report-periods`
    );
    return response || [];
}

/**
 * Get a specific period by ID
 */
export async function getPeriod(
    classroomId: string,
    periodId: string
): Promise<PeriodResponse> {
    const response = await apiFetch<ApiResponse<PeriodResponse>>(
        `/api/periods/classrooms/${classroomId}/${periodId}`
    );
    return response.data!;
}

/**
 * Create a new period
 */
export async function createPeriod(
    data: CreatePeriodRequest,
    classroomId: string,
    roundId: string
): Promise<PeriodResponse> {
    const response = await apiFetch<PeriodResponse>(
        `/api/projects/classrooms/${classroomId}/rounds/${roundId}/report-periods`,
        {
            method: "POST",
            body: JSON.stringify(data),
        }
    );
    return response;
}

/**
 * Create multiple periods
 */
export async function createPeriods(
    data: CreatePeriodsRequest
): Promise<{ created_count: number }> {
    const response = await apiFetch<ApiResponse<{ created_count: number }>>(
        "/api/periods/bulk",
        {
            method: "POST",
            body: JSON.stringify(data),
        }
    );
    return response.data!;
}

/**
 * Update a period
 */
export async function updatePeriod(
    data: UpdatePeriodRequest
): Promise<PeriodResponse> {
    const response = await apiFetch<PeriodResponse>(
        "/api/periods",
        {
            method: "PUT",
            body: JSON.stringify(data),
        }
    );
    return response;
}

/**
 * Delete a period
 */
export async function deletePeriod(
    classroomId: string,
    roundId: string,
    periodId: string
): Promise<void> {
    await apiFetch(`/api/projects/classrooms/${classroomId}/rounds/${roundId}/report-periods/${periodId}`, {
        method: "DELETE",
    });
}
