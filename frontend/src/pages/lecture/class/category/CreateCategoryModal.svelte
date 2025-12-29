<script lang="ts">
    import { X } from "lucide-svelte";
    import type { ClassData } from "../../../../types/class";
    import type { ProjectCategory } from "../../../../types/category";

    // Định nghĩa Props theo cú pháp Svelte 5
    let { onClose, onSubmit, classData, editingCategory } = $props<{
        onClose: () => void;
        onSubmit: (category: any) => void;
        classData: ClassData;
        editingCategory?: ProjectCategory;
    }>();
    // Helper function để đảm bảo định dạng YYYY-MM-DD
    const formatToDateInput = (dateStr: string | undefined) => {
        if (!dateStr) return "";
        // Tách lấy phần trước chữ T (YYYY-MM-DD)
        return dateStr.split("T")[0];
    };

    // Khởi tạo state cho form
    let formData = $state({
        name: editingCategory?.name || "",
        description: editingCategory?.description || "",
        // Dùng helper để convert '2026-01-01T00:00:00Z' thành '2026-01-01'
        startDate: formatToDateInput(editingCategory?.startDate),
        endDate: formatToDateInput(editingCategory?.endDate),
        status: editingCategory?.status || "upcoming",
    });
    let errors = $state<Record<string, string>>({});

    function handleSubmit(e: Event) {
        e.preventDefault();

        const newErrors: Record<string, string> = {};

        if (!formData.name.trim()) {
            newErrors.name = "Vui lòng nhập tên hạng mục";
        }

        if (!formData.startDate) {
            newErrors.startDate = "Vui lòng chọn ngày bắt đầu";
        }

        if (!formData.endDate) {
            newErrors.endDate = "Vui lòng chọn ngày kết thúc";
        }

        if (
            formData.startDate &&
            formData.endDate &&
            new Date(formData.startDate) > new Date(formData.endDate)
        ) {
            newErrors.endDate = "Ngày kết thúc phải sau ngày bắt đầu";
        }

        if (Object.keys(newErrors).length > 0) {
            errors = newErrors;
            return;
        }

        // Gửi dữ liệu đi (Svelte tự động bóc tách proxy của $state khi truyền ra ngoài)
        onSubmit({ ...formData });
    }

    // Xử lý xóa lỗi khi người dùng nhập liệu
    function clearError(field: string) {
        if (errors[field]) {
            delete errors[field];
        }
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4">
        <div
            class="flex justify-between items-center p-6 border-b border-gray-200"
        >
            <h2 class="text-xl">
                {editingCategory
                    ? "Chỉnh sửa hạng mục"
                    : "Tạo hạng mục đề tài mới"}
            </h2>
            <button
                onclick={onClose}
                class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
                <X class="w-5 h-5" />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="p-6 space-y-6">
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p class="text-sm text-blue-800">
                    <strong>Lớp học:</strong>
                    {classData.name} - {classData.semester}
                </p>
            </div>

            <div>
                <label class="block text-sm mb-2" for="name">
                    Tên hạng mục <span class="text-red-500">*</span>
                </label>
                <input
                    id="name"
                    type="text"
                    bind:value={formData.name}
                    oninput={() => clearError("name")}
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.name
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Ví dụ: Hạng mục 1 - Đề tài Web Application"
                />
                {#if errors.name}
                    <p class="text-red-500 text-sm mt-1">{errors.name}</p>
                {/if}
            </div>

            <div>
                <label class="block text-sm mb-2" for="description">Mô tả</label
                >
                <textarea
                    id="description"
                    bind:value={formData.description}
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows={3}
                    placeholder="Mô tả ngắn về hạng mục này"
                ></textarea>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label class="block text-sm mb-2" for="startDate">
                        Ngày bắt đầu <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="startDate"
                        type="date"
                        bind:value={formData.startDate}
                        oninput={() => clearError("startDate")}
                        class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.startDate
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                    />
                    {#if errors.startDate}
                        <p class="text-red-500 text-sm mt-1">
                            {errors.startDate}
                        </p>
                    {/if}
                </div>

                <div>
                    <label class="block text-sm mb-2" for="endDate">
                        Ngày kết thúc <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="endDate"
                        type="date"
                        bind:value={formData.endDate}
                        oninput={() => clearError("endDate")}
                        class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.endDate
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                    />
                    {#if errors.endDate}
                        <p class="text-red-500 text-sm mt-1">
                            {errors.endDate}
                        </p>
                    {/if}
                </div>
            </div>

            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <p class="text-sm text-yellow-800">
                    <strong>Lưu ý:</strong> Sau khi tạo hạng mục, bạn có thể thêm
                    đề tài, thiết lập báo cáo và cài đặt cho hạng mục này.
                </p>
            </div>

            <div class="flex justify-end gap-3 pt-4 border-t border-gray-200">
                <button
                    type="button"
                    onclick={onClose}
                    class="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                >
                    Hủy
                </button>
                <button
                    type="submit"
                    class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                >
                    {editingCategory ? "Cập nhật" : "Tạo hạng mục"}
                </button>
            </div>
        </form>
    </div>
</div>
