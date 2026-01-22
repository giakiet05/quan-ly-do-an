export interface CreatePostRequest {
    title: string;
    content: string;
    files: File | File[]; // Một hoặc nhiều file từ input type="file"
}

export interface UpdatePostRequest {
    title: string;
    content: string;
    files: File | File[];
    filesToRemove: File | File[];
}

export interface PostResponse {
    id: string;
    classroom_id: string;
    author: {
        userId: string;
        fullName: string;
        email: string;
        avatar: {
            url: string;
            publicId: string;
            uploadedAt: string;
        };
    };
    title: string;
    content: string;
    isPinned: boolean;
    createdAt: string;
    updatedAt: string;
}

export interface PaginatedPostResponse {
    page: number;
    pageSize: number;
    posts: PostResponse[];
    total: number;
}

export interface PinPostRequest {
    isPinned: boolean;
}