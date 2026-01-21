// src/stores/notification-store.ts
import { writable } from "svelte/store";
import type { NotificationItem } from "../types/notification";
import {
    getMyNotifications,
    getClassroomNotifications, // Import hàm mới
    markAsRead as markAsReadApi, // Đổi tên để tránh trùng
    markAllNotificationsAsRead
} from "../services/notification-service";
import { mapNotifications } from "../mappers/notification-mapper";

function createNotificationStore() {
    const notifications = writable<NotificationItem[]>([]);
    const loading = writable(false);
    const totalItems = writable(0);

    // Fetch thông báo cá nhân (Route cũ)
    async function fetchNotifications(page = 1, limit = 20) {
        loading.set(true);
        try {
            const res = await getMyNotifications(page, limit);
            notifications.set(mapNotifications(res.notifications));
            totalItems.set(res.pagination.totalItems);
        } catch (error) {
            console.error("Fetch notifications failed", error);
        } finally {
            loading.set(false);
        }
    }

    // Fetch thông báo của một lớp (Route mới)
    async function fetchClassroomNotifications(classroomId: string, page = 1, limit = 10) {
        loading.set(true);
        try {
            const res = await getClassroomNotifications(classroomId, page, limit);
            notifications.set(mapNotifications(res.notifications));
            totalItems.set(res.pagination.totalItems);
        } catch (error) {
            console.error("Fetch classroom notifications failed", error);
        } finally {
            loading.set(false);
        }
    }

    async function markAllAsRead() {
        try {
            await markAllNotificationsAsRead();
            notifications.update(list =>
                list.map(n => ({ ...n, read: true }))
            );
        } catch (error) {
            console.error("Mark all as read failed", error);
        }
    }

    async function markAsRead(id: string) {
        notifications.update(list =>
            list.map(n => n.id === id ? { ...n, read: true } : n)
        );

        try {
            await markAsReadApi(id);
        } catch (err) {
            console.error("Mark as read API failed", err);
        }
    }

    async function remove(id: string) {
        // Tương tự xóa, nếu BE có API thì gọi ở đây
        notifications.update(list => list.filter(n => n.id !== id));
    }

    return {
        notifications,
        loading,
        totalItems,
        fetchNotifications,
        fetchClassroomNotifications,
        markAllAsRead,
        markAsRead,
        remove,
    };
}

export const notificationStore = createNotificationStore();