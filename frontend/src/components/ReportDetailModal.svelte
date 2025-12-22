<script lang="ts">
  interface Comment {
    id: string;
    authorId: string;
    authorName: string;
    authorType: "teacher" | "student";
    content: string;
    timestamp: string;
    parentId?: string;
    attachmentFile?: string;
  }

  let { reportTitle, submittedFile, grade, feedback, onClose } = $props<{
    reportTitle: string;
    submittedFile?: {
      name: string;
      size: string;
      uploadTime: string;
    };
    grade?: number;
    feedback?: {
      text: string;
      file?: string;
    };
    onClose: () => void;
  }>();

  let comments = $state<Comment[]>([
    {
      id: "1",
      authorId: "gv1",
      authorName: "TS. Trần Văn B",
      authorType: "teacher",
      content:
        "Các em xem lại comment của thầy trong file đính kèm nhé. Chú ý sửa lại phần thiết kế CSDL cho tuần sau.",
      timestamp: "10:45, 25/03/2024",
    },
    {
      id: "2",
      authorId: "sv1",
      authorName: "Nguyễn Văn An (Bạn)",
      authorType: "student",
      content:
        "Dạ chúng em cảm ơn thầy ạ. Nhóm sẽ sửa lại phần này trong báo cáo tuần 2 ạ.",
      timestamp: "11:02, 25/03/2024",
      parentId: "1",
    },
    {
      id: "3",
      authorId: "gv1",
      authorName: "TS. Trần Văn B",
      authorType: "teacher",
      content: "Ok em.",
      timestamp: "09:15, 26/03/2024",
      parentId: "2",
    },
  ]);

  let newComment = $state("");
  let replyingTo = $state<string | null>(null);

  function handleSubmitComment() {
    if (!newComment.trim()) return;

    const comment: Comment = {
      id: Date.now().toString(),
      authorId: "sv1",
      authorName: "Nguyễn Văn An (Bạn)",
      authorType: "student",
      content: newComment,
      timestamp: new Date().toLocaleString("vi-VN"),
      parentId: replyingTo || undefined,
    };

    comments = [...comments, comment];
    newComment = "";
    replyingTo = null;
  }

  function getReplyQuote(commentId: string) {
    const comment = comments.find((c) => c.id === commentId);
    if (!comment) return "";
    return comment.content.length > 50
      ? `"${comment.content.substring(0, 50)}..."`
      : `"${comment.content}"`;
  }
</script>

