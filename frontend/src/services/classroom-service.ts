import { apiFetch } from "../utils/api-fetch";
import type {
  ClassroomResponse,
  CreateClassroomRequest,
  UpdateClassroomRequest,
  UpdateClassroomStatusRequest,
  RegenerateCodeResponse,
  UpdateWhitelistStudentCodeRequest,
  UploadWhitelistStudentCodeRequest,
} from "../dtos/classroom-dto";
import type {
  ClassroomPreviewResponse,
  JoinClassroomRequest,
  JoinClassroomResponse,
  JoinRequestResponse,
} from "../dtos/classroom-join-dto";
import type {
  ClassPostResponse,
  CreateClassPostRequest,
  UpdateClassPostRequest,
} from "../dtos/class-post-dto";
import type { ApiResponse, PaginatedResponse } from "../dtos/api-response-dto";

// ==================== LECTURER: CLASSROOM MANAGEMENT ====================

/**
 * Get all classrooms where the user is the lecturer
 */
export async function getMyClassrooms(): Promise<ClassroomResponse[]> {
  const response = await apiFetch<ApiResponse<ClassroomResponse[]>>(
    "/api/classrooms/my"
  );
  return response.data || [];
}

/**
 * Get a specific classroom by ID
 */
export async function getClassroom(
  classroomId: string
): Promise<ClassroomResponse> {
  return await apiFetch<ClassroomResponse>(
    `/api/classrooms/${classroomId}`
  );
}

/**
 * Create a new classroom
 */
