export type NotificationType =
    | "new_message"
    | "project_registration_deadline"
    | "project_registration_expired"
    | "project_registration_opened"
    | "report_graded"
    | "report_expired"
    | "report_submitted"
    | "report_deadline"
    | "report_opened"
    | "class_updated"
    | "system";

export interface NotificationResponse {
    id: string;
    type: NotificationType;
    message: string;
    link: string;
    isRead: boolean;
    createdAt: string;
}

export interface Pagination {
    page: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
}

export interface PaginatedNotificationsResponse {
    notifications: NotificationResponse[];
    pagination: Pagination;
}
