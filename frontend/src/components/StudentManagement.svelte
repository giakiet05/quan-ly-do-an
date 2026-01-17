<script lang="ts">
    import { Upload } from "@lucide/svelte";
    import {
        Plus,
        Search,
        Trash2,
        Edit2 as Edit3,
        Download,
    } from "../libs/Icons";
    import { fade } from "svelte/transition";
    import { downloadWhitelistTemplate } from "../services/classroom-service";

    type ActiveTab = "list" | "excel" | "upload";

    interface Student {
        fullName: string;
        studentCode: string;
        email: string;
        selected: boolean;
    }

    let allCheckbox = $state<HTMLInputElement>();
    let uploadedFile: File | null = null;

    const props = $props<{
        activeTab: ActiveTab;
        setActiveTab: (tab: ActiveTab) => void;
        searchTerm: string;
        students: Student[];
        setSearchTerm: (value: string) => void;
        newStudent: { fullName: string; studentCode: string; email: string };
        selectedCount: number;
        filteredStudents: Student[];
        allSelected: boolean;
        indeterminate: boolean;
        toggleAll: () => void;
        addStudent: () => void;
        removeSelected: () => void;
    }>();

    async function downloadExcel() {
        try {
            await downloadWhitelistTemplate();
        } catch (err) {
            console.error(err);
            alert("Không thể tải file mẫu");
        }
    }
    const handleFileUpload = () => {
        if (uploadedFile) {
            const reader = new FileReader();
            reader.onload = (event) => {
                const content = event.target?.result;
                console.log("Uploaded file content:", content);
                // Process the uploaded file content here
            };
            reader.readAsText(uploadedFile);
        }
    };

    $effect(() => {
        if (allCheckbox) {
            allCheckbox.indeterminate = props.indeterminate;
        }
    });
</script>

