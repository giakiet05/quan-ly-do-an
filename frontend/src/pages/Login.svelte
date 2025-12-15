<script lang="ts">
  import AuthLayout from '../layouts/AuthLayout.svelte'; 
  import { link, push } from "svelte-spa-router";
  import Button from "../components/Button.svelte";
  
  // Import service object
  import { authService } from "../services/auth-service";
  
  // Chỉ cần import setAuth để cập nhật UI state (Store), không cần import storage-service nữa
  import { setAuth } from "../stores/auth-store";

  let Email = "";
  let password = "";
  let showPassword = false;
  let loading = false;
  let error = "";

  const togglePassword = () => {
    showPassword = !showPassword;
  };

  const clearError = () => {
    if(error) error = "";
  }

  const handleGoogleLogin = () => {
    authService.loginWithGoogle();
  };

  const handleSubmit = async () => {
    error = "";
    if (!Email || !password) {
      error = "Vui lòng điền đầy đủ thông tin.";
      return;
    }
    
    loading = true;
    
    try {
      // 1. Gọi API Login
      const res = await authService.login({
        identifier: Email,
        password: password,
      });

      console.log("Login thành công:", res);

      // 2. Cập nhật Store (để Header/Sidebar cập nhật giao diện ngay lập tức)
      // res.user và res.accessToken lấy từ kết quả trả về của API
      setAuth(res.user, res.accessToken);

      // 3. Chuyển hướng
      push("/"); 

    } catch (err: any) {
      console.error("Lỗi đăng nhập:", err);
      
      // Xử lý hiển thị lỗi linh hoạt tùy theo backend trả về format gì
      if (err?.message) {
        error = err.message;
      } else if (err?.error?.message) {
        error = err.error.message; // Strapi format thường gặp
      } else if (typeof err === "string") {
        error = err;
      } else {
        error = "Đăng nhập thất bại. Vui lòng kiểm tra lại thông tin.";
      }
    } finally {
      loading = false;
    }
  };
</script>