export async function createClassroom(
  data: CreateClassroomRequest
): Promise<ClassroomResponse> {
  const response = await apiFetch<ApiResponse<ClassroomResponse>>(
    "/api/classrooms",
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Update classroom information
 */
export async function updateClassroom(
  classroomId: string,
  data: UpdateClassroomRequest
): Promise<ClassroomResponse> {
  const response = await apiFetch<ApiResponse<ClassroomResponse>>(
    `/api/classrooms/${classroomId}`,
    {
      method: "PUT",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Update classroom status
 */
export async function updateClassroomStatus(
  classroomId: string,
  data: UpdateClassroomStatusRequest
): Promise<ClassroomResponse> {
  const response = await apiFetch<ApiResponse<ClassroomResponse>>(
    `/api/classrooms/${classroomId}/status`,
    {
      method: "PATCH",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Delete a classroom
 */
export async function deleteClassroom(classroomId: string): Promise<void> {
  await apiFetch(`/api/classrooms/${classroomId}`, {
    method: "DELETE",
  });
}

/**
 * Regenerate invitation code
 */
export async function regenerateInvitationCode(
  classroomId: string
): Promise<RegenerateCodeResponse> {
  const response = await apiFetch<ApiResponse<RegenerateCodeResponse>>(
    `/api/classrooms/${classroomId}/regenerate-code`,
    {
      method: "POST",
    }
  );
  return response.data!;
}

/**
 * Upload whitelist of student codes
 */
export async function uploadWhitelistStudentCodes(
  classroomId: string,
  data: UploadWhitelistStudentCodeRequest
): Promise<{ message: string }> {
  const response = await apiFetch<ApiResponse<{ message: string }>>(
    `/api/classrooms/${classroomId}/whitelist`,
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Update whitelist student codes (add/remove)
 */
export async function updateWhitelistStudentCodes(
  classroomId: string,
  data: UpdateWhitelistStudentCodeRequest
): Promise<{ message: string }> {
  const response = await apiFetch<ApiResponse<{ message: string }>>(
    `/api/classrooms/${classroomId}/whitelist`,
    {
      method: "PATCH",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Remove a student from classroom
 */
export async function removeStudentFromClassroom(
  classroomId: string,
  studentId: string
): Promise<void> {
  await apiFetch(`/api/classrooms/${classroomId}/students/${studentId}`, {
    method: "DELETE",
  });
}

// ==================== STUDENT: JOIN CLASSROOM ====================

/**
 * Preview classroom information before joining
 */
export async function previewClassroom(
  invitationCode: string
): Promise<ClassroomPreviewResponse> {
  const response = await apiFetch<ApiResponse<ClassroomPreviewResponse>>(
    `/api/classrooms/preview?code=${invitationCode}`
  );
  return response.data!;
}

/**
 * Student requests to join a classroom
 */
export async function joinClassroom(
  data: JoinClassroomRequest
): Promise<JoinClassroomResponse> {
  const response = await apiFetch<ApiResponse<JoinClassroomResponse>>(
    "/api/classrooms/join",
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Student leaves a classroom
 */
export async function leaveClassroom(classroomId: string): Promise<void> {
  await apiFetch(`/api/classrooms/${classroomId}/leave`, {
    method: "POST",
  });
}

/**
 * Get all classrooms the student has joined
 */
export async function getMyJoinedClassrooms(): Promise<ClassroomResponse[]> {
  const response = await apiFetch<{classrooms: ClassroomResponse[], total: number, page: number, page_size: number}>(
    "/api/classrooms/joined"
  );
  return response.classrooms || [];
}

// ==================== LECTURER: JOIN REQUESTS ====================

/**
 * Get all pending join requests for a classroom
 */
export async function getPendingJoinRequests(
  classroomId: string
): Promise<JoinRequestResponse[]> {
  const response = await apiFetch<ApiResponse<JoinRequestResponse[]>>(
    `/api/classrooms/${classroomId}/join-requests`
  );
  return response.data || [];
}

/**
 * Approve a join request
 */
export async function approveJoinRequest(
  classroomId: string,
  requestId: string
): Promise<{ message: string }> {
  const response = await apiFetch<ApiResponse<{ message: string }>>(
    `/api/classrooms/${classroomId}/join-requests/${requestId}/approve`,
    {
      method: "POST",
    }
  );
  return response.data!;
}

/**
 * Reject a join request
 */
export async function rejectJoinRequest(
  classroomId: string,
  requestId: string
): Promise<{ message: string }> {
  const response = await apiFetch<ApiResponse<{ message: string }>>(
    `/api/classrooms/${classroomId}/join-requests/${requestId}/reject`,
    {
      method: "POST",
    }
  );
  return response.data!;
}

// ==================== CLASS POSTS (Announcements) ====================

/**
 * Get all posts in a classroom
 */
export async function getClassPosts(
  classroomId: string
): Promise<ClassPostResponse[]> {
  return await apiFetch<ClassPostResponse[]>(
    `/api/classrooms/${classroomId}/posts`
  ) || [];
}

/**
 * Get a specific post
 */
export async function getClassPost(
  classroomId: string,
  postId: string
): Promise<ClassPostResponse> {
  const response = await apiFetch<ApiResponse<ClassPostResponse>>(
    `/api/classrooms/${classroomId}/posts/${postId}`
  );
  return response.data!;
}

/**
 * Create a new post (announcement)
 */
export async function createClassPost(
  classroomId: string,
  data: CreateClassPostRequest
): Promise<ClassPostResponse> {
  const response = await apiFetch<ApiResponse<ClassPostResponse>>(
    `/api/classrooms/${classroomId}/posts`,
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Update a post
 */
export async function updateClassPost(
  classroomId: string,
  postId: string,
  data: UpdateClassPostRequest
): Promise<ClassPostResponse> {
  const response = await apiFetch<ApiResponse<ClassPostResponse>>(
    `/api/classrooms/${classroomId}/posts/${postId}`,
    {
      method: "PUT",
      body: JSON.stringify(data),
    }
  );
  return response.data!;
}

/**
 * Delete a post
 */
export async function deleteClassPost(
  classroomId: string,
  postId: string
): Promise<void> {
  await apiFetch(`/api/classrooms/${classroomId}/posts/${postId}`, {
    method: "DELETE",
  });
}

/**
 * Pin/Unpin a post
 */
export async function togglePinClassPost(
  classroomId: string,
  postId: string,
  isPinned: boolean
): Promise<ClassPostResponse> {
  const response = await apiFetch<ApiResponse<ClassPostResponse>>(
    `/api/classrooms/${classroomId}/posts/${postId}/pin`,
    {
      method: "PATCH",
      body: JSON.stringify({ isPinned }),
    }
  );
  return response.data!;
}
