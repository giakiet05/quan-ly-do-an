<script lang="ts">
    import {
        Plus,
        Edit2,
        Trash2,
        Eye,
        Search,
        ChevronLeft,
        ChevronRight,
    } from "../../../libs/Icons";
    import { onMount } from "svelte";

    import CreateClassModal from "./CreateClassModal.svelte";
    import EditClassModal from "./EditClassModal.svelte";
    import { classStore } from "../../../stores/class-store";
    import { mockClasses } from "../../../mocks/classes.mock";
    import type { ClassItem, CreateClassRequest } from "../../../types/class";
    import { push } from "svelte-spa-router";

    let {
        classes,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredClasses,
        paginatedClasses,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        addClass,
        updateClass,
        removeClass,
    } = classStore;

    onMount(async () => {
        try {
            classStore.fetchMyClasses();
        } catch (error) {
            console.error("Failed to load classrooms", error);
        }
    });
    let showCreateModal = $state(false);
    let isProcessing = $state(false);
    let showEditModal = $state(false);
    let editingClass: ClassItem | null = null;

    async function handleCreateClass(
        classData: CreateClassRequest,
    ): Promise<void> {
        console.log(">>> PAGE ĐÃ NHẬN LỆNH TẠO LỚP:", classData);
        isProcessing = true;
        try {
            await classStore.addClass(classData);
            showCreateModal = false;
        } catch (err) {
            console.error("Failed to create class", err);
        } finally {
            isProcessing = false;
        }
    }

    function openEditModal(cls: ClassItem): void {
        editingClass = cls;
        showEditModal = true;
    }
</script>

