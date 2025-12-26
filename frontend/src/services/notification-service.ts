import { apiFetch } from "./api";
import type { Notification } from "../models";

// ==================== NOTIFICATION SERVICES ====================

/**
 * Get all notifications for current user
 */
export async function getNotifications(options?: {
  limit?: number;
  offset?: number;
  unreadOnly?: boolean;
}): Promise<{
  notifications: Notification[];
  total: number;
  unreadCount: number;
}> {
  const params = new URLSearchParams();
  if (options?.limit) params.append("limit", options.limit.toString());
  if (options?.offset) params.append("offset", options.offset.toString());
  if (options?.unreadOnly) params.append("unread_only", "true");

  const query = params.toString() ? `?${params.toString()}` : "";
  return apiFetch(`/api/notifications${query}`);
}

/**
 * Mark a notification as read
 */
export async function markNotificationAsRead(
  notificationId: string
): Promise<void> {
  return apiFetch<void>(`/api/notifications/${notificationId}/read`, {
    method: "POST",
  });
}

/**
 * Mark all notifications as read
 */
export async function markAllNotificationsAsRead(): Promise<void> {
  return apiFetch<void>(`/api/notifications/read-all`, {
    method: "POST",
  });
}

/**
 * Delete a notification
 */
export async function deleteNotification(notificationId: string): Promise<void> {
  return apiFetch<void>(`/api/notifications/${notificationId}`, {
    method: "DELETE",
  });
}

/**
 * Get unread notification count
 */
export async function getUnreadNotificationCount(): Promise<number> {
  const result = await apiFetch<{ count: number }>(`/api/notifications/unread-count`);
  return result.count;
}
