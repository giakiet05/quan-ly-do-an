<script lang="ts">
    import {
        Plus,
        Edit2,
        Trash2,
        Eye,
        Search,
        ChevronLeft,
        ChevronRight,
    } from "../../libs/Icons";

    import CreateClassModal from "./CreateClassModal.svelte";
    import EditClassModal from "./EditClassModal.svelte";

    // Props
    let { onSelectClass } = $props();

    // State
    let searchTerm = $state("");
    let showCreateModal = $state(false);
    let showEditModal = $state(false);
    let editingClass = $state(null);
    let currentPage = $state(1);
    let itemsPerPage = $state(5); // Make itemsPerPage reactive

    let classes = $state([
        {
            name: "Phát triển ứng dụng Web",
            school: "Trường Đại học Bách Khoa",
            description: "Lớp học về phát triển ứng dụng web hiện đại.",
            studentCount: 45,
            semester: "HK2 2023-2024",
            status: "active",
        },
        {
            name: "An toàn và bảo mật thông tin",
            school: "Trường Đại học Bách Khoa",
            description: "Lớp học về bảo mật thông tin và an toàn mạng.",
            studentCount: 58,
            semester: "HK2 2023-2024",
            status: "inactive",
        },
        {
            name: "Lập trình hướng đối tượng",
            school: "Trường Đại học Bách Khoa",
            description: "Lớp học về các nguyên lý lập trình hướng đối tượng.",
            studentCount: 55,
            semester: "HK1 2023-2024",
            status: "active",
        },
        {
            name: "Đồ án tốt nghiệp",
            school: "Trường Đại học Bách Khoa",
            description:
                "Lớp học dành cho sinh viên thực hiện đồ án tốt nghiệp.",
            studentCount: 15,
            semester: "HK2 2023-2024",
            status: "inactive",
        },
        {
            name: "Trí tuệ nhân tạo",
            school: "Trường Đại học Bách Khoa",
            description: "Lớp học về các khái niệm và ứng dụng của AI.",
            studentCount: 40,
            semester: "HK1 2023-2024",
            status: "active",
        },
    ]);

    // Computed (Derived state)
    let filteredClasses = $derived(
        classes.filter(
            (cls) =>
                cls.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                cls.semester.toLowerCase().includes(searchTerm.toLowerCase()),
        ),
    );

    let totalPages = $derived(Math.ceil(filteredClasses.length / itemsPerPage));
    let startIndex = $derived((currentPage - 1) * itemsPerPage);
    let endIndex = $derived(startIndex + itemsPerPage);
    let currentClasses = $derived(filteredClasses.slice(startIndex, endIndex));

    // Methods
    function handleCreateClass(newClass) {
        showCreateModal = false;
    }

    function handleEditClass(updatedClass) {
        showEditModal = false;
        editingClass = null;
    }

    function handleDeleteClass(classId) {
        confirm("Bạn có chắc chắn muốn xóa lớp học này?");
    }

    function openEditModal(cls) {
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
                bind:value={searchTerm}
                oninput={() => (currentPage = 1)}
                placeholder="Tìm kiếm theo mã học - tên lớp..."
            />
        </div>
    </div>

    <div class="overflow-x-auto">
        <table class="w-full">
            <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                    <th
                        class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                        >Tên Lớp học</th
                    >
                    <th
                        class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                        >Trường</th
                    >
                    <th
                        class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                        >Mô tả</th
                    >
                    <th
                        class="px-6 py-3 text-center text-sm font-medium text-gray-600"
                        >Số lượng SV</th
                    >
                    <th
                        class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                        >Học kỳ</th
                    >
                    <th
                        class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                        >Trạng thái</th
                    >
                    <th
                        class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                        >Hành động</th
                    >
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
                {#each currentClasses as cls (cls.name)}
                    <tr class="hover:bg-gray-50 transition-colors">
                        <td class="px-6 py-4 font-medium">{cls.name}</td>
                        <td class="px-6 py-4">{cls.school}</td>
                        <td class="px-6 py-4 text-gray-600"
                            >{cls.description}</td
                        >
                        <td class="px-6 py-4 text-center">{cls.studentCount}</td
                        >
                        <td class="px-6 py-4 text-gray-600">{cls.semester}</td>
                        <td class="px-6 py-4">
                            <span
                                class="status-tag {cls.status === 'active'
                                    ? 'bg-green-100 text-green-700'
                                    : 'bg-red-100 text-red-700'}"
                            >
                                {cls.status === "active"
                                    ? "Đang hoạt động"
                                    : "Ngừng hoạt động"}
                            </span>
                        </td>
                        <td class="px-6 py-4">
                            <div class="flex items-center gap-2">
                                <button
                                    onclick={() => onSelectClass(cls)}
                                    class="flex items-center gap-1 text-blue-600 hover:text-blue-800"
                                >
                                    <Eye size={16} />
                                </button>
                                <span class="text-gray-300">|</span>
                                <button
                                    onclick={() => openEditModal(cls)}
                                    class="flex items-center gap-1 text-green-600 hover:text-green-800"
                                >
                                    <Edit2 size={16} />
                                </button>
                                <span class="text-gray-300">|</span>
                                <button
                                    onclick={() => handleDeleteClass(cls.name)}
                                    class="flex items-center gap-1 text-red-600 hover:text-red-800"
                                >
                                    <Trash2 size={16} />
                                </button>
                            </div>
                        </td>
                    </tr>
                {:else}
                    <tr>
                        <td
                            colspan="7"
                            class="px-6 py-12 text-center text-gray-500"
                        >
                            {searchTerm
                                ? "Không tìm thấy lớp học nào"
                                : "Chưa có lớp học nào"}
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>

    {#if filteredClasses.length > 0}
        <div
            class="p-6 border-t border-gray-200 flex justify-between items-center"
        >
            <div class="text-sm text-gray-600">
                Hiển thị {startIndex + 1}-{Math.min(
                    endIndex,
                    filteredClasses.length,
                )} của {filteredClasses.length} lớp
            </div>

            <div class="flex items-center gap-4">
                <select
                    bind:value={itemsPerPage}
                    onchange={() => (currentPage = 1)}
                    class="px-3 py-1 border border-gray-300 rounded"
                >
                    <option value={5}>5 hàng</option>
                    <option value={10}>10 hàng</option>
                    <option value={20}>20 hàng</option>
                </select>

                <div class="flex items-center gap-2">
                    <button
                        onclick={() =>
                            (currentPage = Math.max(1, currentPage - 1))}
                        disabled={currentPage === 1}
                        class="px-3 py-1 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
                    >
                        <ChevronLeft size={20} />
                    </button>

                    {#each Array.from({ length: totalPages }, (_, i) => i + 1) as page}
                        <button
                            onclick={() => (currentPage = page)}
                            class="px-3 py-1 rounded transition-colors {currentPage ===
                            page
                                ? 'bg-blue-600 text-white'
                                : 'border border-gray-300 hover:bg-gray-50'}"
                        >
                            {page}
                        </button>
                    {/each}

                    <button
                        onclick={() =>
                            (currentPage = Math.min(
                                totalPages,
                                currentPage + 1,
                            ))}
                        disabled={currentPage === totalPages}
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
        onSubmit={handleEditClass}
    />
{/if}

<style>
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
