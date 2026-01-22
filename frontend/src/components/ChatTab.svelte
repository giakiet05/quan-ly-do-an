<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { getMessages } from "../services/message-service";
  import { getChannelById } from "../services/channel-service";
  import { wsService } from "../services/websocket-service";
  import { authStore } from "../stores/auth-store";
  import type { Message as APIMessage } from "../models";

  let { channelId, currentUserRole, currentUserName } = $props<{
    channelId: string;
    currentUserRole: "teacher" | "student";
    currentUserName: string;
  }>();

  interface Message {
    id: string;
    userId: string;
    userName: string;
    userRole: "teacher" | "student";
    content: string;
    timestamp: string;
    attachments?: string[];
  }

  let messages = $state<Message[]>([]);
  let loading = $state(true);
  let wsConnected = $state(false);
  let pendingMessages = $state<Map<string, Message>>(new Map());

  let newMessage = $state("");
  let messagesEndRef: HTMLDivElement;
  let textareaRef: HTMLTextAreaElement;

  function scrollToBottom() {
    messagesEndRef?.scrollIntoView({ behavior: "smooth" });
  }

  async function loadMessages() {
    try {
      loading = true;
      console.log("📥 Loading messages for channel:", channelId);

      // Check if channelId is valid
      if (!channelId || channelId === "") {
        console.warn("⚠️ No channel ID provided for ChatTab");
        messages = [];
        loading = false;
        return;
      }

      // TEMPORARY: Use mock data if channelId is mock
      if (channelId.startsWith("channel") && !channelId.includes("-")) {
        // Mock channel ID detected
        messages = [
          {
            id: "1",
            userId: "teacher1",
            userName: "TS. Nguyễn Văn A",
            userRole: "teacher",
            content:
              "Chào các em! Đây là kênh chat chung của lớp. Các em có thể trao đổi, thảo luận về đề tài ở đây.",
            timestamp: new Date().toLocaleString("vi-VN"),
            attachments: [],
          },
          {
            id: "2",
            userId: "student1",
            userName: "Nguyễn Văn Minh",
            userRole: "student",
            content: "Dạ em chào thầy ạ!",
            timestamp: new Date().toLocaleString("vi-VN"),
            attachments: [],
          },
        ];
        loading = false;
        scrollToBottom();
        return;
      }

      // Fetch messages from API
      const response = await getMessages({ channel_id: channelId });
      messages = response.messages.map((msg) => ({
        id: msg.id,
        userId: msg.sender_id,
        userName: msg.sender_username,
        userRole: currentUserRole,
        content: msg.content,
        timestamp: new Date(msg.created_at).toLocaleString("vi-VN"),
        attachments: [],
      }));
      console.log("✅ Messages loaded:", messages.length);
    } catch (err) {
      console.error("❌ Error loading messages:", err);
    } finally {
      loading = false;
      scrollToBottom();
    }
  }

  onMount(async () => {
    try {
      // Connect to WebSocket
      const token = localStorage.getItem("accessToken");
      if (token) {
        await wsService.connect(token);
        wsConnected = wsService.isConnected();

        // Listen for incoming messages
        wsService.on("send_message", handleIncomingMessage);
        wsService.on("ack_message", handleMessageAck);
      }
    } catch (err) {
      console.error("Error connecting WebSocket:", err);
    }
  });

  // Reload messages when channelId changes
  $effect(() => {
    if (channelId) {
      console.log("🔄 Channel changed, reloading messages:", channelId);
      loadMessages();
    }
  });

  onDestroy(() => {
    wsService.off("send_message", handleIncomingMessage);
    wsService.off("ack_message", handleMessageAck);
  });

  function handleIncomingMessage(data: any) {
    // Only add message if it belongs to current channel
    if (data.channel_id !== channelId) {
      return;
    }
    const newMsg: Message = {
      id: data.message_id,
      userId: data.sender_id,
      userName: data.sender_name || "Unknown",
      userRole:
        data.sender_id === $authStore.user?.id ? currentUserRole : "student",
      content: data.content,
      timestamp: new Date().toLocaleString("vi-VN"),
      attachments: [],
    };
    messages = [...messages, newMsg];
    scrollToBottom();
  }

  function handleMessageAck(data: any) {
    const tempMsg = pendingMessages.get(data.temp_id);
    if (tempMsg) {
      const index = messages.findIndex((m) => m.id === data.temp_id);
      if (index !== -1) {
        messages[index] = {
          ...tempMsg,
          id: data.message_id,
        };
      }
      pendingMessages.delete(data.temp_id);
    }
  }

  $effect(() => {
    if (messages.length) {
      scrollToBottom();
    }
  });

  function handleSendMessage() {
    if (!newMessage.trim()) return;

    const tempId = `temp_${Date.now()}`;
    const message: Message = {
      id: tempId,
      userId: $authStore.user?.id || "unknown",
      userName: currentUserName,
      userRole: currentUserRole,
      content: newMessage,
      timestamp: new Date().toLocaleString("vi-VN", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
      }),
      attachments: [],
    };

    // Add to messages immediately (optimistic UI)
    messages = [...messages, message];

    // If mock channelId, just keep local message
    if (!channelId.startsWith("channel") || channelId.includes("-")) {
      // Real channel - send via WebSocket
      if (wsConnected) {
        pendingMessages.set(tempId, message);
        wsService.send("new_message", {
          channel_id: channelId,
          content: newMessage,
          temp_message_id: tempId,
        });
      }
    }

    newMessage = "";

    // Reset textarea height
    if (textareaRef) {
      textareaRef.style.height = "auto";
    }
  }

  function handleKeyPress(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  }

  function handleTextareaChange() {
    // Auto-resize textarea
    if (textareaRef) {
      textareaRef.style.height = "auto";
      textareaRef.style.height = textareaRef.scrollHeight + "px";
    }
  }

  function groupMessagesByDate(messages: Message[]) {
    const groups: { [key: string]: Message[] } = {};

    messages.forEach((message) => {
      const date = message.timestamp.split(" ")[0];
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(message);
    });

    return groups;
  }

  let messageGroups = $derived(groupMessagesByDate(messages));
