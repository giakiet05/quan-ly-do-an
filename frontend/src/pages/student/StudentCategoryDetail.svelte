<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { getClassroom } from "../../services/classroom-service";
  import { getProjectRounds } from "../../services/project-round-service";
  import { getProjects } from "../../services/project-service";
  import { getGroupsFilter } from "../../services/group-service";
  import { authStore } from "../../stores/auth-store";
  import type { ProjectRound, Group } from "../../models";
  import type { ProjectResponse } from "../../dtos/project-dto";

  let { params } = $props<{
    params: { id: string; categoryId: string };
  }>();

  let activeTab = $state<"projects" | "reports">("projects");
  let loading = $state(true);
  let error = $state<string | null>(null);

  let projectRound = $state<ProjectRound | null>(null);
  let projects = $state<ProjectResponse[]>([]);
  let groups = $state<Group[] | null>(null);
  let myGroup = $derived(
    groups?.find((g) => g.members.some((m) => m.id === $authStore.user?.id)),
  );

  // Mock data - Replace with API calls
  let category = $state({
    id: params.categoryId,
    name: "Đồ án chuyên ngành",
    description: "Phát triển ứng dụng hoàn chỉnh với đầy đủ tính năng",
    startDate: "2024-01-15",
    endDate: "2024-05-30",
    status: "ongoing",
  });

  // ...

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

  onMount(async () => {
    try {
      loading = true;
      error = null;

      // TEMPORARY: Use mock data if classroom ID starts with "mock-"
      if (params.id.startsWith("mock-")) {
        await new Promise((resolve) => setTimeout(resolve, 500));

        const mockClassroom = {
          id: params.id,
          projectRounds:
            params.id === "mock-class-1"
              ? [
                  {
                    id: "round1",
                    name: "Đợt 1 - Đồ án cuối kỳ",
                    description: "Phát triển ứng dụng web hoàn chỉnh",
                    startDate: "2024-02-01T00:00:00Z",
                    endDate: "2024-05-31T23:59:59Z",
                    reportPeriods: [
                      {
                        id: "rp1",
                        title: "Báo cáo đề cương",
                        description: "Nộp báo cáo đề cương dự án",
                        fileType: ["pdf", "docx"],
                        startDate: "2024-02-01T00:00:00Z",
                        endDate: "2024-02-15T23:59:59Z",
                      },
                    ],
                    createdAt: "2024-01-15T00:00:00Z",
                    isDeleted: false,
                  },
                ]
              : [
                  {
                    id: "round2",
                    name: "Đợt 1 - AI Research",
                    description: "Nghiên cứu và ứng dụng AI",
                    startDate: "2024-02-01T00:00:00Z",
                    endDate: "2024-05-31T23:59:59Z",
                    reportPeriods: [],
                    createdAt: "2024-01-15T00:00:00Z",
                    isDeleted: false,
                  },
                ],
        };

        const foundRound = mockClassroom.projectRounds.find(
          (round) => round.id === params.categoryId,
        );

        if (!foundRound) {
          error = "Không tìm thấy đợt đồ án này";
          return;
        }

        projectRound = foundRound as unknown as ProjectRound;
        groups = []; // Empty groups for mock data
      } else {
        // Fetch classroom info first
        const classroom = await getClassroom(params.id);
        console.log("📚 Classroom data:", classroom);

        // Fetch project rounds separately
        const projectRounds = await getProjectRounds(params.id);
        console.log("📋 Project rounds:", projectRounds);

        const foundRound = projectRounds.find(
          (round) => round.id === params.categoryId,
        );

        if (!foundRound) {
          error = "Không tìm thấy đợt đồ án này";
          return;
        }

        projectRound = foundRound;

        // Fetch projects in this round
        projects = await getProjects(params.id, params.categoryId);
        console.log("📁 Projects in this round:", projects);

        // Try to fetch groups (may fail if backend not fixed)
        try {
          const groupsData = await getGroupsFilter({ classroom_id: params.id });
          groups = groupsData || [];
          console.log("👥 Groups in this classroom:", groups);
        } catch (groupErr) {
          console.warn("⚠️ Could not load groups (backend issue):", groupErr);
          groups = []; // Continue without groups
        }
      }
    } catch (err) {
      console.error("Error loading project round:", err);
      error = err instanceof Error ? err.message : "Không thể tải dữ liệu";
    } finally {
      loading = false;
    }
  });

  function handleBack() {
    push(`/classes/${params.id}`);
  }

  function handleProjectClick(projectId: string) {
    push(
      `/student/classes/${params.id}/categories/${params.categoryId}/projects/${projectId}`,
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
  {:else if error}
    <div class="error-message">
      <p>{error}</p>
      <button onclick={handleBack} class="back-button">Quay lại</button>
    </div>
  {:else if projectRound}
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
          <h1 class="category-title">{projectRound.name}</h1>
          <p class="category-description">{projectRound.description}</p>
          <p class="category-dates">
            Thời gian: {formatDate(projectRound.startDate)} - {formatDate(
              projectRound.endDate,
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
        {@const myProjectId = myGroup?.project_id}
        <div class="projects-tab">
          {#if myProjectId}
            {@const myProject = projects.find((p) => p.id === myProjectId)}
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
                <p>Đề tài: {myProject?.title || "N/A"}</p>
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

          {#if projects.length === 0}
            <div class="empty-state">
              <svg
                width="64"
                height="64"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                ></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
                <polyline points="10 9 9 9 8 9"></polyline>
              </svg>
              <h3>Chưa có đề tài nào</h3>
              <p>
                Giảng viên chưa tạo đề tài cho đợt này. Vui lòng quay lại sau.
              </p>
            </div>
          {:else}
            <div class="projects-list">
              {#each projects as project}
                {@const projectGroups = groups.filter(
                  (g) => g.project_id === project.id,
                )}
                {@const currentMembers = projectGroups.reduce(
                  (sum, g) => sum + g.members.length,
                  0,
                )}
                {@const maxMembers = project.amount * project.max_member}
                {@const isMyProject = myProjectId === project.id}
                {@const isFull = currentMembers >= maxMembers}

                <div
                  class="project-card"
                  onclick={() => handleProjectClick(project.id)}
                >
                  <div class="project-header">
                    <h3 class="project-name">{project.title}</h3>
                    <div class="project-badges">
                      {#if isMyProject}
                        <span class="badge badge-my-project">
                          Đề tài của tôi
                        </span>
                      {/if}
                      {#if project.status === "approved"}
                        <span class="badge badge-available">
                          {isFull ? "Đã đủ" : "Còn chỗ"}
                        </span>
                      {:else if project.status === "ongoing"}
                        <span class="badge badge-ongoing">Đang thực hiện</span>
                      {:else if project.status === "completed"}
                        <span class="badge badge-closed">Đã hoàn thành</span>
                      {:else}
                        <span class="badge badge-pending">Chờ duyệt</span>
                      {/if}
                    </div>
                  </div>

                  <p class="project-description">{project.description}</p>

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
                      <span>{currentMembers}/{maxMembers} sinh viên</span>
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
                        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"
                        ></path>
                        <circle cx="9" cy="7" r="4"></circle>
                      </svg>
                      <span>Số lượng nhóm: {project.amount}</span>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "reports"}
        <div class="reports-tab">
          {#if myGroup}
            <div class="reports-timeline">
              {#each projectRound.reportPeriods as reportPeriod}
                {@const myReport = myGroup.reports.find(
                  (r) => r.reportPeriodId === reportPeriod.id,
                )}
                {@const daysUntil = getDaysUntil(reportPeriod.endDate)}
                {@const isSubmitted = !!myReport}

                <div class="timeline-item" class:completed={isSubmitted}>
                  <div class="timeline-marker">
                    {#if isSubmitted}
                      <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                      >
                        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                        <polyline points="22 4 12 14.01 9 11.01"></polyline>
                      </svg>
                    {:else if daysUntil >= 0 && daysUntil <= 7}
                      <div class="marker-pending"></div>
                    {:else}
                      <div class="marker-upcoming"></div>
                    {/if}
                  </div>

                  <div class="timeline-content">
                    <div class="report-header">
                      <h3 class="report-title">{reportPeriod.title}</h3>
                      {#if isSubmitted}
                        <span class="status-badge status-submitted">
                          Đã nộp
                        </span>
                      {:else if daysUntil >= 0}
                        <span class="status-badge status-pending">
                          Chưa nộp
                        </span>
                      {:else}
                        <span class="status-badge status-overdue">
                          Quá hạn
                        </span>
                      {/if}
                    </div>

                    <p class="report-description">{reportPeriod.description}</p>

                    <div class="report-info">
                      <div class="info-item">
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
                        <span>Hạn nộp: {formatDate(reportPeriod.endDate)}</span>
                      </div>

                      {#if isSubmitted && myReport}
                        <div class="info-item success">
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
                          <span>Đã nộp: {formatDate(myReport.createdAt)}</span>
                        </div>
                      {:else if daysUntil >= 0}
                        <div class="info-item warning">
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
                          <span>Còn {daysUntil} ngày</span>
                        </div>
                      {:else}
                        <div class="info-item error">
                          <span>Đã quá hạn {Math.abs(daysUntil)} ngày</span>
                        </div>
                      {/if}
                    </div>

                    {#if isSubmitted && myReport?.feedback}
                      <div class="feedback-box">
                        <h4>Nhận xét từ giảng viên:</h4>
                        <p>{myReport.feedback.content}</p>
                        {#if myReport.feedback.grade}
                          <p class="grade">Điểm: {myReport.feedback.grade}</p>
                        {/if}
                      </div>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <div class="alert alert-info">
              <p>Bạn chưa tham gia nhóm nào trong hạng mục này</p>
            </div>
          {/if}
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

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
    background: white;
    border-radius: 0.75rem;
  }

  .empty-state svg {
    color: #d1d5db;
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #111827;
    margin-bottom: 0.5rem;
  }

  .empty-state p {
    font-size: 0.875rem;
    color: #6b7280;
  }
</style>
