<script lang="ts">
    import { push, querystring } from "svelte-spa-router";
    import Button from "../../components/Button.svelte";

    // Khai báo các loại lỗi thân thiện với người dùng
    const errorMap: Record<string, { title: string; hint: string }> = {
        missing_auth_code: {
            title: "Thiếu mã xác thực",
            hint: "Không nhận được phản hồi từ Google. Vui lòng thử lại.",
        },
        mismatch_state: {
            title: "Lỗi bảo mật (State)",
            hint: "Phiên làm việc đã hết hạn hoặc yêu cầu bị giả mạo. Hãy đăng nhập lại.",
        },
        google_api_error: {
            title: "Lỗi kết nối Google",
            hint: "Hệ thống không thể lấy thông tin từ Google.",
        },
        unauthorized_domain: {
            title: "Tài khoản không hợp lệ",
            hint: "Vui lòng sử dụng email sinh viên (@gm.uit.edu.vn).",
        },
    };

    let errorTitle = "Lỗi xác thực";
    let displayMessage = "Đã xảy ra lỗi không xác định.";

    // Reactive: Tự động chạy lại mỗi khi query trên URL thay đổi
    $: if ($querystring) {
        const params = new URLSearchParams($querystring);
        const code = params.get("message"); // Đây là mã lỗi từ Backend gửi về (ví dụ: missing_auth_code)

        if (code) {
            // Nếu mã lỗi nằm trong danh sách định nghĩa sẵn
            if (errorMap[code]) {
                errorTitle = errorMap[code].title;
                displayMessage = errorMap[code].hint;
            } else {
                // Nếu là message tự do từ Backend gửi xuống (đã encode)
                errorTitle = "Thông báo lỗi";
                displayMessage = decodeURIComponent(code).replace(/\+/g, " ");
            }
        }
    }
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
    <h2 class="title" style="color: #991b1b;">{errorTitle}</h2>
    <p class="subtitle">Hệ thống không thể hoàn tất xác thực</p>
</div>

<div class="error-container">
    <div class="error-alert">
        {displayMessage}
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
    /* Giữ nguyên style của bạn */
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
</style>
