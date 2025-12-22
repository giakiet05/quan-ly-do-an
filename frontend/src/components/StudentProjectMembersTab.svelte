<script lang="ts">
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

  let { team, isLeader } = $props<{
    team: Team;
    isLeader: boolean;
  }>();
</script>

<div class="members-container">
  <!-- Team Info -->
  <div class="team-info-card">
    <h3 class="section-title">Thông tin nhóm</h3>
    <div class="info-grid">
      <div class="info-item">
        <span class="info-label">Tên nhóm:</span>
        <span class="info-value">{team.name}</span>
      </div>
      <div class="info-item">
        <span class="info-label">Số thành viên:</span>
        <span class="info-value">{team.members.length + 1} người</span>
      </div>
      <div class="info-item">
        <span class="info-label">Trạng thái:</span>
        <span class="info-value">
          {#if team.status === "registered"}
            Đã đăng ký
          {:else if team.status === "ready"}
            Sẵn sàng
          {:else}
            Đang tuyển
          {/if}
        </span>
      </div>
      {#if team.createdAt}
        <div class="info-item">
          <span class="info-label">Ngày tạo:</span>
          <span class="info-value">
            {new Date(team.createdAt).toLocaleDateString("vi-VN")}
          </span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Members List -->
  <div class="members-section">
    <h3 class="section-title">Danh sách thành viên</h3>
    <div class="members-list">
      <!-- Leader -->
      <div class="member-card leader-card">
        <div class="member-content">
          <div class="member-avatar leader-avatar">
            <svg
              width="32"
              height="32"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"
              ></path>
            </svg>
          </div>
          <div class="member-info">
            <div class="member-header">
              <h4 class="member-name">{team.leaderName}</h4>
              <span class="role-badge leader-badge">Nhóm trưởng</span>
            </div>
            <div class="member-details">
              <div class="detail-item">
                <img
                  src="/Calendar_duotone.svg"
                  alt="Calendar"
                  width="16"
                  height="16"
                />
                <span>
                  Tham gia: {team.createdAt
                    ? new Date(team.createdAt).toLocaleDateString("vi-VN")
                    : "N/A"}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Members -->
      {#each team.members as member (member.id)}
        <div class="member-card">
          <div class="member-content">
            <div class="member-avatar">
              <span>{member.name.charAt(0)}</span>
            </div>
            <div class="member-info">
              <div class="member-header">
                <h4 class="member-name">{member.name}</h4>
                <span class="role-badge">Thành viên</span>
              </div>
              <div class="member-details">
                {#if member.email}
                  <div class="detail-item">
                    <svg
                      width="16"
                      height="16"
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
                    <span>{member.email}</span>
                  </div>
                {/if}
                <div class="detail-item">
                  <span class="mssv-label">MSSV: {member.studentId}</span>
                </div>
                {#if member.joinedAt}
                  <div class="detail-item">
                    <img
                      src="/Calendar_duotone.svg"
                      alt="Calendar"
                      width="16"
                      height="16"
                    />
                    <span>
                      Tham gia: {new Date(member.joinedAt).toLocaleDateString(
                        "vi-VN"
                      )}
                    </span>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- Empty State -->
  {#if team.members.length === 0}
    <div class="empty-state">
      <p class="empty-message">Chỉ có nhóm trưởng</p>
      <p class="empty-hint">Hãy mời thêm thành viên vào nhóm!</p>
    </div>
  {/if}
</div>

<style>
  .members-container {
    max-width: 1000px;
  }

  /* Team Info Card */
  .team-info-card {
    background: rgba(59, 130, 246, 0.05);
    border: 1px solid rgba(59, 130, 246, 0.2);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 24px;
  }

  .section-title {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 16px;
  }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .info-item {
    font-size: 14px;
  }

  .info-label {
    color: #64748b;
  }

  .info-value {
    margin-left: 8px;
    color: #1e293b;
  }

  /* Members Section */
  .members-section {
  }

  .members-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  /* Member Card */
  .member-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 20px;
    transition: box-shadow 0.2s;
  }

  .member-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .leader-card {
    border: 2px solid #fbbf24;
    background: #fffbeb;
  }

  .member-content {
    display: flex;
    align-items: flex-start;
    gap: 16px;
  }

  .member-avatar {
    width: 64px;
    height: 64px;
    background: #dbeafe;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #1e40af;
    font-size: 20px;
    font-weight: 600;
    flex-shrink: 0;
  }

  .leader-avatar {
    background: #fef3c7;
    color: #b45309;
  }

  .leader-avatar svg {
    fill: #b45309;
    stroke: #b45309;
  }

  .member-info {
    flex: 1;
  }

  .member-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
  }

  .member-name {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
  }

  .role-badge {
    padding: 4px 12px;
    background: #f1f5f9;
    color: #334155;
    border-radius: 16px;
    font-size: 12px;
    font-weight: 500;
  }

  .leader-badge {
    background: #fef3c7;
    color: #b45309;
  }

  .member-details {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .detail-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #64748b;
  }

  .detail-item svg {
    flex-shrink: 0;
  }

  .mssv-label {
    font-weight: 500;
  }

  /* Empty State */
  .empty-state {
    text-align: center;
    padding: 48px 24px;
    color: #64748b;
  }

  .empty-message {
    font-size: 16px;
    margin-bottom: 4px;
  }

  .empty-hint {
    font-size: 14px;
    margin-top: 4px;
  }
</style>
