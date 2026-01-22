<script lang="ts">
    import {
        ArrowLeft,
        Download,
        FileText,
        CheckCircle,
        Clock,
        XCircle,
        Search,
        Filter,
        Edit,
    } from "lucide-svelte";
    import type { ReportStage } from "../../../../../types/report";
    import GradeReportModal from "./GradeReportModal.svelte";
    import type { PeriodResponse } from "../../../../../dtos/period-dto";
    import { periodStore } from "../../../../../stores/period-store";

    // Định nghĩa Interface địa phương
    interface GroupReport {
        id: string;
        groupName: string;
        groupNumber: number;
        members: string[];
        submittedDate: string | null;
        status: "submitted" | "pending" | "late";
        grade: number | null;
        reportFile?: string;
        feedback?: string;
    }

    // 1. Nhận Props
    let { period, onBack } = $props<{
        period: PeriodResponse;
        onBack: () => void;
    }>();
    // period store

    // 2. Local State
    let searchTerm = $state("");
    let filterStatus = $state<"all" | "submitted" | "pending" | "late">("all");
    let gradingReport = $state<GroupReport | null>(null);

    // Mock data khởi tạo
    let groupReports = $state<GroupReport[]>([
        {
            id: "1",
            groupName: "Nhóm Alpha",
            groupNumber: 1,
            members: ["Nguyễn Văn A", "Trần Thị B", "Lê Văn C"],
            submittedDate: "2024-01-28 14:30",
            status: "submitted",
            grade: 8.5,
            reportFile: "bao-cao-tuan-1-2-nhom-1.pdf",
            feedback: "Báo cáo tốt, nội dung đầy đủ",
        },
        {
            id: "2",
            groupName: "Nhóm Beta",
            groupNumber: 2,
            members: ["Phạm Văn D", "Hoàng Thị E"],
            submittedDate: "2024-01-29 16:45",
            status: "late",
            grade: null,
            reportFile: "bao-cao-tuan-1-2-nhom-2.pdf",
        },
        {
            id: "3",
            groupName: "Nhóm Gamma",
            groupNumber: 3,
            members: ["Vũ Văn F", "Đặng Thị G", "Bùi Văn H"],
            submittedDate: null,
            status: "pending",
            grade: null,
        },
        {
            id: "4",
            groupName: "Nhóm Delta",
            groupNumber: 4,
            members: ["Ngô Văn I", "Trương Thị K"],
            submittedDate: "2024-01-27 10:20",
            status: "submitted",
            grade: 9.0,
            reportFile: "bao-cao-tuan-1-2-nhom-4.pdf",
            feedback: "Xuất sắc!",
        },
        {
            id: "5",
            groupName: "Nhóm Epsilon",
            groupNumber: 5,
            members: ["Lý Văn L", "Phan Thị M", "Võ Văn N"],
            submittedDate: null,
            status: "pending",
            grade: null,
        },
    ]);

    // 3. Derived State (Computed) - Tự động cập nhật khi state phụ thuộc thay đổi
    const filteredReports = $derived(
        groupReports.filter((report) => {
            const matchesSearch =
                report.groupName
                    .toLowerCase()
                    .includes(searchTerm.toLowerCase()) ||
                report.groupNumber.toString().includes(searchTerm) ||
                report.members.some((m) =>
                    m.toLowerCase().includes(searchTerm.toLowerCase()),
                );
            const matchesFilter =
                filterStatus === "all" || report.status === filterStatus;
            return matchesSearch && matchesFilter;
        }),
    );

    const stats = $derived({
        total: groupReports.length,
        submitted: groupReports.filter(
            (r) => r.status === "submitted" || r.status === "late",
        ).length,
        pending: groupReports.filter((r) => r.status === "pending").length,
        late: groupReports.filter((r) => r.status === "late").length,
    });

    // 4. Helper Functions
    const getStatusConfig = (status: GroupReport["status"]) => {
        const configs = {
            submitted: {
                icon: CheckCircle,
                label: "Đã nộp",
                color: "text-green-600",
                bg: "bg-green-50",
                border: "border-green-200",
            },
            late: {
                icon: Clock,
                label: "Nộp muộn",
                color: "text-orange-600",
                bg: "bg-orange-50",
                border: "border-orange-200",
            },
            pending: {
                icon: XCircle,
                label: "Chưa nộp",
                color: "text-red-600",
                bg: "bg-red-50",
                border: "border-red-200",
            },
        };
        return configs[status];
    };

    const handleGradeSubmit = (
        reportId: string,
        grade: number,
        feedback: string,
    ) => {
        groupReports = groupReports.map((report) =>
            report.id === reportId ? { ...report, grade, feedback } : report,
        );
        gradingReport = null;
    };

    const handleDeletePeriod = async () => {
        if (confirm("Bạn có chắc chắn muốn xóa giai đoạn này không?")) {
            try {
                await periodStore.removePeriod(period.classroomId, period.id);
                onBack();
            } catch (error) {
                console.error("Failed to delete period:", error);
                alert("Xóa giai đoạn thất bại. Vui lòng thử lại.");
            }
        }
    };
