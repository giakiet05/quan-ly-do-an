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
    import type { ProjectCategory } from "../../../../../types/category";
    import { mockProjects, type Project } from "../../../../../types/project";
    import { push } from "svelte-spa-router";

    let projects = mockProjects;

    // 1. Định nghĩa Props
    let { category, onSelectProject, onUpdateProjects } = $props<{
        category: ProjectCategory;
    }>();

    // 2. Local State
    let showCreateModal = $state(false);
    let editingProject = $state<Project | null>(null);
    let searchTerm = $state("");
    let statusFilter = $state<"all" | Project["status"]>("all");
    let sortBy = $state<"name" | "students" | "recent">("recent");
    let selectedProjectIds = $state(new Set<string>());

    // 3. Computed State ($derived)
    const stats = $derived({
        total: projects.length,
        available: projects.filter((p: Project) => p.status === "available")
            .length,
        full: projects.filter((p: Project) => p.status === "full").length,
        registered: projects.reduce(
            (sum: number, p: Project) => sum + p.currentStudents,
            0,
        ),
    });

    const filteredProjects = $derived(
        projects.filter((project: Project) => {
            const matchesSearch =
                project.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                project.description
                    .toLowerCase()
                    .includes(searchTerm.toLowerCase()) ||
                project.tags.some((tag) =>
                    tag.toLowerCase().includes(searchTerm.toLowerCase()),
                );

            const matchesStatus =
                statusFilter === "all" || project.status === statusFilter;

            return matchesSearch && matchesStatus;
        }),
    );

    const sortedProjects = $derived(
        [...filteredProjects].sort((a, b) => {
            if (sortBy === "name") return a.name.localeCompare(b.name);
            if (sortBy === "students")
                return (
                    b.maxStudents -
                    b.currentStudents -
                    (a.maxStudents - a.currentStudents)
                );
            return 0; // "recent" - giả định mảng đã sắp xếp theo thời gian
        }),
    );

    // 4. Helper Functions
    const getStatusConfig = (status: Project["status"]) => {
        const configs = {
            available: {
                label: "Còn chỗ",
                color: "text-green-600",
                bg: "bg-green-50",
                border: "border-green-200",
            },
            full: {
                label: "Đã đủ",
                color: "text-orange-600",
                bg: "bg-orange-50",
                border: "border-orange-200",
            },
            closed: {
                label: "Đã khóa",
                color: "text-red-600",
                bg: "bg-red-50",
                border: "border-red-200",
            },
            forming: {
                label: "Đang hình thành",
                color: "text-blue-600",
                bg: "bg-blue-50",
                border: "border-blue-200",
            },
        };
        return (
            configs[status] || {
                label: "Không xác định",
                color: "text-gray-600",
                bg: "bg-gray-50",
                border: "border-gray-200",
            }
        );
    };

    // 5. Event Handlers
    const handleCreateProject = (
        newProject: Omit<Project, "id" | "currentStudents">,
    ) => {
        const project: Project = {
            ...newProject,
            id: Date.now().toString(),
            currentStudents: 0,
        };
        onUpdateProjects([project, ...projects]);
        showCreateModal = false;
    };

    const handleUpdateProject = (
        updatedProject: Omit<Project, "id" | "currentStudents">,
    ) => {
        if (!editingProject) return;
        onUpdateProjects(
            projects.map((p: Project) =>
                p.id === editingProject!.id ? { ...p, ...updatedProject } : p,
            ),
        );
        editingProject = null;
    };

    const handleDeleteProject = (id: string) => {
        if (confirm("Bạn có chắc chắn muốn xóa đề tài này?")) {
            onUpdateProjects(projects.filter((p: Project) => p.id !== id));
        }
    };

    const handleDuplicateProject = (project: Project) => {
        const duplicated: Project = {
            ...project,
            id: Date.now().toString(),
            name: `${project.name} (Bản sao)`,
            currentStudents: 0,
            status: "available",
        };
        onUpdateProjects([duplicated, ...projects]);
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
        if (selectedProjectIds.size === sortedProjects.length) {
            selectedProjectIds = new Set();
        } else {
            selectedProjectIds = new Set(sortedProjects.map((p) => p.id));
        }
    };

    const handleBulkDelete = () => {
        if (selectedProjectIds.size === 0)
            return alert("Vui lòng chọn ít nhất 1 đề tài");
        if (confirm(`Xóa ${selectedProjectIds.size} đề tài đã chọn?`)) {
            onUpdateProjects(
                projects.filter((p: Project) => !selectedProjectIds.has(p.id)),
            );
            selectedProjectIds = new Set();
        }
    };

    const toggleAllowEdit = (projectId: string, currentValue: boolean) => {
        onUpdateProjects(
            projects.map((p: Project) =>
                p.id === projectId
                    ? { ...p, allowStudentEdit: !currentValue }
                    : p,
            ),
        );
    };
</script>

