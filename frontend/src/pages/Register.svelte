<script lang="ts">
  import { link, replace } from "svelte-spa-router";
  import OtpVerification from "../components/OTPVerification.svelte";
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

  const togglePassword = () => (showPassword = !showPassword);
  const toggleConfirmPassword = () =>
    (showConfirmPassword = !showConfirmPassword);

  const clearError = () => {
    if (error) error = "";
  };

  const handleGoogleRegister = () => {
    authService.loginWithGoogle();
  };

  const handleOtpSuccess = async (
    event: CustomEvent<{ verification_token: string }>,
  ) => {
    loading = true;
    error = "";

    try {
      const input_token = event.detail.verification_token;
      const res = await authService.verifyEmail({ email, otp: input_token });
      const verificationToken = res.verificationToken || input_token;
      await completeRegistration({
        verificationToken,
        fullName: fullname,
        password,
      });
      replace("/auth/login");
    } catch (e: any) {
      error = e?.message || "Đăng ký thất bại";
      step = "verify";
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
    <p class="subtitle">Tham gia hệ thống quản lý đồ án ngay</p>
  </div>

  <form on:submit|preventDefault={handleSubmit} class="auth-form">
    {#if error}
      <div class="error-alert">{error}</div>
    {/if}

    <div class="form-group">
      <label for="fullname">Tên hiển thị <span class="required">*</span></label>
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
        <button
          type="button"
          class="icon-button"
          on:click|preventDefault={togglePassword}
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
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
              /><path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
              />
            </svg>
          {:else}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="eye-icon"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88"
              />
            </svg>
          {/if}
        </button>
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
        <button
          type="button"
          class="icon-button"
          on:click|preventDefault={toggleConfirmPassword}
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
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
              /><path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
              />
            </svg>
          {:else}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="eye-icon"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88"
              />
            </svg>
          {/if}
        </button>
      </div>
    </div>

    <Button variant="primary" type="submit" disabled={loading} class="w-full">
      {#if loading}
        <span class="loader"></span> Đang xử lý...
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
    <a href="/auth/login" class="link-text" use:link>Đăng nhập</a>
  </div>
{:else}
  <OtpVerification
    {email}
    bind:error
    on:success={handleOtpSuccess}
    on:back={() => (step = "register")}
  />
{/if}

<style>
  /* Header section - Căn giữa đồng bộ */
  .auth-header {
    margin-bottom: 32px;
    text-align: center;
  }

  .icon-wrapper {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 16px;
  }

  .hat-icon {
    width: 56px;
    height: 56px;
    color: #1a56db;
  }

  .title {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0 0 8px 0;
  }

  .subtitle {
    font-size: 16px;
    color: #6b7280;
    margin: 0;
  }

  /* Form styles */
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
    color: #111827;
    margin-bottom: 8px;
  }

  .required {
    color: #ef4444;
    margin-left: 2px;
  }

  input {
    width: 100%;
    padding: 14px 16px;
    background-color: #f9fafb;
    border: 2px solid #e5e7eb;
    border-radius: 12px;
    font-size: 15px;
    color: #111827;
    transition: all 0.2s;
    box-sizing: border-box;
  }

  input::placeholder {
    color: #9ca3af;
  }

  input:focus {
    outline: none;
    background-color: #fff;
    border-color: #1a56db;
  }

  input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  /* Password toggle */
  .input-wrapper-with-icon {
    position: relative;
  }

  input.has-icon-right {
    padding-right: 48px;
  }

  .icon-button {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    background: transparent;
    border: none;
    padding: 0;
    margin: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    width: 32px;
    height: 32px;
    color: #6b7280;
    transition: color 0.2s;
  }

  .icon-button:hover {
    color: #1a56db;
  }

  .eye-icon {
    width: 20px;
    height: 20px;
  }

  /* Error alert */
  .error-alert {
    background-color: #fee2e2;
    color: #991b1b;
    padding: 12px 16px;
    border-radius: 10px;
    margin-bottom: 20px;
    font-size: 14px;
    border-left: 4px solid #ef4444;
  }

  /* Loader */
  .loader {
    width: 18px;
    height: 18px;
    border: 3px solid #fff;
    border-bottom-color: transparent;
    border-radius: 50%;
    display: inline-block;
    animation: rotation 1s linear infinite;
    margin-right: 8px;
  }

  @keyframes rotation {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  /* Divider */
  .divider {
    display: flex;
    align-items: center;
    margin: 24px 0;
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
    padding: 0 16px;
    font-weight: 500;
    color: #9ca3af;
  }

  /* Footer */
  .auth-footer {
    margin-top: 28px;
    text-align: center;
    font-size: 15px;
    color: #6b7280;
  }

  .link-text {
    color: #1a56db;
    font-weight: 700;
    text-decoration: none;
    transition: color 0.2s;
    margin-left: 4px;
  }

  .link-text:hover {
    color: #1e40af;
    text-decoration: underline;
  }
  .w-full {
    width: 100% !important;
    display: flex !important;
  }
  /* Utility */
</style>
