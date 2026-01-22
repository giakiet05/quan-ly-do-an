<script lang="ts">
    import {
        X,
        FileText,
        FileCode,
        Archive,
        Calendar,
        Info,
        Check,
    } from "lucide-svelte";
    import type { CreatePeriodRequest } from "../../../../../dtos/period-dto";

    let { onClose, onSubmit, totalStudents, maxStudents } = $props<{
        onClose: () => void;
        onSubmit: (data: CreatePeriodRequest) => void;
        totalStudents: number;
        maxStudents: number;
    }>();

    let formData = $state({
        title: "",
        description: "",
        startDate: "",
        endDate: "",
        fileType: [] as string[],
    });

    let errors = $state<Record<string, string>>({});

    // Danh sách file type phong phú hơn, phân nhóm để người dùng dễ chọn
    const suggestedTypes = [
        {
            label: "Tài liệu",
            types: [
                {
                    v: ".pdf",
                    l: "PDF",
                    color: "bg-red-100 text-red-700 border-red-200",
                },
                {
                    v: ".docx",
                    l: "Word",
                    color: "bg-blue-100 text-blue-700 border-blue-200",
                },
                {
                    v: ".pptx",
                    l: "PowerPoint",
                    color: "bg-orange-100 text-orange-700 border-orange-200",
                },
                {
                    v: ".xlsx",
                    l: "Excel",
                    color: "bg-green-100 text-green-700 border-green-200",
                },
                {
                    v: ".txt",
                    l: "Text",
                    color: "bg-gray-100 text-gray-700 border-gray-200",
                },
            ],
        },
        {
            label: "Lập trình/Dữ liệu",
            types: [
                {
                    v: ".zip",
                    l: "ZIP/RAR",
                    color: "bg-purple-100 text-purple-700 border-purple-200",
                },
                {
                    v: ".sql",
                    l: "SQL",
                    color: "bg-indigo-100 text-indigo-700 border-indigo-200",
                },
                {
                    v: ".json",
                    l: "JSON",
                    color: "bg-yellow-100 text-yellow-700 border-yellow-200",
                },
                {
                    v: ".csv",
                    l: "CSV",
                    color: "bg-emerald-100 text-emerald-700 border-emerald-200",
                },
            ],
        },
    ];

    const formatDate = (date: string) => {
        if (!date) return "";
        const d = new Date(date);
        return d.toISOString().split("T")[0];
    };

    const toggleFileType = (val: string) => {
        if (formData.fileType.includes(val)) {
            formData.fileType = formData.fileType.filter((t) => t !== val);
        } else {
            formData.fileType = [...formData.fileType, val];
        }
        if (errors.fileType) delete errors.fileType;
    };

    const handleSubmit = (e: Event) => {
        e.preventDefault();
        const newErrors: Record<string, string> = {};

        if (!formData.title.trim())
            newErrors.title = "Vui lòng nhập tên giai đoạn";
        if (!formData.startDate)
            newErrors.startDate = "Vui lòng chọn ngày bắt đầu";
        if (!formData.endDate)
            newErrors.endDate = "Vui lòng chọn ngày kết thúc";
        if (formData.fileType.length === 0)
            newErrors.fileType = "Vui lòng chọn ít nhất một định dạng tệp";

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
            startDate: formatDate(formData.startDate),
            endDate: formatDate(formData.endDate),
        });
    };
</script>

<div
    class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4"
