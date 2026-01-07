<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";

  let { params } = $props<{
    params: { id: string; categoryId: string };
  }>();

  let activeTab = $state<"projects" | "reports">("projects");
  let loading = $state(true);

  // Mock data - Replace with API calls
  let category = $state({
    id: params.categoryId,
    name: "Đồ án chuyên ngành",
    description: "Phát triển ứng dụng hoàn chỉnh với đầy đủ tính năng",
    startDate: "2024-01-15",
    endDate: "2024-05-30",
    status: "ongoing",
  });

  let projects = $state([
    {
      id: "p1",
      name: "Hệ thống quản lý thư viện",
      description:
        "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
      instructor: "TS. Nguyễn Văn A",
      currentStudents: 2,
      maxStudents: 3,
      status: "available",
      tags: ["Web", "React", "Node.js"],
      isMyProject: false,
    },
    {
      id: "p2",
      name: "Ứng dụng quản lý chi tiêu",
      description: "Ứng dụng mobile giúp theo dõi thu chi cá nhân",
      instructor: "ThS. Trần Thị B",
      currentStudents: 2,
      maxStudents: 2,
      status: "full",
      tags: ["Mobile", "React Native"],
      isMyProject: true,
    },
  ]);

  let reports = $state([
    {
      id: "r1",
      title: "Báo cáo đề cương",
      deadline: "2024-02-15",
      status: "submitted",
      submittedDate: "2024-02-14",
    },
    {
      id: "r2",
      title: "Báo cáo tiến độ giữa kỳ",
      deadline: "2024-03-30",
      status: "pending",
    },
    {
      id: "r3",
      title: "Báo cáo cuối kỳ",
      deadline: "2024-05-25",
      status: "upcoming",
    },
  ]);

  onMount(() => {
    loading = false;
  });

  function handleBack() {
    push(`/classes/${params.id}`);
  }

  function handleProjectClick(projectId: string) {
    push(
      `/classes/${params.id}/categories/${params.categoryId}/projects/${projectId}`
    );
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString("vi-VN");
  }

  function getDaysUntil(dateString: string): number {
    const deadline = new Date(dateString);
    const today = new Date();
    const diff = deadline.getTime() - today.getTime();
    return Math.ceil(diff / (1000 * 60 * 60 * 24));
  }
</script>

