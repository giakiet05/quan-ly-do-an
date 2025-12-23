<script lang="ts">
  import { push } from "svelte-spa-router";

  interface Project {
    id: string;
    name: string;
    description: string;
    instructor: string;
    status: "available" | "full" | "closed";
    currentStudents: number;
    maxStudents: number;
    tags: string[];
    className: string;
    categoryName: string;
    categoryStatus: "upcoming" | "ongoing" | "completed";
    allowStudentEdit: boolean;
  }

  // Mock data - sẽ thay bằng API call
  let myProjects = $state<Project[]>([
    {
      id: "p1",
      name: "Hệ thống quản lý thư viện trực tuyến",
      description:
        "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
      instructor: "TS. Nguyễn Văn A",
      status: "available",
      currentStudents: 2,
      maxStudents: 3,
      tags: ["Web", "React", "Node.js", "MongoDB"],
      className: "Công nghệ phần mềm - K18",
      categoryName: "Đồ án chuyên ngành - HK1 2024-2025",
      categoryStatus: "ongoing",
      allowStudentEdit: true,
    },
    {
      id: "p2",
      name: "Ứng dụng di động quản lý chi tiêu cá nhân",
      description:
        "Ứng dụng mobile giúp người dùng theo dõi thu chi, lập kế hoạch tài chính",
      instructor: "ThS. Trần Thị B",
      status: "full",
      currentStudents: 2,
      maxStudents: 2,
      tags: ["Mobile", "React Native", "Firebase"],
      className: "Lập trình di động - K18",
      categoryName: "Đồ án chuyên ngành - HK1 2024-2025",
      categoryStatus: "ongoing",
      allowStudentEdit: false,
    },
  ]);

  function handleSelectProject(project: Project) {
    push(`/my-projects/${project.id}`);
  }

  function goToClasses() {
    push("/classes");
  }
</script>

<div class="projects-container">
  <div class="header">
    <h1 class="title">Đề tài của tôi</h1>
    <p class="subtitle">Danh sách các đề tài bạn đã đăng ký</p>
  </div>

  {#if myProjects.length === 0}
    <div class="empty-state">
      <div class="empty-icon">
        <svg
          width="64"
          height="64"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
          <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
        </svg>
      </div>
      <h3 class="empty-title">Chưa có đề tài nào</h3>
      <p class="empty-text">Bạn chưa đăng ký đề tài nào</p>
      <button onclick={goToClasses} class="btn-primary">
        Xem danh sách lớp học
      </button>
    </div>
  {:else}
    <div class="projects-grid">
      {#each myProjects as project (project.id)}
        {@const statusConfig = {
          available: {
            label: "Còn chỗ",
            color: "status-available",
          },
          full: {
            label: "Đã đủ",
            color: "status-full",
          },
          closed: {
            label: "Đã khóa",
            color: "status-closed",
          },
        }[project.status]}
        <button
          onclick={() => handleSelectProject(project)}
          class="project-card"
        >
          <div class="project-header">
            <div class="project-main">
              <div class="project-title-row">
                <h3 class="project-title">{project.name}</h3>
                <span class="status-badge {statusConfig.color}">
                  {statusConfig.label}
                </span>
              </div>
              <p class="project-description">{project.description}</p>
            </div>
          </div>

          <div class="project-tags">
            {#each project.tags as tag}
              <span class="tag">
                <svg
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"
                  ></path>
                  <line x1="7" y1="7" x2="7.01" y2="7"></line>
                </svg>
                {tag}
              </span>
            {/each}
          </div>

          <div class="project-info">
            <div class="info-item">
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
              <span>{project.className}</span>
            </div>
            <span class="separator">•</span>
            <div class="info-item">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
              <span>{project.categoryName}</span>
            </div>
          </div>

          <div class="project-meta">
            <div class="info-item">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
              </svg>
              <span
                >{project.currentStudents}/{project.maxStudents} sinh viên</span
              >
            </div>
            <span class="separator">•</span>
            <span>GVHD: {project.instructor}</span>
            {#if project.allowStudentEdit}
              <span class="separator">•</span>
              <span class="editable">Được phép chỉnh sửa</span>
            {/if}
          </div>
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .projects-container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .header {
    margin-bottom: 24px;
  }

  .title {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .subtitle {
    color: #64748b;
  }

  /* Empty State */
  .empty-state {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 48px;
    text-align: center;
  }

  .empty-icon {
    width: 64px;
    height: 64px;
    background: #f3f4f6;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    color: #d1d5db;
  }

  .empty-title {
    font-size: 18px;
    color: #111827;
    margin-bottom: 8px;
  }

  .empty-text {
    color: #6b7280;
    margin-bottom: 16px;
  }

  .btn-primary {
    padding: 8px 16px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-primary:hover {
    background: #2563eb;
  }

  /* Projects Grid */
  .projects-grid {
    display: grid;
    gap: 16px;
  }

  .project-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 24px;
    border: 1px solid #e5e7eb;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    width: 100%;
  }

  .project-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    border-color: #93c5fd;
  }

  .project-header {
    margin-bottom: 16px;
  }

  .project-main {
    flex: 1;
  }

  .project-title-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .project-title {
    font-size: 20px;
    font-weight: 600;
    color: #1e293b;
  }

  .status-badge {
    padding: 4px 12px;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 500;
  }

  .status-available {
    background: #dcfce7;
    color: #16a34a;
  }

  .status-full {
    background: #fed7aa;
    color: #ea580c;
  }

  .status-closed {
    background: #fee2e2;
    color: #dc2626;
  }

  .project-description {
    color: #64748b;
    font-size: 14px;
    margin-bottom: 12px;
  }

  /* Tags */
  .project-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
  }

  .tag {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    background: #f3f4f6;
    color: #374151;
    border-radius: 4px;
    font-size: 12px;
  }

  /* Info */
  .project-info {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 12px;
    font-size: 14px;
    color: #64748b;
    margin-bottom: 12px;
  }

  .project-meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 12px;
    font-size: 14px;
    color: #64748b;
  }

  .info-item {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .separator {
    color: #d1d5db;
  }

  .editable {
    color: #16a34a;
  }
</style>
