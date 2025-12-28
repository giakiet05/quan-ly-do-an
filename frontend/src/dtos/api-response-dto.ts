// ==================== API RESPONSE ====================

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error_code?: string;
}

// ==================== PAGINATION ====================

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: Pagination;
}

export interface PaginatedUsersResponse {
  users: UserResponse[];
  pagination: Pagination;
}

export interface PaginatedChannelsResponse {
  channels: ChannelResponse[];
  pagination: Pagination;
}

export interface PaginatedMessagesResponse {
  messages: MessageResponse[];
  pagination: Pagination;
}

export interface PaginatedNotificationsResponse {
  notifications: NotificationResponse[];
  pagination: Pagination;
}
