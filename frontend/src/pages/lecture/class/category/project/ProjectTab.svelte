<script lang="ts">
    import {
        Plus,
        Search,
        Users,
        Tag,
        Edit,
        Trash2,
        Eye,
        Copy,
        CheckSquare,
        Square,
    } from "lucide-svelte";
    import CreateProjectModal from "./CreateProjectModal.svelte";
    import { projectStore } from "../../../../../stores/project-store";
    import { push } from "svelte-spa-router";
    import type {
        CreateProjectRequest,
        UpdateProjectRequest,
    } from "../../../../../dtos/project-dto";

    // Props
    let { classroomId, projectRoundId } = $props<{
        classroomId: string;
        projectRoundId: string;
    }>();

    // State
    const {
        projects,
        searchTerm,
        fetchProjects,
        addProject,
        addProjects,
        updateProjectData,
        removeProject,
    } = projectStore;
    //log projects
    let showCreateModal = $state(false);
    let editingProject = $state<UpdateProjectRequest | null>(null);
    let selectedProjectIds = $state<Set<string>>(new Set<string>());

    $effect(() => {
        fetchProjects(classroomId, projectRoundId);
    });

    const handleCreateProject = async (newProject: CreateProjectRequest) => {
        try {
            if (newProject.amount > 1) {
                await addProjects({
                    projects: Array(newProject.amount).fill({
                        ...newProject,
                        amount: undefined, // Remove the amount property for individual projects
                    }),
                    classroomId,
                    projectRoundId,
                });
                console.log(
                    `Created ${newProject.amount} projects successfully.`,
                );
            } else {
                const createdProject = await addProject({
                    ...newProject,
                    classroomId,
                    projectRoundId,
                });
                console.log("Created project:", createdProject);
            }
            showCreateModal = false; // Close modal after successful creation
        } catch (error) {
            console.error("Error while adding project(s):", error);
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
        selectedProjectIds = new Set(selectedProjectIds); // Trigger reactivity
    };

    const handleSelectAll = () => {
        if (selectedProjectIds.size === $projects.length) {
            selectedProjectIds = new Set();
        } else {
            selectedProjectIds = new Set($projects.map((p) => p.id));
        }
    };

    const handleBulkDelete = async () => {
        if (selectedProjectIds.size === 0)
            return alert("Vui lòng chọn ít nhất 1 đề tài");
        if (confirm(`Xóa ${selectedProjectIds.size} đề tài đã chọn?`)) {
            for (const id of selectedProjectIds) {
                await removeProject(classroomId, id);
            }
            selectedProjectIds = new Set();
        }
    };
</script>

<div class="rounded-lg bg-white p-6 shadow-sm">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-4">
            <h2 class="text-xl font-semibold">Danh sách đề tài</h2>
            {#if $projects.length > 0}
                <div class="flex gap-2">
                    <button
                        onclick={handleSelectAll}
                        class="flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-1.5 text-sm transition-colors hover:bg-gray-50"
                    >
                        {#if selectedProjectIds.size === $projects.length && $projects.length > 0}
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
                            class="flex items-center gap-2 rounded-lg bg-red-600 px-3 py-1.5 text-sm text-white transition-colors hover:bg-red-700"
                        >
                            <Trash2 class="h-4 w-4" /> Xóa
                        </button>
                    {/if}
                </div>
            {/if}
        </div>
        <button
            onclick={() => (showCreateModal = true)}
            class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
        >
            <Plus class="h-5 w-5" /> Thêm đề tài mới
        </button>
    </div>

    <div class="mb-6 space-y-4">
        <div class="relative">
            <Search
                class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
                type="text"
                bind:value={$searchTerm}
                placeholder="Tìm kiếm đề tài..."
                class="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
        </div>
    </div>

    {#if $projects.length === 0}
        <div class="py-16 text-center">
            <div
                class="mx-auto flex h-16 w-16 items-center justify-center rounded-full border-2 border-dashed border-gray-300 text-gray-400"
            >
                <Search class="h-8 w-8" />
            </div>
            <p class="mt-4 text-gray-600">Không tìm thấy đề tài nào phù hợp.</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 gap-4">
            {#each $projects as project (project.id)}
                <div
                    class="group relative rounded-lg border p-5 transition-all hover:shadow-md"
                >
                    <div class="flex items-start justify-between">
                        <div class="flex items-start gap-3">
                            <button
                                onclick={() => toggleSelectProject(project.id)}
                                class="mt-1"
                            >
                                {#if selectedProjectIds.has(project.id)}
                                    <CheckSquare
                                        class="h-5 w-5 text-blue-600"
                                    />
                                {:else}
                                    <Square class="h-5 w-5 text-gray-300" />
                                {/if}
                            </button>
                            <div>
                                <div class="mb-2 flex items-center gap-3">
                                    <h3 class="text-lg font-medium">
                                        {project.title}
                                        <!-- Display title -->
                                    </h3>
                                </div>
                                <p class="mb-3 text-sm text-gray-600">
                                    {project.description}
                                </p>
                                <div
                                    class="flex items-center gap-4 text-xs text-gray-500"
                                >
                                    <div class="flex items-center gap-1">
                                        <Users class="h-4 w-4" />
                                        <span
                                            >{project.minMember}-{project.maxMember}
                                            thành viên</span
                                        >
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="flex flex-col items-end gap-4">
                            <div class="flex gap-1">
                                <button
                                    onclick={() =>
                                        (editingProject = {
                                            ...project,
                                            projectId: project.id,
                                        })}
                                    class="rounded p-2 text-gray-600 hover:bg-gray-100"
                                    title="Sửa"><Edit class="h-5 w-5" /></button
                                >
                                <button
                                    onclick={() =>
                                        handleDeleteProject(project.id)}
                                    class="rounded p-2 text-red-600 hover:bg-red-50"
                                    title="Xóa"
                                    ><Trash2 class="h-5 w-5" /></button
                                >
                                <button
                                    onclick={() =>
                                        push(
                                            `/lecture/my-classes/${classroomId}/categories/${projectRoundId}/projects/${project.id}`,
                                        )}
                                    class="rounded p-2 text-blue-600 hover:bg-blue-50"
                                    title="Xem chi tiết"
                                    ><Eye class="h-5 w-5" /></button
                                >
                            </div>
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
        onClose={() => (showCreateModal = false)}
        onSubmit={handleCreateProject}
    />
{/if}

{#if editingProject}
    <CreateProjectModal
        {projectRoundId}
        {editingProject}
        onClose={() => (editingProject = null)}
        onSubmit={handleUpdateProject}
    />
{/if}
