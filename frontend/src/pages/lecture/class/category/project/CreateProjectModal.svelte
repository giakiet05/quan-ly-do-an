<script lang="ts">
    import { X, Plus } from "lucide-svelte";
    import type { UpdateProjectRequest } from "../../../../../dtos/project-dto";
    // 1. Props
    let { onClose, onSubmit, projectRoundId, editingProject } = $props<{
        onClose: () => void;
        onSubmit: (project: any) => void;
        projectRoundId: string;
        editingProject?: UpdateProjectRequest;
    }>();

    // 2. Form State sử dụng Rune $state
    let formData = $derived({
        projectRoundId,
        title: editingProject?.title || "",
        amount: editingProject?.amount || 1, // Map to "amount" field
        description: editingProject?.description || "",
        minMember: editingProject?.minMember || 1, // Map to "minMember"
        maxMember: editingProject?.maxMember || 3, // Map to "maxMember"
        tags: [...(editingProject?.tags || [])],
        status: editingProject?.status || "available",
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

        if (!formData.title.trim()) newErrors.title = "Vui lòng nhập tiêu đề";
        if (formData.amount < 1) newErrors.amount = "Số lượng phải lớn hơn 0"; // Validate "amount"
        if (!formData.description.trim())
            newErrors.description = "Vui lòng nhập mô tả";
        if (formData.minMember < 1)
            newErrors.minMember = "Số thành viên tối thiểu phải lớn hơn 0";
        if (formData.maxMember < formData.minMember)
            newErrors.maxMember =
                "Số thành viên tối đa phải lớn hơn hoặc bằng tối thiểu";

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
                <label class="mb-2 block text-sm font-medium" for="title">
                    Tiêu đề <span class="text-red-500">*</span>
                </label>
                <input
                    id="title"
                    type="text"
                    bind:value={formData.title}
                    oninput={() => clearError("title")}
                    class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.title
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Ví dụ: Đề tài nghiên cứu AI"
                />
                {#if errors.title}
                    <p class="mt-1 text-sm text-red-500">{errors.title}</p>
                {/if}
            </div>

            <div>
                <label class="mb-2 block text-sm font-medium" for="amount">
                    Số lượng <span class="text-red-500">*</span>
                </label>
                <input
                    id="amount"
                    type="number"
                    bind:value={formData.amount}
                    oninput={() => clearError("amount")}
                    class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.amount
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Nhập số lượng đề tài"
                />
                {#if errors.amount}
                    <p class="mt-1 text-sm text-red-500">{errors.amount}</p>
                {/if}
            </div>

            <div>
                <label class="mb-2 block text-sm font-medium" for="desc">
                    Mô tả <span class="text-red-500">*</span>
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
                    <label
                        class="mb-2 block text-sm font-medium"
                        for="minMember"
                    >
                        Số thành viên tối thiểu <span class="text-red-500"
                            >*</span
                        >
                    </label>
                    <input
                        id="minMember"
                        type="number"
                        bind:value={formData.minMember}
                        oninput={() => clearError("minMember")}
                        class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.minMember
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                        placeholder="Nhập số thành viên tối thiểu"
                    />
                    {#if errors.minMember}
                        <p class="mt-1 text-sm text-red-500">
                            {errors.minMember}
                        </p>
                    {/if}
                </div>
                <div>
                    <label
                        class="mb-2 block text-sm font-medium"
                        for="maxMember"
                    >
                        Số thành viên tối đa <span class="text-red-500">*</span>
                    </label>
                    <input
                        id="maxMember"
                        type="number"
                        bind:value={formData.maxMember}
                        oninput={() => clearError("maxMember")}
                        class="w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.maxMember
                            ? 'border-red-500'
                            : 'border-gray-300'}"
                        placeholder="Nhập số thành viên tối đa"
                    />
                    {#if errors.maxMember}
                        <p class="mt-1 text-sm text-red-500">
                            {errors.maxMember}
                        </p>
                    {/if}
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
