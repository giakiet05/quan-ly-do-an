<script lang="ts">
  import { location, push } from "svelte-spa-router";
  import { authStore, clearAuth } from "../stores/auth-store";
  import Button from "../components/Button.svelte";

  // --- Logic Breadcrumbs giữ nguyên ---
  const routeNameMap: Record<string, string> = {
    "/": "Tổng quan",
    "/settings": "Cài đặt",
    "/profile": "Hồ sơ cá nhân",
    "/projects": "Đồ án",
    "/login": "Đăng nhập",
    "/register": "Đăng ký"
  };

  let currentPathName = $derived(routeNameMap[$location] || "Trang hiện tại");

  // --- Logic Auth giữ nguyên ---
  let user = $derived($authStore.user);
  let isAuthenticated = $derived($authStore.isAuthenticated);

  let avatarLetter = $derived(
    user?.username?.charAt(0).toUpperCase() || user?.email?.charAt(0).toUpperCase() || "?"
  );

  function handleLogout() {
    if (confirm("Bạn muốn đăng xuất?")) {
      clearAuth();
      push("/login");
    }
  }

  function goToLogin() {
    push("/login");
  }
</script>

<header class="topbar">
  <div class="breadcrumbs">
    <span class="breadcrumb-root hidden-mobile">Trang chủ</span>
    <span class="separator hidden-mobile">/</span>
    <h1 class="breadcrumb-current">{currentPathName}</h1>
  </div>

  <div class="actions">
    {#if isAuthenticated}
      <Button
        variant="share"
        class="hidden-xs"
        title="Chia sẻ trang này"
        onclick={() => console.log('Share')}
      >
        <span>Share</span>
      </Button>

      <Button variant="icon" title="Thông báo">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
          <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
        </svg>
      </Button>

      <div class="user-profile">
        <Button variant="avatar" onclick={handleLogout} title="Click để đăng xuất">
          {avatarLetter}
        </Button>
      </div>

    {:else}
      <Button variant="primary" size="sm" onclick={goToLogin}>
        Đăng nhập
      </Button>
    {/if}
  </div>
</header>

<style>
  /* --- Layout chung --- */
  .topbar {
    height: 64px;
    background-color: #ffffff;
    border-bottom: 1px solid #e5e7eb;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    position: sticky;
    top: 0;
    z-index: 40;
    box-sizing: border-box;
  }

  /* --- Breadcrumbs --- */
  .breadcrumbs {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 15px;
    white-space: nowrap;
    color: #6b7280;
  }

  .breadcrumb-root { 
    font-weight: 500;
    cursor: default;
  }
  
  .separator { 
    color: #d1d5db; 
    font-size: 16px; 
  }

  .breadcrumb-current { 
    font-size: 18px;
    font-weight: 700; 
    color: #111827; 
    margin: 0;
    line-height: 1.2;
  }

  /* --- Actions (Bên phải) --- */
  .actions { display: flex; align-items: center; gap: 12px; }

  /* ĐÃ XÓA: Đoạn :global(.nav-btn) vì không còn cần thiết nữa */

  /* --- RESPONSIVE --- */
  @media (max-width: 768px) {
    .topbar { padding: 0 16px; height: 60px; }
    .hidden-mobile { display: none; }
    .breadcrumb-current { font-size: 16px; } 
  }

  @media (max-width: 480px) {
    .hidden-xs { display: none; }
    .breadcrumbs { gap: 8px; }
  }
</style>
