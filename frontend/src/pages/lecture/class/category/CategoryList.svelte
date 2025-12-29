<script lang="ts">
    import { derived } from "svelte/store";
    import { push } from "svelte-spa-router";
    import { Edit2, Eye, FileText, Plus, Trash2 } from "lucide-svelte";
    import { Users } from "@lucide/svelte";
    import { Calendar } from "@lucide/svelte";
    import CreateCategoryModal from "./CreateCategoryModal.svelte";
    import type { ClassData } from "../../../../types/class";
    import { mockCategoriesList } from "../../../../types/category";
    let { classData, onOpen, onCLose, openEditCategoryModal } = $props<{
        classData: ClassData;
        onOpen: () => void;
        onCLose: () => void;
        openEditCategoryModal: (category: any) => void;
    }>();
    let showCreateModal = $state(false);
    let handleCreateClass = (data: any) => {
        showCreateModal = false;
    };
    let categories = $derived(
        mockCategoriesList.map((cat) => ({
            ...cat,
            status:
                new Date(cat.endDate) < new Date()
                    ? "đã kết thúc"
                    : new Date(cat.startDate) > new Date()
                      ? "săp diễn ra"
                      : "đang diễn ra",
        })),
    );
</script>

<div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
        <h2 class="text-xl font-semibold">Danh sách hạng mục đề tài</h2>
        <button
            onclick={onOpen}
            class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
            <Plus class="w-5 h-5" />
            Tạo hạng mục mới
        </button>
    </div>

    <!-- Empty State -->
    {#if categories.length === 0}
        <div class="flex flex-col items-center justify-center py-16">
            <button
                class="inline-flex items-center justify-center w-16 h-16 rounded-full border-2 border-dashed border-gray-300 hover:border-blue-500 hover:bg-blue-50 group"
            >
                <Plus class="w-8 h-8 text-gray-400 group-hover:text-blue-500" />
            </button>
            <p class="mt-4 text-gray-600">Chưa có hạng mục nào</p>
            <p class="text-sm text-gray-500">
                Nhấn vào dấu + để tạo hạng mục mới
            </p>
        </div>
    {:else}
        <!-- Category List -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each categories as category (category.id)}
                <button
                    onclick={() =>
                        push(
                            `/lecture/my-classes/${classData.id}/categories/${category.id}`,
                        )}
                >
                    <div
                        class={`rounded-lg p-4 shadow-sm transition-shadow hover:shadow-md hover:scale-105 duration-300 ${
                            category.status === "đã kết thúc"
                                ? "border-red-500 bg-red-50"
                                : category.status === "săp diễn ra"
                                  ? "border-yellow-500 bg-yellow-50"
                                  : "border-green-500 bg-green-50"
                        }`}
                    >
                        <!-- Header -->
                        <div class="flex justify-between items-start mb-4">
                            <div>
                                <h3 class="text-lg font-medium">
                                    {category.name}
                                </h3>
                                <p class="text-sm text-gray-500">
                                    {category.description}
                                </p>
                            </div>
                            <span
                                class={`px-2 py-1 text-xs font-semibold rounded-lg ${
                                    category.status === "đã kết thúc"
                                        ? "bg-red-100 text-red-600"
                                        : category.status === "săp diễn ra"
                                          ? "bg-yellow-100 text-yellow-600"
                                          : "bg-green-100 text-green-600"
                                }`}
                            >
                                {category.status}
                            </span>
                        </div>

                        <!-- Dates -->
                        <div
                            class="flex justify-between items-center text-sm text-gray-600 mb-4"
                        >
                            <div class="flex items-center gap-2">
                                <Calendar class="w-4 h-4" />
                                <span>
                                    {new Date(
                                        category.startDate,
                                    ).toLocaleDateString("vi-VN")}
                                </span>
                            </div>
                            <div class="flex items-center gap-2">
                                <Calendar class="w-4 h-4" />
                                <span>
                                    {new Date(
                                        category.endDate,
                                    ).toLocaleDateString("vi-VN")}
                                </span>
                            </div>
                        </div>

                        <!-- Stats -->
                        <div
                            class="flex justify-between items-center text-sm text-gray-600 mb-4"
                        >
                            <div class="flex items-center gap-2">
                                <FileText class="w-4 h-4" />
                                <span>Số đề tài: {category.projectCount}</span>
                            </div>
                            <div class="flex items-center gap-2">
                                <Users class="w-4 h-4" />
                                <span>
                                    Đã đăng ký: {category.registeredCount}/{classData.studentCount}
                                </span>
                            </div>
                        </div>
                    </div>
                </button>
            {/each}
        </div>
    {/if}
</div>
