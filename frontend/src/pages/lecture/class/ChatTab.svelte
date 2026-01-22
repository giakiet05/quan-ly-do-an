<script lang="ts">
    import { onMount, onDestroy, tick } from "svelte";
    import { Send, Paperclip, Smile, Hash, Loader2 } from "lucide-svelte";
    import { authStore } from "../../../stores/auth-store";
    import { wsService } from "../../../services/websocket-service";
    import { getMessages } from "../../../services/message-service";
    import { getChannelById } from "../../../services/channel-service";

    // --- Props ---
    let { generalChannelId } = $props<{ generalChannelId: string }>();

    // --- States (Svelte 5 Runes) ---
    let messageText = $state("");
    let loading = $state(true);
    let wsConnected = $state(false);
    let channelInfo = $state<any>(null);
    let messages = $state<any[]>([]);
    let scrollContainer = $state<HTMLDivElement | undefined>(undefined);

    // Lấy ID người dùng hiện tại từ store
    const currentUserId = $derived($authStore.user?.id || "");

    // --- Logic cuộn xuống cuối ---
    async function scrollToBottom() {
        await tick();
        if (scrollContainer) {
            scrollContainer.scrollTo({
                top: scrollContainer.scrollHeight,
                behavior: "smooth",
            });
        }
    }

    // --- Handlers cho WebSocket ---
    const handleIncomingMessage = (payload: any) => {
        const msg = payload.message;
        if (msg && String(msg.channelId) === String(generalChannelId)) {
            if (!messages.find((m) => String(m.id) === String(msg.id))) {
                messages = [...messages, msg];
                scrollToBottom();
            }
        }
    };

    const handleMessageAck = (payload: any) => {
        const msg = payload.message;
        if (msg && String(msg.channel_id) === String(generalChannelId)) {
            if (!messages.find((m) => String(m.id) === String(msg.id))) {
                messages = [...messages, msg];
                scrollToBottom();
            }
        }
    };

    // --- API & WebSocket ---
    onMount(async () => {
        if (!generalChannelId) return;

        try {
            loading = true;
            channelInfo = await getChannelById(generalChannelId);

            const res = await getMessages({
                channel_id: generalChannelId,
                page: 1,
                page_size: 50,
            });
            messages = res.messages.reverse();

            const token = localStorage.getItem("access_token");
            if (token) {
                await wsService.connect(token);
                wsConnected = true;
                wsService.on("send_message", handleIncomingMessage);
                wsService.on("ack_message", handleMessageAck);
            }

            await scrollToBottom();
        } catch (err) {
            console.error("Lỗi khởi tạo chat:", err);
        } finally {
            loading = false;
        }
    });

    onDestroy(() => {
        wsService.off("send_message", handleIncomingMessage);
        wsService.off("ack_message", handleMessageAck);
    });

    function handleSendMessage() {
        if (!messageText.trim() || !wsConnected) return;

        const success = wsService.send("send_message", {
            channel_id: generalChannelId,
            content: messageText,
            type: "text",
            temp_message_id: Date.now().toString(),
        });

        if (success) {
            messageText = "";
        }
    }

    function handleKeyPress(e: KeyboardEvent) {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    }
</script>

<div class="bg-white rounded-lg shadow-sm p-6 flex flex-col h-[700px]">
    {#if loading}
        <div
            class="flex-1 flex flex-col items-center justify-center bg-gray-50"
        >
            <Loader2 class="w-8 h-8 text-blue-600 animate-spin mb-2" />
            <p class="text-sm text-gray-500">Đang tải tin nhắn...</p>
        </div>
    {:else}
        <header
            class="h-16 px-6 border-b border-gray-100 flex items-center justify-between bg-white z-10"
        >
            <div class="flex items-center gap-3">
                <div
                    class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white shadow-sm"
                >
                    <Hash class="w-5 h-5" />
                </div>
                <div>
                    <h3 class="font-bold text-gray-900">
                        {channelInfo?.name || "Kênh thảo luận chung"}
                    </h3>
                    <div class="flex items-center gap-2">
                        <span class="flex h-2 w-2 rounded-full bg-green-500"
                        ></span>
                        <p class="text-[11px] text-gray-500 font-medium">
                            {channelInfo?.members?.length || 0} thành viên
                        </p>
                    </div>
                </div>
            </div>
        </header>

        <div
            bind:this={scrollContainer}
            class="flex-1 overflow-y-auto p-6 space-y-4 bg-[#f8f9fa]"
        >
            {#each messages as msg (msg.id || msg._id)}
                {@const msgSenderId = msg.senderId || msg.sender_id}
                {@const isOwn = String(msgSenderId) === String(currentUserId)}

                <div
                    class="flex {isOwn
                        ? 'justify-end'
                        : 'justify-start'} w-full"
                >
                    <div
                        class="flex flex-col {isOwn
                            ? 'items-end'
                            : 'items-start'} max-w-[80%]"
                    >
                        {#if !isOwn}
                            <span
                                class="text-[11px] font-bold text-gray-500 mb-1 ml-1 uppercase tracking-wider"
                            >
                                {msg.senderUsername ||
                                    msg.sender_username ||
                                    msg.senderName ||
                                    msg.sender_name ||
                                    "Người dùng"}
                            </span>
                        {/if}

                        <div
                            class="px-4 py-2.5 rounded-2xl text-sm shadow-sm
                            {isOwn
                                ? 'bg-blue-600 text-white rounded-tr-none'
                                : 'bg-white border border-gray-200 text-gray-800 rounded-tl-none'}"
                        >
                            {msg.content}
                        </div>

                        <span class="text-[9px] text-gray-400 mt-1 uppercase">
                            {new Date(
                                msg.createdAt || msg.created_at || Date.now(),
                            ).toLocaleTimeString("vi-VN", {
                                hour: "2-digit",
                                minute: "2-digit",
                            })}
                        </span>
                    </div>
                </div>
            {/each}
        </div>

        <footer class="p-4 bg-white border-t border-gray-100">
            <div
                class="flex items-end gap-2 bg-gray-50 rounded-2xl p-2 border border-gray-200 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 transition-all"
            >
                <button class="p-2 text-gray-400 hover:text-blue-600"
                    ><Paperclip class="w-5 h-5" /></button
                >
                <textarea
                    bind:value={messageText}
                    onkeypress={handleKeyPress}
                    placeholder="Nhập nội dung thảo luận..."
                    class="flex-1 bg-transparent border-none focus:ring-0 text-sm py-2 resize-none max-h-32"
                    rows="1"
                ></textarea>
                <button class="p-2 text-gray-400 hover:text-yellow-500"
                    ><Smile class="w-5 h-5" /></button
                >
                <button
                    onclick={handleSendMessage}
                    disabled={!messageText.trim()}
                    class="p-2.5 bg-blue-600 text-white rounded-xl disabled:bg-gray-300 transition-all active:scale-95 shadow-md shadow-blue-100"
                >
                    <Send class="w-4 h-4" />
                </button>
            </div>
        </footer>
    {/if}
</div>

<style>
    div::-webkit-scrollbar {
        width: 8px;
    }
    div::-webkit-scrollbar-track {
        background: #f1f1f1;
    }
    div::-webkit-scrollbar-thumb {
        background: #c1c1c1;
        border-radius: 4px;
    }
    div::-webkit-scrollbar-thumb:hover {
        background: #a8a8a8;
    }
</style>
