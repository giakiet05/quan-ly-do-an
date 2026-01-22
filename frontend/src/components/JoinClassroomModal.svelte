<script lang="ts">
  import {
    previewClassroom,
    joinClassroom,
  } from "../services/classroom-service";
  import {
    getClassroomErrorMessage,
    CLASSROOM_ERROR_CODES,
  } from "../constants/classroom-error-codes";
  import type { ClassroomPreviewResponse } from "../dtos/classroom-join-dto";

  let { show = $bindable(false), onSuccess } = $props<{
    show: boolean;
    onSuccess?: () => void;
  }>();

  let invitationCode = $state("");
  let preview = $state<ClassroomPreviewResponse | null>(null);
  let loading = $state(false);
  let error = $state("");
  let step = $state<"input" | "preview">("input");

  function handleClose() {
    show = false;
    invitationCode = "";
    preview = null;
    error = "";
    step = "input";
  }

  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      handleClose();
    }
  }

  async function handlePreview() {
    if (!invitationCode.trim()) {
      error = "Vui lòng nhập mã mời";
      return;
    }

    try {
      loading = true;
      error = "";
      console.log("🔍 Preview classroom with code:", invitationCode.trim());
      preview = await previewClassroom(invitationCode.trim());
      console.log("✅ Preview success:", preview);
      step = "preview";
    } catch (err: any) {
      console.error("❌ Preview failed:", err);
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
      console.log("🔗 Join classroom with code:", invitationCode.trim());
      const result = await joinClassroom({
        invitationCode: invitationCode.trim(),
      });
      console.log("✅ Join result:", result);

      if (result.status === "pending") {
        alert(
          "✅ Yêu cầu tham gia đã được gửi!\nVui lòng chờ giảng viên phê duyệt.",
        );
      } else {
        alert("✅ Tham gia lớp học thành công!");
      }

      handleClose();
      onSuccess?.();
    } catch (err: any) {
      console.error("❌ Join failed:", err);
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

{#if show}
  <div class="modal-overlay" onclick={handleOverlayClick}>
    <div class="modal">
      <div class="modal-header">
        <h3>
          {step === "input" ? "Tham gia lớp học" : "Xác nhận tham gia"}
        </h3>
        <button onclick={handleClose} class="btn-close" aria-label="Đóng">
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
        {#if step === "input"}
          <!-- Step 1: Nhập mã mời -->
          <div class="input-section">
            <p class="description">
              Nhập mã mời từ giảng viên để tham gia lớp học
            </p>

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
          </div>
        {:else if preview}
          <!-- Step 2: Preview -->
          <div class="preview-section">
            <div class="preview-header">
              <div class="preview-icon">
                <svg
                  width="48"
                  height="48"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                  <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                </svg>
              </div>
              <h4>{preview.name}</h4>
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
                <h5>Yêu cầu tham gia:</h5>
                <div class="requirements-list">
                  {#if preview.enableWhitelist}
                    <div class="requirement-item">
                      <svg
                        width="16"
                        height="16"
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
                      <span>Chỉ MSSV trong danh sách</span>
                    </div>
                  {/if}
                  {#if preview.enableEmailRestriction && preview.allowedEmailDomains}
                    <div class="requirement-item">
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
                      <span>
                        Chỉ email {preview.allowedEmailDomains.join(", ")}
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
          </div>
        {/if}
      </div>

      <div class="modal-footer">
        {#if step === "input"}
          <button
            onclick={handleClose}
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
        {:else}
          <button onclick={handleBack} class="btn-secondary" disabled={loading}>
            Quay lại
          </button>
          <button
            onclick={handleJoin}
            class="btn-primary"
            disabled={loading ||
              (preview && preview.studentCount >= preview.maxStudents)}
          >
            {#if loading}
              <span class="spinner"></span>
              Đang xử lý...
            {:else}
              Tham gia lớp học
            {/if}
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .modal {
    background: white;
    border-radius: 16px;
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    animation: slideUp 0.3s ease-out;
  }

  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .modal-header h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }

  .btn-close {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #6b7280;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-close:hover {
    background: #f3f4f6;
    color: #374151;
  }

  .modal-body {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
  }

  .description {
    color: #6b7280;
    font-size: 0.875rem;
    margin-bottom: 1.5rem;
    text-align: center;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-group label {
    display: block;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.5rem;
    font-size: 0.875rem;
  }

  .input-code {
    width: 100%;
    padding: 0.75rem;
    font-size: 1rem;
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
    border-color: #2563eb;
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
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

  .preview-section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .preview-header {
    text-align: center;
  }

  .preview-icon {
    display: flex;
    justify-content: center;
    margin-bottom: 1rem;
  }

  .preview-icon svg {
    color: #2563eb;
  }

  .preview-header h4 {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 0.5rem 0;
  }

  .lecturer-name {
    color: #6b7280;
    font-size: 0.875rem;
  }

  .preview-stats {
    padding: 1rem;
    background: #f9fafb;
    border-radius: 8px;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.875rem;
    color: #374151;
  }

  .stat-item svg {
    color: #2563eb;
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

  .requirements-section h5 {
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
    margin: 0 0 0.75rem 0;
  }

  .requirements-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .requirement-item {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    padding: 0.75rem;
    background: #fef3c7;
    border-left: 3px solid #f59e0b;
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
    border-left: 3px solid #ef4444;
    border-radius: 4px;
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

  .modal-footer {
    padding: 1.5rem;
    border-top: 1px solid #e5e7eb;
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
    font-weight: 600;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn-primary {
    background: #2563eb;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #1d4ed8;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
  }

  .btn-primary:disabled {
    background: #9ca3af;
    cursor: not-allowed;
    transform: none;
  }

  .btn-secondary {
    background: #f3f4f6;
    color: #374151;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #e5e7eb;
  }

  .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
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
    .modal {
      max-height: 95vh;
    }

    .modal-header,
    .modal-body,
    .modal-footer {
      padding: 1rem;
    }
  }
</style>
