<!-- 
  EXAMPLE: Tích hợp Backend cho My Projects Page
  
  So sánh Mock Data vs API Integration
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import type { Group, Project } from "../models";
  import { getClassroomGroups } from "../services/group-service";
  import { getClassroomProjects } from "../services/project-service";

  // ==================== STATE ====================
  let projects = $state<
    Array<{
      project: Project;
      group?: Group;
    }>
  >([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  // Assume we have classroomId from somewhere (route params or context)
  const classroomId = "classroom-id-here";

  // ==================== API INTEGRATION ====================
  onMount(async () => {
    try {
      loading = true;

      // Fetch projects and groups in parallel
      const [allProjects, myGroups] = await Promise.all([
        getClassroomProjects(classroomId),
        getClassroomGroups(classroomId),
      ]);

      // Map projects with their corresponding groups
      projects = allProjects.map((project) => ({
        project,
        group: myGroups.find((g) => g.projectId === project.id),
      }));
    } catch (err: any) {
      error = err.message || "Không thể tải dự án";
      console.error("Error loading projects:", err);
    } finally {
      loading = false;
    }
  });

  // ==================== HANDLERS ====================
  function handleProjectClick(projectId: string) {
    push(`/my-projects/${projectId}`);
  }

  // ==================== HELPERS ====================
  function getStatusBadge(status: Project["status"]) {
    const badges = {
      pending: { text: "Chờ duyệt", class: "badge-warning" },
      approved: { text: "Đã duyệt", class: "badge-success" },
      ongoing: { text: "Đang thực hiện", class: "badge-info" },
      completed: { text: "Hoàn thành", class: "badge-completed" },
    };
    return badges[status] || badges.pending;
  }

  function getAvailableSlots(project: Project, group?: Group) {
    if (!group) return project.maxMember;
    return project.maxMember - (group.members.length + 1); // +1 for leader
  }
</script>

<!-- ==================== TEMPLATE ==================== -->
<div class="my-projects-page">
  <h1 class="page-title">Đề tài của tôi</h1>

  <!-- LOADING STATE -->
  {#if loading}
    <div class="loading-container">
      <div class="spinner"></div>
      <p>Đang tải dự án...</p>
    </div>

    <!-- ERROR STATE -->
  {:else if error}
    <div class="error-container">
      <svg
        width="48"
        height="48"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
      >
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <p class="error-text">{error}</p>
      <button onclick={() => window.location.reload()} class="retry-button">
        Thử lại
      </button>
    </div>

    <!-- EMPTY STATE -->
  {:else if projects.length === 0}
    <div class="empty-container">
      <svg
        width="64"
        height="64"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
      >
        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
        <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
      </svg>
      <p class="empty-text">Bạn chưa tham gia đề tài nào</p>
      <button onclick={() => push("/classes")} class="primary-button">
        Xem danh sách lớp học
      </button>
    </div>

    <!-- PROJECT LIST -->
  {:else}
    <div class="projects-grid">
      {#each projects as { project, group } (project.id)}
        {@const badge = getStatusBadge(project.status)}
        {@const availableSlots = getAvailableSlots(project, group)}

        <div
          class="project-card"
          onclick={() => handleProjectClick(project.id)}
        >
          <div class="card-header">
            <h3 class="project-title">{project.title}</h3>
            <span class="status-badge {badge.class}">{badge.text}</span>
          </div>

          <p class="project-description">{project.description}</p>

          <div class="project-meta">
            <div class="meta-item">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
              </svg>
              <span>
                {#if group}
                  Nhóm {group.members.length + 1}/{project.maxMember}
                {:else}
                  Chưa có nhóm
                {/if}
              </span>
            </div>

            {#if availableSlots > 0}
              <div class="meta-item success">
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
                <span>Còn {availableSlots} chỗ</span>
              </div>
            {:else}
              <div class="meta-item warning">
                <span>Đã đủ thành viên</span>
              </div>
            {/if}
          </div>

          {#if group}
            <div class="group-badge">
              Bạn là {group.leaderId === "current-user-id"
                ? "Nhóm trưởng"
                : "Thành viên"}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .my-projects-page {
    padding: 24px;
    max-width: 1200px;
    margin: 0 auto;
  }

  .page-title {
    font-size: 28px;
    font-weight: 600;
    margin-bottom: 24px;
    color: #1e293b;
  }

  /* Loading State */
  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 20px;
    color: #64748b;
  }

  .spinner {
    width: 48px;
    height: 48px;
    border: 4px solid #e2e8f0;
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-bottom: 16px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  /* Error State */
  .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 80px 20px;
    color: #ef4444;
  }

  .error-text {
    margin: 16px 0;
    font-size: 16px;
  }

  .retry-button {
    padding: 8px 16px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }

  /* Empty State */
  .empty-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 80px 20px;
    color: #64748b;
  }

  .empty-text {
    margin: 16px 0;
    font-size: 16px;
  }

  .primary-button {
    padding: 10px 20px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
  }

  /* Projects Grid */
  .projects-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 24px;
  }

  .project-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 20px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .project-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }

  .card-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 12px;
  }

  .project-title {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
    flex: 1;
  }

  .status-badge {
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
  }

  .badge-warning {
    background: #fef3c7;
    color: #92400e;
  }
  .badge-success {
    background: #dcfce7;
    color: #166534;
  }
  .badge-info {
    background: #dbeafe;
    color: #1e40af;
  }
  .badge-completed {
    background: #f3f4f6;
    color: #374151;
  }

  .project-description {
    color: #64748b;
    font-size: 14px;
    margin-bottom: 16px;
    line-height: 1.5;
  }

  .project-meta {
    display: flex;
    gap: 16px;
    margin-bottom: 12px;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    color: #64748b;
  }

  .meta-item.success {
    color: #22c55e;
  }

  .meta-item.warning {
    color: #f59e0b;
  }

  .group-badge {
    padding: 8px 12px;
    background: #eff6ff;
    color: #1e40af;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    text-align: center;
  }
</style>
