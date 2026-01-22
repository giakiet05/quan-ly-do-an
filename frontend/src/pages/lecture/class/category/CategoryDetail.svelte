<script lang="ts">
    import { ArrowLeft, Download } from "lucide-svelte";

    import {
        mockCategoriesList,
        type ProjectCategory,
    } from "../../../../types/category";
    import { mockClasses } from "../../../../mocks/classes.mock";
    import ProjectTab from "./project/ProjectTab.svelte";
    import ReportsTab from "./report/ReportsTab.svelte";
    import { projectRoundStore } from "../../../../stores/project-round-store";
    import type { ProjectRound } from "../../../../types/project-round";
    import { getClassroom } from "../../../../services/classroom-service";
    import type { ClassroomResponse } from "../../../../dtos";

    // Lấy params từ URL
    let { params } = $props();
    let classroomId: string = $derived(params.id);
    let projectRoundId: string = $derived(params.categoryId);

    // 2. Tạo một state để chứa dữ liệu sau khi fetch xong
    let projectRound = $state<ProjectRound | null>(null);
    let isLoading = $state(true);

    // 3. Sử dụng $effect để tự động fetch lại khi classroomId hoặc projectRoundId thay đổi
    $effect(() => {
        if (classroomId && projectRoundId) {
            isLoading = true;
            projectRoundStore
                .getRoundById(classroomId, projectRoundId)
                .then((data) => {
                    projectRound = data;
                })
                .catch((err) => {
                    console.error("Lỗi khi lấy dữ liệu:", err);
                })
                .finally(() => {
                    isLoading = false;
                });
        }
    });
    // state cho class
    let classDetail = $state<ClassroomResponse | null>(null);
    $effect(() => {
        if (classroomId) {
            getClassroom(classroomId)
                .then((data) => {
                    classDetail = data;
                })
                .catch((err) => {
                    console.error("Lỗi khi lấy chi tiết lớp học:", err);
                });
        }
    });

    // 2. State quản lý tab hiện tại
    let activeTab = $state("projects");

    const tabs = [
        { id: "projects", label: "Đề tài" },
        { id: "reports", label: "Báo cáo" },
        { id: "settings", label: "Cài đặt" },
    ];

    // Helper format ngày
    const formatDate = (date: string) =>
        new Date(date).toLocaleDateString("vi-VN");
</script>

