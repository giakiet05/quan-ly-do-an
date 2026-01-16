export type NotificationUIType =
    | "message"
    | "deadline"
    | "submission"
    | "class"
    | "system";

export interface NotificationItem {
    id: string;
    ui_type: NotificationUIType;
    title: string;
    content: string;
    link?: string;
    read: boolean;
    createdAt: string;
    meta?: string;
}