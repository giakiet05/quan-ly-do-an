<script lang="ts">
  import AuthLayout from "../layouts/AuthLayout.svelte";
  import { link, replace } from "svelte-spa-router";
  import OtpVerification from "../components/OtpVerification.svelte";
  import Button from "../components/Button.svelte";
  import {
    authService,
    sendEmailVerification,
    completeRegistration,
  } from "../services/auth-service";

  // --- STATE ---
  let fullname = "";
  let email = "";
  let password = "";
  let confirmPassword = "";
  let step = "register";

  // UI States
  let showPassword = false;
  let showConfirmPassword = false;
  let loading = false;
  let error = "";
  let verificationToken = "";

  const togglePassword = () => (showPassword = !showPassword);
  const toggleConfirmPassword = () =>
    (showConfirmPassword = !showConfirmPassword);

  // THÊM: Hàm xóa lỗi khi người dùng nhập lại (cho mượt)
  const clearError = () => {
    if (error) error = "";
  };

  const handleGoogleRegister = () => {
    authService.loginWithGoogle();
  };

  const handleOtpSuccess = async (
    event: CustomEvent<{ verificationToken: string }>,
  ) => {
    loading = true;
    error = "";

    try {
      verificationToken = event.detail.verificationToken;

      await completeRegistration({
        verificationToken,
        username: fullname,
        password,
      });
      replace("/");
    } catch (e: any) {
      error = e?.message || "Đăng ký thất bại";
      // Có thể không cần reset step về register nếu muốn cho user nhập lại OTP
    } finally {
      loading = false;
    }
  };

  const handleSubmit = async () => {
    error = "";

    if (!fullname || !email || !password || !confirmPassword) {
      error = "Vui lòng điền đầy đủ các trường bắt buộc.";
      return;
    }

    if (password !== confirmPassword) {
      error = "Mật khẩu xác nhận không khớp.";
      return;
    }

    if (password.length < 6) {
      error = "Mật khẩu phải có ít nhất 6 ký tự.";
      return;
    }

    loading = true;

    try {
      await sendEmailVerification({ email });
      step = "verify";
    } catch (e: any) {
      error = e?.message || "Không thể gửi mã xác thực";
    } finally {
      loading = false;
    }
  };
</script>

