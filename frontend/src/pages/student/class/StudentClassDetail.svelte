<script lang="ts">
  import { push } from "svelte-spa-router";
  import type { CategoryItem } from "../../../types/category";
  import type { StudentInClass } from "../../../types/student";

  type TabType = "overview" | "students" | "notifications" | "chat";
  type RegistrationFilter = "all" | "registered" | "not-registered";

  interface StudentClassDetailProps {
    classId: string;
  }

  let { classId }: StudentClassDetailProps = $props();

  let activeTab = $state<TabType>("overview");
  let selectedStudent = $state<StudentInClass | null>(null);
  let registrationFilter = $state<RegistrationFilter>("all");

  // Mock data - Thông tin lớp học
  const classData = {
    id: classId,
    code: "IT4441",
    name: "Phát triển ứng dụng Web",
    semester: "HK1 2024-2025",
    studentCount: 45,
    projectCount: 15,
    instructor: "TS. Nguyễn Văn A",
    university: "Đại học Bách Khoa Hà Nội",
    description: "Môn học về phát triển ứng dụng web hiện đại",
  };

  // Mock categories
  const categories: CategoryItem[] = [
    {
      id: "1",
      name: "Đồ án 1",
      description: "Đồ án phát triển ứng dụng web cơ bản",
      projectCount: 5,
      startDate: "01/01/2025",
      endDate: "31/01/2025",
      registeredCount: 10,
    },
    {
      id: "2",
      name: "Đồ án 2",
      description: "Đồ án phát triển ứng dụng web nâng cao",
      projectCount: 6,
      startDate: "01/02/2025",
      endDate: "28/02/2025",
      registeredCount: 12,
    },
    {
      id: "3",
      name: "Đồ án cuối kỳ",
      description: "Đồ án tổng hợp kiến thức cả học kỳ",
      projectCount: 4,
      startDate: "01/03/2025",
      endDate: "31/03/2025",
      registeredCount: 8,
    },
  ];

  // Mock students
  const students: StudentInClass[] = [
    {
      id: "1",
      name: "Nguyễn Văn A",
      studentCode: "SV001",
      email: "nguyenvana@student.edu.vn",
      enrolledProjects: [],
    },
    {
      id: "2",
      name: "Trần Thị B",
      studentCode: "SV002",
      email: "tranthib@student.edu.vn",
      enrolledProjects: [],
    },
  ];

  // Mock notifications
  const classNotifications = [
    {
      id: 1,
      title: "Thông báo mở đăng ký đề tài",
      content: "Hạng mục Đồ án 1 đã mở đăng ký",
      time: "2 giờ trước",
      isRead: false,
    },
    {
      id: 2,
      title: "Deadline báo cáo",
      content: "Nhắc nhở nộp báo cáo giữa kỳ trước 15/12",
      time: "1 ngày trước",
      isRead: true,
    },
    {
      id: 3,
      title: "Cập nhật thông tin lớp",
      content: "Đã cập nhật danh sách sinh viên",
      time: "3 ngày trước",
      isRead: true,
    },
  ];

  // Mock - sinh viên đã đăng ký hạng mục nào
  const currentStudentRegisteredCategories = ["1"];

  const filteredCategories = $derived(() => {
    if (registrationFilter === "all") return categories;
    const registered = categories.filter((c) =>
      currentStudentRegisteredCategories.includes(c.id)
    );
    const notRegistered = categories.filter(
      (c) => !currentStudentRegisteredCategories.includes(c.id)
    );
    return registrationFilter === "registered" ? registered : notRegistered;
  });

  const registeredCount = $derived(
    categories.filter((c) => currentStudentRegisteredCategories.includes(c.id))
      .length
  );
  const notRegisteredCount = $derived(categories.length - registeredCount);

  function handleBack(): void {
    push("/classes");
  }

  function handleSelectCategory(category: CategoryItem): void {
    push(`/classes/${classId}/categories/${category.id}`);
  }

  function handleChatWithInstructor(): void {
    // TODO: Implement chat with instructor
    console.log("Chat with instructor");
  }
</script>

