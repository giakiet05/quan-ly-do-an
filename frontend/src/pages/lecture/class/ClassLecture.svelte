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
    const statusMap = {
        active: {
            label: "Đang chạy",
            class: "bg-green-100 text-green-700",
        },
        inactive: { label: "Tạm dừng", class: "bg-red-100 text-red-700" },
        archived: { label: "Lưu trữ", class: "bg-gray-100 text-gray-700" },
    };
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
    let editingClass = $state<ClassItem | null>(null);

    async function handleCreateClass(
        classData: CreateClassRequest,
        file?: File,
    ): Promise<void> {
        isProcessing = true;
        try {
            await classStore.addClass(classData, file);
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

    async function handleDeleteClass(id: string) {
        if (
            confirm(
                "Bạn có chắc chắn muốn xóa lớp học này? Hành động này không thể hoàn tác.",
            )
        ) {
            await removeClass(id);
        }
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
                    <th class="w-2/15 text-center">Năm</th>
                    <th class="w-2/15 text-center">Sĩ số</th>
                    <th class="w-2/15">Học kỳ</th>
                    <th class="w-2/15 text-center">Trạng thái</th>
                    <th class="w-2/15 text-center">Hành động</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
                {#each $paginatedClasses as cls (cls.id)}
                    {@const status =
                        statusMap[cls.status as keyof typeof statusMap] ||
                        statusMap.inactive}

                    <tr
                        class="hover:bg-gray-50 transition-colors align-top cursor-pointer"
                        onclick={() => push(`/lecture/my-classes/${cls.id}`)}
                    >
                        <td class="px-6 py-4">
                            <div class="flex items-start gap-3">
                                <img
                                    src={cls.avatar ||
                                        `https://ui-avatars.com/api/?name=${cls.name}&background=random`}
                                    alt=""
                                    class="w-10 h-10 rounded-lg object-cover bg-gray-100 flex-shrink-0"
                                />
                                <span
                                    class="font-medium text-gray-900 line-clamp-3 leading-tight"
                                >
                                    {cls.name}
                                </span>
                            </div>
                        </td>

                        <td class="px-6 py-4 text-gray-600">
                            <p class="text-sm line-clamp-3 leading-normal">
                                {cls.description || "Không có thông tin mô tả"}
                            </p>
                        </td>

                        <td class="px-6 py-4 text-center text-gray-600 text-sm">
                            {cls.year}
                        </td>

                        <td class="px-6 py-4 text-center text-gray-600 text-sm">
                            {cls.studentCount}/{cls.maxStudents}
                        </td>

                        <td class="px-6 py-4 text-gray-600 text-sm text-center">
                            {cls.semester}
                        </td>

                        <td class="px-6 py-4 text-center">
                            <span class="status-tag {status.class}">
                                {status.label}
                            </span>
                        </td>

                        <td class="px-6 py-4">
                            <div class="flex items-center justify-center gap-3">
                                <button
                                    onclick={(e) => {
                                        e.stopPropagation();
                                        push(`/lecture/my-classes/${cls.id}`);
                                    }}
                                    class="text-blue-600 hover:scale-110 transition-transform"
                                    title="Xem chi tiết"
                                >
                                    <Eye size={18} />
                                </button>

                                <button
                                    onclick={(e) => {
                                        e.stopPropagation();
                                        openEditModal(cls);
                                    }}
                                    class="text-green-600 hover:scale-110 transition-transform"
                                    title="Chỉnh sửa"
                                >
                                    <Edit2 size={18} />
                                </button>

                                <button
                                    onclick={(e) => {
                                        e.stopPropagation();
                                        handleDeleteClass(cls.id);
                                    }}
                                    class="text-red-600 hover:scale-110 transition-transform"
                                    title="Xóa lớp"
                                >
                                    <Trash2 size={18} />
                                </button>
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
        onSubmit={(updatedClass, file) => {
            updateClass(editingClass?.id || "", updatedClass, file);
            showEditModal = false;
            editingClass = null;
        }}
    />
{/if}

<style>
    .line-clamp-3 {
        display: -webkit-box;
        --webkit-line-clamp: 3; /* Cho trình duyệt cũ (Chrome, Safari, Edge cũ) */
        line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
        white-space: normal;
    }
    .table-header th {
        padding: 0.75rem 1.5rem;
        text-align: left;
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: #4b5563;
        background-color: #f3f4f6;
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
        border-collapse: collapse;
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
    img {
        aspect-ratio: 1/1;
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
