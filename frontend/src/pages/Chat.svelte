<script lang="ts">
  import { authStore } from "../stores/auth-store";

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

  const currentUserId = $derived($authStore.user?.id || "current-user");
  const currentUserName = $derived($authStore.user?.fullname || "User");

  // Mock conversations data
  let conversations = $state<Conversation[]>([
    {
      id: "c1",
      userId: "sv1",
      userName: "Nguyễn Văn An",
      userRole: "MSSV: 20200001",
      lastMessage: "Em cảm ơn thầy ạ!",
      lastMessageTime: "10:30",
      unread: 2,
      messages: [
        {
          id: "m1",
          senderId: "sv1",
          senderName: "Nguyễn Văn An",
          content: "Thầy ơi, em có thắc mắc về đề tài ạ",
          timestamp: "10:25",
          read: true,
        },
        {
          id: "m2",
          senderId: currentUserId,
          senderName: currentUserName,
          content: "Chào em, em cứ hỏi thầy nhé",
          timestamp: "10:27",
          read: true,
        },
        {
          id: "m3",
          senderId: "sv1",
          senderName: "Nguyễn Văn An",
          content: "Em muốn hỏi về phần công nghệ sử dụng trong đề tài ạ",
          timestamp: "10:28",
          read: true,
        },
        {
          id: "m4",
          senderId: currentUserId,
          senderName: currentUserName,
          content:
            "Em có thể sử dụng React + Node.js, hoặc Vue.js tùy em thích",
          timestamp: "10:29",
          read: true,
        },
        {
          id: "m5",
          senderId: "sv1",
          senderName: "Nguyễn Văn An",
          content: "Em cảm ơn thầy ạ!",
          timestamp: "10:30",
          read: false,
        },
      ],
    },
    {
      id: "c2",
      userId: "sv2",
      userName: "Trần Thị Bình",
      userRole: "MSSV: 20200002",
      lastMessage: "Vâng ạ, em sẽ làm theo hướng dẫn của thầy",
      lastMessageTime: "Hôm qua",
      unread: 0,
      messages: [
        {
          id: "m1",
          senderId: "sv2",
          senderName: "Trần Thị Bình",
          content: "Thầy cho em hỏi deadline nộp báo cáo là khi nào ạ?",
          timestamp: "Hôm qua 14:20",
          read: true,
        },
        {
          id: "m2",
          senderId: currentUserId,
          senderName: currentUserName,
          content: "Deadline là 25/12 em nhé",
          timestamp: "Hôm qua 14:25",
          read: true,
        },
        {
          id: "m3",
          senderId: "sv2",
          senderName: "Trần Thị Bình",
          content: "Vâng ạ, em sẽ làm theo hướng dẫn của thầy",
          timestamp: "Hôm qua 14:30",
          read: true,
        },
      ],
    },
    {
      id: "c3",
      userId: "sv3",
      userName: "Lê Văn Cường",
      userRole: "MSSV: 20200003",
      lastMessage: "Em hiểu rồi ạ, cảm ơn thầy",
      lastMessageTime: "2 ngày trước",
      unread: 0,
      messages: [
        {
          id: "m1",
          senderId: "sv3",
          senderName: "Lê Văn Cường",
          content: "Thầy ơi, em có thể đổi đề tài không ạ?",
          timestamp: "2 ngày trước 09:00",
          read: true,
        },
        {
          id: "m2",
          senderId: currentUserId,
          senderName: currentUserName,
          content: "Em liên hệ với thầy trong giờ hành chính nhé",
          timestamp: "2 ngày trước 09:15",
          read: true,
        },
        {
          id: "m3",
          senderId: "sv3",
          senderName: "Lê Văn Cường",
          content: "Em hiểu rồi ạ, cảm ơn thầy",
          timestamp: "2 ngày trước 09:20",
          read: true,
        },
      ],
    },
  ]);

  const filteredConversations = $derived(
    conversations.filter(
      (conv) =>
        conv.userName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        conv.userRole.toLowerCase().includes(searchTerm.toLowerCase())
    )
  );

  const currentConversation = $derived(
    conversations.find((c) => c.id === selectedConversation)
  );

  function handleSendMessage() {
    if (!messageText.trim() || !selectedConversation) return;

    const newMessage: Message = {
      id: `m${Date.now()}`,
      senderId: currentUserId,
      senderName: currentUserName,
      content: messageText,
      timestamp: new Date().toLocaleTimeString("vi-VN", {
        hour: "2-digit",
        minute: "2-digit",
      }),
      read: false,
    };

    conversations = conversations.map((conv) => {
      if (conv.id === selectedConversation) {
        return {
          ...conv,
          messages: [...conv.messages, newMessage],
          lastMessage: messageText,
          lastMessageTime: "Vừa xong",
        };
      }
      return conv;
    });

    messageText = "";
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
      {#if filteredConversations.length === 0}
        <div class="empty-state">Không tìm thấy cuộc trò chuyện</div>
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
        {#each currentConversation.messages as message (message.id)}
          {@const isOwn = message.senderId === currentUserId}
          <div class="message-wrapper" class:own={isOwn}>
            <div class="message-bubble" class:own={isOwn}>
              {#if !isOwn}
                <div class="sender-name">{message.senderName}</div>
              {/if}
              <div class="message-content">
                <p>{message.content}</p>
              </div>
              <div class="message-footer" class:own={isOwn}>
                <span>{message.timestamp}</span>
                {#if isOwn}
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
    gap: 4px;
    margin-top: 4px;
    font-size: 12px;
    color: #6b7280;
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