<div class="rounded-lg bg-white p-6 shadow-sm">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-4">
            <h2 class="text-xl font-semibold">Danh sách đề tài</h2>
            {#if sortedProjects.length > 0}
                <div class="flex gap-2">
                    <button
                        onclick={handleSelectAll}
                        class="flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-1.5 text-sm transition-colors hover:bg-gray-50"
                    >
                        {#if selectedProjectIds.size === sortedProjects.length && sortedProjects.length > 0}
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
                bind:value={searchTerm}
                placeholder="Tìm kiếm đề tài..."
                class="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
        </div>

        <div class="flex flex-wrap items-center gap-4">
            <div class="flex items-center gap-2">
                <span class="text-sm text-gray-600">Trạng thái:</span>
                {#each ["all", "available", "full", "closed"] as status}
                    <button
                        onclick={() => (statusFilter = status as any)}
                        class="rounded-lg px-3 py-1 text-sm transition-colors {statusFilter ===
                        status
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
                    >
                        {status === "all"
                            ? `Tất cả (${projects.length})`
                            : status === "available"
                              ? "Còn chỗ"
                              : status === "full"
                                ? "Đã đủ"
                                : "Đã khóa"}
                    </button>
                {/each}
            </div>
            <div class="ml-auto flex items-center gap-2">
                <span class="text-sm text-gray-600">Sắp xếp:</span>
                <select
                    bind:value={sortBy}
                    class="rounded-lg border border-gray-300 px-3 py-1 text-sm focus:outline-none"
                >
                    <option value="recent">Mới nhất</option>
                    <option value="name">Tên A-Z</option>
                    <option value="students">Còn nhiều chỗ</option>
                </select>
            </div>
        </div>
    </div>

    {#if sortedProjects.length === 0}
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
            {#each sortedProjects as project (project.id)}
                {@const config = getStatusConfig(project.status)}
                <div
                    class="group relative rounded-lg border {config.border} p-5 transition-all hover:shadow-md"
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
                                        {project.name}
                                    </h3>
                                    <span
                                        class="rounded-full {config.bg} {config.color} px-3 py-0.5 text-xs font-medium"
                                    >
                                        {config.label}
                                    </span>
                                </div>
                                <p class="mb-3 text-sm text-gray-600">
                                    {project.description}
                                </p>
                                <div class="mb-3 flex flex-wrap gap-2">
                                    {#each project.tags as tag}
                                        <span
                                            class="flex items-center gap-1 rounded bg-gray-100 px-2 py-0.5 text-xs text-gray-700"
                                        >
                                            <Tag class="h-3 w-3" />
                                            {tag}
                                        </span>
                                    {/each}
                                </div>
                                <div
                                    class="flex items-center gap-4 text-xs text-gray-500"
                                >
                                    <div class="flex items-center gap-1">
                                        <Users class="h-4 w-4" />
                                        <span
                                            >{project.currentStudents}/{project.maxStudents}
                                            sinh viên</span
                                        >
                                    </div>
                                    <span>•</span>
                                    <span>GVHD: {project.instructor}</span>
                                </div>
                            </div>
                        </div>

                        <div class="flex flex-col items-end gap-4">
                            <div class="flex items-center gap-2">
                                <span class="text-xs text-gray-500"
                                    >Cho phép SV sửa:</span
                                >
                                <button
                                    onclick={() =>
                                        toggleAllowEdit(
                                            project.id,
                                            project.allowStudentEdit ?? false,
                                        )}
                                    class="relative h-5 w-9 rounded-full transition-colors {project.allowStudentEdit
                                        ? 'bg-blue-600'
                                        : 'bg-gray-300'}"
                                >
                                    <span
                                        class="absolute top-0.5 h-4 w-4 rounded-full bg-white transition-all {project.allowStudentEdit
                                            ? 'left-4.5'
                                            : 'left-0.5'}"
                                    ></span>
                                </button>
                            </div>
                            <div class="flex gap-1">
                                <button
                                    onclick={() => {
                                        push(
                                            `/lecture/my-classes/1/categories/${category.id}/projects/${project.id}`,
                                        );
                                    }}
                                    class="rounded p-2 text-blue-600 hover:bg-blue-50"
                                    title="Xem chi tiết"
                                    ><Eye class="h-5 w-5" /></button
                                >
                                <button
                                    onclick={() => (editingProject = project)}
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
                                        handleDuplicateProject(project)}
                                    class="rounded p-2 text-gray-600 hover:bg-gray-100"
                                    title="Sao chép"
                                    ><Copy class="h-5 w-5" /></button
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
        onClose={() => (showCreateModal = false)}
        onSubmit={handleCreateProject}
        categoryId={category.id}
    />
{/if}

{#if editingProject}
    <CreateProjectModal
        {editingProject}
        onClose={() => (editingProject = null)}
        onSubmit={handleUpdateProject}
        categoryId={category.id}
    />
{/if}
