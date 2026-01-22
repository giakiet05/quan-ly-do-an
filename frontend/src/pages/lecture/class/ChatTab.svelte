<script lang="ts">
    import { onMount, onDestroy, tick } from "svelte";
    import {
        Send,
        Search,
        Paperclip,
        Smile,
        MoreVertical,
        MessageSquare,
        User,
        Hash,
        Check,
        CheckCheck,
    } from "lucide-svelte";
    import { authStore } from "../../../stores/auth-store";
    import { wsService } from "../../../services/websocket-service";
    import { getMessages } from "../../../services/message-service";

    // --- Props ---
    let { classroomId } = $props<{ classroomId: string }>();

    // --- Types nội bộ ---
    interface Message {
        id: string;
        senderId: string;
        senderName: string;
        content: string;
        timestamp: string;
        read: boolean;
    }

    interface Channel {
        id: string;
        name: string;
        type: "group" | "direct";
        lastMessage?: string;
        unreadCount: number;
    }

    // --- States (Svelte 5 Runes) ---
    let messageText = $state("");
    let searchTerm = $state("");
    let loading = $state(true);
    let wsConnected = $state(false);
    let selectedChannelId = $state<string | null>(null);

    // Dữ liệu danh sách kênh và tin nhắn
    let channels = $state<Channel[]>([]);
    let messages = $state<Message[]>([]);

    let scrollContainer: HTMLDivElement;

    // --- Derived ---
    const currentUser = $derived($authStore.user);
    const currentUserId = $derived(currentUser?.id || "");

    const filteredChannels = $derived(
        channels.filter((c) =>
            c.name.toLowerCase().includes(searchTerm.toLowerCase()),
        ),
    );

    const activeChannel = $derived(
        channels.find((c) => c.id === selectedChannelId),
    );

    // --- Logic cuộn xuống cuối ---
    async function scrollToBottom() {
        await tick();
        if (scrollContainer) {
            scrollContainer.scrollTop = scrollContainer.scrollHeight;
        }
    }

    // --- API & WebSocket ---
    onMount(async () => {
        try {
            loading = true;
            // GIẢ VỜ: Gọi API lấy danh sách channel của lớp
            // const data = await getChannelsByClassId(classId);

            // Mock dữ liệu channel
            channels = [
                {
                    id: "class_general",
                    name: "Kênh chung của lớp",
                    type: "group",
                    unreadCount: 2,
                    lastMessage: "Thông báo lịch thi...",
                },
                {
                    id: "group_1",
                    name: "Nhóm đồ án 1",
                    type: "group",
                    unreadCount: 0,
                    lastMessage: "Đã nộp bài chưa?",
                },
            ];

            // Kết nối WebSocket
            const token = $authStore.accessToken;
            if (token) {
                await wsService.connect(token);
                wsConnected = true;
                setupWebSocket();
            }
        } catch (err) {
            console.error("Lỗi khởi tạo chat:", err);
        } finally {
            loading = false;
        }
    });

    onDestroy(() => {
        wsService.disconnect();
    });

    function setupWebSocket() {
        wsService.on("send_message", (payload: any) => {
            const msg = payload.message;
            if (msg.channel_id === selectedChannelId) {
                const newMessage: Message = {
                    id: msg.id,
                    senderId: msg.sender_id,
                    senderName: msg.sender_username,
                    content: msg.content,
                    timestamp: new Date(msg.created_at).toLocaleTimeString(
                        "vi-VN",
                        { hour: "2-digit", minute: "2-digit" },
                    ),
                    read: false,
                };
                messages = [...messages, newMessage];
                scrollToBottom();
            }
        });
    }

    async function selectChannel(id: string) {
        selectedChannelId = id;
        // GIẢ VỜ: Load tin nhắn cũ từ API
        // const data = await getMessages({ channel_id: id, page: 1 });
        messages = []; // Clear cũ
        await scrollToBottom();
    }

    function handleSendMessage() {
        if (!messageText.trim() || !selectedChannelId) return;

        const tempMsg: Message = {
            id: Date.now().toString(),
            senderId: currentUserId,
            senderName: currentUser?.fullname || "Tôi",
            content: messageText,
            timestamp: new Date().toLocaleTimeString("vi-VN", {
                hour: "2-digit",
                minute: "2-digit",
            }),
            read: false,
        };

        // Optimistic Update
        messages = [...messages, tempMsg];
        wsService.send("new_message", {
            channel_id: selectedChannelId,
            content: messageText,
        });

        messageText = "";
        scrollToBottom();
    }

    function handleKeyPress(e: KeyboardEvent) {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    }
