<script lang="ts">
  import { push } from "svelte-spa-router";
  import {
    previewClassroom,
    joinClassroom,
  } from "../../services/classroom-service";
  import {
    getClassroomErrorMessage,
    CLASSROOM_ERROR_CODES,
  } from "../../constants/classroom-error-codes";
  import ClassroomRequirements from "../../components/ClassroomRequirements.svelte";
  import type { ClassroomPreviewResponse } from "../../dtos/classroom-join-dto";

  let invitationCode = $state("");
  let preview = $state<ClassroomPreviewResponse | null>(null);
  let loading = $state(false);
  let error = $state("");
  let step = $state<"input" | "preview">("input");

  async function handlePreview() {
    if (!invitationCode.trim()) {
      error = "Vui lòng nhập mã mời";
      return;
    }

    try {
      loading = true;
      error = "";
      preview = await previewClassroom(invitationCode.trim());
      step = "preview";
    } catch (err: any) {
      error = getClassroomErrorMessage(
        err.error_code || CLASSROOM_ERROR_CODES.CODE_INVALID,
      );
      preview = null;
    } finally {
      loading = false;
    }
  }

  async function handleJoin() {
    try {
      loading = true;
      error = "";
      const result = await joinClassroom({
        invitationCode: invitationCode.trim(),
      });

      if (result.status === "pending") {
        alert(
          "✅ Yêu cầu tham gia đã được gửi!\nVui lòng chờ giảng viên phê duyệt.",
        );
      } else {
        alert("✅ Tham gia lớp học thành công!");
      }

      // Chuyển về danh sách lớp
      push("/student/classes");
    } catch (err: any) {
      error = getClassroomErrorMessage(err.error_code);
    } finally {
      loading = false;
    }
  }

  function handleBack() {
    step = "input";
    preview = null;
    error = "";
  }
</script>

