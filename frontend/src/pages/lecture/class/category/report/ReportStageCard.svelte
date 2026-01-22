<script lang="ts">
    import {
        Calendar,
        Users,
        Clock,
        XCircle,
        Trash2,
        Eye,
        ChevronRight,
    } from "lucide-svelte";
    import type { PeriodResponse } from "../../../../../dtos/period-dto";

    let { period, onDelete, onViewDetail } = $props<{
        period: PeriodResponse;
        onDelete: (id: string) => void;
        onViewDetail: (period: PeriodResponse) => void;
    }>();

    const statusConfig = $derived.by(() => {
        switch (period.status) {
            case "ongoing":
                return {
                    icon: Clock,
                    label: "Đang diễn ra",
                    color: "text-emerald-700",
                    bg: "bg-emerald-50",
                    border: "border-emerald-200",
                    bar: "bg-emerald-500",
                };
            case "upcoming":
                return {
                    icon: Calendar,
                    label: "Sắp tới",
                    color: "text-amber-700",
                    bg: "bg-amber-50",
                    border: "border-amber-200",
                    bar: "bg-amber-500",
                };
            case "overdue":
                return {
                    icon: XCircle,
                    label: "Quá hạn",
                    color: "text-rose-700",
                    bg: "bg-rose-50",
                    border: "border-rose-200",
                    bar: "bg-rose-500",
                };
            default:
                return {
                    icon: Calendar,
                    label: "N/A",
                    color: "text-slate-600",
                    bg: "bg-slate-50",
                    border: "border-slate-200",
                    bar: "bg-slate-400",
                };
        }
    });

    const submissionPercentage = $derived(
        period.totalStudents > 0
            ? (period.submittedCount / period.totalStudents) * 100
            : 0,
    );
    const StatusIcon = $derived(statusConfig.icon);
</script>

<div
    class="bg-white border-2 {statusConfig.border} rounded-xl p-5 shadow-sm hover:shadow-md transition-all duration-200"
>
    <div class="flex flex-wrap items-center justify-between gap-4 mb-6">
        <button
            onclick={() => onViewDetail(period)}
            class="flex items-center gap-4 flex-1 min-w-0 hover:bg-slate-50 p-2 rounded-lg transition-colors text-left"
        >
            <div class="flex items-center gap-3">
                <div class="p-2 rounded-lg {statusConfig.bg}">
                    <StatusIcon class="w-6 h-6 {statusConfig.color}" />
                </div>
                <div>
                    <h3 class="text-xl font-bold text-slate-800">
                        {period.title}
                    </h3>
                    <span
                        class="text-xs font-bold uppercase tracking-widest {statusConfig.color}"
                        >{statusConfig.label}</span
                    >
                </div>
            </div>
        </button>

        <div class="flex items-center gap-2">
            <button
                onclick={() => onViewDetail(period)}
                class="flex items-center gap-2 px-6 py-2.5 bg-slate-900 hover:bg-slate-800 text-white rounded-lg font-semibold transition-colors shadow-sm"
            >
                <Eye class="w-5 h-5" />
                Chi tiết
            </button>
            <button
                onclick={() => onDelete(period.id)}
                class="p-2.5 text-rose-600 hover:bg-rose-50 border border-rose-100 rounded-lg transition-colors"
                title="Xóa"
            >
                <Trash2 class="w-5 h-5" />
            </button>
        </div>
    </div>

    <div
        class="grid grid-cols-1 lg:grid-cols-2 gap-6 p-4 bg-slate-50 rounded-xl border border-slate-100 mb-6"
    >
        <div class="space-y-1">
            <span
                class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
                >Mô tả đợt</span
            >
            <p class="text-slate-600 text-sm leading-relaxed">
                {period.description ||
                    "Chưa có thông tin mô tả chi tiết cho đợt này."}
            </p>
        </div>

        <div
            class="flex items-center justify-between lg:border-l lg:pl-6 border-slate-200"
        >
            <div class="space-y-3 w-full">
                <span
                    class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
                    >Thời hạn</span
                >
                <div
                    class="flex items-center justify-between bg-white p-2 rounded-lg border border-slate-200"
                >
                    <div class="text-center px-4">
                        <span class="block text-[10px] text-slate-400 italic"
                            >Bắt đầu</span
                        >
                        <span class="text-sm font-bold text-slate-700"
                            >{new Date(period.startDate).toLocaleDateString(
                                "vi-VN",
                            )}</span
                        >
                    </div>
                    <ChevronRight class="w-4 h-4 text-slate-300" />
                    <div class="text-center px-4">
                        <span class="block text-[10px] text-slate-400 italic"
                            >Kết thúc</span
                        >
                        <span class="text-sm font-bold text-rose-600"
                            >{new Date(period.endDate).toLocaleDateString(
                                "vi-VN",
                            )}</span
                        >
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="space-y-3">
        <div class="flex justify-between items-end">
            <div class="flex items-center gap-2">
                <div class="flex -space-x-2">
                    <div
                        class="w-7 h-7 rounded-full border-2 border-white bg-slate-200 flex items-center justify-center text-[10px] font-bold"
                    >
                        1
                    </div>
                    <div
                        class="w-7 h-7 rounded-full border-2 border-white bg-slate-300 flex items-center justify-center text-[10px] font-bold"
                    >
                        2+
                    </div>
                </div>
                <span class="text-sm text-slate-500 ml-2">
                    Đã nộp: <span class="font-bold text-slate-800"
                        >{period.submittedCount}</span
                    >/{period.totalStudents} sinh viên
                </span>
            </div>
            <span
                class="px-2 py-0.5 bg-white border border-slate-200 rounded text-sm font-black text-slate-700"
            >
                {submissionPercentage.toFixed(0)}%
            </span>
        </div>

        <div
            class="w-full bg-slate-200 rounded-full h-3 overflow-hidden shadow-inner"
        >
            <div
                class="{statusConfig.bar} h-full rounded-full transition-all duration-1000 relative"
                style="width: {submissionPercentage}%"
            >
                <div class="absolute inset-0 bg-white/20 animate-pulse"></div>
            </div>
        </div>
    </div>
</div>
