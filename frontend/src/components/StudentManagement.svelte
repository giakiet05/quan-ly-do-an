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
    let manualCode = $state("");
    let domainInput = $state("");
    let localWhitelist = $state<string[]>([]);

    const {
        classDetail,
        activeTab,
        setActiveTab,
        uploadedFile,
        setUploadedFile,
        onWhitelistChange,
        onDomainChange,
        initialDomain,
    } = $props<{
        classDetail?: ClassroomResponse;
        activeTab: ActiveTab;
        setActiveTab: (tab: ActiveTab) => void;
        uploadedFile?: File | null;
        setUploadedFile?: (file: File) => void;
        onWhitelistChange?: (whitelist: string[]) => void;
        onDomainChange?: (domain: string) => void;
        initialDomain?: string;
    }>();
    const whitelistStudents = $derived(
        classDetail?.whitelistStudentCode?.map(
            (entry: { studentCode: string }) => entry.studentCode,
        ) ?? [],
    );

    // Sync data when component mounts or classDetail changes
    $effect(() => {
        if (whitelistStudents.length > 0 && localWhitelist.length === 0) {
            localWhitelist = [...whitelistStudents];
        }
    });

    // Initialize domain input
    $effect(() => {
        if (initialDomain && !domainInput) {
            domainInput = initialDomain;
        }
    });

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

    function mergeUploadedToWhitelist() {
        // Thêm uploadedData vào localWhitelist (avoid duplicates)
        const newCodes = uploadedData.filter(
            (code) => !localWhitelist.includes(code),
        );
        if (newCodes.length > 0) {
            localWhitelist = [...localWhitelist, ...newCodes];
            onWhitelistChange?.(localWhitelist);
            console.log("✅ Merged uploaded data:", localWhitelist);
        }
    }
    function addManualStudent() {
        const code = manualCode.trim().toUpperCase();
        if (code && !localWhitelist.includes(code)) {
            localWhitelist = [...localWhitelist, code];
            onWhitelistChange?.(localWhitelist);
            manualCode = "";
        }
    }

    function removeStudent(code: string) {
        localWhitelist = localWhitelist.filter((c) => c !== code);
        onWhitelistChange?.(localWhitelist);
    }

    function updateDomain(value: string) {
        domainInput = value;
        onDomainChange?.(value);
    }
</script>

<div class="content-area">
    <div class="tabs">
        <button
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
            <div class="helper-note">
                <div class="helper-icon">i</div>
                <p>
                    <strong>Cơ chế Whitelist:</strong> Khi bật, chỉ sinh viên có
                    MSSV trong danh sách này mới có thể vào lớp. Bạn có thể giới
                    hạn email theo domain (ví dụ: @hcmute.edu.vn) bên dưới.
                </p>
            </div>

            <div class="config-row">
                <div class="input-group">
                    <label for="domain-email">Giới hạn Domain Email</label>
                    <div class="flex-row">
                        <input
                            id="domain-email"
                            type="text"
                            placeholder="VD: uit.edu.vn"
                            value={domainInput}
                            onchange={(e) =>
                                updateDomain(
                                    (e.target as HTMLInputElement).value,
                                )}
                        />
                    </div>
                </div>
            </div>

            <div class="quick-add">
                <div class="input-with-icon">
                    <span class="icon-wrapper">
                        <Search size={18} />
                    </span>
                    <input
                        type="text"
                        placeholder="Nhập MSSV thủ công..."
                        bind:value={manualCode}
                        onkeydown={(e) =>
                            e.key === "Enter" && addManualStudent()}
                    />
                </div>
                <button
                    type="button"
                    class="add-btn"
                    onclick={addManualStudent}
                >
                    <Plus size={18} /> Thêm
                </button>
            </div>

            <div class="table-wrapper">
                <div class="table-container">
                    <table>
                        <thead>
                            <tr>
                                <th>MSSV</th>
                                <th style="text-align: right;">Hành động</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each localWhitelist as studentCode (studentCode)}
                                <tr>
                                    <td class="code">{studentCode}</td>
                                    <td style="text-align: right;">
                                        <button
                                            class="delete-btn"
                                            onclick={() =>
                                                removeStudent(studentCode)}
                                        >
                                            <Trash2 size={16} />
                                        </button>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>

                    {#if localWhitelist.length === 0}
                        <div class="empty-table">
                            <div class="empty-state">
                                <p>
                                    Chưa có sinh viên nào trong danh sách
                                    Whitelist
                                </p>
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
                        <div class="upload-actions">
                            <button
                                type="button"
                                class="confirm-btn"
                                onclick={mergeUploadedToWhitelist}
                            >
                                ✓ Xác nhận thêm
                            </button>
                        </div>
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
        align-items: stretch; /* Đảm bảo con cao bằng nhau */
        width: 100%;
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
        padding: 0 20px;
        border-radius: 8px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        height: 42px;
        font-size: 14px;
        font-weight: 600;
        white-space: nowrap;
        transition: background 0.2s;
    }

    .add-btn:hover {
        background: #00358a;
    }
    .helper-note {
        display: flex;
        gap: 12px;
        background: #f0f7ff;
        border: 1px solid #bee3f8;
        padding: 12px 16px;
        border-radius: 8px;
        margin-bottom: 20px;
    }

    .helper-icon {
        background: #3182ce;
        color: white;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        font-weight: bold;
        flex-shrink: 0;
    }

    .helper-note p {
        margin: 0;
        font-size: 13px;
        color: #2c5282;
        line-height: 1.4;
    }

    /* Email Config */
    .config-row {
        margin-bottom: 20px;
    }

    .input-group label {
        display: block;
        font-size: 13px;
        font-weight: 700;
        color: #4a5568;
        margin-bottom: 6px;
    }

    .input-group input {
        width: 100%;
        padding: 10px;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        font-size: 14px;
        background: #fff;
    }

    .input-with-icon {
        position: relative;
        flex: 1;
        display: flex;
        align-items: center;
    }

    .icon-wrapper {
        position: absolute;
        left: 14px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #a0aec0;
        pointer-events: none; /* Không cho icon cản trở click vào input */
        z-index: 10;
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

    .code {
        font-family: monospace;
        font-weight: 600;
        color: #2d3748;
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

    .upload-actions {
        display: flex;
        justify-content: flex-end;
        padding: 12px 16px;
        border-top: 1px solid #edf2f7;
        background: #f8fafc;
    }

    .confirm-btn {
        background: #0045b1;
        color: white;
        border: none;
        padding: 8px 16px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 13px;
        font-weight: 600;
        transition: background 0.2s;
    }

    .confirm-btn:hover {
        background: #00358a;
    }

    .input-with-icon input {
        width: 100%;
        height: 42px; /* Cố định chiều cao */
        padding: 0 12px 0 40px; /* Padding trái né cái Icon ra */
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        font-size: 14px;
        box-sizing: border-box; /* Quan trọng để padding không làm tăng size input */
    }
    .flex-row {
        display: flex;
        gap: 8px;
    }

    .delete-btn {
        background: none;
        border: none;
        color: #e53e3e;
        cursor: pointer;
        padding: 4px;
        border-radius: 4px;
    }

    .delete-btn:hover {
        background: #fff5f5;
    }
</style>
