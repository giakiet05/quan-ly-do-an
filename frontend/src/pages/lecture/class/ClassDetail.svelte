<script lang="ts">
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import { params } from "svelte-spa-router";
    import { ArrowLeft } from "../../../libs/Icons";

    // Components
    import CategoryList from "./category/CategoryList.svelte";
    import CreateCategoryModal from "./category/CreateCategoryModal.svelte";
    import StudentTab from "./student/StudentTab.svelte";
    import AnnouncementTab from "./announcement/AnnouncementTab.svelte";
    import CreateAnnouncementModal from "./announcement/CreateAnnouncementModal.svelte";

    // Stores & Types
    import { projectRoundStore } from "../../../stores/project-round-store";
    import type { Announcement } from "../../../types/announcement";
    import type { ProjectRound } from "../../../types/project-round";
    import { mockStudentInClassData } from "../../../types/student";
    import { classStore } from "../../../stores/class-store";
    import type { ClassItem } from "../../../types/class";
    import type { UpdateProjectRoundRequest } from "../../../dtos/project-dto";

    const classroomId = $derived($params?.id);
    const classesStore = classStore.classes;
    const classData = $derived(
        $classesStore.find((c: ClassItem) => c.id === classroomId),
    );

    $effect(() => {
        if (classroomId) {
            projectRoundStore.fetchRounds(classroomId);
            if ($classesStore.length === 0) {
                classStore.fetchMyClasses();
            }
        }
    });
    let students = $state(mockStudentInClassData);

    // UI States
    let activeTab = $state<"overview" | "students" | "announcements" | "chat">(
        "overview",
    );
    let isCreateCategoryModalOpen = $state(false);
    let editingCategory = $state<ProjectRound | null>(null);
    let editingAnnouncement = $state<Announcement | null>(null);
    let showCreateAnnouncementModal = $state(false);

    const tabs = [
        { id: "overview", label: "Tổng quan" },
        { id: "students", label: "Sinh viên" },
        { id: "announcements", label: "Thông báo" },
        { id: "chat", label: "Chat" },
    ] as const;

    // 4. Handlers
    async function handleCreateRound(data: any) {
        try {
            await projectRoundStore.addRound({
                ...data,
                classroomId: classroomId,
            });
            isCreateCategoryModalOpen = false;
        } catch (error) {
            alert("Có lỗi khi tạo hạng mục!");
        }
    }

    async function handleUpdateRound(data: any) {
        if (!editingCategory) return;
        try {
            const payload: UpdateProjectRoundRequest = {
                roundId: editingCategory.id,
                projectRoundId: editingCategory.id,
                classroomId: classroomId,
                name: data.name,
                description: data.description,
                startDate: data.startDate,
                endDate: data.endDate,
            };

            await projectRoundStore.editRound(payload);
            editingCategory = null;
        } catch (error) {
            alert("Có lỗi khi cập nhật!");
        }
    }
</script>

<div class="bg-white rounded-lg shadow-sm p-6">
    <div class="p-2">
        <button
            onclick={() => history.back()}
            class="flex items-center text-gray-600 hover:text-gray-900 transition-colors mb-4"
        >
            <ArrowLeft class="w-5 h-5 mr-1" />
            Quay lại danh sách lớp học
        </button>

        {#if classData}
            <h1 class="text-2xl font-bold text-gray-800">
                Lớp: {classData.name}
            </h1>
            <p class="text-gray-500">
                Mã lớp: {classData.id} • {classData.semester} • {classData.studentCount}
                sinh viên
            </p>
        {:else}
            <div class="animate-pulse">
                <div class="h-8 bg-gray-200 rounded w-1/3 mb-2"></div>
                <div class="h-4 bg-gray-200 rounded w-1/4"></div>
            </div>
        {/if}
    </div>

    <div class="flex gap-6 border-b border-gray-200 mt-6">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="pb-3 px-1 transition-all relative font-medium
                {activeTab === tab.id
                    ? 'text-blue-600'
                    : 'text-gray-500 hover:text-gray-700'}"
            >
                {tab.label}
                {#if activeTab === tab.id}
                    <div
                        class="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600"
                        in:fade
                    ></div>
                {/if}
            </button>
        {/each}
    </div>
</div>

<div class="mt-4">
    {#key activeTab}
        <div in:fade={{ duration: 200 }}>
            {#if activeTab === "overview"}
                <div class="bg-white rounded-lg shadow-sm">
                    <CategoryList
                        {classData}
                        openEditCategoryModal={(category) =>
                            (editingCategory = category)}
                        onOpen={() => (isCreateCategoryModalOpen = true)}
                        onClose={() => (isCreateCategoryModalOpen = false)}
                    />
                </div>
            {:else if activeTab === "students"}
                <div class="bg-white rounded-lg shadow-sm">
                    <StudentTab {students} />
                </div>
            {:else if activeTab === "announcements"}
                <div class="bg-white rounded-lg shadow-sm">
                    <AnnouncementTab
                        onEdit={(announcement) =>
                            (editingAnnouncement = announcement)}
                        onOpen={() => (showCreateAnnouncementModal = true)}
                        onClose={() => (showCreateAnnouncementModal = false)}
                    />
                </div>
            {:else if activeTab === "chat"}
                <div
                    class="bg-white rounded-lg shadow-sm p-20 text-center text-gray-400"
                >
                    Tính năng Chat đang được phát triển...
                </div>
            {/if}
        </div>
    {/key}
</div>

{#if classData}
    {#if isCreateCategoryModalOpen}
        <CreateCategoryModal
            {classData}
            onClose={() => (isCreateCategoryModalOpen = false)}
            onSubmit={handleCreateRound}
        />
    {/if}

    {#if editingCategory}
        <CreateCategoryModal
            isEdit={true}
            initialData={editingCategory}
            {classData}
            onClose={() => (editingCategory = null)}
            onSubmit={handleUpdateRound}
        />
    {/if}
{/if}

{#if showCreateAnnouncementModal || editingAnnouncement}
    <CreateAnnouncementModal
        {editingAnnouncement}
        onClose={() => {
            showCreateAnnouncementModal = false;
            editingAnnouncement = null;
        }}
        onSubmit={(data) => {
            console.log("Announcement data:", data);
            showCreateAnnouncementModal = false;
            editingAnnouncement = null;
        }}
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
