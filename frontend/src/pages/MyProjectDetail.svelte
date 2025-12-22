<script lang="ts">
  import { push } from "svelte-spa-router";
  import StudentProjectReportsTab from "../components/StudentProjectReportsTab.svelte";
  import StudentProjectTeamChatTab from "../components/StudentProjectTeamChatTab.svelte";
  import StudentProjectMembersTab from "../components/StudentProjectMembersTab.svelte";

  interface TeamMember {
    id: string;
    name: string;
    studentId: string;
    email?: string;
    joinedAt?: string;
  }

  interface Team {
    id: string;
    name: string;
    leaderId: string;
    leaderName: string;
    members: TeamMember[];
    status?: "registered" | "ready" | "recruiting";
    createdAt?: string;
  }

  interface Project {
    id: string;
    name: string;
    description: string;
    instructor: string;
    className?: string;
    categoryName?: string;
  }

  let { params } = $props();
  let activeTab = $state<"reports" | "team" | "members">("reports");

  // Mock data
  const currentStudentId = "s1";
  const project: Project = {
    id: params.id,
    name: "Hệ thống quản lý thư viện trực tuyến",
    description:
      "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
    instructor: "TS. Nguyễn Văn A",
    className: "Công nghệ phần mềm - K18",
    categoryName: "Đồ án chuyên ngành - HK1 2024-2025",
  };

  const team: Team = {
    id: "t1",
    name: "Nhóm 1",
    leaderId: "s1",
    leaderName: "Nguyễn Văn An",
    status: "registered",
    createdAt: "2024-02-15T08:00:00Z",
    members: [
      {
        id: "s2",
        name: "Trần Văn B",
        studentId: "20200002",
        email: "tranvanb@student.edu.vn",
        joinedAt: "2024-02-16T10:30:00Z",
      },
      {
        id: "s3",
        name: "Lê Thị C",
        studentId: "20200003",
        email: "lethic@student.edu.vn",
        joinedAt: "2024-02-17T14:20:00Z",
      },
    ],
  };

  const isLeader = $derived(team.leaderId === currentStudentId);

  function handleBack() {
    push("/my-projects");
  }
</script>

<div class="detail-container">
  <!-- Header -->
  <div class="header-card">
    <button onclick={handleBack} class="back-button">
      <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <line x1="19" y1="12" x2="5" y2="12"></line>
        <polyline points="12 19 5 12 12 5"></polyline>
      </svg>
      Quay lại
    </button>

    <div class="header-content">
      <div class="header-main">
        <h1 class="project-title">{project.name}</h1>
        <p class="project-description">{project.description}</p>

        <div class="project-info">
          {#if project.className}
            <span>{project.className}</span>
            <span class="separator">•</span>
          {/if}
          {#if project.categoryName}
            <span>{project.categoryName}</span>
            <span class="separator">•</span>
          {/if}
          <span>GVHD: {project.instructor}</span>
        </div>

        <div class="badges">
          <span class="badge badge-team">Nhóm: {team.name}</span>
          {#if isLeader}
            <span class="badge badge-leader">Nhóm trưởng</span>
          {/if}
        </div>
      </div>
    </div>
  </div>

  <!-- Tabs -->
  <div class="tabs-card">
    <div class="tabs-header">
      <button
        onclick={() => (activeTab = "reports")}
        class="tab-button"
        class:active={activeTab === "reports"}
      >
        <img src="/myproject_report.svg" alt="Báo cáo" width="16" height="16" />
        Báo cáo
      </button>
      <button
        onclick={() => (activeTab = "team")}
        class="tab-button"
        class:active={activeTab === "team"}
      >
        <img src="/myproject_chat.svg" alt="Nhóm" width="16" height="16" />
        Nhóm
      </button>
      <button
        onclick={() => (activeTab = "members")}
        class="tab-button"
        class:active={activeTab === "members"}
      >
        <img src="/users.svg" alt="Thành viên" width="16" height="16" />
        Thành viên ({team.members.length + 1})
      </button>
    </div>

    <div class="tab-content">
      {#if activeTab === "reports"}
        <StudentProjectReportsTab
          projectName={project.name}
          teamName={team.name}
        />
      {:else if activeTab === "team"}
        <StudentProjectTeamChatTab {team} {currentStudentId} />
      {:else if activeTab === "members"}
        <StudentProjectMembersTab {team} {isLeader} />
      {/if}
    </div>
  </div>
</div>

<style>
  .detail-container {
    max-width: 1200px;
    margin: 0 auto;
  }

  /* Header Card */
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
    margin-bottom: 16px;
    transition: color 0.2s;
    padding: 0;
  }

  .back-button:hover {
    color: #111827;
  }

  .header-content {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
  }

  .header-main {
    flex: 1;
  }

  .project-title {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .project-description {
    color: #64748b;
    margin-bottom: 12px;
  }

  .project-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    font-size: 14px;
    color: #64748b;
    margin-bottom: 12px;
  }

  .separator {
    color: #d1d5db;
  }

  .badges {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .badge {
    padding: 4px 12px;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 500;
  }

  .badge-team {
    background: #dbeafe;
    color: #1e40af;
  }

  .badge-leader {
    background: #fef3c7;
    color: #92400e;
  }

  /* Tabs */
  .tabs-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .tabs-header {
    display: flex;
    gap: 32px;
    padding: 0 24px;
    border-bottom: 1px solid #e5e7eb;
  }

  .tab-button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 16px 0;
    border: none;
    background: none;
    color: #64748b;
    cursor: pointer;
    border-bottom: 2px solid transparent;
    transition: all 0.2s;
    font-size: 14px;
  }

  .tab-button img {
    opacity: 0.6;
    transition: opacity 0.2s;
  }

  .tab-button:hover {
    color: #111827;
  }

  .tab-button:hover img {
    opacity: 0.8;
  }

  .tab-button.active {
    color: #3b82f6;
    border-bottom-color: #3b82f6;
  }

  .tab-button.active img {
    opacity: 1;
    filter: brightness(0) saturate(100%) invert(45%) sepia(96%) saturate(2270%)
      hue-rotate(201deg) brightness(98%) contrast(94%);
  }

  .tab-content {
    padding: 24px;
  }

  .tab-panel {
  }

  .panel-title {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .text-gray {
    color: #64748b;
    margin-bottom: 24px;
  }

  .empty-content {
    padding: 48px;
    text-align: center;
    color: #9ca3af;
  }

  /* Members List */
  .members-list {
    display: grid;
    gap: 12px;
  }

  .member-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    transition: all 0.2s;
  }

  .member-card:hover {
    background: #f9fafb;
  }

  .member-card.leader {
    background: #eff6ff;
    border-color: #93c5fd;
  }

  .member-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: #3b82f6;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    flex-shrink: 0;
  }

  .member-info {
    flex: 1;
  }

  .member-name {
    font-weight: 500;
    color: #1e293b;
    margin-bottom: 4px;
  }

  .member-id {
    font-size: 14px;
    color: #64748b;
  }

  .leader-badge {
    padding: 4px 12px;
    background: #fef3c7;
    color: #92400e;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
  }
</style>
