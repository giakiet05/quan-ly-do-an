<script lang="ts">
    import { derived } from "svelte/store";
    import { push } from "svelte-spa-router";
    import { Edit2, Eye, FileText, Plus, Trash2 } from "lucide-svelte";
    import { Users } from "@lucide/svelte";
    import { Calendar } from "@lucide/svelte";
    import CreateCategoryModal from "./CreateCategoryModal.svelte";
    import type { ClassData } from "../../../../types/class";
    import { mockCategoriesList } from "../../../../types/category";
    let { classData, onOpen, onCLose } = $props<{
        classData: ClassData;
        onOpen: () => void;
        onCLose: () => void;
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

<div class="p-6">
    <div class="flex justify-between items-center mb-6 p-6">
        <h2 class="text-xl">Danh sách hạng mục đề tài</h2>

        <button
            onclick={onOpen}
            class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
            <Plus class="w-5 h-5" />
            Tạo hạng mục mới
        </button>
    </div>
    {#if categories.length === 0}s
        <div class="text-center py-16">
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
        <div class="grid grid-cols-1 gap-4">
            {#each categories as category (category.id)}
                <div
                    class={`border rounded-lg p-2 hover:shadow-md ${category.status === "đã kết thúc" ? "border-red-500" : category.status === "săp diễn ra" ? "border-yellow-500" : "border-green-500"}`}
                >
                    <div class="flex justify-between items-start mb-4">
                        <div class="flex-1">
                            <div class="flex items-center gap-3 mb-2">
                                <h3 class="text-lg">{category.name}</h3>
                            </div>

                            <p class="text-gray-600 text-sm">
                                {category.description}
                            </p>
                        </div>
                        <div class="flex items-center gap-2">
                            <button
                                class="flex items-center gap-1 text-blue-600 hover:text-blue-800"
                            >
                                <Eye size={16} />
                            </button>
                            <span class="text-gray-300">|</span>
                            <button
                                class="flex items-center gap-1 text-green-600 hover:text-green-800"
                            >
                                <Edit2 size={16} />
                            </button>
                            <span class="text-gray-300">|</span>
                            <button
                                class="flex items-center gap-1 text-red-600 hover:text-red-800"
                            >
                                <Trash2 size={16} />
                            </button>
                        </div>
                    </div>

                    <div class="grid grid-cols-2 gap-4 mb-4">
                        <div
                            class="flex items-center gap-2 text-sm text-gray-600"
                        >
                            <Calendar class="w-4 h-4" />
                            <span>
                                Từ {new Date(
                                    category.startDate,
                                ).toLocaleDateString("vi-VN")}
                            </span>
                        </div>

                        <div
                            class="flex items-center gap-2 text-sm text-gray-600"
                        >
                            <Calendar class="w-4 h-4" />
                            <span>
                                Đến {new Date(
                                    category.endDate,
                                ).toLocaleDateString("vi-VN")}
                            </span>
                        </div>
                    </div>

                    <div class="flex items-center gap-6 mb-4 text-sm">
                        <div class="flex items-center gap-2 text-gray-600">
                            <FileText class="w-4 h-4" />
                            <span>Số đề tài: {category.projectCount}</span>
                        </div>

                        <div class="flex items-center gap-2 text-gray-600">
                            <Users class="w-4 h-4" />
                            <span>
                                Đã đăng ký: {category.registeredCount} sinh viên/{classData.studentCount}
                                sinh viên
                            </span>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>
