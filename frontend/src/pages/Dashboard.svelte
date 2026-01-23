<script lang="ts">
    import { onMount } from "svelte";
    import { authStore } from "../stores/auth-store";
    import { classStore } from "../stores/class-store";
    import { notificationStore } from "../stores/notification-store";
    import { BookOpen, GraduationCap, Bell, ArrowRight } from "lucide-svelte";

    // Lấy dữ liệu từ các store
    const { classes } = classStore;
    const { unreadCount } = notificationStore;

    let enrolledCount = $state(0);
    let loading = $state(true);

    onMount(async () => {
        try {
            await Promise.all([
                classStore.fetchMyClasses(),
                notificationStore.loadUnreadCount(),
            ]);
            // Giả lập dữ liệu lớp đang học của sinh viên
            enrolledCount = 3;
        } catch (e) {
            console.error("Lỗi load dashboard:", e);
        } finally {
            loading = false;
        }
    });

    const displayName = $derived($authStore.user?.fullname || "Người dùng");
</script>

<div class="dashboard-wrapper">
    {#if loading}
        <div class="loading-state">
            <div class="spinner"></div>
            <p>Đang tải dữ liệu...</p>
        </div>
    {:else}
        <div class="welcome-banner">
            <div class="banner-content">
                <div class="banner-text">
                    <h1>Xin chào, {displayName} 👋</h1>
                    <p>
                        Hệ thống đã sẵn sàng. Bạn có thể <strong
                            >quản lý giảng dạy</strong
                        >
                        và
                        <strong>theo dõi học tập</strong> ngay trên cùng một tài
                        khoản.
                    </p>
                </div>
                <div class="role-badges">
                    <span class="badge gv">Giảng viên</span>
                    <span class="badge sv">Sinh viên</span>
                </div>
            </div>
        </div>

        <div class="stats-grid">
            <div class="stat-card gv-theme">
                <div class="card-header">
                    <div class="icon-box"><GraduationCap size={24} /></div>
                    <span class="role-tag">QUẢN TRỊ VIÊN</span>
                </div>
                <div class="card-body">
                    <span class="number">{$classes.length}</span>
                    <span class="label">Lớp học bạn đang dạy</span>
                </div>
                <button class="action-btn">
                    Vào trang quản lý <ArrowRight size={16} />
                </button>
            </div>

            <div class="stat-card sv-theme">
                <div class="card-header">
                    <div class="icon-box"><BookOpen size={24} /></div>
                    <span class="role-tag">HỌC VIÊN</span>
                </div>
                <div class="card-body">
                    <span class="number">{enrolledCount}</span>
                    <span class="label">Lớp học bạn đang tham gia</span>
                </div>
                <button class="action-btn">
                    Xem lớp đang học <ArrowRight size={16} />
                </button>
            </div>

            <div class="stat-card notification-theme">
                <div class="card-header">
                    <div class="icon-box"><Bell size={24} /></div>
                    <span class="role-tag">HỆ THỐNG</span>
                </div>
                <div class="card-body">
                    <span class="number">{$unreadCount}</span>
                    <span class="label">Thông báo chưa xem</span>
                </div>
                <button class="action-btn">
                    Mở hộp thư <ArrowRight size={16} />
                </button>
            </div>
        </div>

        <div class="quick-guide">
            <div class="guide-item">
                <span class="dot gv"></span>
                <p>Màu xanh dương: Chức năng cho <strong>Giảng viên</strong></p>
            </div>
            <div class="guide-item">
                <span class="dot sv"></span>
                <p>Màu xanh lá: Chức năng cho <strong>Sinh viên</strong></p>
            </div>
        </div>
    {/if}
</div>

<style>
    .dashboard-wrapper {
        padding: 2rem;
        max-width: 1240px;
        margin: 0 auto;
        font-family: "Inter", system-ui, sans-serif;
    }

    /* Banner */
    .welcome-banner {
        background: linear-gradient(135deg, #1e40af 0%, #1a56db 100%);
        border-radius: 20px;
        padding: 3rem 2rem;
        color: white;
        margin-bottom: 2.5rem;
        box-shadow: 0 10px 25px -5px rgba(26, 86, 219, 0.3);
    }

    .banner-content {
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-wrap: wrap;
        gap: 1.5rem;
    }

    .banner-text h1 {
        font-size: 2rem;
        font-weight: 700;
        margin-bottom: 0.75rem;
    }
    .banner-text p {
        font-size: 1.1rem;
        opacity: 0.9;
        max-width: 600px;
        line-height: 1.6;
    }

    .role-badges {
        display: flex;
        gap: 0.75rem;
    }
    .badge {
        padding: 8px 16px;
        border-radius: 99px;
        font-size: 0.85rem;
        font-weight: 600;
        border: 1px solid rgba(255, 255, 255, 0.4);
        background: rgba(255, 255, 255, 0.1);
    }

    /* Stats Grid */
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
        gap: 1.5rem;
    }

    .stat-card {
        background: white;
        border-radius: 18px;
        padding: 2rem;
        border: 1px solid #e2e8f0;
        display: flex;
        flex-direction: column;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;
    }

    .stat-card:hover {
        transform: translateY(-8px);
        box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.05);
    }

    .card-header {
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-bottom: 2rem;
    }

    .icon-box {
        width: 50px;
        height: 50px;
        border-radius: 14px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .role-tag {
        font-size: 0.75rem;
        font-weight: 800;
        letter-spacing: 0.1em;
        color: #64748b;
    }

    .card-body {
        flex-grow: 1;
    }
    .number {
        font-size: 3.5rem;
        font-weight: 800;
        color: #0f172a;
        line-height: 1;
    }
    .label {
        display: block;
        color: #64748b;
        margin-top: 0.75rem;
        font-size: 1rem;
        font-weight: 500;
    }

    .action-btn {
        margin-top: 2rem;
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 0.95rem;
        font-weight: 600;
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
        transition: gap 0.2s;
    }

    .action-btn:hover {
        gap: 12px;
    }

    /* Themes */
    .gv-theme .icon-box {
        background: #eff6ff;
        color: #2563eb;
    }
    .gv-theme .action-btn {
        color: #2563eb;
    }
    .gv-theme:hover {
        border-color: #2563eb;
    }

    .sv-theme .icon-box {
        background: #f0fdf4;
        color: #16a34a;
    }
    .sv-theme .action-btn {
        color: #16a34a;
    }
    .sv-theme:hover {
        border-color: #16a34a;
    }

    .notification-theme .icon-box {
        background: #f5f3ff;
        color: #7c3aed;
    }
    .notification-theme .action-btn {
        color: #7c3aed;
    }
    .notification-theme:hover {
        border-color: #7c3aed;
    }

    /* Guide */
    .quick-guide {
        margin-top: 3rem;
        display: flex;
        gap: 2rem;
        padding: 1.5rem;
        background: #f8fafc;
        border-radius: 12px;
        border: 1px dashed #cbd5e1;
    }

    .guide-item {
        display: flex;
        align-items: center;
        gap: 10px;
    }
    .dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
    }
    .dot.gv {
        background: #2563eb;
    }
    .dot.sv {
        background: #16a34a;
    }
    .guide-item p {
        font-size: 0.9rem;
        color: #475569;
        margin: 0;
    }

    /* Loading */
    .loading-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 400px;
        color: #64748b;
    }

    .spinner {
        width: 40px;
        height: 40px;
        border: 4px solid #f3f3f3;
        border-top: 4px solid #1a56db;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 1rem;
    }

    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }

    @media (max-width: 640px) {
        .dashboard-wrapper {
            padding: 1rem;
        }
        .welcome-banner {
            padding: 2rem 1.5rem;
        }
        .banner-text h1 {
            font-size: 1.5rem;
        }
        .quick-guide {
            flex-direction: column;
            gap: 1rem;
        }
    }
</style>
