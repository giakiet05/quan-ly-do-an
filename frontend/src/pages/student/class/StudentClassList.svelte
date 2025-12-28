<script lang="ts">
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";

  interface ClassData {
    id: string;
    code: string;
    name: string;
    semester: string;
    projectCount: number;
    studentCount: number;
    instructor: string;
    university: string;
    description: string;
    requiresProjectRegistration: boolean;
  }

  // Mock data - sinh viên chỉ thấy lớp mình tham gia
  const enrolledClasses: ClassData[] = [
    {
      id: "1",
      code: "IT4441",
      name: "Phát triển ứng dụng Web",
      semester: "HK1 2024-2025",
      projectCount: 15,
      studentCount: 45,
      instructor: "TS. Nguyễn Văn A",
      university: "Đại học Bách Khoa Hà Nội",
      description: "Môn học về phát triển ứng dụng web hiện đại",
      requiresProjectRegistration: true,
    },
    {
      id: "2",
      code: "IT4408",
      name: "An toàn và bảo mật thông tin",
      semester: "HK1 2024-2025",
      projectCount: 12,
      studentCount: 58,
      instructor: "TS. Trần Thị B",
      university: "Đại học Bách Khoa Hà Nội",
      description: "Môn học về an toàn thông tin và bảo mật mạng",
      requiresProjectRegistration: false,
    },
    {
      id: "3",
      code: "IT4992",
      name: "Đồ án tốt nghiệp",
      semester: "HK1 2024-2025",
      projectCount: 21,
      studentCount: 15,
      instructor: "TS. Lê Văn C",
      university: "Đại học Bách Khoa Hà Nội",
      description: "Đồ án tốt nghiệp cho sinh viên năm cuối",
      requiresProjectRegistration: true,
    },
  ];

  function handleSelectClass(classData: ClassData): void {
    push(`/classes/${classData.id}`);
  }
</script>

<div class="page-container">
  <div class="header-section">
    <div>
      <h1 class="page-title">Lớp học của tôi</h1>
      <p class="page-subtitle">Danh sách các lớp học bạn đang tham gia</p>
    </div>
  </div>

  <div class="class-grid">
    {#each enrolledClasses as classData (classData.id)}
      <div
        class="class-card"
        onclick={() => handleSelectClass(classData)}
        role="button"
        tabindex="0"
        onkeydown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            handleSelectClass(classData);
          }
        }}
      >
        <div class="card-header">
          <div class="header-content">
            <div class="badge-group">
              <span class="code-badge">{classData.code}</span>
              <span class="semester-text">{classData.semester}</span>
              {#if classData.requiresProjectRegistration}
                <span class="registration-badge">
                  <img
                    src="/alert-circle.svg"
                    alt="alert"
                    width="14"
                    height="14"
                  />
                  Cần đăng ký đề tài
                </span>
              {/if}
            </div>
            <h3 class="class-name">{classData.name}</h3>
            <p class="class-description">{classData.description}</p>
          </div>
        </div>

        <div class="card-footer">
          <div class="info-item">
            <img src="/users.svg" alt="users" width="16" height="16" />
            <span>{classData.studentCount} sinh viên</span>
          </div>
          <div class="info-item">
            <img
              src="/folder-kanban.svg"
              alt="projects"
              width="16"
              height="16"
            />
            <span>{classData.projectCount} đề tài</span>
          </div>
          <div class="info-item">
            <img
              src="/book-open-text.svg"
              alt="teacher"
              width="16"
              height="16"
            />
            <span>GV: {classData.instructor}</span>
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .page-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 8px;
  }

  .page-subtitle {
    color: #64748b;
    font-size: 0.875rem;
  }

  .class-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .class-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 24px;
    cursor: pointer;
    transition: all 0.2s;
    border: 1px solid #e2e8f0;
  }

  .class-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    border-color: #3b82f6;
  }

  .class-card:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 2px;
  }

  .card-header {
    margin-bottom: 16px;
  }

  .header-content {
    flex: 1;
  }

  .badge-group {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .code-badge {
    padding: 4px 12px;
    background: #dbeafe;
    color: #1e40af;
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .semester-text {
    color: #64748b;
    font-size: 0.875rem;
  }

  .registration-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 4px 12px;
    background: #fff7ed;
    color: #c2410c;
    border-radius: 4px;
    font-size: 0.875rem;
    border: 1px solid #fed7aa;
  }

  .class-name {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 8px;
  }

  .class-description {
    color: #64748b;
    font-size: 0.875rem;
    margin-bottom: 12px;
    line-height: 1.5;
  }

  .card-footer {
    display: flex;
    align-items: center;
    gap: 24px;
    flex-wrap: wrap;
  }

  .info-item {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 0.875rem;
  }

  .info-item img {
    opacity: 0.7;
  }

  @media (max-width: 768px) {
    .page-container {
      padding: 16px;
    }

    .card-footer {
      gap: 16px;
    }
  }
</style>
