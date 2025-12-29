<script lang="ts">
    import { X } from "lucide-svelte";
    // Định nghĩa Interface (giữ nguyên logic từ React)
    interface GroupReport {
        id: string;
        groupName: string;
        groupNumber: number;
        members: string[];
        submittedDate: string | null;
        status: "submitted" | "pending" | "late";
        grade: number | null;
        reportFile?: string;
        feedback?: string;
    }

    // Khai báo Props trong Svelte 5
    let { report, onClose, onSubmit } = $props<{
        report: GroupReport;
        onClose: () => void;
        onSubmit: (reportId: string, grade: number, feedback: string) => void;
    }>();

    // State (Runes)
    let grade = $state(report.grade?.toString() || "");
    let feedback = $state(report.feedback || "");
    let errors = $state<Record<string, string>>({});

    const quickFeedbacks = [
        "Báo cáo tốt, nội dung đầy đủ",
        "Cần bổ sung thêm phần phân tích",
        "Trình bày rõ ràng, logic tốt",
        "Thiếu một số nội dung quan trọng",
        "Xuất sắc! Tiếp tục phát huy",
        "Nộp muộn, cần cải thiện thời gian",
    ];

    const gradeScales = [
        {
            range: "9-10",
            label: "Xuất sắc",
            color: "bg-green-100 text-green-700",
        },
        { range: "8-8.9", label: "Giỏi", color: "bg-blue-100 text-blue-700" },
        { range: "7-7.9", label: "Khá", color: "bg-cyan-100 text-cyan-700" },
        {
            range: "5-6.9",
            label: "Trung bình",
            color: "bg-yellow-100 text-yellow-700",
        },
        { range: "0-4.9", label: "Yếu", color: "bg-red-100 text-red-700" },
    ];

    function handleSubmit(e: Event) {
        e.preventDefault();
        const newErrors: Record<string, string> = {};

        if (!grade.trim()) {
            newErrors.grade = "Vui lòng nhập điểm";
        } else {
            const gradeNum = parseFloat(grade);
            if (isNaN(gradeNum) || gradeNum < 0 || gradeNum > 10) {
                newErrors.grade = "Điểm phải từ 0 đến 10";
            }
        }

        if (Object.keys(newErrors).length > 0) {
            errors = newErrors;
            return;
        }

        onSubmit(report.id, parseFloat(grade), feedback);
    }

    function handleQuickFeedback(text: string) {
        if (feedback && !feedback.endsWith(". ")) {
            feedback = `${feedback}. ${text}`;
        } else {
            feedback = feedback + text;
        }
    }
</script>

<div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    role="dialog"
    aria-modal="true"
>
    <div
        class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto"
    >
        <div
            class="flex justify-between items-center p-6 border-b border-gray-200 sticky top-0 bg-white"
        >
            <h2 class="text-xl font-semibold">Chấm điểm báo cáo</h2>
            <button
                onclick={onClose}
                class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
                <X size={20} />
            </button>
        </div>

        <form onsubmit={handleSubmit} class="p-6">
            <div class="bg-gray-50 rounded-lg p-4 mb-6">
                <div class="mb-2">
                    <span class="text-sm text-gray-600">Nhóm:</span>
                    <span class="ml-2 font-medium"
                        >{report.groupName} (Nhóm {report.groupNumber})</span
                    >
                </div>
                <div class="mb-2">
                    <span class="text-sm text-gray-600">Thành viên:</span>
                    <span class="ml-2">{report.members.join(", ")}</span>
                </div>
                {#if report.submittedDate}
                    <div>
                        <span class="text-sm text-gray-600">Thời gian nộp:</span
                        >
                        <span class="ml-2"
                            >{new Date(report.submittedDate).toLocaleString(
                                "vi-VN",
                            )}</span
                        >
                    </div>
                {/if}
                {#if report.status === "late"}
                    <div class="mt-2 text-sm text-orange-600 font-medium">
                        ⚠️ Báo cáo này được nộp muộn
                    </div>
                {/if}
            </div>

            <div class="mb-6">
                <label class="block text-sm mb-2" for="grade-input">
                    Điểm số <span class="text-red-500">*</span>
                </label>
                <div class="flex items-center gap-4">
                    <div class="flex-1">
                        <input
                            id="grade-input"
                            bind:value={grade}
                            oninput={() => {
                                if (errors.grade) errors.grade = "";
                            }}
                            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 {errors.grade
                                ? 'border-red-500'
                                : 'border-gray-300'}"
                            placeholder="0"
                        />
                        {#if errors.grade}
                            <p class="text-red-500 text-sm mt-1">
                                {errors.grade}
                            </p>
                        {/if}
                    </div>
                </div>
            </div>

            <div class="mb-6">
                <label class="block text-sm mb-2" for="feedback-input">
                    Nhận xét
                </label>
                <textarea
                    id="feedback-input"
                    bind:value={feedback}
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows={4}
                    placeholder="Nhập nhận xét cho báo cáo..."
                ></textarea>
            </div>

            <div class="mb-6">
                <label class="block text-sm mb-2 text-gray-600">
                    Mẫu nhận xét nhanh (click để thêm vào):
                </label>
                <div class="flex flex-wrap gap-2">
                    {#each quickFeedbacks as text}
                        <button
                            type="button"
                            onclick={() => handleQuickFeedback(text)}
                            class="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-full transition-colors"
                        >
                            {text}
                        </button>
                    {/each}
                </div>
            </div>

            <div class="flex justify-end gap-3 pt-6 border-t border-gray-200">
                <button
                    type="button"
                    onclick={onClose}
                    class="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                >
                    Hủy
                </button>
                <button
                    type="submit"
                    class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                >
                    Lưu điểm
                </button>
            </div>
        </form>
    </div>
</div>
