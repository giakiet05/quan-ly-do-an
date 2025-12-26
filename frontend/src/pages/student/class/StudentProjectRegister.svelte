<script lang="ts">
  import { push } from "svelte-spa-router";
  import type { Project } from "../../../types/project";

  interface Team {
    id: string;
    leaderId: string;
    leaderName: string;
    members: {
      id: string;
      name: string;
      studentCode: string;
    }[];
    status: "forming" | "registered";
  }

  interface StudentProjectRegisterProps {
    classId: string;
    categoryId: string;
    projectId: string;
  }

  let { classId, categoryId, projectId }: StudentProjectRegisterProps =
    $props();

  let showInviteModal = $state(false);
  let inviteEmail = $state("");

  // Mock current student
  const currentStudentId = "student-1";
  const currentStudentName = "Nguyễn Văn A";
  const currentStudentCode = "SV001";
  const currentStudentEmail = "student1@example.com";

  // Mock project
  const project: Project & {
    className?: string;
    categoryName?: string;
    categoryStatus?: "upcoming" | "ongoing" | "completed";
    projectNumber?: number;
  } = {
    id: projectId,
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
    className: "IT4441 - Phát triển ứng dụng Web",
    categoryName: "Đồ án 1",
    categoryStatus: "ongoing",
    projectNumber: 1,
  };

  // Mock team - null = chưa đăng ký, có data = đã có nhóm
  let myTeam = $state<Team | null>(null);
  /*
  {
    id: "team-1",
    leaderId: currentStudentId,
    leaderName: currentStudentName,
    members: [
      { id: "student-2", name: "Trần Thị B", studentCode: "SV002" },
    ],
    status: "forming",
  }
  */

  // Mock pending invites
  let pendingInvites = $state<{ id: string; email: string; sentAt: string }[]>(
    []
  );

  const isLeader = $derived(myTeam?.leaderId === currentStudentId);
  const memberCount = $derived(myTeam ? myTeam.members.length + 1 : 0);
  const minMembers = project.minTeamMembers || 1;
  const maxMembers = project.maxTeamMembers || project.maxStudents;
  const isTeamReady = $derived(
    memberCount >= minMembers && memberCount <= maxMembers
  );
  const isRegistered = $derived(myTeam?.status === "registered");

  function handleBack(): void {
    push(`/classes/${classId}/categories/${categoryId}`);
  }

  function handleRegisterProject(): void {
    // Tạo nhóm mới với sinh viên hiện tại là leader
    myTeam = {
      id: "team-" + Date.now(),
      leaderId: currentStudentId,
      leaderName: currentStudentName,
      members: [],
      status: "forming",
    };
    alert("Đăng ký thành công! Bạn là nhóm trưởng.");
  }

  function handleInvite(): void {
    if (!inviteEmail.trim()) {
      alert("Vui lòng nhập email sinh viên");
      return;
    }
    if (memberCount >= maxMembers) {
      alert("Nhóm đã đủ số lượng thành viên tối đa");
      return;
    }

    // Thêm vào pending invites
    pendingInvites = [
      ...pendingInvites,
      {
        id: "invite-" + Date.now(),
        email: inviteEmail,
        sentAt: "Vừa xong",
      },
    ];

    alert(`Đã gửi lời mời tới ${inviteEmail}`);
    inviteEmail = "";
    showInviteModal = false;
  }

  function handleRemoveMember(memberId: string): void {
    if (!myTeam) return;
    myTeam.members = myTeam.members.filter((m) => m.id !== memberId);
  }

  function handleConfirmRegistration(): void {
    if (memberCount < minMembers) {
      const confirmMsg = `Nhóm hiện có ${memberCount} thành viên, chưa đủ tối thiểu (${minMembers}). Bạn vẫn muốn đăng ký chính thức?`;
      if (!confirm(confirmMsg)) return;
    }
    if (
      confirm(
        "Sau khi đăng ký chính thức, bạn không thể thay đổi thành viên. Xác nhận đăng ký?"
      )
    ) {
      if (myTeam) {
        myTeam.status = "registered";
      }
      alert("Đăng ký chính thức thành công!");
    }
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
          <h1 class="project-title">{project.name}</h1>
          {#if project.projectNumber}
            <span class="number-badge">#{project.projectNumber}</span>
          {/if}
          {#if isRegistered}
            <span class="registered-badge">
              <img src="/check-circle.svg" alt="check" width="16" height="16" />
              Đã đăng ký
            </span>
          {/if}
        </div>
        <p class="project-description">{project.description}</p>

        <div class="project-meta">
          {#if project.className}
            <div class="meta-item">
              <img
                src="/book-open-text.svg"
                alt="class"
                width="16"
                height="16"
              />
              <span>{project.className}</span>
            </div>
            <span>•</span>
          {/if}
          {#if project.categoryName}
            <div class="meta-item">
              <img src="/calendar.svg" alt="category" width="16" height="16" />
              <span>{project.categoryName}</span>
            </div>
            <span>•</span>
          {/if}
          <span>GVHD: {project.instructor}</span>
        </div>

        <div class="project-tags">
          {#each project.tags as tag}
            <span class="tag">
              <img src="/tag.svg" alt="tag" width="12" height="12" />
              {tag}
            </span>
          {/each}
        </div>
      </div>
    </div>

    <!-- Yêu cầu nhóm -->
    <div class="requirement-box">
      <div class="requirement-title">Yêu cầu nhóm</div>
      <div class="requirement-text">
        <span>Tối thiểu: {minMembers} người</span>
        <span>•</span>
        <span>Tối đa: {maxMembers} người</span>
      </div>
    </div>
  </div>

  <!-- Nếu chưa đăng ký -->
  {#if !myTeam && project.status === "available" && project.categoryStatus === "ongoing"}
    <div class="action-card">
      <div class="action-content">
        <img
          src="/users.svg"
          alt="users"
          width="64"
          height="64"
          class="action-icon"
        />
        <h3 class="action-title">Đăng ký đề tài này</h3>
        <p class="action-text">
          Bạn sẽ tự động trở thành nhóm trưởng và có thể mời thêm thành viên vào
          nhóm
        </p>
        <button onclick={handleRegisterProject} class="register-button">
          Đăng ký đề tài
        </button>
      </div>
    </div>
  {/if}

  <!-- Nếu đề tài đã bị đăng ký -->
  {#if !myTeam && (project.status === "forming" || project.status === "full" || project.status === "closed")}
    <div class="action-card">
      <div class="action-content">
        <img
          src="/alert-circle.svg"
          alt="alert"
          width="64"
          height="64"
          class="action-icon warning"
        />
        <h3 class="action-title">Đề tài này đã có người đăng ký</h3>
        <p class="action-text">
          {#if project.status === "forming"}
            Một nhóm đang hình thành.
          {:else if project.status === "full"}
            Đề tài đã đủ người.
          {/if}
          Vui lòng chọn đề tài khác
        </p>
      </div>
    </div>
  {/if}

  <!-- Nếu đã có nhóm -->
  {#if myTeam}
    <div class="team-card">
      <div class="team-header">
        <div>
          <h2 class="team-title">Nhóm của bạn</h2>
          <div class="team-stats">
            <div class="stat-item">
              <img src="/users.svg" alt="users" width="16" height="16" />
              <span>{memberCount}/{maxMembers} thành viên</span>
            </div>
            {#if memberCount < minMembers}
              <span>•</span>
              <span class="warning-text"
                >Cần thêm {minMembers - memberCount} người</span
              >
            {/if}
            {#if isTeamReady && !isRegistered}
              <span>•</span>
              <span class="success-text">Sẵn sàng đăng ký</span>
            {/if}
          </div>
        </div>
        {#if isLeader && !isRegistered}
          <button
            onclick={() => (showInviteModal = true)}
            class="invite-button"
            disabled={memberCount >= maxMembers}
          >
            <img src="/user-plus.svg" alt="invite" width="16" height="16" />
            Mời thành viên
          </button>
        {/if}
      </div>

      <!-- Team Members -->
      <div class="members-list">
        <!-- Leader -->
        <div class="member-card leader">
          <div class="member-avatar leader-avatar">
            <img src="/crown.svg" alt="leader" width="24" height="24" />
          </div>
          <div class="member-info">
            <div class="member-name">
              {myTeam.leaderName}
              <span class="leader-badge">Nhóm trưởng</span>
            </div>
            <div class="member-code">{currentStudentCode}</div>
          </div>
          {#if myTeam.leaderId === currentStudentId}
            <span class="you-badge">(Bạn)</span>
          {/if}
        </div>

        <!-- Members -->
        {#each myTeam.members as member (member.id)}
          <div class="member-card">
            <div class="member-avatar">{member.name.charAt(0)}</div>
            <div class="member-info">
              <div class="member-name">
                {member.name}
                {#if member.id === currentStudentId}
                  <span class="you-badge">(Bạn)</span>
                {/if}
              </div>
              <div class="member-code">{member.studentCode}</div>
            </div>
            {#if isLeader && !isRegistered}
              <button
                onclick={() => {
                  if (confirm(`Xóa ${member.name} khỏi nhóm?`)) {
                    handleRemoveMember(member.id);
                  }
                }}
                class="remove-button"
              >
                <img src="/x.svg" alt="remove" width="20" height="20" />
              </button>
            {/if}
          </div>
        {/each}
      </div>

      <!-- Pending Invites -->
      {#if isLeader && pendingInvites.length > 0 && !isRegistered}
        <div class="pending-invites-box">
          <h4 class="pending-title">
            Lời mời đang chờ ({pendingInvites.length})
          </h4>
          <div class="pending-list">
            {#each pendingInvites as invite (invite.id)}
              <div class="pending-item">
                <span class="pending-email">{invite.email}</span>
                <span class="pending-time">Đã gửi {invite.sentAt}</span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Warning -->
      {#if memberCount < minMembers && !isRegistered}
        <div class="warning-box">
          <img src="/alert-circle.svg" alt="warning" width="20" height="20" />
          <div>
            <div class="warning-title">
              Chưa đủ số lượng thành viên tối thiểu
            </div>
            <div class="warning-text">
              Nhóm cần ít nhất {minMembers} thành viên. Hiện tại có {memberCount}
              người. Bạn vẫn có thể đăng ký nhưng nên bổ sung thêm thành viên.
            </div>
          </div>
        </div>
      {/if}

      <!-- Confirm Registration Button -->
      {#if isLeader && !isRegistered}
        <button onclick={handleConfirmRegistration} class="confirm-button">
          Đăng ký chính thức
        </button>
      {/if}

      <!-- Registered Status -->
      {#if isRegistered}
        <div class="success-box">
          <img src="/check-circle.svg" alt="success" width="20" height="20" />
          <div>
            <div class="success-title">Đã đăng ký thành công</div>
            <p class="success-text">
              Nhóm của bạn đã đăng ký đề tài này. Không thể thay đổi thành viên.
            </p>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Invite Modal -->
{#if showInviteModal}
  <div class="modal-overlay" onclick={() => (showInviteModal = false)}>
    <div
      class="modal-content"
      onclick={(e) => e.stopPropagation()}
      role="dialog"
    >
      <h3 class="modal-title">Mời thành viên vào nhóm</h3>
      <div class="modal-body">
        <label for="invite-email" class="input-label">Email sinh viên</label>
        <input
          id="invite-email"
          type="email"
          bind:value={inviteEmail}
          class="email-input"
          placeholder="student@example.com"
        />
      </div>
      <p class="modal-hint">
        Sinh viên sẽ nhận được lời mời tham gia nhóm của bạn
      </p>
      <div class="modal-actions">
        <button
          onclick={() => {
            showInviteModal = false;
            inviteEmail = "";
          }}
          class="cancel-button"
        >
          Hủy
        </button>
        <button onclick={handleInvite} class="send-button">
          Gửi lời mời
        </button>
      </div>
    </div>
  </div>
{/if}

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

  .header-badges {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .project-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1e293b;
  }

  .number-badge {
    padding: 4px 12px;
    background: #f1f5f9;
    color: #475569;
    border-radius: 9999px;
    font-size: 0.875rem;
  }

  .registered-badge {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 12px;
    background: #dcfce7;
    color: #166534;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .project-description {
    color: #64748b;
    margin-bottom: 16px;
    line-height: 1.5;
  }

  .project-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
    margin-bottom: 12px;
    flex-wrap: wrap;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .project-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .tag {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    background: #f1f5f9;
    color: #475569;
    border-radius: 9999px;
    font-size: 0.875rem;
  }

  .requirement-box {
    margin-top: 16px;
    padding: 16px;
    background: #f9fafb;
    border-radius: 8px;
  }

  .requirement-title {
    font-size: 0.875rem;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .requirement-text {
    display: flex;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .action-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 24px;
  }

  .action-content {
    text-align: center;
    padding: 48px 24px;
  }

  .action-icon {
    margin: 0 auto 16px;
    opacity: 0.7;
  }

  .action-icon.warning {
    filter: brightness(0) saturate(100%) invert(60%) sepia(98%) saturate(473%)
      hue-rotate(345deg) brightness(98%) contrast(98%);
  }

  .action-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 8px;
  }

  .action-text {
    color: #64748b;
    margin-bottom: 24px;
    line-height: 1.5;
  }

  .register-button {
    padding: 12px 24px;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.9375rem;
    font-weight: 500;
    transition: background 0.2s;
  }

  .register-button:hover {
    background: #1d4ed8;
  }

  .team-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 24px;
  }

  .team-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
    gap: 16px;
  }

  .team-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 4px;
  }

  .team-stats {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 4px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .warning-text {
    color: #c2410c;
    font-size: 0.875rem;
  }

  .success-text {
    color: #166534;
    font-size: 0.875rem;
  }

  .invite-button {
    display: flex;
    align-items: center;
    gap: 8px;
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

  .invite-button:hover:not(:disabled) {
    background: #1d4ed8;
  }

  .invite-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .members-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 24px;
  }

  .member-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
  }

  .member-card.leader {
    background: #fffbeb;
    border-color: #fde68a;
  }

  .member-avatar {
    width: 48px;
    height: 48px;
    background: #dbeafe;
    color: #1e40af;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    flex-shrink: 0;
  }

  .leader-avatar {
    background: #fef3c7;
    color: #b45309;
  }

  .member-info {
    flex: 1;
  }

  .member-name {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 500;
    color: #1e293b;
    margin-bottom: 4px;
  }

  .leader-badge {
    padding: 2px 8px;
    background: #fde68a;
    color: #b45309;
    border-radius: 4px;
    font-size: 0.75rem;
  }

  .you-badge {
    color: #64748b;
    font-size: 0.875rem;
    font-weight: 400;
  }

  .member-code {
    color: #64748b;
    font-size: 0.875rem;
  }

  .remove-button {
    padding: 8px;
    color: #dc2626;
    background: none;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;
  }

  .remove-button:hover {
    background: #fef2f2;
  }

  .pending-invites-box {
    padding: 16px;
    background: #fff7ed;
    border: 1px solid #fed7aa;
    border-radius: 8px;
    margin-bottom: 24px;
  }

  .pending-title {
    font-size: 0.875rem;
    color: #475569;
    margin-bottom: 12px;
  }

  .pending-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .pending-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.875rem;
  }

  .pending-email {
    color: #475569;
  }

  .pending-time {
    color: #64748b;
    font-size: 0.75rem;
  }

  .warning-box {
    display: flex;
    gap: 12px;
    padding: 16px;
    background: #fff7ed;
    border: 1px solid #fed7aa;
    border-radius: 8px;
    margin-bottom: 24px;
  }

  .warning-title {
    color: #9a3412;
    font-weight: 500;
    margin-bottom: 4px;
  }

  .warning-box .warning-text {
    color: #c2410c;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  .confirm-button {
    width: 100%;
    padding: 12px;
    background: #16a34a;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.9375rem;
    font-weight: 500;
    transition: background 0.2s;
  }

  .confirm-button:hover {
    background: #15803d;
  }

  .success-box {
    display: flex;
    gap: 12px;
    padding: 16px;
    background: #dcfce7;
    border: 1px solid #86efac;
    border-radius: 8px;
  }

  .success-title {
    color: #166534;
    font-weight: 500;
    margin-bottom: 4px;
  }

  .success-text {
    color: #166534;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 50;
    padding: 16px;
  }

  .modal-content {
    background: white;
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    max-width: 28rem;
    width: 100%;
    padding: 24px;
  }

  .modal-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 16px;
  }

  .modal-body {
    margin-bottom: 16px;
  }

  .input-label {
    display: block;
    font-size: 0.875rem;
    color: #475569;
    margin-bottom: 8px;
  }

  .email-input {
    width: 100%;
    padding: 12px;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    font-size: 0.9375rem;
  }

  .email-input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
  }

  .modal-hint {
    color: #64748b;
    font-size: 0.875rem;
    margin-bottom: 16px;
    line-height: 1.5;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }

  .cancel-button {
    padding: 8px 16px;
    background: #e5e7eb;
    color: #374151;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    transition: background 0.2s;
  }

  .cancel-button:hover {
    background: #d1d5db;
  }

  .send-button {
    padding: 8px 16px;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    transition: background 0.2s;
  }

  .send-button:hover {
    background: #1d4ed8;
  }

  @media (max-width: 768px) {
    .page-container {
      padding: 16px;
    }

    .team-header {
      flex-direction: column;
    }

    .invite-button {
      width: 100%;
      justify-content: center;
    }
  }
</style>
