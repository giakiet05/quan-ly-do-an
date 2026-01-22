<script lang="ts">
    import { onMount } from "svelte";
    import {
        Plus,
        Calendar,
        Edit,
        Trash2,
        ChevronDown,
        ChevronUp,
        Pin,
        User,
    } from "lucide-svelte"; // Sửa lại import cho chuẩn Lucide
    import { postStore } from "../../../../stores/post-store";

    let { onOpen, onEdit, classroomId } = $props<{
        onOpen: () => void;
        onEdit: (announcement: any) => void;
        classroomId: string;
    }>();

    let expandedId = $state<string | null>(null);
    const { posts, fetchPosts } = postStore;

    onMount(() => {
        if (classroomId) fetchPosts(classroomId);
    });

    const formatDate = (dateString: string): string => {
        const date = new Date(dateString);
        return date.toLocaleDateString("vi-VN", {
            day: "2-digit",
            month: "short",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    function toggleExpand(id: string) {
        expandedId = expandedId === id ? null : id;
    }

    function togglePin(announcementId: string) {
        const post = $posts.find((p) => p.id === announcementId);
        if (post) {
            post.isPinned = !post.isPinned;
            $posts.sort((a, b) => Number(b.isPinned) - Number(a.isPinned));
        }
    }
</script>

<div class="max-w-4xl mx-auto space-y-6 p-4">
    <div
        class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8"
    >
        <div>
            <h2 class="text-3xl font-extrabold text-gray-900 tracking-tight">
                Bảng tin <span class="text-blue-600">Lớp học</span>
            </h2>
            <p class="text-gray-500 mt-1">
                Nơi cập nhật những thông tin quan trọng nhất
            </p>
        </div>

        <button
            onclick={onOpen}
            class="group flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-full font-semibold transition-all duration-300 shadow-[0_10px_20px_-10px_rgba(37,99,235,0.4)] hover:shadow-blue-500/40 active:scale-95"
        >
            <Plus class="w-5 h-5 transition-transform group-hover:rotate-90" />
            Tạo thông báo
        </button>
    </div>

    {#if $posts.length === 0}
        <div
            class="flex flex-col items-center justify-center py-24 bg-gray-50/50 rounded-3xl border-2 border-dashed border-gray-200"
        >
            <div class="bg-gray-100 p-4 rounded-full mb-4">
                <Calendar class="w-8 h-8 text-gray-400" />
            </div>
            <p class="text-gray-500 font-medium text-lg">Hộp thư đang trống</p>
        </div>
    {:else}
        <div class="grid gap-4">
            {#each $posts as announcement (announcement.id)}
                {@const isExpanded = expandedId === announcement.id}

                <div
                    class="relative group bg-white rounded-2xl border border-gray-100 transition-all duration-300
                    {announcement.isPinned
                        ? 'ring-2 ring-yellow-400/30'
                        : 'hover:shadow-xl hover:shadow-gray-200/50'}"
                >
                    {#if announcement.isPinned}
                        <div
                            class="absolute -top-3 left-6 px-3 py-1 bg-yellow-400 text-yellow-900 text-xs font-bold uppercase tracking-wider rounded-full shadow-sm flex items-center gap-1"
                        >
                            <Pin class="w-3 h-3" /> Ghim
                        </div>
                    {/if}

                    <div class="p-6">
                        <div
                            class="flex flex-col md:flex-row justify-between gap-4"
                        >
                            <div class="flex items-start gap-4">
                                <div class="relative">
                                    <img
                                        src={announcement.author.avatar.url}
                                        alt="Avatar"
                                        class="w-12 h-12 rounded-2xl object-cover ring-2 ring-gray-50 shadow-sm"
                                    />
                                    <div
                                        class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 border-2 border-white rounded-full"
                                    ></div>
                                </div>

                                <div class="flex-1 min-w-0">
                                    <button
                                        onclick={() =>
                                            toggleExpand(announcement.id)}
                                        class="text-left block group"
                                    >
                                        <h3
                                            class="text-xl font-bold text-gray-800 leading-snug group-hover:text-blue-600 transition-colors line-clamp-1"
                                        >
                                            {announcement.title}
                                        </h3>
                                    </button>

                                    <div
                                        class="flex flex-wrap items-center gap-x-4 gap-y-1 mt-2 text-sm text-gray-500"
                                    >
                                        <span
                                            class="flex items-center gap-1.5 font-medium text-gray-700"
                                        >
                                            <User
                                                class="w-4 h-4 text-blue-500"
                                            />
                                            {announcement.author.fullName}
                                        </span>
                                        <span class="flex items-center gap-1.5">
                                            <Calendar class="w-4 h-4" />
                                            {formatDate(announcement.createdAt)}
                                        </span>
                                    </div>
                                </div>
                            </div>

                            <div
                                class="flex items-center justify-end gap-2 bg-gray-50 p-1.5 rounded-xl self-end md:self-start"
                            >
                                <button
                                    onclick={() => togglePin(announcement.id)}
                                    title="Ghim thông báo"
                                    class="p-2.5 {announcement.isPinned
                                        ? 'text-yellow-600 bg-white shadow-sm'
                                        : 'text-gray-400 hover:text-yellow-500'} rounded-lg transition-all"
                                >
                                    <Pin
                                        class="w-4 h-4 {announcement.isPinned
                                            ? 'fill-current'
                                            : ''}"
                                    />
                                </button>

                                <div class="w-[1px] h-4 bg-gray-200 mx-1"></div>

                                <button
                                    onclick={() => onEdit(announcement)}
                                    class="p-2.5 text-gray-400 hover:text-blue-600 hover:bg-white hover:shadow-sm rounded-lg transition-all"
                                >
                                    <Edit class="w-4 h-4" />
                                </button>
                                <button
                                    onclick={() =>
                                        console.log("Delete", announcement.id)}
                                    class="p-2.5 text-gray-400 hover:text-red-600 hover:bg-white hover:shadow-sm rounded-lg transition-all"
                                >
                                    <Trash2 class="w-4 h-4" />
                                </button>

                                <button
                                    onclick={() =>
                                        toggleExpand(announcement.id)}
                                    class="ml-1 p-2.5 bg-gray-200/50 text-gray-600 hover:bg-gray-200 rounded-lg transition-all"
                                >
                                    <ChevronDown
                                        class="w-5 h-5 transition-transform duration-300 {isExpanded
                                            ? 'rotate-180'
                                            : ''}"
                                    />
                                </button>
                            </div>
                        </div>

                        {#if isExpanded}
                            <div class="mt-6 pt-6 border-t border-gray-100">
                                <div
                                    class="prose prose-blue max-w-none text-gray-700 leading-relaxed bg-blue-50/30 p-5 rounded-2xl border border-blue-100/50 whitespace-pre-wrap"
                                >
                                    {announcement.content}
                                </div>
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    /* Animation cho việc mở rộng content */
    div {
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
</style>
