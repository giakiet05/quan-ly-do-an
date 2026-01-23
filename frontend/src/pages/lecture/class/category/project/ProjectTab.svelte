<script lang="ts">
    import {
        Plus,
        Search,
        Users,
        Edit,
        Trash2,
        Eye,
        CheckSquare,
        Square,
        Download,
        Upload,
        FileSpreadsheet,
    } from "lucide-svelte";
    import CreateProjectModal from "./CreateProjectModal.svelte";
    import { projectStore } from "../../../../../stores/project-store";
    import {
        downloadProjectTemplate,
        uploadProjectsExcel,
    } from "../../../../../services/project-service";
    import { push } from "svelte-spa-router";
    import type {
        CreateProjectRequest,
        UpdateProjectRequest,
    } from "../../../../../dtos/project-dto";
    import type { ProjectRound } from "../../../../../types/project-round";

    // Props
    let { classroomId, projectRoundId, projectRound } = $props<{
        classroomId: string;
        projectRoundId: string;
        projectRound: ProjectRound | null;
    }>();

    // Store
    const {
        projects,
        searchTerm,
        fetchProjects,
        addProject,
        addProjects,
        updateProjectData,
        removeProject,
    } = projectStore;

    // State
    let showCreateModal = $state(false);
    let editingProject = $state<UpdateProjectRequest | null>(null);
    let selectedProjectIds = $state<Set<string>>(new Set());
    let fileInput: HTMLInputElement;
    $effect(() => {
        fetchProjects(classroomId, projectRoundId);
    });

    const handleCreateProject = async (newProject: CreateProjectRequest) => {
        try {
            if (newProject.amount > 1) {
                await addProjects({
                    projects: Array(newProject.amount).fill({
                        ...newProject,
                        amount: undefined,
                    }),
                    classroomId,
                    projectRoundId,
                });
            } else {
                await addProject({
                    ...newProject,
                    classroomId,
                    projectRoundId,
                });
            }
            showCreateModal = false;
        } catch (error) {
            console.error(error);
            alert("Không thể thêm đề tài. Vui lòng thử lại.");
        }
    };

    const handleUpdateProject = async (
        updatedProject: UpdateProjectRequest,
    ) => {
        if (!editingProject) return;
        await updateProjectData({
            ...updatedProject,
            classroomId,
            projectId: editingProject.projectId,
        });
        editingProject = null;
    };

    const handleDeleteProject = async (id: string) => {
        if (confirm("Bạn có chắc chắn muốn xóa đề tài này?")) {
            await removeProject(classroomId, id);
        }
    };

    const toggleSelectProject = (id: string) => {
        if (selectedProjectIds.has(id)) {
            selectedProjectIds.delete(id);
        } else {
            selectedProjectIds.add(id);
        }
        selectedProjectIds = new Set(selectedProjectIds);
    };

    const handleSelectAll = () => {
        if (selectedProjectIds.size === $projects.length) {
            selectedProjectIds = new Set();
        } else {
            selectedProjectIds = new Set($projects.map((p) => p.id));
        }
    };

    const handleBulkDelete = async () => {
        if (selectedProjectIds.size === 0) {
            alert("Vui lòng chọn ít nhất 1 đề tài");
            return;
        }
        if (confirm(`Xóa ${selectedProjectIds.size} đề tài đã chọn?`)) {
            for (const id of selectedProjectIds) {
                await removeProject(classroomId, id);
            }
            selectedProjectIds = new Set();
        }
    };
    const handleDownloadTemplate = async () => {
        try {
            await downloadProjectTemplate();
        } catch (error) {
            console.error(error);
            alert("Lỗi khi tải xuống file mẫu.");
        }
    };

    // ✨ NEW: Hàm kích hoạt input file khi bấm nút Upload
    const triggerFileUpload = () => {
        fileInput.click();
    };

    // ✨ NEW: Hàm xử lý khi người dùng chọn file Excel xong
    const handleFileChange = async (event: Event) => {
        const target = event.target as HTMLInputElement;
        const file = target.files?.[0];

        if (!file) return;

        // Reset value để nếu chọn lại file cũ vẫn kích hoạt event
        target.value = "";

        try {
            if (confirm(`Bạn muốn import đề tài từ file "${file.name}"?`)) {
                await uploadProjectsExcel(classroomId, projectRoundId, file);
                alert("Import thành công!");
                // Refresh lại danh sách
                fetchProjects(classroomId, projectRoundId);
            }
        } catch (error) {
            console.error(error);
            alert("Lỗi khi upload file. Vui lòng kiểm tra lại định dạng.");
        }
    };
</script>

