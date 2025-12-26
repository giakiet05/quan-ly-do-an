import { apiFetch } from "./api";
import type { Message } from "../models";

// ==================== MESSAGE SERVICES ====================

/**
 * Get messages for a channel with pagination
 */
export async function getChannelMessages(
  channelId: string,
  options?: {
    limit?: number;
    before?: string; // message ID
    after?: string; // message ID
  }
): Promise<Message[]> {
  const params = new URLSearchParams();
  if (options?.limit) params.append("limit", options.limit.toString());
  if (options?.before) params.append("before", options.before);
  if (options?.after) params.append("after", options.after);

  const query = params.toString() ? `?${params.toString()}` : "";
  return apiFetch<Message[]>(`/api/channels/${channelId}/messages${query}`);
}

/**
 * Send a message to a channel
 */
export async function sendMessage(
  channelId: string,
  content: string,
  type: "user" | "system" = "user"
): Promise<Message> {
  return apiFetch<Message>(`/api/channels/${channelId}/messages`, {
    method: "POST",
    body: JSON.stringify({
      content,
      type,
    }),
  });
}

/**
 * Delete a message
 */
export async function deleteMessage(
  channelId: string,
  messageId: string
): Promise<void> {
  return apiFetch<void>(`/api/channels/${channelId}/messages/${messageId}`, {
    method: "DELETE",
  });
}

/**
 * Mark message as read
 */
export async function markMessageAsRead(
  channelId: string,
  messageId: string
): Promise<void> {
  return apiFetch<void>(
    `/api/channels/${channelId}/messages/${messageId}/read`,
    {
      method: "POST",
    }
  );
}

/**
 * Mark all messages in channel as read
 */
export async function markAllMessagesAsRead(channelId: string): Promise<void> {
  return apiFetch<void>(`/api/channels/${channelId}/messages/read-all`, {
    method: "POST",
  });
}
