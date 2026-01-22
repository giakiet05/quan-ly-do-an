<script lang="ts">
    import { Plus } from "lucide-svelte";
    import { periodStore } from "../../../../../stores/period-store";

    import type { ReportStage } from "../../../../../types/report";
    import ReportStageCard from "./ReportStageCard.svelte";
    import ReportStageDetail from "./ReportStageDetail.svelte";
    import type {
        CreatePeriodRequest,
        PeriodResponse,
    } from "../../../../../dtos/period-dto";
    import CreateReportStageModal from "./CreateReportStageModal.svelte";

    // Props
    let { classroomId, projectRoundId } = $props<{
        classroomId: string;
        projectRoundId: string;
    }>();

    // State
    let showCreateModal = $state(false);
    let selectedPeriod: PeriodResponse | null = $state(null);

    // Derived store data
    const { periods, fetchPeriods } = periodStore;
    $effect(() => {
        fetchPeriods(classroomId, projectRoundId);
    });

    // Handlers

    const handleDeleteStage = (id: string) => {
        if (confirm("Bạn có chắc chắn muốn xóa giai đoạn báo cáo này?")) {
            periods.update((list) => list.filter((stage) => stage.id !== id));
        }
    };

    const handleViewDetail = (period: PeriodResponse) => {
        selectedPeriod = period;
    };

    const handleCreatePeriod = async (newPeriod: CreatePeriodRequest) => {
        try {
            console.log(">>> PAGE ĐÃ NHẬN LỆNH TẠO GIAI ĐỔN:", newPeriod);
            const createdPeriod = await periodStore.addPeriod(
                newPeriod,
                classroomId,
                projectRoundId,
            );
            if (createdPeriod) {
                showCreateModal = false;
            }
        } catch (error) {
            console.error("Failed to create period:", error);
            alert("Đã xảy ra lỗi khi tạo giai đoạn báo cáo. Vui lòng thử lại.");
        }
    };
</script>

{#if selectedPeriod}
    <ReportStageDetail
        {classroomId}
        {projectRoundId}
        period={selectedPeriod}
        onBack={() => (selectedPeriod = null)}
    />
{:else}
    <div class="rounded-lg bg-white p-6 shadow-sm">
        {#if $periods.length === 0}
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
                    {#each $periods as period (period.id)}
                        <ReportStageCard
                            {period}
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
        onSubmit={handleCreatePeriod}
        totalStudents={45}
        maxStudents={60}
    />
{/if}