<div class="rounded-xl bg-white p-6 shadow-sm ring-1 ring-gray-100">
    <!-- HEADER -->
    <div class="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div>
            <h2 class="text-2xl font-semibold text-gray-800">
                Danh sách đề tài
            </h2>
            <p class="text-sm text-gray-500">
                Quản lý và phân công đề tài cho lớp học
            </p>
        </div>

        <div class="flex items-center gap-2">
            {#if $projects.length > 0}
                <button
                    onclick={handleSelectAll}
                    class="flex items-center gap-2 rounded-lg border px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                >
                    {#if selectedProjectIds.size === $projects.length}
                        <CheckSquare class="h-4 w-4 text-blue-600" />
                    {:else}
                        <Square class="h-4 w-4" />
                    {/if}
                    {selectedProjectIds.size > 0
                        ? `Đã chọn ${selectedProjectIds.size}`
                        : "Chọn tất cả"}
                </button>

                {#if selectedProjectIds.size > 0}
                    <button
                        onclick={handleBulkDelete}
                        class="flex items-center gap-2 rounded-lg bg-red-600 px-3 py-2 text-sm text-white hover:bg-red-700"
                    >
                        <Trash2 class="h-4 w-4" />
                        Xóa
                    </button>
                {/if}
            {/if}
        </div>
        <button
            onclick={() => (showCreateModal = true)}
            class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
        >
            <Plus class="h-5 w-5" /> Thêm đề tài mới
        </button>
    </div>

    <!-- SEARCH -->
    <div class="relative mb-6">
        <Search
            class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
        />
        <input
            type="text"
            bind:value={$searchTerm}
            placeholder="Tìm theo tên hoặc mô tả đề tài..."
            class="w-full rounded-lg border border-gray-300 bg-gray-50 py-2 pl-10 pr-4 text-sm focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
    </div>

    <!-- CONTENT -->
    {#if $projects.length === 0}
        <div class="py-20 text-center">
            <div
                class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100"
            >
                <Search class="h-8 w-8 text-gray-400" />
            </div>
            <p class="text-gray-600">Chưa có đề tài nào</p>
            <p class="text-sm text-gray-400">
                Hãy thêm đề tài đầu tiên cho vòng này
            </p>
        </div>
    {:else}
        <div class="space-y-4">
            {#each $projects as project (project.id)}
                <div
                    class="group relative rounded-xl border border-gray-200 bg-white p-5 transition hover:shadow-md"
                >
                    <div class="flex justify-between gap-4">
                        <div class="flex gap-4">
                            <button
                                onclick={() => toggleSelectProject(project.id)}
                                class="pt-1"
                            >
                                {#if selectedProjectIds.has(project.id)}
                                    <CheckSquare
                                        class="h-5 w-5 text-blue-600"
                                    />
                                {:else}
                                    <Square
                                        class="h-5 w-5 text-gray-300 group-hover:text-gray-400"
                                    />
                                {/if}
                            </button>

                            <div>
                                <h3 class="text-lg font-semibold text-gray-800">
                                    {project.title}
                                </h3>
                                <p
                                    class="mt-1 line-clamp-2 text-sm text-gray-600"
                                >
                                    {project.description}
                                </p>

                                <div
                                    class="mt-3 flex items-center gap-4 text-xs text-gray-500"
                                >
                                    <div class="flex items-center gap-1">
                                        <Users class="h-4 w-4" />
                                        {project.minMember}-{project.maxMember}
                                        thành viên
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div
                            class="flex items-start gap-1 opacity-0 transition group-hover:opacity-100"
                        >
                            <button
                                onclick={() =>
                                    (editingProject = {
                                        ...project,
                                        projectId: project.id,
                                    })}
                                class="rounded-lg p-2 text-gray-600 hover:bg-gray-100"
                                title="Sửa"
                            >
                                <Edit class="h-5 w-5" />
                            </button>

                            <button
                                onclick={() => handleDeleteProject(project.id)}
                                class="rounded-lg p-2 text-red-600 hover:bg-red-50"
                                title="Xóa"
                            >
                                <Trash2 class="h-5 w-5" />
                            </button>

                            <button
                                onclick={() =>
                                    push(
                                        `/lecture/my-classes/${classroomId}/categories/${projectRoundId}/projects/${project.id}`,
                                    )}
                                class="rounded-lg p-2 text-blue-600 hover:bg-blue-50"
                                title="Chi tiết"
                            >
                                <Eye class="h-5 w-5" />
                            </button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

{#if showCreateModal}
    <CreateProjectModal
        {projectRoundId}
        {projectRound}
        onClose={() => (showCreateModal = false)}
        onSubmit={handleCreateProject}
    />
{/if}

{#if editingProject}
    <CreateProjectModal
        {projectRound}
        {projectRoundId}
        {editingProject}
        onClose={() => (editingProject = null)}
        onSubmit={handleUpdateProject}
    />
{/if}
