export interface CreatePostRequest {
    title: string;
    content: string;
    files: File[]; // Ensure this is an array of File objects
}

export interface UpdatePostRequest {
    title: string;
    content: string;
    files: File | File[];
    filesToRemove: string[];
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
    attachments: {
        fileName: string;
        fileUrl: string;
        publicId: string;
        fileSize: number;
        mimeType: string;
    }[];
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
