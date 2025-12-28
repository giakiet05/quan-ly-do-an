import type { NotificationType } from "../models";

// ==================== RESPONSE DTOs ====================

export interface NotificationResponse {
  id: string;
  type: NotificationType;
  message: string;
  link: string;
  isRead: boolean;
  createdAt: string;
}
