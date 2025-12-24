<script lang="ts">
    import {
        Plus,
        X,
        Trash2,
        Search,
        Download,
        Upload,
    } from "../../libs/Icons";
    import type { CreateClassRequest, StudentInfo } from "../../types/class";

    export let onClose: () => void;
    export let onSubmit: (classData: CreateClassRequest) => void;

    type ActiveTab = "list" | "excel" | "upload";

    interface StudentInModal {
        id: string;
        studentCode: string;
        name: string;
        email: string;
        selected?: boolean;
    }

    let activeTab: ActiveTab = "list";
    let searchTerm = "";

    let formData: CreateClassRequest = {
        name: "",
        semester: "",
        school: "",
        description: "",
        students: [],
    };

    let students: StudentInfo[] = [
        {
            fullName: "Nguyễn Văn A",
            studentCode: "SV001",
            email: "nguyenvana@student.edu.vn",
        },
        {
            fullName: "Trần Thị B",
            studentCode: "SV002",
            email: "tranthib@student.edu.vn",
        },
        {
            fullName: "Lê Văn C",
            studentCode: "SV003",
            email: "levanc@student.edu.vn",
        },
        {
            fullName: "Phạm Thị D",
            studentCode: "SV004",
            email: "phamthid@student.edu.vn",
        },
        {
            fullName: "Hoàng Văn E",
            studentCode: "SV005",
            email: "hoangvane@student.edu.vn",
        },
    ];

    let errors: Record<string, string> = {};

    function handleSubmit() {
        errors = {};

        if (!formData.name.trim()) errors.name = "Vui lòng nhập tên lớp học";
        if (!formData.semester) errors.semester = "Vui lòng chọn học kỳ";
        if (!formData.school) errors.school = "Vui lòng chọn trường đại học";

        if (Object.keys(errors).length) return;

        onSubmit({
            ...formData,
        });
    }

    let filteredStudents = students.filter(
        (s) =>
            s.fullName.toLowerCase().includes(searchTerm.toLowerCase()) ||
            s.studentCode.toLowerCase().includes(searchTerm.toLowerCase()) ||
            s.email.toLowerCase().includes(searchTerm.toLowerCase()),
    );

    function toggleStudent(student: StudentInfo) {
        student.selected = !student.selected;
    }
    function toggleAllStudents() {
        const allSelected = students.every((s) => s.selected);
        students = students.map((s) => ({ ...s, selected: !allSelected }));
    }

    function deleteSelected() {
        if (confirm("Bạn có chắc chắn muốn xóa các sinh viên đã chọn?")) {
            students = students.filter((s) => !s.selected);
        }
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
            <button on:click={onClose} class="p-2 hover:bg-gray-100 rounded-lg">
                <X size={20} />
            </button>
        </div>

        <!-- Form -->
        <form
            on:submit|preventDefault={handleSubmit}
            class="flex-1 overflow-y-auto"
        >
            <div class="p-6 space-y-6">
                <!-- Basic info -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm mb-2">Tên lớp học</label>
                        <input
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
                        <label class="block text-sm mb-2">Học kỳ</label>
                        <input
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
                </div>

                <!-- Tabs -->
                <div class="flex gap-2">
                    {#each ["list", "excel", "upload"] as tab}
                        <button
                            type="button"
                            on:click={() => (activeTab = tab as ActiveTab)}
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

                        <button class="btn-add">
                            <Plus size={16} /> Thêm
                        </button>

                        <button class="btn-delete" on:click={deleteSelected}>
                            <Trash2 size={16} /> Xóa
                        </button>
                    </div>

                    <table
                        class="w-full border border-gray-200 border-rounded-lg mt-4"
                    >
                        <thead class="bg-gray-50 border-b border-gray-200">
                            <tr>
                                <th
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                                    >STT</th
                                >
                                <th
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                                    >MSSV</th
                                >
                                <th
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                                    >Tên sinh viên</th
                                >
                                <th
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600"
                                    >Email</th
                                >
                                <th
                                    class="px-6 py-3 text-center text-sm font-medium text-gray-600"
                                >
                                    <!-- <input
                                        type="checkbox"
                                        on:change={() => toggleAllStudents()}
                                    /> -->
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
                                            on:change={() => toggleStudent(s)}
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
                <button type="button" on:click={onClose} class="btn-cancel"
                    >Hủy</button
                >
                <button type="submit" class="btn-add">Tạo lớp</button>
            </div>
        </form>
    </div>
</div>
