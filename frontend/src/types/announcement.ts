// Định nghĩa Interface
export interface AttachedFile {
    id: string;
    name: string;
    size: number;
    url: string;
}

export interface Announcement {
    id: string;
    title: string;
    content: string;
    createdAt: string;
    updatedAt?: string;
    files: AttachedFile[];
    author: string;
}