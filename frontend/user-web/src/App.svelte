<script lang="ts">
  import { onMount } from "svelte";
  import Router, { push, replace } from "svelte-spa-router";
  import { authService } from "./services/auth-service";
  import routes from "./routes";
  import Topbar from "./components/Topbar.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import Footer from "./components/Footer.svelte"; 
  import { authStore, setAuth } from "./stores/auth-store";

  const sidebarItems = [ /* ... */ ];
  let isSidebarCompact = false; // Biến này đang điều khiển độ rộng Sidebar
  let topbarUser: any = undefined;
  $: topbarUser = $authStore.user;
  // ... (Logic giữ nguyên)
  onMount(() => {
    // Gọi service kiểm tra URL xem có token do Google trả về không
    const result = authService.handleLoginCallback();

    if (result.success) {
      // 1. Cập nhật Store (lúc này biến $authStore thay đổi -> topbarUser tự cập nhật)
      setAuth(result.user, result.accessToken);

      // 2. Xóa token trên URL cho sạch sẽ (quan trọng)
      replace('/'); 
      
      console.log("✅ Đã login Google và cập nhật Auth Store");
    }
  });
  function handleLogout() { /* ... */ }
  function handleNavigate(item: any) { push(item.to); }
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
    --transition-speed: 0.2s; /* Thêm biến này cho đồng bộ chuyển động */
  }

  .app-layout {
    position: relative;
    min-height: 100vh;
    background-color: white;
  }

  /* --- 2. CSS CHO TOPBAR WRAPPER (PHẦN QUAN TRỌNG NHẤT) --- */
  .layout-topbar {
    position: fixed;
    top: 0;
    right: 0;
    z-index: 50; /* Thấp hơn Sidebar (thường là 100) nhưng cao hơn nội dung */
    
    /* Mấu chốt: Bắt đầu từ vị trí kết thúc của Sidebar */
    left: var(--sidebar-width); 
    
    /* Chiều cao cố định */
    height: var(--topbar-height);
    
    /* Hiệu ứng trượt mượt mà khi Sidebar co giãn */
    transition: left var(--transition-speed) ease;
  }

  /* Khi Sidebar thu nhỏ -> Topbar giãn ra sang trái */
  .layout-topbar[data-compact="true"] {
    left: var(--sidebar-compact-width);
  }

  /* --- 3. CSS CHO MAIN CONTENT (Giữ nguyên logic cũ của bạn, chỉnh lại chút) --- */
  .main-content {
    /* Đẩy sang phải né Sidebar */
    margin-left: var(--sidebar-width);
    
    /* Đẩy xuống dưới né Topbar */
    padding-top: var(--topbar-height);
    
    transition: margin-left var(--transition-speed) ease;
    
    /* Flexbox cho Footer */
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

  /* --- 4. RESPONSIVE (MOBILE) --- */
  @media (max-width: 768px) {
    /* Trên mobile, Sidebar thường ẩn đi (left: -100%) */
    
    .layout-topbar {
      left: 0; /* Topbar về full màn hình */
    }
    
    .main-content {
      margin-left: 0; /* Nội dung về full màn hình */
    }
  }
</style>