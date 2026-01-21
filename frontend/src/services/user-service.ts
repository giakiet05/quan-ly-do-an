import { apiFetch } from "../utils/api-fetch";
import type { UserProfile, UpdateProfileRequest, ChangePasswordRequest } from "../types/user";

// ==================== SERVICE FUNCTIONS ====================

/**
 * Get current user profile
 * @returns User profile data
 */
export async function getUserProfile(): Promise<UserProfile> {
  const response = await apiFetch<UserProfile>("/api/users/me", {
    method: "GET",
  });
  return response;
}

/**
 * Update user profile (full_name, student_code)
 * @param data - Profile data to update
 * @returns Updated user profile
 */
export async function updateUserProfile(
  data: UpdateProfileRequest
): Promise<UserProfile> {
  const response = await apiFetch<UserProfile>("/api/users/me", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
  return response;
}

/**
 * Change user password
 * @param oldPassword - Current password
 * @param newPassword - New password (min 6 characters)
 * @returns Success message
 */
export async function changePassword(
  oldPassword: string,
  newPassword: string
): Promise<void> {
  await apiFetch<void>("/api/users/me/password", {
    method: "PUT",
    body: JSON.stringify({
      oldPassword,
      newPassword,
    }),
  });
}

/**
 * Upload user avatar
 * @param file - Image file to upload
 * @returns Updated user profile with new avatar
 */
export async function uploadAvatar(file: File): Promise<UserProfile> {
  const formData = new FormData();
  formData.append("avatar", file);

  const response = await apiFetch<UserProfile>("/api/users/me/avatar", {
    method: "POST",
    body: formData,
  });
  return response;
}

/**
 * Delete user avatar
 * @returns Updated user profile without avatar
 */
export async function deleteAvatar(): Promise<UserProfile> {
  const response = await apiFetch<UserProfile>("/api/users/me/avatar", {
    method: "DELETE",
  });
  return response;
}
