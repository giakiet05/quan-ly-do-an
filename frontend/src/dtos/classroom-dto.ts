// ==================== INTERFACES ====================

import type { StudentInfo } from "../types/class";

export interface WhitelistEntry {
  studentCode: string;
  joinedBy: string | null;
  joinedAt: string | null;
}

// ==================== REQUEST DTOs ====================

export interface CreateClassroomRequest {
  name: string;
  description: string;
  avatar: string;
  semester: string;
  year: number;
  maxStudents?: number;
  autoApprove?: boolean;
  allowedEmailDomains: string[];
  enableWhitelist: boolean;
  enableEmailRestriction: boolean;
  whitelistStudentCodes?: string[];
}

export interface WhitelistEntry {
  studentCode: string;
  joinedBy: string | null; // BE trả về pointer string nên có thể null
  joinedAt: string | null; // BE trả về pointer string định dạng RFC3339
}

// UserInfoResponse dùng chung cho Lecturer, CoLecturers, Students
export interface UserInfo {
  userId: string;
  fullName: string;
  avatar: string;
  studentCode?: string;
  email?: string;
}

export interface UpdateClassroomRequest {
  name?: string;
  description?: string;
  avatar?: string;
  semester?: string;
  year?: number;
  maxStudents?: number;
  autoApprove?: boolean;
  canStudentDeleteGroup?: boolean;
  allowedEmailDomains: string[];
  enableWhitelist: boolean;
  enableEmailRestriction: boolean;
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
  generalChannelId: string;
  lecturer: {
    userId: string;
    fullName: string;
    avatar: string;
  };
  coLecturers?: Array<{
    userId: string;
    fullName: string;
    avatar: string;
  }>;
  students: Array<UserInfo>;
  projectRounds: Array<{
    id: string;
    name: string;
    startDate: string;
    endDate: string;
    description: string;
    reportPeriods: Array<{
      id: string;
      title: string;
      description: string;
      fileType: string[];
      startDate: string;
      endDate: string;
    }>;
    createdAt: string;
    isDeleted: boolean;
  }>;
  invitationCode: string;
  maxStudents: number;
  autoApprove: boolean;
  canStudentDeleteGroup: boolean;
  whitelistStudentCode?: WhitelistEntry[];
  allowedEmailDomains?: string[];
  enableWhitelist?: boolean;
  enableEmailRestriction?: boolean;
  createdAt: string;
}

export interface RegenerateCodeResponse {
  newCode: string;
  codeExpiresAt: string;
}