</script>

<div class="bg-white rounded-lg shadow-sm overflow-hidden">
    <div class="p-6 border-b border-gray-200">
        <button
            onclick={onBack}
            class="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-4 transition-colors font-medium"
        >
            <ArrowLeft class="w-5 h-5" />
            Quay lại danh sách
        </button>

        <div class="flex justify-between items-start">
            <div>
                <h2 class="text-2xl font-bold text-gray-800 mb-2">
                    {period.title}
                </h2>
                <p class="text-gray-600 mb-4">{period.description}</p>
                <div class="flex items-center gap-6 text-sm text-gray-500">
                    <span class="flex items-center gap-2">
                        Từ {new Date(period.startDate).toLocaleDateString(
                            "vi-VN",
                        )} - Đến {new Date(period.endDate).toLocaleDateString(
                            "vi-VN",
                        )}
                    </span>
                </div>
            </div>
            <button
                onclick={handleDeletePeriod}
                class="rounded-lg bg-red-600 px-6 py-2 text-white transition-colors hover:bg-red-700"
            >
                Xóa giai đoạn
            </button>
        </div>
    </div>

    <div
        class="grid grid-cols-2 md:grid-cols-4 gap-4 p-6 border-b border-gray-200 bg-gray-50/50"
    >
        <div class="bg-blue-50 rounded-lg p-4 border border-blue-100">
            <div class="text-blue-600 text-sm font-medium mb-1">
                Tổng số nhóm
            </div>
            <div class="text-2xl font-bold text-blue-900">{stats.total}</div>
        </div>
        <div class="bg-green-50 rounded-lg p-4 border border-green-100">
            <div class="text-green-600 text-sm font-medium mb-1">Đã nộp</div>
            <div class="text-2xl font-bold text-green-900">
                {stats.submitted}
            </div>
        </div>
        <div class="bg-red-50 rounded-lg p-4 border border-red-100">
            <div class="text-red-600 text-sm font-medium mb-1">Chưa nộp</div>
            <div class="text-2xl font-bold text-red-900">{stats.pending}</div>
        </div>
        <div class="bg-orange-50 rounded-lg p-4 border border-orange-100">
            <div class="text-orange-600 text-sm font-medium mb-1">Nộp muộn</div>
            <div class="text-2xl font-bold text-orange-900">{stats.late}</div>
        </div>
    </div>

    <div class="p-6 border-b border-gray-200">
        <div class="flex flex-col md:flex-row gap-4">
            <div class="flex-1 relative">
                <Search
                    class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
                />
                <input
                    type="text"
                    bind:value={searchTerm}
                    placeholder="Tìm kiếm theo tên nhóm, số nhóm hoặc thành viên..."
                    class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                />
            </div>
            <div class="flex items-center gap-2">
                <Filter class="w-5 h-5 text-gray-400" />
                <select
                    bind:value={filterStatus}
                    class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                >
                    <option value="all">Tất cả trạng thái</option>
                    <option value="submitted">Đã nộp</option>
                    <option value="pending">Chưa nộp</option>
                    <option value="late">Nộp muộn</option>
                </select>
            </div>
        </div>
    </div>

    <div class="p-6 bg-gray-50/30 min-h-[400px]">
        {#if filteredReports.length === 0}
            <div class="text-center py-20 text-gray-500 italic">
                Không tìm thấy nhóm nào phù hợp với điều kiện tìm kiếm.
            </div>
        {:else}
            <div class="space-y-4">
                {#each filteredReports as report (report.id)}
                    {@const config = getStatusConfig(report.status)}
                    {@const StatusIcon = config.icon}

                    <div
                        class="border {config.border} rounded-xl p-5 hover:shadow-lg transition-all bg-white"
                    >
                        <div class="flex justify-between items-start mb-4">
                            <div class="flex-1">
                                <div class="flex items-center gap-3 mb-2">
                                    <h3 class="text-lg font-bold text-gray-800">
                                        {report.groupName} (Nhóm {report.groupNumber})
                                    </h3>
                                    <span
                                        class={`flex items-center gap-1.5 px-3 py-1 ${config.bg} ${config.color} rounded-full text-xs font-bold uppercase tracking-wider`}
                                    >
                                        <StatusIcon class="w-3.5 h-3.5" />
                                        {config.label}
                                    </span>
                                </div>
                                <div class="text-sm text-gray-600">
                                    <span
                                        class="font-semibold text-gray-500 uppercase text-[10px]"
                                        >Thành viên:</span
                                    >
                                    {report.members.join(", ")}
                                </div>
                            </div>

                            {#if report.grade !== null}
                                <div
                                    class="ml-4 text-center bg-blue-50 border border-blue-100 rounded-lg px-4 py-2"
                                >
                                    <div
                                        class="text-[10px] text-blue-600 font-bold uppercase mb-0.5"
                                    >
                                        Điểm
                                    </div>
                                    <div
                                        class="text-2xl font-black text-blue-700 leading-none"
                                    >
                                        {report.grade}
                                    </div>
                                </div>
                            {/if}
                        </div>

                        {#if report.submittedDate}
                            <div class="mb-4 flex items-center gap-2">
                                <Clock class="w-4 h-4 text-gray-400" />
                                <span class="text-xs text-gray-500">
                                    Thời gian nộp: {new Date(
                                        report.submittedDate,
                                    ).toLocaleString("vi-VN")}
                                </span>
                            </div>
                        {/if}

                        {#if report.reportFile}
                            <div
                                class="flex items-center gap-3 mb-4 p-3 bg-gray-50 rounded-lg border border-gray-100"
                            >
                                <FileText class="w-5 h-5 text-blue-600" />
                                <span
                                    class="flex-1 text-sm font-medium text-gray-700 truncate"
                                    >{report.reportFile}</span
                                >
                                <div class="flex gap-1">
                                    <button
                                        class="px-3 py-1.5 text-xs font-bold text-blue-600 hover:bg-white rounded-md transition-all"
                                        >Tải xuống</button
                                    >
                                    <button
                                        class="px-3 py-1.5 text-xs font-bold text-blue-600 hover:bg-white rounded-md transition-all"
                                        >Xem</button
                                    >
                                </div>
                            </div>
                        {/if}

                        {#if report.feedback}
                            <div
                                class="p-4 bg-blue-50 rounded-lg mb-4 border-l-4 border-blue-400"
                            >
                                <div
                                    class="text-sm text-blue-900 leading-relaxed"
                                >
                                    <span class="font-bold">Nhận xét:</span>
                                    {report.feedback}
                                </div>
                            </div>
                        {/if}

                        <div
                            class="pt-4 border-t border-gray-100 flex justify-end"
                        >
                            {#if report.status !== "pending"}
                                <button
                                    onclick={() => (gradingReport = report)}
                                    class="flex items-center gap-2 px-4 py-2 {report.grade !==
                                    null
                                        ? 'border border-gray-300 text-gray-700 hover:bg-gray-50'
                                        : 'bg-blue-600 text-white hover:bg-blue-700'} rounded-lg transition-all text-sm font-bold shadow-sm"
                                >
                                    <Edit class="w-4 h-4" />
                                    {report.grade !== null
                                        ? "Chỉnh sửa điểm"
                                        : "Chấm điểm"}
                                </button>
                            {:else}
                                <div class="text-xs text-gray-400 italic">
                                    Nhóm chưa nộp bài
                                </div>
                            {/if}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

{#if gradingReport}
    <GradeReportModal
        report={gradingReport}
        onClose={() => (gradingReport = null)}
        onSubmit={handleGradeSubmit}
    />
{/if}
