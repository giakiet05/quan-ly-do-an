<script lang="ts">
    /* =======================
     * Imports
     * ======================= */
    import {
        Plus,
        X,
        Trash2,
        Search,
        Download,
        Upload,
    } from "../../../libs/Icons";
    import { mockStudents } from "../../../mocks/classes.mock";
    import type { CreateClassRequest } from "../../../types/class";
    import { z } from "zod";

    /* =======================
     * Props
     * ======================= */
    const { onClose, onSubmit } = $props<{
        onClose: () => void;
        onSubmit: (data: CreateClassRequest) => void;
    }>();

    /* =======================
     * Types & Constants
     * ======================= */
    type ActiveTab = "list" | "excel" | "upload";

    /* =======================
     * Validation schema
     * ======================= */
    const ClassSchema = z.object({
        name: z.string().min(1, "Tên lớp học là bắt buộc"),
        semester: z.string().min(1, "Học kỳ là bắt buộc"),
        school: z.string().min(1, "Trường là bắt buộc"),
        description: z.string().optional(),
    });

    /* =======================
     * State: UI
     * ======================= */
    let activeTab = $state<ActiveTab>("list");
    let searchTerm = $state("");
    let selectAll = $state(false);
    let errors = $state<Record<string, string>>({});

    /* =======================
     * State: Form
     * ======================= */
    let formData: CreateClassRequest = {
        name: "",
        semester: "",
        school: "",
        description: "",
        students: [],
    };

    /* =======================
     * State: Students
     * ======================= */
    const students = $state(
        mockStudents.map((s) => ({
            ...s,
            selected: false,
        })),
    );

    /* =======================
     * Derived
     * ======================= */
    const filteredStudents = $derived(
        students.filter((s) =>
            s.fullName.toLowerCase().includes(searchTerm.toLowerCase()),
        ),
    );

    /* =======================
     * Handlers: Form
     * ======================= */
    function handleSubmit(
        event: SubmitEvent & { currentTarget: HTMLFormElement },
    ) {
        event.preventDefault();
        errors = {};

        const result = ClassSchema.safeParse(formData);
        if (!result.success) {
            result.error.issues.forEach((issue) => {
                const field = issue.path[0];
                if (field) errors[field as string] = issue.message;
            });
            return;
        }

        onSubmit(formData);
    }

    /* =======================
     * Handlers: Students
     * ======================= */
    function toggleStudent() {
        syncSelectedStudents();
    }

    function toggleAllStudents() {
        students.forEach((s) => (s.selected = selectAll));
        syncSelectedStudents();
    }

    function syncSelectedStudents() {
        formData.students = students
            .filter((s) => s.selected)
            .map(({ fullName, studentCode, email }) => ({
                fullName,
                studentCode,
                email,
            }));
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div
        class="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden"
    >
        <!-- Header -->
        <div
            class="flex justify-between items-center p-6 border-b border-gray-200"
        >
            <div>
                <h2 class="text-xl mb-1">Tạo Lớp Học Mới</h2>
                <p class="text-sm text-gray-600">
                    Vui lòng điền đầy đủ thông tin để tạo lớp học mới cho sinh
                    viên.
                </p>
            </div>
            <button onclick={onClose} class="p-2 hover:bg-gray-100 rounded-lg">
                <X size={20} />
            </button>
        </div>

        <!-- Form -->
        <form onsubmit={handleSubmit} class="flex-1 overflow-y-auto">
            <div class="p-6 space-y-6">
                <!-- Basic info -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label for="class-name" class="block text-sm mb-2"
                            >Tên lớp học</label
                        >
                        <input
                            id="class-name"
                            class="w-full px-4 py-2 border rounded-lg"
                            class:border-red-500={errors.name}
                            bind:value={formData.name}
                        />
                        {#if errors.name}
                            <p class="text-red-500 text-sm mt-1">
                                {errors.name}
                            </p>
                        {/if}
                    </div>

                    <div>
                        <label for="semester" class="block text-sm mb-2"
                            >Học kỳ</label
                        >
                        <input
                            id="semester"
                            class="w-full px-4 py-2 border rounded-lg"
                            class:border-red-500={errors.semester}
                            bind:value={formData.semester}
                        />
                        {#if errors.semester}
                            <p class="text-red-500 text-sm mt-1">
                                {errors.semester}
                            </p>
                        {/if}
                    </div>
                    <div>
                        <label for="school" class="block text-sm mb-2"
                            >Trường</label
                        >
                        <input
                            id="school"
                            class="w-full px-4 py-2 border rounded-lg"
                            class:border-red-500={errors.school}
                            bind:value={formData.school}
                        />
                        {#if errors.school}
                            <p class="text-red-500 text-sm mt-1">
                                {errors.school}
                            </p>
                        {/if}
                    </div>
                    <div>
                        <label for="description" class="block text-sm mb-2"
                            >Mô tả</label
                        >
                        <input
                            id="description"
                            class="w-full px-4 py-2 border rounded-lg"
                            class:border-red-500={errors.description}
                            bind:value={formData.description}
                        />
                        {#if errors.description}
                            <p class="text-red-500 text-sm mt-1">
                                {errors.description}
                            </p>
                        {/if}
                    </div>
                </div>

                <!-- Tabs -->
                <div class="flex gap-2">
                    {#each ["list", "excel", "upload"] as tab}
                        <button
                            type="button"
                            onclick={() => (activeTab = tab as ActiveTab)}
                            class="px-4 py-2 rounded-lg"
                            class:bg-blue-50={activeTab === tab}
                            class:text-blue-600={activeTab === tab}
                        >
                            {tab === "list"
                                ? "Danh sách sinh viên"
                                : tab === "excel"
                                  ? "Tải mẫu excel DSSV"
                                  : "Tải lên DSSV"}
                        </button>
                    {/each}
                </div>

                <!-- List -->
                {#if activeTab === "list"}
                    <div class="flex gap-2">
                        <div class="flex-1 relative">
                            <Search
                                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
                                size={16}
                            />
                            <input
                                class="w-full pl-10 pr-4 py-2 border rounded-lg"
                                placeholder="Tìm kiếm"
                                bind:value={searchTerm}
                            />
                        </div>

                        <button
                            class="btn-add"
                            onclick={() => {
                                alert("Thêm sinh viên");
                            }}
                        >
                            <Plus size={16} /> Thêm
                        </button>

                        <button class="btn-delete" onclick={() => {}}>
                            <Trash2 size={16} /> Xóa
                        </button>
                    </div>

                    <table
                        class="w-full border border-gray-200 border-rounded-lg mt-4"
                    >
                        <thead class="table-header">
                            <tr>
                                <th>STT</th>
                                <th>MSSV</th>
                                <th>Tên sinh viên</th>
                                <th>Email</th>
                                <th>
                                    <input
                                        type="checkbox"
                                        bind:checked={selectAll}
                                        onchange={() => toggleAllStudents()}
                                    />
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each filteredStudents as s, i}
                                <tr class="hover:bg-gray-50">
                                    <td class="px-3 py-2">{i + 1}</td>
                                    <td class="px-3 py-2">{s.studentCode}</td>
                                    <td class="px-3 py-2">{s.fullName}</td>
                                    <td class="px-3 py-2">{s.email}</td>
                                    <td class="px-3 py-2 text-center">
                                        <input
                                            type="checkbox"
                                            bind:checked={s.selected}
                                            onchange={() => toggleStudent()}
                                        />
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}

                {#if activeTab === "excel"}
                    <div
                        class="flex flex-col items-center text-center p-8 border border-gray-200 rounded-lg"
                    >
                        <Download size={48} class="mx-auto text-gray-400" />
                        <button class="btn-add bg-green-500"
                            >Tải mẫu Excel</button
                        >
                    </div>
                {/if}

                {#if activeTab === "upload"}
                    <div
                        class="flex flex-col items-center text-center p-8 border-dashed border border-gray-200 rounded-lg"
                    >
                        <Upload size={48} class="mx-auto text-gray-400" />
                        <button class="btn-add">Chọn file</button>
                    </div>
                {/if}
            </div>

            <!-- Footer -->
            <div
                class="flex justify-end gap-3 p-4 border-t border-gray-200 bg-gray-50"
            >
                <button type="button" onclick={onClose} class="btn-cancel"
                    >Hủy</button
                >
                <button type="submit" class="btn-add">Tạo lớp</button>
            </div>
        </form>
    </div>
</div>
