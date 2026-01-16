// src/stores/notification-store.ts
import { writable } from "svelte/store";
import type { NotificationItem } from "../types/notification";
import { getMyNotifications, markAllNotificationsAsRead } from "../services/notification-service";
import { mapNotifications } from "../mappers/notification-mapper";

function createNotificationStore() {
    const notifications = writable<NotificationItem[]>([]);
    const loading = writable(false);

    async function fetchNotifications() {
        loading.set(true);
        try {
            const res = await getMyNotifications();
            console.log("Dữ liệu thô từ API:", res);
            notifications.set(mapNotifications(res.notifications));
        } finally {
            loading.set(false);
        }
    }

    async function markAllAsRead() {
        await markAllNotificationsAsRead();
        notifications.update(list =>
            list.map(n => ({ ...n, read: true }))
        );
    }

    async function markAsRead(id: string) {
        console.log("BE chưa có API, đang giả lập đọc cho ID:", id);

        notifications.update(list =>
            list.map(n => n.id === id ? { ...n, read: true } : n)
        );

        /* Khi nào BE có API thì mở cái này:
        try {
            await markNotificationAsRead(id);
        } catch (err) {
            // Nếu lỗi thì hoàn tác (rollback) trạng thái cũ ở đây
        }
        */
    }

    async function remove(id: string) {
        console.log("BE chưa có API, đang giả lập xóa ID:", id);

        notifications.update(list => list.filter(n => n.id !== id));
        /* Khi nào BE xong thì thêm:
        await deleteNotification(id);
        */
    }

    return {
        notifications,
        loading,
        fetchNotifications,
        markAllAsRead,
        markAsRead,
        remove,
    };
}

export const notificationStore = createNotificationStore();
