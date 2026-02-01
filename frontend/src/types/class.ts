import type { UserInfo } from "./user.ts";

export interface WhitelistEntry {
    studentCode: string;
    joinedBy: string | null;
    joinedAt: string | null;
}

export interface ReportPeriod {
    id: string;
    title: string;
    description: string;
    fileType: string[];
    startDate: string;
    endDate: string;
}

export interface ProjectRound {
    id: string;
    name: string;
    description: string;
    startDate: string;
    endDate: string;
    defaultMinMember: number;
    defaultMaxMember: number;
    reportPeriods: ReportPeriod[];
    createdAt: string;
    isDeleted: boolean;
}

export interface ClassItem {
    id: string;
    name: string;
    description: string;
    avatar: string;
    semester: string;
    year: number;
    status: "active" | "inactive" | "archived";
    studentCount: number;

    // Thông tin định danh & Bảo mật
    invitationCode: string;
    maxStudents: number;
    autoApprove: boolean;
    allowedEmailDomains: string[];
    enableWhitelist: boolean;
    enableEmailRestriction: boolean;

    lecturer: UserInfo;
    coLecturers: UserInfo[];
    students: UserInfo[];
    whitelistStudentCode?: WhitelistEntry[];
    projectRounds: ProjectRound[];

    createdAt: string;
    generalChannelId: string;
}

// Giúp code cũ không bị lỗi nếu mày đang dùng tên ClassData
export type ClassData = ClassItem;

export interface CreateClassRequest {
    name: string;
    description: string;
    avatar: string;
    semester: string;
    year: number;
    maxStudents: number;
    autoApprove: boolean;
    allowedEmailDomains: string[];
    enableWhitelist: boolean;
    enableEmailRestriction: boolean;
}

export interface UpdateClassRequest extends Partial<CreateClassRequest> {
    canStudentDeleteGroup?: boolean;
}

export type ClassStatus = "active" | "inactive" | "archived";

export interface StudentInfo {
    fullName: string;
    email: string;
    studentCode: string;
    selected?: boolean;
}