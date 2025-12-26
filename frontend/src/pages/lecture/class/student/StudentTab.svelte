<script lang="ts">
  import {
    Search,
    Users,
    Mail,
    Phone,
    CheckCircle,
    XCircle,
    Filter,
    FileText,
  } from "@lucide/svelte";
  import type { ProjectCategory } from "../../../../types/category";
  import type { StudentInClass } from "../../../../types/student";

  // Khai báo Props bằng $props()
  let { categories, students } = $props<{
    categories: ProjectCategory[];
    students: StudentInClass[];
  }>();

  // State quản lý UI
  let searchTerm = $state("");
  let selectedCategoryId = $state(categories[0]?.id || "");
  let statusFilter = $state<"all" | "enrolled" | "not-enrolled">("all");

  // Hàm bổ trợ: Tìm dự án của sinh viên trong hạng mục đang chọn
  function getProjectInCategory(student: StudentInClass) {
    return student.enrolledProjects.find(
      (p) => p.categoryId === selectedCategoryId,
    );
  }

  // Sử dụng $derived để tự động tính toán danh sách sinh viên khi state thay đổi
  const filteredStudents = $derived(
    students.filter((student: StudentInClass) => {
      const matchesSearch =
        student.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        student.studentCode.toLowerCase().includes(searchTerm.toLowerCase()) ||
        student.email.toLowerCase().includes(searchTerm.toLowerCase());

      if (!matchesSearch) return false;

      const projectInCategory = getProjectInCategory(student);
      if (statusFilter === "enrolled" && !projectInCategory) return false;
      if (statusFilter === "not-enrolled" && projectInCategory) return false;

      return true;
    }),
  );

  // Thống kê dựa trên hạng mục đang chọn
  const enrolledCount = $derived(
    students.filter((s: StudentInClass) => getProjectInCategory(s)).length,
  );
  const notEnrolledCount = $derived(students.length - enrolledCount);

  // Helper lấy tên viết tắt
  const getInitials = (name: string): string => {
    const words = name.split(" ");
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  };
</script>

