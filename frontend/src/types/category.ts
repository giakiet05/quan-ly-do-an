
export const mockCategoriesList: CategoryItem[] = [
    {
        id: "1",
        name: "Web Development",
        description: "Learn how to build web applications with HTML, CSS, JavaScript and other technologies.",
        projectCount: 5,
        startDate: "2026-01-01T00:00:00Z",
        endDate: "2026-01-31T23:59:00Z",
        registeredCount: 10,
    },
    {
        id: "2",
        name: "Mobile App Development",
        description: "Learn how to build mobile applications with React Native, Flutter, and other technologies.",
        projectCount: 3,
        startDate: "2025-02-01T00:00:00Z",
        endDate: "2026-02-28T23:59:00Z",
        registeredCount: 8,
    },
    {
        id: "3",
        name: "Artificial Intelligence",
        description: "Learn how to build AI applications with Python, TensorFlow, and other technologies.",
        projectCount: 2,
        startDate: "2024-03-01T00:00:00Z",
        endDate: "2024-03-31T23:59:00Z",
        registeredCount: 5,
    },
    {
        id: "4",
        name: "Data Science",
        description: "Learn how to analyze and visualize data with Python, Pandas, Matplotlib, and other technologies.",
        projectCount: 4,
        startDate: "2024-04-01T00:00:00Z",
        endDate: "2024-04-30T23:59:00Z",
        registeredCount: 12,
    },
];
export interface CategoryItem {
    id: string;
    name: string;
    description?: string;
    projectCount: number;
    startDate: string;
    endDate: string;
    registeredCount: number;// số lượng sinh viên đã đăng ký
}
export interface ProjectCategory {
    id: string;
    name: string;
    description?: string;
    startDate: string;
    endDate: string;
}
