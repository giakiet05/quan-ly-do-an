<!-- <script lang="ts">
  import { onMount } from "svelte";
  import Router, { push, replace } from "svelte-spa-router";
  import { authService } from "./services/auth-service";
  import routes from "./routes";
  import Topbar from "./components/Topbar.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import Footer from "./components/Footer.svelte"; 
  import { authStore, setAuth } from "./stores/auth-store";

  let topbarUser: any = undefined;
  let isSidebarCompact = false;

  $: topbarUser = $authStore.user;

  authStore.subscribe((state) => {
    if (state.user) {
      topbarUser = {
        name: state.user.username || state.user.email || "User",
        avatar: state.user.avatar || undefined
      };
    } else {
      topbarUser = undefined;
    }
  });

  const sidebarItems = [
    { id: "home", label: "Home", to: "/" },
    { id: "popular", label: "Popular", to: "/popular" },
    { id: "explore", label: "Explore", to: "/explore" },
    { id: "all", label: "All", to: "/all" }
  ];

  onMount(() => {
    const result = authService.handleLoginCallback();
    if (result.success) {
      setAuth(result.user, result.accessToken);
      replace("/");
      console.log("✅ Đã login Google và cập nhật Auth Store");
    }
  });

  function handleLogout() {
    console.log("Logout triggered");
  }

  function handleNavigate(item: any) {
    push(item.to);
  }
</script>

<div class="app-layout">
  <div class="layout-topbar" data-compact={isSidebarCompact}>
    <Topbar user={topbarUser} onLogout={handleLogout} />
  </div>

  <Sidebar
    items={sidebarItems}
    onNavigate={handleNavigate}
    bind:compact={isSidebarCompact}
  />

  <main class="main-content" data-compact={isSidebarCompact}>
    <div class="page-content">
      <Router {routes} />
    </div>
    <Footer />
  </main>
</div>

<style>
  :root {
    --sidebar-width: 256px;
    --sidebar-compact-width: 64px;
    --topbar-height: 56px;
    --transition-speed: 0.2s;
  }

  .app-layout {
    position: relative;
    min-height: 100vh;
    background-color: white;
  }

  .layout-topbar {
    position: fixed;
    top: 0;
    right: 0;
    z-index: 50;
    left: var(--sidebar-width);
    height: var(--topbar-height);
    transition: left var(--transition-speed) ease;
  }

  .layout-topbar[data-compact="true"] {
    left: var(--sidebar-compact-width);
  }

  .main-content {
    margin-left: var(--sidebar-width);
    padding-top: var(--topbar-height);
    transition: margin-left var(--transition-speed) ease;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    box-sizing: border-box;
  }

  .page-content {
    flex: 1;
    width: 100%;
  }

  .main-content[data-compact="true"] {
    margin-left: var(--sidebar-compact-width);
  }

  @media (max-width: 768px) {
    .layout-topbar {
      left: 0;
    }
    .main-content {
      margin-left: 0;
    }
  }
</style> -->

<script>
  import Router from "svelte-spa-router";
  import { location } from "svelte-spa-router";

  import MainLayout from "./layouts/MainLayout.svelte";
  import AuthLayout from "./layouts/AuthLayout.svelte";
  import Login from "./pages/Login.svelte";
  import routes from "./routes";
</script>

{#if $location.startsWith("/auth")}
  <AuthLayout>
    <Router {routes} />
  </AuthLayout>
{:else}
  <MainLayout>
    <Router {routes} />
  </MainLayout>
{/if}
