export interface StudentInClass {
    userId: string;
    fullName: string;
    studentCode?: string;
    email?: string;
    phone?: string;
    enrolledProjects?: {
        categoryId: string;
        categoryName: string;
        projectId: string;
        projectName: string;
        role: 'leader' | 'member';
    }[];
}