<div class="mb-6 rounded-lg bg-white p-6 shadow-sm">
    <button
        onclick={() => history.back()}
        class="mb-4 flex items-center gap-2 text-gray-600 transition-colors hover:text-gray-900"
    >
        <ArrowLeft class="h-5 w-5" />
        Quay lại tổng quan lớp học
    </button>

    <div class="mb-6 flex items-start justify-between">
        <div>
            <h1 class="mb-2 text-2xl font-semibold">{projectRound?.name}</h1>
            <p class="mb-2 text-gray-600">{projectRound?.description}</p>
            <p class="text-sm text-gray-500">
                Lớp: {classDetail?.name} - {classDetail?.semester} | Thời gian: {formatDate(
                    projectRound?.startDate.toDateString() || "",
                )} - {formatDate(projectRound?.endDate.toDateString() || "")}
            </p>
        </div>
    </div>

    <div class="flex gap-6 border-b border-gray-200">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="relative px-1 pb-3 transition-colors {activeTab ===
                tab.id
                    ? 'text-blue-600'
                    : 'text-gray-600 hover:text-gray-900'}"
            >
                {tab.label}
                {#if activeTab === tab.id}
                    <div
                        class="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600"
                    ></div>
                {/if}
            </button>
        {/each}
    </div>
</div>

{#if activeTab === "projects"}
    <ProjectTab {classroomId} {projectRoundId} />
{:else if activeTab === "reports"}
    <ReportsTab {classroomId} {projectRoundId} />
{:else if activeTab === "settings"}
    <div class="rounded-lg bg-white p-6 shadow-sm">
        <h2 class="mb-6 text-xl font-semibold">Cài đặt hạng mục</h2>

        <div class="space-y-6">
            <div class="border-b border-gray-200 pb-6">
                <h3 class="mb-4 font-medium">Thông tin chung</h3>
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label class="mb-2 block text-sm text-gray-600"
                            >Tên hạng mục</label
                        >
                        <input
                            type="text"
                            value={projectRound?.name}
                            class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
            </div>

            <div class="border-b border-gray-200 pb-6">
                <h3 class="mb-4 font-medium">Thời gian đăng ký</h3>
                <div class="space-y-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <div class="font-medium">Mở đăng ký đề tài</div>
                            <div class="text-sm text-gray-600">
                                Cho phép sinh viên đăng ký đề tài trong hạng mục
                                này
                            </div>
                        </div>
                        <label
                            class="relative inline-flex cursor-pointer items-center"
                        >
                            <input
                                type="checkbox"
                                checked
                                class="peer sr-only"
                            />
                            <div
                                class="peer h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300"
                            ></div>
                        </label>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="mb-2 block text-sm text-gray-600"
                                >Ngày bắt đầu</label
                            >
                            <input
                                type="date"
                                value={projectRound?.startDate}
                                class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label class="mb-2 block text-sm text-gray-600"
                                >Ngày kết thúc</label
                            >
                            <input
                                type="date"
                                value={projectRound?.endDate}
                                class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>
                </div>
            </div>

            <div class="border-b border-gray-200 pb-6">
                <h3 class="mb-4 font-medium">Quy định nhóm</h3>
                <div class="space-y-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <div class="font-medium">
                                Cho phép sinh viên chỉnh sửa nhóm
                            </div>
                            <div class="text-sm text-gray-600">
                                Sinh viên có thể thêm/xóa thành viên trong nhóm
                                của mình
                            </div>
                        </div>
                        <label
                            class="relative inline-flex cursor-pointer items-center"
                        >
                            <input
                                type="checkbox"
                                checked
                                class="peer sr-only"
                            />
                            <div
                                class="peer h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300"
                            ></div>
                        </label>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="mb-2 block text-sm text-gray-600"
                                >Số thành viên tối thiểu</label
                            >
                            <input
                                type="number"
                                value={1}
                                min={1}
                                max={10}
                                class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label class="mb-2 block text-sm text-gray-600"
                                >Số thành viên tối đa</label
                            >
                            <input
                                type="number"
                                value={3}
                                min={1}
                                max={10}
                                class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>

                    <div
                        class="rounded-lg border border-blue-200 bg-blue-50 p-4"
                    >
                        <div class="flex items-start justify-between gap-4">
                            <div class="flex-1">
                                <div class="mb-1 font-medium text-blue-900">
                                    Áp dụng giới hạn cho tất cả đề tài
                                </div>
                                <div class="text-sm text-blue-700">
                                    Khi bật, tất cả đề tài trong hạng mục sẽ sử
                                    dụng giới hạn thành viên mặc định này. Các
                                    thiết lập riêng của từng đề tài sẽ bị ghi
                                    đè.
                                </div>
                            </div>
                            <label
                                class="relative inline-flex flex-shrink-0 cursor-pointer items-center"
                            >
                                <input
                                    type="checkbox"
                                    class="peer sr-only"
                                    onchange={(e) => {
                                        if (e.currentTarget.checked) {
                                            if (
                                                !confirm(
                                                    "Bạn có chắc chắn muốn áp dụng giới hạn này cho TẤT CẢ đề tài? Các thiết lập riêng của từng đề tài sẽ bị ghi đè.",
                                                )
                                            ) {
                                                e.currentTarget.checked = false;
                                            }
                                        }
                                    }}
                                />
                                <div
                                    class="peer h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300"
                                ></div>
                            </label>
                        </div>
                    </div>
                </div>
            </div>

            <div class="flex justify-end gap-3">
                <button
                    class="rounded-lg border border-gray-300 px-6 py-2 transition-colors hover:bg-gray-50"
                >
                    Hủy
                </button>
                <button
                    class="rounded-lg bg-blue-600 px-6 py-2 text-white transition-colors hover:bg-blue-700"
                >
                    Lưu thay đổi
                </button>
                <button
                    class="rounded-lg bg-red-600 px-6 py-2 text-white transition-colors hover:bg-red-700"
                >
                    Xóa hạng mục
                </button>
            </div>
        </div>
    </div>
{/if}
