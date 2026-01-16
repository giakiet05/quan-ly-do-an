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
    is_read: boolean;
    created_at: string;
}

export interface Pagination {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
}

export interface PaginatedNotificationsResponse {
    notifications: NotificationResponse[];
    pagination: Pagination;
}
