<script lang="ts">
    import { X } from "lucide-svelte";
    import type { ClassData } from "../../../../types/class";
    import type { ProjectRound } from "../../../../types/project-round";

    // 1. Định nghĩa Props
    let {
        onClose,
        onSubmit,
        classData,
        isEdit = false,
        initialData = null,
    } = $props<{
        onClose: () => void;
        onSubmit: (data: any) => void;
        classData: ClassData | undefined;
        isEdit?: boolean;
        initialData?: ProjectRound | null;
    }>();

    // 2. Helper function để định dạng date input (YYYY-MM-DD)
    const formatToDateInput = (dateStr: string | undefined | null | Date) => {
        if (!dateStr) return "";
        const dateString =
            typeof dateStr === "string"
                ? dateStr
                : new Date(dateStr).toISOString();
        return dateString.split("T")[0];
    };

    let formData = $state({
        name: initialData?.name ?? "",
        description: initialData?.description ?? "",
        startDate: formatToDateInput(initialData?.startDate),
        endDate: formatToDateInput(initialData?.endDate),
        status: initialData?.status ?? "upcoming",
    });

    let errors = $state<Record<string, string>>({});

    // 4. Xử lý logic Submit
    function handleSubmit(e: Event) {
        e.preventDefault();
        const newErrors: Record<string, string> = {};

        if (!formData.name.trim())
            newErrors.name = "Vui lòng nhập tên hạng mục";
        if (!formData.startDate)
            newErrors.startDate = "Vui lòng chọn ngày bắt đầu";
        if (!formData.endDate)
            newErrors.endDate = "Vui lòng chọn ngày kết thúc";

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

        // Trả về dữ liệu kèm ID nếu là đang edit
        onSubmit({
            ...formData,
            id: initialData?.id, // Giữ lại ID để API biết là update
        });
    }

    function clearError(field: string) {
        if (errors[field]) delete errors[field];
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4">
        <div
            class="flex justify-between items-center p-6 border-b border-gray-200"
        >
            <h2 class="text-xl font-semibold">
                {isEdit ? "Chỉnh sửa hạng mục" : "Tạo hạng mục đề tài mới"}
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
                    {classData?.name} - {classData?.semester}
                </p>
            </div>

            <div>
                <label class="block text-sm font-medium mb-2" for="name">
                    Tên hạng mục <span class="text-red-500">*</span>
                </label>
                <input
                    id="name"
                    type="text"
                    bind:value={formData.name}
                    oninput={() => clearError("name")}
                    class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none {errors.name
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Ví dụ: Hạng mục 1 - Đề tài Web Application"
                />
                {#if errors.name}
                    <p class="text-red-500 text-sm mt-1">{errors.name}</p>
                {/if}
            </div>

            <div>
                <label class="block text-sm font-medium mb-2" for="description"
                    >Mô tả</label
                >
                <textarea
                    id="description"
                    bind:value={formData.description}
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                    rows={3}
                    placeholder="Mô tả ngắn về hạng mục này"
                ></textarea>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        class="block text-sm font-medium mb-2"
                        for="startDate"
                    >
                        Ngày bắt đầu <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="startDate"
                        type="date"
                        bind:value={formData.startDate}
                        oninput={() => clearError("startDate")}
                        class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none {errors.startDate
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
                    <label class="block text-sm font-medium mb-2" for="endDate">
                        Ngày kết thúc <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="endDate"
                        type="date"
                        bind:value={formData.endDate}
                        oninput={() => clearError("endDate")}
                        class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none {errors.endDate
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
                    {isEdit ? "Cập nhật" : "Tạo hạng mục"}
                </button>
            </div>
        </form>
    </div>
</div>
