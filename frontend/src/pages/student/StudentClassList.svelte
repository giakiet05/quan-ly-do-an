<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { getMyJoinedClassrooms } from "../../services/classroom-service";
  import type { ClassroomResponse } from "../../dtos/classroom-dto";

  let classes = $state<ClassroomResponse[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      loading = true;
      error = null;

      // Fetch real data from API
      classes = await getMyJoinedClassrooms();

      // TEMPORARY: Add mock data if empty for UI testing
      if (classes.length === 0) {
        classes = [
          {
            id: "mock-class-1",
            name: "Đồ án Phát triển ứng dụng web",
            description:
              "Lớp học đồ án cuối kỳ HK2 2024-2025 - Phát triển ứng dụng web full-stack",
            semester: "HK2",
            year: 2024,
            invitationCode: "WEB2024",
            status: "active",
            generalChannelId: "channel1",
            lecturer: {
              userId: "lecturer1",
              fullName: "TS. Nguyễn Văn A",
              avatar: "",
            },
            students: Array(45).fill({
              userId: "s1",
              fullName: "Student",
              avatar: "",
            }),
            projectRounds: [
              {
                id: "round1",
                name: "Đợt 1 - Đồ án cuối kỳ",
                description: "Phát triển ứng dụng web hoàn chỉnh",
                startDate: "2024-02-01T00:00:00Z",
                endDate: "2024-05-31T23:59:59Z",
                projects: [
                  {
                    id: "proj1",
                    classroomId: "mock-class-1",
                    projectRoundId: "round1",
                    title: "Hệ thống quản lý thư viện",
                    amount: 5,
                    description:
                      "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách",
                    minMember: 2,
                    maxMember: 4,
                    status: "approved",
                  },
                  {
                    id: "proj2",
                    classroomId: "mock-class-1",
                    projectRoundId: "round1",
                    title: "Ứng dụng quản lý chi tiêu",
                    amount: 3,
                    description:
                      "Ứng dụng mobile giúp theo dõi thu chi cá nhân",
                    minMember: 2,
                    maxMember: 3,
                    status: "approved",
                  },
                ],
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
            ],
            maxStudents: 50,
            autoApprove: false,
            canStudentDeleteGroup: false,
            avatar: "",
            createdAt: "2024-01-15T00:00:00Z",
          },
          {
            id: "mock-class-2",
            name: "Đồ án Trí tuệ nhân tạo",
            description: "Áp dụng AI vào bài toán thực tế",
            semester: "HK2",
            year: 2024,
            invitationCode: "AI2024",
            status: "active",
            generalChannelId: "channel2",
            lecturer: {
              userId: "lecturer2",
              fullName: "TS. Trần Thị B",
              avatar: "",
            },
            students: Array(38).fill({
              userId: "s1",
              fullName: "Student",
              avatar: "",
            }),
            projectRounds: [
              {
                id: "round2",
                name: "Đợt 1 - AI Research",
                description: "Nghiên cứu và ứng dụng AI",
                startDate: "2024-02-01T00:00:00Z",
                endDate: "2024-05-31T23:59:59Z",
                projects: [],
                reportPeriods: [],
                createdAt: "2024-01-15T00:00:00Z",
                isDeleted: false,
              },
            ],
            maxStudents: 40,
            autoApprove: true,
            canStudentDeleteGroup: false,
            avatar: "",
            createdAt: "2024-01-15T00:00:00Z",
          },
        ] as any;
      }

      /* MOCK DATA - Uncomment để test UI
      await new Promise((resolve) => setTimeout(resolve, 800));
      classes = [
        {
          id: "class1",
          name: "Đồ án Phát triển ứng dụng web",
          description: "Lớp học đồ án cuối kỳ HK2 2024-2025",
          semester: "HK2",
          year: 2024,
          invitation_code: "WEB2024",
          status: "active",
          requiresProjectRegistration: true, // Thêm field này
          lecturer: {
            userId: "lecturer1",
            fullName: "TS. Nguyễn Văn A",
            avatar: "",
          },
          students: Array(45).fill(null),
          project_rounds: [
            {
              id: "round1",
              name: "Đợt 1",
              description: "",
              start_date: "2024-02-01",
              end_date: "2024-05-31",
              created_at: "",
              updated_at: "",
            },
          ],
          avatar: "",
          created_at: "2024-01-15",
          updated_at: "2024-01-15",
        },
        {
          id: "class2",
          name: "Đồ án Trí tuệ nhân tạo",
          description: "Áp dụng AI vào bài toán thực tế",
          semester: "HK2",
          year: 2024,
          invitation_code: "AI2024",
          status: "active",
          requiresProjectRegistration: false,
          lecturer: {
            userId: "lecturer2",
            fullName: "TS. Trần Thị B",
            avatar: "",
          },
          students: Array(38).fill(null),
          project_rounds: [
            {
              id: "round2",
              name: "Đợt 1",
              description: "",
              start_date: "2024-02-01",
              end_date: "2024-05-31",
              created_at: "",
              updated_at: "",
            },
            {
              id: "round3",
              name: "Đợt 2",
              description: "",
              start_date: "2024-09-01",
              end_date: "2024-12-31",
              created_at: "",
              updated_at: "",
            },
          ],
          avatar: "",
          created_at: "2024-01-20",
          updated_at: "2024-01-20",
        },
      ];
      */
    } catch (err: any) {
      error = err.message || "Không thể tải danh sách lớp học";
      console.error("Error loading classes:", err);
    } finally {
      loading = false;
    }
  });

  function handleClassClick(classId: string) {
    push(`/classes/${classId}`);
  }

  function getInitials(name: string): string {
    const words = name.split(" ");
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  }
</script>

