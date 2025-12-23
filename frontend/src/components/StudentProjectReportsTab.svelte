<script lang="ts">
  import ReportDetailModal from "./ReportDetailModal.svelte";

  interface Report {
    id: string;
    title: string;
    deadline: string;
    status: "completed" | "grading" | "pending" | "late";
    submittedFile?: {
      name: string;
      size: string;
      uploadTime: string;
    };
    grade?: number;
    feedback?: {
      text: string;
      file?: string;
    };
  }

  let { projectName, teamName } = $props<{
    projectName: string;
    teamName: string;
  }>();

  let selectedReport = $state<Report | null>(null);

  function showReportDetail(report: Report) {
    if (report.submittedFile) {
      selectedReport = report;
    }
  }

  let reports = $state<Report[]>([
    {
      id: "1",
      title: "Báo cáo tiến độ Tuần 1",
      deadline: "23/03/2024",
      status: "completed",
      submittedFile: {
        name: "BaoCaoTuan1_NhomA.pdf",
        size: "1.2 MB",
        uploadTime: "10:30, 22/03/2024",
      },
      grade: 8.0,
      feedback: {
        text: "Nhóm làm tốt phần phân tích yêu cầu, tuy nhiên cần chi tiết hơn ở mục thiết kế database.",
        file: "GV_Feedback_T1.docx",
      },
    },
    {
      id: "2",
      title: "Báo cáo tiến độ Tuần 2",
      deadline: "30/03/2024",
      status: "grading",
      submittedFile: {
        name: "BaoCaoTuan2_NhomA.pdf",
        size: "2.5 MB",
        uploadTime: "18:00, 29/03/2024",
      },
    },
    {
      id: "3",
      title: "Báo cáo giữa kỳ",
      deadline: "15/04/2024",
      status: "pending",
    },
    {
      id: "4",
      title: "Báo cáo cuối kỳ",
      deadline: "30/04/2024",
      status: "pending",
    },
  ]);

  const completedReports = $derived(
    reports.filter((r) => r.status === "completed").length
  );
  const totalReports = $derived(reports.length);
  const progressPercentage = $derived((completedReports / totalReports) * 100);
  const averageGrade = $derived(
    reports.filter((r) => r.grade).reduce((sum, r) => sum + (r.grade || 0), 0) /
      reports.filter((r) => r.grade).length || 0
  );

  const nextReport = $derived(reports.find((r) => r.status === "pending"));

  function getStatusConfig(status: Report["status"]) {
    switch (status) {
      case "completed":
        return {
          label: "Đã hoàn thành",
          color: "status-completed",
          icon: "✓",
        };
      case "grading":
        return {
          label: "Đang chấm",
          color: "status-grading",
          icon: "/dangcham_icon.svg",
        };
      case "pending":
        return {
          label: "Chưa nộp",
          color: "status-pending",
          icon: "/chuanop_icon.svg",
        };
      case "late":
        return { label: "Trễ hạn", color: "status-late", icon: "✕" };
    }
  }

  function getIconColor(status: Report["status"]) {
    switch (status) {
      case "completed":
        return "icon-completed";
      case "grading":
        return "icon-grading";
      case "pending":
        return "icon-pending";
      case "late":
        return "icon-late";
    }
  }
</script>

