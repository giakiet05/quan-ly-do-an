// ==================== REQUEST DTOS (Gửi lên Backend) ====================

export interface CreateProjectRoundRequest {
    classroomId: string;
    name: string;
    startDate: string;
    endDate: string;
    description?: string;
    defaultMinMember: number;
    defaultMaxMember: number;
}

export interface CreateProjectRoundsRequest {
    classroomId: string;
    rounds: CreateProjectRoundRequest[];
}

export interface UpdateProjectRoundRequest {
    classroomId: string;
    projectRoundId: string;
    name?: string;
    startDate?: string;
    endDate?: string;
    description?: string;
    defaultMinMember: number;
    defaultMaxMember: number;
}

export interface CreateProjectRequest {
    classroomId: string;
    projectRoundId: string;
    title: string;
    amount: number;
    description?: string;
    minMember: number;
    maxMember: number;
}

export interface CreateProjectsRequest {
    classroomId: string;
    projectRoundId: string;
    projects: CreateProjectRequest[];
}

export interface UpdateProjectRequest {
    classroomId: string;
    projectId: string;
    title?: string;
    amount?: number;
    description?: string;
    minMember?: number;
    maxMember?: number;
}

export interface CreateReportPeriodRequest {
    classroomId: string;
    projectRoundId: string;
    title: string;
    description?: string;
    fileType: string[];
    startDate: string;
    endDate: string;
}

// ==================== RESPONSE DTOS (Nhận từ Backend) ====================

export interface ReportPeriodResponse {
    id: string;
    title: string;
    description: string;
    fileType: string[];
    startDate: string; // ISO Date String
    endDate: string;   // ISO Date String
}

export interface ProjectRoundResponse {
    id: string;
    name: string;
    startDate: string;
    endDate: string;
    description: string;
    reportPeriods: ReportPeriodResponse[];
    createdAt: string;
    isDeleted: boolean;
    defaultMinMember: number;
    defaultMaxMember: number;
}

export interface ProjectResponse {
    id: string;
    classroomId: string;
    projectRoundId: string;
    title: string;
    amount: number;
    description: string;
    minMember: number;
    maxMember: number;
    status: string;
}

export interface ReportPeriodWithProjectRound {
    classroomId: string;
    classroomName: string;
    projectRoundId: string;
    projectRoundName: string;
    reportPeriods: ReportPeriodResponse; // Map với model.ReportPeriod bên Go
}

export interface ProjectRoundWithClassroom {
    classroomId: string;
    classroomName: string;
    round: ProjectRoundResponse; // Map với model.ProjectRound bên Go
}