<div class="container">
  <div class="header">
    <div>
      <h1 class="title">Lớp học của tôi</h1>
      <p class="subtitle">Danh sách các lớp học bạn đã tham gia</p>
    </div>
  </div>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Đang tải danh sách lớp học...</p>
    </div>
  {:else if error}
    <div class="error-state">
      <svg
        width="64"
        height="64"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <h3>{error}</h3>
      <button onclick={() => window.location.reload()} class="btn-retry">
        Thử lại
      </button>
    </div>
  {:else if classes.length === 0}
    <div class="empty-state">
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
      <h3>Chưa có lớp học nào</h3>
      <p>Bạn chưa tham gia lớp học nào</p>
    </div>
  {:else}
    <div class="classes-grid">
      {#each classes as classData (classData.id)}
        <div class="class-card" onclick={() => handleClassClick(classData.id)}>
          <div class="card-header">
            {#if classData.avatar}
              <img
                src={classData.avatar}
                alt={classData.name}
                class="class-avatar"
              />
            {:else}
              <div class="class-avatar-placeholder">
                {getInitials(classData.name)}
              </div>
            {/if}
          </div>

          <div class="card-body">
            <div class="badges-top">
              <span class="badge badge-semester">
                {classData.semester}
                {classData.year}
              </span>
              {#if classData.requiresProjectRegistration}
                <span class="badge badge-warning">
                  <svg
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                  </svg>
                  Cần đăng ký đề tài
                </span>
              {/if}
            </div>

            <h3 class="class-name">{classData.name}</h3>
            {#if classData.description}
              <p class="class-description">{classData.description}</p>
            {/if}
          </div>

          <div class="card-footer">
            <div class="stat-item">
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
              <span>{classData.students?.length || 0} sinh viên</span>
            </div>

            <div class="stat-item">
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
              <span>{classData.project_rounds?.length || 0} đợt đồ án</span>
            </div>

            <div class="stat-item instructor">
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
              <span>GV: {classData.lecturer.fullName}</span>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  .header {
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }

  .title {
    font-size: 1.875rem;
    font-weight: 700;
    color: #111827;
    margin-bottom: 0.5rem;
  }

  .subtitle {
    font-size: 1rem;
    color: #6b7280;
  }

  .loading-state,
  .error-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
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

  .error-state svg,
  .empty-state svg {
    margin-bottom: 1.5rem;
    color: #9ca3af;
  }

  .error-state h3,
  .empty-state h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.5rem;
  }

  .error-state p,
  .empty-state p {
    margin-bottom: 1.5rem;
  }

  .btn-retry {
    padding: 0.625rem 1.5rem;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-retry:hover {
    background: #2563eb;
  }

  .classes-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
  }

  .class-card {
    background: white;
    border-radius: 0.75rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.2s;
  }

  .class-card:hover {
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }

  .card-header {
    height: 140px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
  }

  .class-avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
    border: 4px solid white;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .class-avatar-placeholder {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    font-weight: 700;
    color: #667eea;
    border: 4px solid white;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .card-body {
    padding: 1.5rem;
  }

  .badges-top {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .class-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #111827;
    margin-bottom: 0.5rem;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .class-description {
    font-size: 0.875rem;
    color: #6b7280;
    margin-bottom: 0.5rem;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: 9999px;
  }

  .badge-semester {
    background: #dbeafe;
    color: #1e40af;
  }

  .badge-code {
    background: #f3f4f6;
    color: #4b5563;
  }

  .badge-warning {
    background: #fff7ed;
    color: #c2410c;
    border: 1px solid #fed7aa;
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .badge-warning svg {
    flex-shrink: 0;
  }

  .card-footer {
    padding: 1rem 1.5rem;
    background: #f9fafb;
    border-top: 1px solid #e5e7eb;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    font-size: 0.875rem;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #6b7280;
  }

  .stat-item svg {
    flex-shrink: 0;
  }

  .stat-item.instructor {
    flex-basis: 100%;
    color: #3b82f6;
    font-weight: 500;
  }
</style>
