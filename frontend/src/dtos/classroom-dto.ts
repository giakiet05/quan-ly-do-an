// ==================== REQUEST DTOs ====================

export interface CreateClassroomRequest {
  name: string;
  description: string;
  avatar: string;
  semester: string; // "Fall 2026", "Spring 2027"
  year: number;
  maxStudents?: number;
  autoApprove?: boolean;
  requireEmailDomain?: string;
}

export interface UpdateClassroomRequest {
  name?: string;
  description?: string;
  avatar?: string;
  semester?: string;
  year?: number;
  maxStudents?: number;
  autoApprove?: boolean;
  requireEmailDomain?: string;
  canStudentDeleteGroup?: boolean;
}

export interface UpdateClassroomStatusRequest {
  status: "active" | "inactive" | "archived";
}

export interface UploadWhitelistStudentCodeRequest {
  studentCodes: string[];
}

export interface UpdateWhitelistStudentCodeRequest {
  addCodes?: string[];
  removeCodes?: string[];
}

// ==================== RESPONSE DTOs ====================

export interface ClassroomResponse {
  id: string;
  name: string;
  description: string;
  avatar: string;
  semester: string;
  year: number;
  status: string;
  lecturer: {
    userId: string;
    fullName: string;
    avatar: string;
  };
  students: Array<{
    userId: string;
    fullName: string;
    avatar: string;
  }>;
  invitationCode: string;
  maxStudents: number;
  autoApprove: boolean;
  whitelistStudentCode?: string[];
  createdAt: string;
}

export interface RegenerateCodeResponse {
  newCode: string;
  codeExpiresAt: string;
}
