export interface ClassItem {
    name: string;
    school: string;
    description?: string;
    studentCount: number;
    semester: string;
    status: "active" | "inactive"; // Add status field
}
export interface CreateClassRequest {
    name: string;
    school: string;
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
