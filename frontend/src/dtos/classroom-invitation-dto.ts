import type { InvitationStatus } from "../models";

// ==================== REQUEST DTOs ====================

export interface InviteRequest {
  emails: string[];
}

// ==================== RESPONSE DTOs ====================

export interface InviteResult {
  invited: string[];
  alreadyMembers: string[];
  failed: InviteError[];
  totalProcessed: number;
}

export interface InviteError {
  row?: number; // For Excel import
  email: string;
  error: string;
}

export interface InvitationResponse {
  id: string;
  email: string;
  classroomId: string;
  classroomName?: string;
  invitedBy: string;
  status: InvitationStatus;
  expiresAt: string;
  createdAt: string;
}

export interface AcceptInvitationResponse {
  message: string;
  classroomId: string;
}
