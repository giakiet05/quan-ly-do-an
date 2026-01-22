import { apiFetch } from "../utils/api-fetch";

// ==================== INTERFACES ====================

export interface Message {
  id: string;
  channel_id: string;
  sender_id: string;
  sender_username: string;
  type: "user" | "system";
  content: string;
  read_by: string[];
  created_at: string;
}

export interface GetMessagesResponse {
  messages: Message[];
  total: number;
  page: number;
  page_size: number;
}

export interface GetMessagesFilterQuery {
  channel_id: string;
  sender_id?: string;
  search_content?: string;
  is_read?: boolean;
  is_send?: boolean;
  is_media?: boolean;
  page?: number;
  page_size?: number;
}

// ==================== SERVICE FUNCTIONS ====================

/**
 * Get messages with filters
 * @param query - Filter query parameters
 */
export async function getMessages(
  query: GetMessagesFilterQuery
): Promise<GetMessagesResponse> {
  const params = new URLSearchParams();

  params.append("channel_id", query.channel_id);
  if (query.sender_id) params.append("sender_id", query.sender_id);
  if (query.search_content) params.append("search_content", query.search_content);
  if (query.is_read !== undefined) params.append("is_read", String(query.is_read));
  if (query.is_send !== undefined) params.append("is_send", String(query.is_send));
  if (query.is_media !== undefined) params.append("is_media", String(query.is_media));
  params.append("page", String(query.page || 1));
  params.append("page_size", String(query.page_size || 50));
  console.log("Fetching messages with params:", params.toString());
  const response = await apiFetch<GetMessagesResponse>(
    `/api/messages/filter?${params.toString()}`,
    {
      method: "GET",
    }
  );
  return response;
}

/**
 * Get a single message by ID
 * @param channelId - Channel ID
 * @param messageId - Message ID
 */
export async function getMessageById(
  channelId: string,
  messageId: string
): Promise<Message> {
  const response = await apiFetch<Message>(
    `/api/messages/${messageId}?channel_id=${channelId}`,
    {
      method: "GET",
    }
  );
  return response;
}

/**
 * Delete a message
 * @param channelId - Channel ID
 * @param messageId - Message ID
 */
export async function deleteMessage(
  channelId: string,
  messageId: string
): Promise<void> {
  await apiFetch<void>(
    `/api/messages/${messageId}?channel_id=${channelId}`,
    {
      method: "DELETE",
    }
  );
}
