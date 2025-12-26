import type { AuthProvider, Image } from "../models";

// ==================== REQUEST DTOs ====================

export interface GetUsersQuery {
  username?: string;
  page?: number;
  pageSize?: number;
}

export interface ChangePasswordRequest {
  oldPassword: string;
  newPassword: string;
}

// ==================== RESPONSE DTOs ====================

export interface UserResponse {
  id: string;
  username: string;
  email?: string;
  provider: AuthProvider;
  avatar?: Image;
}
