<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import ChatTab from "../../components/ChatTab.svelte";
  import ClassroomRequirements from "../../components/ClassroomRequirements.svelte";
  import CoLecturersModal from "../../components/CoLecturersModal.svelte";
  import type { ProjectRound } from "../../types/project-round";
  import {
    getClassroom,
    getClassPosts,
  } from "../../services/classroom-service";
  import {
    createChannel,
    getChannelBetweenUsers,
  } from "../../services/channel-service";
  import { authStore } from "../../stores/auth-store";
  import type { ClassroomResponse } from "../../dtos/classroom-dto";
  import type { ClassPostResponse } from "../../dtos/class-post-dto";
  import { projectRoundStore } from "../../stores/project-round-store";
  import { classStore } from "../../stores/class-store";

  let { params } = $props<{ params: { id: string } }>();

  let classroom = $state<ClassroomResponse | null>(null);
  let classPosts = $state<ClassPostResponse[]>([]);
  let projectRounds = $state<ProjectRound[]>([]); // Tách riêng khỏi classroom
  let loading = $state(true);
  let showCoLecturersModal = $state(false);
  let lecturerChannelId = $state<string | null>(null);
  let activeTab = $state<"overview" | "students" | "notifications" | "chat">(
    "overview",
  );
  let registrationFilter = $state<"all" | "registered" | "not-registered">(
    "all",
  );
  let unreadNotifications = $state(1);

  // Mock data - sinh viên đã đăng ký category nào
  let registeredCategoryIds = $state<string[]>(["round1"]);
  onMount(async () => {
    try {
      loading = true;
      classroom = await getClassroom(params.id);

      // 1. Gọi API để đổ dữ liệu vào roundsStore
      await projectRoundStore.fetchRounds(params.id);

      // 2. Gán danh sách từ getter "currentRounds" (có chữ 's' ở cuối)
      // chứ không phải "currentRound"
      projectRounds = projectRoundStore.currentRounds;

      classPosts = await getClassPosts(params.id);
    } catch (err) {
      // ... handle error
    } finally {
      loading = false;
    }
  });

  function handleBack() {
    push("/student/classes");
  }

  function handleClassroomChat() {
    if (!classroom?.generalChannelId) {
      alert("Kênh chat lớp chưa được tạo");
      return;
    }

    // Store channel ID and navigate
    localStorage.setItem("openChannelId", classroom.generalChannelId);
    push("/chats");
  }

  async function handleChatWithLecturer() {
    if (!classroom?.lecturer) return;

    console.log("🔍 Lecturer object:", classroom.lecturer);

    // Skip API call for mock data
    if (
      params.id.startsWith("mock-") ||
      classroom.lecturer.userId.startsWith("lecturer")
    ) {
      console.log("Using mock data - fallback to general channel");
      lecturerChannelId = classroom?.generalChannelId || null;
      activeTab = "chat";
      return;
    }

    try {
      // Get current user info
      const currentUser = $authStore.user;
      if (!currentUser) {
        console.error("No current user found");
        return;
      }

      console.log("📤 Current user:", JSON.stringify(currentUser, null, 2));

      // First, check if a channel already exists between these two users
      console.log("🔍 Checking for existing channel...");
      const existingChannel = await getChannelBetweenUsers(
        currentUser.id,
        classroom.lecturer.userId,
      );

      let channelId: string;

      if (existingChannel) {
        console.log("✅ Found existing channel:", existingChannel.id);
        channelId = existingChannel.id;
      } else {
        console.log("➕ No existing channel, creating new one...");
        // Create new DM channel with lecturer (requires 2 members)
        const members = [
          {
            id: currentUser.id,
            full_name: currentUser.fullName || currentUser.fullname || "User",
            email: currentUser.email,
          },
          {
            id: classroom.lecturer.userId,
            full_name: classroom.lecturer.fullName,
            email: classroom.lecturer.email,
          },
        ];

        console.log(
          "📤 Creating channel with members:",
          JSON.stringify(members, null, 2),
        );
        const channel = await createChannel(members);
        console.log("✅ Channel created:", channel);
        channelId = channel.id;
      }

      // Store channel ID in localStorage to pass to Chat page
      localStorage.setItem("openChannelId", channelId);

      // Navigate to Messages page
      push("/chats");
    } catch (err) {
      console.error("Error creating/getting channel with lecturer:", err);
      // If error, still try to navigate to Messages
      push("/chats");
    }
  }

  function handleCategoryClick(categoryId: string) {
    push(`/student/classes/${params.id}/categories/${categoryId}`);
  }

  async function handleLeaveClass() {
    if (
      !confirm(
        "Bạn có chắc chắn muốn rời khỏi lớp học này? Bạn sẽ mất quyền truy cập vào tất cả tài liệu và đề tài.",
      )
    ) {
      return;
    }

    try {
      // MOCK - Comment API call
      // await leaveClassroom(params.id);
      await new Promise((resolve) => setTimeout(resolve, 500));
      alert("Đã rời khỏi lớp học thành công");
      push("/student/classes");
    } catch (err: any) {
      alert(err.message || "Không thể rời khỏi lớp học");
    }
  }

  function getInitials(name: string): string {
    const words = name.split(" ");
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString("vi-VN");
  }

  function downloadAttachment(url: string, fileName: string) {
    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    link.target = "_blank";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<div class="container">
  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Đang tải thông tin lớp học...</p>
    </div>
  {:else if classroom}
    <div class="header-section">
      <button onclick={handleBack} class="back-button">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
        Quay lại danh sách
      </button>

      <div class="header-content">
        <div class="header-left">
          <div class="badges-header">
            <span class="badge badge-code">{classroom.invitationCode}</span>
            <span class="badge badge-semester"
              >{classroom.semester} {classroom.year}</span
            >
          </div>
          <h1 class="class-title">{classroom.name}</h1>
          {#if classroom.description}
            <p class="class-description">{classroom.description}</p>
          {/if}
          <div class="class-meta">
            <span class="meta-item">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
              </svg>
              {projectRounds.length || 0} đề tài
            </span>
            <span class="separator">•</span>
            <span class="meta-item">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="1"></circle>
                <circle cx="12" cy="5" r="1"></circle>
                <circle cx="12" cy="19" r="1"></circle>
              </svg>
              GV: {classroom.lecturer.fullName}
              {#if classroom.coLecturers && classroom.coLecturers.length > 0}
                <button
                  onclick={() => (showCoLecturersModal = true)}
                  class="co-lecturers-btn"
                >
                  + {classroom.coLecturers.length} trợ giảng
                </button>
              {/if}
            </span>
          </div>
        </div>

        <div class="header-actions">
          <button onclick={handleChatWithLecturer} class="btn-chat">
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
              ></path>
            </svg>
            Chat với giảng viên
          </button>
        </div>
      </div>

      <div class="tabs">
        {#each [{ id: "overview", label: "Hạng mục đề tài" }, { id: "students", label: "Sinh viên" }, { id: "notifications", label: "Thông báo" }, { id: "chat", label: "Trò chuyện" }] as tab}
          <button
            onclick={() => (activeTab = tab.id as any)}
            class="tab"
            class:active={activeTab === tab.id}
          >
            {tab.label}
            {#if tab.id === "notifications" && unreadNotifications > 0}
              <span class="badge-unread">{unreadNotifications}</span>
            {/if}
          </button>
        {/each}
      </div>
    </div>

    <div class="content-section">
      {#if activeTab === "overview"}
        <div class="overview-tab">
          <!-- Hiển thị yêu cầu tham gia lớp (whitelist, email restriction, etc.) -->
          <ClassroomRequirements {classroom} />

          <h2 class="section-title">Danh sách đợt đồ án</h2>

          {#if !projectRounds || projectRounds.length === 0}
            <div class="empty-message">
              <p>Chưa có đợt đồ án nào trong lớp học này</p>
            </div>
          {:else}
            <div class="categories-grid">
              {#each projectRounds as round}
                <div
                  class="category-card"
                  onclick={() => handleCategoryClick(round.id)}
                >
                  <div class="category-header">
                    <h3 class="category-name">{round.name}</h3>
                    <span class="status-badge status-active">Đang diễn ra</span>
                  </div>

                  {#if round.description}
                    <p class="category-description">{round.description}</p>
                  {/if}

                  <div class="category-footer">
                    <div class="footer-item">
                      <svg
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"
                        ></rect>
                        <line x1="16" y1="2" x2="16" y2="6"></line>
                        <line x1="8" y1="2" x2="8" y2="6"></line>
                        <line x1="3" y1="10" x2="21" y2="10"></line>
                      </svg>
                      <span
                        >{formatDate(round.startDate)} - {formatDate(
                          round.endDate,
                        )}</span
                      >
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "students"}
        <div class="students-tab">
          <h2 class="section-title">
            Danh sách sinh viên ({classroom.students?.length || 0})
          </h2>

          {#if !classroom.students || classroom.students.length === 0}
            <div class="empty-message">
              <p>Chưa có sinh viên nào trong lớp</p>
            </div>
          {:else}
            <div class="students-list">
              {#each classroom.students as student, index}
                <div class="student-card-horizontal">
                  <div class="student-main">
                    <div class="student-avatar-section">
                      {#if student.avatar && typeof student.avatar === "string" && student.avatar.trim()}
                        <img
                          src={student.avatar}
                          alt={student.fullName}
                          class="student-avatar-large"
                          onerror={(e) => {
                            e.currentTarget.style.display = "none";
                            e.currentTarget.nextElementSibling.style.display =
                              "flex";
                          }}
                        />
                        <div
                          class="student-avatar-large"
                          style="display: none;"
                        >
                          {getInitials(student.fullName)}
                        </div>
                      {:else}
                        <div class="student-avatar-large">
                          {getInitials(student.fullName)}
                        </div>
                      {/if}
                    </div>

                    <div class="student-details">
                      <div class="student-header-row">
                        <span class="student-code-badge"
                          >MSSV: 2052{String(index + 1).padStart(3, "0")}</span
                        >
                      </div>
                      <h4 class="student-name-large">{student.fullName}</h4>
                      <div class="student-meta">
                        <div class="meta-item">
                          <svg
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
                            ></path>
                            <circle cx="12" cy="7" r="4"></circle>
                          </svg>
                          <span>Sinh viên</span>
                        </div>
                        <div class="meta-item">
                          <svg
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
                            <path
                              d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"
                            ></path>
                          </svg>
                          <span>Công nghệ thông tin</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "notifications"}
        <div class="notifications-tab">
          <h2 class="section-title">Thông báo từ giảng viên</h2>

          {#if classPosts.length === 0}
            <div class="empty-message">
              <p>Chưa có thông báo nào</p>
            </div>
          {:else}
            <div class="posts-list">
              {#each classPosts as post}
                <div
                  class="post-card"
                  class:pinned={post.isPinned}
                  class:unread={!(post as any).isRead}
                >
                  <div class="post-badges">
                    {#if !(post as any).isRead}
                      <div class="unread-indicator">
                        <svg
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"
                          ></path>
                          <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
                        </svg>
                        <span>Mới</span>
                      </div>
                    {/if}

                    {#if post.isPinned}
                      <div class="pinned-badge">
                        <svg
                          width="14"
                          height="14"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                        >
                          <path
                            d="M16 9V4h1c.55 0 1-.45 1-1s-.45-1-1-1H7c-.55 0-1 .45-1 1s.45 1 1 1h1v5c0 1.66-1.34 3-3 3v2h5.97v7l1 1 1-1v-7H19v-2c-1.66 0-3-1.34-3-3z"
                          />
                        </svg>
                        Đã ghim
                      </div>
                    {/if}
                  </div>

                  <div class="post-header">
                    <div class="post-author">
                      <div class="author-avatar">
                        {getInitials(post.author.fullName)}
                      </div>
                      <div>
                        <h4 class="author-name">{post.author.fullName}</h4>
                        <p class="post-date">{formatDate(post.createdAt)}</p>
                      </div>
                    </div>
                  </div>

                  <h3 class="post-title">{post.title}</h3>
                  <p class="post-content">{post.content}</p>

                  {#if post.attachments && post.attachments.length > 0}
                    <div class="attachments">
                      <p class="attachments-label">File đính kèm:</p>
                      {#each post.attachments as attachment}
                        <button
                          onclick={() =>
                            downloadAttachment(
                              attachment.fileURL,
                              attachment.fileName,
                            )}
                          class="attachment-item"
                        >
                          <svg
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path
                              d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"
                            />
                          </svg>
                          <span>{attachment.fileName}</span>
                          <span class="file-size">
                            ({Math.round(attachment.fileSize / 1024)} KB)
                          </span>
                        </button>
                      {/each}
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "chat"}
        <div class="chat-tab">
          {#if classroom?.generalChannelId}
            <ChatTab generalChannelId={classroom.generalChannelId} />
          {:else}
            <div class="no-chat">
              <p>Kênh chat chưa được tạo</p>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  {/if}

  <!-- Modal trợ giảng -->
  <CoLecturersModal {classroom} bind:show={showCoLecturersModal} />
</div>

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem;
    color: #6b7280;
  }

  .spinner {
    width: 48px;
    height: 48px;
    border: 4px solid #e5e7eb;
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .header-section {
    background: white;
    border-radius: 0.75rem;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .back-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: none;
    border: none;
    color: #6b7280;
    font-size: 0.875rem;
    cursor: pointer;
    border-radius: 0.5rem;
    transition: all 0.2s;
    margin-bottom: 1rem;
  }

  .back-button:hover {
    background: #f3f4f6;
    color: #111827;
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1.5rem;
  }

  .class-title {
    font-size: 1.875rem;
    font-weight: 700;
    color: #111827;
    margin-bottom: 0.5rem;
  }

  .class-description {
    font-size: 1rem;
    color: #6b7280;
    margin-bottom: 0.75rem;
  }

  .class-meta {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.875rem;
    color: #6b7280;
    flex-wrap: wrap;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .meta-item svg {
    flex-shrink: 0;
  }

  .co-lecturers-btn {
    display: inline-block;
    padding: 2px 8px;
    background: #dbeafe;
    color: #1e40af;
    border: 1px solid #93c5fd;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
    margin-left: 8px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .co-lecturers-btn:hover {
    background: #bfdbfe;
    border-color: #60a5fa;
    transform: translateY(-1px);
  }

  .separator {
    color: #d1d5db;
  }

  .btn-chat {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-chat:hover {
    background: #6d28d9;
  }

  .tabs {
    display: flex;
    gap: 0.5rem;
    border-bottom: 2px solid #e5e7eb;
  }

  .tab {
    position: relative;
    padding: 0.75rem 1.25rem;
    background: none;
    border: none;
    color: #6b7280;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
  }

  .tab:hover {
    color: #3b82f6;
  }

  .tab.active {
    color: #3b82f6;
    border-bottom-color: #3b82f6;
  }

  .badge-unread {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    margin-left: 0.5rem;
    background: #ef4444;
    color: white;
    font-size: 0.7rem;
    font-weight: 600;
    border-radius: 9999px;
  }

  .content-section {
    background: white;
    border-radius: 0.75rem;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #111827;
    margin-bottom: 1.5rem;
  }

  .empty-message {
    text-align: center;
    padding: 3rem 1rem;
    color: #9ca3af;
  }

  .categories-grid {
    display: grid;
    gap: 1rem;
  }

  .category-card {
    padding: 1.25rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .category-card:hover {
    border-color: #3b82f6;
    box-shadow: 0 4px 6px rgba(59, 130, 246, 0.1);
  }

  .category-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.75rem;
  }

  .category-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
  }

  .status-badge {
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: 9999px;
  }

  .status-active {
    background: #dcfce7;
    color: #166534;
  }

  .category-description {
    font-size: 0.875rem;
    color: #6b7280;
    margin-bottom: 1rem;
  }

  .category-footer {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .footer-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  /* Students Tab - Horizontal Cards */
  .students-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .student-card-horizontal {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    padding: 1.5rem;
    transition: all 0.2s;
  }

  .student-card-horizontal:hover {
    border-color: #3b82f6;
    box-shadow: 0 4px 6px rgba(59, 130, 246, 0.1);
  }

  .student-main {
    display: flex;
    gap: 1.5rem;
    align-items: flex-start;
  }

  .student-avatar-section {
    flex-shrink: 0;
  }

  .student-avatar-large {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: #3b82f6;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .student-details {
    flex: 1;
    min-width: 0;
  }

  .student-header-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
  }

  .student-code-badge {
    padding: 0.25rem 0.75rem;
    background: #dbeafe;
    color: #1e40af;
    font-size: 0.875rem;
    border-radius: 0.25rem;
    font-weight: 500;
  }

  .student-name-large {
    font-size: 1.25rem;
    font-weight: 600;
    color: #111827;
    margin-bottom: 0.75rem;
  }

  .student-meta {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .post-card {
    padding: 1.5rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    position: relative;
  }

  .post-card.pinned {
    border-color: #fbbf24;
    background: #fffbeb;
  }

  .post-card.unread {
    background: #eff6ff;
    border-color: #3b82f6;
  }

  .post-badges {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1rem;
    min-height: 28px;
  }

  .unread-indicator {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.25rem 0.75rem;
    background: #3b82f6;
    color: white;
    font-size: 0.75rem;
    font-weight: 600;
    border-radius: 0.25rem;
  }

  .pinned-badge {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.25rem 0.75rem;
    background: #fbbf24;
    color: white;
    font-size: 0.75rem;
    font-weight: 600;
    border-radius: 0.25rem;
    margin-left: auto;
  }

  .post-header {
    margin-bottom: 1rem;
  }

  .post-author {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .author-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #3b82f6;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    font-weight: 600;
  }

  .author-name {
    font-size: 0.875rem;
    font-weight: 600;
    color: #111827;
  }

  .post-date {
    font-size: 0.75rem;
    color: #9ca3af;
  }

  .post-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
    margin-bottom: 0.5rem;
  }

  .post-content {
    font-size: 0.875rem;
    color: #4b5563;
    line-height: 1.6;
    margin-bottom: 1rem;
  }

  .attachments {
    padding-top: 1rem;
    border-top: 1px solid #e5e7eb;
  }

  .attachments-label {
    font-size: 0.75rem;
    font-weight: 600;
    color: #6b7280;
    text-transform: uppercase;
    margin-bottom: 0.5rem;
  }

  .attachment-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    background: #f3f4f6;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    color: #374151;
    cursor: pointer;
    transition: all 0.2s;
    margin-bottom: 0.5rem;
  }

  .attachment-item:hover {
    background: #e5e7eb;
    border-color: #3b82f6;
  }

  .file-size {
    color: #9ca3af;
    font-size: 0.75rem;
  }
</style>