<div class="page-container">
  <!-- Header -->
  <div class="header-card">
    <button onclick={handleBack} class="back-button">
      <img src="/arrow-left.svg" alt="back" width="20" height="20" />
      Quay lại danh sách lớp học
    </button>

    <div class="header-content">
      <div class="header-info">
        <div class="header-badges">
          <span class="code-badge">{classData.code}</span>
          <span class="semester-text">{classData.semester}</span>
        </div>
        <h1 class="class-title">{classData.name}</h1>
        <p class="class-description">{classData.description}</p>

        <div class="stats-row">
          <div class="stat-item">
            <img src="/users.svg" alt="students" width="16" height="16" />
            <span>{classData.studentCount} sinh viên</span>
          </div>
          <div class="stat-item">
            <img
              src="/folder-kanban.svg"
              alt="projects"
              width="16"
              height="16"
            />
            <span>{classData.projectCount} đề tài</span>
          </div>
          <div class="stat-item">
            <img src="/info.svg" alt="instructor" width="16" height="16" />
            <span>GV: {classData.instructor}</span>
          </div>
        </div>
      </div>

      <button onclick={handleChatWithInstructor} class="chat-button">
        <img src="/message-square.svg" alt="chat" width="20" height="20" />
        Chat với giảng viên
      </button>
    </div>
  </div>

  <!-- Tabs -->
  <div class="tabs-card">
    <div class="tabs-header">
      <button
        onclick={() => (activeTab = "overview")}
        class="tab-button"
        class:active={activeTab === "overview"}
      >
        Hạng mục đề tài
      </button>
      <button
        onclick={() => (activeTab = "students")}
        class="tab-button"
        class:active={activeTab === "students"}
      >
        Sinh viên ({classData.studentCount})
      </button>
      <button
        onclick={() => (activeTab = "notifications")}
        class="tab-button"
        class:active={activeTab === "notifications"}
      >
        Thông báo
        {#if classNotifications.filter((n) => !n.isRead).length > 0}
          <span class="notification-badge">
            {classNotifications.filter((n) => !n.isRead).length}
          </span>
        {/if}
      </button>
      <button
        onclick={() => (activeTab = "chat")}
        class="tab-button"
        class:active={activeTab === "chat"}
      >
        Trò chuyện
      </button>
    </div>

    <div class="tab-content">
      {#if activeTab === "overview"}
        <div class="overview-tab">
          <div class="tab-header">
            <h2 class="tab-title">Danh sách hạng mục</h2>

            <div class="filter-group">
              <button
                onclick={() => (registrationFilter = "all")}
                class="filter-button"
                class:active={registrationFilter === "all"}
              >
                Tất cả ({categories.length})
              </button>
              <button
                onclick={() => (registrationFilter = "registered")}
                class="filter-button"
                class:active={registrationFilter === "registered"}
              >
                Đã đăng ký ({registeredCount})
              </button>
              <button
                onclick={() => (registrationFilter = "not-registered")}
                class="filter-button"
                class:active={registrationFilter === "not-registered"}
              >
                Chưa đăng ký ({notRegisteredCount})
              </button>
            </div>
          </div>

          <div class="categories-list">
            {#each filteredCategories() as category (category.id)}
              {@const isRegistered =
                currentStudentRegisteredCategories.includes(category.id)}
              <div
                class="category-card"
                onclick={() => handleSelectCategory(category)}
                role="button"
                tabindex="0"
              >
                <div class="category-header">
                  <div class="category-info">
                    <div class="category-badges">
                      <h3 class="category-name">{category.name}</h3>
                      <span class="status-badge ongoing">Đang diễn ra</span>
                      {#if isRegistered}
                        <span class="registered-badge">Đã đăng ký</span>
                      {/if}
                    </div>
                    <p class="category-description">{category.description}</p>
                  </div>
                </div>

                <div class="category-footer">
                  <span>Từ {category.startDate} đến {category.endDate}</span>
                  <span>•</span>
                  <span>{category.projectCount} đề tài</span>
                  <span>•</span>
                  <span>{category.registeredCount} sinh viên đã đăng ký</span>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else if activeTab === "students"}
        <div class="students-tab">
          <h2 class="tab-title">Danh sách sinh viên</h2>
          <div class="students-list">
            {#each students as student (student.id)}
              <div class="student-card">
                <div class="student-avatar">{student.name.charAt(0)}</div>
                <div class="student-info">
                  <div class="student-name">{student.name}</div>
                  <div class="student-code">MSSV: {student.studentCode}</div>
                </div>
                <div class="student-email">{student.email}</div>
              </div>
            {/each}
          </div>
        </div>
      {:else if activeTab === "notifications"}
        <div class="notifications-tab">
          <h2 class="tab-title">Thông báo trong lớp</h2>
          <div class="notifications-list">
            {#each classNotifications as notification (notification.id)}
              <div
                class="notification-card"
                class:unread={!notification.isRead}
              >
                <img
                  src="/bell.svg"
                  alt="notification"
                  width="20"
                  height="20"
                  class="notification-icon"
                  class:unread={!notification.isRead}
                />
                <div class="notification-content">
                  <div class="notification-header">
                    <h4
                      class="notification-title"
                      class:unread={!notification.isRead}
                    >
                      {notification.title}
                    </h4>
                    <span class="notification-time">{notification.time}</span>
                  </div>
                  <p class="notification-text">{notification.content}</p>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else if activeTab === "chat"}
        <div class="chat-tab">
          <p class="chat-placeholder">Chat feature coming soon...</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px;
  }

  .header-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 24px;
    margin-bottom: 24px;
  }

  .back-button {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 0.875rem;
    margin-bottom: 16px;
    transition: color 0.2s;
  }

  .back-button:hover {
    color: #1e293b;
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 24px;
  }

  .header-info {
    flex: 1;
  }

  .header-badges {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
  }

  .code-badge {
    padding: 4px 12px;
    background: #dbeafe;
    color: #1e40af;
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .semester-text {
    color: #64748b;
    font-size: 0.875rem;
  }

  .class-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 8px;
  }

  .class-description {
    color: #64748b;
    margin-bottom: 16px;
    line-height: 1.5;
  }

  .stats-row {
    display: flex;
    align-items: center;
    gap: 24px;
    flex-wrap: wrap;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .chat-button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: #7c3aed;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    white-space: nowrap;
    transition: background 0.2s;
  }

  .chat-button:hover {
    background: #6d28d9;
  }

  .tabs-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .tabs-header {
    display: flex;
    gap: 32px;
    padding: 0 24px;
    border-bottom: 1px solid #e2e8f0;
  }

  .tab-button {
    position: relative;
    padding: 16px 0;
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    color: #64748b;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .tab-button:hover {
    color: #1e293b;
  }

  .tab-button.active {
    color: #2563eb;
    border-bottom-color: #2563eb;
  }

  .notification-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    padding: 0 6px;
    background: #ef4444;
    color: white;
    border-radius: 10px;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .tab-content {
    padding: 24px;
  }

  .tab-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 16px;
  }

  .tab-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .filter-group {
    display: flex;
    gap: 8px;
    background: #f1f5f9;
    padding: 4px;
    border-radius: 8px;
  }

  .filter-button {
    padding: 8px 16px;
    background: transparent;
    border: none;
    border-radius: 6px;
    color: #64748b;
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.2s;
  }

  .filter-button:hover {
    color: #1e293b;
  }

  .filter-button.active {
    background: white;
    color: #2563eb;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .categories-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .category-card {
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 20px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .category-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    border-color: #3b82f6;
  }

  .category-header {
    margin-bottom: 12px;
  }

  .category-badges {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .category-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1e293b;
  }

  .status-badge {
    padding: 4px 12px;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .status-badge.ongoing {
    background: #dcfce7;
    color: #166534;
  }

  .registered-badge {
    padding: 4px 12px;
    background: #2563eb;
    color: white;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .category-description {
    color: #64748b;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  .category-footer {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .students-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .student-card {
    display: flex;
    align-items: center;
    gap: 12px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 16px;
    transition: border-color 0.2s;
  }

  .student-card:hover {
    border-color: #3b82f6;
  }

  .student-avatar {
    width: 40px;
    height: 40px;
    background: #dbeafe;
    color: #1e40af;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
  }

  .student-info {
    flex: 1;
  }

  .student-name {
    font-weight: 500;
    color: #1e293b;
  }

  .student-code {
    font-size: 0.875rem;
    color: #64748b;
  }

  .student-email {
    color: #64748b;
    font-size: 0.875rem;
  }

  .notifications-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .notification-card {
    display: flex;
    gap: 12px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 16px;
    background: white;
  }

  .notification-card.unread {
    background: #eff6ff;
    border-color: #bfdbfe;
  }

  .notification-icon {
    flex-shrink: 0;
    opacity: 0.5;
  }

  .notification-icon.unread {
    opacity: 1;
    filter: brightness(0) saturate(100%) invert(40%) sepia(93%) saturate(1664%)
      hue-rotate(204deg) brightness(100%) contrast(91%);
  }

  .notification-content {
    flex: 1;
  }

  .notification-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 4px;
  }

  .notification-title {
    font-weight: 500;
    color: #1e293b;
    font-size: 0.9375rem;
  }

  .notification-title.unread {
    color: #1e3a8a;
  }

  .notification-time {
    font-size: 0.75rem;
    color: #64748b;
    white-space: nowrap;
    margin-left: 8px;
  }

  .notification-text {
    color: #64748b;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  .chat-placeholder {
    text-align: center;
    color: #64748b;
    padding: 48px 24px;
  }

  @media (max-width: 768px) {
    .page-container {
      padding: 16px;
    }

    .header-content {
      flex-direction: column;
    }

    .chat-button {
      width: 100%;
      justify-content: center;
    }

    .tabs-header {
      gap: 16px;
      overflow-x: auto;
    }

    .tab-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .filter-group {
      width: 100%;
    }

    .filter-button {
      flex: 1;
    }
  }
</style>