<div class="modal-overlay" onclick={onClose}>
  <div class="modal-container" onclick={(e) => e.stopPropagation()}>
    <!-- Header -->
    <div class="modal-header">
      <h3 class="modal-title">{reportTitle}</h3>
      <button onclick={onClose} class="close-button">
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <!-- Content - Scrollable -->
    <div class="modal-content">
      <!-- File & Grade Section -->
      <div class="info-section">
        <h4 class="section-title">Thông tin nộp bài & Kết quả</h4>

        <!-- Submitted File -->
        {#if submittedFile}
          <div class="file-info-card">
            <div class="file-header">
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                class="file-icon"
              >
                <path
                  d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                ></path>
                <polyline points="14 2 14 8 20 8"></polyline>
              </svg>
              <div class="file-details">
                <div class="file-name">{submittedFile.name}</div>
                <div class="file-meta">
                  {submittedFile.size} - Nộp lúc {submittedFile.uploadTime}
                </div>
              </div>
            </div>

            <!-- Grade & Feedback -->
            {#if grade !== undefined}
              <div class="grade-feedback-section">
                <div class="grade-display">
                  <div class="grade-label">Điểm</div>
                  <div class="grade-value">{grade.toFixed(1)}</div>
                </div>

                <div class="feedback-content">
                  <div class="feedback-label">Nhận xét của Giảng viên</div>
                  {#if feedback}
                    <p class="feedback-text">{feedback.text}</p>
                    {#if feedback.file}
                      <button class="feedback-file-button">
                        <svg
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                          ></path>
                          <polyline points="7 10 12 15 17 10"></polyline>
                          <line x1="12" y1="15" x2="12" y2="3"></line>
                        </svg>
                        {feedback.file}
                      </button>
                    {/if}
                  {/if}
                </div>
              </div>
            {/if}
          </div>
        {/if}
      </div>

      <!-- Comments Section -->
      <div class="comments-section">
        <h4 class="section-title">Khu vực Trao đổi / Phản biện</h4>

        <div class="comments-list">
          {#each comments as comment (comment.id)}
            {@const isTeacher = comment.authorType === "teacher"}
            {@const isReply = !!comment.parentId}
            <div class="comment-wrapper" class:reply={isReply}>
              <div class="comment-container">
                <div class="comment-avatar" class:teacher={isTeacher}>
                  <span>{isTeacher ? "GV" : comment.authorName.charAt(0)}</span>
                </div>

                <div class="comment-content">
                  <div class="comment-bubble" class:teacher={isTeacher}>
                    <div class="comment-meta">
                      <span class="comment-author" class:teacher={isTeacher}>
                        {comment.authorName}
                      </span>
                      <span class="comment-time">{comment.timestamp}</span>
                    </div>

                    <!-- Reply Quote -->
                    {#if comment.parentId}
                      <div class="reply-quote">
                        <p>{getReplyQuote(comment.parentId)}</p>
                      </div>
                    {/if}

                    <p class="comment-text">{comment.content}</p>
                  </div>

                  <button
                    onclick={() => (replyingTo = comment.id)}
                    class="reply-button"
                  >
                    Trả lời
                  </button>

                  <!-- Pending Reply Badge -->
                  {#if isReply && comment.authorType === "student" && !comments.some((c) => c.parentId === comment.id && c.authorType === "teacher")}
                    <div class="pending-badge">
                      <svg
                        width="14"
                        height="14"
                        fill="currentColor"
                        viewBox="0 0 15 16"
                      >
                        <path
                          d="M7.5 0C3.36 0 0 3.36 0 7.5S3.36 15 7.5 15 15 11.64 15 7.5 11.64 0 7.5 0zm.75 11.25h-1.5v-1.5h1.5v1.5zm0-3h-1.5v-4.5h1.5v4.5z"
                        />
                      </svg>
                      Pending Reply
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="modal-footer">
      {#if replyingTo}
        <div class="replying-to">
          <span>Đang trả lời: {getReplyQuote(replyingTo)}</span>
          <button onclick={() => (replyingTo = null)} class="cancel-reply">
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      {/if}

      <div class="input-area">
        <textarea
          bind:value={newComment}
          placeholder="Viết phản hồi của bạn..."
          class="comment-input"
          rows="2"
        ></textarea>
        <div class="action-buttons">
          <button class="attach-button">
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
          <button
            onclick={handleSubmitComment}
            disabled={!newComment.trim()}
            class="send-button"
            class:disabled={!newComment.trim()}
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
    </div>
  </div>
</div>

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 16px;
  }

  .modal-container {
    background: white;
    border-radius: 8px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 600px;
    height: 90vh;
    display: flex;
    flex-direction: column;
  }

  /* Header */
  .modal-header {
    border-bottom: 1px solid #e2e8f0;
    padding: 16px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .modal-title {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
  }

  .close-button {
    padding: 4px;
    color: #64748b;
    background: none;
    border: none;
    cursor: pointer;
    transition: color 0.2s;
  }

  .close-button:hover {
    color: #334155;
  }

  /* Content */
  .modal-content {
    flex: 1;
    overflow-y: auto;
  }

  .info-section,
  .comments-section {
    padding: 24px;
  }

  .section-title {
    font-size: 16px;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 16px;
  }

  .file-info-card {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    padding: 16px;
  }

  .file-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 24px;
  }

  .file-icon {
    color: #3b82f6;
    flex-shrink: 0;
  }

  .file-details {
    flex: 1;
  }

  .file-name {
    font-size: 14px;
    color: #1e293b;
    margin-bottom: 2px;
  }

  .file-meta {
    font-size: 12px;
    color: #64748b;
  }

  .grade-feedback-section {
    border-top: 1px solid #e2e8f0;
    padding-top: 16px;
    display: flex;
    gap: 16px;
  }

  .grade-display {
    text-align: center;
  }

  .grade-label {
    font-size: 12px;
    color: #64748b;
    margin-bottom: 4px;
  }

  .grade-value {
    font-size: 48px;
    font-weight: 600;
    color: #3b82f6;
  }

  .feedback-content {
    flex: 1;
    border-left: 1px solid #e2e8f0;
    padding-left: 16px;
  }

  .feedback-label {
    font-size: 12px;
    color: #64748b;
    margin-bottom: 8px;
  }

  .feedback-text {
    font-size: 14px;
    color: #334155;
    margin-bottom: 12px;
  }

  .feedback-file-button {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #3b82f6;
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }

  /* Comments */
  .comments-section {
    border-top: 1px solid #e2e8f0;
  }

  .comments-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .comment-wrapper {
  }

  .comment-wrapper.reply {
    margin-left: 48px;
  }

  .comment-container {
    display: flex;
    gap: 12px;
  }

  .comment-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: #f1f5f9;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .comment-avatar.teacher {
    background: #e2e8f0;
  }

  .comment-avatar span {
    font-size: 12px;
    color: #64748b;
  }

  .comment-avatar.teacher span {
    color: #3b82f6;
  }

  .comment-content {
    flex: 1;
  }

  .comment-bubble {
    border-radius: 8px;
    border-top-left-radius: 0;
    padding: 12px;
    background: #f1f5f9;
  }

  .comment-bubble.teacher {
    background: rgba(59, 130, 246, 0.1);
  }

  .comment-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }

  .comment-author {
    font-size: 14px;
    color: #1e293b;
  }

  .comment-author.teacher {
    color: #3b82f6;
  }

  .comment-time {
    font-size: 12px;
    color: #64748b;
  }

  .reply-quote {
    border-left: 2px solid #3b82f6;
    padding-left: 8px;
    margin-bottom: 8px;
    opacity: 0.7;
  }

  .reply-quote p {
    font-size: 12px;
    color: #334155;
    margin: 0;
  }

  .comment-text {
    font-size: 14px;
    color: #1e293b;
    margin: 0;
  }

  .reply-button {
    font-size: 12px;
    color: #64748b;
    background: none;
    border: none;
    cursor: pointer;
    margin-top: 4px;
    margin-left: 12px;
    transition: color 0.2s;
  }

  .reply-button:hover {
    color: #334155;
  }

  .pending-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    background: #fed7aa;
    color: #ea580c;
    border-radius: 16px;
    font-size: 12px;
    margin-top: 8px;
    margin-left: 12px;
  }

  /* Footer */
  .modal-footer {
    border-top: 1px solid #e2e8f0;
    background: white;
    padding: 16px;
  }

  .replying-to {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #64748b;
    margin-bottom: 8px;
  }

  .cancel-reply {
    color: #94a3b8;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    transition: color 0.2s;
  }

  .cancel-reply:hover {
    color: #64748b;
  }

  .input-area {
    display: flex;
    gap: 8px;
  }

  .comment-input {
    flex: 1;
    padding: 12px;
    background: #f8fafc;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    resize: none;
    font-family: inherit;
    font-size: 14px;
    outline: none;
  }

  .comment-input:focus {
    border-color: #3b82f6;
  }

  .action-buttons {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .attach-button,
  .send-button {
    padding: 8px;
    border-radius: 8px;
    border: none;
    cursor: pointer;
    transition: all 0.2s;
  }

  .attach-button {
    color: #64748b;
    background: none;
  }

  .attach-button:hover {
    background: #f1f5f9;
  }

  .send-button {
    background: #3b82f6;
    color: white;
  }

  .send-button:hover:not(.disabled) {
    background: #2563eb;
  }

  .send-button.disabled {
    background: #e5e7eb;
    color: #9ca3af;
    cursor: not-allowed;
  }
</style>
