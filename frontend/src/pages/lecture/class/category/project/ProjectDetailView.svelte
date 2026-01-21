<script lang="ts">
    import {
        ArrowLeft,
        Users,
        Plus,
        Trash2,
        Mail,
        Phone,
        Crown,
        UserCircle,
    } from "lucide-svelte";
    import { projectStore } from "../../../../../stores/project-store";
    import {
        projectRoundStore,
        currentProjectRound,
    } from "../../../../../stores/project-round-store";
    import { classStore } from "../../../../../stores/class-store"; // Import classStore
    import type { ClassroomResponse } from "../../../../../dtos";

    let { params } = $props();

    let classroomId: string = $derived(params.id);
    let projectRoundId: string = $derived(params.categoryId);
    let projectId: string = $derived(params.projectId);

    const { selectedProject, fetchProject } = projectStore;
    const { getRoundById } = projectRoundStore;

    let currentRound = $currentProjectRound;

    let classDetails = $state<ClassroomResponse | null>(null);

    $effect(() => {
        fetchProject(classroomId, projectId);
        getRoundById(classroomId, projectRoundId); // Fetch the current round
    });

    $effect(() => {
        (async () => {
            classDetails = await classStore.fetchClassById(classroomId);
        })();
    });

    //state
    //students in the selected project mockdata

    let students = $state([
        {
            id: "1",
            name: "Nguyễn Văn A",
            studentCode: "B1234567",
            email: "a@b.com",
            role: "leader",
            phone: "0123456789",
        },
        {
            id: "2",
            name: "Trần Thị B",
            studentCode: "B2345678",
            email: "b@c.com",
            role: "member",
            phone: "0987654321",
        },
    ]);
    let showAddStudentModal = $state(false);
    let showStudentProfileModal = $state(false);
</script>

<div class="space-y-6">
    <div class="bg-white rounded-lg shadow-sm p-6">
        <button
            onclick={() => history.back()}
            class="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-4 transition-colors"
        >
            <ArrowLeft class="w-5 h-5" />
            Quay lại danh sách đề tài
        </button>

        <div class="mb-4">
            <h1 class="text-2xl font-bold mb-2">{$selectedProject?.title}</h1>
            <p class="text-sm text-gray-500">
                Hạng mục: {currentRound?.name} | Lớp: {classDetails?.name} -
                {classDetails?.semester}
            </p>
            {#if classDetails}
                <p class="text-sm text-gray-500">
                    Giáo viên: {classDetails.lecturer.fullName} | Số lượng sinh viên:
                    {classDetails.students.length}
                </p>
            {/if}
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div class="bg-blue-50 rounded-lg p-4">
                <div class="text-blue-600 text-sm mb-1 font-medium">
                    Số lượng sinh viên
                </div>
                <div class="text-2xl font-bold text-blue-900">
                    <!-- {students.length}/{$selectedProject?.maxMember} -->
                </div>
            </div>

            <div class="bg-purple-50 rounded-lg p-4">
                <div class="text-purple-600 text-sm mb-1 font-medium">
                    Trạng thái
                </div>
                <div class="text-lg font-semibold text-purple-900">
                    {#if $selectedProject?.status === "available"}
                        Còn chỗ
                    {:else if $selectedProject?.status === "full"}
                        Đã đủ
                    {:else}
                        Đã khóa
                    {/if}
                </div>
            </div>
        </div>

        <div class="mt-4">
            <h3 class="font-bold mb-2 text-gray-800">Mô tả đề tài:</h3>
            <p class="text-gray-600 leading-relaxed">
                {$selectedProject?.description}
            </p>
        </div>

        <!-- {#if $selectedProject?.tags && $selectedProject.tags.length > 0}
            <div class="mt-4">
                <h3 class="text-sm font-bold mb-2 text-gray-700">Tags:</h3>
                <div class="flex flex-wrap gap-2">  
                    {#each $selectedProject.tags as tag}
                        <span
                            class="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm font-medium"
                        >
                            {tag}
                        </span>
                    {/each}
                </div>
            </div>
        {/if} -->
    </div>

    <div class="bg-white rounded-lg shadow-sm p-6">
        <div class="flex justify-between items-center mb-6">
            <h2 class="text-xl font-bold">
                Danh sách sinh viên ({students.length})
            </h2>
        </div>

        {#if students.length === 0}
            <div class="text-center py-12 text-gray-400 italic">
                Chưa có sinh viên nào trong đề tài này
            </div>
        {:else}
            <div class="space-y-3">
                {#each students as student (student.id)}
                    <div
                        class="border border-gray-100 rounded-xl p-4 hover:shadow-md transition-all duration-200 bg-white"
                    >
                        <div class="flex justify-between items-start">
                            <div class="flex-1">
                                <div class="flex items-center gap-3 mb-3">
                                    <h3
                                        class="text-lg font-semibold text-gray-900"
                                    >
                                        {student.name}
                                    </h3>
                                    {#if student.role === "leader"}
                                        <span
                                            class="flex items-center gap-1 px-3 py-1 bg-amber-50 text-amber-700 rounded-full text-xs font-bold border border-amber-100"
                                        >
                                            <Crown class="w-3 h-3" />
                                            Nhóm trưởng
                                        </span>
                                    {/if}
                                </div>

                                <div
                                    class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-sm text-gray-600"
                                >
                                    <div class="flex items-center gap-2">
                                        <Users class="w-4 h-4 text-gray-400" />
                                        <span
                                            >MSSV: <span
                                                class="font-medium text-gray-800"
                                                >{student.studentCode}</span
                                            ></span
                                        >
                                    </div>
                                    <div
                                        class="flex items-center gap-2 col-span-2 sm:col-span-1"
                                    >
                                        <Mail class="w-4 h-4 text-gray-400" />
                                        <span class="truncate"
                                            >{student.email}</span
                                        >
                                    </div>
                                    {#if student.phone}
                                        <div class="flex items-center gap-2">
                                            <Phone
                                                class="w-4 h-4 text-gray-400"
                                            />
                                            <span>{student.phone}</span>
                                        </div>
                                    {/if}
                                </div>
                            </div>

                            <div class="flex items-center gap-1 ml-4">
                                {#if student.role !== "leader"}
                                    <button
                                        class="px-3 py-1.5 text-xs font-medium border border-gray-200 rounded-lg hover:bg-gray-50 text-gray-700 transition-colors"
                                    >
                                        Làm nhóm trưởng
                                    </button>
                                {/if}
                                <button
                                    class="p-2 hover:bg-red-50 rounded-lg transition-colors group"
                                    title="Xóa sinh viên"
                                >
                                    <Trash2
                                        class="w-5 h-5 text-gray-400 group-hover:text-red-600"
                                    />
                                </button>
                                <button
                                    class="p-2 hover:bg-blue-50 rounded-lg transition-colors group"
                                    title="Hồ sơ sinh viên"
                                >
                                    <UserCircle
                                        class="w-5 h-5 text-gray-400 group-hover:text-blue-600"
                                    />
                                </button>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<!-- {#if showAddStudentModal}
    <AddStudentModal
        onClose={() => (showAddStudentModal = false)}
        onSubmit={handleAddStudent}
        existingStudents={students}
    />
{/if}

{#if showStudentProfileModal && selectedStudent}
    <StudentProfileModal
        onClose={() => (showStudentProfileModal = false)}
        student={selectedStudent}
    />
{/if} -->
