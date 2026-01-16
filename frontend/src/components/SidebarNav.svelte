<script lang="ts">
  import { link, location, push } from "svelte-spa-router";
  import active from "svelte-spa-router/active";
  import { authStore } from "../stores/auth-store";

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

  const teacherMenuItems: MenuItem[] = [
    {
      id: "dashboard",
      label: "Tổng quan",
      route: "/dashboard",
      icon: "/layout-dashboard.svg",
    },
    {
      id: "classes",
      label: "Lớp học",
      route: "/lecture/my-classes",
      icon: "/book-open-text.svg",
    },
    {
      id: "students",
      label: "Sinh viên",
      route: "/lecture/students-management",
      icon: "/users.svg",
    },
    {
      id: "notifications",
      label: "Thông báo",
      route: "/notifications",
      icon: "/bell.svg",
      badge: 0,
    },
    {
      id: "chats",
      label: "Tin nhắn",
      route: "/chats",
      icon: "/message-circle.svg",
      badge: 0,
    },
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
      route: "/classes*",
      icon: "/book-open-text.svg",
    },
    {
      id: "my-projects",
      label: "Đề tài của tôi",
      route: "/my-projects",
      icon: "/my_project.svg",
    },
    {
      id: "notifications",
      label: "Thông báo",
      route: "/notifications",
      icon: "/bell.svg",
      badge: 0,
    },
    {
      id: "chats",
      label: "Tin nhắn",
      route: "/chats",
      icon: "/message-circle.svg",
      badge: 0,
    },
  ];

  const menuItems = $derived(
    role === "LECTURER" ? teacherMenuItems : studentMenuItems
  );

  // Helper to check if route is active (supports wildcards)
  function isActiveRoute(route: string, currentLocation: string): boolean {
    if (route.endsWith("*")) {
      const baseRoute = route.slice(0, -1); // Remove *
      return currentLocation.startsWith(baseRoute);
    }
    return currentLocation === route;
  }
</script>

<aside class="sidebar" class:collapsed>
  <div class="sidebar-header">
    <div class="logo-wrapper">
      {#if !collapsed}
        <span class="logo">MY APP</span>
      {/if}
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

  <!-- Logout Button -->
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
        {#if $authStore.user?.avatar}
          <img src={$authStore.user.avatar} alt="User" />
        {:else}
          <div class="avatar-placeholder">
            <span class="icon">
              <img src="/user.svg" alt="user" width="22" height="22" />
            </span>
          </div>
        {/if}
      </div>
      {#if !collapsed}
        <div class="user-info">
          <span class="username">{$authStore.user?.fullname || "User"}</span>
          <span class="role"
            >{role === "LECTURER" ? "Giáo viên" : "Sinh viên"}</span
          >
        </div>
      {/if}
    </button>
  </div>
</aside>

<style>
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

  .logo {
    font-weight: 700;
    letter-spacing: 1px;
    color: #1e293b;
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
    transition: background 0.2s;
  }

  .toggle-btn:hover {
    background: rgba(0, 0, 0, 0.1);
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
    position: relative;
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

  .badge {
    background: #ef4444;
    color: white;
    font-size: 11px;
    padding: 2px 6px;
    border-radius: 10px;
    min-width: 20px;
    text-align: center;
  }

  .sidebar-footer {
    padding: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
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
    transition: background 0.2s;
    gap: 12px;
  }

  .user-profile:hover {
    background: #f1f5f9;
  }

  .avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    overflow: hidden;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background-color: #2563eb;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
  }

  .avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
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
    color: #1e293b;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .role {
    font-size: 0.75rem;
    color: #64748b;
  }

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
    color: #1e293b;
    text-decoration: none;
    border: none;
    background: none;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
    gap: 8px;
  }

  .logout-menu-item:hover {
    background: #f1f5f9;
  }

  @media (max-width: 768px) {
    .sidebar {
      width: 100%;
      max-width: 320px;
    }
  }
</style>
