import { writable, derived } from "svelte/store";
import type { NotificationItem } from "../types/notification";
import {
    getMyNotifications,
    getClassroomNotifications,
    markAsRead as markAsReadApi,
    markAllNotificationsAsRead
} from "../services/notification-service";
import { mapNotifications } from "../mappers/notification-mapper";

function createNotificationStore() {
    const notifications = writable<NotificationItem[]>([]);
    const loading = writable(false);
    const totalItems = writable(0);
    // BIẾN MỚI: Giữ số lượng chưa đọc toàn cục (không bị ảnh hưởng khi chuyển tab lớp)
    const unreadCount = writable(0);

    // Hàm lấy số badge (Gắn vào Sidebar)
    async function loadUnreadCount() {
        try {
            const res = await getMyNotifications(1, 100);
            const count = res.notifications.filter(n => !n.isRead).length;
            unreadCount.set(count);
        } catch (e) { console.error(e); }
    }

    async function fetchNotifications(page = 1, limit = 20) {
        loading.set(true);
        try {
            const res = await getMyNotifications(page, limit);
            notifications.set(mapNotifications(res.notifications));
            totalItems.set(res.pagination.totalItems);
            // Cập nhật badge count khi ở tab thông báo chung
            unreadCount.set(res.notifications.filter(n => !n.isRead).length);
        } catch (error) {
            console.error("Fetch notifications failed", error);
        } finally {
            loading.set(false);
        }
    }

    async function fetchClassroomNotifications(classroomId: string, page = 1, limit = 10) {
        loading.set(true);
        try {
            const res = await getClassroomNotifications(classroomId, page, limit);
            notifications.set(mapNotifications(res.notifications));
            totalItems.set(res.pagination.totalItems);
            // KHÔNG set unreadCount ở đây để tránh bị mất số ở Sidebar
        } catch (error) {
            console.error("Fetch classroom notifications failed", error);
        } finally {
            loading.set(false);
        }
    }

    // GIỮ NGUYÊN HÀM CŨ CỦA BẠN VÀ THÊM LOGIC UPDATE BADGE
    async function markAllAsRead() {
        try {
            await markAllNotificationsAsRead();
            notifications.update(list =>
                list.map(n => ({ ...n, read: true }))
            );
            // Sau khi đọc hết, set số badge về 0 ngay lập tức
            unreadCount.set(0);
        } catch (error) {
            console.error("Mark all as read failed", error);
        }
    }

    async function markAsRead(id: string) {
        notifications.update(list =>
            list.map(n => {
                if (n.id === id && !n.read) {
                    // Trừ 1 vào số badge toàn cục ngay trên UI
                    unreadCount.update(c => Math.max(0, c - 1));
                    return { ...n, read: true };
                }
                return n;
            })
        );
        try {
            await markAsReadApi(id);
        } catch (err) {
            console.error("Mark as read API failed", err);
        }
    }

    async function remove(id: string) {
        notifications.update(list => list.filter(n => n.id !== id));
    }

    return {
        notifications,
        loading,
        totalItems,
        unreadCount, // Biến này để Sidebar lắng nghe
        loadUnreadCount,
        fetchNotifications,
        fetchClassroomNotifications,
        markAllAsRead,
        markAsRead,
        remove,
    };
}

export const notificationStore = createNotificationStore();

// SỬA LẠI: Lấy số từ biến unreadCount đã tách biệt
export const unreadNotificationCount = derived(
    notificationStore.unreadCount,
    ($count) => $count
);