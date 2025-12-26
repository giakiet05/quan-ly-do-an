import type { ChannelStatus, UserInfo, Image } from "../models";

// ==================== REQUEST DTOs ====================

export interface CreateChannelRequest {
  members: UserInfo[];
}

export interface GetChannelByUserIDQuery {
  userId: string;
  page?: number;
  pageSize?: number;
}

export interface UpdateChannelRequest {
  channelId: string;
  nickname?: string;
  background?: string;
  notification?: boolean;
  typingIndicator?: boolean;
  status?: ChannelStatus;
}

// ==================== RESPONSE DTOs ====================

export interface ChannelResponse {
  id: string;
  members: ChannelMemberResponse[];
  settings: ChannelSettingResponse[];
  background?: string;
  status: ChannelStatus;
  unreadMessageCount?: number;
  createdAt: string;
  updatedAt: string;
}

export interface ChannelMemberResponse {
  userId: string;
  username: string;
  avatar?: Image;
}

export interface ChannelSettingResponse {
  userId: string;
  nickname?: string;
  notification: boolean;
  typingIndicator: boolean;
}
