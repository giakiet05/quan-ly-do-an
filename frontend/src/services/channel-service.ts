import { apiFetch } from "../utils/api-fetch";

// ==================== INTERFACES ====================

export interface ChannelMember {
  userId: string;
  username: string;
  avatar?: {
    url: string;
    public_id: string;
  };
}

export interface ChannelSetting {
  userId: string;
  nickname?: string;
  notification: boolean;
  typingIndicator: boolean;
}

export interface Channel {
  id: string;
  members: ChannelMember[];
  settings: ChannelSetting[];
  background?: string;
  status: "active" | "block";
  unread_message_count?: number;
  created_at: string;
  updated_at: string;
}

export interface GetChannelsResponse {
  channels: Channel[];
  total: number;
  page: number;
}

export interface CreateChannelRequest {
  members: {
    id: string;
    full_name: string;
    email: string;
    avatar?: {
      url: string;
      public_id: string;
    } | null;
    student_code?: string;
  }[];
}

// ==================== SERVICE FUNCTIONS ====================

/**
 * Get all channels for a user
 * @param userId - User ID to get channels for
 * @param page - Page number (default 1)
 * @param pageSize - Number of items per page (default 20)
 */
export async function getChannelsByUserId(
  userId: string,
  page: number = 1,
  pageSize: number = 20
): Promise<GetChannelsResponse> {
  const response = await apiFetch<GetChannelsResponse>(
    `/api/channels/user?user_id=${userId}&page=${page}&page_size=${pageSize}`,
    {
      method: "GET",
    }
  );
  return response;
}

/**
 * Get a specific channel by ID
 * @param channelId - Channel ID
 */
export async function getChannelById(channelId: string): Promise<Channel> {
  const response = await apiFetch<Channel>(`/api/channels/${channelId}`, {
    method: "GET",
  });
  return response;
}

/**
 * Get channel between two users (for direct messages)
 * @param user1Id - First user ID
 * @param user2Id - Second user ID
 */
export async function getChannelBetweenUsers(
  user1Id: string,
  user2Id: string
): Promise<Channel | null> {
  try {
    const response = await apiFetch<Channel>(
      `/api/channels/between/${user1Id}/${user2Id}`,
      {
        method: "GET",
      }
    );
    return response;
  } catch (error) {
    // If channel doesn't exist, return null
    return null;
  }
}

/**
 * Create a new channel
 * @param members - Array of members to add to the channel
 */
export async function createChannel(
  members: CreateChannelRequest["members"]
): Promise<Channel> {
  console.log("📤 createChannel request body:", JSON.stringify({ members }, null, 2));
  const response = await apiFetch<Channel>("/api/channels", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ members }),
  });
  return response;
}

/**
 * Update channel settings
 * @param channelId - Channel ID
 * @param updates - Settings to update
 */
export async function updateChannel(
  channelId: string,
  updates: {
    nickname?: string;
    background?: string;
    notification?: boolean;
    typingIndicator?: boolean;
    status?: "active" | "block";
  }
): Promise<Channel> {
  const response = await apiFetch<Channel>("/api/channels", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      channel_id: channelId,
      ...updates,
    }),
  });
  return response;
}

/**
 * Delete a channel
 * @param channelId - Channel ID
 */
export async function deleteChannel(channelId: string): Promise<void> {
  await apiFetch<void>(`/api/channels/${channelId}`, {
    method: "DELETE",
  });
}
