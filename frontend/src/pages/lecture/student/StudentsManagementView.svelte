<script lang="ts">
    import {
        Plus,
        Search,
        Edit,
        Trash2,
        Eye,
        MessageSquare,
        ChevronLeft,
        ChevronRight,
        MessageCircle,
    } from "lucide-svelte";
    import { studentStore } from "../../../stores/student-store";
    import { mockStudentInClassData } from "../../../types/student";
    import { onMount } from "svelte";

    let {
        students,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredStudents,
        paginatedStudents,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        addStudent,
        updateStudent,
        removeStudent,
    } = studentStore;

    let showCreateModal = false;
    let editingStudent = null;

    onMount(() => {
        setData(mockStudentInClassData);
    });
</script>

<div class="bg-white rounded-lg shadow-sm">
    <div class="p-6 border-b border-gray-200">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-semibold">Danh sách Sinh viên</h1>
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
                placeholder="Tìm kiếm theo tên, MSSV, email..."
            />
        </div>
    </div>

    <div class="overflow-x-auto">
        <table>
            <thead class="table-header">
                <tr>
                    <th class="w-1/13">MSSV</th>
                    <th class="w-2/13">Sinh viên</th>
                    <th class="w-4/13">Email</th>
                    <th class="w-4/13">Số điện thoại</th>
                    <th class="w-2/13">Hành động</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
                {#each $paginatedStudents as student (student.id)}
                    <tr class="hover:bg-gray-50 transition-colors">
                        <td class=" px-6 py-4 font-medium"
                            >{student.studentCode}</td
                        >
                        <td class=" px-6 py-4">{student.name}</td>
                        <td class="px-6 py-4 text-gray-600"
                            >{student.email || "-"}</td
                        >
                        <td class=" px-6 py-4 text-center"
                            >{student.phone || "-"}</td
                        >
                        <td class=" px-6 py-4">
                            <div class="flex items-center gap-2">
                                <button
                                    class="flex items-center gap-1 text-green-600 hover:text-green-800"
                                    title="Xem chi tiết"
                                >
                                    <Eye size={16} />
                                </button>

                                <span class="text-gray-300">|</span>
                                <button
                                    onclick={() => removeStudent(student.id)}
                                    class="flex items-center gap-1 text-red-600 hover:text-red-800"
                                    title="Ngưng quản lý"
                                >
                                    <Trash2 size={16} />
                                </button>
                                <span class="text-gray-300">|</span>
                                <button
                                    class="flex items-center gap-1 text-purple-600 hover:text-purple-800"
                                    title="Nhắn tin"
                                >
                                    <MessageCircle size={16} />
                                </button>
                            </div>
                        </td>
                    </tr>
                {:else}
                    <tr>
                        <td
                            colspan="5"
                            class="px-6 py-12 text-center text-gray-500"
                        >
                            {$searchTerm
                                ? "Không tìm thấy sinh viên nào"
                                : "Chưa có sinh viên nào"}
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>

    {#if $filteredStudents.length > 0}
        <div
            class="p-6 border-t border-gray-200 flex justify-between items-center"
        >
            <div class="text-sm text-gray-600">
                Hiển thị {($currentPage - 1) * $itemsPerPage + 1}-{Math.min(
                    $currentPage * $itemsPerPage,
                    $filteredStudents.length,
                )} của {$filteredStudents.length} sinh viên
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
</style>
