<script lang="ts">
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
    import type { SidebarItem } from "../types/sidebar";
    import { Icons } from "../assets/icons/icons";
    import { Roles } from "../types/role";

    let { children } = $props();
    let isCollapsed = $state(false);

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }

    // menu item có thuoccj tính role để dựa vào ai đang đăng nhập sẽ hiện menu item đó
    const menuItems: SidebarItem[] = [
        /* ... mảng menuItems của bạn ... */
        {
            id: "dashboard",
            label: "Tổng quan",
            route: "/dashboard",
            icon: Icons.dashboard,
            roles: [Roles.LECTURER, Roles.STUDENT],
        },
        {
            id: "classes",
            label: "Lớp học",
            route: "/classes",
            icon: Icons.class,
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
        {
            id: "students",
            label: "Sinh viên",
            route: "/students",
            icon: Icons.students,
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
        {
            id: "notifications",
            label: "Thông báo",
            route: "/notifications",
            icon: Icons.notification,
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
        {
            id: "messages",
            label: "Tin nhắn",
            route: "/messages",
            icon: Icons.message,
            roles: [Roles.STUDENT, Roles.LECTURER],
        },
    ];
    let currentRole = $state<string | null>("LECTURER");
    const filteredMenuItems = $derived(
        menuItems.filter((item) =>
            item.roles.includes(currentRole as keyof typeof Roles),
        ),
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
                <svelte:component this={Icons.resize} size={16} />
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
                                <svelte:component
                                    this={item.icon}
                                    size={22}
                                    strokeWidth={2}
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
            <a href="/auth/login" use:link>
                <span class="icon">
                    <svelte:component
                        this={Icons.logout}
                        size={22}
                        strokeWidth={2}
                    />
                </span>

                {#if !isCollapsed}
                    <span class="label">Logout</span>
                {/if}
            </a>
        </div>
    </aside>

    <main class="content">
        <header class="top-bar">
            <div class="header-left">
                <h2>Dashboard Section</h2>
            </div>
        </header>
        <div class="page-body">
            {@render children()}
        </div>
    </main>
</div>

<style>
    :root {
        --sidebar-width: 260px;
        --sidebar-collapsed-width: 80px;
        --primary-bg: #1e293b;
        --active-bg: #3b82f6;
        --transition-speed: 0.3s;
    }
    .sidebar-footer {
        height: 50px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-top: 1px solid rgba(255, 255, 255, 0.05);
    }
    .main-layout {
        display: flex;
        height: 100vh;
        background-color: #f1f5f9;
    }

    .sidebar {
        width: var(--sidebar-width);
        background: var(--primary-bg);
        color: #cbd5e1;
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
        color: white;
        letter-spacing: 1px;
    }

    .toggle-btn {
        background: rgba(255, 255, 255, 0.1);
        border: none;
        color: white;
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
        background: rgba(255, 255, 255, 0.2);
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
        color: #94a3b8;
        text-decoration: none;
        border-radius: 10px;
        transition: all 0.2s;
        position: relative;
    }

    .sidebar-nav a:hover {
        background: rgba(255, 255, 255, 0.05);
        color: white;
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
        background: var(--active-bg) !important;
        color: white !important;
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
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
</style>
