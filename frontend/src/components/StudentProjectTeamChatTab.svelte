<script lang="ts">
  interface Message {
    id: string;
    senderId: string;
    senderName: string;
    content: string;
    timestamp: string;
    isMe: boolean;
  }

  interface Team {
    id: string;
    name: string;
    leaderId: string;
    members: { id: string; name: string; studentId: string }[];
  }

  let { team, currentStudentId } = $props<{
    team: Team;
    currentStudentId: string;
  }>();

  let message = $state("");
  let messages = $state<Message[]>([
    {
      id: "1",
      senderId: team.leaderId,
      senderName: "Nhóm trưởng",
      content: "Chào mọi người! Chúng ta hãy thảo luận về tiến độ dự án nhé.",
      timestamp: "10:30",
      isMe: team.leaderId === currentStudentId,
    },
    {
      id: "2",
      senderId: team.members[0]?.id || "2",
      senderName: team.members[0]?.name || "Member 1",
      content: "Em đã hoàn thành phần thiết kế database rồi ạ!",
      timestamp: "10:35",
      isMe: team.members[0]?.id === currentStudentId,
    },
    {
      id: "3",
      senderId: currentStudentId,
      senderName: "Bạn",
      content: "Tuyệt vời! Em sẽ bắt đầu code phần backend.",
      timestamp: "10:40",
      isMe: true,
    },
  ]);

  function handleSendMessage() {
    if (!message.trim()) return;

    const newMessage: Message = {
      id: Date.now().toString(),
      senderId: currentStudentId,
      senderName: "Bạn",
      content: message,
      timestamp: new Date().toLocaleTimeString("vi-VN", {
        hour: "2-digit",
        minute: "2-digit",
      }),
      isMe: true,
    };

    messages = [...messages, newMessage];
    message = "";
  }

  function handleKeyPress(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  }
</script>

<div class="chat-container">
  <!-- Chat Header -->
  <div class="chat-header">
    <h3 class="chat-title">Chat nhóm: {team.name}</h3>
    <p class="chat-subtitle">{team.members.length + 1} thành viên</p>
  </div>

  <!-- Messages -->
  <div class="messages-container">
    {#each messages as msg (msg.id)}
      <div class="message-wrapper" class:own={msg.isMe}>
        <div class="message-content" class:own={msg.isMe}>
          <!-- Avatar -->
          <div
            class="message-avatar"
            class:leader={msg.senderId === team.leaderId}
          >
            {msg.senderName.charAt(0)}
          </div>

          <!-- Message Bubble -->
          <div class="message-bubble-wrapper">
            <div class="message-info" class:own={msg.isMe}>
              <span class="message-sender">{msg.senderName}</span>
              <span class="message-time">{msg.timestamp}</span>
            </div>
            <div class="message-bubble" class:own={msg.isMe}>
              {msg.content}
            </div>
          </div>
        </div>
      </div>
    {/each}
  </div>

  <!-- Input Area -->
  <div class="input-container">
    <div class="input-wrapper">
      <textarea
        bind:value={message}
        onkeypress={handleKeyPress}
        placeholder="Nhập tin nhắn..."
        class="message-input"
        rows="1"
      ></textarea>
      <div class="input-actions">
        <button class="action-btn">
          <svg
            width="16"
            height="16"
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
        <button class="action-btn">
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <circle cx="9" cy="9" r="2"></circle>
            <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
          </svg>
        </button>
        <button class="action-btn">
          <svg
            width="16"
            height="16"
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
      </div>
    </div>
    <button
      onclick={handleSendMessage}
      disabled={!message.trim()}
      class="send-btn"
      class:disabled={!message.trim()}
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

<style>
  .chat-container {
    display: flex;
    flex-direction: column;
    height: 600px;
  }

  /* Chat Header */
  .chat-header {
    padding-bottom: 16px;
    border-bottom: 1px solid #e5e7eb;
  }

  .chat-title {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
  }

  .chat-subtitle {
    font-size: 14px;
    color: #64748b;
  }

  /* Messages */
  .messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 16px 0;
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

  .message-content {
    display: flex;
    gap: 12px;
    max-width: 70%;
  }

  .message-content.own {
    flex-direction: row-reverse;
  }

  .message-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #dbeafe;
    color: #1e40af;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    flex-shrink: 0;
  }

  .message-avatar.leader {
    background: #fef3c7;
    color: #92400e;
  }

  .message-bubble-wrapper {
    flex: 1;
  }

  .message-info {
    display: flex;
    align-items: baseline;
    gap: 8px;
    margin-bottom: 4px;
  }

  .message-info.own {
    flex-direction: row-reverse;
  }

  .message-sender {
    font-size: 14px;
    font-weight: 500;
    color: #1e293b;
  }

  .message-time {
    font-size: 12px;
    color: #94a3b8;
  }

  .message-bubble {
    padding: 12px 16px;
    border-radius: 12px;
    background: #f1f5f9;
    color: #111827;
    word-wrap: break-word;
  }

  .message-bubble.own {
    background: #3b82f6;
    color: white;
  }

  /* Input Area */
  .input-container {
    padding-top: 16px;
    border-top: 1px solid #e5e7eb;
    display: flex;
    align-items: flex-end;
    gap: 8px;
  }

  .input-wrapper {
    flex: 1;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 12px;
  }

  .message-input {
    width: 100%;
    background: transparent;
    border: none;
    outline: none;
    font-family: inherit;
    font-size: 14px;
    resize: none;
    min-height: 24px;
    max-height: 120px;
  }

  .input-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }

  .action-btn {
    padding: 6px;
    color: #64748b;
    background: none;
    border: none;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .action-btn:hover {
    color: #334155;
    background: #f1f5f9;
  }

  .send-btn {
    padding: 12px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;
    flex-shrink: 0;
  }

  .send-btn:hover:not(.disabled) {
    background: #2563eb;
  }

  .send-btn.disabled {
    background: #e5e7eb;
    color: #9ca3af;
    cursor: not-allowed;
  }
</style>
