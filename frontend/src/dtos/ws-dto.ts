import type { MessageType } from "../models";
import type { MessageResponse } from "./message-dto";

// ==================== WEBSOCKET MESSAGE TYPES ====================

export type WebSocketMessageType =
  | "new_notification"
  | "ack_message"
  | "new_message"
  | "send_message"
  | "typing"
  | "in_chat"
  | "error";

// ==================== WEBSOCKET MESSAGE ====================

export interface WebSocketMessage<T = any> {
  type: WebSocketMessageType;
  payload: T;
}

// ==================== PAYLOAD DTOs ====================

export interface NewMessagePayload {
  tempMessageId: string;
  channelId: string;
  senderId: string;
  senderUsername: string;
  type: MessageType;
  content: string;
}

export interface SendMessagePayload {
  message: MessageResponse;
}

export interface ACKMessagePayload {
  tempMessageId: string;
  message: MessageResponse;
}

export interface TypingIndicatorPayload {
  channelId: string;
  senderId: string;
  isTyping: boolean;
}

export interface InChatIndicatorPayload {
  channelId: string;
  isInChat: boolean;
}

export interface ErrorPayload {
  tempMessageId?: string;
  errorCode?: string;
  errorMsg: string;
}
