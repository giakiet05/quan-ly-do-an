<script lang="ts">
    import {
        Plus,
        X,
        Trash2,
        Search,
        Download,
        Upload,
        Edit2 as Edit3,
    } from "../../../libs/Icons";
    import { mockStudents } from "../../../mocks/classes.mock";
    import type { ClassData, CreateClassRequest } from "../../../types/class";
    import { z } from "zod";
    import ClassInfoSidebar from "../../../components/ClassInfoSidebar.svelte";
    import StudentManagement from "../../../components/StudentManagement.svelte";
    import { classStore } from "../../../stores/class-store";
    import { getClassroom } from "../../../services/classroom-service";
    import type { Classroom } from "../../../models";
    import type { ClassroomResponse } from "../../../dtos";

    const { onClose, onSubmit, classData } = $props<{
        classData: ClassData;
        onClose: () => void;
        onSubmit: (data: CreateClassRequest) => void;
    }>();

    type ActiveTab = "list" | "excel" | "upload";

    interface Student {
        fullName: string;
        studentCode: string;
        email: string;
        selected: boolean;
    }

    const ClassSchema = z.object({
        name: z.string().min(1, "Tên lớp học là bắt buộc"),
        semester: z.string().min(1, "Học kỳ là bắt buộc"),
        avatar: z.string().optional(),
        description: z.string().optional(),
        students: z
            .array(
                z.object({
                    fullName: z.string(),
                    studentCode: z.string(),
                    email: z.string().optional(),
                }),
            )
            .optional(),
    });
    let detail = $state<ClassroomResponse>({} as ClassroomResponse);

    $effect(() => {
        getClassroom(classData.id)
            .then((res) => {
                detail = res; // Đảm bảo `res` có dữ liệu hợp lệ
                console.log(
                    "Fetched class detail:",
                    res.whitelist_student_code,
                );
            })
            .catch((err) => {
                console.error("Failed to fetch class detail:", err);
            });
    });

    let activeTab = $state<ActiveTab>("list");
    let errors = $state<Record<string, string>>({});

    let formData = $state<CreateClassRequest>({
        name: "",
        semester: "",
        avatar: "",
        description: "",
        students: [],
    });

    $effect(() => {
        formData.name = classData.name || "";
        formData.semester = classData.semester || "";
        formData.description = classData.description || "";
        formData.students = classData.students || [];
        formData.avatar = classData.avatar || "";
    });

    let students = $state<Student[]>(
        mockStudents.map((s) => ({
            ...s,
            selected: formData.students.some(
                (student) => student.studentCode === s.studentCode,
            ),
        })),
    );

    // Cập nhật formData.students khi selected thay đổi
    $effect(() => {
        formData.students = students
            .filter((s) => s.selected)
            .map(({ fullName, studentCode, email }) => ({
                fullName,
                studentCode,
                email: email || `${studentCode}@student.edu.vn`,
            }));
    });

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files?.[0]) return;

        const reader = new FileReader();
        reader.onload = (e) => {
            if (typeof e.target?.result === "string") {
                formData.avatar = e.target.result;
            }
        };
        reader.readAsDataURL(input.files[0]);
    }

    function handleSubmit(event: SubmitEvent) {
        event.preventDefault();
        errors = {};

        const result = ClassSchema.safeParse(formData);
        if (!result.success) {
            result.error.issues.forEach((issue) => {
                const field = issue.path[0] as string;
                if (field) errors[field] = issue.message;
            });
            return;
        }
        classStore
            .updateClass(classData.id, formData)
            .then(() => {
                onSubmit(formData);
                onClose();
            })
            .catch((err) => {
                console.error(err);
                alert("Lưu thay đổi thất bại!");
            });
    }

    function stopPropagation(event: MouseEvent) {
        event.stopPropagation();
    }

    function handleModalKeydown(event: KeyboardEvent) {
        event.stopPropagation();
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Escape") {
            onClose();
        }
        if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
        }
    }
