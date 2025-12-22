<script lang="ts">
  import { location, push } from "svelte-spa-router";
  import { authStore } from "../stores/auth-store";
  import SidebarNav from "../components/SidebarNav.svelte";

  let { children } = $props();
  let isCollapsed = $state(false);
  let currentRole = $state<"LECTURER" | "STUDENT">("LECTURER");
  let isLoadingRole = $state(false);

  function toggleSidebar() {
    isCollapsed = !isCollapsed;
  }

  function goToMyProfile() {
    push("/my-profile");
  }

  function handleLogout() {
    if (confirm("Bạn có chắc muốn đăng xuất?")) {
      authStore.set({
        user: null,
        accessToken: null,
        refreshToken: null,
        isAuthenticated: false,
      });
      localStorage.removeItem("auth");
      push("/auth/login");
    }
  }

  function switchRole(role: "LECTURER" | "STUDENT") {
    isLoadingRole = true;
    currentRole = role;
    setTimeout(() => {
      push("/home");
      setTimeout(() => {
        isLoadingRole = false;
      }, 300);
    }, 200);
  }

  const currentBreadcrumb = $derived("Trang chủ");
</script>

<div class="main-layout">
  {#if isLoadingRole}
    <div class="loading-overlay">
      <div class="loading-spinner"></div>
      <p class="loading-text">Đang chuyển đổi vai trò...</p>
    </div>
  {/if}

  <SidebarNav
    bind:collapsed={isCollapsed}
    role={currentRole}
    onToggle={toggleSidebar}
    onProfileClick={goToMyProfile}
    onLogout={handleLogout}
  />

  <main class="content">
    <header class="top-bar">
      <div class="header-left">
        <div class="breadcrumb">
          <span class="breadcrumb-item">Trang chủ</span>
          <span class="separator">/</span>
          <span class="breadcrumb-item active">{currentBreadcrumb}</span>
        </div>
      </div>

      <!-- Role Switcher -->
      <div class="role-switcher">
        <button
          onclick={() => switchRole("LECTURER")}
          class="role-btn {currentRole === 'LECTURER' ? 'active' : ''}"
        >
          <img src="/teacher_role.svg" alt="Teacher" width="16" height="16" />
          <span>Giáo viên</span>
        </button>
        <button
          onclick={() => switchRole("STUDENT")}
          class="role-btn {currentRole === 'STUDENT' ? 'active' : ''}"
        >
          <img src="/user.svg" alt="Student" width="16" height="16" />
          <span>Sinh viên</span>
        </button>
      </div>
    </header>
    <div class="page-body">
      {@render children()}
    </div>
  </main>
</div>

<style>
  .main-layout {
    display: flex;
    height: 100vh;
    background-color: #f1f5f9;
    position: relative;
  }

  .loading-overlay {
    position: fixed;
    inset: 0;
    background: rgba(15, 23, 42, 0.7);
    backdrop-filter: blur(4px);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    z-index: 9999;
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

  .loading-spinner {
    width: 48px;
    height: 48px;
    border: 4px solid rgba(255, 255, 255, 0.2);
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .loading-text {
    margin-top: 16px;
    color: white;
    font-size: 16px;
    font-weight: 500;
    letter-spacing: 0.5px;
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow-x: hidden;
  }

  .top-bar {
    background: white;
    height: 70px;
    padding: 0 30px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #e2e8f0;
  }

  .role-switcher {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 4px;
    display: inline-flex;
    gap: 4px;
  }

  .role-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    border-radius: 6px;
    border: none;
    background: transparent;
    cursor: pointer;
    font-size: 14px;
    color: #374151;
    transition: all 0.2s;
  }

  .role-btn:hover {
    background: #f3f4f6;
  }

  .role-btn.active {
    background: #2563eb;
    color: white;
  }

  .role-btn.active img {
    filter: brightness(0) invert(1);
  }

  .role-btn:not(.active) img {
    filter: brightness(0) saturate(100%) invert(22%) sepia(8%) saturate(1234%)
      hue-rotate(176deg) brightness(95%) contrast(89%);
  }

  .page-body {
    flex: 1;
    padding: 30px;
    overflow-y: auto;
  }

  .breadcrumb {
    display: flex;
    align-items: center;
    font-size: 0.95rem;
    color: #64748b;
  }
  .breadcrumb-item.active {
    color: #1e293b;
    font-weight: 600;
  }
  .separator {
    margin: 0 10px;
    color: #cbd5e1;
  }
</style>