<div class="bg-white rounded-lg shadow-sm p-6">
  <div class="mb-6">
    <h2 class="text-xl mb-4">Danh sách sinh viên theo hạng mục</h2>

    <div class="flex gap-2 mb-6 overflow-x-auto pb-2">
      {#each categories as category (category.id)}
        <button
          onclick={() => {
            selectedCategoryId = category.id;
            statusFilter = "all";
          }}
          class="px-4 py-2 rounded-lg text-sm whitespace-nowrap transition-colors {selectedCategoryId ===
          category.id
            ? 'bg-blue-600 text-white'
            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
        >
          {category.name}
        </button>
      {/each}
    </div>

    <div class="space-y-4">
      <div class="relative">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
        />
        <input
          type="text"
          bind:value={searchTerm}
          placeholder="Tìm kiếm theo tên, MSSV hoặc email..."
          class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div class="flex gap-2 items-center">
        <Filter class="w-4 h-4 text-gray-600" />
        <span class="text-sm text-gray-600">Trạng thái trong hạng mục này:</span
        >
        <div class="flex gap-2">
          <button
            onclick={() => (statusFilter = "all")}
            class="px-3 py-1 rounded-lg text-sm transition-colors {statusFilter ===
            'all'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
          >
            Tất cả ({students.length})
          </button>
          <button
            onclick={() => (statusFilter = "enrolled")}
            class="px-3 py-1 rounded-lg text-sm transition-colors {statusFilter ===
            'enrolled'
              ? 'bg-green-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
          >
            Đã tham gia ({enrolledCount})
          </button>
          <button
            onclick={() => (statusFilter = "not-enrolled")}
            class="px-3 py-1 rounded-lg text-sm transition-colors {statusFilter ===
            'not-enrolled'
              ? 'bg-orange-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
          >
            Chưa tham gia ({notEnrolledCount})
          </button>
        </div>
      </div>
    </div>
  </div>

  {#if filteredStudents.length === 0}
    <div class="text-center py-16">
      <Users class="w-16 h-16 text-gray-300 mx-auto mb-4" />
      <p class="text-gray-600">
        {searchTerm || statusFilter !== "all"
          ? "Không tìm thấy sinh viên nào"
          : "Chưa có sinh viên nào trong lớp"}
      </p>
      {#if searchTerm || statusFilter !== "all"}
        <p class="text-sm text-gray-500">Thử thay đổi bộ lọc hoặc tìm kiếm</p>
      {/if}
    </div>
  {:else}
    <div class="overflow-x-auto border border-gray-200 rounded-lg">
      <table class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th
              class="px-4 py-3 text-left text-sm text-gray-700 border-b border-gray-200"
              >STT</th
            >
            <th
              class="px-4 py-3 text-left text-sm text-gray-700 border-b border-gray-200"
              >Sinh viên</th
            >
            <th
              class="px-4 py-3 text-left text-sm text-gray-700 border-b border-gray-200"
              >MSSV</th
            >
            <th
              class="px-4 py-3 text-left text-sm text-gray-700 border-b border-gray-200"
              >Liên hệ</th
            >
            <th
              class="px-4 py-3 text-left text-sm text-gray-700 border-b border-gray-200"
              >Trạng thái</th
            >
            <th
              class="px-4 py-3 text-left text-sm text-gray-700 border-b border-gray-200"
              >Đề tài đã đăng ký</th
            >
          </tr>
        </thead>
        <tbody>
          {#each filteredStudents as student, index (student.id)}
            {@const projectInCategory = getProjectInCategory(student)}
            <tr class="hover:bg-gray-50 transition-colors">
              <td
                class="px-4 py-4 text-sm text-gray-900 border-b border-gray-200"
                >{index + 1}</td
              >
              <td class="px-4 py-4 border-b border-gray-200">
                <div class="flex items-center gap-3">
                  <div
                    class="flex-shrink-0 w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center text-white text-sm"
                  >
                    {getInitials(student.name)}
                  </div>
                  <div class="text-sm text-gray-900">{student.name}</div>
                </div>
              </td>
              <td
                class="px-4 py-4 text-sm text-gray-900 border-b border-gray-200"
                >{student.studentCode}</td
              >
              <td class="px-4 py-4 text-sm border-b border-gray-200">
                <div class="space-y-1">
                  <div class="flex items-center gap-2 text-gray-600">
                    <Mail class="w-4 h-4" />
                    <span class="text-xs">{student.email}</span>
                  </div>
                  {#if student.phone}
                    <div class="flex items-center gap-2 text-gray-600">
                      <Phone class="w-4 h-4" />
                      <span class="text-xs">{student.phone}</span>
                    </div>
                  {/if}
                </div>
              </td>
              <td class="px-4 py-4 border-b border-gray-200">
                {#if projectInCategory}
                  <span
                    class="inline-flex items-center gap-1 px-3 py-1 bg-green-50 text-green-700 rounded-full text-sm"
                  >
                    <CheckCircle class="w-4 h-4" /> Đã tham gia
                  </span>
                {:else}
                  <span
                    class="inline-flex items-center gap-1 px-3 py-1 bg-orange-50 text-orange-700 rounded-full text-sm"
                  >
                    <XCircle class="w-4 h-4" /> Chưa tham gia
                  </span>
                {/if}
              </td>
              <td class="px-4 py-4 border-b border-gray-200">
                {#if projectInCategory}
                  <div class="text-sm flex items-start gap-2">
                    <FileText
                      class="w-4 h-4 text-blue-600 flex-shrink-0 mt-0.5"
                    />
                    <div class="flex-1 min-w-0">
                      <div class="text-gray-900">
                        {projectInCategory.projectName}
                        {#if projectInCategory.role === "leader"}
                          <span
                            class="ml-2 px-2 py-0.5 bg-yellow-100 text-yellow-800 rounded text-xs"
                            >Nhóm trưởng</span
                          >
                        {/if}
                      </div>
                    </div>
                  </div>
                {:else}
                  <span class="text-sm text-gray-400">Chưa đăng ký</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
