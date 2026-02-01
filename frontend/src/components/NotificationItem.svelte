<script lang="ts">
    import { createEventDispatcher, onMount } from "svelte";
    import type { NotificationItem } from "../types/notification";
    import {
        getClassroomByChannelId,
        getClassroom,
    } from "../services/classroom-service";
    import {
        getGroup,
        acceptGroupInvitation,
        rejectGroupInvitation,
    } from "../services/group-service";
    import Swal from "sweetalert2";
    import { push } from "svelte-spa-router";
    import { getProject } from "../services/project-service";

    export let item: NotificationItem;

    const dispatch = createEventDispatcher();

    // --- LOGIC XỬ LÝ MỚI ---
    let displayContent = item.content;
    let targetClassId = ""; // Lưu Classroom ID để điều hướng
    let isLoadingClass = false;
    const isGroupInvite = item.content.includes("mời tham gia nhóm");

    onMount(async () => {
        // Regex tìm ID 24 ký tự ở cuối chuỗi
        console.log("Item trong onMount:", item);
        const match = item.content.match(/(.*)\s([a-f\d]{24})$/i);

        if (item.ui_type === "message" && match) {
            const prefix = match[1].trim(); // "Tin nhắn mới từ"
            const channelId = match[2];

            try {
                isLoadingClass = true;
                const res = await getClassroomByChannelId(channelId);
                if (res) {
                    displayContent = `${prefix} lớp ${res.name}`;
                    targetClassId = res.id; // ID thật của lớp học
                }
            } catch (e) {
                console.error("Lỗi lấy thông tin lớp:", e);
            } finally {
                isLoadingClass = false;
            }
        }
    });

    const handleItemClick = () => {
        if (isGroupInvite) {
            handleGroupInviteClick();
        } else if (targetClassId) {
            push(`/lecture/my-classes/${targetClassId}`);
            if (!item.read) doMark();
        }
    };
    // -----------------------
    const handleGroupInviteClick = async () => {
        const link = item.link || "";
        const regex =
            /\/classrooms\/([a-f\d]{24})\/projects\/([a-f\d]{24})\/groups\/([a-f\d]{24})\/invitations\/([a-f\d]{24})/;
        const match = link.match(regex);

        if (!match) return;

        const [_, classroomId, projectId, groupId, invitationId] = match;

        try {
            // Sử dụng Promise.all để gọi 2 API song song, tiết kiệm thời gian chờ
            const [classroom, project] = await Promise.all([
                getClassroom(classroomId),
                getProject(classroomId, projectId),
            ]);
            const result = await Swal.fire({
                title: '<span style="font-size: 20px; font-weight: 800; color: #1e293b;">Lời mời tham gia nhóm</span>',
                html: `
                <div style="margin-top: 15px; text-align: left; font-family: 'Inter', sans-serif;">
                    <div style="background: #f1f5f9; padding: 12px 16px; border-radius: 12px; margin-bottom: 12px; border-left: 4px solid #64748b;">
                        <div style="font-size: 11px; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 4px;">Lớp học</div>
                        <div style="font-size: 15px; font-weight: 600; color: #334155;">${classroom.name}</div>
                    </div>

                    <div style="background: #f5f3ff; padding: 12px 16px; border-radius: 12px; margin-bottom: 16px; border-left: 4px solid #7c3aed;">
                        <div style="font-size: 11px; font-weight: 700; color: #7c3aed; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 4px;">Đề tài dự án</div>
                        <div style="font-size: 15px; font-weight: 600; color: #4c1d95;">${project.title}</div>
                    </div>

                    <div style="padding: 0 4px;">
                        <p style="font-size: 14px; color: #475569; line-height: 1.6; margin: 0;">
                            <span style="color: #1e293b; font-weight: 600;">Thông báo:</span> 
                            ${item.content.split(" ").slice(0, -1).join(" ")}
                        </p>
                    </div>
                </div>
            `,
                icon: "info",
                iconColor: "#7c3aed",
                showCancelButton: true,
                confirmButtonText: "Đồng ý tham gia",
                cancelButtonText: "Để sau",
                confirmButtonColor: "#7c3aed",
                cancelButtonColor: "#f1f5f9",
                customClass: {
                    confirmButton: "swal-confirm-btn",
                    cancelButton: "swal-cancel-btn",
                },
            });

            if (result.isConfirmed) {
                await acceptGroupInvitation(groupId, invitationId);
                Swal.fire({
                    title: "Thành công!",
                    text: `Bạn đã tham gia nhóm của đề tài: ${project.title}`,
                    icon: "success",
                    timer: 2000,
                    showConfirmButton: false,
                });
                doMark();
                push(`/lecture/my-classes/${classroomId}`);
            }
        } catch (error: any) {
            console.error("Lỗi lấy thông tin:", error);
            Swal.fire(
                "Lỗi",
                "Không thể tải thông tin lớp học hoặc đề tài này.",
                "error",
            );
        }
    };
    const doMark = () => {
        dispatch("mark", { id: item.id });
    };

    const doDelete = () => {
        dispatch("delete", { id: item.id });
    };

    function formatTime(dateStr: string) {
        const date = new Date(dateStr);
        const now = new Date();

        if (date.toDateString() === now.toDateString()) {
            return date.toLocaleTimeString("vi-VN", {
                hour: "2-digit",
                minute: "2-digit",
            });
        }
        if (date.getFullYear() === now.getFullYear()) {
            return date.toLocaleDateString("vi-VN", {
                day: "numeric",
                month: "short",
            });
        }
        return date.toLocaleDateString("vi-VN");
    }
