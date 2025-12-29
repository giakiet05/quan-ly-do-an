<script lang="ts">
    import {
        Plus,
        Calendar,
        Paperclip,
        Edit,
        Trash2,
        ChevronDown,
        ChevronUp,
    } from "@lucide/svelte";
    import CreateAnnouncementModal from "./CreateAnnouncementModal.svelte";
    import type { Announcement } from "../../../../types/announcement";

    // Định nghĩa Props bằng $props()
    let { onOpen, onCLose, onEdit } = $props<{
        onOpen: () => void;
        onCLose: () => void;
        onEdit: (announcement: Announcement) => void;
    }>();
    // State sử dụng Runes
    let showCreateModal = $state(false);
    let editingAnnouncement = $state<Announcement | null>(null);
    let expandedId = $state<string | null>(null);

    let announcements = $state<Announcement[]>([
        {
            id: "1",
            title: "Thông báo về lịch nộp báo cáo giữa kỳ",
            content:
                "Các nhóm sinh viên lưu ý nộp báo cáo giữa kỳ trước ngày 15/03/2024. Báo cáo cần bao gồm:\n\n1. Phân tích yêu cầu hệ thống\n2. Thiết kế cơ sở dữ liệu\n3. Giao diện mockup\n4. Tiến độ thực hiện\n\nNộp qua email hoặc upload lên hệ thống.",
            createdAt: "2024-03-01T10:00:00",
            files: [
                {
                    id: "f1",
                    name: "Mau_bao_cao_giua_ky.docx",
                    size: 245760,
                    url: "#",
                },
                {
                    id: "f2",
                    name: "Huong_dan_viet_bao_cao.pdf",
                    size: 1048576,
                    url: "#",
                },
            ],
            author: "TS. Nguyễn Văn A",
        },
        // ... các thông báo khác
    ]);

    // Helper: Format dung lượng file
    const formatFileSize = (bytes: number): string => {
        if (bytes === 0) return "0 Bytes";
        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return (
            Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i]
        );
    };

    // Helper: Format ngày tháng (Svelte có thể dùng trực tiếp trong template)
    const formatDate = (dateString: string): string => {
        const date = new Date(dateString);
        const now = new Date();
        const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60);

        if (diffInHours < 24) {
            const diffInMinutes = Math.floor(diffInHours * 60);
            if (diffInMinutes < 60) return `${diffInMinutes} phút trước`;
            return `${Math.floor(diffInHours)} giờ trước`;
        } else if (diffInHours < 48) {
            return "Hôm qua";
        }

        return date.toLocaleDateString("vi-VN", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    // Logic xử lý
    function handleCreateAnnouncement(newAnnouncement: any) {
        const announcement: Announcement = {
            ...newAnnouncement,
            id: Date.now().toString(),
            createdAt: new Date().toISOString(),
            author: "TS. Nguyễn Văn A",
        };
        announcements = [announcement, ...announcements];
        showCreateModal = false;
    }

    function handleUpdateAnnouncement(updatedData: any) {
        if (!editingAnnouncement) return;

        const index = announcements.findIndex(
            (a) => a.id === editingAnnouncement?.id,
        );
        if (index !== -1) {
            announcements[index] = {
                ...announcements[index],
                ...updatedData,
                updatedAt: new Date().toISOString(),
            };
        }
        editingAnnouncement = null;
    }

    function handleDeleteAnnouncement(id: string) {
        if (confirm("Bạn có chắc chắn muốn xóa thông báo này?")) {
            announcements = announcements.filter((a) => a.id !== id);
            if (expandedId === id) expandedId = null;
        }
    }

    function toggleExpand(id: string) {
        expandedId = expandedId === id ? null : id;
    }
</script>

<div class="bg-white rounded-lg shadow-sm p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-semibold">Thông báo lớp học</h2>
        <button
            onclick={onOpen}
            class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
            <Plus class="w-5 h-5" />
            Tạo thông báo mới
        </button>
    </div>

    {#if announcements.length === 0}
        <div class="text-center py-16">
            <button
                onclick={() => (showCreateModal = true)}
                class="inline-flex items-center justify-center w-16 h-16 rounded-full border-2 border-dashed border-gray-300 hover:border-blue-500 hover:bg-blue-50 transition-colors group"
            >
                <Plus class="w-8 h-8 text-gray-400 group-hover:text-blue-500" />
            </button>
            <p class="mt-4 text-gray-600">Chưa có thông báo nào</p>
            <p class="text-sm text-gray-500">
                Nhấn vào dấu + để tạo thông báo mới
            </p>
        </div>
    {:else}
        <div class="space-y-4">
            {#each announcements as announcement (announcement.id)}
                {@const isExpanded = expandedId === announcement.id}

                <div
                    class="border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
                >
                    <div class="p-5">
                        <div class="flex justify-between items-start mb-3">
                            <div class="flex-1">
                                <h3 class="text-lg font-medium mb-2">
                                    {announcement.title}
                                </h3>
                                <div
                                    class="flex items-center gap-4 text-sm text-gray-500"
                                >
                                    <div class="flex items-center gap-1">
                                        <Calendar class="w-4 h-4" />
                                        <span
                                            >{formatDate(
                                                announcement.createdAt,
                                            )}</span
                                        >
                                    </div>
                                    {#if announcement.updatedAt}
                                        <span class="text-gray-400"
                                            >(Đã chỉnh sửa)</span
                                        >
                                    {/if}
                                    <span>•</span>
                                    <span>{announcement.author}</span>
                                    {#if announcement.files.length > 0}
                                        <span>•</span>
                                        <div class="flex items-center gap-1">
                                            <Paperclip class="w-4 h-4" />
                                            <span
                                                >{announcement.files.length} file</span
                                            >
                                        </div>
                                    {/if}
                                </div>
                            </div>

                            <div class="flex items-center gap-2 ml-4">
                                <button
                                    onclick={() =>
                                        toggleExpand(announcement.id)}
                                    class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                                    aria-label={isExpanded
                                        ? "Thu gọn"
                                        : "Xem chi tiết"}
                                >
                                    {#if isExpanded}
                                        <ChevronUp
                                            class="w-5 h-5 text-gray-600"
                                        />
                                    {:else}
                                        <ChevronDown
                                            class="w-5 h-5 text-gray-600"
                                        />
                                    {/if}
                                </button>
                                <button
                                    onclick={() => onEdit(announcement)}
                                    class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                                >
                                    <Edit class="w-5 h-5 text-gray-600" />
                                </button>
                                <button
                                    onclick={() =>
                                        handleDeleteAnnouncement(
                                            announcement.id,
                                        )}
                                    class="p-2 hover:bg-red-50 rounded-lg transition-colors"
                                >
                                    <Trash2 class="w-5 h-5 text-red-600" />
                                </button>
                            </div>
                        </div>

                        {#if !isExpanded}
                            <p class="text-gray-600 text-sm line-clamp-2">
                                {announcement.content}
                            </p>
                        {/if}
                    </div>

                    {#if isExpanded}
                        <div class="px-5 pb-5 border-t border-gray-200 pt-4">
                            <div class="mb-4">
                                <h4 class="text-sm font-semibold mb-2">
                                    Nội dung:
                                </h4>
                                <div
                                    class="text-gray-700 whitespace-pre-wrap bg-gray-50 rounded-lg p-4"
                                >
                                    {announcement.content}
                                </div>
                            </div>

                            {#if announcement.files.length > 0}
                                <div>
                                    <h4 class="text-sm font-semibold mb-3">
                                        File đính kèm:
                                    </h4>
                                    <div
                                        class="grid grid-cols-1 sm:grid-cols-2 gap-2"
                                    >
                                        {#each announcement.files as file (file.id)}
                                            <a
                                                href={file.url}
                                                class="flex items-center gap-3 p-3 bg-blue-50 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors group"
                                            >
                                                <Paperclip
                                                    class="w-5 h-5 text-blue-600"
                                                />
                                                <div class="flex-1 min-w-0">
                                                    <div
                                                        class="text-sm text-blue-900 group-hover:underline truncate"
                                                    >
                                                        {file.name}
                                                    </div>
                                                    <div
                                                        class="text-xs text-blue-600"
                                                    >
                                                        {formatFileSize(
                                                            file.size,
                                                        )}
                                                    </div>
                                                </div>
                                            </a>
                                        {/each}
                                    </div>
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</div>
