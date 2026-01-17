export interface ClassItem {
    id: string;
    name: string;
    description?: string;
    studentCount: number;
    avatar?: string;
    semester: string;
    status: "active" | "inactive"; // Add status field
}
export interface CreateClassRequest {
    name: string;
    avatar?: string;
    description?: string;
    semester: string;
    students: StudentInfo[];
}
export interface StudentInfo {
    fullName: string;
    email: string;
    studentCode: string;
    selected?: boolean;
}
export interface ClassInfo {
    name: string;
    avatar?: string;
    description?: string;
    semester: string;
}
export interface ClassData {
    id: string;
    name: string;
    avatar?: string;
    description?: string;
    studentCount: number;
    semester: string;
    status: "active" | "inactive";
}