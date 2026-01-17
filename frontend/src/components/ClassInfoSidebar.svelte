<script lang="ts">
    import { Upload } from "../libs/Icons";
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
        {#if errors.name}
            <span class="err-text">{errors.name}</span>
        {/if}
    </div>

    <div class="form-group">
        <label for="semester">Học kỳ <span class="req">*</span></label>
        <input
            id="semester"
            bind:value={formData.semester}
            placeholder="VD: HK1"
            class:error={errors.semester}
        />
        {#if errors.semester}
            <span class="err-text">{errors.semester}</span>
        {/if}
    </div>

    <div class="form-group">
        <label for="description">Mô tả</label>
        <textarea
            id="description"
            bind:value={formData.description}
            placeholder="Mô tả ngắn về lớp học..."
        >
        </textarea>
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
</style>
