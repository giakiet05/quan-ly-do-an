export interface Draft {
  id: number;
  title: string;
  editedTime: string;
}

export interface DraftDetail {
  title: string;
  community: string;
  bodyText?: string;
  linkUrl?: string;
  tags: string[];
  tab: "text" | "images" | "link";
}

// Mock data cleared - kept for reference
export const mockDraftsList: Draft[] = [];

export const totalDrafts = 0;

// Original data archived
export const mockDraftsDetails: Record<number, DraftDetail> = {};