</script>

<div
    class="notification-item"
    on:click={handleItemClick}
    style:cursor={targetClassId ? "pointer" : "default"}
>
    <div class="icon-wrapper {item.ui_type}">
        {#if item.ui_type === "message"}
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
                ><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" /><circle
                    cx="12"
                    cy="7"
                    r="4"
                /></svg
            >
        {:else if item.ui_type === "deadline"}
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
                ><circle cx="12" cy="12" r="10" /><polyline
                    points="12 6 12 12 16 14"
                /></svg
            >
        {:else if item.ui_type === "submission"}
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
                ><path
                    d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                /><polyline points="14 2 14 8 20 8" /></svg
            >
        {:else if item.ui_type === "class"}
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
                ><path d="M22 10v6M2 10v6M6 6h12M6 18h12" /></svg
            >
        {:else}
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
                ><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" /><path
                    d="M13.73 21a2 2 0 0 1-3.46 0"
                /></svg
            >
        {/if}
    </div>

    <div class="content-body">
        <div class="header">
            <span class="title">{item.title}</span>
            {#if !item.read}
                <span class="dot-unread"></span>
            {/if}
        </div>

        <p class="description">
            {#if isLoadingClass}
                <span style="color: #94a3b8; font-style: italic;"
                    >Đang xác định lớp học...</span
                >
            {:else}
                {displayContent}
            {/if}
        </p>

        {#if item.meta}
            <div class="meta-info">{item.meta}</div>
        {/if}
    </div>

    <div class="side-panel">
        <span class="time">
            {formatTime(item.createdAt)}
        </span>

        <div class="tool-panel">
            {#if !item.read}
                <button
                    class="tool-btn mark"
                    on:click|stopPropagation={doMark}
                    title="Đánh dấu đã đọc"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><polyline points="20 6 9 17 4 12"></polyline></svg
                    >
                </button>
            {/if}
            <button
                class="tool-btn delete"
                on:click|stopPropagation={doDelete}
                title="Xóa"
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
                    ><path d="M3 6h18" /><path
                        d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"
                    /><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" /></svg
                >
            </button>
        </div>
    </div>
</div>

<style>
    /* TOÀN BỘ CSS CỦA BẠN GIỮ NGUYÊN */
    .notification-item {
        display: flex;
        gap: 16px;
        padding: 18px 20px;
        background: white;
        border-radius: 16px;
        border: 1px solid #f1f5f9;
        position: relative;
        transition: all 0.2s ease;
        margin-bottom: 12px;
    }

    .notification-item:hover {
        border-color: #cbd5e1;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
        transform: translateY(-1px);
    }

    /* Icon */
    .icon-wrapper {
        width: 42px;
        height: 42px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }
    .icon-wrapper.message {
        background: #eff6ff;
        color: #2563eb;
    }
    .icon-wrapper.deadline {
        background: #fff1f2;
        color: #e11d48;
    }
    .icon-wrapper.submission {
        background: #f0fdf4;
        color: #16a34a;
    }
    .icon-wrapper.system {
        background: #f8fafc;
        color: #64748b;
    }
    .icon-wrapper.class {
        background: #fefce8;
        color: #ca8a04;
    }

    .content-body {
        flex: 1;
        min-width: 0;
    }
    .header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 4px;
    }
    .title {
        font-weight: 700;
        color: #1e293b;
        font-size: 15px;
    }
    .description {
        margin: 0;
        font-size: 14px;
        color: #475569;
        line-height: 1.5;
    }
    .meta-info {
        margin-top: 6px;
        font-size: 12px;
        color: #94a3b8;
    }

    .side-panel {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        justify-content: space-between;
        min-width: 80px;
        flex-shrink: 0;
    }

    .time {
        font-size: 12px;
        color: #94a3b8;
        white-space: nowrap;
    }

    .tool-panel {
        display: flex;
        gap: 4px;
        opacity: 0;
        transition: opacity 0.2s ease;
        margin-top: 8px;
    }

    .notification-item:hover .tool-panel {
        opacity: 1;
    }

    .tool-btn {
        width: 28px;
        height: 28px;
        border-radius: 8px;
        border: none;
        background: #f8fafc;
        color: #94a3b8;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
    }
    .tool-btn:hover {
        background: #f1f5f9;
        color: #1e293b;
    }
    .tool-btn.delete:hover {
        background: #fef2f2;
        color: #ef4444;
    }

    .dot-unread {
        width: 8px;
        height: 8px;
        background: #3b82f6;
        border-radius: 50%;
        flex-shrink: 0;
    }
    @media (max-width: 768px) {
        .tool-panel {
            opacity: 1;
        }
    }
    .icon-wrapper.group {
        background: #f5f3ff;
        color: #7c3aed;
    }

    .invite-bg {
        border-left: 4px solid #7c3aed !important;
    }

    .action-hint {
        display: block;
        font-size: 12px;
        color: #7c3aed;
        margin-top: 4px;
        font-weight: 500;
    }
</style>
