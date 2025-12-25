export interface QueuePost {
  id: string;
  title: string;
  author: string;
  authorAvatar?: string;
  community: string;
  content: string;
  type: "text" | "image" | "link";
  createdAt: string;
  queueType: "unmoderated" | "edited" | "removed" | "reported";
  reportReason?: string;
  reportCount?: number;
  editedAt?: string;
  removedBy?: string;
  removedReason?: string;
}

// Mock data cleared - kept for reference
export const mockQueuePosts: QueuePost[] = [];

// Original data archived
