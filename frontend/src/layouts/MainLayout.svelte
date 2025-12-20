<script lang="ts">
    import { link, location, push } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
    import type { SidebarItem } from "../types/sidebar";
    import { Roles } from "../types/role";
    import { authStore } from "../stores/auth-store";

    let { children } = $props();
    let isCollapsed = $state(false);

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }

    function goToMyProfile() {
        push("/my-profile");
    }

    const menuItems: SidebarItem[] = [
        {
            id: "dashboard",
            label: "Tổng quan",
            route: "/dashboard",
            icon: "/layout-dashboard.svg",
            roles: [Roles.LECTURER, Roles.STUDENT],
        },
        {
            id: "classes",
            label: "Lớp học",
            route: "/classes",
            icon: "/book-open-text.svg",
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
        {
            id: "students",
            label: "Sinh viên",
            route: "/students",
            icon: "/users.svg",
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
        {
            id: "notifications",
            label: "Thông báo",
            route: "/notifications",
            icon: "/bell.svg",
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
        {
            id: "messages",
            label: "Tin nhắn",
            route: "/messages",
            icon: "/message-circle.svg",
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
    ];
    let currentRole = $state<string | null>("LECTURER");
    const filteredMenuItems = $derived(
        menuItems.filter((item) =>
            item.roles.includes(currentRole as keyof typeof Roles),
        ),
    );

    const currentBreadcrumb = $derived(
        menuItems.find((item) => item.route === $location)?.label ||
            "Trang chủ",
    );
</script>

<div class="main-layout">
    <aside class="sidebar" class:collapsed={isCollapsed}>
        <div class="sidebar-header">
            <div class="logo-wrapper">
                {#if !isCollapsed}
                    <span class="logo">MY APP</span>
                {/if}
            </div>
            <button
                onclick={toggleSidebar}
                class="toggle-btn"
                aria-label="Toggle Sidebar"
            >
                <span class="icon">
                    <img
                        src={"/arrow-left-right.svg"}
                        alt={"resize"}
                        width="22"
                        height="22"
                    />
                </span>
            </button>
        </div>

        <nav class="sidebar-nav">
            <ul>
                {#each filteredMenuItems as item}
                    <li>
                        <a
                            href={item.route}
                            use:link
                            use:active={{
                                path: item.route,
                                className: "active-link",
                            }}
                            title={isCollapsed ? item.label : ""}
                        >
                            <span class="icon">
                                <img
                                    src={item.icon}
                                    alt={item.label}
                                    width="22"
                                    height="22"
                                />
                            </span>
                            <span class="label" class:hidden={isCollapsed}
                                >{item.label}</span
                            >
                        </a>
                    </li>
                {/each}
            </ul>
        </nav>
        <div class="sidebar-footer">
            <button class="user-profile" onclick={goToMyProfile}>
                <div class="avatar">
                    {#if $authStore.user?.avatar}
                        <img src={$authStore.user.avatar} alt="User" />
                    {:else}
                        <div class="avatar-placeholder">
                            <span class="icon">
                                <img
                                    src={"/user.svg"}
                                    alt={"user"}
                                    width="22"
                                    height="22"
                                />
                            </span>
                        </div>
                    {/if}
                </div>
                {#if !isCollapsed}
                    <div class="user-info">
                        <span class="username"
                            >{$authStore.user?.fullname || "User"}</span
                        >
                        <span class="role"
                            >{$authStore.user?.role || "Member"}</span
                        >
                    </div>
                {/if}
            </button>
        </div>
    </aside>

    <main class="content">
        <header class="top-bar">
            <div class="header-left">
                <div class="breadcrumb">
                    <span class="breadcrumb-item">Trang chủ</span>
                    <span class="separator">/</span>
                    <span class="breadcrumb-item active"
                        >{currentBreadcrumb}</span
                    >
                </div>
            </div>
        </header>
        <div class="page-body">
            {@render children()}
        </div>
    </main>
</div>

<style>
    .main-layout {
        --sidebar-width: 260px;
        --sidebar-collapsed-width: 80px;
        --active-text: #3b82f6;
        --active-item: #b5d3fa;
        --hover-item: #cccccc;
        --transition-speed: 0.3s;
    }
    .sidebar-footer {
        padding: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-top: 1px solid #e2e8f0;
    }
    .sidebar:not(.collapsed) .sidebar-footer {
        justify-content: flex-start;
    }
    .user-profile {
        display: flex;
        align-items: center;
        width: 100%;
        padding: 8px;
        border-radius: 10px;
        border: 0;
        background: none;
        cursor: pointer;
        transition: background 0.2s;
    }
    .user-profile:hover {
        background: var(--hover-item);
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
        background-color: var(--active-text);
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
        margin-left: 12px;
        display: flex;
        flex-direction: column;
    }
    .username {
        font-weight: 600;
        font-size: 0.9rem;
        color: #1e293b;
    }
    .role {
        font-size: 0.75rem;
        color: #64748b;
    }
    .main-layout {
        display: flex;
        height: 100vh;
        background-color: #f1f5f9;
    }

    .sidebar {
        width: var(--sidebar-width);
        display: flex;
        flex-direction: column;
        transition: width var(--transition-speed) cubic-bezier(0.4, 0, 0.2, 1);
        box-shadow: 4px 0 10px rgba(0, 0, 0, 0.1);
        z-index: 100;
    }

    .sidebar.collapsed {
        width: var(--sidebar-collapsed-width);
    }

    .sidebar-header {
        height: 70px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 20px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }

    .collapsed .sidebar-header {
        justify-content: center;
        padding: 0;
    }

    .logo {
        font-weight: 700;
        letter-spacing: 1px;
    }

    .toggle-btn {
        background: rgba(255, 255, 255, 0.1);
        border: none;
        cursor: pointer;
        width: 30px;
        height: 30px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: transform 0.3s;
    }

    .toggle-btn:hover {
        background: var(--hover-item);
    }

    .sidebar-nav {
        flex: 1;
        padding: 15px 12px;
    }

    .sidebar-nav ul {
        list-style: none;
        padding: 0;
    }

    .sidebar-nav a {
        display: flex;
        align-items: center;
        padding: 12px;
        margin-bottom: 8px;
        color: rgb(0, 0, 0);
        text-decoration: none;
        border-radius: 10px;
        transition: all 0.2s;
        position: relative;
    }

    .sidebar-nav a:hover {
        background: var(--hover-item);
    }

    .icon {
        min-width: 44px; /* Cố định độ rộng icon để không bị lệch khi thu nhỏ */
        display: flex;
        justify-content: center;
        transition: margin 0.3s;
    }

    .label {
        font-size: 0.95rem;
        font-weight: 500;
        white-space: nowrap;
        opacity: 1;
        transition:
            opacity 0.2s,
            visibility 0.2s;
    }

    .label.hidden {
        opacity: 0;
        visibility: hidden;
        width: 0;
    }

    /* Active Link Style */
    :global(.active-link) {
        background: var(--active-item) !important;
        color: var(--active-text) !important;
        box-shadow: 0 4px 12px var(--active-item);
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
