<script lang="ts">
    import { params } from "svelte-spa-router";
    import type { ClassData } from "../../../types/class";
    import { ArrowLeft } from "../../../libs/Icons";
    import { fade } from "svelte/transition";
    import CategoryList from "./category/CategoryList.svelte";
    import CreateCategoryModal from "./category/CreateCategoryModal.svelte";
    import { mockCategoriesList } from "../../../types/category";
    import StudentTab from "./student/StudentTab.svelte";
    import { mockStudentInClassData } from "../../../types/student";
    import AnnouncementTab from "./announcement/AnnouncementTab.svelte";
    import CreateAnnouncementModal from "./announcement/CreateAnnouncementModal.svelte";

    let categories = $derived(
        mockCategoriesList.map((cat) => ({
            ...cat,
            status:
                new Date(cat.endDate) < new Date()
                    ? "đã kết thúc"
                    : new Date(cat.startDate) > new Date()
                      ? "săp diễn ra"
                      : "đang diễn ra",
        })),
    );
    let students = $state(mockStudentInClassData);

    let classData = $state<ClassData>({
        id: "1",
        name: "Phát triển ứng dụng Web",
        school: "Trường Đại học Bách Khoa",
        description:
            "Lớp học tập trung vào phát triển web với React, Node.js và REST API.",
        studentCount: 45,
        semester: "HK2 2023–2024",
        status: "active",
    });
    let activeTab = $state<"overview" | "students" | "announcements" | "chat">(
        "overview",
    );
    let isCreateCategoryModalOpen = $state(false);
    let showCreateAnnouncementModal = $state(false);

    const tabs = [
        { id: "overview", label: "Tổng quan" },
        { id: "students", label: "Sinh viên" },
        { id: "announcements", label: "Thông báo" },
        { id: "chat", label: "Chat" },
    ] as const;
</script>

<div class="bg-white rounded-lg shadow-sm">
    <div class="p-2 border-b border-gray-200">
        <button
            onclick={() => history.back()}
            class="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
        >
            <ArrowLeft class="w-5 h-5" />
            Quay lại danh sách lớp học
        </button>
        <h1 class="text-2xl mb-2">
            Chi tiết Lớp: {classData.name} - {classData.semester}
        </h1>
    </div>
    <!-- Tabs -->
    <div class="flex gap-6 border-b border-gray-200">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="pb-3 px-1 transition-colors relative
            {activeTab === tab.id
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
    <!-- Tab Content -->
</div>
<div class="bg-white rounded-lg mt-4 shadow-sm">
    <div class="tab-content" transition:fade>
        {#if activeTab === "overview"}
            <div class="tab-content-enter-active">
                <!-- Content for "Tổng quan" -->
                <CategoryList
                    onOpen={() => (isCreateCategoryModalOpen = true)}
                    onCLose={() => (isCreateCategoryModalOpen = false)}
                    {classData}
                />
            </div>
        {:else if activeTab === "students"}
            <div class="tab-content-enter-active">
                <!-- Content for "Sinh viên" -->
                <StudentTab {students} {categories} />
            </div>
        {:else if activeTab === "announcements"}
            <div class="tab-content-enter-active">
                <!-- Content for "Thông báo" -->
                <AnnouncementTab
                    onOpen={() => (showCreateAnnouncementModal = true)}
                    onCLose={() => (showCreateAnnouncementModal = false)}
                />
            </div>
        {:else if activeTab === "chat"}
            <div class="tab-content-enter-active">
                <!-- Content for "Chat" -->
            </div>
        {/if}
    </div>
</div>
{#if isCreateCategoryModalOpen}
    <CreateCategoryModal
        onClose={() => (isCreateCategoryModalOpen = false)}
        {classData}
        onSubmit={() => {}}
    />
{/if}
{#if showCreateAnnouncementModal}
    <CreateAnnouncementModal
        onClose={() => (showCreateAnnouncementModal = false)}
        onSubmit={() => {}}
    />
{/if}

<style>
    .tab-content {
        transition:
            opacity 0.3s ease,
            transform 0.3s ease;
    }
    .tab-content-enter {
        opacity: 0;
        transform: translateY(10px);
    }
    .tab-content-enter-active {
        opacity: 1;
        transform: translateY(0);
    }
    .tab-content-leave {
        opacity: 1;
        transform: translateY(0);
    }
    .tab-content-leave-active {
        opacity: 0;
        transform: translateY(-10px);
    }
</style>
