<script lang="ts">
    import {
        notifications,
        markAsRead,
        removeNotification,
    } from "../stores/notification-store";
    import { derived } from "svelte/store";
    import NotificationItem from "./NotificationItem.svelte";

    // Khai báo kiểu dữ liệu cho filter
    type FilterType =
        | "all"
        | "unread"
        | "student"
        | "deadline"
        | "submission"
        | "system";

    export let filter: FilterType = "all";

    // Logic lọc thông báo (đã tối ưu)
    const filteredList = derived(notifications, ($n) => {
        if (filter === "all") return $n;
        if (filter === "unread") return $n.filter((x) => !x.read);
        if (filter === "student")
            return $n.filter((x) => x.type === "student_action");
        if (filter === "deadline")
            return $n.filter((x) => x.type === "deadline");
        if (filter === "submission")
            return $n.filter((x) => x.type === "submission");
        return $n.filter((x) => x.type === "system");
    });

    // Handlers
    function handleMark(e: CustomEvent<{ id: string }>) {
        markAsRead(e.detail.id);
    }
    function handleDelete(e: CustomEvent<{ id: string }>) {
        removeNotification(e.detail.id);
    }

    // Các nhãn hiển thị cho bộ lọc
    const filterOptions: { label: string; value: FilterType }[] = [
        { label: "Tất cả", value: "all" },
        { label: "Chưa đọc", value: "unread" },
        { label: "Sinh viên", value: "student" },
        { label: "Deadline", value: "deadline" },
        { label: "Bài nộp", value: "submission" },
        { label: "Hệ thống", value: "system" },
    ];
</script>

<div class="list-wrap">
    <div class="filters">
        {#each filterOptions as opt}
            <button
                class="chip"
                class:active={filter === opt.value}
                on:click={() => (filter = opt.value)}
            >
                {opt.label}
            </button>
        {/each}
    </div>

    <div class="cards-stack">
        {#if $filteredList.length === 0}
            <div class="empty-state">
                <div class="empty-icon">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="48"
                        height="48"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><path
                            d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"
                        /><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" /></svg
                    >
                </div>
                <p>Không có thông báo nào trong mục này</p>
            </div>
        {:else}
            {#each $filteredList as item (item.id)}
                <NotificationItem
                    {item}
                    on:mark={handleMark}
                    on:delete={handleDelete}
                    on:accept={handleMark}
                    on:reject={handleDelete}
                />
            {/each}
        {/if}
    </div>
</div>

<style>
    .list-wrap {
        width: 100%;
    }

    /* Bộ lọc Chips */
    .filters {
        display: flex;
        gap: 8px;
        margin-bottom: 24px;
        flex-wrap: wrap;
    }

    .chip {
        padding: 8px 18px;
        border-radius: 12px;
        background: white;
        border: 1px solid #e2e8f0;
        cursor: pointer;
        font-size: 14px;
        font-weight: 600;
        color: #64748b;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .chip:hover {
        background: #f1f5f9;
        border-color: #cbd5e1;
        color: #1e293b;
    }

    .chip.active {
        background: #1a56db;
        border-color: #1a56db;
        color: white;
        box-shadow: 0 4px 12px rgba(26, 86, 219, 0.25);
    }

    /* Stack danh sách */
    .cards-stack {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    /* Trạng thái trống */
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 60px 24px;
        text-align: center;
        background: white;
        border-radius: 16px;
        border: 1px dashed #cbd5e1;
    }

    .empty-icon {
        color: #cbd5e1;
        margin-bottom: 16px;
    }

    .empty-state p {
        color: #94a3b8;
        font-size: 15px;
        margin: 0;
    }

    /* Utility */
    :global(.cards-stack > *) {
        animation: slideUp 0.3s ease-out forwards;
    }

    @keyframes slideUp {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
