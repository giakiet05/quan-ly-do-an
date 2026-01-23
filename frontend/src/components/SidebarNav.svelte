<script lang="ts">
  import { link, location } from "svelte-spa-router";
  import { push } from "svelte-spa-router";
  import { authStore, clearAuth } from "../stores/auth-store";
  import {
    unreadNotificationCount,
    notificationStore,
  } from "../stores/notification-store";
  import { onMount } from "svelte"; // Thêm import này
  import { GraduationCap } from "lucide-svelte";

  type MenuItem = {
    id: string;
    label: string;
    route: string;
    icon: string;
    badge?: number;
  };

  type SidebarProps = {
    collapsed?: boolean;
    role: "LECTURER" | "STUDENT";
    onToggle?: () => void;
    onProfileClick?: () => void;
    onLogout?: () => void;
  };

  let {
    collapsed = $bindable(false),
    role,
    onToggle,
    onProfileClick,
    onLogout,
  }: SidebarProps = $props();

  // Danh sách menu gốc
  const teacherMenuItems: MenuItem[] = [
    // {
    //   id: "dashboard",
    //   label: "Tổng quan",
    //   route: "/dashboard",
    //   icon: "/layout-dashboard.svg",
    // },
    {
      id: "classes",
      label: "Lớp học",
      route: "/lecture/my-classes",
      icon: "/book-open-text.svg",
    },
    {
      id: "notifications",
      label: "Thông báo",
      route: "/notifications",
      icon: "/bell.svg",
    },
    // ❌ HIDDEN FOR TEST - TẠM ẨN ĐỂ KIỂM TRA FRONTEND CẬP NHẬT ❌
    // {
    //   id: "chats",
    //   label: "Tin nhắn",
    //   route: "/chats",
    //   icon: "/message-circle.svg",
    // },
  ];

  const studentMenuItems: MenuItem[] = [
    {
      id: "dashboard",
      label: "Tổng quan",
      route: "/dashboard",
      icon: "/layout-dashboard.svg",
    },
    {
      id: "classes",
      label: "Lớp học",
      route: "/student/classes*",
      icon: "/book-open-text.svg",
    },
    {
      id: "my-projects",
      label: "Đề tài của tôi",
      route: "/my-projects*",
      icon: "/my_project.svg",
    },
    {
      id: "notifications",
      label: "Thông báo",
      route: "/notifications",
      icon: "/bell.svg",
    },
    // ❌ HIDDEN FOR TEST - TẠM ẨN ĐỂ KIỂM TRA FRONTEND CẬP NHẬT ❌
    // {
    //   id: "chats",
    //   label: "Tin nhắn",
    //   route: "/chats",
    //   icon: "/message-circle.svg",
    // },
  ];

  // Helper check active route
  function isActiveRoute(route: string, currentLocation: string): boolean {
    if (route.endsWith("*")) {
      const baseRoute = route.slice(0, -1);
      return currentLocation.startsWith(baseRoute);
    }
    return currentLocation === route;
  }

  // --- LOGIC PHẢN XẠ (REACTIVE) ---
  const currentUser = $derived($authStore.user);
  const userName = $derived(
    currentUser?.fullname || "Refresh tài khoản tại đây..",
  );
  const userAvatar = $derived(currentUser?.avatar);

  // Lấy số lượng từ Store
  const unreadCount = $derived($unreadNotificationCount);

  // Map lại menuItems để gắn số badge động
  const menuItems = $derived(
    (role === "LECTURER" ? teacherMenuItems : studentMenuItems).map((item) => {
      if (item.id === "notifications") {
        return { ...item, badge: unreadCount };
      }
      return item;
    }),
  );
  onMount(() => {
    notificationStore.fetchNotifications();

    // Nếu muốn tự động cập nhật mỗi 1 phút để hiện badge mới:
    const interval = setInterval(() => {
      notificationStore.fetchNotifications();
    }, 60000); // 60.000ms = 1 phút

    return () => clearInterval(interval); // Xóa bộ đợi khi logout/hủy component
  });
  async function handleLogout() {
    const confirmed = window.confirm("Bạn có chắc chắn muốn đăng xuất?");
    if (!confirmed) return;

    // 1. Xóa dữ liệu (Store + LocalStorage)
    clearAuth();

    // 2. Chuyển về trang login
    push("/login");
  }
</script>

