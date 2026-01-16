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

    console.log("Dữ liệu thực tế sau apiFetch:", response);

    if (!response) {
        return {
            notifications: [],
            pagination: { page, page_size: limit, total_items: 0, total_pages: 0 }
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
