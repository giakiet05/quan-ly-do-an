<script lang="ts">
    import { ArrowLeft, Download } from "lucide-svelte";

    import {
        mockCategoriesList,
        type ProjectCategory,
    } from "../../../../types/category";
    import { mockClasses } from "../../../../mocks/classes.mock";
    import ProjectTab from "./project/ProjectTab.svelte";
    import ReportsTab from "./report/ReportsTab.svelte";
    import { params } from "svelte-spa-router";
    import { projectRoundStore } from "../../../../stores/project-round-store";
    import { classStore } from "../../../../stores/class-store";

    const classId = $derived($params?.id);
    const categoryId = $derived($params?.categoryId);
    const classesStore = classStore.classes;

    const category = $derived(
        $projectRoundStore.find((r) => r.id === categoryId),
    );

    const classData = $derived($classesStore.find((c) => c.id === classId));

    $effect(() => {
        if (classId) {
            if ($classesStore.length === 0) {
                classStore.fetchMyClasses();
            }
            if ($projectRoundStore.length === 0) {
                projectRoundStore.fetchRounds(classId);
            }
        }
    });

    let activeTab = $state("projects");

    const tabs = [
        { id: "projects", label: "Đề tài" },
        { id: "reports", label: "Báo cáo" },
        { id: "settings", label: "Cài đặt" },
    ];

    // FIX LỖI 2: Chỉnh lại hàm formatDate để nhận cả Date object và string
    const formatDate = (date: any) => {
        if (!date) return "";
        const d = typeof date === "string" ? new Date(date) : date;
        return d.toLocaleDateString("vi-VN");
    };

    // Hàm bổ trợ để hiển thị lên input date (YYYY-MM-DD)
    const toInputDate = (date: any) => {
        if (!date) return "";
        const d = typeof date === "string" ? new Date(date) : date;
        return d.toISOString().split("T")[0];
    };
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
        {#if category && classData}
            <div>
                <h1 class="mb-2 text-2xl font-semibold">{category.name}</h1>
                <p class="mb-2 text-gray-600">{category.description}</p>
                <p class="text-sm text-gray-500">
                    Lớp: {classData.name} - {classData.semester} | Thời gian: {formatDate(
                        category.startDate,
                    )} - {formatDate(category.endDate)}
                </p>
            </div>
        {:else}
            <div class="flex h-64 items-center justify-center">
                <div class="text-gray-500">Đang tải dữ liệu thực tế...</div>
            </div>
        {/if}
        <button
            class="flex items-center gap-2 rounded-lg border border-gray-300 px-4 py-2 transition-colors hover:bg-gray-50"
        >
            <Download class="h-4 w-4" />
            Xuất dữ liệu
        </button>
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

{#if category}
    {#if activeTab === "projects"}
        <ProjectTab category={category as any} />
    {:else if activeTab === "reports"}
        <ReportsTab />
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
                                bind:value={category.name}
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
                                    Cho phép sinh viên đăng ký đề tài trong hạng
                                    mục này
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
                                    value={toInputDate(category.startDate)}
                                    onchange={(e) =>
                                        (category.startDate = new Date(
                                            e.currentTarget.value,
                                        ))}
                                    class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                            </div>
                            <div>
                                <label class="mb-2 block text-sm text-gray-600"
                                    >Ngày kết thúc</label
                                >
                                <input
                                    type="date"
                                    value={toInputDate(category.endDate)}
                                    onchange={(e) =>
                                        (category.endDate = new Date(
                                            e.currentTarget.value,
                                        ))}
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
                                    Sinh viên có thể thêm/xóa thành viên trong
                                    nhóm của mình
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
                                        Khi bật, tất cả đề tài trong hạng mục sẽ
                                        sử dụng giới hạn thành viên mặc định
                                        này. Các thiết lập riêng của từng đề tài
                                        sẽ bị ghi đè.
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
                        Xóa hạng mục
                    </button>
                </div>
            </div>
        </div>
    {/if}
{/if}
