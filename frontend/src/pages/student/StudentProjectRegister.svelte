<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { getProject } from "../../services/project-service";
  import { getGroupsFilter, createGroup } from "../../services/group-service";
  import { authStore } from "../../stores/auth-store";
  import type { Group } from "../../models";
  import type { ProjectResponse } from "../../dtos/project-dto";
  import { getClassroomStudents } from "../../services/classroom-service";
  import {
    sendGroupInvitation,
    getGroupInvitations,
  } from "../../services/group-service";
  let { params } = $props<{
    params: { id: string; categoryId: string; projectId: string };
  }>();

  let loading = $state(true);
  let error = $state<string | null>(null);
  let registrationState = $state<
    "not-registered" | "forming-team" | "registered"
  >("not-registered");
  let showInviteModal = $state(false);
  let inviteEmail = $state("");

  let project = $state<ProjectResponse | null>(null);
  let myGroup = $state<Group | null>(null);
  let allGroups = $state<Group[]>([]);
  let pendingInvitations = $state<any[]>([]);

  // Hàm load dữ liệu (có thể gọi lại khi cần)
  async function loadData() {
    try {
      loading = true;
      error = null;

      // 1. Load chi tiết đề tài
      try {
        const projectData = await getProject(params.id, params.projectId);
        console.log("Project loaded:", projectData);
        project = projectData;
      } catch (projectErr) {
        console.warn("Không load được chi tiết đề tài:", projectErr);
        error = "Không thể tải thông tin đề tài. Bạn vẫn có thể đăng ký nhóm.";
      }

      // 2. Load tất cả nhóm của project
      const groups = await getGroupsFilter({
        classroom_id: params.id,
        project_id: params.projectId,
      });

      allGroups = groups || [];
      console.log("Groups for project:", allGroups);

      // 3. Tìm nhóm của user hiện tại
      const userGroup = allGroups.find((g) =>
        g.members?.some((m) => m.user_id === $authStore.user?.id),
      );

      myGroup = userGroup || null;
      console.log("My group:", myGroup);

      // 4. Cập nhật trạng thái
      if (myGroup) {
        registrationState = "forming-team";
      } else {
        registrationState = "not-registered";
      }
    } catch (err) {
      console.error("Error loading data:", err);
      error = err instanceof Error ? err.message : "Không thể tải dữ liệu";
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  async function handleStartRegistration() {
    try {
      if (!$authStore.user) {
        alert("Vui lòng đăng nhập");
        return;
      }

      // Kiểm tra đề tài đã đủ nhóm chưa
      if (project) {
        const totalGroups = allGroups.length;
        const maxGroups = project.amount || 0;
        if (totalGroups >= maxGroups) {
          alert(`Đề tài đã đủ số lượng nhóm (${totalGroups}/${maxGroups}).`);
          return;
        }
      }

      // Tạo nhóm mới
      const newGroup = await createGroup(
        params.id,
        params.projectId,
        params.categoryId,
      );

      // Cập nhật lại dữ liệu sau khi tạo nhóm
      myGroup = newGroup;
      allGroups = [...allGroups, newGroup];
      registrationState = "forming-team";

      alert("Tạo nhóm thành công! Bạn có thể mời thêm thành viên.");
    } catch (err: any) {
      console.error("Error creating group:", err);
      alert(err?.message || "Không thể tạo nhóm. Vui lòng thử lại.");
    }
  }

  function handleOpenInviteModal() {
    showInviteModal = true;
    inviteEmail = "";
  }

  async function handleSendInvite() {
    if (!inviteEmail) {
      alert("Vui lòng nhập email");
      return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(inviteEmail)) {
      alert("Email không hợp lệ");
      return;
    }

    // TODO: Thay bằng API gửi lời mời thật khi backend hỗ trợ
    const newInvitation = {
      id: `inv${Date.now()}`,
      email: inviteEmail,
      status: "pending",
      sentAt: new Date().toISOString(),
    };

    pendingInvitations = [...pendingInvitations, newInvitation];
    showInviteModal = false;
    alert(`Đã gửi lời mời đến ${inviteEmail}`);
  }

  function handleCancelInvite(inviteId: string) {
    if (!confirm("Hủy lời mời này?")) return;
    pendingInvitations = pendingInvitations.filter(
      (inv) => inv.id !== inviteId,
    );
  }

  async function handleCompleteRegistration() {
    if (!myGroup || myGroup.members.length === 0) {
      alert("Nhóm cần ít nhất 1 thành viên");
      return;
    }

    if (!confirm("Hoàn tất đăng ký? Bạn sẽ không thể thay đổi nhóm nữa.")) {
      return;
    }

    registrationState = "registered";
    alert("Đăng ký đề tài thành công!");
  }

  function handleViewMyProject() {
    push("/my-projects");
  }

  function handleBack() {
    push(`/student/classes/${params.id}/categories/${params.categoryId}`);
  }

  function getInitials(name: string | undefined): string {
    if (!name) return "?";
    const words = name.split(" ");
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  }

  function formatDate(dateString: string | undefined): string {
    if (!dateString) return "N/A";
    return new Date(dateString).toLocaleDateString("vi-VN", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

<div class="container">
  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Đang tải thông tin đề tài...</p>
    </div>
  {:else if error}
    <div class="error-message">
      <p>{error}</p>
      <button onclick={handleBack}>Quay lại</button>
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
        Quay lại danh sách đề tài
      </button>

      <div class="project-info">
        {#if project}
          <h1 class="project-title">{project.title}</h1>
          <p class="project-description">
            {project.description || "Không có mô tả"}
          </p>

          <div class="project-meta">
            <div class="meta-item">
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
                >Nhóm: {project.minMember || 0}–{project.maxMember || 3} thành viên</span
              >
            </div>
          </div>
        {:else}
          <div class="alert alert-warning">
            <strong>Không tải được thông tin đề tài</strong>
            <p>Bạn vẫn có thể tạo nhóm và đăng ký nếu biết thông tin.</p>
          </div>
        {/if}
      </div>
    </div>

    <div class="content-section">
      {#if registrationState === "not-registered"}
        <div class="registration-prompt">
          <div class="prompt-icon">👥</div>
          <h2>Bạn chưa đăng ký đề tài này</h2>
          <p>
            Để đăng ký, bạn cần tạo nhóm và mời thành viên. Nhóm cần từ
            {project?.minMember || 2} đến {project?.maxMember || 3} người.
          </p>
          <button onclick={handleStartRegistration} class="btn-primary-large">
            Tạo nhóm & đăng ký
          </button>
        </div>
      {:else if registrationState === "forming-team"}
        <div class="team-formation">
          <div class="alert alert-info">
            <strong>Đang tạo nhóm</strong>
            <p>Bạn là trưởng nhóm. Mời thêm thành viên để hoàn tất.</p>
          </div>

          {#if myGroup && myGroup.members.length < (project?.minMember || 2)}
            <div class="alert alert-warning">
              <strong>Chưa đủ thành viên</strong>
              <p>
                Cần ít nhất {project?.minMember || 2} người. Hiện tại: {myGroup
                  .members.length}
              </p>
            </div>
          {/if}

          <div class="section-card">
            <div class="section-header">
              <h3>
                Thành viên nhóm ({myGroup?.members.length ||
                  0}/{project?.maxMember || 3})
              </h3>
              {#if myGroup && myGroup.members.length < (project?.maxMember || 3)}
                <button onclick={handleOpenInviteModal} class="btn-secondary">
                  Mời thành viên
                </button>
              {/if}
            </div>

            <div class="members-list">
              {#each myGroup?.members || [] as member}
                <div class="member-card">
                  <div class="member-avatar">
                    {getInitials(member.full_name)}
                  </div>
                  <div class="member-info">
                    <h4 class="member-name">{member.full_name}</h4>
                    {#if member.role === "leader"}
                      <span class="leader-badge">Trưởng nhóm</span>
                    {/if}
                    <p class="member-email">{member.email}</p>
                  </div>
                </div>
              {/each}
            </div>
          </div>

          {#if pendingInvitations.length > 0}
            <div class="section-card">
              <h3>Lời mời đang chờ ({pendingInvitations.length})</h3>
              <div class="invitations-list">
                {#each pendingInvitations as invitation}
                  <div class="invitation-card">
                    <div>
                      <p class="invitation-email">{invitation.email}</p>
                      <p class="invitation-date">
                        Gửi lúc: {formatDate(invitation.sentAt)}
                      </p>
                    </div>
                    <button
                      onclick={() => handleCancelInvite(invitation.id)}
                      class="btn-cancel"
                    >
                      Hủy
                    </button>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <div class="action-buttons">
            <button
              onclick={handleCompleteRegistration}
              class="btn-primary-large"
              disabled={!myGroup ||
                myGroup.members.length < (project?.minMember || 2)}
            >
              Hoàn tất đăng ký
            </button>
          </div>
        </div>
      {:else if registrationState === "registered"}
        <div class="registration-success">
          <div class="success-icon">✅</div>
          <h2>Đăng ký đề tài thành công!</h2>
          <p>Bạn và nhóm đã đăng ký đề tài "{project?.title || "này"}".</p>

          <div class="registered-team-info">
            <h3>Thông tin nhóm</h3>
            <div class="members-list">
              {#each myGroup?.members || [] as member}
                <div class="member-card">
                  <div class="member-avatar">
                    {getInitials(member.full_name)}
                  </div>
                  <div class="member-info">
                    <h4>{member.full_name}</h4>
                    {#if member.role === "leader"}<span class="leader-badge"
                        >Trưởng nhóm</span
                      >{/if}
                    <p>{member.email}</p>
                  </div>
                </div>
              {/each}
            </div>
          </div>

          <button onclick={handleViewMyProject} class="btn-primary-large">
            Xem đề tài của tôi
          </button>
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Modal mời thành viên -->
<!-- Invite Modal -->
{#if showInviteModal}
  <div class="modal-overlay" onclick={() => (showInviteModal = false)}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3>Mời thành viên tham gia nhóm</h3>
        <button onclick={() => (showInviteModal = false)} class="modal-close">
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

      <div class="modal-body">
        <label for="invite-email" class="form-label">Email sinh viên</label>
        <input
          type="email"
          id="invite-email"
          bind:value={inviteEmail}
          placeholder="example@student.edu.vn"
          class="form-input"
        />
        <p class="form-help">Nhập email của sinh viên bạn muốn mời vào nhóm</p>
      </div>

      <div class="modal-footer">
        <button onclick={() => (showInviteModal = false)} class="btn-secondary">
          Hủy
        </button>
        <button onclick={handleSendInvite} class="btn-primary">
          Gửi lời mời
        </button>
      </div>
    </div>
  </div>
{/if}

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

  .project-title {
    font-size: 1.875rem;
    font-weight: 700;
    color: #111827;
    margin-bottom: 0.75rem;
  }

  .project-description {
    font-size: 1rem;
    color: #6b7280;
    margin-bottom: 1rem;
    line-height: 1.6;
  }

  .project-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .project-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .tag {
    padding: 0.25rem 0.75rem;
    background: #f3f4f6;
    color: #4b5563;
    font-size: 0.75rem;
    border-radius: 0.25rem;
  }

  .content-section {
    background: white;
    border-radius: 0.75rem;
    padding: 2rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .registration-prompt,
  .registration-success {
    text-align: center;
    padding: 3rem 2rem;
  }

  .prompt-icon,
  .success-icon {
    margin: 0 auto 1.5rem;
    color: #3b82f6;
  }

  .success-icon {
    color: #10b981;
  }

  .registration-prompt h2,
  .registration-success h2 {
    font-size: 1.5rem;
    font-weight: 700;
    color: #111827;
    margin-bottom: 0.75rem;
  }

  .registration-prompt p,
  .registration-success p {
    font-size: 1rem;
    color: #6b7280;
    margin-bottom: 2rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
  }

  .btn-primary-large {
    padding: 0.875rem 2rem;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-primary-large:hover:not(:disabled) {
    background: #2563eb;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px rgba(59, 130, 246, 0.3);
  }

  .btn-primary-large:disabled {
    background: #9ca3af;
    cursor: not-allowed;
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

  .alert-info {
    background: #dbeafe;
    color: #1e40af;
    border: 1px solid #3b82f6;
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

  .section-card {
    margin-bottom: 1.5rem;
    padding: 1.5rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.75rem;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .section-header h3,
  .section-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
  }

  .btn-secondary {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1rem;
    background: white;
    border: 1px solid #d1d5db;
    color: #374151;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-secondary:hover {
    background: #f9fafb;
    border-color: #9ca3af;
  }

  .members-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .member-card {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    background: #f9fafb;
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
    font-size: 1rem;
    font-weight: 600;
    flex-shrink: 0;
  }

  .member-info {
    flex: 1;
  }

  .member-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.25rem;
  }

  .member-name {
    font-size: 1rem;
    font-weight: 600;
    color: #111827;
  }

  .leader-badge {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.125rem 0.5rem;
    background: #fbbf24;
    color: white;
    font-size: 0.75rem;
    font-weight: 600;
    border-radius: 9999px;
  }

  .member-email,
  .member-joined {
    font-size: 0.875rem;
    color: #6b7280;
  }

  .invitations-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .invitation-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    background: #fffbeb;
  }

  .invitation-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .invitation-email {
    font-size: 0.875rem;
    font-weight: 500;
    color: #111827;
  }

  .invitation-date {
    font-size: 0.75rem;
    color: #6b7280;
  }

  .btn-cancel {
    padding: 0.375rem 0.875rem;
    background: white;
    border: 1px solid #ef4444;
    color: #ef4444;
    border-radius: 0.375rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-cancel:hover {
    background: #ef4444;
    color: white;
  }

  .action-buttons {
    display: flex;
    justify-content: center;
    margin-top: 2rem;
  }

  .registered-team-info {
    max-width: 600px;
    margin: 2rem auto;
    text-align: left;
  }

  .registered-team-info h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
    margin-bottom: 1rem;
  }

  /* Modal styles */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
  }

  .modal {
    background: white;
    border-radius: 0.75rem;
    max-width: 500px;
    width: 100%;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .modal-header h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #111827;
  }

  .modal-close {
    padding: 0.5rem;
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    border-radius: 0.375rem;
    transition: all 0.2s;
  }

  .modal-close:hover {
    background: #f3f4f6;
    color: #111827;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .form-label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: #374151;
    margin-bottom: 0.5rem;
  }

  .form-input {
    width: 100%;
    padding: 0.625rem 0.875rem;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    transition: all 0.2s;
  }

  .form-input:focus {
    outline: none;
    border-color: #3b82f6;
    ring: 2px;
    ring-color: rgba(59, 130, 246, 0.2);
  }

  .form-help {
    margin-top: 0.5rem;
    font-size: 0.75rem;
    color: #6b7280;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1.5rem;
    border-top: 1px solid #e5e7eb;
  }

  .btn-primary {
    padding: 0.625rem 1.25rem;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-primary:hover {
    background: #2563eb;
  }
</style>