>
    <div
        class="w-full max-w-2xl overflow-hidden rounded-2xl bg-white shadow-2xl animate-in fade-in zoom-in duration-200"
    >
        <div
            class="flex items-center justify-between border-b border-gray-100 bg-gray-50/50 px-8 py-5"
        >
            <div>
                <h2 class="text-xl font-bold text-gray-800">
                    Cấu hình giai đoạn báo cáo
                </h2>
                <p class="text-sm text-gray-500">
                    Thiết lập thời gian và yêu cầu nộp bài
                </p>
            </div>
            <button
                onclick={onClose}
                class="rounded-full p-2 text-gray-400 transition-all hover:bg-white hover:text-gray-600 hover:shadow-sm"
            >
                <X class="h-6 w-6" />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="max-h-[80vh] overflow-y-auto p-8">
            <div class="space-y-6">
                <div class="space-y-2">
                    <label
                        class="text-sm font-semibold text-gray-700"
                        for="title"
                        >Tên giai đoạn <span class="text-red-500">*</span
                        ></label
                    >
                    <input
                        id="title"
                        type="text"
                        bind:value={formData.title}
                        placeholder="Ví dụ: Báo cáo giữa kỳ, Đồ án chính thức..."
                        class="w-full rounded-xl border px-4 py-3 transition-all focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 {errors.title
                            ? 'border-red-500 bg-red-50'
                            : 'border-gray-200 hover:border-gray-300'}"
                    />
                    {#if errors.title}<p
                            class="text-xs font-medium text-red-500"
                        >
                            {errors.title}
                        </p>{/if}
                </div>

                <div class="space-y-2">
                    <label
                        class="text-sm font-semibold text-gray-700"
                        for="description">Yêu cầu cụ thể</label
                    >
                    <textarea
                        id="description"
                        bind:value={formData.description}
                        rows="2"
                        placeholder="Hướng dẫn sinh viên những gì cần nộp trong giai đoạn này..."
                        class="w-full rounded-xl border border-gray-200 px-4 py-3 transition-all focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 hover:border-gray-300"
                    ></textarea>
                </div>

                <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                    <div class="space-y-2">
                        <label
                            class="flex items-center gap-2 text-sm font-semibold text-gray-700"
                            for="startDate"
                        >
                            <Calendar class="h-4 w-4 text-blue-500" /> Ngày bắt đầu
                        </label>
                        <input
                            id="startDate"
                            type="date"
                            bind:value={formData.startDate}
                            class="w-full rounded-xl border px-4 py-3 transition-all focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 {errors.startDate
                                ? 'border-red-500'
                                : 'border-gray-200'}"
                        />
                    </div>
                    <div class="space-y-2">
                        <label
                            class="flex items-center gap-2 text-sm font-semibold text-gray-700"
                            for="endDate"
                        >
                            <Calendar class="h-4 w-4 text-red-500" /> Ngày kết thúc
                        </label>
                        <input
                            id="endDate"
                            type="date"
                            bind:value={formData.endDate}
                            class="w-full rounded-xl border px-4 py-3 transition-all focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 {errors.endDate
                                ? 'border-red-500'
                                : 'border-gray-200'}"
                        />
                    </div>
                </div>

                <div class="space-y-3">
                    <label
                        class="flex items-center justify-between text-sm font-semibold text-gray-700"
                    >
                        <span
                            >Định dạng tệp cho phép <span class="text-red-500"
                                >*</span
                            ></span
                        >
                        <span class="text-xs font-normal text-gray-400"
                            >Đã chọn: {formData.fileType.length}</span
                        >
                    </label>

                    <div
                        class="rounded-2xl border border-gray-100 bg-gray-50/30 p-4 space-y-4"
                    >
                        {#each suggestedTypes as group}
                            <div class="space-y-2">
                                <h4
                                    class="text-[10px] font-bold uppercase tracking-wider text-gray-400"
                                >
                                    {group.label}
                                </h4>
                                <div class="flex flex-wrap gap-2">
                                    {#each group.types as type}
                                        <button
                                            type="button"
                                            onclick={() =>
                                                toggleFileType(type.v)}
                                            class="flex items-center gap-2 rounded-lg border px-3 py-2 text-sm font-medium transition-all
                                            {formData.fileType.includes(type.v)
                                                ? 'border-blue-500 bg-blue-600 text-white shadow-md shadow-blue-200 scale-105'
                                                : `bg-white ${type.color} hover:scale-105 hover:shadow-sm`}"
                                        >
                                            {#if formData.fileType.includes(type.v)}
                                                <Check class="h-3 w-3" />
                                            {/if}
                                            {type.l}
                                        </button>
                                    {/each}
                                </div>
                            </div>
                        {/each}
                    </div>
                    {#if errors.fileType}<p
                            class="text-xs font-medium text-red-500"
                        >
                            {errors.fileType}
                        </p>{/if}
                </div>

                <div
                    class="flex gap-3 rounded-xl bg-amber-50 border border-amber-100 p-4"
                >
                    <Info class="h-5 w-5 shrink-0 text-amber-600" />
                    <p class="text-xs leading-relaxed text-amber-800">
                        Hệ thống sẽ tự động khóa nộp bài khi hết hạn. Giảng viên
                        có thể gia hạn thêm thời gian sau khi giai đoạn kết thúc
                        nếu cần thiết.
                    </p>
                </div>
            </div>

            <div
                class="mt-8 flex items-center justify-end gap-3 border-t border-gray-100 pt-6"
            >
                <button
                    type="button"
                    onclick={onClose}
                    class="rounded-xl px-6 py-3 text-sm font-bold text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700"
                >
                    Hủy bỏ
                </button>
                <button
                    type="submit"
                    class="rounded-xl bg-blue-600 px-8 py-3 text-sm font-bold text-white shadow-lg shadow-blue-200 transition-all hover:bg-blue-700 hover:shadow-blue-300 active:scale-95"
                >
                    Xác nhận tạo
                </button>
            </div>
        </form>
    </div>
</div>

<style>
    /* Tuỳ chỉnh thanh cuộn cho đẹp hơn trên Chrome/Safari */
    form::-webkit-scrollbar {
        width: 6px;
    }
    form::-webkit-scrollbar-track {
        background: transparent;
    }
    form::-webkit-scrollbar-thumb {
        background: #e2e8f0;
        border-radius: 10px;
    }
    form::-webkit-scrollbar-thumb:hover {
        background: #cbd5e1;
    }
</style>