</script>

<div class="chat-container">
  <!-- Chat Header -->
  <div class="chat-header">
    <div class="header-content">
      <div>
        <h3 class="header-title">Chat lớp học</h3>
        <p class="header-subtitle">{messages.length} tin nhắn</p>
      </div>
      <button class="btn-more" aria-label="Tùy chọn khác">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle cx="12" cy="12" r="1"></circle>
          <circle cx="12" cy="5" r="1"></circle>
          <circle cx="12" cy="19" r="1"></circle>
        </svg>
      </button>
    </div>
  </div>

  <!-- Messages Area -->
  <div class="messages-area">
    {#each Object.entries(messageGroups) as [date, groupMessages]}
      <div>
        <!-- Date Divider -->
        <div class="date-divider">
          <div class="date-badge">{date}</div>
        </div>

        <!-- Messages -->
        {#each groupMessages as message, index}
          {@const isCurrentUser =
            (currentUserRole === "teacher" && message.userRole === "teacher") ||
            (currentUserRole === "student" &&
              message.userName === currentUserName &&
              message.userRole === "student")}
          {@const showAvatar =
            index === 0 || groupMessages[index - 1].userId !== message.userId}
          {@const isLastInGroup =
            index === groupMessages.length - 1 ||
            groupMessages[index + 1].userId !== message.userId}

          <div
            class="message-row"
            class:reverse={isCurrentUser}
            class:show-avatar={showAvatar}
          >
            <!-- Avatar -->
            <div class="avatar-wrapper" class:invisible={!showAvatar}>
              <div
                class="avatar"
                class:teacher={message.userRole === "teacher"}
                class:student={message.userRole === "student"}
              >
                {message.userName.charAt(0)}
              </div>
            </div>

            <!-- Message Content -->
            <div class="message-content" class:align-end={isCurrentUser}>
              {#if showAvatar}
                <div class="message-header" class:reverse={isCurrentUser}>
                  <span class="message-author">{message.userName}</span>
                  {#if message.userRole === "teacher"}
                    <span class="badge-teacher">Giáo viên</span>
                  {/if}
                </div>
              {/if}

              <div
                class="message-bubble"
                class:current-user={isCurrentUser}
                class:teacher={isCurrentUser && message.userRole === "teacher"}
                class:student={isCurrentUser && message.userRole === "student"}
                class:other-user={!isCurrentUser}
              >
                <p class="message-text">{message.content}</p>
              </div>

              {#if isLastInGroup}
                <span class="message-time">
                  {message.timestamp.split(" ")[1]}
                </span>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/each}
    <div bind:this={messagesEndRef}></div>
  </div>

  <!-- Input Area -->
  <div class="input-area">
    <div class="input-wrapper">
      <button class="btn-attach" aria-label="Đính kèm file">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"
          />
        </svg>
      </button>

      <div class="textarea-wrapper">
        <textarea
          bind:this={textareaRef}
          bind:value={newMessage}
          oninput={handleTextareaChange}
          onkeypress={handleKeyPress}
          placeholder="Nhập tin nhắn..."
          rows="1"
        ></textarea>
      </div>

      <button class="btn-emoji" aria-label="Chọn biểu tượng cảm xúc">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
          <line x1="9" y1="9" x2="9.01" y2="9"></line>
          <line x1="15" y1="9" x2="15.01" y2="9"></line>
        </svg>
      </button>

      <button
        class="btn-send"
        disabled={!newMessage.trim()}
        onclick={handleSendMessage}
        aria-label="Gửi tin nhắn"
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <line x1="22" y1="2" x2="11" y2="13"></line>
          <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
        </svg>
      </button>
    </div>
    <p class="input-hint">Nhấn Enter để gửi, Shift + Enter để xuống dòng</p>
  </div>
</div>

<style>
  .chat-container {
    display: flex;
    flex-direction: column;
    height: 600px;
    background: white;
    border-radius: 0.75rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  /* Chat Header */
  .chat-header {
    border-bottom: 1px solid #e5e7eb;
    padding: 1rem;
  }

  .header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
  }

  .header-subtitle {
    font-size: 0.875rem;
    color: #6b7280;
  }

  .btn-more {
    padding: 0.5rem;
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    border-radius: 0.5rem;
    transition: background 0.2s;
  }

  .btn-more:hover {
    background: #f3f4f6;
  }

  /* Messages Area */
  .messages-area {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
  }

  .date-divider {
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 1rem 0;
  }

  .date-badge {
    padding: 0.25rem 0.75rem;
    background: #f3f4f6;
    border-radius: 9999px;
    font-size: 0.75rem;
    color: #6b7280;
  }

  .message-row {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 0.25rem;
  }

  .message-row.reverse {
    flex-direction: row-reverse;
  }

  .message-row.show-avatar {
    margin-top: 1rem;
  }

  .avatar-wrapper {
    flex-shrink: 0;
  }

  .avatar-wrapper.invisible {
    visibility: hidden;
  }

  .avatar {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    color: white;
  }

  .avatar.teacher {
    background: #3b82f6;
  }

  .avatar.student {
    background: #d1d5db;
    color: #4b5563;
  }

  .message-content {
    display: flex;
    flex-direction: column;
    max-width: 70%;
  }

  .message-content.align-end {
    align-items: flex-end;
  }

  .message-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.25rem;
  }

  .message-header.reverse {
    flex-direction: row-reverse;
  }

  .message-author {
    font-size: 0.875rem;
    font-weight: 500;
    color: #111827;
  }

  .badge-teacher {
    padding: 0.125rem 0.5rem;
    background: #dbeafe;
    color: #1e40af;
    font-size: 0.75rem;
    border-radius: 0.25rem;
  }

  .message-bubble {
    padding: 0.75rem 1rem;
    border-radius: 1rem;
    word-wrap: break-word;
  }

  .message-bubble.current-user {
    border-top-right-radius: 0.25rem;
  }

  .message-bubble.teacher {
    background: #3b82f6;
    color: white;
  }

  .message-bubble.student {
    background: #3b82f6;
    color: white;
  }

  .message-bubble.other-user {
    background: #f3f4f6;
    color: #111827;
    border-top-left-radius: 0.25rem;
  }

  .message-text {
    font-size: 0.875rem;
    white-space: pre-wrap;
  }

  .message-time {
    font-size: 0.75rem;
    color: #9ca3af;
    margin-top: 0.25rem;
    padding: 0 0.5rem;
  }

  /* Input Area */
  .input-area {
    border-top: 1px solid #e5e7eb;
    padding: 1rem;
  }

  .input-wrapper {
    display: flex;
    align-items: flex-end;
    gap: 0.75rem;
  }

  .btn-attach,
  .btn-emoji {
    flex-shrink: 0;
    padding: 0.5rem;
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    border-radius: 0.5rem;
    transition: background 0.2s;
  }

  .btn-attach:hover,
  .btn-emoji:hover {
    background: #f3f4f6;
  }

  .textarea-wrapper {
    flex: 1;
    position: relative;
  }

  textarea {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    resize: none;
    max-height: 8rem;
    font-family: inherit;
    font-size: 0.875rem;
    line-height: 1.5;
    transition: border-color 0.2s;
  }

  textarea:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  textarea::placeholder {
    color: #9ca3af;
  }

  .btn-send {
    flex-shrink: 0;
    padding: 0.5rem;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 0.5rem;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-send:hover:not(:disabled) {
    background: #2563eb;
  }

  .btn-send:disabled {
    background: #e5e7eb;
    color: #9ca3af;
    cursor: not-allowed;
  }

  .input-hint {
    margin-top: 0.5rem;
    font-size: 0.75rem;
    color: #9ca3af;
  }
</style>
