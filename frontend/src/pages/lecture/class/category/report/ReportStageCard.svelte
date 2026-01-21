<script lang="ts">
    import {
        Calendar,
        Users,
        CheckCircle,
        Clock,
        XCircle,
        Trash2,
        Eye,
    } from "lucide-svelte";
    import type { PeriodResponse } from "../../../../../dtos/period-dto";

    // Props
    let { period, onDelete, onViewDetail } = $props<{
        period: PeriodResponse;
        onDelete: (id: string) => void;
        onViewDetail: (period: PeriodResponse) => void;
    }>();

    // Derived state for status configuration
    const statusConfig = $derived.by(() => {
        switch (period.status) {
            case "completed":
                return {
                    icon: CheckCircle,
                    label: "Hoàn thành",
                    color: "text-green-600",
                    bg: "bg-green-50",
                    border: "border-green-200",
                };
            case "ongoing":
                return {
                    icon: Clock,
                    label: "Đang diễn ra",
                    color: "text-blue-600",
                    bg: "bg-blue-50",
                    border: "border-blue-200",
                };
            case "upcoming":
                return {
                    icon: Calendar,
                    label: "Sắp tới",
                    color: "text-gray-600",
                    bg: "bg-gray-50",
                    border: "border-gray-200",
                };
            case "overdue":
                return {
                    icon: XCircle,
                    label: "Quá hạn",
                    color: "text-red-600",
                    bg: "bg-red-50",
                    border: "border-red-200",
                };
            default:
                return {
                    icon: Calendar,
                    label: "Không xác định",
                    color: "text-gray-600",
                    bg: "bg-gray-50",
                    border: "border-gray-200",
                };
        }
    });

    // Derived variables
    const submissionPercentage = $derived(
        (period.submittedCount / period.totalStudents) * 100,
    );
    const StatusIcon = $derived(statusConfig.icon);
</script>

<button onclick={() => onViewDetail(period)}>
    <div
        class="border {statusConfig.border} rounded-lg p-6 hover:shadow-md transition-shadow bg-white"
    >
        <div class="flex justify-between items-start mb-4">
            <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                    <h3 class="text-lg font-medium">{period.title}</h3>
                    <span
                        class={`flex items-center gap-1 px-3 py-1 ${statusConfig.bg} ${statusConfig.color} rounded-full text-sm font-medium`}
                    >
                        <StatusIcon class="w-4 h-4" />
                        {statusConfig.label}
                    </span>
                </div>
                <p class="text-gray-600 text-sm">{period.description}</p>
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="flex items-center gap-2 text-sm text-gray-600">
                <Calendar class="w-4 h-4" />
                <span
                    >Từ {new Date(period.startDate).toLocaleDateString(
                        "vi-VN",
                    )}</span
                >
            </div>
            <div class="flex items-center gap-2 text-sm text-gray-600">
                <Calendar class="w-4 h-4" />
                <span
                    >Đến {new Date(period.endDate).toLocaleDateString(
                        "vi-VN",
                    )}</span
                >
            </div>
        </div>

        <div class="space-y-2">
            <div class="flex justify-between items-center text-sm">
                <div class="flex items-center gap-2 text-gray-600">
                    <Users class="w-4 h-4" />
                    <span
                        >Đã nộp: {period.submittedCount}/{period.totalStudents}</span
                    >
                </div>
                <span class="text-gray-600 font-medium"
                    >{submissionPercentage.toFixed(0)}%</span
                >
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div
                    class="bg-blue-600 h-2 rounded-full transition-all duration-500"
                    style="width: {submissionPercentage}%"
                ></div>
            </div>
        </div>
    </div>
</button>
