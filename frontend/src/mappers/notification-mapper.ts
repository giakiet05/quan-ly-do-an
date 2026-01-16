import type {
    NotificationResponse,
} from "../dtos/notification-dto";
import type { NotificationItem } from "../types/notification";

function mapType(
    type: NotificationResponse["type"]
): NotificationItem["ui_type"] {
    switch (type) {
        case "new_message":
            return "message";

        case "project_registration_deadline":
        case "project_registration_expired":
        case "project_registration_opened":
        case "report_deadline":
        case "report_expired":
            return "deadline";

        case "report_submitted":
        case "report_graded":
            return "submission";

        case "class_updated":
            return "class";

        case "system":
            return "system";

        default:
            return "system";
    }
}
function mapTitle(type: NotificationResponse["type"]): string {
    switch (type) {
        case "new_message":
            return "Tin nhắn mới";

        case "project_registration_deadline":
            return "Sắp hết hạn đăng ký đồ án";

        case "project_registration_expired":
            return "Đã hết hạn đăng ký đồ án";

        case "project_registration_opened":
            return "Mở đăng ký đồ án";

        case "report_submitted":
            return "Bài báo cáo đã được nộp";

        case "report_graded":
            return "Bài báo cáo đã được chấm";

        case "report_deadline":
            return "Sắp đến hạn nộp báo cáo";

        case "report_expired":
            return "Đã quá hạn nộp báo cáo";

        case "class_updated":
            return "Lớp học có thay đổi";

        case "system":
            return "Thông báo hệ thống";

        default:
            return "Thông báo";
    }
}

export function mapNotification(
    dto: NotificationResponse
): NotificationItem {
    return {
        id: dto.id,
        ui_type: mapType(dto.type),
        title: mapTitle(dto.type),
        content: dto.message,
        link: dto.link || undefined,
        read: dto.is_read,
        createdAt: new Date(dto.created_at).toISOString(),
    };
}


export function mapNotifications(
    list: NotificationResponse[]
): NotificationItem[] {
    return list.map(mapNotification);
}