<div class="bg-white rounded-lg shadow-sm">
    <div class="p-6 border-b border-gray-200">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-semibold">Danh sách Lớp học của tôi</h1>
            <button onclick={() => (showCreateModal = true)} class="btn-add">
                <Plus size={20} />
                Tạo Lớp học mới
            </button>
        </div>

        <div class="relative">
            <div class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
                <Search size={20} />
            </div>
            <input
                type="text"
                bind:value={$searchTerm}
                oninput={(e) =>
                    setSearch((e.target as HTMLInputElement).value || "")}
                placeholder="Tìm kiếm theo mã lớp - tên lớp..."
            />
        </div>
    </div>

    <div class="overflow-x-auto">
        <table>
            <thead class="table-header">
                <tr>
                    <th class="w-3/15">Thông tin Lớp học</th>
                    <th class="w-4/15">Mô tả</th>
                    <th class="w-2/15 text-center">Số lượng SV</th>
                    <th class="w-2/15">Học kỳ</th>
                    <th class="w-2/15 text-center">Trạng thái</th>
                    <th class="w-2/15 text-center">Hành động</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
                {#each $paginatedClasses as cls (cls.id)}
                    <tr class="hover:bg-gray-50 transition-colors">
                        <td class="px-6 py-4">
                            <div class="flex items-center gap-3">
                                <img
                                    src={cls.avatar ||
                                        `https://ui-avatars.com/api/?name=${cls.name}&background=random`}
                                    alt=""
                                    class="w-10 h-10 rounded-lg object-cover bg-gray-100 flex-shrink-0"
                                />
                                <span class="font-medium text-gray-900"
                                    >{cls.name}</span
                                >
                            </div>
                        </td>

                        <td class="px-6 py-4 text-gray-600 truncate"
                            >{cls.description || "Không có thông tin mô tả"}</td
                        >

                        <td class="px-6 py-4 text-center">
                            <span
                                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700"
                            >
                                {cls.studentCount} SV
                            </span>
                        </td>

                        <td class="px-6 py-4 text-gray-600 text-sm"
                            >{cls.semester}</td
                        >

                        <td class="px-6 py-4 text-center">
                            <span
                                class="status-tag {cls.status === 'active'
                                    ? 'bg-green-100 text-green-700'
                                    : 'bg-red-100 text-red-700'}"
                            >
                                {cls.status === "active"
                                    ? "Đang chạy"
                                    : "Tạm dừng"}
                            </span>
                        </td>

                        <td class="px-6 py-4">
                            <div class="flex items-center justify-center gap-2">
                                <button
                                    onclick={() =>
                                        push(`/lecture/my-classes/${cls.id}`)}
                                    class="text-blue-600"
                                    ><Eye size={18} /></button
                                >
                                <button
                                    onclick={() => openEditModal(cls)}
                                    class="text-green-600"
                                    ><Edit2 size={18} /></button
                                >
                                <button
                                    onclick={() => removeClass(cls.id)}
                                    class="text-red-600"
                                    ><Trash2 size={18} /></button
                                >
                            </div>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
    <!-- pagination -->
    {#if $filteredClasses.length > 0}
        <div
            class="p-6 border-t border-gray-200 flex justify-between items-center"
        >
            <div class="text-sm text-gray-600">
                Hiển thị {($currentPage - 1) * $itemsPerPage + 1}-{Math.min(
                    $currentPage * $itemsPerPage,
                    $filteredClasses.length,
                )} của {$filteredClasses.length} lớp
            </div>

            <div class="flex items-center gap-4">
                <select
                    bind:value={$itemsPerPage}
                    onchange={(e) =>
                        changePageSize(+(e.target as HTMLSelectElement).value)}
                    class="px-3 py-1 border border-gray-300 rounded"
                >
                    <option value={5}>5 hàng</option>
                    <option value={10}>10 hàng</option>
                    <option value={20}>20 hàng</option>
                </select>

                <div class="flex items-center gap-2">
                    <button
                        onclick={() => changePage($currentPage - 1)}
                        disabled={$currentPage === 1}
                        class="px-3 py-1 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
                    >
                        <ChevronLeft size={20} />
                    </button>

                    {#each Array.from({ length: $totalPages }, (_, i) => i + 1) as page}
                        <button
                            onclick={() => changePage(page)}
                            class="px-3 py-1 rounded transition-colors {page ===
                            $currentPage
                                ? 'bg-blue-600 text-white'
                                : 'border border-gray-300 hover:bg-gray-50'}"
                        >
                            {page}
                        </button>
                    {/each}

                    <button
                        onclick={() => changePage($currentPage + 1)}
                        disabled={$currentPage === $totalPages}
                        class="px-3 py-1 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
                    >
                        <ChevronRight size={20} />
                    </button>
                </div>
            </div>
        </div>
    {/if}
</div>

{#if showCreateModal}
    <CreateClassModal
        onClose={() => (showCreateModal = false)}
        onSubmit={handleCreateClass}
    />
{/if}

{#if showEditModal && editingClass}
    <EditClassModal
        classData={editingClass}
        onClose={() => (showEditModal = false)}
        onSubmit={(updatedClass) => {
            // updateClass(updatedClass);
            showEditModal = false;
            editingClass = null;
        }}
    />
{/if}

<style>
    .table-header th {
        padding: 0.75rem 1.5rem;
        text-align: left;
        font-size: 0.875rem;
        font-weight: 600;
        text-transform: uppercase;
        color: #6b7280;
        background-color: #f9fafb;
    }

    .btn-add {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        background-color: #2563eb;
        color: white;
        border-radius: 0.375rem;
        font-size: 0.875rem;
        font-weight: 500;
        transition: background-color 0.2s;
    }

    .btn-add:hover {
        background-color: #1d4ed8;
    }

    table {
        width: 100%;
        table-layout: fixed;
    }

    th,
    td {
        word-wrap: break-word;
        text-align: left;
    }

    .status-tag {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        font-size: 0.875rem;
        font-weight: 500;
        border-radius: 0.375rem;
    }
    .bg-green-100 {
        background-color: rgba(16, 185, 129, 0.2);
    }
    .text-green-700 {
        color: #047857;
    }
    .bg-red-100 {
        background-color: rgba(239, 68, 68, 0.2);
    }
    .text-red-700 {
        color: #b91c1c;
    }
</style>
