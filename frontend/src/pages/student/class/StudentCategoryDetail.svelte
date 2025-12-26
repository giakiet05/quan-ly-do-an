<script lang="ts">
  import { push } from "svelte-spa-router";
  import type { Project } from "../../../types/project";

  type TabType = "projects" | "reports";

  interface StudentCategoryDetailProps {
    classId: string;
    categoryId: string;
  }

  let { classId, categoryId }: StudentCategoryDetailProps = $props();

  let activeTab = $state<TabType>("projects");

  // Mock data - thông tin lớp học
  const classData = {
    id: classId,
    code: "IT4441",
    name: "Phát triển ứng dụng Web",
  };

  // Mock data - thông tin hạng mục
  const category = {
    id: categoryId,
    name: "Đồ án 1",
    description: "Đồ án phát triển ứng dụng web cơ bản",
    startDate: "01/01/2025",
    endDate: "31/01/2025",
    status: "ongoing" as "upcoming" | "ongoing" | "completed",
  };

  // Mock projects
  const projects: Project[] = [
    {
      id: "p1",
      categoryId: categoryId,
      name: "Hệ thống quản lý thư viện trực tuyến",
      description:
        "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
      maxStudents: 3,
      currentStudents: 0,
      instructor: "TS. Nguyễn Văn A",
      tags: ["Web", "React", "Node.js", "MongoDB"],
      status: "available",
      minTeamMembers: 2,
      maxTeamMembers: 3,
    },
    {
      id: "p2",
      categoryId: categoryId,
      name: "Ứng dụng di động quản lý chi tiêu cá nhân",
      description:
        "Ứng dụng mobile giúp người dùng theo dõi thu chi, lập kế hoạch tài chính",
      maxStudents: 2,
      currentStudents: 2,
      instructor: "ThS. Trần Thị B",
      tags: ["Mobile", "React Native", "Firebase"],
      status: "full",
    },
    {
      id: "p3",
      categoryId: categoryId,
      name: "Hệ thống quản lý học sinh",
      description: "Xây dựng hệ thống quản lý học sinh với đầy đủ tính năng",
      maxStudents: 4,
      currentStudents: 0,
      instructor: "TS. Nguyễn Văn C",
      tags: ["Web", "Vue.js", "Spring Boot"],
      status: "available",
    },
  ];

  // Mock reports
  const reports = [
    {
      id: 1,
      title: "Báo cáo tiến độ lần 1",
      deadline: "15/12/2024",
      status: "submitted" as "submitted" | "pending" | "upcoming",
      submittedAt: "14/12/2024",
    },
    {
      id: 2,
      title: "Báo cáo giữa kỳ",
      deadline: "20/12/2024",
      status: "pending" as "submitted" | "pending" | "upcoming",
      submittedAt: null,
    },
    {
      id: 3,
      title: "Báo cáo cuối kỳ",
      deadline: "10/01/2025",
      status: "upcoming" as "submitted" | "pending" | "upcoming",
      submittedAt: null,
    },
  ];

  // Mock - giả sử sinh viên đã đăng ký đề tài đầu tiên
  const myRegisteredProject = null; // projects[0];
  const hasRegistered = !!myRegisteredProject;
  const isRegistrationOpen = category.status === "ongoing";

  const statusConfig = {
    upcoming: {
      label: "Sắp diễn ra",
      color: "text-gray-600",
      bg: "bg-gray-50",
    },
    ongoing: {
      label: "Đang diễn ra",
      color: "text-green-600",
      bg: "bg-green-50",
    },
    completed: {
      label: "Đã kết thúc",
      color: "text-blue-600",
      bg: "bg-blue-50",
    },
  }[category.status];

  function handleBack(): void {
    push(`/classes/${classId}`);
  }

  function handleSelectProject(project: Project): void {
    push(`/classes/${classId}/categories/${categoryId}/projects/${project.id}`);
  }
</script>

