<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";

  let { params } = $props<{
    params: { id: string; categoryId: string; projectId: string };
  }>();

  let loading = $state(true);
  let registrationState = $state<
    "not-registered" | "forming-team" | "registered"
  >("not-registered");
  let showInviteModal = $state(false);
  let inviteEmail = $state("");

  // Mock data
  let project = $state({
    id: params.projectId,
    name: "Hệ thống quản lý thư viện trực tuyến",
    description:
      "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
    instructor: "TS. Nguyễn Văn A",
    minTeamMembers: 2,
    maxTeamMembers: 3,
    tags: ["Web", "React", "Node.js", "MongoDB"],
  });

  let currentUser = $state({
    id: "user1",
    name: "Nguyễn Văn An",
    email: "nguyenvanan@student.edu.vn",
  });

  let team = $state<any>(null);

  let pendingInvitations = $state<any[]>([]);

  onMount(() => {
    loading = false;
    // Simulate checking registration status
    // registrationState = 'not-registered' | 'forming-team' | 'registered'
  });

  function handleBack() {
    push(`/classes/${params.id}/categories/${params.categoryId}`);
  }

  async function handleStartRegistration() {
    // Create group and set as leader
    team = {
      id: "team1",
      leaderId: currentUser.id,
      members: [
        {
          id: currentUser.id,
          name: currentUser.name,
          email: currentUser.email,
          role: "leader",
          joinedAt: new Date().toISOString(),
        },
      ],
      createdAt: new Date().toISOString(),
    };
    registrationState = "forming-team";
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

    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(inviteEmail)) {
      alert("Email không hợp lệ");
      return;
    }

    // Send invitation
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
    if (!confirm("Bạn có chắc chắn muốn hủy lời mời này?")) return;
    pendingInvitations = pendingInvitations.filter(
      (inv) => inv.id !== inviteId
    );
  }

  async function handleCompleteRegistration() {
    if (team.members.length < project.minTeamMembers) {
      alert(
        `Nhóm cần ít nhất ${project.minTeamMembers} thành viên để đăng ký đề tài`
      );
      return;
    }

    if (
      !confirm(
        "Bạn có chắc chắn muốn hoàn tất đăng ký? Sau khi đăng ký, bạn sẽ không thể thay đổi thành viên nhóm."
      )
    ) {
      return;
    }

    // Submit registration
    registrationState = "registered";
    alert("Đăng ký đề tài thành công!");
  }

  function handleViewMyProject() {
    push("/my-projects");
  }

  function getInitials(name: string): string {
    const words = name.split(" ");
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  }

  function formatDate(dateString: string): string {
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
        <h1 class="project-title">{project.name}</h1>
        <p class="project-description">{project.description}</p>

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
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            <span>GVHD: {project.instructor}</span>
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
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            <span
              >Nhóm: {project.minTeamMembers}-{project.maxTeamMembers} thành viên</span
            >
          </div>
        </div>

        <div class="project-tags">
          {#each project.tags as tag}
            <span class="tag">{tag}</span>
          {/each}
        </div>
      </div>
    </div>

    <div class="content-section">
      {#if registrationState === "not-registered"}
        <!-- State 1: Not registered yet -->
        <div class="registration-prompt">
          <div class="prompt-icon">
            <svg
              width="64"
              height="64"
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
          </div>
          <h2>Bạn chưa đăng ký đề tài này</h2>
          <p>
            Để đăng ký đề tài, bạn cần tạo nhóm và mời thành viên tham gia. Nhóm
            cần có từ {project.minTeamMembers} đến {project.maxTeamMembers}
            thành viên.
          </p>
          <button onclick={handleStartRegistration} class="btn-primary-large">
            Đăng ký đề tài
          </button>
        </div>
      {:else if registrationState === "forming-team"}
        <!-- State 2: Forming team -->
        <div class="team-formation">
          <div class="alert alert-info">
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            <div>
              <strong>Đang tạo nhóm</strong>
              <p>
                Bạn đang là trưởng nhóm. Hãy mời thêm thành viên để hoàn tất
                đăng ký đề tài.
              </p>
            </div>
          </div>

          {#if team.members.length < project.minTeamMembers}
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
                <strong>Chưa đủ thành viên</strong>
                <p>
                  Nhóm cần ít nhất {project.minTeamMembers} thành viên. Hiện tại:
                  {team.members.length}/{project.minTeamMembers}
                </p>
              </div>
            </div>
          {/if}

          <div class="section-card">
            <div class="section-header">
              <h3>
                Thành viên nhóm ({team.members.length}/{project.maxTeamMembers})
              </h3>
              {#if team.members.length < project.maxTeamMembers}
                <button onclick={handleOpenInviteModal} class="btn-secondary">
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                  Mời thành viên
                </button>
              {/if}
            </div>

            <div class="members-list">
              {#each team.members as member}
                <div class="member-card">
                  <div class="member-avatar">
                    {getInitials(member.name)}
                  </div>
                  <div class="member-info">
                    <div class="member-header">
                      <h4 class="member-name">{member.name}</h4>
                      {#if member.role === "leader"}
                        <span class="leader-badge">
                          <svg
                            width="14"
                            height="14"
                            viewBox="0 0 24 24"
                            fill="currentColor"
                          >
                            <path
                              d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                            />
                          </svg>
                          Trưởng nhóm
                        </span>
                      {/if}
                    </div>
                    <p class="member-email">{member.email}</p>
                    <p class="member-joined">
                      Tham gia: {formatDate(member.joinedAt)}
                    </p>
                  </div>
                </div>
              {/each}
            </div>
          </div>

          {#if pendingInvitations.length > 0}
            <div class="section-card">
              <h3 class="section-title">
                Lời mời đang chờ ({pendingInvitations.length})
              </h3>

              <div class="invitations-list">
                {#each pendingInvitations as invitation}
                  <div class="invitation-card">
                    <div class="invitation-info">
                      <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path
                          d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
                        ></path>
                        <polyline points="22,6 12,13 2,6"></polyline>
                      </svg>
                      <div>
                        <p class="invitation-email">{invitation.email}</p>
                        <p class="invitation-date">
                          Đã gửi: {formatDate(invitation.sentAt)}
                        </p>
                      </div>
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
              disabled={team.members.length < project.minTeamMembers}
            >
              Hoàn tất đăng ký
            </button>
          </div>
        </div>
      {:else if registrationState === "registered"}
        <!-- State 3: Registered -->
        <div class="registration-success">
          <div class="success-icon">
            <svg
              width="64"
              height="64"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
          </div>
          <h2>Đã đăng ký đề tài thành công!</h2>
          <p>
            Bạn và nhóm đã đăng ký đề tài "{project.name}" thành công. Giảng
            viên hướng dẫn sẽ liên hệ với nhóm trong thời gian sớm nhất.
          </p>

          <div class="registered-team-info">
            <h3>Thông tin nhóm</h3>
            <div class="members-list">
              {#each team.members as member}
                <div class="member-card">
                  <div class="member-avatar">
                    {getInitials(member.name)}
                  </div>
                  <div class="member-info">
                    <div class="member-header">
                      <h4 class="member-name">{member.name}</h4>
                      {#if member.role === "leader"}
                        <span class="leader-badge">
                          <svg
                            width="14"
                            height="14"
                            viewBox="0 0 24 24"
                            fill="currentColor"
                          >
                            <path
                              d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                            />
                          </svg>
                          Trưởng nhóm
                        </span>
                      {/if}
                    </div>
                    <p class="member-email">{member.email}</p>
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
        <label for="invite-email" class="form-label"> Email sinh viên </label>
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
