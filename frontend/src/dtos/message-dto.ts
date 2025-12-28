import type { MessageType } from "../models";

// ==================== REQUEST DTOs ====================

export interface CreateMessageRequest {
  channelId: string;
  senderId: string;
  type: MessageType;
  content: string;
}

export interface GetMessageFilterQuery {
  channelId: string;
  senderId?: string;
  searchContent?: string;
  isRead?: boolean;
  isSend?: boolean;
  isMedia?: boolean;
  page?: number;
  pageSize?: number;
}

// ==================== RESPONSE DTOs ====================

export interface MessageResponse {
  id: string;
  channelId: string;
  senderId: string;
  senderUsername: string;
  type: MessageType;
  content: string;
  readBy: string[];
  createdAt: string;
}
