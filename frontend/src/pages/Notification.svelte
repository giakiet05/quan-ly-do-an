<script lang="ts">
    import { onMount } from "svelte";
    import NotificationList from "../components/NotificationList.svelte";
    import { notificationStore } from "../stores/notification-store";

    // Khai báo các biến trạng thái cho Lịch
    let selectedDate = new Date();
    let viewDate = new Date();

    // Lấy dữ liệu từ store
    const {
        notifications,
        markAllAsRead,
        fetchNotifications,
        markAsRead,
        remove,
    } = notificationStore;

    onMount(() => {
        fetchNotifications();
    });

    // Tối ưu hóa việc đánh dấu các ngày có thông báo trên lịch (dùng Set để tìm kiếm O(1))
    $: datesWithNotifications = new Set(
        $notifications.map((n) => new Date(n.createdAt).toDateString()),
    );

    // Logic tạo ma trận ngày trong tháng
    function getDaysInMonth(date: Date) {
        const year = date.getFullYear();
        const month = date.getMonth();
        const firstDay = new Date(year, month, 1).getDay();
        const daysInMonth = new Date(year, month + 1, 0).getDate();

        const days = [];
        // Căn chỉnh để Thứ 2 là ngày đầu tuần
        const offset = firstDay === 0 ? 6 : firstDay - 1;

        for (let i = 0; i < offset; i++) {
            days.push(null);
        }
        for (let i = 1; i <= daysInMonth; i++) {
            days.push(new Date(year, month, i));
        }
        return days;
    }

    $: days = getDaysInMonth(viewDate);

    // Chuyển tháng
    const nextMonth = () =>
        (viewDate = new Date(
            viewDate.getFullYear(),
            viewDate.getMonth() + 1,
            1,
        ));
    const prevMonth = () =>
        (viewDate = new Date(
            viewDate.getFullYear(),
            viewDate.getMonth() - 1,
            1,
        ));

    function selectDate(d: Date | null) {
        if (d) selectedDate = d;
    }

    // Xử lý sự kiện từ List bắn lên
    function handleMark(e: CustomEvent<{ id: string }>) {
        markAsRead(e.detail.id);
    }

    function handleDelete(e: CustomEvent<{ id: string }>) {
        remove(e.detail.id);
    }
</script>

<div class="notification-page">
    <main class="main-content">
        <header class="header-section">
            <div class="title-group">
                <h1>Thông báo</h1>
                <p class="count">
                    {selectedDate.toLocaleDateString("vi-VN", {
                        day: "numeric",
                        month: "long",
                        year: "numeric",
                    })}
                </p>
            </div>
            <button class="btn-mark-read" on:click={markAllAsRead}>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M18 6 7 17l-5-5" /><path
                        d="m22 10-7.5 7.5L13 16"
                    />
                </svg>
                Đánh dấu tất cả đã đọc
            </button>
        </header>

        <div class="list-container">
            <NotificationList
                filterDate={selectedDate}
                on:mark={handleMark}
                on:delete={handleDelete}
            />
        </div>
    </main>

    <aside class="sidebar">
        <div class="calendar-card">
            <div class="calendar-header">
                <h3>
                    {viewDate.toLocaleDateString("vi-VN", {
                        month: "long",
                        year: "numeric",
                    })}
                </h3>
                <div class="nav-btns">
                    <button on:click={prevMonth} aria-label="Tháng trước"
                        >&lt;</button
                    >
                    <button on:click={nextMonth} aria-label="Tháng sau"
                        >&gt;</button
                    >
                </div>
            </div>

            <div class="calendar-grid">
                {#each ["T2", "T3", "T4", "T5", "T6", "T7", "CN"] as day}
                    <span class="weekday">{day}</span>
                {/each}

                {#each days as d}
                    {#if d === null}
                        <div class="day empty"></div>
                    {:else}
                        <button
                            class="day"
                            class:active={selectedDate.toDateString() ===
                                d.toDateString()}
                            on:click={() => selectDate(d)}
                        >
                            {d.getDate()}
                            {#if datesWithNotifications.has(d.toDateString())}
                                <span class="dot"></span>
                            {/if}
                        </button>
                    {/if}
                {/each}
            </div>
        </div>

        <div class="info-box">
            <h4>Mẹo nhỏ</h4>
            <p>
                Chọn một ngày trên lịch để xem các hoạt động và thông báo cụ thể
                của ngày đó.
            </p>
        </div>
    </aside>
</div>

<style>
    .notification-page {
        display: grid;
        grid-template-columns: 1fr 340px;
        gap: 32px;
        padding: 40px;
        max-width: 1300px;
        margin: 0 auto;
        background-color: #f8fafc;
        min-height: 100vh;
    }

    .header-section {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;
    }

    h1 {
        font-size: 28px;
        font-weight: 800;
        margin: 0;
        color: #0f172a;
    }

    .count {
        color: #64748b;
        font-size: 14px;
        margin-top: 4px;
    }

    .btn-mark-read {
        display: flex;
        align-items: center;
        gap: 8px;
        background: white;
        border: 1px solid #e2e8f0;
        padding: 10px 16px;
        border-radius: 12px;
        font-weight: 600;
        color: #475569;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-mark-read:hover {
        background: #f1f5f9;
        border-color: #cbd5e1;
        color: #1a56db;
    }

    /* Lịch */
    .calendar-card {
        background: white;
        border-radius: 20px;
        padding: 24px;
        border: 1px solid #e2e8f0;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
        position: sticky;
        top: 40px;
    }

    .calendar-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;
    }

    .calendar-header h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 700;
        text-transform: capitalize;
    }

    .nav-btns button {
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        cursor: pointer;
        font-weight: bold;
        color: #64748b;
        padding: 4px 10px;
        transition: all 0.2s;
    }

    .nav-btns button:hover {
        background: #f1f5f9;
        color: #1a56db;
    }

    .calendar-grid {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        gap: 4px;
        text-align: center;
    }

    .weekday {
        font-size: 12px;
        font-weight: 700;
        color: #94a3b8;
        padding-bottom: 8px;
    }

    .day {
        aspect-ratio: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 10px;
        border: none;
        background: transparent;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        position: relative;
        transition: all 0.2s;
        color: #475569;
    }

    .day:hover:not(.active) {
        background: #f1f5f9;
        color: #1e293b;
    }

    .day.active {
        background: #1a56db;
        color: white;
        font-weight: 700;
        box-shadow: 0 4px 10px rgba(26, 86, 219, 0.3);
    }

    .dot {
        width: 5px;
        height: 5px;
        background: #ef4444;
        border-radius: 50%;
        position: absolute;
        bottom: 5px;
    }

    .day.active .dot {
        background: white;
    }

    .info-box {
        margin-top: 24px;
        padding: 20px;
        background: #eff6ff;
        border-radius: 16px;
        color: #1e40af;
        border: 1px solid #dbeafe;
    }

    .info-box h4 {
        margin: 0 0 8px 0;
        font-size: 14px;
        font-weight: 700;
    }

    .info-box p {
        margin: 0;
        font-size: 13px;
        line-height: 1.6;
        opacity: 0.9;
    }

    /* Responsive */
    @media (max-width: 1024px) {
        .notification-page {
            grid-template-columns: 1fr;
            padding: 20px;
        }
        .sidebar {
            order: -1; /* Đưa lịch lên đầu trên mobile */
        }
        .calendar-card {
            position: relative;
            top: 0;
        }
    }
</style>