<div class="page-container">
  <!-- Header -->
  <div class="header-card">
    <button onclick={handleBack} class="back-button">
      <img src="/arrow-left.svg" alt="back" width="20" height="20" />
      Quay lại
    </button>

    <div class="header-content">
      <div class="header-info">
        <div class="header-badges">
          <h1 class="category-title">{category.name}</h1>
          <span class="status-badge {statusConfig.bg} {statusConfig.color}">
            {statusConfig.label}
          </span>
        </div>
        <p class="category-description">{category.description}</p>
        <div class="category-meta">
          <span>Lớp: {classData.code} - {classData.name}</span>
          <span>•</span>
          <span>Từ {category.startDate} đến {category.endDate}</span>
        </div>
      </div>
    </div>

    {#if !isRegistrationOpen}
      <div class="alert-box warning">
        <img src="/lock.svg" alt="lock" width="20" height="20" />
        <div>
          <div class="alert-title">Chưa mở đăng ký</div>
          <div class="alert-text">
            Giáo viên chưa mở đăng ký hoặc đã hết thời gian đăng ký cho hạng mục
            này
          </div>
        </div>
      </div>
    {/if}

    {#if hasRegistered && myRegisteredProject}
      <div class="alert-box success">
        <img src="/folder-kanban.svg" alt="project" width="20" height="20" />
        <div>
          <div class="alert-title">Đã đăng ký đề tài</div>
          <div class="alert-text">
            Bạn đã đăng ký đề tài: <strong>{myRegisteredProject.name}</strong>
          </div>
        </div>
      </div>
    {/if}
  </div>

  <!-- Tabs -->
  <div class="tabs-card">
    <div class="tabs-header">
      <button
        onclick={() => (activeTab = "projects")}
        class="tab-button"
        class:active={activeTab === "projects"}
      >
        Đề tài ({projects.length})
      </button>
      <button
        onclick={() => (activeTab = "reports")}
        class="tab-button"
        class:active={activeTab === "reports"}
      >
        Báo cáo ({reports.filter((r) => r.status === "pending").length} chờ nộp)
      </button>
    </div>

    <div class="tab-content">
      {#if activeTab === "projects"}
        <div class="projects-list">
          {#each projects as project (project.id)}
            {@const projectStatus = {
              available: {
                label: "Còn chỗ",
                color: "text-green-600",
                bg: "bg-green-50",
              },
              full: {
                label: "Đã đủ",
                color: "text-orange-600",
                bg: "bg-orange-50",
              },
              closed: {
                label: "Đã khóa",
                color: "text-red-600",
                bg: "bg-red-50",
              },
              forming: {
                label: "Đang hình thành",
                color: "text-blue-600",
                bg: "bg-blue-50",
              },
            }[project.status]}
            {@const isMyProject =
              hasRegistered && myRegisteredProject?.id === project.id}

            <div class="project-card" class:my-project={isMyProject}>
              <div class="project-header">
                <div class="project-info">
                  <div class="project-badges">
                    <h3 class="project-name">{project.name}</h3>
                    <span
                      class="status-badge {projectStatus.bg} {projectStatus.color}"
                    >
                      {projectStatus.label}
                    </span>
                    {#if isMyProject}
                      <span class="my-project-badge">Đề tài của tôi</span>
                    {/if}
                  </div>
                  <p class="project-description">{project.description}</p>

                  <div class="project-tags">
                    {#each project.tags as tag}
                      <span class="tag">
                        <img src="/tag.svg" alt="tag" width="12" height="12" />
                        {tag}
                      </span>
                    {/each}
                  </div>

                  <div class="project-meta">
                    <div class="meta-item">
                      <img
                        src="/users.svg"
                        alt="users"
                        width="16"
                        height="16"
                      />
                      <span
                        >{project.currentStudents}/{project.maxStudents} sinh viên</span
                      >
                    </div>
                    <span>•</span>
                    <span>GVHD: {project.instructor}</span>
                  </div>
                </div>

                <button
                  onclick={() => handleSelectProject(project)}
                  class="view-button"
                >
                  {isMyProject ? "Xem chi tiết" : "Xem & Đăng ký"}
                </button>
              </div>
            </div>
          {/each}
        </div>
      {:else if activeTab === "reports"}
        <div class="reports-tab">
          <div class="tab-header">
            <h2 class="tab-title">Danh sách báo cáo</h2>
            {#if !hasRegistered}
              <div class="warning-text">
                Bạn cần đăng ký đề tài để nộp báo cáo
              </div>
            {/if}
          </div>

          <div class="reports-list">
            {#each reports as report (report.id)}
              {@const reportStatus = {
                submitted: {
                  label: "Đã nộp",
                  color: "text-green-600",
                  bg: "bg-green-50",
                  icon: "/file-text.svg",
                },
                pending: {
                  label: "Chờ nộp",
                  color: "text-orange-600",
                  bg: "bg-orange-50",
                  icon: "/clock.svg",
                },
                upcoming: {
                  label: "Sắp tới",
                  color: "text-gray-600",
                  bg: "bg-gray-50",
                  icon: "/clock.svg",
                },
              }[report.status]}

              <div
                class="report-card"
                class:pending={report.status === "pending"}
              >
                <div class="report-content">
                  <img
                    src={reportStatus.icon}
                    alt="status"
                    width="20"
                    height="20"
                    class="report-icon {reportStatus.color}"
                  />
                  <div class="report-info">
                    <div class="report-header-row">
                      <h4 class="report-title">{report.title}</h4>
                      <span
                        class="status-badge {reportStatus.bg} {reportStatus.color}"
                      >
                        {reportStatus.label}
                      </span>
                    </div>
                    <div class="report-meta">
                      Hạn nộp: {report.deadline}
                    </div>
                    {#if report.submittedAt}
                      <div class="submitted-text">
                        Đã nộp: {report.submittedAt}
                      </div>
                    {/if}
                  </div>
                </div>

                {#if hasRegistered && report.status === "pending"}
                  <button class="submit-button">Nộp báo cáo</button>
                {:else if hasRegistered && report.status === "submitted"}
                  <button class="view-detail-button">Xem chi tiết</button>
                {/if}
              </div>
            {/each}
          </div>
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
    margin-bottom: 16px;
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

  .category-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1e293b;
  }

  .status-badge {
    padding: 4px 12px;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .bg-gray-50 {
    background: #f9fafb;
  }
  .text-gray-600 {
    color: #4b5563;
  }
  .bg-green-50 {
    background: #dcfce7;
  }
  .text-green-600 {
    color: #166534;
  }
  .bg-blue-50 {
    background: #dbeafe;
  }
  .text-blue-600 {
    color: #1e40af;
  }
  .bg-orange-50 {
    background: #fff7ed;
  }
  .text-orange-600 {
    color: #c2410c;
  }
  .bg-red-50 {
    background: #fef2f2;
  }
  .text-red-600 {
    color: #dc2626;
  }

  .category-description {
    color: #64748b;
    margin-bottom: 12px;
    line-height: 1.5;
  }

  .category-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .alert-box {
    display: flex;
    gap: 12px;
    padding: 16px;
    border-radius: 8px;
    margin-top: 16px;
  }

  .alert-box.warning {
    background: #fff7ed;
    border: 1px solid #fed7aa;
  }

  .alert-box.success {
    background: #eff6ff;
    border: 1px solid #bfdbfe;
  }

  .alert-title {
    font-weight: 500;
    margin-bottom: 4px;
  }

  .alert-box.warning .alert-title {
    color: #9a3412;
  }

  .alert-box.success .alert-title {
    color: #1e3a8a;
  }

  .alert-text {
    font-size: 0.875rem;
  }

  .alert-box.warning .alert-text {
    color: #c2410c;
  }

  .alert-box.success .alert-text {
    color: #3b82f6;
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
  }

  .tab-button:hover {
    color: #1e293b;
  }

  .tab-button.active {
    color: #2563eb;
    border-bottom-color: #2563eb;
  }

  .tab-content {
    padding: 24px;
  }

  .projects-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .project-card {
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 20px;
    transition: all 0.2s;
  }

  .project-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .project-card.my-project {
    background: #eff6ff;
    border-color: #93c5fd;
  }

  .project-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16px;
  }

  .project-info {
    flex: 1;
  }

  .project-badges {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .project-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1e293b;
  }

  .my-project-badge {
    padding: 4px 12px;
    background: #2563eb;
    color: white;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .project-description {
    color: #64748b;
    font-size: 0.875rem;
    margin-bottom: 12px;
    line-height: 1.5;
  }

  .project-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 12px;
  }

  .tag {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    background: #f1f5f9;
    color: #475569;
    border-radius: 4px;
    font-size: 0.75rem;
  }

  .project-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .view-button {
    padding: 10px 16px;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    white-space: nowrap;
    transition: background 0.2s;
  }

  .view-button:hover {
    background: #1d4ed8;
  }

  .tab-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .tab-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1e293b;
  }

  .warning-text {
    color: #64748b;
    font-size: 0.875rem;
  }

  .reports-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .report-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 16px;
  }

  .report-card.pending {
    border-color: #fed7aa;
  }

  .report-content {
    display: flex;
    gap: 12px;
    flex: 1;
  }

  .report-icon {
    flex-shrink: 0;
  }

  .report-info {
    flex: 1;
  }

  .report-header-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 4px;
  }

  .report-title {
    font-weight: 500;
    color: #1e293b;
  }

  .report-meta {
    color: #64748b;
    font-size: 0.875rem;
  }

  .submitted-text {
    color: #166534;
    font-size: 0.875rem;
    margin-top: 4px;
  }

  .submit-button {
    padding: 8px 16px;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    white-space: nowrap;
    transition: background 0.2s;
  }

  .submit-button:hover {
    background: #1d4ed8;
  }

  .view-detail-button {
    padding: 8px 16px;
    background: #e5e7eb;
    color: #374151;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    white-space: nowrap;
    transition: background 0.2s;
  }

  .view-detail-button:hover {
    background: #d1d5db;
  }

  @media (max-width: 768px) {
    .page-container {
      padding: 16px;
    }

    .project-header {
      flex-direction: column;
    }

    .view-button {
      width: 100%;
    }

    .report-card {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .submit-button,
    .view-detail-button {
      width: 100%;
    }
  }
</style>
