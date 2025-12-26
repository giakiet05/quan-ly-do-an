<script lang="ts">
    import { X, Paperclip, Trash2, Upload } from "lucide-svelte";
    import type {
        Announcement,
        AttachedFile,
    } from "../../../../types/announcement";

    // Định nghĩa Props bằng $props()
    let { onClose, onSubmit, editingAnnouncement } = $props<{
        onClose: () => void;
        onSubmit: (announcement: any) => void;
        editingAnnouncement?: Announcement;
    }>();

    // Khởi tạo state của form
    let formData = $state({
        title: editingAnnouncement?.title || "",
        content: editingAnnouncement?.content || "",
        files: editingAnnouncement?.files || ([] as AttachedFile[]),
    });

    let errors = $state<Record<string, string>>({});
    let fileInput: HTMLInputElement; // Thay thế cho useRef

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
        const files = target.files;
        if (!files || files.length === 0) return;

        const newFiles: AttachedFile[] = Array.from(files).map((file) => ({
            id: `${Date.now()}-${Math.random()}`,
            name: file.name,
            size: file.size,
            url: URL.createObjectURL(file),
        }));

        // Cập nhật mảng trực tiếp nhờ proxy của Svelte 5
        formData.files = [...formData.files, ...newFiles];
        target.value = ""; // Reset input
    }

    // Xử lý xóa file
    function handleRemoveFile(fileId: string) {
        formData.files = formData.files.filter(
            (f: AttachedFile) => f.id !== fileId,
        );
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

        onSubmit({ ...formData });
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
                <p class="text-sm text-gray-500 mt-1">
                    Bạn có thể nhập văn bản nhiều dòng. Xuống dòng sẽ được giữ
                    nguyên khi hiển thị.
                </p>
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

                <p class="text-sm text-gray-500 mt-2">
                    Hỗ trợ: PDF, Word, Excel, PowerPoint, Text, ZIP, RAR
                </p>

                {#if formData.files.length > 0}
                    <div class="mt-4 space-y-2">
                        <div class="text-sm font-medium">
                            File đã chọn ({formData.files.length}):
                        </div>
                        {#each formData.files as file (file.id)}
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
                                    onclick={() => handleRemoveFile(file.id)}
                                    class="p-1 hover:bg-red-50 rounded transition-colors flex-shrink-0"
                                >
                                    <Trash2 class="w-4 h-4 text-red-600" />
                                </button>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>

            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p class="text-sm text-blue-800">
                    <strong>Lưu ý:</strong> Thông báo sẽ được hiển thị cho tất cả
                    sinh viên trong lớp.
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
                    {editingAnnouncement
                        ? "Cập nhật thông báo"
                        : "Đăng thông báo"}
                </button>
            </div>
        </form>
    </div>
</div>
