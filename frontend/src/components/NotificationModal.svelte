<script lang="ts">
    import NotificationList from "./NotificationList.svelte";
    import { writable } from "svelte/store";
    import { markAllAsRead, notifications } from "../stores/notification-store";
    import { derived } from "svelte/store";

    export let open = writable(false);

    // Tính toán số thông báo chưa đọc từ store
    $: unreadCount = $notifications.filter((n) => !n.read).length;

    function close() {
        open.set(false);
    }
</script>

{#if $open}
    <div class="backdrop" on:click={close}></div>
    <div class="modal" on:click|stopPropagation>
        <header class="modal-header">
            <div class="title-group">
                <h3>Thông báo</h3>
                {#if unreadCount > 0}
                    <span class="unread-badge">{unreadCount} mới</span>
                {/if}
            </div>

            <div class="actions">
                <button
                    class="btn-mark-all"
                    on:click={markAllAsRead}
                    title="Đánh dấu tất cả đã đọc"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><path d="M18 6 7 17l-5-5" /><path
                            d="m22 10-7.5 7.5L13 16"
                        /></svg
                    >
                </button>
                <button class="close-btn" on:click={close}>
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><line x1="18" y1="6" x2="6" y2="18"></line><line
                            x1="6"
                            y1="6"
                            x2="18"
                            y2="18"
                        ></line></svg
                    >
                </button>
            </div>
        </header>

        <div class="modal-body">
            <NotificationList />
        </div>

        <footer class="modal-footer">
            <a href="/notifications" on:click={close}>Xem tất cả thông báo</a>
        </footer>
    </div>
{/if}

<style>
    /* Lớp nền mờ */
    .backdrop {
        position: fixed;
        inset: 0;
        background: rgba(
            15,
            23,
            42,
            0.1
        ); /* Màu tối nhẹ hơn để trông hiện đại */
        backdrop-filter: blur(2px); /* Làm mờ nhẹ phía sau */
        z-index: 40;
    }

    /* Modal dạng Dropdown */
    .modal {
        position: fixed;
        right: 24px;
        top: 80px;
        width: 400px;
        background: white;
        border-radius: 20px;
        display: flex;
        flex-direction: column;
        box-shadow:
            0 20px 25px -5px rgba(0, 0, 0, 0.1),
            0 8px 10px -6px rgba(0, 0, 0, 0.1);
        border: 1px solid #f1f5f9;
        z-index: 50;
        overflow: hidden;
        animation: slideIn 0.2s ease-out;
    }

    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateY(-10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    /* Header của Modal */
    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 20px;
        border-bottom: 1px solid #f1f5f9;
    }

    .title-group {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .title-group h3 {
        margin: 0;
        font-size: 18px;
        font-weight: 800;
        color: #0f172a;
    }

    .unread-badge {
        background: #eff6ff;
        color: #2563eb;
        font-size: 12px;
        font-weight: 700;
        padding: 2px 8px;
        border-radius: 99px;
    }

    .actions {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .btn-mark-all,
    .close-btn {
        background: transparent;
        border: none;
        color: #64748b;
        cursor: pointer;
        padding: 4px;
        display: flex;
        align-items: center;
        border-radius: 6px;
        transition: all 0.2s;
    }

    .btn-mark-all:hover {
        background: #f1f5f9;
        color: #1a56db;
    }

    .close-btn:hover {
        background: #fef2f2;
        color: #ef4444;
    }

    /* Thân Modal (Vùng cuộn) */
    .modal-body {
        max-height: 480px; /* Độ cao tối đa để không che hết màn hình */
        overflow-y: auto;
        background: #fcfdfe; /* Nền hơi xanh nhẹ giống trang chính */
    }

    /* Tùy chỉnh thanh cuộn cho đẹp */
    .modal-body::-webkit-scrollbar {
        width: 6px;
    }
    .modal-body::-webkit-scrollbar-thumb {
        background: #e2e8f0;
        border-radius: 10px;
    }

    /* Footer Modal */
    .modal-footer {
        padding: 16px;
        text-align: center;
        border-top: 1px solid #f1f5f9;
        background: white;
    }

    .modal-footer a {
        font-size: 14px;
        font-weight: 600;
        color: #2563eb;
        text-decoration: none;
        transition: color 0.2s;
    }

    .modal-footer a:hover {
        color: #1e429f;
        text-decoration: underline;
    }
</style>
