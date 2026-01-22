<script lang="ts">
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import { params } from "svelte-spa-router";
    import { ArrowLeft } from "../../../libs/Icons";

    // Components
    import CategoryList from "./category/CategoryList.svelte";
    import CategoryDetailModal from "./category/CategoryDetailModal.svelte";
    import StudentTab from "./student/StudentTab.svelte";
    import AnnouncementTab from "./announcement/AnnouncementTab.svelte";
    import CreateAnnouncementModal from "./announcement/CreateAnnouncementModal.svelte";
    import ChatTab from "./ChatTab.svelte";

    // Stores & Type
    import { projectRoundStore } from "../../../stores/project-round-store";
    import type { Announcement } from "../../../types/announcement";
    import type { ProjectRound } from "../../../types/project-round";
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
    const students = $derived(classData?.students ?? []);

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
            console.log("Creating round with data:", data);
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
                projectRoundId: editingCategory.id,
                classroomId: classroomId,
                name: data.name,
                description: data.description,
                startDate: data.startDate,
                endDate: data.endDate,
                defaultMinMember: data.minStudents,
                defaultMaxMember: data.maxStudents,
            };
            await projectRoundStore.editRound(payload);
            editingCategory = null;
        } catch (error) {
            console.error("❌ Error updating category:", error);
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
            <div
                class="flex flex-col md:flex-row md:items-end justify-between gap-4"
            >
                <div>
                    <h1
                        class="text-3xl font-bold text-gray-800 flex items-center gap-3"
                    >
                        {classData.name}
                        <span
                            class="text-sm font-normal px-2 py-1 bg-blue-100 text-blue-700 rounded"
                        >
                            {classData.status === "active"
                                ? "Đang chạy"
                                : "Lưu trữ"}
                        </span>
                    </h1>

                    <div
                        class="flex flex-wrap items-center gap-y-2 gap-x-4 mt-2 text-gray-600"
                    >
                        <div class="flex items-center gap-1">
                            <span class="font-medium text-gray-900"></span>
                            {classData.semester} — Năm {classData.year}
                        </div>

                        <div class="hidden md:block w-px h-4 bg-gray-300"></div>

                        <div class="flex items-center gap-2">
                            <span class="font-regular text-gray-900"
                                >Mã mời:</span
                            >
                            <code
                                class="bg-gray-100 px-2 py-0.5 rounded font-mono text-blue-600 font-bold border border-gray-200"
                            >
                                {classData.invitationCode}
                            </code>
                            <button
                                onclick={() => {
                                    navigator.clipboard.writeText(
                                        classData.invitationCode,
                                    );
                                    // Mày có thể thêm toast thông báo ở đây nếu muốn
                                }}
                                class="text-xs text-blue-500 hover:underline"
                            >
                                Sao chép
                            </button>
                        </div>

                        <div class="hidden md:block w-px h-4 bg-gray-300"></div>

                        <div class="flex items-center gap-1">
                            <span class="font-regular text-gray-900"
                                >Sĩ số:</span
                            >
                            {classData.studentCount}/{classData.maxStudents} sinh
                            viên
                        </div>
                    </div>
                </div>

                {#if classData.avatar}
                    <img
                        src={classData.avatar}
                        alt="Class Avatar"
                        class="w-16 h-16 rounded-xl object-cover border-2 border-white shadow-sm"
                    />
                {/if}
            </div>
        {:else}
            <div class="animate-pulse">
                <div class="h-8 bg-gray-200 rounded w-1/3 mb-2"></div>
                <div class="h-4 bg-gray-200 rounded w-1/2"></div>
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
                        id={classroomId}
                    />
                </div>
            {:else if activeTab === "chat"}
                <div class="w-full h-full min-h-[650px]">
                    <ChatTab {classroomId} />
                </div>
            {/if}
        </div>
    {/key}
</div>

{#if classData}
    {#if isCreateCategoryModalOpen}
        <CategoryDetailModal
            {classData}
            onClose={() => (isCreateCategoryModalOpen = false)}
            onSubmit={handleCreateRound}
        />
    {/if}

    {#if editingCategory}
        <CategoryDetailModal
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