<div class="join-classroom-page">
  <div class="container">
    {#if step === "input"}
      <!-- Step 1: Nhập mã mời -->
      <div class="input-card">
        <div class="card-header">
          <svg
            width="48"
            height="48"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M13.8 12H3"
            ></path>
          </svg>
          <h1>Tham gia lớp học</h1>
          <p>Nhập mã mời từ giảng viên để tham gia lớp học</p>
        </div>

        <div class="form-group">
          <label for="code">Mã mời lớp học</label>
          <input
            id="code"
            type="text"
            bind:value={invitationCode}
            placeholder="VD: ABC123XYZ"
            class="input-code"
            disabled={loading}
            onkeydown={(e) => {
              if (e.key === "Enter") handlePreview();
            }}
          />
          {#if error}
            <p class="error-message">{error}</p>
          {/if}
        </div>

        <div class="button-group">
          <button
            onclick={() => push("/student/classes")}
            class="btn-secondary"
            disabled={loading}
          >
            Hủy
          </button>
          <button
            onclick={handlePreview}
            class="btn-primary"
            disabled={loading || !invitationCode.trim()}
          >
            {#if loading}
              <span class="spinner"></span>
              Đang kiểm tra...
            {:else}
              Tiếp tục
            {/if}
          </button>
        </div>
      </div>
    {:else}
      <!-- Step 2: Preview và Join -->
      <div class="preview-card">
        <button onclick={handleBack} class="btn-back">
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
          Quay lại
        </button>

        {#if preview}
          <div class="preview-header">
            <div class="preview-icon">
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
            </div>
            <h2>{preview.name}</h2>
            <p class="lecturer-name">Giảng viên: {preview.lecturer}</p>
          </div>

          <div class="preview-stats">
            <div class="stat-item">
              <svg
                width="20"
                height="20"
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
              <span>
                {preview.studentCount} / {preview.maxStudents} sinh viên
              </span>
              {#if preview.studentCount >= preview.maxStudents}
                <span class="badge-full">Đã đầy</span>
              {/if}
            </div>
          </div>

          {#if preview.enableWhitelist || preview.enableEmailRestriction}
            <div class="requirements-section">
              <h3>Yêu cầu tham gia:</h3>
              <div class="requirements-list">
                {#if preview.enableWhitelist}
                  <div class="requirement-item">
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path d="M9 11l3 3L22 4"></path>
                      <path
                        d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"
                      ></path>
                    </svg>
                    <span
                      >Chỉ mã sinh viên trong danh sách mới có thể tham gia</span
                    >
                  </div>
                {/if}
                {#if preview.enableEmailRestriction && preview.allowedEmailDomains}
                  <div class="requirement-item">
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
                    <span>
                      Chỉ email {preview.allowedEmailDomains.join(", ")} mới có thể
                      tham gia
                    </span>
                  </div>
                {/if}
              </div>
            </div>
          {/if}

          {#if error}
            <div class="error-box">
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <p>{error}</p>
            </div>
          {/if}

          <div class="button-group">
            <button
              onclick={handleBack}
              class="btn-secondary"
              disabled={loading}
            >
              Quay lại
            </button>
            <button
              onclick={handleJoin}
              class="btn-primary"
              disabled={loading || preview.studentCount >= preview.maxStudents}
            >
              {#if loading}
                <span class="spinner"></span>
                Đang xử lý...
              {:else}
                Tham gia lớp học
              {/if}
            </button>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .join-classroom-page {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .container {
    width: 100%;
    max-width: 600px;
  }

  .input-card,
  .preview-card {
    background: white;
    border-radius: 16px;
    padding: 3rem;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .card-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .card-header svg {
    color: #667eea;
    margin-bottom: 1rem;
  }

  .card-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 0.5rem;
  }

  .card-header p {
    color: #6b7280;
    font-size: 1rem;
  }

  .form-group {
    margin-bottom: 2rem;
  }

  .form-group label {
    display: block;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.5rem;
  }

  .input-code {
    width: 100%;
    padding: 1rem;
    font-size: 1.125rem;
    border: 2px solid #e5e7eb;
    border-radius: 8px;
    text-align: center;
    text-transform: uppercase;
    letter-spacing: 2px;
    font-weight: 600;
    transition: all 0.2s;
  }

  .input-code:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .input-code:disabled {
    background: #f9fafb;
    cursor: not-allowed;
  }

  .error-message {
    color: #ef4444;
    font-size: 0.875rem;
    margin-top: 0.5rem;
  }

  .button-group {
    display: flex;
    gap: 1rem;
  }

  .btn-primary,
  .btn-secondary {
    flex: 1;
    padding: 1rem;
    font-size: 1rem;
    font-weight: 600;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }

  .btn-primary {
    background: #667eea;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #5568d3;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }

  .btn-primary:disabled {
    background: #9ca3af;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: #f3f4f6;
    color: #374151;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #e5e7eb;
  }

  .btn-back {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: transparent;
    border: none;
    color: #6b7280;
    font-size: 0.875rem;
    cursor: pointer;
    margin-bottom: 1.5rem;
    transition: color 0.2s;
  }

  .btn-back:hover {
    color: #374151;
  }

  .preview-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .preview-icon {
    display: flex;
    justify-content: center;
    margin-bottom: 1rem;
  }

  .preview-icon svg {
    color: #667eea;
  }

  .preview-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 0.5rem;
  }

  .lecturer-name {
    color: #6b7280;
    font-size: 1rem;
  }

  .preview-stats {
    padding: 1.5rem;
    background: #f9fafb;
    border-radius: 8px;
    margin-bottom: 1.5rem;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 1rem;
    color: #374151;
  }

  .stat-item svg {
    color: #667eea;
    flex-shrink: 0;
  }

  .badge-full {
    display: inline-block;
    padding: 2px 8px;
    background: #ef4444;
    color: white;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 600;
    margin-left: auto;
  }

  .requirements-section {
    margin-bottom: 1.5rem;
  }

  .requirements-section h3 {
    font-size: 1rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 1rem;
  }

  .requirements-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .requirement-item {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 1rem;
    background: #fef3c7;
    border-left: 4px solid #f59e0b;
    border-radius: 4px;
  }

  .requirement-item svg {
    color: #f59e0b;
    flex-shrink: 0;
    margin-top: 2px;
  }

  .requirement-item span {
    color: #78350f;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  .error-box {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 1rem;
    background: #fee2e2;
    border-left: 4px solid #ef4444;
    border-radius: 4px;
    margin-bottom: 1.5rem;
  }

  .error-box svg {
    color: #ef4444;
    flex-shrink: 0;
  }

  .error-box p {
    color: #7f1d1d;
    font-size: 0.875rem;
    margin: 0;
  }

  .spinner {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 0.6s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  @media (max-width: 640px) {
    .join-classroom-page {
      padding: 1rem;
    }

    .input-card,
    .preview-card {
      padding: 2rem;
    }

    .card-header h1 {
      font-size: 1.5rem;
    }

    .button-group {
      flex-direction: column;
    }
  }
</style>