<AuthLayout>
  {#if step === "register"}
    <div class="auth-header">
      <div class="icon-wrapper">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="hat-icon"
        >
          <path
            d="M5.25 6.375a4.125 4.125 0 1 1 8.25 0 4.125 4.125 0 0 1-8.25 0ZM2.25 19.125a7.125 7.125 0 0 1 14.25 0v.003l-.001.119a.75.75 0 0 1-.363.63 13.067 13.067 0 0 1-6.761 1.873c-2.472 0-4.786-.684-6.76-1.873a.75.75 0 0 1-.364-.63l-.001-.122ZM18.75 7.5a.75.75 0 0 0-1.5 0v2.25H15a.75.75 0 0 0 0 1.5h2.25v2.25a.75.75 0 0 0 1.5 0v-2.25H21a.75.75 0 0 0 0-1.5h-2.25V7.5Z"
          />
        </svg>
      </div>
      <h2 class="title">Tạo tài khoản mới</h2>
      <p class="subtitle">Tham gia hệ thống quản lý đồ án ngay hôm nay</p>
    </div>

    <form on:submit|preventDefault={handleSubmit} class="auth-form">
      {#if error}
        <div class="error-alert">{error}</div>
      {/if}

      <div class="form-group">
        <label for="fullname"
          >Tên hiển thị <span class="required">*</span></label
        >
        <input
          id="fullname"
          type="text"
          bind:value={fullname}
          placeholder="Ví dụ: Nguyễn Văn A"
          disabled={loading}
          on:input={clearError}
        />
      </div>

      <div class="form-group">
        <label for="email">Email <span class="required">*</span></label>
        <input
          id="email"
          type="email"
          bind:value={email}
          placeholder="name@example.com"
          disabled={loading}
          on:input={clearError}
        />
      </div>

      <div class="form-group">
        <label for="password">Mật khẩu <span class="required">*</span></label>
        <div class="input-wrapper-with-icon">
          <input
            id="password"
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
              tabindex="-1"
            >
              {#if showPassword}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="eye-icon"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
                  /><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                  /></svg
                >
              {:else}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="eye-icon"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88"
                  /></svg
                >
              {/if}
            </Button>
          </div>
        </div>
      </div>

      <div class="form-group">
        <label for="confirmPassword"
          >Xác nhận mật khẩu <span class="required">*</span></label
        >
        <div class="input-wrapper-with-icon">
          <input
            id="confirmPassword"
            type={showConfirmPassword ? "text" : "password"}
            bind:value={confirmPassword}
            placeholder="Nhập lại mật khẩu"
            disabled={loading}
            class="has-icon-right"
            on:input={clearError}
          />
          <div class="toggle-password-wrapper">
            <Button
              variant="icon"
              type="button"
              onclick={toggleConfirmPassword}
              tabindex="-1"
            >
              {#if showConfirmPassword}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="eye-icon"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
                  /><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                  /></svg
                >
              {:else}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="eye-icon"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88"
                  /></svg
                >
              {/if}
            </Button>
          </div>
        </div>
      </div>

      <Button
        variant="primary"
        type="submit"
        disabled={loading}
        class="w-full mt-4"
      >
        {#if loading}
          <span class="loader"></span> Đang tạo tài khoản...
        {:else}
          Đăng ký
        {/if}
      </Button>
    </form>

    <div class="divider">
      <span>hoặc</span>
    </div>

    <Button
      variant="google"
      type="button"
      onclick={handleGoogleRegister}
      disabled={loading}
      class="w-full"
    >
      Đăng ký bằng Google
    </Button>

    <div class="auth-footer">
      Đã có tài khoản?
      <a href="/login" class="link-text" use:link>Đăng nhập</a>
    </div>
  {:else}
    <OtpVerification
      {email}
      on:success={handleOtpSuccess}
      on:back={() => (step = "register")}
    />
  {/if}
</AuthLayout>

<style>
  :root {
    --primary-blue: #1a56db;
    --text-dark: #111827;
    --text-gray: #6b7280;
    --input-bg: #f3f5f7;
  }

  /* Style cho dấu sao bắt buộc */
  .required {
    color: #ef4444; /* Màu đỏ */
    margin-left: 4px;
    font-weight: bold;
  }

  .auth-header {
    margin-bottom: 28px;
  }
  .icon-wrapper {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 16px;
  }
  .hat-icon {
    width: 48px;
    height: 48px;
    color: var(--primary-blue);
  }
  .title {
    font-size: 28px;
    font-weight: 700;
    color: var(--text-dark);
    margin: 0 0 8px 0;
  }
  .subtitle {
    font-size: 16px;
    color: var(--text-gray);
    margin: 0;
  }

  .auth-form {
    text-align: left;
  }
  .form-group {
    margin-bottom: 20px;
  }
  label {
    display: block;
    font-size: 15px;
    font-weight: 600;
    color: var(--text-dark);
    margin-bottom: 6px;
  }

  input {
    width: 100%;
    padding: 14px 16px;
    background-color: var(--input-bg);
    border: 2px solid transparent;
    border-radius: 12px;
    font-size: 16px;
    color: var(--text-dark);
    transition: all 0.2s;
    box-sizing: border-box;
  }
  input::placeholder {
    color: #9ca3af;
  }
  input:focus {
    outline: none;
    background-color: #fff;
    border-color: var(--primary-blue);
  }

  .input-wrapper-with-icon {
    position: relative;
  }
  input.has-icon-right {
    padding-right: 48px;
  }

  .toggle-password-wrapper {
    position: absolute;
    right: 8px;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
  }
  .eye-icon {
    width: 20px;
    height: 20px;
  }

  .w-full {
    width: 100% !important;
    display: flex !important;
  }
  .mt-4 {
    margin-top: 16px !important;
  }

  .auth-footer {
    margin-top: 24px;
    font-size: 16px;
    color: var(--text-dark);
  }
  .link-text {
    color: var(--primary-blue);
    font-weight: 700;
    margin-left: 4px;
    text-decoration: none;
  }
  .link-text:hover {
    text-decoration: underline;
  }

  .error-alert {
    background-color: #fee2e2;
    color: #991b1b;
    padding: 12px;
    border-radius: 8px;
    margin-bottom: 20px;
    font-size: 14px;
    text-align: left;
  }

  .loader {
    width: 20px;
    height: 20px;
    border: 3px solid #fff;
    border-bottom-color: transparent;
    border-radius: 50%;
    display: inline-block;
    animation: rotation 1s linear infinite;
  }
  @keyframes rotation {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .divider {
    display: flex;
    align-items: center;
    margin: 20px 0;
    color: #6b7280;
    font-size: 14px;
  }
  .divider::before,
  .divider::after {
    content: "";
    flex: 1;
    height: 1px;
    background-color: #e5e7eb;
  }
  .divider span {
    padding: 0 12px;
    font-weight: 500;
    color: #9ca3af;
    text-transform: lowercase;
  }
</style>
