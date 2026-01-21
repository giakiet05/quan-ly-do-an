<script lang="ts">
    import { X } from "lucide-svelte";
    import type { ReportStage } from "../../../../../types/report";
    import type { CreatePeriodRequest } from "../../../../../dtos/period-dto";
    // 1. Định nghĩa Props
    let { onClose, onSubmit, totalStudents, maxStudents } = $props<{
        onClose: () => void;
        onSubmit: (data: CreatePeriodRequest) => void;
        totalStudents: number;
        maxStudents: number;
    }>();

    // 2. State cho Form và Errors
    let formData = $state({
        title: "",
        description: "",
        startDate: "",
        endDate: "",
    });

    let errors = $state<Record<string, string>>({});

    // 3. Logic xử lý Submit
    const handleSubmit = (e: Event) => {
        e.preventDefault();

        const newErrors: Record<string, string> = {};

        if (!formData.title.trim()) {
            newErrors.title = "Vui lòng nhập tên";
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

        onSubmit({
            ...formData,
            startDate: new Date(formData.startDate),
            endDate: new Date(formData.endDate),
        });
    };

    // Helper để xóa lỗi khi người dùng nhập liệu (thay thế logic trong onChange của React)
    function clearError(field: string) {
        if (errors[field]) {
            delete errors[field];
        }
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div class="mx-4 w-full max-w-2xl rounded-lg bg-white shadow-xl">
        <div
            class="flex items-center justify-between border-b border-gray-200 p-6"
        >
            <h2 class="text-xl font-semibold">Tạo giai đoạn báo cáo mới</h2>
            <button
                onclick={onClose}
                class="rounded-lg p-2 transition-colors hover:bg-gray-100"
                aria-label="Close"
            >
                <X class="h-5 w-5" />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="space-y-6 p-6">
            <div>
                <label class="mb-2 block text-sm font-medium" for="title">
                    Tên giai đoạn <span class="text-red-500">*</span>
                </label>
                <input
                    id="title"
                    type="text"
                    bind:value={formData.title}
                    oninput={() => clearError("name")}
                    class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.name
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Ví dụ: Báo cáo tuần 1-2"
                />
                {#if errors.name}
                    <p class="mt-1 text-sm text-red-500">{errors.name}</p>
                {/if}
            </div>

            <div>
                <label class="mb-2 block text-sm font-medium" for="description">
                    Mô tả
                </label>
                <textarea
                    id="description"
                    bind:value={formData.description}
                    class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows="3"
                    placeholder="Mô tả ngắn về nội dung báo cáo trong giai đoạn này"
                ></textarea>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        class="mb-2 block text-sm font-medium"
                        for="startDate"
                    >
                        Ngày bắt đầu <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="startDate"
                        type="date"
                        bind:value={formData.startDate}
                        oninput={() => clearError("startDate")}
                        class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.startDate
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                    />
                    {#if errors.startDate}
                        <p class="mt-1 text-sm text-red-500">
                            {errors.startDate}
                        </p>
                    {/if}
                </div>

                <div>
                    <label class="mb-2 block text-sm font-medium" for="endDate">
                        Ngày kết thúc <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="endDate"
                        type="date"
                        bind:value={formData.endDate}
                        oninput={() => clearError("endDate")}
                        class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.endDate
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                    />
                    {#if errors.endDate}
                        <p class="mt-1 text-sm text-red-500">
                            {errors.endDate}
                        </p>
                    {/if}
                </div>
            </div>

            <div class="rounded-lg border border-blue-200 bg-blue-50 p-4">
                <p class="text-sm text-blue-800">
                    <strong>Lưu ý:</strong> Sau khi tạo giai đoạn báo cáo, sinh viên
                    sẽ có thể nộp báo cáo trong khoảng thời gian đã định.
                </p>
            </div>

            <div class="flex justify-end gap-3 border-t border-gray-200 pt-4">
                <button
                    type="button"
                    onclick={onClose}
                    class="rounded-lg border border-gray-300 px-6 py-2 transition-colors hover:bg-gray-50"
                >
                    Hủy
                </button>
                <button
                    type="submit"
                    class="rounded-lg bg-blue-600 px-6 py-2 text-white transition-colors hover:bg-blue-700"
                >
                    Tạo giai đoạn
                </button>
            </div>
        </form>
    </div>
</div>
