export const mockStudentInClassData: StudentInClass[] = [
    {
        id: '1',
        name: 'Nguyễn Văn A',
        studentCode: 'SV001',
        email: 'nguyenvana@student.edu.vn',
        phone: '0123456789',
        enrolledProjects: [
            {
                categoryId: '1',
                categoryName: 'Category 1',
                projectId: '1',
                projectName: 'Project 1',
                role: 'leader',
            },
        ],
    },
    {
        id: '2',
        name: 'Trần Thị B',
        studentCode: 'SV002',
        email: 'tranthib@student.edu.vn',
        phone: '0987654321',
        enrolledProjects: [
            {
                categoryId: '2',
                categoryName: 'Category 2',
                projectId: '2',
                projectName: 'Project 2',
                role: 'member',
            },
        ],
    },
    {
        id: '3',
        name: 'Lê Văn C',
        studentCode: 'SV003',
        email: 'levanc@student.edu.vn',
        phone: '0912345678',
        enrolledProjects: [
            {
                categoryId: '3',
                categoryName: 'Category 3',
                projectId: '3',
                projectName: 'Project 3',
                role: 'leader',
            },
        ],
    },
    {
        id: '4',
        name: 'Phạm Thị D',
        studentCode: 'SV004',
        email: 'phamthid@student.edu.vn',
        phone: '0934567890',
        enrolledProjects: [
            {
                categoryId: '1',
                categoryName: 'Category 1',
                projectId: '4',
                projectName: 'Project 4',
                role: 'member',
            },
        ],
    },
    {
        id: '5',
        name: 'Hoàng Văn E',
        studentCode: 'SV005',
        email: 'hoangvane@student.edu.vn',
        phone: '0967890123',
        enrolledProjects: [
            {
                categoryId: '2',
                categoryName: 'Category 2',
                projectId: '5',
                projectName: 'Project 5',
                role: 'leader',
            },
        ],
    }
];
export interface StudentInClass {
    id: string;
    name: string;
    studentCode: string;
    email: string;
    phone?: string;
    enrolledProjects: {
        categoryId: string;
        categoryName: string;
        projectId: string;
        projectName: string;
        role: 'leader' | 'member';
    }[];
}
