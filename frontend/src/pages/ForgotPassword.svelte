<script lang="ts">
    import AuthLayout from "../layouts/AuthLayout.svelte";
    import { link, push } from "svelte-spa-router";
    import Button from "../components/Button.svelte";
    import OtpVerification from "../components/OtpVerification.svelte";

    // --- STATE ---
    let step = "email"; // Các bước: 'email' -> 'otp' -> 'reset'
    let email = "";

    // State cho form Reset Password
    let newPassword = "";
    let confirmNewPassword = "";
    let showPassword = false;
    let showConfirmPassword = false;

    let loading = false;
    let error = "";

    // --- HANDLERS ---

    const clearError = () => {
        if (error) error = "";
    };
    const togglePassword = () => (showPassword = !showPassword);
    const toggleConfirmPassword = () =>
        (showConfirmPassword = !showConfirmPassword);

    // BƯỚC 1: Gửi yêu cầu reset
    const handleSendOtp = async () => {
        error = "";
        if (!email || !/\S+@\S+\.\S+/.test(email)) {
            error = "Vui lòng nhập địa chỉ email hợp lệ.";
            return;
        }

        loading = true;
        // Giả lập API gửi mail
        setTimeout(() => {
            console.log("Sending OTP to:", email);
            loading = false;
            step = "otp"; // Chuyển sang bước nhập OTP
        }, 1500);
    };

    // BƯỚC 2: Xử lý khi OTP thành công (nhận sự kiện từ component con)
    const handleOtpSuccess = (e) => {
        console.log("OTP Verified code:", e.detail.code);
        step = "reset"; // Chuyển sang bước đặt lại mật khẩu
    };

    // BƯỚC 3: Đổi mật khẩu mới
    const handleResetPassword = async () => {
        error = "";
        if (!newPassword || !confirmNewPassword) {
            error = "Vui lòng nhập đầy đủ thông tin.";
            return;
        }
        if (newPassword !== confirmNewPassword) {
            error = "Mật khẩu xác nhận không khớp.";
            return;
        }
        if (newPassword.length < 6) {
            error = "Mật khẩu phải có ít nhất 6 ký tự.";
            return;
        }

        loading = true;
        // Giả lập API đổi pass
        setTimeout(() => {
            loading = false;
            alert("Đổi mật khẩu thành công! Vui lòng đăng nhập lại.");
            push("/login"); // Chuyển hướng về trang Login
        }, 2000);
    };
</script>

