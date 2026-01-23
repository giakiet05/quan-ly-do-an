<script lang="ts">
  import { onMount, onDestroy, tick } from "svelte";
  import { authStore } from "../stores/auth-store";
  import { wsService } from "../services/websocket-service";
  import { getMessages } from "../services/message-service";
  import { getChannelById } from "../services/channel-service";

  interface Team {
    id: string;
    name: string;
    leaderId: string;
    members: { id: string; name: string; studentId: string }[];
    groupChannelId?: string;
  }

  let { team, currentStudentId } = $props<{
    team: Team;
    currentStudentId: string;
  }>();

  let messageText = $state("");
  let loading = $state(true);
  let wsConnected = $state(false);
  let channelInfo = $state<any>(null);
  let messages = $state<any[]>([]);
  let scrollContainer = $state<HTMLDivElement | undefined>(undefined);

  const currentUserId = $derived($authStore.user?.id || "");

  async function scrollToBottom() {
    await tick();
    if (scrollContainer) {
      scrollContainer.scrollTo({
        top: scrollContainer.scrollHeight,
        behavior: "smooth",
      });
    }
  }

  const handleIncomingMessage = (payload: any) => {
    console.log("[Team Chat] Incoming message:", payload);
    const msg = payload.message;

    if (msg && team.groupChannelId) {
      const incomingId = String(msg.channelId || msg.channel_id);
      const currentId = String(team.groupChannelId);

      if (incomingId === currentId) {
        const exists = messages.some(
          (m) => String(m.id || m._id) === String(msg.id || msg._id),
        );

        if (!exists) {
          messages = [...messages, msg];
          console.log("[Team Chat] Added message, total:", messages.length);
          scrollToBottom();
        }
      }
    }
  };

  const handleMessageAck = (payload: any) => {
    const msg = payload.message;
    console.log("[Team Chat] Message ACK:", payload);
    if (
      msg &&
      team.groupChannelId &&
      String(msg.channel_id) === String(team.groupChannelId)
    ) {
      if (!messages.find((m) => String(m.id) === String(msg.id))) {
        messages = [...messages, msg];
        scrollToBottom();
      }
    }
  };

  onMount(async () => {
    if (!team.groupChannelId) {
      console.warn("[Team Chat] No groupChannelId provided");
      loading = false;
      return;
    }

    try {
      loading = true;
      console.log("[Team Chat] Loading channel:", team.groupChannelId);

      channelInfo = await getChannelById(team.groupChannelId);
      console.log("[Team Chat] Channel info:", channelInfo);

      const res = await getMessages({
        channel_id: team.groupChannelId,
        page: 1,
        page_size: 50,
      });
      messages = res.messages.reverse();
      console.log("[Team Chat] Loaded messages:", messages.length);

      const token = localStorage.getItem("access_token");
      if (token) {
        await wsService.connect(token);
        wsConnected = true;
        wsService.on("send_message", handleIncomingMessage);
        wsService.on("ack_message", handleMessageAck);
        console.log("[Team Chat] WebSocket connected");
      }

      await scrollToBottom();
    } catch (err) {
      console.error("[Team Chat] Error loading:", err);
    } finally {
      loading = false;
    }
  });

  onDestroy(() => {
    wsService.off("send_message", handleIncomingMessage);
    wsService.off("ack_message", handleMessageAck);
  });

  function handleSendMessage() {
    if (!messageText.trim() || !wsConnected || !team.groupChannelId) return;

    const success = wsService.send("send_message", {
      channel_id: team.groupChannelId,
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

  const memberMap = $derived(() => {
    const map = new Map<string, string>();
    if (!channelInfo) return map;

    channelInfo.members?.forEach((m: any) => {
      const id = String(m.userId || m.id || m._id);
      map.set(id, m.fullName || m.fullname || "Thành viên");
    });

    channelInfo.adminIds?.forEach((adminId: string) => {
      map.set(String(adminId), "Nhóm trưởng");
    });

    return map;
  });
</script>

<div class="bg-white rounded-lg shadow-sm p-6 flex flex-col h-[700px]">
  {#if loading}
    <div class="flex-1 flex flex-col items-center justify-center bg-gray-50">
      <div
        class="w-10 h-10 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin"
      ></div>
      <p class="text-sm text-gray-500 mt-4">Đang tải tin nhắn...</p>
    </div>
  {:else if !team.groupChannelId}
    <div class="flex-1 flex flex-col items-center justify-center bg-gray-50">
      <p class="text-gray-500">Channel nhóm chưa được tạo</p>
    </div>
  {:else}
    <header
      class="h-16 px-6 border-b border-gray-100 flex items-center justify-between bg-white z-10"
    >
      <div class="flex items-center gap-3">
        <div
          class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white shadow-sm"
        >
          💬
        </div>
        <div>
          <h3 class="font-bold text-gray-900">Chat nhóm</h3>
          <div class="flex items-center gap-2">
            <span
              class="flex h-2 w-2 rounded-full {wsConnected
                ? 'bg-green-500'
                : 'bg-gray-400'}"
            ></span>
            <p class="text-[11px] text-gray-500 font-medium">
              {wsConnected ? "Đang kết nối" : "Ngoại tuyến"}
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

        <div class="flex {isOwn ? 'justify-end' : 'justify-start'} w-full">
          <div
            class="flex flex-col {isOwn
              ? 'items-end'
              : 'items-start'} max-w-[80%]"
          >
            {#if !isOwn}
              <span
                class="text-[11px] font-bold text-gray-500 mb-1 ml-1 uppercase tracking-wider"
              >
                {memberMap().get(String(msgSenderId)) || "Người dùng"}
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
        <textarea
          bind:value={messageText}
          onkeypress={handleKeyPress}
          placeholder="Nhập tin nhắn..."
          class="flex-1 bg-transparent border-none focus:ring-0 text-sm py-2 resize-none max-h-32"
          rows="1"
          disabled={!wsConnected}
        ></textarea>
        <button
          onclick={handleSendMessage}
          disabled={!messageText.trim() || !wsConnected}
          class="p-2.5 bg-blue-600 text-white rounded-xl disabled:bg-gray-300 transition-all active:scale-95 shadow-md shadow-blue-100"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="m22 2-7 20-4-9-9-4Z"></path>
            <path d="M22 2 11 13"></path>
          </svg>
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
