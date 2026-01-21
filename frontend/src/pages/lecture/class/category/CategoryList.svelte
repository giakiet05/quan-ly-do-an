<script lang="ts">
    import { push } from "svelte-spa-router";
    import {
        Edit2,
        Plus,
        FileText,
        Calendar,
        Users,
        Trash2,
    } from "lucide-svelte";
    import { projectRoundStore } from "../../../../stores/project-round-store";
    import type { ProjectRound } from "../../../../types/project-round";
    import { fade, scale } from "svelte/transition";

    // 1. Nhận Props
    let { classData, onOpen, openEditCategoryModal } = $props<{
        classData: any;
        onOpen: () => void;
        openEditCategoryModal: (category: ProjectRound) => void;
        onClose?: () => void;
    }>();

    // 2. Logic xử lý trạng thái (Status) và Derived Data
    let categories = $derived(
        $projectRoundStore.map((round) => {
            let statusConfig = {
                text: "Đang diễn ra",
                colorClass: "border-green-500",
                badgeClass: "bg-green-100 text-green-600",
            };

            if (round.status === "ended") {
                statusConfig = {
                    text: "Đã kết thúc",
                    colorClass: "border-red-500",
                    badgeClass: "bg-red-100 text-red-600",
                };
            } else if (round.status === "upcoming") {
                statusConfig = {
                    text: "Sắp diễn ra",
                    colorClass: "border-yellow-500",
                    badgeClass: "bg-yellow-100 text-yellow-600",
                };
            }
            return { ...round, statusConfig };
        }),
    );

    // 3. State quản lý xóa
    let deletingId = $state<string | null>(null);
    let confirmDeleteId = $state<string | null>(null);

    async function processDelete(id: string) {
        if (deletingId) return;
        deletingId = id;
        try {
            // Đã truyền đủ id và classroomId
            await projectRoundStore.removeRound(id, classData.id);
        } catch (error) {
            alert("Có lỗi xảy ra khi xóa!");
        } finally {
            deletingId = null;
            confirmDeleteId = null;
        }
    }
</script>

<div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
        <h2 class="text-xl font-semibold text-gray-800">
            Danh sách hạng mục đề tài
        </h2>
        <button
            onclick={onOpen}
            class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-sm"
        >
            <Plus class="w-5 h-5" />
            Tạo hạng mục mới
        </button>
    </div>

    {#if categories.length === 0}
        <div
            class="flex flex-col items-center justify-center py-20 bg-gray-50 rounded-xl border-2 border-dashed border-gray-200"
        >
            <div
                class="w-16 h-16 bg-white rounded-full flex items-center justify-center shadow-sm mb-4"
            >
                <Plus class="w-8 h-8 text-gray-400" />
            </div>
            <p class="text-gray-600 font-medium">
                Chưa có hạng mục nào được tạo
            </p>
            <p class="text-sm text-gray-400">
                Bắt đầu bằng cách tạo hạng mục đầu tiên
            </p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each categories as category (category.id)}
                <div
                    class="group relative rounded-xl border-t-4 p-5 shadow-sm transition-all hover:shadow-md bg-white
                    {category.statusConfig.colorClass} {deletingId ===
                    category.id
                        ? 'opacity-50 pointer-events-none'
                        : ''}"
                >
                    <div class="absolute top-4 right-4 flex gap-2">
                        {#if confirmDeleteId === category.id}
                            <div
                                class="flex items-center gap-1 bg-white p-1 rounded-full shadow-lg border border-red-100"
                                in:scale={{ duration: 150, start: 0.95 }}
                            >
                                <button
                                    onclick={() => processDelete(category.id)}
                                    class="px-3 py-1 bg-red-600 text-[10px] text-white font-bold rounded-full hover:bg-red-700 flex items-center gap-1"
                                >
                                    {#if deletingId === category.id}
                                        <div
                                            class="w-2 h-2 border-2 border-white/30 border-t-white rounded-full animate-spin"
                                        ></div>
                                    {/if}
                                    XÁC NHẬN
                                </button>
                                <button
                                    onclick={() => (confirmDeleteId = null)}
                                    class="px-3 py-1 bg-gray-100 text-[10px] text-gray-600 font-bold rounded-full hover:bg-gray-200"
                                >
                                    HỦY
                                </button>
                            </div>
                        {:else}
                            <div
                                class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity"
                            >
                                <button
                                    onclick={() =>
                                        openEditCategoryModal(category)}
                                    class="p-2 bg-white text-blue-600 rounded-full shadow-sm hover:bg-blue-50 border border-blue-100"
                                    title="Chỉnh sửa"
                                >
                                    <Edit2 class="w-4 h-4" />
                                </button>
                                <button
                                    onclick={() =>
                                        (confirmDeleteId = category.id)}
                                    class="p-2 bg-white text-red-600 rounded-full shadow-sm hover:bg-red-50 border border-red-100"
                                    title="Xóa"
                                >
                                    <Trash2 class="w-4 h-4" />
                                </button>
                            </div>
                        {/if}
                    </div>

                    <div class="mb-4">
                        <span
                            class="inline-block px-2 py-1 text-[10px] font-bold uppercase tracking-wider rounded-md mb-2 {category
                                .statusConfig.badgeClass}"
                        >
                            {category.statusConfig.text}
                        </span>
                        <h3
                            class="text-lg font-bold text-gray-800 cursor-pointer hover:text-blue-600 transition-colors line-clamp-1"
                            onclick={() =>
                                push(
                                    `/lecture/my-classes/${classData.id}/categories/${category.id}`,
                                )}
                        >
                            {category.name}
                        </h3>
                        <p
                            class="text-sm text-gray-500 line-clamp-2 mt-1 min-h-[40px]"
                        >
                            {category.description || "Không có mô tả."}
                        </p>
                    </div>

                    <div
                        class="grid grid-cols-2 gap-2 text-xs text-gray-600 mb-4 bg-white/50 p-2 rounded-lg"
                    >
                        <div class="flex items-center gap-1.5">
                            <Calendar class="w-3.5 h-3.5 text-gray-400" />
                            <span
                                >{new Date(
                                    category.startDate,
                                ).toLocaleDateString("vi-VN")}</span
                            >
                        </div>
                        <div
                            class="flex items-center gap-1.5 text-right justify-end"
                        >
                            <Calendar class="w-3.5 h-3.5 text-gray-400" />
                            <span
                                >{new Date(category.endDate).toLocaleDateString(
                                    "vi-VN",
                                )}</span
                            >
                        </div>
                    </div>

                    <div
                        class="flex justify-between items-center pt-4 border-t border-gray-100"
                    >
                        <div
                            class="flex items-center gap-1.5 text-sm text-gray-700"
                        >
                            <FileText class="w-4 h-4 text-blue-500" />
                            <span class="font-medium"
                                >{category.projectCount || 0}</span
                            >
                            <span class="text-gray-400 text-xs">đề tài</span>
                        </div>
                        <div
                            class="flex items-center gap-1.5 text-sm text-gray-700"
                        >
                            <Users class="w-4 h-4 text-purple-500" />
                            <span>
                                {category.registeredCount ||
                                    0}/{classData?.studentCount || 0}
                            </span>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>
