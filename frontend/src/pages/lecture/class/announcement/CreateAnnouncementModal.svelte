<script lang="ts">
    import { X, Paperclip, Trash2, Upload } from "lucide-svelte";
    import type { Announcement } from "../../../../types/announcement";
    import type {
        CreatePostRequest,
        PostResponse,
        UpdatePostRequest,
    } from "../../../../dtos/post-dto";

    let {
        onClose,
        onSubmit,
        editingAnnouncement = null,
    } = $props<{
        onClose: () => void;
        onSubmit: (post: UpdatePostRequest | CreatePostRequest) => void;
        editingAnnouncement?: PostResponse | null;
    }>();

    // Khởi tạo state của form
    let formData = $state({
        title: editingAnnouncement?.title || "",
        content: editingAnnouncement?.content || "",
        // Lưu trữ trực tiếp đối tượng File của trình duyệt
        files: [] as File[],
        existingFiles:
            editingAnnouncement?.attachments ||
            ([] as typeof editingAnnouncement.attachments),
        filesToRemove: [] as string[], // Track files to remove
    });

    let errors = $state<Record<string, string>>({});
    let fileInput: HTMLInputElement;

    // Helper: Định dạng kích thước file
    const formatFileSize = (bytes: number): string => {
        if (bytes === 0) return "0 Bytes";
        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return (
            Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i]
        );
    };

    // Xử lý chọn file
    function handleFileSelect(e: Event) {
        const target = e.target as HTMLInputElement;
        if (!target.files) return;

        // Thêm các file mới vào mảng hiện tại
        const selectedFiles = Array.from(target.files);
        formData.files = [...formData.files, ...selectedFiles];

        target.value = ""; // Reset input để có thể chọn lại cùng 1 file nếu muốn
    }

    // Xử lý xóa file (dựa vào tên và size hoặc index)
    function handleRemoveFile(index: number) {
        formData.files = formData.files.filter((_, i) => i !== index);
    }

    function handleRemoveExistingFile(fileUrl: string) {
        formData.existingFiles = formData.existingFiles.filter(
            (file: any) => file.fileUrl !== fileUrl,
        );
        formData.filesToRemove.push(fileUrl);

        // Hide the file from the UI
        const fileElement = document.querySelector(
            `[data-file-id="${fileUrl}"]`,
        );
        if (fileElement) {
            fileElement.classList.add("hidden");
        }
    }

    // Xử lý nộp form
    function handleSubmit(e: SubmitEvent) {
        e.preventDefault();
        const newErrors: Record<string, string> = {};

        if (!formData.title.trim())
            newErrors.title = "Vui lòng nhập tiêu đề thông báo";
        if (!formData.content.trim())
            newErrors.content = "Vui lòng nhập nội dung thông báo";

        if (Object.keys(newErrors).length > 0) {
            errors = newErrors;
            return;
        }

        // Gửi trực tiếp formData vì cấu trúc đã khớp với CreatePostRequest
        onSubmit({
            title: formData.title,
            content: formData.content,
            files: formData.files, // Đây đã là mảng File[]
            filesToRemove: formData.filesToRemove, // Include files to remove
        });
        console.log(formData);
    }

    function clearError(field: string) {
        if (errors[field]) delete errors[field];
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
>
    <div
        class="bg-white rounded-lg shadow-xl w-full max-w-3xl mx-4 max-h-[90vh] overflow-y-auto"
    >
        <div
            class="flex justify-between items-center p-6 border-b border-gray-200 sticky top-0 bg-white z-10"
        >
            <h2 class="text-xl font-semibold">
                {editingAnnouncement
                    ? "Chỉnh sửa thông báo"
                    : "Tạo thông báo mới"}
            </h2>
            <button
                onclick={onClose}
                class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
                <X class="w-5 h-5" />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="p-6 space-y-6">
            <div>
                <label class="block text-sm mb-2" for="title">
                    Tiêu đề thông báo <span class="text-red-500">*</span>
                </label>
                <input
                    id="title"
                    type="text"
                    bind:value={formData.title}
                    oninput={() => clearError("title")}
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.title
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    placeholder="Ví dụ: Thông báo về lịch nộp báo cáo giữa kỳ"
                />
                {#if errors.title}
                    <p class="text-red-500 text-sm mt-1">{errors.title}</p>
                {/if}
            </div>

            <div>
                <label class="block text-sm mb-2" for="content">
                    Nội dung thông báo <span class="text-red-500">*</span>
                </label>
                <textarea
                    id="content"
                    bind:value={formData.content}
                    oninput={() => clearError("content")}
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.content
                        ? 'border-red-500'
                        : 'border-gray-300'}"
                    rows={8}
                    placeholder="Nhập nội dung chi tiết thông báo..."
                ></textarea>
                {#if errors.content}
                    <p class="text-red-500 text-sm mt-1">{errors.content}</p>
                {/if}
            </div>

            <div>
                <label class="block text-sm mb-2" for="file-upload"
                    >File đính kèm</label
                >
                <input
                    bind:this={fileInput}
                    type="file"
                    multiple
                    onchange={handleFileSelect}
                    class="hidden"
                    id="file-upload"
                    accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.zip,.rar"
                />

                <button
                    type="button"
                    onclick={() => fileInput.click()}
                    class="flex items-center gap-2 px-4 py-2 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors w-full justify-center text-gray-600 hover:text-blue-600"
                >
                    <Upload class="w-5 h-5" />
                    <span>Chọn file để đính kèm</span>
                </button>

                {#if formData.files.length > 0}
                    <div class="mt-4 space-y-2">
                        <div class="text-sm font-medium">
                            File đã chọn ({formData.files.length}):
                        </div>
                        {#each formData.files as file, i}
                            <div
                                class="flex items-center gap-3 p-3 bg-gray-50 border border-gray-200 rounded-lg"
                            >
                                <Paperclip
                                    class="w-5 h-5 text-gray-600 flex-shrink-0"
                                />
                                <div class="flex-1 min-w-0">
                                    <div class="text-sm text-gray-900 truncate">
                                        {file.name}
                                    </div>
                                    <div class="text-xs text-gray-500">
                                        {formatFileSize(file.size)}
                                    </div>
                                </div>
                                <button
                                    type="button"
                                    onclick={() => handleRemoveFile(i)}
                                    class="p-1 hover:bg-red-50 rounded transition-colors flex-shrink-0"
                                >
                                    <Trash2 class="w-4 h-4 text-red-600" />
                                </button>
                            </div>
                        {/each}
                    </div>
                {/if}

                {#if formData.existingFiles.length > 0}
                    <div class="mt-4 space-y-2">
                        <div class="text-sm font-medium">File hiện có:</div>
                        {#each formData.existingFiles as file}
                            <div
                                class="flex items-center gap-3 p-3 bg-gray-50 border border-gray-200 rounded-lg"
                                data-file-id={file.publicId}
                            >
                                <Paperclip
                                    class="w-5 h-5 text-gray-600 flex-shrink-0"
                                />
                                <div class="flex-1 min-w-0">
                                    <div class="text-sm text-gray-900 truncate">
                                        {file.fileName}
                                    </div>
                                    <div class="text-xs text-gray-500">
                                        {file.mimeType}
                                    </div>
                                </div>
                                <button
                                    type="button"
                                    onclick={() =>
                                        handleRemoveExistingFile(file.fileUrl)}
                                    class="p-1 hover:bg-red-50 rounded transition-colors flex-shrink-0"
                                >
                                    <Trash2 class="w-4 h-4 text-red-600" />
                                </button>
                            </div>
                        {/each}
                    </div>
                {/if}
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
                    {editingAnnouncement
                        ? "Cập nhật thông báo"
                        : "Đăng thông báo"}
                </button>
            </div>
        </form>
    </div>
</div>
