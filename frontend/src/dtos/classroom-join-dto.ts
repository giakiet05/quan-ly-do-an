// ==================== REQUEST DTOs ====================

export interface JoinClassroomRequest {
  invitationCode: string;
}

export interface UploadWhitelistRequest {
  studentCodes: string[];
}

// ==================== RESPONSE DTOs ====================

export interface ClassroomPreviewResponse {
  id: string;
  name: string;
  lecturer: string;
  studentCount: number;
  maxStudents: number;
}

export interface JoinClassroomResponse {
  classroomId: string;
  status: string; // "approved" or "pending"
  message: string;
}

export interface JoinRequestResponse {
  id: string;
  studentCode: string;
  fullName: string;
  email: string;
  status: string;
  createdAt: string;
}