</script>

<div
    class="bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden flex h-[650px]"
>
    <aside class="w-80 border-r border-gray-100 flex flex-col bg-gray-50/30">
        <div class="p-4 bg-white border-b border-gray-100">
            <h2
                class="text-lg font-bold text-gray-800 flex items-center gap-2 mb-3"
            >
                <MessageSquare class="w-5 h-5 text-blue-600" />
                Trò chuyện
            </h2>
            <div class="relative">
                <Search
                    class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
                />
                <input
                    type="text"
                    bind:value={searchTerm}
                    placeholder="Tìm phòng chat..."
                    class="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                />
            </div>
        </div>

        <div class="flex-1 overflow-y-auto">
            {#if loading}
                <div
                    class="p-8 text-center text-gray-400 animate-pulse text-sm"
                >
                    Đang tải phòng chat...
                </div>
            {:else if filteredChannels.length === 0}
                <div class="p-8 text-center text-gray-400 text-sm">
                    Không tìm thấy phòng chat
                </div>
            {:else}
                {#each filteredChannels as channel (channel.id)}
                    <button
                        onclick={() => selectChannel(channel.id)}
                        class="w-full p-4 flex items-start gap-3 transition-all border-b border-gray-50 hover:bg-white group {selectedChannelId ===
                        channel.id
                            ? 'bg-white border-l-4 border-l-blue-600'
                            : ''}"
                    >
                        <div
                            class="w-10 h-10 rounded-full {selectedChannelId ===
                            channel.id
                                ? 'bg-blue-600'
                                : 'bg-gray-200 group-hover:bg-blue-100'} flex items-center justify-center transition-colors"
                        >
                            {#if channel.type === "group"}
                                <Hash
                                    class="w-5 h-5 {selectedChannelId ===
                                    channel.id
                                        ? 'text-white'
                                        : 'text-gray-500'}"
                                />
                            {:else}
                                <User
                                    class="w-5 h-5 {selectedChannelId ===
                                    channel.id
                                        ? 'text-white'
                                        : 'text-gray-500'}"
                                />
                            {/if}
                        </div>
                        <div class="flex-1 text-left overflow-hidden">
                            <div class="flex justify-between items-center mb-1">
                                <span
                                    class="font-semibold text-sm text-gray-900 truncate"
                                    >{channel.name}</span
                                >
                                {#if channel.unreadCount > 0}
                                    <span
                                        class="bg-blue-600 text-white text-[10px] px-1.5 py-0.5 rounded-full"
                                        >{channel.unreadCount}</span
                                    >
                                {/if}
                            </div>
                            <p class="text-xs text-gray-500 truncate">
                                {channel.lastMessage || "Chưa có tin nhắn"}
                            </p>
                        </div>
                    </button>
                {/each}
            {/if}
        </div>
    </aside>

    <main class="flex-1 flex flex-col bg-white">
        {#if activeChannel}
            <header
                class="h-16 px-6 border-b border-gray-100 flex items-center justify-between bg-white"
            >
                <div class="flex items-center gap-3">
                    <div class="p-2 bg-blue-50 rounded-lg">
                        <Hash class="w-5 h-5 text-blue-600" />
                    </div>
                    <div>
                        <h3 class="font-bold text-gray-900">
                            {activeChannel.name}
                        </h3>
                        <p
                            class="text-[11px] text-green-500 flex items-center gap-1"
                        >
                            <span class="w-1.5 h-1.5 rounded-full bg-green-500"
                            ></span> Trực tuyến
                        </p>
                    </div>
                </div>
                <button
                    class="p-2 hover:bg-gray-100 rounded-full text-gray-400 transition-colors"
                >
                    <MoreVertical class="w-5 h-5" />
                </button>
            </header>

            <div
                bind:this={scrollContainer}
                class="flex-1 overflow-y-auto p-6 space-y-6 bg-gray-50/50"
            >
                {#each messages as msg (msg.id)}
                    {@const isOwn = msg.senderId === currentUserId}
                    <div class="flex {isOwn ? 'justify-end' : 'justify-start'}">
                        <div
                            class="flex flex-col {isOwn
                                ? 'items-end'
                                : 'items-start'} max-w-[75%]"
                        >
                            {#if !isOwn}
                                <span
                                    class="text-xs font-medium text-gray-500 ml-1 mb-1"
                                    >{msg.senderName}</span
                                >
                            {/if}
                            <div class="relative group">
                                <div
                                    class="px-4 py-2.5 rounded-2xl text-sm shadow-sm
                  {isOwn
                                        ? 'bg-blue-600 text-white rounded-tr-none'
                                        : 'bg-white border border-gray-200 text-gray-800 rounded-tl-none'}"
                                >
                                    {msg.content}
                                </div>
                            </div>
                            <div class="flex items-center gap-1 mt-1 px-1">
                                <span class="text-[10px] text-gray-400"
                                    >{msg.timestamp}</span
                                >
                                {#if isOwn}
                                    {#if msg.read}
                                        <CheckCheck
                                            class="w-3 h-3 text-blue-500"
                                        />
                                    {:else}
                                        <Check class="w-3 h-3 text-gray-300" />
                                    {/if}
                                {/if}
                            </div>
                        </div>
                    </div>
                {/each}
            </div>

            <footer class="p-4 bg-white border-t border-gray-100">
                <div class="max-w-4xl mx-auto flex items-end gap-2">
                    <div
                        class="flex-1 bg-gray-50 border border-gray-200 rounded-2xl px-4 py-2 focus-within:ring-2 focus-within:ring-blue-500/20 focus-within:border-blue-500 transition-all flex items-end gap-2"
                    >
                        <button
                            class="p-1.5 text-gray-400 hover:text-blue-600 transition-colors"
                            ><Paperclip class="w-5 h-5" /></button
                        >
                        <textarea
                            bind:value={messageText}
                            onkeypress={handleKeyPress}
                            placeholder="Viết tin nhắn cho lớp..."
                            class="flex-1 bg-transparent border-none focus:ring-0 text-sm py-2 resize-none max-h-32 min-h-[20px]"
                            rows="1"
                        ></textarea>
                        <button
                            class="p-1.5 text-gray-400 hover:text-yellow-500 transition-colors"
                            ><Smile class="w-5 h-5" /></button
                        >
                    </div>
                    <button
                        onclick={handleSendMessage}
                        disabled={!messageText.trim()}
                        class="p-3 bg-blue-600 text-white rounded-xl shadow-blue-200 shadow-lg hover:bg-blue-700 hover:-translate-y-0.5 disabled:bg-gray-200 disabled:shadow-none transition-all active:scale-95"
                    >
                        <Send class="w-5 h-5" />
                    </button>
                </div>
            </footer>
        {:else}
            <div
                class="flex-1 flex flex-col items-center justify-center p-12 text-center bg-gray-50/30"
            >
                <div
                    class="w-20 h-20 bg-blue-50 rounded-full flex items-center justify-center mb-4"
                >
                    <MessageSquare class="w-10 h-10 text-blue-200" />
                </div>
                <h3 class="text-lg font-bold text-gray-800">
                    Chọn một cuộc trò chuyện
                </h3>
                <p class="text-sm text-gray-500 max-w-xs mt-2">
                    Kết nối với giảng viên và sinh viên khác trong lớp học để
                    trao đổi kiến thức.
                </p>
            </div>
        {/if}
    </main>
</div>

<style>
    /* Tùy chỉnh thanh cuộn cho chuyên nghiệp */
    ::-webkit-scrollbar {
        width: 6px;
    }
    ::-webkit-scrollbar-track {
        background: transparent;
    }
    ::-webkit-scrollbar-thumb {
        background: #e5e7eb;
        border-radius: 10px;
    }
    ::-webkit-scrollbar-thumb:hover {
        background: #d1d5db;
    }
</style>
