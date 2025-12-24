<script lang="ts">
  // Không cần props title/subtitle nữa
</script>

<div class="auth-wrapper">
  <div class="auth-container">
    <div class="left-panel">
      <div class="brand-logo">
        <svg width="40" height="40" viewBox="0 0 32 32" fill="none">
          <circle cx="16" cy="16" r="16" fill="#1a56db" />
          <path
            d="M16 16L8 8C8 8 8 24 16 24C24 24 24 8 24 8L16 16Z"
            fill="white"
          />
        </svg>
        <span class="brand-name">DoAnHub</span>
      </div>

      <div class="branding-content">
        <div class="cyan-shape"></div>
        <h3>Quản lý Đồ án hiệu quả</h3>
        <p>
          Nền tảng hỗ trợ sinh viên và giảng viên tối ưu hóa quy trình thực hiện
          đồ án.
        </p>
      </div>

      <div class="footer-note">© 2024 DoAnHub Team</div>
    </div>

    <div class="right-panel">
      <div class="form-content">
        <!-- Slot nhận toàn bộ nội dung từ page con -->
        <slot></slot>
      </div>
    </div>
  </div>
</div>

<style>
  .auth-wrapper {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f3f4f6;
    padding: 20px;
  }

  .auth-container {
    display: flex;
    width: 100%;
    max-width: 1000px;
    min-height: 600px;
    background: white;
    border-radius: 24px;
    overflow: hidden;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  }

  /* Panel xanh ngọc bên trái */
  .left-panel {
    flex: 0 0 45%;
    background: linear-gradient(135deg, #b2dfdb 0%, #009688 100%);
    padding: 40px 30px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    color: #004d40;
  }

  .brand-logo {
    display: flex;
    align-items: center;
    gap: 12px;
    font-weight: 800;
    font-size: 24px;
  }

  .brand-name {
    color: #004d40;
  }

  .branding-content {
    text-align: center;
  }

  .cyan-shape {
    width: 140px;
    height: 140px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    margin: 0 auto 24px;
    backdrop-filter: blur(5px);
  }

  .branding-content h3 {
    font-size: 22px;
    margin-bottom: 12px;
    font-weight: 700;
    color: #004d40;
  }

  .branding-content p {
    font-size: 14px;
    opacity: 0.9;
    line-height: 1.5;
    color: #004d40;
  }

  .footer-note {
    font-size: 12px;
    opacity: 0.7;
    text-align: center;
  }

  /* Panel nội dung bên phải */
  .right-panel {
    flex: 1;
    padding: 40px 50px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    overflow-y: auto;
  }

  .form-content {
    width: 100%;
    max-width: 420px;
    margin: 0 auto;
  }

  /* Mobile responsive */
  @media (max-width: 900px) {
    .auth-container {
      flex-direction: column;
      max-width: 500px;
    }

    .left-panel {
      flex: 0 0 auto;
      padding: 30px 20px;
      min-height: 200px;
    }

    .cyan-shape {
      width: 100px;
      height: 100px;
    }

    .branding-content h3 {
      font-size: 20px;
    }

    .branding-content p {
      font-size: 13px;
    }

    .right-panel {
      padding: 40px 24px;
    }
  }

  @media (max-width: 640px) {
    .auth-wrapper {
      padding: 12px;
    }

    .auth-container {
      border-radius: 16px;
    }

    .right-panel {
      padding: 30px 20px;
    }
  }
</style>
