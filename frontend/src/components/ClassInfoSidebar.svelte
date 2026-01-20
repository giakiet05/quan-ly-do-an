<script lang="ts">
    import { Upload, Users, Calendar } from "../libs/Icons"; // Thêm icon cho đẹp
    import type { CreateClassRequest } from "../types/class";

    const {
        formData = $bindable(),
        errors = $bindable(),
        onAvatarChange,
    } = $props<{
        formData: CreateClassRequest;
        errors: Record<string, string>;
        onAvatarChange: (e: Event) => void;
    }>();

    // Danh sách năm để chọn cho nhanh
    const currentYear = new Date().getFullYear();
    const years = [currentYear - 1, currentYear, currentYear + 1];

    // Đảm bảo nếu năm của lớp cũ không có trong danh sách mặc định thì vẫn hiển thị được
    $effect(() => {
        if (formData.year && !years.includes(formData.year)) {
            years.push(formData.year);
            years.sort((a, b) => b - a);
        }
    });
</script>

<div class="sidebar">
    <div class="section-label">THÔNG TIN CƠ BẢN</div>

    <div class="avatar-upload-zone">
        <div class="avatar-preview">
            <img
                src={formData.avatar ||
                    `https://ui-avatars.com/api/?name=${formData.name || "C"}&background=0045B1&color=fff&size=128`}
                alt="Avatar lớp học"
            />
            <label class="upload-overlay">
                <Upload size={20} />
                <input
                    type="file"
                    accept="image/*"
                    class="hidden"
                    onchange={onAvatarChange}
                />
            </label>
        </div>
        <span>Ảnh đại diện lớp</span>
    </div>

    <div class="form-group">
        <label for="class-name">Tên lớp <span class="req">*</span></label>
        <input
            id="class-name"
            bind:value={formData.name}
            class:error={errors.name}
            placeholder="VD: Lớp K64 CNTT1"
        />
        {#if errors.name}<span class="err-text">{errors.name}</span>{/if}
    </div>

    <div class="grid-cols-2">
        <div class="form-group">
            <label for="semester">Học kỳ <span class="req">*</span></label>
            <select
                id="semester"
                bind:value={formData.semester}
                class:error={errors.semester}
            >
                <option value="HK1">Học kỳ 1</option>
                <option value="HK2">Học kỳ 2</option>
                <option value="HK3">Học kỳ 3</option>
            </select>
        </div>
        <div class="form-group">
            <label for="year">Năm học <span class="req">*</span></label>
            <select id="year" bind:value={formData.year}>
                {#each years as y}
                    <option value={y}>{y}</option>
                {/each}
            </select>
        </div>
    </div>

    <div class="form-group">
        <label for="max-students">Sĩ số tối đa</label>
        <div class="input-with-icon">
            <input
                id="max-students"
                type="number"
                bind:value={formData.maxStudents}
                min="1"
                max="200"
            />
        </div>
        <p class="helper-text">Giới hạn số lượng SV tham gia lớp</p>
    </div>

    <div class="toggle-group">
        <label class="switch">
            <input type="checkbox" bind:checked={formData.autoApprove} />
            <span class="slider"></span>
        </label>
        <div class="toggle-info">
            <span class="toggle-label">Tự động duyệt</span>
            <span class="toggle-desc">SV có mã sẽ vào thẳng lớp</span>
        </div>
    </div>

    <div class="form-group">
        <label for="description">Mô tả</label>
        <textarea
            id="description"
            bind:value={formData.description}
            placeholder="Mô tả ngắn về lớp học..."
        ></textarea>
    </div>
</div>

<style>
    /* ──────────────── SIDEBAR (CỘT TRÁI) ──────────────── */
    .sidebar {
        background: #f8fafc;
        padding: 32px;
        border-right: 1px solid #edf2f7;
    }

    .section-label {
        font-size: 11px;
        font-weight: 800;
        color: #4a5568;
        letter-spacing: 1px;
        margin-bottom: 24px;
    }

    .avatar-upload-zone {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;
        margin-bottom: 32px;
    }

    .avatar-preview {
        position: relative;
        width: 100px;
        height: 100px;
        border-radius: 20px;
        overflow: hidden;
        border: 4px solid white;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
    }

    .avatar-preview img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .upload-overlay {
        position: absolute;
        inset: 0;
        background: rgba(0, 0, 0, 0.4);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        opacity: 0;
        cursor: pointer;
        transition: 0.2s;
    }

    .avatar-preview:hover .upload-overlay {
        opacity: 1;
    }

    .avatar-upload-zone span {
        font-size: 12px;
        font-weight: 600;
        color: #718096;
    }

    .form-group {
        margin-bottom: 20px;
    }

    .form-group label {
        display: block;
        font-size: 13px;
        font-weight: 600;
        color: #4a5568;
        margin-bottom: 8px;
    }

    .form-group input,
    .form-group textarea {
        width: 100%;
        padding: 10px 14px;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        font-size: 14px;
        transition: 0.2s;
    }

    .form-group input:focus,
    .form-group textarea:focus {
        outline: none;
        border-color: #0045b1;
        box-shadow: 0 0 0 3px rgba(0, 69, 177, 0.1);
    }

    .form-group textarea {
        height: 100px;
        resize: none;
    }

    .req {
        color: #e53e3e;
    }

    .err-text {
        font-size: 12px;
        color: #e53e3e;
        margin-top: 4px;
        display: block;
    }
    .grid-cols-2 {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 12px;
    }

    .input-with-icon {
        position: relative;
    }

    .input-with-icon input {
        padding-left: 36px !important;
    }

    .helper-text {
        font-size: 11px;
        color: #718096;
        margin-top: 4px;
    }

    /* Style cho cái Toggle/Switch */
    .toggle-group {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 12px;
        background: white;
        border-radius: 12px;
        border: 1px solid #e2e8f0;
        margin-bottom: 20px;
    }

    .toggle-info {
        display: flex;
        flex-direction: column;
    }

    .toggle-label {
        font-size: 13px;
        font-weight: 700;
        color: #2d3748;
    }

    .toggle-desc {
        font-size: 11px;
        color: #718096;
    }

    .switch {
        position: relative;
        display: inline-block;
        width: 34px;
        height: 20px;
        flex-shrink: 0;
    }

    .switch input {
        opacity: 0;
        width: 0;
        height: 0;
    }

    .slider {
        position: absolute;
        cursor: pointer;
        inset: 0;
        background-color: #cbd5e0;
        transition: 0.4s;
        border-radius: 20px;
    }

    .slider:before {
        position: absolute;
        content: "";
        height: 14px;
        width: 14px;
        left: 3px;
        bottom: 3px;
        background-color: white;
        transition: 0.4s;
        border-radius: 50%;
    }

    input:checked + .slider {
        background-color: #0045b1;
    }
    input:checked + .slider:before {
        transform: translateX(14px);
    }

    select {
        width: 100%;
        padding: 10px 14px;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        background: white;
        font-size: 14px;
    }
</style>
