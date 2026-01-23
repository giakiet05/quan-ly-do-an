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
    import ChatTab from "../../../components/ChatTab.svelte";

    // Stores & Type
    import { projectRoundStore } from "../../../stores/project-round-store";
    import type { Announcement } from "../../../types/announcement";
    import type { ProjectRound } from "../../../types/project-round";
    import { classStore } from "../../../stores/class-store";
    import type { ClassItem } from "../../../types/class";
    import type { UpdateProjectRoundRequest } from "../../../dtos/project-dto";
    import type { ClassroomResponse } from "../../../dtos";
    import { getClassroom } from "../../../services/classroom-service";
    import { postStore } from "../../../stores/post-store";
    import type {
        CreatePostRequest,
        PostResponse,
        UpdatePostRequest,
    } from "../../../dtos/post-dto";

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

    const students = $derived(classDetail?.students ?? []);
    const generalChannelId = $derived(classDetail?.generalChannelId ?? "");
    $effect(() => {
        console.log("Students list updated:", classDetail);
    });
    // UI States
    let activeTab = $state<"overview" | "students" | "announcements" | "chat">(
        "overview",
    );
    let isCreateCategoryModalOpen = $state(false);
    let editingCategory = $state<ProjectRound | null>(null);
    let editingAnnouncement = $state<PostResponse | null>(null);
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
            // Chụp ảnh data từ modal gửi về (đã bao gồm các thay đổi từ UI)
            // Chúng ta truyền cả 'data' (chứa thông tin mới)
            // và 'editingCategory.id' để Mapper biết đang update record nào
            const updateData = {
                ...data,
                id: editingCategory.id, // Đảm bảo luôn có ID chuẩn từ entity cũ
            };

            await projectRoundStore.editRound(updateData, classroomId);

            editingCategory = null;
            alert("✅ Cập nhật hạng mục thành công!");
        } catch (error) {
            alert("Có lỗi khi cập nhật: " + (error as Error).message);
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

                <div class="flex items-center justify-between gap-2">
                    <div class="flex items-center gap-1">
                        <span class="font-regular text-gray-900"
                            >Trạng thái:</span
                        >
                    </div>
                    <button
                        class="relative w-20 h-10 rounded-full transition-colors duration-300
                        focus:outline-none focus:ring-2 focus:ring-offset-2
                        {classData.status === 'active'
                            ? 'bg-green-500 hover:bg-green-600 focus:ring-green-400'
                            : 'bg-gray-300 hover:bg-gray-400 focus:ring-gray-400'}"
                        onclick={async () => {
                            const newStatus =
                                classData.status === "active"
                                    ? "inactive"
                                    : "active";
                            try {
                                await classStore.updateClassroomStatus(
                                    classroomId,
                                    newStatus,
                                );
                                classData.status = newStatus;
                            } catch (err) {
                                alert("Không thể cập nhật trạng thái lớp học!");
                                console.error(err);
                            }
                        }}
                    >
                        <!-- Nút tròn trượt -->
                        <span
                            class="absolute top-1 left-1 w-8 h-8 bg-white rounded-full shadow-md
                            transition-transform duration-300"
                            class:translate-x-10={classData.status === "active"}
                        ></span>

                        <!-- Text trạng thái -->
                        <span
                            class="absolute inset-0 flex items-center justify-center text-xs font-semibold
                            text-white select-none pointer-events-none"
                        >
                            {classData.status === "active" ? "ON" : "OFF"}
                        </span>
                    </button>
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
                    <StudentTab {classroomId} {students} />
                </div>
            {:else if activeTab === "announcements"}
                <div class="bg-white rounded-lg shadow-sm">
                    <AnnouncementTab
                        onEdit={(announcement) =>
                            (editingAnnouncement = announcement)}
                        onOpen={() => (showCreateAnnouncementModal = true)}
                        {classroomId}
                    />
                </div>
            {:else if activeTab === "chat"}
                <div class="w-full h-full min-h-[650px]">
                    <ChatTab {generalChannelId} />
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
        onSubmit={(post) => {
            if (editingAnnouncement) {
                postStore.updatePostData(
                    editingAnnouncement.id,
                    post as UpdatePostRequest,
                );
            } else {
                postStore.addPost(post as CreatePostRequest, classroomId);
            }
            showCreateAnnouncementModal = false;
            editingAnnouncement = null;
        }}
    />
{/if}
