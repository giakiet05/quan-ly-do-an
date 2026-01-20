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
    import * as XLSX from "xlsx";
    import type { ClassroomResponse } from "../dtos";

    type ActiveTab = "list" | "excel" | "upload";

    let uploadedData = $state<string[]>([]);

    const {
        classDetail,
        activeTab,
        setActiveTab,
        uploadedFile,
        setUploadedFile,
    } = $props<{
        classDetail?: ClassroomResponse;
        activeTab: ActiveTab;
        setActiveTab: (tab: ActiveTab) => void;
        uploadedFile?: File | null;
        setUploadedFile?: (file: File) => void;
    }>();
    const whitelistStudents = $derived(classDetail?.whitelistStudentCode ?? []);

    async function downloadExcel() {
        try {
            await downloadWhitelistTemplate();
        } catch (err) {
            console.error(err);
            alert("Không thể tải file mẫu");
        }
    }

    const loadFileUpload = () => {
        const input = document.createElement("input");
        input.type = "file";
        input.accept = ".xlsx, .xls";
        input.onchange = async (e: Event) => {
            const target = e.target as HTMLInputElement;
            if (target.files && target.files.length > 0) {
                const file = target.files[0];
                setUploadedFile(file);

                const data = await file.arrayBuffer();
                const workbook = XLSX.read(data, { type: "array" });
                const firstSheetName = workbook.SheetNames[0];
                const worksheet = workbook.Sheets[firstSheetName];
                const jsonData = XLSX.utils.sheet_to_json(worksheet, {
                    header: 1,
                });

                // Extract only the "Student Code" column
                uploadedData = jsonData
                    .slice(1)
                    .map((row: any) => row[0])
                    .filter(Boolean);
            }
        };
        input.click();
    };

    const clearUploadedFile = () => {
        setUploadedFile(null);
        uploadedData = [];
    };
</script>

<div class="content-area">
    <div class="tabs">
        <button
            hidden={!classDetail}
            type="button"
            class:active={activeTab === "list"}
            onclick={() => setActiveTab("list")}
        >
            Danh sách sinh viên
        </button>
        <button
            type="button"
            class:active={activeTab === "excel"}
            onclick={() => setActiveTab("excel")}
        >
            Tải Excel
        </button>
        <button
            type="button"
            class:active={activeTab === "upload"}
            onclick={() => setActiveTab("upload")}
        >
            Upload Excel
        </button>
    </div>
    {#if activeTab === "list"}
        <div class="student-manager">
            <!-- <div class="search-bar">
                <Search size={18} />
                <input
                    type="text"
                    placeholder="Tìm kiếm MSSV..."
                    bind:value={searchTerm}
                />
            </div> -->

            <div class="table-wrapper">
                <div class="table-container">
                    <table>
                        <thead>
                            <tr>
                                <th>MSSV</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each whitelistStudents as studentCode (studentCode)}
                                <tr>
                                    <td class="code"
                                        >{studentCode.studentCode}</td
                                    >
                                </tr>
                            {/each}
                        </tbody>
                    </table>

                    {#if whitelistStudents.length === 0}
                        <div class="empty-table">
                            <!-- {searchTerm
                                ? "Không tìm thấy MSSV nào khớp"
                                : "Lớp học chưa có sinh viên nào trong danh sách"} -->
                            <div class="empty-state">
                                <p>Không có sinh viên nào trong danh sách</p>
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {:else if activeTab === "excel"}
        <div class="excel-tab">
            <button type="button" onclick={downloadExcel}>
                <Download size={20} /> Tải xuống mẫu Excel
            </button>
        </div>
    {:else if activeTab === "upload"}
        <div class="upload-tab">
            {#if uploadedFile}
                <div class="uploaded-file">
                    <span>{uploadedFile.name}</span>
                    <button type="button" onclick={loadFileUpload}>
                        Thay file
                    </button>
                    <button type="button" onclick={clearUploadedFile}>
                        Xóa
                    </button>
                </div>
                {#if uploadedData.length > 0}
                    <div class="table-container">
                        <table>
                            <thead>
                                <tr>
                                    <th>Student Code</th>
                                </tr>
                            </thead>
                            <tbody>
                                {#each uploadedData as studentCode}
                                    <tr>
                                        <td>{studentCode}</td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                {:else}
                    <p class="empty-table">Không có dữ liệu để hiển thị.</p>
                {/if}
            {:else}
                <button type="button" onclick={loadFileUpload}>
                    <Upload size={20} /> Upload
                </button>
            {/if}
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
        margin-top: 16px;
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
        margin-top: 16px;
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

    .uploaded-file {
        display: flex;
        align-items: center;
        gap: 16px;
        font-size: 14px;
        color: #4a5568;
    }

    .uploaded-file button {
        background: #0045b1;
        color: white;
        border: none;
        padding: 6px 12px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 12px;
    }

    .uploaded-file button:hover {
        background: #00358a;
    }
</style>
