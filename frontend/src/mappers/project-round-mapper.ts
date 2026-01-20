import type { ProjectRoundResponse } from "../dtos/project-dto";
import type { ProjectRound, ReportPeriod } from "../types/project-round";

export const projectRoundMapper = {
    toEntity(dto: ProjectRoundResponse): ProjectRound {
        const now = new Date();
        const start = new Date(dto.startDate);
        const end = new Date(dto.endDate);

        const mappedPeriods: ReportPeriod[] = (dto.reportPeriods || []).map(rp => {
            const rpStart = new Date(rp.startDate);
            const rpEnd = new Date(rp.endDate);
            return {
                id: rp.id,
                title: rp.title,
                description: rp.description,
                fileTypes: rp.fileType,
                startDate: rpStart,
                endDate: rpEnd,
                isOpen: now >= rpStart && now <= rpEnd
            };
        });

        // Xác định status enum
        const status: 'upcoming' | 'ongoing' | 'ended' =
            now < start ? 'upcoming' : (now > end ? 'ended' : 'ongoing');

        return {
            id: dto.id,
            name: dto.name,
            description: dto.description || "",
            startDate: start,
            endDate: end,
            createdAt: new Date(dto.createdAt),
            isDeleted: dto.isDeleted,
            status: status,

            // 1. Thêm statusText để UI dùng trực tiếp (thay thế logic trong category.ts)
            statusText: status === 'ended' ? "đã kết thúc" : (status === 'upcoming' ? "sắp diễn ra" : "đang diễn ra"),

            // 2. Map các trường từ DTO (nếu backend chưa có thì để mặc định theo mock cũ)
            // Lưu ý: Bạn cần kiểm tra lại Backend có trả về 2 trường này không
            projectCount: (dto as any).projectCount ?? 0,
            registeredCount: (dto as any).registeredCount ?? 0,

            progress: calculateProgress(start, end),
            reportPeriods: mappedPeriods
        };
    }
};

function calculateProgress(start: Date, end: Date): number {
    const total = end.getTime() - start.getTime();
    if (total <= 0) return 0;
    const elapsed = new Date().getTime() - start.getTime();
    return Math.min(Math.max(Math.round((elapsed / total) * 100), 0), 100);
}