<aside class="sidebar" class:collapsed>
  <div class="sidebar-header">
    <div class="logo-wrapper" class:centered={collapsed} onclick={handleLogout}>
      <div class="logo-container">
        <div class="brand-icon">
          <GraduationCap size={28} strokeWidth={2.5} />
        </div>
        {#if !collapsed}
          <span class="logo">DoAnHUB</span>
        {/if}
      </div>
    </div>
    <button onclick={onToggle} class="toggle-btn" aria-label="Toggle Sidebar">
      <span class="icon">
        <img src="/arrow-left-right.svg" alt="resize" width="22" height="22" />
      </span>
    </button>
  </div>

  <nav class="sidebar-nav">
    <ul>
      {#each menuItems as item}
        <li>
          <a
            href={item.route.replace("*", "")}
            use:link
            class:active-link={isActiveRoute(item.route, $location)}
            title={collapsed ? item.label : ""}
          >
            <span class="icon">
              <img src={item.icon} alt={item.label} width="22" height="22" />

              {#if collapsed && item.badge && item.badge > 0}
                <div class="mini-badge-dot"></div>
              {/if}
            </span>

            <span class="label" class:hidden={collapsed}>{item.label}</span>

            {#if item.badge && item.badge > 0 && !collapsed}
              <span class="badge">
                {item.badge > 99 ? "99+" : item.badge}
              </span>
            {/if}
          </a>
        </li>
      {/each}
    </ul>
  </nav>

  <div class="logout-wrapper">
    <button
      onclick={onLogout}
      class="logout-menu-item"
      title={collapsed ? "Đăng xuất" : ""}
    >
      <span class="icon">
        <img src="/logout_icon.svg" alt="Logout" width="22" height="22" />
      </span>
      <span class="label" class:hidden={collapsed}>Đăng xuất</span>
    </button>
  </div>

  <div class="sidebar-footer">
    <button class="user-profile" onclick={onProfileClick}>
      <div class="avatar">
        {#if userAvatar}
          <img src={userAvatar} alt="User" />
        {:else}
          <div class="avatar-placeholder">
            <img src="/user.svg" />
          </div>
        {/if}
      </div>

      {#if !collapsed}
        <div class="user-info">
          <span class="username">{userName}</span>
          <span class="role">
            {role === "LECTURER" ? "Giáo viên" : "Sinh viên"}
          </span>
        </div>
      {/if}
    </button>
  </div>
</aside>

<style>
  /* --- GIỮ NGUYÊN CSS CŨ CỦA BẠN --- */
  .sidebar {
    height: 100vh;
    width: 256px;
    background: #ffffff;
    border-right: 1px solid #e2e8f0;
    display: flex;
    flex-direction: column;
    transition: width 0.3s ease;
    flex-shrink: 0;
  }
  .sidebar.collapsed {
    width: 80px;
  }
  .sidebar-header {
    height: 70px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    border-bottom: 1px solid #e2e8f0;
  }
  .collapsed .sidebar-header {
    justify-content: center;
    padding: 0;
  }
  .toggle-btn {
    background: rgba(0, 0, 0, 0.05);
    border: none;
    cursor: pointer;
    width: 30px;
    height: 30px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .sidebar-nav {
    flex: 1;
    overflow-y: auto;
    padding: 16px 12px;
  }
  .sidebar-nav ul {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .sidebar-nav li {
    margin-bottom: 4px;
  }
  .sidebar-nav a {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border-radius: 8px;
    color: #1e293b;
    text-decoration: none;
    transition: all 0.2s;
    position: relative; /* Quan trọng để căn badge */
  }
  .sidebar-nav a:hover {
    background: #f1f5f9;
  }
  .active-link {
    background: #dbeafe !important;
    color: #2563eb !important;
  }
  .label {
    font-size: 14px;
    font-weight: 500;
    flex: 1;
  }
  .label.hidden {
    display: none;
  }

  /* --- CSS CHO BADGE SỐ (GIỮ CŨ & CẢI TIẾN) --- */
  .badge {
    background: #ef4444;
    color: white;
    font-size: 11px;
    padding: 2px 6px;
    border-radius: 10px;
    min-width: 20px;
    text-align: center;
    font-weight: 600;
  }

  /* --- THÊM MỚI: CSS CHO CHẤM ĐỎ POP-UP --- */
  .icon {
    position: relative;
    display: flex;
    align-items: center;
  }

  .mini-badge-dot {
    position: absolute;
    top: -2px;
    right: -2px;
    width: 10px;
    height: 10px;
    background-color: #ef4444;
    border: 2px solid white;
    border-radius: 50%;
    z-index: 10;
  }

  /* --- PHẦN CÒN LẠI CỦA SIDEBAR FOOTER (GIỮ CŨ) --- */
  .logout-wrapper {
    padding: 0 12px 12px 12px;
    border-top: 1px solid #e2e8f0;
    margin-top: auto;
  }
  .logout-menu-item {
    width: 100%;
    display: flex;
    align-items: center;
    padding: 12px;
    margin-top: 8px;
    border-radius: 10px;
    border: none;
    background: none;
    cursor: pointer;
    gap: 8px;
  }
  .sidebar-footer {
    padding: 10px;
    display: flex;
    align-items: center;
    border-top: 1px solid #e2e8f0;
  }
  .user-profile {
    display: flex;
    align-items: center;
    flex: 1;
    padding: 8px;
    border-radius: 10px;
    border: 0;
    background: none;
    cursor: pointer;
    gap: 12px;
  }
  .avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
  }
  .avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background-color: #2563eb;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .user-info {
    display: flex;
    flex-direction: column;
    text-align: left;
    min-width: 0;
  }
  .username {
    font-weight: 600;
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .role {
    font-size: 0.75rem;
    color: #64748b;
  }
  .logo-wrapper {
    display: flex;
    align-items: center;
    gap: 10px;
    transition: all 0.3s;
  }

  .logo-wrapper.centered {
    justify-content: center;
    width: 100%;
  }

  .logo-container {
    display: flex;
    align-items: center;
    gap: 8px; /* Khoảng cách giữa icon và chữ */
  }
  .brand-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #2563eb; /* Màu xanh chủ đạo */
    background: #eff6ff; /* Nền xanh nhạt tạo khối */
    padding: 6px;
    border-radius: 10px;
    box-shadow:
      0 4px 6px -1px rgba(37, 99, 235, 0.1),
      0 2px 4px -1px rgba(37, 99, 235, 0.06);
  }

  .logo {
    font-size: 20px;
    font-weight: 800;
    letter-spacing: -0.5px;
    /* Hiệu ứng Gradient cho chữ */
    background: linear-gradient(135deg, #1e293b 30%, #2563eb 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    white-space: nowrap;
  }

  /* Hiệu ứng khi Hover vào Logo */
  .logo-container:hover .brand-icon {
    transform: rotate(-10deg) scale(1.1);
    transition: transform 0.2s ease;
  }
</style>
