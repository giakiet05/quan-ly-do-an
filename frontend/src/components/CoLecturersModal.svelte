<script lang="ts">
  import type { ClassroomResponse } from "../dtos/classroom-dto";

  let { classroom, show = $bindable(false) } = $props<{
    classroom: ClassroomResponse | null;
    show: boolean;
  }>();

  function handleClose() {
    show = false;
  }

  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      handleClose();
    }
  }
</script>

{#if show && classroom && classroom.coLecturers && classroom.coLecturers.length > 0}
  <div class="modal-overlay" onclick={handleOverlayClick}>
    <div class="modal">
      <div class="modal-header">
        <h3>Danh sách trợ giảng</h3>
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
        <div class="lecturers-list">
          {#each classroom.coLecturers as coLecturer}
            <div class="lecturer-item">
              <div class="lecturer-avatar">
                {#if coLecturer.avatar}
                  <img src={coLecturer.avatar} alt={coLecturer.fullName} />
                {:else}
                  <div class="avatar-placeholder">
                    {coLecturer.fullName.charAt(0).toUpperCase()}
                  </div>
                {/if}
              </div>
              <div class="lecturer-info">
                <h4>{coLecturer.fullName}</h4>
                <span class="lecturer-badge">Trợ giảng</span>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <div class="modal-footer">
        <button onclick={handleClose} class="btn-secondary">Đóng</button>
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
    border-radius: 12px;
    width: 100%;
    max-width: 500px;
    max-height: 80vh;
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

  .lecturers-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .lecturer-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: #f9fafb;
    border-radius: 8px;
    transition: background 0.2s;
  }

  .lecturer-item:hover {
    background: #f3f4f6;
  }

  .lecturer-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
  }

  .lecturer-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    font-size: 1.25rem;
    font-weight: 600;
  }

  .lecturer-info {
    flex: 1;
  }

  .lecturer-info h4 {
    font-size: 1rem;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 0.25rem 0;
  }

  .lecturer-badge {
    display: inline-block;
    padding: 2px 8px;
    background: #dbeafe;
    color: #1e40af;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .modal-footer {
    padding: 1.5rem;
    border-top: 1px solid #e5e7eb;
    display: flex;
    justify-content: flex-end;
  }

  .btn-secondary {
    padding: 0.625rem 1.5rem;
    background: #f3f4f6;
    color: #374151;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-secondary:hover {
    background: #e5e7eb;
  }

  @media (max-width: 640px) {
    .modal {
      max-height: 90vh;
    }

    .modal-header,
    .modal-body,
    .modal-footer {
      padding: 1rem;
    }
  }
</style>
