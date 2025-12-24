import { writable, derived } from "svelte/store";

export type NotificationType =
    | "student_action"
    | "deadline"
    | "submission"
    | "system";

export type Notification = {
    id: string;
    type: NotificationType;
    title: string;
    content?: string;
    meta?: string;
    tag?: string;
    date: string;
    datetime: string;
    read: boolean;
    actionable?: boolean;
};

const initial: Notification[] = [
    {
        id: "n1",
        type: "student_action",
        title: 'Sinh viên đăng ký đề tài',
        content: 'Nguyễn Văn A đã đăng ký đề tài "Hệ thống quản lý thư viện"',
        meta: 'Sinh viên: Nguyễn Văn A · Đề tài: Hệ thống quản lý thư viện',
        tag: "Hành động sinh viên",
        date: "10:30 30/11/2024",
        datetime: "2024-11-30T10:30:00Z",
        read: false,
        actionable: true
    },
    {
        id: "n2",
        type: "deadline",
        title: "Sắp đến hạn nộp báo cáo",
        content: 'Hạng mục "Đề tài Web Application" sẽ đến hạn nộp báo cáo vào 02/12/2024',
        tag: "Nhắc nhở deadline",
        date: "09:15 30/11/2024",
        datetime: "2024-11-30T09:15:00Z",
        read: false,
        actionable: false
    },
    {
        id: "n3",
        type: "submission",
        title: "Báo cáo mới được nộp",
        content: 'Trần Thị B đã nộp báo cáo cho đề tài "Website quản lý sinh viên"',
        tag: "Nộp bài",
        date: "16:45 29/11/2024",
        datetime: "2024-11-29T16:45:00Z",
        read: true,
        actionable: false
    }
];

export const notifications = writable<Notification[]>(initial);

export const unreadCount = derived(notifications, $n =>
    $n.filter((x) => !x.read).length
);

function updateItem(id: string, patch: Partial<Notification>) {
    notifications.update(list => list.map(it => it.id === id ? { ...it, ...patch } : it));
}

export const markAsRead = (id: string) => updateItem(id, { read: true });
export const markAsUnread = (id: string) => updateItem(id, { read: false });
export const markAllAsRead = () => notifications.update(list => list.map(n => ({ ...n, read: true })));
export const removeNotification = (id: string) => notifications.update(list => list.filter(n => n.id !== id));
export const addNotification = (n: Notification) => notifications.update(list => [n, ...list]);

export default {
    subscribe: notifications.subscribe,
    markAsRead,
    markAsUnread,
    markAllAsRead,
    removeNotification,
    addNotification
};
