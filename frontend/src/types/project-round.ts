// Thêm export vào trước mỗi interface/type
export type RoundStatus = 'upcoming' | 'ongoing' | 'ended';

export interface ReportPeriod {
    id: string;
    title: string;
    description: string;
    fileTypes: string[];
    startDate: Date;
    endDate: Date;
    isOpen: boolean;
}

export interface ProjectRound {
    id: string;
    name: string;
    description: string;
    startDate: Date;
    endDate: Date;
    createdAt: Date;
    isDeleted: boolean;
    status: RoundStatus;
    statusText: string;
    projectCount: number;
    registeredCount: number;
    progress: number;
    minStudents: number; // Số sinh viên tối thiểu trong nhóm
    maxStudents: number; // Số sinh viên tối đa trong nhóm
    reportPeriods: ReportPeriod[]; // <--- Sử dụng interface đã export ở trên
}