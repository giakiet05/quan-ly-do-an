export const mockProjects: Project[] = [
    {
        id: "p1",
        categoryId: "c1",
        name: "Hệ thống quản lý thư viện trực tuyến",
        description: "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
        maxStudents: 3,
        currentStudents: 2,
        instructor: "TS. Nguyễn Văn A",
        tags: ["Web", "React", "Node.js", "MongoDB"],
        status: "available",
        allowStudentEdit: true,
        maxTeams: 2,
        minTeamMembers: 2,
        maxTeamMembers: 3,
    },
    {
        id: "p2",
        categoryId: "c2",
        name: "Ứng dụng di động quản lý chi tiêu cá nhân",
        description: "Ứng dụng mobile giúp người dùng theo dõi thu chi, lập kế hoạch tài chính",
        maxStudents: 2,
        currentStudents: 2,
        instructor: "ThS. Trần Thị B",
        tags: ["Mobile", "React Native", "Firebase"],
        status: "forming",
        allowStudentEdit: false,
    },
    {
        id: "p3",
        categoryId: "c3",
        name: "Hệ thống quản lý học sinh",
        description: "Xây dựng hệ thống quản lý học sinh",
        maxStudents: 4,
        currentStudents: 3,
        instructor: "TS. Nguyễn Văn C",
        tags: ["Web", "React", "Node.js", "MongoDB"],
        status: "full",
        allowStudentEdit: true,
    },
    {
        id: "p4",
        categoryId: "c4",
        name: "Ứng dụng di động quản lý công việc",
        description: "Ứng dụng mobile giúp người dùng theo dõi công việc",
        maxStudents: 2,
        currentStudents: 2,
        instructor: "ThS. Trần Thị D",
        tags: ["Mobile", "React Native", "Firebase"],
        status: "closed",
        allowStudentEdit: false,
    },
];
export interface Project {
    id: string;
    categoryId: string;
    name: string;
    description: string;
    maxStudents: number;
    currentStudents: number;
    instructor: string;
    tags: string[];
    status: 'available' | 'forming' | 'full' | 'closed'; // Added 'forming' status
    allowStudentEdit?: boolean;
    maxTeams?: number;
    minTeamMembers?: number;
    maxTeamMembers?: number;
}
