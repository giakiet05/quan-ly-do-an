<script lang="ts">
    import { onMount } from "svelte";
    import { push, querystring } from "svelte-spa-router";
    import Button from "../components/Button.svelte";
    import { authService } from "../services/auth-service";

    let fullName = "";
    let studentCode = "";
    let setupToken = "";
    let loading = false;
    let error = "";

    $: if ($querystring) {
        const params = new URLSearchParams($querystring);
        setupToken = params.get("setup_token") || "";
    }

    onMount(() => {
        if (!setupToken) {
            const hashParts = window.location.hash.split("?");
            if (hashParts.length > 1) {
                const params = new URLSearchParams(hashParts[1]);
                setupToken = params.get("setup_token") || "";
            }
        }

        if (!setupToken) {
            error = "Không tìm thấy mã thiết lập. Vui lòng đăng nhập lại.";
        }
    });

    const clearError = () => {
        if (error) error = "";
    };

    const handleComplete = async () => {
        error = "";
        if (!fullName) {
            error = "Vui lòng nhập họ tên của bạn.";
            return;
        }

        loading = true;
        try {
            const res = await authService.completeGoogleSetup({
                setupToken: setupToken,
                fullName: fullName,
                //studentCode: studentCode,
            });

            if (res?.accessToken) {
                // Xóa token trên URL để bảo mật sau khi thành công
                window.history.replaceState(
                    {},
                    document.title,
                    window.location.pathname +
                        window.location.hash.split("?")[0],
                );
                push("/home");
            }
        } catch (err: any) {
            console.error("Lỗi hoàn tất đăng ký:", err);
            error =
                err?.message || "Hoàn tất đăng ký thất bại. Vui lòng thử lại.";
        } finally {
            loading = false;
        }
    };
</script>

<div class="auth-header">
    <div class="icon-wrapper">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="currentColor"
            class="hat-icon"
        >
            <path
                d="M11.7 2.805a.75.75 0 0 1 .6 0l9.75 4.2a.75.75 0 0 1 0 1.39l-9.75 4.2a.75.75 0 0 1-.6 0l-9.75-4.2a.75.75 0 0 1 0-1.39l9.75-4.2ZM4.32 7.5 12 10.8l7.68-3.3-7.68-3.3-7.68 3.3ZM12 13.2l5.556-2.393a3.75 3.75 0 0 1 1.694 3.332v3.181a.75.75 0 0 1-1.5 0v-3.181c0-.725-.43-1.382-1.096-1.668L12 14.587l-4.654-2.017a1.875 1.875 0 0 0-2.192.334L3.75 15.086v.756a3.75 3.75 0 1 0 7.5 0v-.756l-1.404-2.106c-.366-.548-1.058-.78-1.654-.522L12 13.2Z"
            />
        </svg>
    </div>
    <h2 class="title">Hoàn tất đăng ký</h2>
    <p class="subtitle">Bổ sung thông tin để bắt đầu sử dụng hệ thống</p>
</div>

<div class="auth-form">
    {#if error}
        <div class="error-alert">{error}</div>
    {/if}

    {#if setupToken}
        <div class="form-group">
            <label for="fullName"
                >Họ và tên <span class="required">*</span></label
            >
            <input
                id="fullName"
                type="text"
                bind:value={fullName}
                placeholder="Nhập họ và tên của bạn"
                disabled={loading}
                on:input={clearError}
            />
        </div>

        <div class="form-group">
            <label for="studentCode">Mã sinh viên</label>
            <input
                id="studentCode"
                type="text"
                bind:value={studentCode}
                placeholder="Nhập mã sinh viên (nếu có)"
                disabled={loading}
            />
        </div>

        <Button
            variant="primary"
            type="button"
            disabled={loading}
            class="w-full"
            onclick={handleComplete}
        >
            {#if loading}
                <span class="loader"></span> Đang xử lý...
            {:else}
                Hoàn tất đăng ký
            {/if}
        </Button>
    {:else}
        <div style="text-align: center; padding: 20px 0;">
            <p style="color: #6b7280; margin-bottom: 20px;">
                Đang xác thực thông tin từ Google...
            </p>
            <a href="/#/auth/login" class="link-text">Quay lại đăng nhập</a>
        </div>
    {/if}
</div>

<style>
    /* Header section */
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
        margin-bottom: 24px;
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
    input:focus {
        outline: none;
        background-color: #fff;
        border-color: #1a56db;
    }
    input:disabled {
        opacity: 0.6;
        cursor: not-allowed;
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

    .link-text {
        color: #1a56db;
        font-weight: 700;
        text-decoration: none;
    }
    .link-text:hover {
        text-decoration: underline;
    }
</style>
