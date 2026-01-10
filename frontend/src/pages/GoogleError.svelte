<script lang="ts">
    import { onMount } from "svelte";
    import { push } from "svelte-spa-router";
    import Button from "../components/Button.svelte";

    let errorMessage = "Đã xảy ra lỗi trong quá trình xác thực.";

    onMount(() => {
        // Lấy message từ Hash (sau dấu ?)
        const hash = window.location.hash;
        const queryString = hash.includes("?") ? hash.split("?")[1] : "";
        const params = new URLSearchParams(queryString);
        const msg = params.get("message");

        if (msg) {
            errorMessage = decodeURIComponent(msg).replace(/\+/g, " ");
        }
    });
</script>

<div class="auth-header">
    <div class="icon-wrapper">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="error-icon"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
            />
        </svg>
    </div>
    <h2 class="title" style="color: #991b1b;">Lỗi đăng nhập</h2>
    <p class="subtitle">Hệ thống không thể hoàn tất xác thực</p>
</div>

<div class="error-container">
    <div class="error-alert">
        {errorMessage}
    </div>

    <Button
        variant="primary"
        class="w-full"
        onclick={() => push("/auth/login")}
    >
        Quay lại Đăng nhập
    </Button>
</div>

<style>
    .auth-header {
        text-align: center;
        margin-bottom: 24px;
    }
    .error-icon {
        width: 64px;
        height: 64px;
        color: #ef4444;
    }
    .title {
        font-size: 24px;
        font-weight: 700;
        margin: 16px 0 8px;
    }
    .subtitle {
        color: #6b7280;
        margin-bottom: 24px;
    }

    .error-alert {
        background-color: #fee2e2;
        color: #991b1b;
        padding: 16px;
        border-radius: 12px;
        margin-bottom: 24px;
        font-size: 15px;
        line-height: 1.5;
        border: 1px solid #fecaca;
        text-align: center;
    }

    :global(.w-full) {
        width: 100% !important;
    }
</style>
