<script lang="ts">
    import { X, Plus } from "lucide-svelte";
    import type { Project } from "../../../../../types/project";
    // 1. Props
    let { onClose, onSubmit, categoryId, editingProject } = $props<{
        onClose: () => void;
        onSubmit: (project: any) => void;
        categoryId: string;
        editingProject?: Project;
    }>();

    // 2. Form State sử dụng Rune $state
    let formData = $state({
        categoryId,
        name: editingProject?.name || "",
        description: editingProject?.description || "",
        maxStudents: editingProject?.maxStudents || 3,
        minTeamMembers: editingProject?.minTeamMembers || 1,
        maxTeamMembers: editingProject?.maxTeamMembers || 3,
        instructor: editingProject?.instructor || "TS. Nguyễn Văn A",
        tags: [...(editingProject?.tags || [])],
        status: editingProject?.status || "available",
        quantity: 1,
    });

    let tagInput = $state("");
    let errors = $state<Record<string, string>>({});

    // 3. Logic xử lý Tags
    const handleAddTag = () => {
        const tag = tagInput.trim();
        if (tag && !formData.tags.includes(tag)) {
            formData.tags.push(tag); // Svelte 5 tự động nhận biết thay đổi khi push
            tagInput = "";
        }
    };

    const handleRemoveTag = (tagToRemove: string) => {
        formData.tags = formData.tags.filter((t) => t !== tagToRemove);
    };

    // 4. Logic xử lý Submit
    const handleSubmit = (e: Event) => {
        e.preventDefault();
        const newErrors: Record<string, string> = {};

        if (!formData.name.trim()) newErrors.name = "Vui lòng nhập tên đề tài";
        if (!formData.description.trim())
            newErrors.description = "Vui lòng nhập mô tả";
        if (formData.maxStudents < 1)
            newErrors.maxStudents = "Số lượng sinh viên phải lớn hơn 0";

        if (Object.keys(newErrors).length > 0) {
            errors = newErrors;
            return;
        }

        onSubmit({ ...formData });
    };

    // Xóa lỗi khi người dùng nhập
    function clearError(field: string) {
        if (errors[field]) delete errors[field];
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div
        class="mx-4 max-h-[90vh] w-full max-w-2xl overflow-y-auto rounded-lg bg-white shadow-xl"
    >
        <div
            class="sticky top-0 z-10 flex items-center justify-between border-b border-gray-200 bg-white p-6"
        >
            <h2 class="text-xl font-semibold">
                {editingProject ? "Chỉnh sửa đề tài" : "Thêm đề tài mới"}
            </h2>
            <button
                onclick={onClose}
                class="rounded-lg p-2 transition-colors hover:bg-gray-100"
            >
                <X class="h-5 w-5" />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="space-y-6 p-6">
            <div>
                <label class="mb-2 block text-sm font-medium" for="name">
                    Tên đề tài <span class="text-red-500">*</span>
                </label>
                <input
                    id="name"
                    type="text"
                    bind:value={formData.name}
                    oninput={() => clearError("name")}
                    class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.name
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Ví dụ: Hệ thống quản lý thư viện trực tuyến"
                />
                {#if errors.name}
                    <p class="mt-1 text-sm text-red-500">{errors.name}</p>
                {/if}
            </div>

            <div>
                <label class="mb-2 block text-sm font-medium" for="desc">
                    Mô tả đề tài <span class="text-red-500">*</span>
                </label>
                <textarea
                    id="desc"
                    bind:value={formData.description}
                    oninput={() => clearError("description")}
                    class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.description
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    rows="4"
                    placeholder="Mô tả chi tiết về đề tài, yêu cầu, công nghệ sử dụng..."
                ></textarea>
                {#if errors.description}
                    <p class="mt-1 text-sm text-red-500">
                        {errors.description}
                    </p>
                {/if}
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label class="mb-2 block text-sm font-medium" for="qty">
                        Số lượng đề tài {!editingProject ? "*" : ""}
                    </label>
                    <input
                        id="qty"
                        type="number"
                        bind:value={formData.quantity}
                        disabled={!!editingProject}
                        class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-50"
                    />
                    {#if !editingProject}
                        <p class="mt-1 text-xs text-gray-500">
                            Hệ thống sẽ tạo {formData.quantity} đề tài giống nhau
                        </p>
                    {/if}
                </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label class="mb-2 block text-sm font-medium" for="minT"
                        >Số thành viên tối thiểu</label
                    >
                    <input
                        id="minT"
                        type="number"
                        bind:value={formData.minTeamMembers}
                        max={formData.maxTeamMembers}
                        class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <div>
                    <label class="mb-2 block text-sm font-medium" for="maxT"
                        >Số thành viên tối đa</label
                    >
                    <input
                        id="maxT"
                        type="number"
                        bind:value={formData.maxTeamMembers}
                        min={formData.minTeamMembers}
                        class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
            </div>

            <div>
                <label class="mb-2 block text-sm font-medium" for="tags"
                    >Tags/Từ khóa</label
                >
                <div class="mb-2 flex gap-2">
                    <input
                        id="tags"
                        type="text"
                        bind:value={tagInput}
                        onkeydown={(e) =>
                            e.key === "Enter" &&
                            (e.preventDefault(), handleAddTag())}
                        class="flex-1 rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="Nhập tag và nhấn Enter"
                    />
                    <button
                        type="button"
                        onclick={handleAddTag}
                        class="rounded-lg bg-gray-100 px-4 py-2 transition-colors hover:bg-gray-200"
                    >
                        <Plus class="h-5 w-5" />
                    </button>
                </div>
                <div class="flex flex-wrap gap-2">
                    {#each formData.tags as tag (tag)}
                        <span
                            class="flex items-center gap-1 rounded-full bg-blue-50 px-3 py-1 text-sm text-blue-700"
                        >
                            {tag}
                            <button
                                type="button"
                                onclick={() => handleRemoveTag(tag)}
                                class="rounded-full p-0.5 hover:bg-blue-100"
                            >
                                <X class="h-3 w-3" />
                            </button>
                        </span>
                    {/each}
                </div>
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
                    class="rounded-lg bg-blue-600 px-6 py-2 text-white transition-colors hover:bg-blue-700 font-medium"
                >
                    {editingProject ? "Cập nhật" : "Thêm đề tài"}
                </button>
            </div>
        </form>
    </div>
</div>
