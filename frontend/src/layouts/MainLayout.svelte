<script>
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";

    let { children } = $props();

    // Svelte 5 state để quản lý thu gọn sidebar
    let isCollapsed = $state(false);

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }

    const menuItems = [
        { name: "Home", path: "/", icon: "🏠" },
        { name: "Dashboard", path: "/dashboard", icon: "📊" },
        { name: "Settings", path: "/settings", icon: "⚙️" },
    ];
</script>

<div class="main-layout">
    <aside class="sidebar" class:collapsed={isCollapsed}>
        <div class="sidebar-header">
            {#if !isCollapsed}
                <span class="logo">MY APP</span>
            {/if}
            <button onclick={toggleSidebar} class="toggle-btn">
                {isCollapsed ? "→" : "←"}
            </button>
        </div>

        <nav class="sidebar-nav">
            <ul>
                {#each menuItems as item}
                    <li>
                        <a
                            href={item.path}
                            use:link
                            use:active={{
                                path: item.path,
                                className: "active-link",
                            }}
                        >
                            <span class="icon">{item.icon}</span>
                            {#if !isCollapsed}
                                <span class="label">{item.name}</span>
                            {/if}
                        </a>
                    </li>
                {/each}
            </ul>
        </nav>

        <div class="sidebar-footer">
            <a href="/auth/login" use:link>
                <span class="icon">🚪</span>
                {#if !isCollapsed}
                    <span class="label">Logout</span>
                {/if}
            </a>
        </div>
    </aside>

    <main class="content">
        <header class="top-bar">
            <h1>Dashboard Section</h1>
        </header>
        <div class="page-body">
            {@render children()}
        </div>
    </main>
</div>

<style>
    :root {
        --sidebar-width: 260px;
        --sidebar-collapsed-width: 70px;
        --primary-color: #34495e;
        --active-color: #3498db;
        --transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .main-layout {
        display: flex;
        height: 100vh;
        overflow: hidden;
    }

    /* Sidebar Base */
    .sidebar {
        width: var(--sidebar-width);
        background: #2c3e50;
        color: white;
        display: flex;
        flex-direction: column;
        transition: var(--transition);
        position: relative;
    }

    .sidebar.collapsed {
        width: var(--sidebar-collapsed-width);
    }

    /* Header & Toggle Button */
    .sidebar-header {
        height: 60px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 20px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .collapsed .sidebar-header {
        justify-content: center;
        padding: 0;
    }

    .toggle-btn {
        background: #34495e;
        border: none;
        color: white;
        cursor: pointer;
        padding: 5px 10px;
        border-radius: 4px;
    }

    /* Navigation */
    .sidebar-nav {
        flex: 1;
        padding: 10px 0;
    }

    .sidebar-nav ul {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .sidebar-nav a {
        display: flex;
        align-items: center;
        padding: 12px 20px;
        color: #bdc3c7;
        text-decoration: none;
        transition: background 0.2s;
        white-space: nowrap; /* Tránh text bị nhảy dòng khi thu hẹp */
    }

    .collapsed .sidebar-nav a {
        justify-content: center;
        padding: 12px 0;
    }

    .sidebar-nav a:hover {
        background: rgba(255, 255, 255, 0.05);
        color: white;
    }

    .icon {
        font-size: 1.2rem;
        min-width: 30px;
        display: flex;
        justify-content: center;
    }

    /* Active Link Style */
    :global(.active-link) {
        background: var(--active-color) !important;
        color: white !important;
        border-left: 4px solid #fff;
    }

    /* Content Area */
    .content {
        flex: 1;
        display: flex;
        flex-direction: column;
        background: #f8f9fa;
    }

    .top-bar {
        background: white;
        padding: 0 2rem;
        height: 60px;
        display: flex;
        align-items: center;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }

    .page-body {
        flex: 1;
        padding: 2rem;
        overflow-y: auto;
    }
</style>