<div class="content-area">
    <div class="tabs">
        <button
            type="button"
            class:active={props.activeTab === "list"}
            onclick={() => props.setActiveTab("list")}
        >
            Thêm thủ công
        </button>
        <button
            type="button"
            class:active={props.activeTab === "excel"}
            onclick={() => props.setActiveTab("excel")}
        >
            Tải Excel
        </button>
        <button
            type="button"
            class:active={props.activeTab === "upload"}
            onclick={() => props.setActiveTab("upload")}
        >
            Upload File
        </button>
    </div>

    {#if props.activeTab === "list"}
        <div class="student-manager">
            <div class="quick-add">
                <input
                    placeholder="MSSV"
                    bind:value={props.newStudent.studentCode}
                />
                <input
                    placeholder="Họ tên"
                    bind:value={props.newStudent.fullName}
                />
                <button
                    type="button"
                    class="add-btn"
                    onclick={props.addStudent}
                >
                    <Plus size={20} />
                </button>
            </div>

            <div class="search-bar">
                <Search size={16} />
                <input
                    placeholder="Tìm theo tên hoặc MSSV..."
                    value={props.searchTerm}
                    oninput={(e) => props.setSearchTerm(e.currentTarget.value)}
                />
            </div>

            <div class="table-wrapper">
                {#if props.selectedCount > 0}
                    <div
                        class="bulk-actions-bar"
                        transition:fade={{ duration: 150 }}
                    >
                        <div class="selected-info">
                            <span class="count">{props.selectedCount}</span> đã chọn
                        </div>
                        <div class="action-group">
                            <button
                                type="button"
                                class="btn-bulk-edit"
                                disabled={props.selectedCount != 1}
                            >
                                <Edit3 size={14} /> Sửa
                            </button>
                            <button
                                type="button"
                                class="btn-bulk-delete"
                                onclick={props.removeSelected}
                            >
                                <Trash2 size={14} /> Xóa
                            </button>
                        </div>
                    </div>
                {/if}

                <div class="table-container">
                    <table>
                        <thead>
                            <tr>
                                <th>
                                    <input
                                        type="checkbox"
                                        bind:this={allCheckbox}
                                        checked={props.allSelected}
                                        onchange={props.toggleAll}
                                    />
                                </th>
                                <th>MSSV</th>
                                <th>Họ tên</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each props.filteredStudents as student (student.studentCode)}
                                <tr class:selected={student.selected}>
                                    <td>
                                        <input
                                            type="checkbox"
                                            bind:checked={student.selected}
                                        />
                                    </td>
                                    <td class="code">{student.studentCode}</td>
                                    <td class="name">{student.fullName}</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>

                    {#if props.filteredStudents.length === 0}
                        <div class="empty-table">
                            Không tìm thấy sinh viên nào
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {:else if props.activeTab === "excel"}
        <div class="excel-tab">
            <button type="button" onclick={downloadExcel}>
                <Download size={20} /> Tải xuống mẫu Excel
            </button>
        </div>
    {:else if props.activeTab === "upload"}
        <div class="upload-tab">
            <!-- <input type="file" accept=".csv, .xlsx" bind:files={uploadedFile} /> -->
            <button type="button" onclick={handleFileUpload}>
                <Upload size={20} /> Upload
            </button>
        </div>
    {/if}
</div>

<style>
    /* ──────────────── CONTENT AREA (CỘT PHẢI) ──────────────── */
    .content-area {
        padding: 32px;
        display: flex;
        flex-direction: column;
    }

    /* Tabs */
    .tabs {
        display: flex;
        gap: 24px;
        border-bottom: 2px solid #edf2f7;
        margin-bottom: 24px;
    }

    .tabs button {
        padding: 12px 4px;
        border: none;
        background: none;
        font-size: 14px;
        font-weight: 600;
        color: #718096;
        cursor: pointer;
        border-bottom: 2px solid transparent;
        transition: 0.2s;
        margin-bottom: -2px;
    }

    .tabs button.active {
        color: #0045b1;
        border-bottom-color: #0045b1;
    }

    /* Student Manager */
    .table-wrapper {
        position: relative;
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .quick-add {
        display: flex;
        gap: 12px;
        margin-bottom: 16px;
    }

    .quick-add input {
        flex: 1;
        padding: 10px;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        font-size: 14px;
    }

    .add-btn {
        background: #0045b1;
        color: white;
        border: none;
        padding: 0 16px;
        border-radius: 8px;
        cursor: pointer;
    }

    .add-btn:hover {
        background: #00358a;
    }

    .search-bar {
        position: relative;
        margin-bottom: 16px;
    }

    .search-bar :global(svg) {
        position: absolute;
        left: 12px;
        top: 50%;
        transform: translateY(-50%);
        color: #a0aec0;
    }

    .search-bar input {
        width: 100%;
        padding: 10px 10px 10px 40px;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        font-size: 13px;
        background: #f8fafc;
    }

    .bulk-actions-bar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        background: #0045b1;
        padding: 8px 16px;
        border-radius: 8px;
        color: white;
        box-shadow: 0 4px 12px rgba(0, 69, 177, 0.25);
    }

    .selected-info {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
    }

    .selected-info .count {
        background: white;
        color: #0045b1;
        padding: 2px 8px;
        border-radius: 4px;
        font-weight: 800;
    }

    .action-group {
        display: flex;
        gap: 8px;
    }

    .action-group button {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 6px 12px;
        border-radius: 6px;
        font-size: 12px;
        font-weight: 600;
        border: none;
        cursor: pointer;
        transition: 0.2s;
    }

    .btn-bulk-edit {
        background: rgba(255, 255, 255, 0.2);
        color: white;
    }
    .btn-bulk-edit:hover {
        background: rgba(255, 255, 255, 0.3);
    }

    .btn-bulk-delete {
        background: #ff4d4f;
        color: white;
    }
    .btn-bulk-delete:hover {
        background: #ff7875;
    }

    .table-container {
        flex: 1;
        border: 1px solid #edf2f7;
        border-radius: 12px;
        overflow: hidden;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        font-size: 14px;
    }

    thead {
        background: #f8fafc;
        border-bottom: 1px solid #edf2f7;
    }

    th {
        text-align: left;
        padding: 12px 16px;
        font-weight: 700;
        color: #4a5568;
        font-size: 12px;
        text-transform: uppercase;
    }

    td {
        padding: 12px 16px;
        border-bottom: 1px solid #f7fafc;
    }

    tr.selected {
        background: #f0f7ff;
    }

    .code {
        font-family: monospace;
        font-weight: 600;
        color: #2d3748;
    }

    .name {
        color: #4a5568;
    }

    .empty-state {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 16px;
        color: #a0aec0;
    }

    .empty-table {
        padding: 40px;
        text-align: center;
        color: #718096;
        font-style: italic;
    }

    input[type="checkbox"] {
        accent-color: #0045b1;
    }

    .btn-bulk-edit[disabled] {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .excel-tab,
    .upload-tab {
        display: flex;
        flex-direction: column;
        gap: 16px;
        align-items: center;
        justify-content: center;
        padding: 32px;
        border: 1px solid #edf2f7;
        border-radius: 12px;
        background: #f8fafc;
    }

    .excel-tab button,
    .upload-tab button {
        background: #0045b1;
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 14px;
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .excel-tab button:hover,
    .upload-tab button:hover {
        background: #00358a;
    }

    .upload-tab input {
        padding: 8px;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        font-size: 14px;
    }
</style>