<div class="reports-container">
  <!-- Progress and Average Grade -->
  <div class="stats-grid">
    <!-- Progress Card -->
    <div class="progress-card">
      <h3 class="card-title">Tiến độ đồ án</h3>
      <p class="card-subtitle">
        Hoàn thành {completedReports}/{totalReports} mốc báo cáo quan trọng.
      </p>
      <div class="progress-wrapper">
        <div class="progress-labels">
          <span class="progress-current">{progressPercentage.toFixed(0)}%</span>
          <span class="progress-max">100%</span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill" style="width: {progressPercentage}%"></div>
        </div>
      </div>
    </div>

    <!-- Average Grade Card -->
    <div class="grade-card">
      <h3 class="card-title">Điểm trung bình tạm tính</h3>
      <div class="grade-value">{averageGrade.toFixed(1)}</div>
      <div class="grade-max">/ 10</div>
    </div>
  </div>

  <!-- Next Deadline Countdown -->
  {#if nextReport}
    <div class="deadline-card">
      <h3 class="deadline-title">Hạn chót tiếp theo: {nextReport.title}</h3>
      <p class="deadline-date">Thời hạn: 23:59, {nextReport.deadline}</p>

      <div class="countdown-wrapper">
        <div class="countdown">
          <div class="countdown-item">
            <div class="countdown-value">08</div>
            <div class="countdown-label">Ngày</div>
          </div>
          <div class="countdown-separator">:</div>
          <div class="countdown-item">
            <div class="countdown-value">15</div>
            <div class="countdown-label">Giờ</div>
          </div>
          <div class="countdown-separator">:</div>
          <div class="countdown-item">
            <div class="countdown-value">30</div>
            <div class="countdown-label">Phút</div>
          </div>
        </div>

        <button class="submit-button">
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="17 8 12 3 7 8"></polyline>
            <line x1="12" y1="3" x2="12" y2="15"></line>
          </svg>
          Nộp bài ngay
        </button>
      </div>
    </div>
  {/if}

  <!-- Timeline -->
  <div class="timeline">
    {#each reports as report}
      {@const statusConfig = getStatusConfig(report.status)}
      {@const iconColor = getIconColor(report.status)}
      <div class="timeline-item">
        <!-- Timeline icon -->
        <div class="timeline-icon {iconColor}">
          {#if report.status === "grading"}
            <img
              src="/dangcham_icon.svg"
              alt="Đang chấm"
              width="20"
              height="20"
            />
          {:else if report.status === "pending"}
            <img
              src="/chuanop_icon.svg"
              alt="Chưa nộp"
              width="20"
              height="20"
            />
          {:else}
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
              ></path>
              <polyline points="14 2 14 8 20 8"></polyline>
            </svg>
          {/if}
        </div>

        <!-- Report Card -->
        <div class="report-card" class:highlight={report.status === "pending"}>
          <!-- Header -->
          <div class="report-header">
            <div class="report-info">
              <h4 class="report-title">{report.title}</h4>
              <p class="report-deadline">Hạn chót: {report.deadline}</p>
            </div>
            <span class="status-badge {statusConfig.color}">
              {#if statusConfig.icon.endsWith(".svg")}
                <img src={statusConfig.icon} alt="" width="14" height="14" />
              {:else}
                <span>{statusConfig.icon}</span>
              {/if}
              {statusConfig.label}
            </span>
          </div>

          <!-- Content -->
          <div class="report-content">
            {#if report.status === "pending"}
              <div class="upload-section">
                <h5 class="section-title">Nộp bài của bạn</h5>
                <div class="upload-area">
                  <svg
                    width="44"
                    height="44"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                    <polyline points="17 8 12 3 7 8"></polyline>
                    <line x1="12" y1="3" x2="12" y2="15"></line>
                  </svg>
                  <p class="upload-text">Kéo và thả file vào đây</p>
                  <p class="upload-or">hoặc</p>
                  <button class="upload-button">Chọn file để tải lên</button>
                  <p class="upload-hint">
                    Hỗ trợ: PDF, DOCX, ZIP (tối đa 25MB)
                  </p>
                </div>
              </div>
            {:else}
              <div class="submitted-section">
                <h5 class="section-title">Bài đã nộp</h5>
                {#if report.submittedFile}
                  <div class="file-card">
                    <svg
                      width="24"
                      height="24"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      class="file-icon"
                    >
                      <path
                        d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                      ></path>
                      <polyline points="14 2 14 8 20 8"></polyline>
                    </svg>
                    <div class="file-info">
                      <div class="file-name">{report.submittedFile.name}</div>
                      <div class="file-meta">
                        {report.submittedFile.size} - Nộp lúc {report
                          .submittedFile.uploadTime}
                      </div>
                    </div>
                    {#if report.status === "grading"}
                      <button class="cancel-button">
                        <svg
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <line x1="18" y1="6" x2="6" y2="18"></line>
                          <line x1="6" y1="6" x2="18" y2="18"></line>
                        </svg>
                        Hủy nộp
                      </button>
                    {:else}
                      <button class="download-button">
                        <svg
                          width="20"
                          height="20"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                          ></path>
                          <polyline points="7 10 12 15 17 10"></polyline>
                          <line x1="12" y1="15" x2="12" y2="3"></line>
                        </svg>
                      </button>
                    {/if}
                  </div>
                {/if}

                <!-- Feedback section -->
                {#if report.status === "completed" && report.grade}
                  <div class="feedback-section">
                    <div class="feedback-header">
                      <h5 class="section-title">Phản hồi của giảng viên</h5>
                      <button
                        onclick={() => showReportDetail(report)}
                        class="detail-button"
                      >
                        <svg
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path
                            d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                          ></path>
                        </svg>
                        Xem chi tiết & Thảo luận
                      </button>
                    </div>
                    <div class="grade-display">{report.grade.toFixed(1)}</div>
                    {#if report.feedback}
                      <p class="feedback-text">{report.feedback.text}</p>
                      {#if report.feedback.file}
                        <button class="feedback-file">
                          <svg
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                            ></path>
                            <polyline points="7 10 12 15 17 10"></polyline>
                            <line x1="12" y1="15" x2="12" y2="3"></line>
                          </svg>
                          {report.feedback.file}
                        </button>
                      {/if}
                    {/if}
                  </div>
                {/if}
              </div>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

{#if selectedReport}
  <ReportDetailModal
    reportTitle={selectedReport.title}
    submittedFile={selectedReport.submittedFile}
    grade={selectedReport.grade}
    feedback={selectedReport.feedback}
    onClose={() => (selectedReport = null)}
  />
{/if}

<style>
  .reports-container {
    max-width: 1000px;
  }

  /* Stats Grid */
  .stats-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 24px;
    margin-bottom: 24px;
  }

  .progress-card,
  .grade-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    padding: 24px;
  }

  .card-title {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .card-subtitle {
    font-size: 14px;
    color: #64748b;
    margin-bottom: 16px;
  }

  .progress-wrapper {
  }

  .progress-labels {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    font-size: 14px;
  }

  .progress-current {
    color: #1e293b;
    font-weight: 500;
  }

  .progress-max {
    color: #94a3b8;
  }

  .progress-bar {
    height: 10px;
    background: #e2e8f0;
    border-radius: 999px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: #3b82f6;
    border-radius: 999px;
    transition: width 0.5s ease;
  }

  .grade-card {
    text-align: center;
  }

  .grade-value {
    font-size: 36px;
    color: #3b82f6;
    font-weight: 600;
    margin: 12px 0 4px 0;
  }

  .grade-max {
    font-size: 14px;
    color: #64748b;
  }

  /* Deadline Card */
  .deadline-card {
    background: rgba(59, 130, 246, 0.1);
    border: 1px solid rgba(59, 130, 246, 0.2);
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 24px;
  }

  .deadline-title {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .deadline-date {
    font-size: 14px;
    color: #475569;
    margin-bottom: 16px;
  }

  .countdown-wrapper {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .countdown {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .countdown-item {
    text-align: center;
  }

  .countdown-value {
    font-size: 32px;
    font-weight: 600;
    color: #1e293b;
  }

  .countdown-label {
    font-size: 12px;
    color: #1e293b;
  }

  .countdown-separator {
    font-size: 32px;
    font-weight: 300;
    color: #1e293b;
  }

  .submit-button {
    margin-left: auto;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 24px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;
  }

  .submit-button:hover {
    background: #2563eb;
  }

  /* Timeline */
  .timeline {
    position: relative;
    padding-left: 64px;
  }

  .timeline::before {
    content: "";
    position: absolute;
    left: 20px;
    top: 0;
    bottom: 0;
    width: 2px;
    background: #e2e8f0;
  }

  .timeline-item {
    position: relative;
    margin-bottom: 32px;
  }

  .timeline-item:last-child {
    margin-bottom: 0;
  }

  .timeline-icon {
    position: absolute;
    left: -44px;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 0 0 8px #f8fafc;
  }

  .timeline-icon svg {
    color: white;
  }

  .timeline-icon img {
    width: 20px;
    height: 20px;
  }

  .icon-completed {
    background: #22c55e;
  }

  .icon-grading {
    background: #3b82f6;
  }

  .icon-pending {
    background: #3b82f6;
  }

  .icon-late {
    background: #ef4444;
  }

  /* Report Card */
  .report-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  .report-card.highlight {
    border: 2px solid #3b82f6;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .report-header {
    padding: 20px;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
  }

  .report-info {
  }

  .report-title {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 4px;
    color: #1e293b;
  }

  .report-deadline {
    font-size: 14px;
    color: #64748b;
  }

  .status-badge {
    padding: 4px 12px;
    border-radius: 16px;
    font-size: 12px;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .status-badge img {
    width: 14px;
    height: 14px;
    flex-shrink: 0;
  }

  .status-completed {
    background: #dcfce7;
    color: #15803d;
  }

  .status-grading {
    background: #dbeafe;
    color: #1e40af;
  }

  .status-pending {
    background: #fef3c7;
    color: #a16207;
  }

  .status-late {
    background: #fee2e2;
    color: #991b1b;
  }

  /* Report Content */
  .report-content {
    padding: 20px;
  }

  .section-title {
    font-size: 14px;
    color: #475569;
    margin-bottom: 12px;
  }

  /* Upload Area */
  .upload-section {
  }

  .upload-area {
    border: 2px dashed #cbd5e1;
    border-radius: 8px;
    padding: 32px;
    text-align: center;
  }

  .upload-area svg {
    color: #94a3b8;
    margin: 0 auto 12px;
  }

  .upload-text {
    font-size: 14px;
    color: #334155;
    margin-bottom: 4px;
  }

  .upload-or {
    font-size: 12px;
    color: #64748b;
    margin-bottom: 8px;
  }

  .upload-button {
    font-size: 14px;
    color: #3b82f6;
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }

  .upload-hint {
    font-size: 12px;
    color: #94a3b8;
    margin-top: 12px;
  }

  /* Submitted Section */
  .submitted-section {
  }

  .file-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .file-icon {
    color: #3b82f6;
    flex-shrink: 0;
  }

  .file-info {
    flex: 1;
  }

  .file-name {
    font-size: 14px;
    color: #1e293b;
    margin-bottom: 2px;
  }

  .file-meta {
    font-size: 12px;
    color: #64748b;
  }

  .cancel-button {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 14px;
    color: #64748b;
    background: none;
    border: none;
    cursor: pointer;
    transition: color 0.2s;
  }

  .cancel-button:hover {
    color: #111827;
  }

  .download-button {
    padding: 8px;
    color: #64748b;
    background: none;
    border: none;
    cursor: pointer;
    transition: color 0.2s;
  }

  .download-button:hover {
    color: #334155;
  }

  /* Feedback Section */
  .feedback-section {
    border-left: 2px solid #e2e8f0;
    padding-left: 24px;
  }

  .feedback-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }

  .detail-button {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 14px;
    color: #3b82f6;
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }

  .grade-display {
    font-size: 36px;
    color: #3b82f6;
    font-weight: 600;
    margin-bottom: 12px;
  }

  .feedback-text {
    font-size: 14px;
    color: #334155;
    margin-bottom: 12px;
  }

  .feedback-file {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #3b82f6;
    background: none;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }
</style>
