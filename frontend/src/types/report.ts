// Định nghĩa Interface
export interface ReportStage {
    id: string;
    title: string;
    description: string;
    startDate: string;
    endDate: string;
    totalStudents: number;
    submittedCount: number;
    status: "upcoming" | "ongoing" | "completed" | "overdue";
}