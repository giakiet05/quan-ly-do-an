// ==================== REQUEST DTOs ====================

export interface AttachmentUpload {
  fileName: string;
  fileURL: string;
  fileSize: number;
  mimeType: string;
}

export interface CreateClassPostRequest {
  title: string;
  content: string;
  attachments?: AttachmentUpload[];
}

export interface UpdateClassPostRequest {
  title?: string;
  content?: string;
  attachmentsToAdd?: AttachmentUpload[];
  attachmentsToRemove?: string[]; // file URLs to remove
}

// ==================== RESPONSE DTOs ====================

export interface ClassPostResponse {
  id: string;
  classroomId: string;
  author: {
    userId: string;
    fullName: string;
    avatar: string;
  };
  title: string;
  content: string;
  attachments?: Array<{
    fileName: string;
    fileURL: string;
    fileSize: number;
    mimeType: string;
  }>;
  isPinned: boolean;
  createdAt: string;
  updatedAt: string;
}
