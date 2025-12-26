<script lang="ts">
    import { Plus } from "lucide-svelte";

    import type { ReportStage } from "../../../../../types/report";
    import ReportStageCard from "./ReportStageCard.svelte";
    import ReportStageDetail from "./ReportStageDetail.svelte";
    import CreateReportStageModal from "./CreateReportStageModal.svelte";

    // 1. Quản lý State bằng Runes
    let showCreateModal = $state(false);
    let selectedStage = $state<ReportStage | null>(null);
    let reportStages = $state<ReportStage[]>([
        {
            id: "1",
            title: "Báo cáo tuần 1-2",
            description: "Báo cáo tiến độ thực hiện đề tài giai đoạn đầu",
            startDate: "2024-01-15",
            endDate: "2024-01-29",
            totalStudents: 45,
            submittedCount: 38,
            status: "completed",
        },
        {
            id: "2",
            title: "Báo cáo tuần 3-4",
            description: "Báo cáo phân tích yêu cầu và thiết kế hệ thống",
            startDate: "2024-02-01",
            endDate: "2024-02-15",
            totalStudents: 45,
            submittedCount: 42,
            status: "ongoing",
        },
        {
            id: "3",
            title: "Báo cáo tuần 5-6",
            description: "Báo cáo triển khai và coding",
            startDate: "2024-02-20",
            endDate: "2024-03-05",
            totalStudents: 45,
            submittedCount: 0,
            status: "upcoming",
        },
    ]);

    // 2. Các hàm xử lý (Handlers)
    const handleCreateStage = (
        newStage: Omit<ReportStage, "id" | "submittedCount" | "status">,
    ) => {
        const stage: ReportStage = {
            ...newStage,
            id: Date.now().toString(),
            submittedCount: 0,
            status:
                new Date(newStage.startDate) > new Date()
                    ? "upcoming"
                    : "ongoing",
        };
        reportStages = [...reportStages, stage];
        showCreateModal = false;
    };

    const handleDeleteStage = (id: string) => {
        if (confirm("Bạn có chắc chắn muốn xóa giai đoạn báo cáo này?")) {
            reportStages = reportStages.filter((stage) => stage.id !== id);
        }
    };

    const handleViewDetail = (stage: ReportStage) => {
        selectedStage = stage;
    };
</script>

{#if selectedStage}
    <ReportStageDetail
        stage={selectedStage}
        onBack={() => (selectedStage = null)}
    />
{:else}
    <div class="rounded-lg bg-white p-6 shadow-sm">
        {#if reportStages.length === 0}
            <div class="py-16 text-center">
                <button
                    onclick={() => (showCreateModal = true)}
                    class="group inline-flex h-16 w-16 items-center justify-center rounded-full border-2 border-dashed border-gray-300 transition-colors hover:border-blue-500 hover:bg-blue-50"
                >
                    <Plus
                        class="h-8 w-8 text-gray-400 group-hover:text-blue-500"
                    />
                </button>
                <p class="mt-4 text-gray-600">Chưa có giai đoạn báo cáo nào</p>
                <p class="text-sm text-gray-500">
                    Nhấn vào dấu + để tạo giai đoạn báo cáo mới
                </p>
            </div>
        {:else}
            <div>
                <div class="mb-6 flex items-center justify-between">
                    <h2 class="text-xl font-semibold">Các giai đoạn báo cáo</h2>
                    <button
                        onclick={() => (showCreateModal = true)}
                        class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
                    >
                        <Plus class="h-5 w-5" />
                        Tạo giai đoạn mới
                    </button>
                </div>

                <div class="grid grid-cols-1 gap-4">
                    {#each reportStages as stage (stage.id)}
                        <ReportStageCard
                            {stage}
                            onDelete={handleDeleteStage}
                            onViewDetail={handleViewDetail}
                        />
                    {/each}
                </div>
            </div>
        {/if}
    </div>
{/if}

{#if showCreateModal}
    <CreateReportStageModal
        onClose={() => (showCreateModal = false)}
        onSubmit={handleCreateStage}
        totalStudents={45}
    />
{/if}
