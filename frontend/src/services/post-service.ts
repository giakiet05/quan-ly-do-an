import { apiFetch } from "../utils/api-fetch";
import type { CreatePostRequest, PaginatedPostResponse, PostResponse } from "../dtos/post-dto";
import type { ApiResponse } from "../dtos/api-response-dto";

// ==================== POST MANAGEMENT ====================

/**
 * Get all posts in a classroom
 */
export async function getPosts(classroomId: string): Promise<PostResponse[]> {
    const response = await apiFetch<PaginatedPostResponse>(
        `/api/classrooms/${classroomId}/posts`
    );
    return response.posts || [];
}

/**
 * Get a specific post by ID
 */
export async function getPost(classroomId: string, postId: string): Promise<PostResponse> {
    const response = await apiFetch<ApiResponse<PostResponse>>(
        `/api/posts/classrooms/${classroomId}/${postId}`
    );
    return response.data!;
}

/**
 * Create a new post
 */
export async function createPost(
    data: CreatePostRequest,
    classroomId: string
): Promise<PostResponse> {
    const formData = new FormData();
    formData.append("title", data.title);
    formData.append("content", data.content);
    if (Array.isArray(data.files)) {
        data.files.forEach((file) => formData.append("files", file));
    } else {
        formData.append("files", data.files);
    }

    const response = await apiFetch<PostResponse>(
        `/api/classrooms/${classroomId}/posts`,
        {
            method: "POST",
            body: formData,
        }
    );
    return response;
}

/**
 * Update a post
 */
export async function updatePost(
    postId: string,
    data: Partial<CreatePostRequest>
): Promise<PostResponse> {
    const formData = new FormData();
    if (data.title) formData.append("title", data.title);
    if (data.content) formData.append("content", data.content);
    if (data.files) {
        if (Array.isArray(data.files)) {
            data.files.forEach((file) => formData.append("files", file));
        } else {
            formData.append("files", data.files);
        }
    }

    const response = await apiFetch<PostResponse>(
        `/api/posts/${postId}`,
        {
            method: "PUT",
            body: formData,
        }
    );
    return response;
}

/**
 * Delete a post
 */
export async function deletePost(classroomId: string, postId: string): Promise<void> {
    await apiFetch(`/api/posts/classrooms/${classroomId}/${postId}`, {
        method: "DELETE",
    });
}
