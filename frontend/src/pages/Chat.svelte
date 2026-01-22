<script lang="ts">
  // ⚠️⚠️⚠️ VERSION TEST: 2026-01-22 22:50 ⚠️⚠️⚠️
  console.log("🔥🔥🔥 CHAT.SVELTE VERSION 22:50 - CODE MỚI ĐÃ LOAD! 🔥🔥🔥");

  import { onMount, onDestroy } from "svelte";
  import { authStore } from "../stores/auth-store";
  import { apiFetch } from "../utils/api-fetch";
  import {
    getChannelsByUserId,
    getChannelById,
    type Channel,
  } from "../services/channel-service";
  import {
    getMessages,
    type Message as APIMessage,
  } from "../services/message-service";
  import {
    wsService,
    type IncomingMessagePayload,
    type ACKMessagePayload,
    type ErrorMessagePayload,
  } from "../services/websocket-service";

  interface Message {
    id: string;
    senderId: string;
    senderName: string;
    content: string;
    timestamp: string;
    read: boolean;
  }

  interface Conversation {
    id: string;
    channelId: string;
    userId: string;
    userName: string;
    userRole: string;
    userAvatar?: string;
    lastMessage: string;
    lastMessageTime: string;
    unread: number;
    messages: Message[];
  }

  let selectedConversation = $state<string | null>(null);
  let messageText = $state("");
  let searchTerm = $state("");
  let loading = $state(true);
  let loadingMessages = $state(false);
  let error = $state<string | null>(null);
  let wsConnected = $state(false);
  let pendingMessages = new Map<string, Message>();

  const currentUserId = $derived($authStore.user?.id || "");
  const currentUserName = $derived($authStore.user?.fullname || "User");

  let conversations = $state<Conversation[]>([]);

  // Load channels on mount
  onMount(async () => {
    console.log("🚀 Chat.svelte onMount started");

    if (!currentUserId) {
      error = "User not authenticated";
      loading = false;
      return;
    }

    try {
      loading = true;
      error = null;
      console.log("📞 Calling API with userId:", currentUserId);
      const channelsData = await getChannelsByUserId(currentUserId, 1, 50);
      console.log("📦 Raw API response:", channelsData);

      if (!channelsData || !channelsData.channels) {
        console.error("❌ Invalid response format:", channelsData);
        conversations = [];
        return;
      }

      // Convert channels to conversations
      conversations = channelsData.channels.map((channel) => {
        // FIXED: Only check member count, ignore adminIds
        // DM channel = exactly 2 members, Group channel = more than 2 members
        const isGroupChannel = channel.members.length > 2;

        let userName = "";
        let userId = "";
        let userAvatar = undefined;

        if (isGroupChannel) {
          userName = "Lớp học";
          userId = channel.id;
        } else {
          // For DM channels, get the other user
          const otherMember = channel.members.find(
            (m) => m.userId !== currentUserId,
          );
          userId = otherMember?.userId || "";
          userName =
            otherMember?.fullName ||
            otherMember?.full_name ||
            otherMember?.username ||
            "Unknown";
          userAvatar = otherMember?.avatar?.url;
        }

        return {
          id: channel.id,
          channelId: channel.id,
          userId,
          userName,
          userRole: isGroupChannel ? "Classroom" : "User",
          userAvatar,
          lastMessage: "Chưa có tin nhắn",
          lastMessageTime: "",
          unread: channel.unread_message_count || 0,
          messages: [],
        };
      });
    } catch (err) {
      error =
        err instanceof Error ? err.message : "Failed to load conversations";
      console.error("Error loading channels:", err);
    } finally {
      loading = false;

      // Auto-select channel from localStorage (set by "Chat với giảng viên")
      const storedChannelId = localStorage.getItem("openChannelId");
      if (storedChannelId) {
        // Don't remove immediately - keep it for refresh persistence
        console.log(
          "🔍 Auto-selecting channel from localStorage:",
          storedChannelId,
        );

        const existingConv = conversations.find(
          (c) => c.channelId === storedChannelId,
        );

        if (existingConv) {
          console.log("✅ Channel found in conversations");
          selectedConversation = storedChannelId;
          // Don't call loadMessages manually - let $effect handle it
        } else {
          console.log("📥 Fetching new channel...");
          try {
            const channel = await getChannelById(storedChannelId);
            console.log("✅ Channel fetched:", channel);

            const isGroupChannel = channel.members.length > 2;

            let userName = "";
            let userId = "";
            let userAvatar = undefined;

            if (isGroupChannel) {
              userName = "Lớp học";
              userId = channel.id;
            } else {
              const otherMember = channel.members.find(
                (m) => m.userId !== currentUserId,
              );
              userId = otherMember?.userId || "";
              userName =
                otherMember?.fullName ||
                otherMember?.full_name ||
                otherMember?.username ||
                "Unknown";
              userAvatar = otherMember?.avatar?.url;
            }

            const newConv: Conversation = {
              id: channel.id,
              channelId: channel.id,
              userId,
              userName,
              userRole: isGroupChannel ? "Classroom" : "User",
              userAvatar,
              lastMessage: "Chưa có tin nhắn",
              lastMessageTime: "",
              unread: 0,
              messages: [],
            };

            conversations = [newConv, ...conversations];
            selectedConversation = storedChannelId;
            console.log("✅ Auto-selecting channel:", storedChannelId);
            // Don't call loadMessages manually - let $effect handle it
          } catch (err) {
            console.error("❌ Failed to fetch channel:", err);
          }
        }
      }
    }

    // Connect WebSocket with retry
    const token = $authStore.accessToken;
    if (token) {
      // Add small delay to ensure page is fully loaded
      setTimeout(() => {
        wsService
          .connect(token)
          .then(() => {
            wsConnected = true;
            setupWebSocketHandlers();
            console.log("✅ WebSocket connected successfully");
          })
          .catch((err) => {
            console.error("❌ WebSocket connection failed:", err);
            // Retry after 2 seconds
            setTimeout(() => {
              console.log("🔄 Retrying WebSocket connection...");
              wsService
                .connect(token)
                .then(() => {
                  wsConnected = true;
                  setupWebSocketHandlers();
                  console.log("✅ WebSocket reconnected successfully");
                })
                .catch((retryErr) => {
                  console.error("❌ WebSocket retry failed:", retryErr);
                });
            }, 2000);
          });
      }, 500);
    }
  });

  // Cleanup WebSocket on destroy
  onDestroy(() => {
    wsService.disconnect();
  });

  // Load messages when conversation is selected
  async function loadMessages(channelId: string) {
    if (!channelId) return;

    try {
      loadingMessages = true;
      error = null;
      const messagesData = await getMessages({
        channel_id: channelId,
        page: 1,
        page_size: 100,
      });

      // Convert API messages to UI messages
      const uiMessages: Message[] = messagesData.messages.map((msg) => ({
        id: msg.id,
        senderId: msg.sender_id,
        senderName: msg.sender_username,
        content: msg.content,
        timestamp: new Date(msg.created_at).toLocaleString("vi-VN", {
          day: "2-digit",
          month: "2-digit",
          year: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        }),
        read: msg.read_by ? msg.read_by.includes(currentUserId) : false,
      }));

      // Update conversation with messages
      conversations = conversations.map((conv) => {
        if (conv.channelId === channelId) {
          return { ...conv, messages: uiMessages };
        }
        return conv;
      });
    } catch (err) {
      console.error("Error loading messages:", err);
      // Don't show error for 500 (likely empty channel), just show empty
      // Update conversation with empty messages
      conversations = conversations.map((conv) => {
        if (conv.channelId === channelId) {
          return { ...conv, messages: [] };
        }
        return conv;
      });
    } finally {
      loadingMessages = false;
    }
  }

  // Watch for conversation selection changes
  let loadingChannelId: string | null = null; // Track which channel is being loaded
  let attemptedChannels = new Set<string>(); // Track channels we already tried to load

  $effect(() => {
    if (selectedConversation) {
      const conv = conversations.find((c) => c.id === selectedConversation);
      console.log(
        "🔄 Effect triggered for conversation:",
        selectedConversation,
        "messages count:",
        conv?.messages?.length,
      );

      // Only load if: 1) conversation exists, 2) no messages, 3) not already loading, 4) haven't attempted yet
      if (
        conv &&
        conv.messages.length === 0 &&
        loadingChannelId !== conv.channelId &&
        !attemptedChannels.has(conv.channelId)
      ) {
        console.log("📥 Loading messages for channel:", conv.channelId);
        loadingChannelId = conv.channelId;
        attemptedChannels.add(conv.channelId);
        loadMessages(conv.channelId).finally(() => {
          loadingChannelId = null;
        });
      }
    }
  });

  // Mock conversations data
  // let conversations = $state<Conversation[]>([
  //   ... mock data removed ...
  // ]);

  const filteredConversations = $derived(
    conversations.filter(
      (conv) =>
        conv.userName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        conv.userRole.toLowerCase().includes(searchTerm.toLowerCase()),
    ),
  );

  const currentConversation = $derived(
    conversations.find((c) => c.id === selectedConversation),
  );

  async function handleDeleteMessage(messageId: string) {
    if (!selectedConversation) return;

    const conv = conversations.find((c) => c.id === selectedConversation);
    if (!conv) return;

    if (!confirm("Bạn có chắc muốn xóa tin nhắn này?")) return;

    try {
      await apiFetch(`/api/messages/${conv.channelId}/${messageId}`, {
        method: "DELETE",
      });

      // Remove message from UI
      conversations = conversations.map((c) => {
        if (c.id === selectedConversation) {
          return {
            ...c,
            messages: c.messages.filter((m) => m.id !== messageId),
          };
        }
        return c;
      });
    } catch (err) {
      console.error("Error deleting message:", err);
      error = "Không thể xóa tin nhắn. Vui lòng thử lại.";
    }
  }

  function handleSendMessage() {
    if (!messageText.trim() || !selectedConversation) return;

    const conv = conversations.find((c) => c.id === selectedConversation);
    if (!conv) return;

    const tempId = `temp_${Date.now()}`;
    const content = messageText;
    messageText = "";

    // Optimistic UI update
    const optimisticMessage: Message = {
      id: tempId,
      senderId: currentUserId,
      senderName: currentUserName,
      content,
      timestamp: new Date().toLocaleString("vi-VN", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      }),
      read: false,
    };

    conversations = conversations.map((c) => {
      if (c.id === selectedConversation) {
        return {
          ...c,
          messages: [...c.messages, optimisticMessage],
          lastMessage: content,
          lastMessageTime: "Vừa xong",
        };
      }
      return c;
    });

    pendingMessages.set(tempId, optimisticMessage);

    // Send via WebSocket
    if (wsConnected) {
      const sent = wsService.send("send_message", {
        channel_id: conv.channelId,
        content,
        temp_message_id: tempId,
      });

      if (!sent) {
        error = "Failed to send message. Please check connection.";
        // Remove optimistic message on failure
        conversations = conversations.map((c) => {
          if (c.id === selectedConversation) {
            return {
              ...c,
              messages: c.messages.filter((m) => m.id !== tempId),
            };
          }
          return c;
        });
        pendingMessages.delete(tempId);
      }
    } else {
      error = "WebSocket not connected";
      // Remove optimistic message
      conversations = conversations.map((c) => {
        if (c.id === selectedConversation) {
          return {
            ...c,
            messages: c.messages.filter((m) => m.id !== tempId),
          };
        }
        return c;
      });
      pendingMessages.delete(tempId);
    }
  }

  // Setup WebSocket event handlers
  function setupWebSocketHandlers() {
    // Handle ACK (message sent confirmation)
    wsService.on("ack_message", (msg) => {
      const payload = msg.payload as ACKMessagePayload;
      const tempId = payload.temp_message_id;
      const realMessage = payload.message;

      if (pendingMessages.has(tempId)) {
        // Replace temp message with real message
        conversations = conversations.map((conv) => {
          return {
            ...conv,
            messages: conv.messages.map((m) =>
              m.id === tempId
                ? {
                    id: realMessage.id,
                    senderId: realMessage.sender_id,
                    senderName: realMessage.sender_username,
                    content: realMessage.content,
                    timestamp: new Date(realMessage.created_at).toLocaleString(
                      "vi-VN",
                      {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      },
                    ),
                    read: realMessage.read_by.includes(currentUserId),
                  }
                : m,
            ),
          };
        });
        pendingMessages.delete(tempId);
      }
    });

    // Handle incoming messages from others
    wsService.on("send_message", (msg) => {
      const payload = msg.payload as IncomingMessagePayload;
      const incomingMessage = payload.message;

      // Don't add our own messages (already handled by ACK)
      if (incomingMessage.sender_id === currentUserId) return;

      const newMessage: Message = {
        id: incomingMessage.id,
        senderId: incomingMessage.sender_id,
        senderName: incomingMessage.sender_username,
        content: incomingMessage.content,
        timestamp: new Date(incomingMessage.created_at).toLocaleString(
          "vi-VN",
          {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
          },
        ),
        read: incomingMessage.read_by.includes(currentUserId),
      };

      // Add message to conversation
      conversations = conversations.map((conv) => {
        if (conv.channelId === incomingMessage.channel_id) {
          return {
            ...conv,
            messages: [...conv.messages, newMessage],
            lastMessage: incomingMessage.content,
            lastMessageTime: "Vừa xong",
            unread:
              conv.id !== selectedConversation ? (conv.unread || 0) + 1 : 0,
          };
        }
        return conv;
      });
    });

    // Handle errors
    wsService.on("error", (msg) => {
      const payload = msg.payload as ErrorMessagePayload;
      error = payload.error_msg;

      // Remove pending message if it failed
      if (
        payload.temp_message_id &&
        pendingMessages.has(payload.temp_message_id)
      ) {
        const tempId = payload.temp_message_id;
        conversations = conversations.map((conv) => {
          return {
            ...conv,
            messages: conv.messages.filter((m) => m.id !== tempId),
          };
        });
        pendingMessages.delete(tempId);
      }
    });
  }

  function handleKeyPress(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  }
