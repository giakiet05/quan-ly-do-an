<script lang="ts">
    import { X } from "lucide-svelte";
    import type { ClassData } from "../../../../types/class";
    import type { ProjectRound } from "../../../../types/project-round";

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

    const formatToDateInput = (dateStr: string | undefined | null | Date) => {
        if (!dateStr) return "";
        const dateString =
            typeof dateStr === "string"
                ? dateStr
                : new Date(dateStr).toISOString();
        return dateString.split("T")[0];
    };

    // 1. Thêm minStudents và maxStudents vào formData
    let formData = $state({
        name: initialData?.name ?? "",
        description: initialData?.description ?? "",
        startDate: formatToDateInput(initialData?.startDate),
        endDate: formatToDateInput(initialData?.endDate),
        minStudents: initialData?.minStudents ?? 1, // Mặc định là 1
        maxStudents: initialData?.maxStudents ?? 3, // Mặc định là 3
        status: initialData?.status ?? "upcoming",
    });

    let errors = $state<Record<string, string>>({});

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

        // 2. Logic kiểm tra lỗi cho số lượng sinh viên
        if (formData.minStudents < 1) {
            newErrors.minStudents = "Tối thiểu phải có 1 sinh viên";
        }
        if (formData.maxStudents < formData.minStudents) {
            newErrors.maxStudents =
                "Số lượng tối đa không được nhỏ hơn tối thiểu";
        }

        if (Object.keys(newErrors).length > 0) {
            errors = newErrors;
            return;
        }

        onSubmit({
            ...formData,
            id: initialData?.id,
        });
    }

    function clearError(field: string) {
        if (errors[field]) delete errors[field];
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div
        class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto"
    >
        <div
            class="flex justify-between items-center p-6 border-b border-gray-200 sticky top-0 bg-white z-10"
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
                />
                {#if errors.name}
                    <p class="text-red-500 text-sm mt-1">{errors.name}</p>
                {/if}
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        class="block text-sm font-medium mb-2"
                        for="minStudents"
                    >
                        Sinh viên tối thiểu/nhóm <span class="text-red-500"
                            >*</span
                        >
                    </label>
                    <input
                        id="minStudents"
                        type="number"
                        min="1"
                        bind:value={formData.minStudents}
                        oninput={() => clearError("minStudents")}
                        class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none {errors.minStudents
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                    />
                    {#if errors.minStudents}
                        <p class="text-red-500 text-sm mt-1">
                            {errors.minStudents}
                        </p>
                    {/if}
                </div>
                <div>
                    <label
                        class="block text-sm font-medium mb-2"
                        for="maxStudents"
                    >
                        Sinh viên tối đa/nhóm <span class="text-red-500">*</span
                        >
                    </label>
                    <input
                        id="maxStudents"
                        type="number"
                        min="1"
                        bind:value={formData.maxStudents}
                        oninput={() => clearError("maxStudents")}
                        class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none {errors.maxStudents
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                    />
                    {#if errors.maxStudents}
                        <p class="text-red-500 text-sm mt-1">
                            {errors.maxStudents}
                        </p>
                    {/if}
                </div>
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

            <div>
                <label class="block text-sm font-medium mb-2" for="description"
                    >Mô tả</label
                >
                <textarea
                    id="description"
                    bind:value={formData.description}
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                    rows={2}
                ></textarea>
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