</script>

<div
    class="modal-overlay"
    onclick={onClose}
    onkeydown={handleKeydown}
    role="button"
    tabindex="0"
>
    <div
        class="modal-content"
        onclick={stopPropagation}
        onkeydown={handleModalKeydown}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        tabindex="-1"
    >
        <div class="modal-header">
            <div>
                <h2 id="modal-title">Chỉnh Sửa Lớp Học</h2>
                <p>Chỉnh sửa thông tin lớp học và danh sách sinh viên</p>
            </div>
            <button class="close-btn" onclick={onClose}>
                <X size={24} />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="modal-body-wrapper">
            <div class="main-layout">
                <ClassInfoSidebar
                    bind:formData
                    {errors}
                    onAvatarChange={handleFileChange}
                />
                <StudentManagement
                    classDetail={detail}
                    {activeTab}
                    setActiveTab={(tab: ActiveTab) => (activeTab = tab)}
                />
            </div>

            <div class="modal-footer">
                <div class="stats">
                    Tổng cộng: <span>{formData.students.length}</span> sinh viên
                </div>
                <div class="actions">
                    <button type="button" class="btn-cancel" onclick={onClose}>
                        Hủy
                    </button>
                    <button type="submit" class="btn-submit">
                        Lưu thay đổi
                    </button>
                </div>
            </div>
        </form>
    </div>
</div>

<style>
    /* ──────────────── MODAL CONTAINER ──────────────── */
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.4);
        display: flex;
        align-items: flex-start;
        justify-content: center;
        z-index: 1000;
        padding: 40px 60px 20px 20px;
        overflow-y: auto;
    }

    .modal-content {
        background: white;
        border-radius: 12px;
        width: 100%;
        max-width: 950px;
        box-shadow:
            0 20px 25px -5px rgba(0, 0, 0, 0.1),
            0 10px 10px -5px rgba(0, 0, 0, 0.04);
        display: flex;
        flex-direction: column;
        animation: slideIn 0.3s ease-out;
    }

    @keyframes slideIn {
        from {
            transform: translateX(30px);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }

    /* Header */
    .modal-header {
        padding: 24px 32px;
        border-bottom: 1px solid #edf2f7;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .modal-header h2 {
        font-size: 24px;
        font-weight: 700;
        color: #1a202c;
        margin: 0;
    }
    .modal-header p {
        font-size: 14px;
        color: #718096;
        margin-top: 4px;
    }

    .close-btn {
        background: none;
        border: none;
        color: #a0aec0;
        cursor: pointer;
        padding: 8px;
        border-radius: 50%;
        transition: 0.2s;
    }
    .close-btn:hover {
        background: #f7fafc;
        color: #4a5568;
    }

    /* Layout chính */
    .main-layout {
        display: grid;
        grid-template-columns: 320px 1fr;
        min-height: 550px;
    }

    /* Footer */
    .modal-footer {
        padding: 24px 32px;
        border-top: 1px solid #edf2f7;
        display: flex;
        justify-content: space-between;
        align-items: center;
        background: white;
        border-radius: 0 0 12px 12px;
    }

    .stats {
        font-size: 14px;
        color: #718096;
    }
    .stats span {
        font-weight: 800;
        color: #0045b1;
    }

    .actions {
        display: flex;
        gap: 12px;
    }

    .btn-cancel {
        padding: 10px 24px;
        border-radius: 10px;
        font-weight: 600;
        color: #718096;
        background: #f7fafc;
        border: none;
        cursor: pointer;
    }

    .btn-submit {
        padding: 10px 24px;
        border-radius: 10px;
        font-weight: 600;
        color: white;
        background: #0045b1;
        border: none;
        cursor: pointer;
        box-shadow: 0 4px 14px rgba(0, 69, 177, 0.3);
    }

    .btn-submit:hover {
        background: #00358a;
        transform: translateY(-1px);
    }
</style>
