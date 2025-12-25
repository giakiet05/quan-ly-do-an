<script>
  // AuthLayout.svelte
</script>

<div class="auth-layout">
  <div class="auth-card">
    <slot />
  </div>
</div>

<style>
  .auth-layout {
    /* --- 1. XÓA BỎ CÁC DÒNG GÂY LỖI --- */
    /* margin-top: ... -> XÓA (Thủ phạm gây dính lên trên) */
    /* padding-top: ... -> XÓA (Không cần bù trừ nữa) */
    
    /* --- 2. CĂN GIỮA CHUẨN --- */
    /* Chiều cao tối thiểu = Full màn hình trừ đi Topbar */
    /* Giúp Card luôn nằm giữa phần không gian còn lại */
    min-height: calc(100vh - var(--topbar-height));
    
    width: 100%;
    
    display: flex;
    flex-direction: column;
    justify-content: center; /* Căn giữa dọc */
    align-items: center;     /* Căn giữa ngang */
    
    background-color: var(--background);
    
    /* Thêm padding để khi màn hình nhỏ, Card không dính sát lề */
    padding: 24px;
    box-sizing: border-box;
  }

  .auth-card {
    background: var(--surface-color);
    width: 100%;
    max-width: 440px;
    padding: 40px;
    border-radius: 16px;
    border: 1px solid var(--border-color);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    text-align: center;
  }

  @media (max-width: 480px) {
    .auth-layout {
        padding: 16px; /* Giảm padding trên mobile */
    }
    .auth-card {
      box-shadow: none;
      border: none;
      background: transparent;
      padding: 0;
    }
  }
</style>