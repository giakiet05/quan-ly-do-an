// ==================== REQUEST DTOS (Gửi lên Backend) ====================

export interface CreatePeriodRequest {
    classroomId: string;
    roundId: string;
    name: string;
    description?: string;
    startDate: string;
    endDate: string;
}

export interface CreatePeriodsRequest {
    classroomId: string;
    roundId: string;
    periods: CreatePeriodRequest[];
}

export interface UpdatePeriodRequest {
    periodId: string;
    classroomId: string;
    roundId: string;
    name?: string;
    description?: string;
    startDate?: string;
    endDate?: string;
}

// ==================== RESPONSE DTOS (Nhận từ Backend) ====================

export interface PeriodResponse {
    id: string;
    name: string;
    description: string;
    startDate: string; // ISO Date String
    endDate: string;   // ISO Date String
    createdAt: string;
    updatedAt: string;
}

export interface PeriodWithRoundResponse {
    classroomId: string;
    classroomName: string;
    roundId: string;
    roundName: string;
    period: PeriodResponse;
}