<AuthLayout>
  <div class="auth-header">
    <div class="icon-wrapper">
       <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="hat-icon">
          <path d="M11.7 2.805a.75.75 0 0 1 .6 0l9.75 4.2a.75.75 0 0 1 0 1.39l-9.75 4.2a.75.75 0 0 1-.6 0l-9.75-4.2a.75.75 0 0 1 0-1.39l9.75-4.2ZM4.32 7.5 12 10.8l7.68-3.3-7.68-3.3-7.68 3.3ZM12 13.2l5.556-2.393a3.75 3.75 0 0 1 1.694 3.332v3.181a.75.75 0 0 1-1.5 0v-3.181c0-.725-.43-1.382-1.096-1.668L12 14.587l-4.654-2.017a1.875 1.875 0 0 0-2.192.334L3.75 15.086v.756a3.75 3.75 0 1 0 7.5 0v-.756l-1.404-2.106c-.366-.548-1.058-.78-1.654-.522L12 13.2Z" />
       </svg>            
    </div>
    <h2 class="title">Chào mừng trở lại</h2>
    <p class="subtitle">Đăng nhập hệ thống quản lý đồ án</p>
  </div>

  <form on:submit|preventDefault={handleSubmit} class="auth-form">
    
    {#if error}
      <div class="error-alert">{error}</div>
    {/if}

    <div class="form-group">
      <label for="identity">Email <span class="required">*</span></label>
      <input
        id="identity"
        type="email"
        bind:value={Email}
        placeholder="Nhập email của bạn"
        disabled={loading}
        on:input={clearError} 
      />
      </div>

    <div class="form-group">
      <label for="password">Mật khẩu <span class="required">*</span></label>
      <div class="input-wrapper-with-icon">
        <input
          id="password"
          type="text" 
          hidden={!showPassword}
          value={password}
          disabled={loading}
          class="has-icon-right"
          on:input={(e) => { password = e.currentTarget.value; clearError(); }}
        />
        <input
          id="password-field"
          type={showPassword ? "text" : "password"}
          bind:value={password}
          placeholder="Nhập mật khẩu"
          disabled={loading}
          class="has-icon-right"
          on:input={clearError}
        />
        
        <div class="toggle-password-wrapper">
            <Button 
                variant="icon" 
                type="button" 
                onclick={togglePassword} 
                title={showPassword ? "Ẩn mật khẩu" : "Hiện mật khẩu"}
                tabindex="-1"
            >
                {#if showPassword}
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="eye-icon"><path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" /></svg>
                {:else}
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="eye-icon"><path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88" /></svg>
                {/if}
            </Button>
        </div>
      </div>
    </div>

    <div class="forgot-password-link">
        <a href="/forgot-password" use:link>Quên mật khẩu?</a>
    </div>

    <Button 
        variant="primary" 
        type="submit" 
        disabled={loading} 
        class="w-full"
    >
      {#if loading}
        <span class="loader"></span> Đang xử lý...
      {:else}
        Đăng nhập
      {/if}
    </Button>

  </form>
  
  <div class="divider">
    <span>hoặc</span>
  </div>

  <Button 
      variant="google" 
      type="button" 
      onclick={handleGoogleLogin} 
      disabled={loading}
      class="w-full"
  >
    Đăng nhập bằng Google
  </Button>

  <div class="auth-footer">
      Chưa có tài khoản? 
      <a href="/register" class="link-text" use:link>Đăng ký ngay</a>
  </div>
</AuthLayout>

<style>
  /* --- Variables & Base Styles --- */
  :root {
    --primary-blue: #1a56db;
    --text-dark: #111827;
    --text-gray: #6b7280;
    --input-bg: #f3f5f7;
  }

  /* --- MỚI: Style cho dấu sao bắt buộc --- */
  .required {
    color: #ef4444; /* Màu đỏ */
    margin-left: 4px;
    font-weight: bold;
  }

  .auth-header { margin-bottom: 32px; }
  .icon-wrapper { display: inline-flex; justify-content: center; align-items: center; margin-bottom: 16px; }
  .hat-icon { width: 48px; height: 48px; color: var(--primary-blue); }
  .title { font-size: 28px; font-weight: 700; color: var(--text-dark); margin: 0 0 8px 0; }
  .subtitle { font-size: 16px; color: var(--text-gray); margin: 0; }

  .auth-form { text-align: left; }
  .form-group { margin-bottom: 24px; }
  label { display: block; font-size: 16px; font-weight: 600; color: var(--text-dark); margin-bottom: 8px; }

  input {
    width: 100%; padding: 14px 16px; background-color: var(--input-bg);
    border: 2px solid transparent; border-radius: 12px; font-size: 16px;
    color: var(--text-dark); transition: all 0.2s; box-sizing: border-box;
  }
  input::placeholder { color: #9ca3af; }
  input:focus { outline: none; background-color: #fff; border-color: var(--primary-blue); }

  .input-wrapper-with-icon { position: relative; }
  input.has-icon-right { padding-right: 48px; }

  .toggle-password-wrapper {
    position: absolute; 
    right: 8px; 
    top: 50%; 
    transform: translateY(-50%);
    z-index: 10;
  }
  .eye-icon { width: 20px; height: 20px; }

  .forgot-password-link { text-align: right; margin-top: -16px; margin-bottom: 24px; }
  .forgot-password-link a { color: var(--primary-blue); font-weight: 600; text-decoration: none; font-size: 14px; }

  /* Utility class cho Button full width */
  :global(.w-full) {
    width: 100% !important;
    display: flex !important; /* Đảm bảo flex để căn giữa nội dung */
  }

  .auth-footer { margin-top: 32px; font-size: 16px; color: var(--text-dark); }
  .link-text { color: var(--primary-blue); font-weight: 700; margin-left: 4px; text-decoration: none; }
  .link-text:hover { text-decoration: underline; }

  .error-alert { background-color: #fee2e2; color: #991b1b; padding: 12px; border-radius: 8px; margin-bottom: 20px; font-size: 14px; text-align: left; }

  .loader {
    width: 20px; height: 20px; border: 3px solid #fff; border-bottom-color: transparent;
    border-radius: 50%; display: inline-block; animation: rotation 1s linear infinite;
  }
  @keyframes rotation { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }

  .divider {
    display: flex; align-items: center; margin: 24px 0; color: #6b7280; font-size: 14px;
  }
  
  .divider::before, .divider::after {
    content: ''; flex: 1; height: 1px; background-color: #e5e7eb;
  }

  .divider span {
    padding: 0 12px; font-weight: 500; color: #9ca3af; text-transform: lowercase;
  }
</style>