</script>

<div class="chat-container">
  <!-- Conversations List -->
  <div class="conversations-panel">
    <!-- Header -->
    <div class="conversations-header">
      <h2 class="title">Tin nhắn</h2>
      <div class="search-box">
        <svg
          class="search-icon"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
        <input
          type="text"
          bind:value={searchTerm}
          placeholder="Tìm kiếm..."
          class="search-input"
        />
      </div>
    </div>

    <!-- Conversations List -->
    <div class="conversations-list">
      {#if loading}
        <div class="loading-state">
          <div class="spinner"></div>
          <p>Đang tải danh sách...</p>
        </div>
      {:else if error && conversations.length === 0}
        <div class="error-state">
          <p>❌ {error}</p>
        </div>
      {:else if filteredConversations.length === 0}
        <div class="empty-state">
          {conversations.length === 0
            ? "Chưa có cuộc trò chuyện nào"
            : "Không tìm thấy cuộc trò chuyện"}
        </div>
      {:else}
        {#each filteredConversations as conv (conv.id)}
          <button
            onclick={() => (selectedConversation = conv.id)}
            class="conversation-item"
            class:active={selectedConversation === conv.id}
          >
            <div class="avatar">
              {conv.userName.charAt(0)}
            </div>
            <div class="conversation-info">
              <div class="conversation-header-row">
                <span class="user-name">{conv.userName}</span>
                <span class="time">{conv.lastMessageTime}</span>
              </div>
              <div class="user-role">{conv.userRole}</div>
              <div class="last-message-row">
                <p class="last-message" class:unread={conv.unread > 0}>
                  {conv.lastMessage}
                </p>
                {#if conv.unread > 0}
                  <span class="unread-badge">{conv.unread}</span>
                {/if}
              </div>
            </div>
          </button>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Chat Area -->
  {#if selectedConversation && currentConversation}
    <div class="chat-panel">
      <!-- Chat Header -->
      <div class="chat-header">
        <div class="chat-header-user">
          <div class="avatar small">
            {currentConversation.userName.charAt(0)}
          </div>
          <div>
            <div class="chat-user-name">{currentConversation.userName}</div>
            <div class="chat-user-role">{currentConversation.userRole}</div>
          </div>
        </div>
        <button class="menu-button">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <circle cx="12" cy="5" r="2"></circle>
            <circle cx="12" cy="12" r="2"></circle>
            <circle cx="12" cy="19" r="2"></circle>
          </svg>
        </button>
      </div>

      <!-- Messages -->
      <div class="messages-container">
        {#if loadingMessages}
          <div class="loading-messages">
            <div class="spinner"></div>
            <p>Đang tải tin nhắn...</p>
          </div>
        {:else}
          {#each currentConversation.messages as message (message.id)}
            {@const isOwn = message.senderId === currentUserId}
            <div
              class="message-wrapper"
              class:own={isOwn}
              data-message-id={message.id}
              data-sender-id={message.senderId}
              data-current-user-id={currentUserId}
              data-is-own={isOwn}
            >
              {console.log(
                `🔍 Message ${message.id.slice(-4)}: senderId=${message.senderId}, currentUser=${currentUserId}, isOwn=${isOwn}`,
              )}
              <div class="message-bubble" class:own={isOwn}>
                {#if !isOwn}
                  <div class="sender-name">{message.senderName}</div>
                {/if}
                <div class="message-content">
                  <p>{message.content}</p>
                </div>
                <div class="message-footer" class:own={isOwn}>
                  <span class="timestamp">{message.timestamp}</span>
                  {#if isOwn}
                    <button
                      class="delete-btn"
                      onclick={() => handleDeleteMessage(message.id)}
                      title="Xóa tin nhắn"
                    >
                      <svg
                        width="14"
                        height="14"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path
                          d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                        ></path>
                        <line x1="10" y1="11" x2="10" y2="17"></line>
                        <line x1="14" y1="11" x2="14" y2="17"></line>
                      </svg>
                    </button>
                    {#if message.read}
                      <!-- CheckCheck icon -->
                      <svg
                        class="check-icon double"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path d="M18 6 7 17l-5-5"></path>
                        <path d="m22 6-11 11-2-2"></path>
                      </svg>
                    {:else}
                      <!-- Check icon -->
                      <svg
                        class="check-icon"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path d="M20 6 9 17l-5-5"></path>
                      </svg>
                    {/if}
                  {/if}
                </div>
              </div>
            </div>
          {/each}
        {/if}
      </div>

      <!-- Message Input -->
      <div class="message-input-container">
        <button class="action-button">
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48"
            ></path>
          </svg>
        </button>
        <div class="input-wrapper">
          <textarea
            bind:value={messageText}
            onkeypress={handleKeyPress}
            placeholder="Nhập tin nhắn..."
            class="message-textarea"
            rows="1"
          ></textarea>
        </div>
        <button class="action-button">
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
            <line x1="9" x2="9.01" y1="9" y2="9"></line>
            <line x1="15" x2="15.01" y1="9" y2="9"></line>
          </svg>
        </button>
        <button
          onclick={handleSendMessage}
          disabled={!messageText.trim()}
          class="send-button"
          class:disabled={!messageText.trim()}
        >
          <svg
            width="20"
            height="20"
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
    </div>
  {:else}
    <div class="empty-chat">
      <div class="empty-chat-content">
        <div class="empty-icon">
          <svg
            width="48"
            height="48"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.35-4.35"></path>
          </svg>
        </div>
        <p class="empty-title">Chọn một cuộc trò chuyện</p>
        <p class="empty-subtitle">
          Chọn một cuộc trò chuyện từ danh sách bên trái để bắt đầu nhắn tin
        </p>
      </div>
    </div>
  {/if}
</div>

<style>
  .chat-container {
    display: flex;
    height: 100%;
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  /* Conversations Panel */
  .conversations-panel {
    width: 320px;
    border-right: 1px solid #e5e7eb;
    display: flex;
    flex-direction: column;
  }

  .conversations-header {
    padding: 16px;
    border-bottom: 1px solid #e5e7eb;
  }

  .title {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 12px;
  }

  .search-box {
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: #9ca3af;
  }

  .search-input {
    width: 100%;
    padding: 8px 12px 8px 40px;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-size: 14px;
    outline: none;
    transition: all 0.2s;
  }

  .search-input:focus {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .conversations-list {
    flex: 1;
    overflow-y: auto;
  }

  .empty-state {
    padding: 32px;
    text-align: center;
    color: #6b7280;
  }

  .loading-state,
  .error-state {
    padding: 48px 24px;
    text-align: center;
    color: #6b7280;
  }

  .error-state {
    color: #dc2626;
  }

  .loading-messages {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    color: #6b7280;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid #e5e7eb;
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-bottom: 12px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .conversation-item {
    width: 100%;
    padding: 16px;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    border: none;
    background: none;
    border-bottom: 1px solid #f3f4f6;
    cursor: pointer;
    transition: background-color 0.2s;
    text-align: left;
  }

  .conversation-item:hover {
    background-color: #f9fafb;
  }

  .conversation-item.active {
    background-color: #eff6ff;
  }

  .avatar {
    width: 48px;
    height: 48px;
    background: #3b82f6;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: 600;
    flex-shrink: 0;
  }

  .avatar.small {
    width: 40px;
    height: 40px;
  }

  .conversation-info {
    flex: 1;
    min-width: 0;
  }

  .conversation-header-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 4px;
  }

  .user-name {
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .time {
    font-size: 12px;
    color: #6b7280;
    flex-shrink: 0;
  }

  .user-role {
    font-size: 14px;
    color: #6b7280;
    margin-bottom: 4px;
  }

  .last-message-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .last-message {
    font-size: 14px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
    margin: 0;
  }

  .last-message.unread {
    font-weight: 500;
    color: #111827;
  }

  .unread-badge {
    margin-left: 8px;
    background: #3b82f6;
    color: white;
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 12px;
    flex-shrink: 0;
  }

  /* Chat Panel */
  .chat-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .chat-header {
    padding: 16px;
    border-bottom: 1px solid #e5e7eb;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .chat-header-user {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .chat-user-name {
    font-weight: 500;
  }

  .chat-user-role {
    font-size: 14px;
    color: #6b7280;
  }

  .menu-button {
    padding: 8px;
    border: none;
    background: none;
    cursor: pointer;
    border-radius: 8px;
    color: #6b7280;
    transition: all 0.2s;
  }

  .menu-button:hover {
    background: #f3f4f6;
  }

  .messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .message-wrapper {
    display: flex;
    justify-content: flex-start;
  }

  .message-wrapper.own {
    justify-content: flex-end;
  }

  .message-bubble {
    max-width: 70%;
  }

  .message-bubble.own .message-content {
    background: #3b82f6;
    color: white;
  }

  .message-bubble.own .message-footer {
    color: rgba(255, 255, 255, 0.7);
  }

  .message-bubble.own .delete-btn {
    color: rgba(255, 255, 255, 0.7);
  }

  .message-bubble.own .delete-btn:hover {
    background: rgba(255, 255, 255, 0.2);
    color: white;
  }

  .sender-name {
    font-size: 12px;
    color: #6b7280;
    margin-bottom: 4px;
  }

  .message-content {
    padding: 12px 16px;
    border-radius: 12px;
    background: #f3f4f6;
    color: #111827;
  }

  .message-content p {
    margin: 0;
    word-wrap: break-word;
  }

  .message-footer {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: 4px;
    font-size: 11px;
    color: #6b7280;
  }

  .message-footer .timestamp {
    flex: 1;
  }

  .delete-btn {
    padding: 2px 4px;
    background: transparent;
    border: none;
    color: #6b7280;
    cursor: pointer;
    border-radius: 4px;
    opacity: 0.6;
    transition: all 0.2s;
  }

  .delete-btn:hover {
    opacity: 1;
    background: rgba(220, 38, 38, 0.1);
    color: #dc2626;
  }

  .message-footer.own {
    justify-content: flex-end;
  }

  .check-icon {
    color: #6b7280;
  }

  .check-icon.double {
    color: #3b82f6;
  }

  .message-input-container {
    padding: 16px;
    border-top: 1px solid #e5e7eb;
    display: flex;
    align-items: flex-end;
    gap: 8px;
  }

  .action-button {
    padding: 8px;
    border: none;
    background: none;
    cursor: pointer;
    border-radius: 8px;
    color: #6b7280;
    transition: all 0.2s;
    flex-shrink: 0;
  }

  .action-button:hover {
    background: #f3f4f6;
  }

  .input-wrapper {
    flex: 1;
    position: relative;
  }

  .message-textarea {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-family: inherit;
    font-size: 14px;
    resize: none;
    outline: none;
    min-height: 44px;
    max-height: 120px;
    transition: all 0.2s;
  }

  .message-textarea:focus {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .send-button {
    padding: 12px;
    border: none;
    background: #3b82f6;
    color: white;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s;
    flex-shrink: 0;
  }

  .send-button:hover:not(.disabled) {
    background: #2563eb;
  }

  .send-button.disabled {
    background: #f3f4f6;
    color: #9ca3af;
    cursor: not-allowed;
  }

  /* Empty Chat */
  .empty-chat {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #6b7280;
  }

  .empty-chat-content {
    text-align: center;
  }

  .empty-icon {
    width: 96px;
    height: 96px;
    background: #f3f4f6;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    color: #9ca3af;
  }

  .empty-title {
    font-size: 18px;
    margin-bottom: 8px;
  }

  .empty-subtitle {
    font-size: 14px;
  }
</style>