<div class="container">
  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Đang tải...</p>
    </div>
  {:else}
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
        Quay lại chi tiết lớp
      </button>

      <div class="header-content">
        <div>
          <h1 class="category-title">{category.name}</h1>
          <p class="category-description">{category.description}</p>
          <p class="category-dates">
            Thời gian: {formatDate(category.startDate)} - {formatDate(
              category.endDate
            )}
          </p>
        </div>
      </div>

      <div class="tabs">
        <button
          onclick={() => (activeTab = "projects")}
          class="tab"
          class:active={activeTab === "projects"}
        >
          Đề tài
        </button>
        <button
          onclick={() => (activeTab = "reports")}
          class="tab"
          class:active={activeTab === "reports"}
        >
          Báo cáo
        </button>
      </div>
    </div>

    <div class="content-section">
      {#if activeTab === "projects"}
        <div class="projects-tab">
          {#if projects.some((p) => p.isMyProject)}
            <div class="alert alert-success">
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
              <div>
                <strong>Bạn đã đăng ký đề tài trong hạng mục này</strong>
                <p>
                  Đề tài: {projects.find((p) => p.isMyProject)?.name}
                </p>
              </div>
            </div>
          {:else}
            <div class="alert alert-warning">
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"
                ></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
              <div>
                <strong>Bạn chưa đăng ký đề tài</strong>
                <p>Vui lòng chọn và đăng ký đề tài phù hợp với nhóm của bạn</p>
              </div>
            </div>
          {/if}

          <div class="projects-list">
            {#each projects as project}
              <div
                class="project-card"
                onclick={() => handleProjectClick(project.id)}
              >
                <div class="project-header">
                  <h3 class="project-name">{project.name}</h3>
                  <div class="project-badges">
                    {#if project.isMyProject}
                      <span class="badge badge-my-project">
                        Đề tài của tôi
                      </span>
                    {/if}
                    {#if project.status === "available"}
                      <span class="badge badge-available">Còn chỗ</span>
                    {:else if project.status === "full"}
                      <span class="badge badge-full">Đã đủ</span>
                    {:else}
                      <span class="badge badge-closed">Đã khóa</span>
                    {/if}
                  </div>
                </div>

                <p class="project-description">{project.description}</p>

                <div class="project-tags">
                  {#each project.tags as tag}
                    <span class="tag">{tag}</span>
                  {/each}
                </div>

                <div class="project-footer">
                  <div class="footer-item">
                    <svg
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"
                      ></path>
                      <circle cx="9" cy="7" r="4"></circle>
                      <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                      <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                    </svg>
                    <span
                      >{project.currentStudents}/{project.maxStudents} sinh viên</span
                    >
                  </div>
                  <div class="footer-item">
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
                    <span>GVHD: {project.instructor}</span>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else if activeTab === "reports"}
        <div class="reports-tab">
          <div class="reports-timeline">
            {#each reports as report, index}
              {@const daysUntil = getDaysUntil(report.deadline)}
              <div
                class="timeline-item"
                class:completed={report.status === "submitted"}
              >
                <div class="timeline-marker">
                  {#if report.status === "submitted"}
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 24 24"
                      fill="currentColor"
                    >
                      <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                      <polyline points="22 4 12 14.01 9 11.01"></polyline>
                    </svg>
                  {:else if report.status === "pending"}
                    <div class="marker-pending"></div>
                  {:else}
                    <div class="marker-upcoming"></div>
                  {/if}
                </div>

                <div class="timeline-content">
                  <div class="report-header">
                    <h3 class="report-title">{report.title}</h3>
                    {#if report.status === "submitted"}
                      <span class="status-badge status-submitted">
                        Đã nộp
                      </span>
                    {:else if report.status === "pending"}
                      <span class="status-badge status-pending">
                        Đang chờ nộp
                      </span>
                    {:else}
                      <span class="status-badge status-upcoming">
                        Sắp tới
                      </span>
                    {/if}
                  </div>

                  <div class="report-details">
                    <div class="detail-item">
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
                      <span>Hạn nộp: {formatDate(report.deadline)}</span>
                    </div>

                    {#if report.status === "submitted" && report.submittedDate}
                      <div class="detail-item">
                        <svg
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                          <polyline points="22 4 12 14.01 9 11.01"></polyline>
                        </svg>
                        <span>Đã nộp: {formatDate(report.submittedDate)}</span>
                      </div>
                    {:else if report.status === "pending"}
                      <div class="detail-item" class:urgent={daysUntil <= 3}>
                        <svg
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <circle cx="12" cy="12" r="10"></circle>
                          <polyline points="12 6 12 12 16 14"></polyline>
                        </svg>
                        <span>
                          {#if daysUntil > 0}
                            Còn {daysUntil} ngày
                          {:else if daysUntil === 0}
                            Hết hạn hôm nay
                          {:else}
                            Đã quá hạn {Math.abs(daysUntil)} ngày
                          {/if}
                        </span>
                      </div>
                    {/if}
                  </div>

                  {#if report.status === "submitted"}
                    <button class="btn-view-report">Xem báo cáo đã nộp</button>
                  {:else if report.status === "pending"}
                    <button class="btn-submit-report">Nộp báo cáo</button>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}
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

  .category-title {
    font-size: 1.875rem;
    font-weight: 700;
    color: #111827;
    margin-bottom: 0.5rem;
  }

  .category-description {
    font-size: 1rem;
    color: #6b7280;
    margin-bottom: 0.5rem;
  }

  .category-dates {
    font-size: 0.875rem;
    color: #9ca3af;
  }

  .tabs {
    display: flex;
    gap: 0.5rem;
    border-bottom: 2px solid #e5e7eb;
    margin-top: 1.5rem;
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

  .content-section {
    background: white;
    border-radius: 0.75rem;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .alert {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .alert svg {
    flex-shrink: 0;
  }

  .alert-success {
    background: #d1fae5;
    color: #065f46;
    border: 1px solid #10b981;
  }

  .alert-warning {
    background: #fef3c7;
    color: #92400e;
    border: 1px solid #f59e0b;
  }

  .alert strong {
    display: block;
    margin-bottom: 0.25rem;
  }

  .alert p {
    font-size: 0.875rem;
  }

  .projects-list {
    display: grid;
    gap: 1rem;
  }

  .project-card {
    padding: 1.25rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .project-card:hover {
    border-color: #3b82f6;
    box-shadow: 0 4px 6px rgba(59, 130, 246, 0.1);
  }

  .project-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.75rem;
  }

  .project-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
  }

  .project-badges {
    display: flex;
    gap: 0.5rem;
  }

  .badge {
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: 9999px;
  }

  .badge-my-project {
    background: #dcfce7;
    color: #166534;
  }

  .badge-available {
    background: #dbeafe;
    color: #1e40af;
  }

  .badge-full {
    background: #fef3c7;
    color: #92400e;
  }

  .badge-closed {
    background: #f3f4f6;
    color: #6b7280;
  }

  .project-description {
    font-size: 0.875rem;
    color: #6b7280;
    margin-bottom: 1rem;
  }

  .project-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .tag {
    padding: 0.25rem 0.75rem;
    background: #f3f4f6;
    color: #4b5563;
    font-size: 0.75rem;
    border-radius: 0.25rem;
  }

  .project-footer {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .footer-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .reports-timeline {
    position: relative;
  }

  .timeline-item {
    position: relative;
    padding-left: 3rem;
    padding-bottom: 2rem;
  }

  .timeline-item:not(:last-child)::before {
    content: "";
    position: absolute;
    left: 0.625rem;
    top: 2rem;
    bottom: 0;
    width: 2px;
    background: #e5e7eb;
  }

  .timeline-item.completed::before {
    background: #10b981;
  }

  .timeline-marker {
    position: absolute;
    left: 0;
    top: 0.25rem;
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    border-radius: 50%;
  }

  .timeline-marker svg {
    color: #10b981;
  }

  .marker-pending,
  .marker-upcoming {
    width: 1rem;
    height: 1rem;
    border-radius: 50%;
    border: 2px solid;
  }

  .marker-pending {
    border-color: #f59e0b;
    background: #fef3c7;
  }

  .marker-upcoming {
    border-color: #9ca3af;
    background: #f3f4f6;
  }

  .timeline-content {
    padding: 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
    background: white;
  }

  .report-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.75rem;
  }

  .report-title {
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

  .status-submitted {
    background: #dcfce7;
    color: #166534;
  }

  .status-pending {
    background: #fef3c7;
    color: #92400e;
  }

  .status-upcoming {
    background: #f3f4f6;
    color: #6b7280;
  }

  .report-details {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .detail-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .detail-item.urgent {
    color: #dc2626;
    font-weight: 600;
  }

  .btn-view-report,
  .btn-submit-report {
    padding: 0.625rem 1.25rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-view-report {
    background: #f3f4f6;
    color: #374151;
  }

  .btn-view-report:hover {
    background: #e5e7eb;
  }

  .btn-submit-report {
    background: #3b82f6;
    color: white;
  }

  .btn-submit-report:hover {
    background: #2563eb;
  }
</style>
