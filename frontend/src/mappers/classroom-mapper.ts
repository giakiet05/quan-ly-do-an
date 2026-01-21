import { email } from "zod";
import type { ClassroomResponse, CreateClassroomRequest } from "../dtos/classroom-dto";
import type { ClassItem, CreateClassRequest } from "../types/class";

export function mapClassroomToClassItem(cls: ClassroomResponse): ClassItem {
    return {
        id: cls.id,
        name: cls.name,
        avatar: cls.avatar || "",
        description: cls.description || "",
        semester: cls.semester || "",
        year: cls.year || new Date().getFullYear(),
        // Chuyển đổi status chuẩn xác
        status: (cls.status as "active" | "inactive" | "archived") || "active",

        // Thông tin bảo mật & định danh
        invitationCode: cls.invitationCode || "",
        maxStudents: cls.maxStudents || 0,
        autoApprove: cls.autoApprove || false,
        allowedEmailDomains: cls.allowedEmailDomains || [],
        enableWhitelist: cls.enableWhitelist || false,
        enableEmailRestriction: cls.enableEmailRestriction || false,

        // Nhân sự (Mapping từ DTO sang UserInfo)
        lecturer: {
            userId: cls.lecturer?.userId || "",
            fullName: cls.lecturer?.fullName || "",
            avatar: cls.lecturer?.avatar || ""
        },
        coLecturers: (cls.coLecturers || []).map(u => ({
            userId: u.userId,
            fullName: u.fullName,
            avatar: u.avatar
        })),
        students: (cls.students || []).map(u => ({
            userId: u.userId,
            fullName: u.fullName,
            avatar: u.avatar,
            email: u.email,
            studentCode: u.studentCode,
        })),

        // Whitelist
        whitelistStudentCode: (cls.whitelistStudentCode || []).map(w => ({
            studentCode: w.studentCode,
            joinedBy: w.joinedBy,
            joinedAt: w.joinedAt
        })),

        studentCount: cls.students?.length ?? 0,
        createdAt: cls.createdAt || new Date().toISOString()
    };
}

export function mapUIRequestToDTO(uiData: CreateClassRequest): CreateClassroomRequest {
    return {
        name: uiData.name,
        description: uiData.description ?? "",
        avatar: uiData.avatar ?? "",
        semester: uiData.semester,
        year: uiData.year || new Date().getFullYear(),

        // Các trường cấu hình lớp học
        maxStudents: uiData.maxStudents || 100,
        autoApprove: uiData.autoApprove || false,
        allowedEmailDomains: uiData.allowedEmailDomains || [],
        enableWhitelist: uiData.enableWhitelist || false,
        enableEmailRestriction: uiData.enableEmailRestriction || false
    };
}