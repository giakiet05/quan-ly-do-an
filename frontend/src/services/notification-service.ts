import { apiFetch } from "../utils/api-fetch";
import type {
    NotificationResponse,
    PaginatedNotificationsResponse,
} from "../dtos/notification-dto";
import type { ApiResponse } from "../dtos/api-response-dto";

// src/services/notification-service.ts
export async function getMyNotifications(
    page = 1,
    limit = 20
): Promise<PaginatedNotificationsResponse> {
    const response = await apiFetch<PaginatedNotificationsResponse>(
        `/api/notifications?page=${page}&pageSize=${limit}`
    );


    if (!response) {
        return {
            notifications: [],
            pagination: { page, pageSize: limit, totalItems: 0, totalPages: 0 }
        };
    }

    return response;
}

export async function markNotificationAsRead(
    notificationId: string
): Promise<void> {
    await apiFetch(`/api/notifications/${notificationId}/read`, {
        method: "PATCH",
    });
}

export async function markAllNotificationsAsRead(): Promise<void> {
    await apiFetch("/api/notifications/read-all", {
        method: "PUT",
    });
}

// Lấy thông báo riêng cho một lớp học
export async function getClassroomNotifications(
    classroomId: string,
    page = 1,
    limit = 10
): Promise<PaginatedNotificationsResponse> {
    const response = await apiFetch<PaginatedNotificationsResponse>(
        `/api/notifications/classroom/${classroomId}?page=${page}&pageSize=${limit}`
    );

    if (!response) {
        return {
            notifications: [],
            pagination: { page, pageSize: limit, totalItems: 0, totalPages: 0 }
        };
    }

    return response;
}

// Mark một thông báo là đã đọc (Dùng chung cho cả thông báo hệ thống và lớp học)
export async function markAsRead(
    notificationId: string
): Promise<void> {
    await apiFetch(`/api/notifications/${notificationId}/read`, {
        method: "POST",
    });
}