<AuthLayout>
    {#if step === "email"}
        <div class="auth-header">
            <div class="icon-wrapper">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    class="key-icon"
                >
                    <path
                        fill-rule="evenodd"
                        d="M15.75 1.5a6.75 6.75 0 0 0-6.651 7.906c-1.067.322-2.02.854-2.841 1.551l-5.819 4.946a.75.75 0 0 0-.066 1.056l2.499 2.996a.75.75 0 0 0 .736.257l2.846-.711a.75.75 0 0 0 .524-.483l.865-2.597.575-.575a6.75 6.75 0 1 0 7.332-14.343ZM15.75 3a5.25 5.25 0 1 1 0 10.5 5.25 5.25 0 0 1 0-10.5Z"
                        clip-rule="evenodd"
                    />
                </svg>
            </div>
            <h2 class="title">Quên mật khẩu?</h2>
            <p class="subtitle">Đừng lo, hãy nhập email để lấy lại mật khẩu.</p>
        </div>

        <form on:submit|preventDefault={handleSendOtp} class="auth-form">
            {#if error}
                <div class="error-alert">{error}</div>
            {/if}

            <div class="form-group">
                <label for="email"
                    >Email đăng ký <span class="required">*</span></label
                >
                <input
                    id="email"
                    type="email"
                    bind:value={email}
                    placeholder="Nhập email của bạn"
                    disabled={loading}
                    on:input={clearError}
                />
            </div>

            <Button
                variant="primary"
                type="submit"
                disabled={loading}
                class="w-full"
            >
                {#if loading}
                    <span class="loader"></span> Đang gửi mã...
                {:else}
                    Gửi mã xác nhận
                {/if}
            </Button>
        </form>

        <div class="back-login">
            <a href="/login" class="back-link" use:link>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="2"
                    stroke="currentColor"
                    class="arrow-icon"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18"
                    />
                </svg>
                Quay lại đăng nhập
            </a>
        </div>
    {:else if step === "otp"}
        <OtpVerification on:success={handleOtpSuccess} />

        <div class="text-center mt-4">
            <button class="link-text-small" on:click={() => (step = "email")}>
                Thay đổi email?
            </button>
        </div>
    {:else if step === "reset"}
        <div class="auth-header">
            <div class="icon-wrapper">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    class="key-icon"
                >
                    <path
                        fill-rule="evenodd"
                        d="M12 1.5a5.25 5.25 0 0 0-5.25 5.25v3a3 3 0 0 0-3 3v6.75a3 3 0 0 0 3 3h10.5a3 3 0 0 0 3-3v-6.75a3 3 0 0 0-3-3v-3c0-2.9-2.35-5.25-5.25-5.25Zm3.75 8.25v-3a3.75 3.75 0 1 0-7.5 0v3h7.5Z"
                        clip-rule="evenodd"
                    />
                </svg>
            </div>
            <h2 class="title">Đặt lại mật khẩu</h2>
            <p class="subtitle">Tạo mật khẩu mới cho tài khoản của bạn.</p>
        </div>

        <form on:submit|preventDefault={handleResetPassword} class="auth-form">
            {#if error}
                <div class="error-alert">{error}</div>
            {/if}

            <div class="form-group">
                <label for="newpass"
                    >Mật khẩu mới <span class="required">*</span></label
                >
                <div class="input-wrapper-with-icon">
                    <input
                        id="newpass"
                        type={showPassword ? "text" : "password"}
                        bind:value={newPassword}
                        placeholder="Ít nhất 6 ký tự"
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
                <label for="confirmpass"
                    >Xác nhận mật khẩu <span class="required">*</span></label
                >
                <div class="input-wrapper-with-icon">
                    <input
                        id="confirmpass"
                        type={showConfirmPassword ? "text" : "password"}
                        bind:value={confirmNewPassword}
                        placeholder="Nhập lại mật khẩu mới"
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
                    <span class="loader"></span> Đang cập nhật...
                {:else}
                    Đổi mật khẩu
                {/if}
            </Button>
        </form>
    {/if}
</AuthLayout>

<style>
    /* --- Copy các biến và style chung từ Login/Register để đồng bộ --- */
    :root {
        --primary-blue: #1a56db;
        --text-dark: #111827;
        --text-gray: #6b7280;
        --input-bg: #f3f5f7;
    }
    .required {
        color: #ef4444;
        margin-left: 4px;
        font-weight: bold;
    }

    .auth-header {
        margin-bottom: 32px;
        text-align: center;
    } /* Căn giữa cho đẹp với form này */
    .icon-wrapper {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        margin-bottom: 16px;
        background: #eff6ff;
        padding: 12px;
        border-radius: 50%;
    }
    .key-icon {
        width: 32px;
        height: 32px;
        color: var(--primary-blue);
    }
    .title {
        font-size: 26px;
        font-weight: 700;
        color: var(--text-dark);
        margin: 0 0 8px 0;
    }
    .subtitle {
        font-size: 15px;
        color: var(--text-gray);
        margin: 0;
        line-height: 1.5;
    }

    .auth-form {
        text-align: left;
    }
    .form-group {
        margin-bottom: 24px;
    }
    label {
        display: block;
        font-size: 15px;
        font-weight: 600;
        color: var(--text-dark);
        margin-bottom: 8px;
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

    .back-login {
        text-align: center;
        margin-top: 32px;
    }
    .back-link {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        color: var(--text-gray);
        text-decoration: none;
        font-weight: 600;
        font-size: 15px;
        transition: color 0.2s;
    }
    .back-link:hover {
        color: var(--text-dark);
    }
    .arrow-icon {
        width: 18px;
        height: 18px;
        margin-right: 6px;
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

    .error-alert {
        background-color: #fee2e2;
        color: #991b1b;
        padding: 12px;
        border-radius: 8px;
        margin-bottom: 20px;
        font-size: 14px;
    }

    .w-full {
        width: 100% !important;
        display: flex !important;
    }
    .mt-4 {
        margin-top: 16px !important;
    }

    .loader {
        width: 18px;
        height: 18px;
        border: 2px solid #fff;
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

    .text-center {
        text-align: center;
    }
    .mt-4 {
        margin-top: 16px;
    }

    .link-text-small {
        background: none;
        border: none;
        color: var(--text-gray);
        font-size: 14px;
        text-decoration: underline;
        cursor: pointer;
        padding: 5px;
    }
    .link-text-small:hover {
        color: var(--text-dark);
    }
</style>
