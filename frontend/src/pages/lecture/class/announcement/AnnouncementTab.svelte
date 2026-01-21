<script lang="ts">
    import { onMount } from "svelte";
    import {
        Plus,
        Calendar,
        Edit,
        Trash2,
        ChevronDown,
        ChevronUp,
        CheckCheck,
    } from "@lucide/svelte";
    import { notificationStore } from "../../../../stores/notification-store";

    // 1. Props & State
    let { onOpen, onEdit, id } = $props<{
        onOpen: () => void;
        onEdit: (announcement: any) => void;
        id: string;
    }>();

    let expandedId = $state<string | null>(null);

    // 2. Kết nối Store
    const {
        notifications,
        loading,
        fetchClassroomNotifications,
        markAsRead,
        markAllAsRead,
        remove,
    } = notificationStore;

    // 3. Khởi tạo dữ liệu
    onMount(() => {
        if (id) {
            fetchClassroomNotifications(id);
        }
    });

    // 4. Helpers
    const formatDate = (dateString: any): string => {
        const date = new Date(dateString);
        return date.toLocaleDateString("vi-VN", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    function toggleExpand(id: string) {
        expandedId = expandedId === id ? null : id;
        const item = $notifications.find((n) => n.id === id);
        if (item && !item.read) {
            markAsRead(id);
        }
    }
</script>

<div
    class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden"
>
    <div
        class="p-6 border-b border-gray-50 flex justify-between items-center bg-white"
    >
        <div>
            <h2 class="text-xl font-bold text-gray-800">Thông báo lớp học</h2>
            <p class="text-sm text-gray-500">Cập nhật tin tức mới nhất</p>
        </div>

        <div class="flex gap-2">
            <button
                onclick={() => markAllAsRead()}
                class="flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 rounded-lg transition-colors border border-gray-200"
            >
                <CheckCheck class="w-4 h-4" />
                Đọc tất cả
            </button>
            <button
                onclick={onOpen}
                class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 shadow-sm transition-all active:scale-95"
            >
                <Plus class="w-5 h-5" />
                Tạo mới
            </button>
        </div>
    </div>

    {#if $loading}
        <div class="p-12 flex justify-center items-center">
            <div
                class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
            ></div>
        </div>
    {:else if $notifications.length === 0}
        <div class="text-center py-20 bg-gray-50/50">
            <div
                class="bg-white w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4 shadow-sm"
            >
                <Plus class="w-8 h-8 text-gray-300" />
            </div>
            <p class="text-gray-600 font-medium">Chưa có thông báo nào</p>
        </div>
    {:else}
        <div class="divide-y divide-gray-100">
            {#each $notifications as announcement (announcement.id)}
                {@const isExpanded = expandedId === announcement.id}

                <div
                    class="group transition-colors {announcement.read
                        ? 'bg-white'
                        : 'bg-blue-50/30'}"
                >
                    <div class="p-5">
                        <div class="flex justify-between items-start">
                            <div
                                class="flex-1 min-w-0 cursor-pointer"
                                onclick={() => toggleExpand(announcement.id)}
                            >
                                <div class="flex items-center gap-2 mb-1">
                                    {#if !announcement.read}
                                        <span
                                            class="w-2.5 h-2.5 bg-blue-600 rounded-full"
                                        ></span>
                                    {/if}
                                    <h3
                                        class="text-lg font-semibold text-gray-900 truncate group-hover:text-blue-600 transition-colors"
                                    >
                                        {announcement.title}
                                    </h3>
                                </div>

                                <div
                                    class="flex items-center gap-4 text-sm text-gray-500"
                                >
                                    <div class="flex items-center gap-1.5">
                                        <Calendar class="w-4 h-4" />
                                        <span
                                            >{formatDate(
                                                announcement.createdAt,
                                            )}</span
                                        >
                                    </div>
                                </div>
                            </div>

                            <div class="flex items-center gap-1 ml-4">
                                <button
                                    onclick={() => onEdit(announcement)}
                                    class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                                >
                                    <Edit class="w-4 h-4" />
                                </button>
                                <button
                                    onclick={() => remove(announcement.id)}
                                    class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                                >
                                    <Trash2 class="w-4 h-4" />
                                </button>
                                <button
                                    onclick={() =>
                                        toggleExpand(announcement.id)}
                                    class="p-2 text-gray-400 hover:bg-gray-100 rounded-lg transition-all"
                                >
                                    {#if isExpanded}
                                        <ChevronUp class="w-5 h-5" />
                                    {:else}
                                        <ChevronDown class="w-5 h-5" />
                                    {/if}
                                </button>
                            </div>
                        </div>

                        {#if !isExpanded}
                            <p class="mt-2 text-gray-600 text-sm line-clamp-1">
                                {announcement.content}
                            </p>
                        {/if}
                    </div>

                    {#if isExpanded}
                        <div
                            class="px-5 pb-6 animate-in fade-in slide-in-from-top-2 duration-200"
                        >
                            <div
                                class="prose prose-sm max-w-none text-gray-700 bg-gray-50 p-4 rounded-xl border border-gray-100 whitespace-pre-wrap"
                            >
                                {announcement.content}
                            </div>
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    .line-clamp-1 {
        display: -webkit-box;
        -webkit-line-clamp: 1;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
</style>
