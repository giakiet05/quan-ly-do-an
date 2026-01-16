<script lang="ts">
    import { notificationStore } from "../stores/notification-store";
    import NotificationItem from "./NotificationItem.svelte";
    import { createEventDispatcher } from "svelte";

    const { notifications } = notificationStore;
    const dispatch = createEventDispatcher();

    type FilterType =
        | "all"
        | "unread"
        | "message"
        | "deadline"
        | "submission"
        | "system";

    export let filter: FilterType = "all";
    export let filterDate: Date | null = null;

    // ✅ REACTIVE FILTER
    $: filteredList = (() => {
        let list = $notifications;

        if (filter === "unread") {
            list = list.filter((n) => !n.read);
        } else if (filter !== "all") {
            list = list.filter((n) => n.ui_type === filter);
        }

        if (filterDate) {
            const target = filterDate.toDateString();
            const today = new Date().toDateString();

            if (target !== today) {
                list = list.filter(
                    (n) => new Date(n.createdAt).toDateString() === target,
                );
            }
        }
        return list;
    })();

    function handleMark(e: CustomEvent<{ id: string }>) {
        dispatch("mark", e.detail); // Bắn tiếp lên Page
    }

    function handleDelete(e: CustomEvent<{ id: string }>) {
        dispatch("delete", e.detail); // Bắn tiếp lên Page
    }

    const filterOptions: { label: string; value: FilterType }[] = [
        { label: "Tất cả", value: "all" },
        { label: "Chưa đọc", value: "unread" },
        { label: "Tin nhắn", value: "message" },
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
        {#if filteredList.length === 0}
            <div class="empty-state">
                <p>Không có thông báo</p>
            </div>
        {:else}
            {#each filteredList as item (item.id)}
                <NotificationItem
                    {item}
                    on:mark={handleMark}
                    on:delete={handleDelete}
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
