<script lang="ts">
  import type { ClassroomResponse } from "../dtos/classroom-dto";

  let { classroom } = $props<{ classroom: ClassroomResponse }>();
</script>

<div class="classroom-requirements">
  {#if classroom.enableWhitelist && classroom.whitelistStudentCode && classroom.whitelistStudentCode.length > 0}
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
        <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"
        ></path>
      </svg>
      <span>Chỉ MSSV trong danh sách mới có thể tham gia</span>
    </div>
  {/if}

  {#if classroom.enableEmailRestriction && classroom.allowedEmailDomains && classroom.allowedEmailDomains.length > 0}
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
        Chỉ email {classroom.allowedEmailDomains.join(", ")} mới có thể tham gia
      </span>
    </div>
  {/if}

  {#if classroom.coLecturers && classroom.coLecturers.length > 0}
    <div class="requirement-item">
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
      <span>
        Trợ giảng: {classroom.coLecturers
          .map((c: any) => c.fullName)
          .join(", ")}
      </span>
    </div>
  {/if}
</div>

<style>
  .classroom-requirements {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background: #f8fafc;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
  }

  .requirement-item {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    font-size: 14px;
    color: #475569;
  }

  .requirement-item svg {
    flex-shrink: 0;
    margin-top: 2px;
    color: #3b82f6;
  }

  .requirement-item span {
    line-height: 1.5;
  }
</style>
