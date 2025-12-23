<script lang="ts">
  interface ButtonProps {
    variant?: "primary" | "google" | "outline" | "icon" | "share" | "avatar";
    size?: "sm" | "md" /* THÊM: Hỗ trợ size nhỏ và vừa */;
    type?: "button" | "submit" | "reset";
    disabled?: boolean;
    class?: string;
    title?: string;
    onclick?: (e: MouseEvent) => void;
    children?: import("svelte").Snippet;
  }

  let {
    variant = "primary",
    size = "md" /* Mặc định là size to (cho form đăng nhập) */,
    type = "button",
    disabled = false,
    class: className = "",
    title = undefined,
    onclick,
    children,
  }: ButtonProps = $props();

  function handleClick(e: MouseEvent) {
    if (!disabled && onclick) {
      onclick(e);
    }
  }
</script>

<button
  {type}
  {title}
  {disabled}
  class="btn {variant} {size} {className}"
  onclick={handleClick}
>
  {#if variant === "google"}
    <div class="icon-wrapper">
      <svg
        width="18"
        height="18"
        viewBox="0 0 18 18"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M17.64 9.2c0-.637-.057-1.251-.164-1.84H9v3.481h4.844a4.14 4.14 0 0 1-1.796 2.716v2.259h2.908c1.702-1.567 2.684-3.875 2.684-6.615z"
          fill="#4285F4"
        />
        <path
          d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.715H.957v2.332A8.997 8.997 0 0 0 9 18z"
          fill="#34A853"
        />
        <path
          d="M3.964 10.71A5.41 5.41 0 0 1 3.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 0 0 0 9c0 1.452.348 2.827.957 4.042l3.007-2.332z"
          fill="#FBBC05"
        />
        <path
          d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 0 0 .957 4.958L3.964 7.272C4.672 5.14 6.656 3.58 9 3.58z"
          fill="#EA4335"
        />
      </svg>
    </div>
  {/if}

  <span class="content">
    <slot />
  </span>
</button>

<style>
  /* Base Style - Chung cho mọi nút */
  .btn {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    gap: 8px; /* Gap chung */
    font-weight: 600;
    line-height: 1.5;
    border-radius: 8px;
    border: 1px solid transparent;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
    font-family: inherit;
    white-space: nowrap;
    width: auto;
    box-sizing: border-box;
  }

  /* --- SIZES --- Định nghĩa kích thước rõ ràng */

  /* Size Medium (Mặc định - Dùng cho Login/Register Form) */
  .md {
    padding: 12px 24px;
    font-size: 15px;
    height: auto;
  }

  /* Size Small (Mới - Dùng cho Navbar, Action nhỏ) */
  .sm {
    padding: 8px 16px;
    font-size: 14px;
    height: 36px; /* Cố định chiều cao cho gọn */
  }

  /* Base disabled */
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    pointer-events: none;
  }

  /* Variants */
  .primary {
    background-color: #1a56db;
    color: #ffffff;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  }
  .primary:hover {
    background-color: #1e429f;
  }

  .google {
    background-color: #ffffff;
    color: #374151;
    border: 1px solid #d1d5db;
  }
  .google:hover {
    background-color: #f9fafb;
    border-color: #9ca3af;
  }

  .outline {
    background-color: transparent;
    border: 1px solid #d1d5db;
    color: #374151;
  }
  .outline:hover {
    background-color: #f3f4f6;
    border-color: #9ca3af;
  }

  /* Icon button và Avatar thì không cần padding theo size, nó có size riêng */
  .icon {
    background: none;
    padding: 8px;
    border-radius: 50%;
    color: #4b5563;
    width: 36px;
    height: 36px;
    gap: 0;
  }
  .icon:hover {
    background-color: #f3f4f6;
    color: #111827;
  }

  .share {
    background-color: #eff6ff;
    color: #1a56db;
    border: 1px solid #dbeafe;
    padding: 6px 12px;
    font-size: 13px;
  }
  .share:hover {
    background-color: #dbeafe;
  }

  .avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background-color: #db2777;
    color: white;
    font-size: 14px;
    padding: 0;
    gap: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .avatar:hover {
    opacity: 0.9;
  }

  .icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Support for w-full class from parent */
  :global(.btn.w-full) {
    width: 100% !important;
    display: flex !important;
  }
</style>
