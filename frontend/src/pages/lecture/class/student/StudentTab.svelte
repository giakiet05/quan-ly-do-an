<script lang="ts">
  import {
    Search,
    Users,
    Mail,
    Phone,
    GraduationCap,
    Trash2,
  } from "@lucide/svelte";
  import type { StudentInClass } from "../../../../types/student";
  import { classStore } from "../../../../stores/class-store";
  import type { UserInfo } from "../../../../dtos";

  // Nhận props từ ClassDetail
  let { students, classroomId } = $props<{
    students: UserInfo[];
    classroomId: string;
  }>();
  // State quản lý tìm kiếm
  let searchTerm = $state("");

  // Sử dụng $derived để lọc sinh viên theo Search Term
  const filteredStudents = $derived(
    students.filter((student: StudentInClass) => {
      const search = searchTerm.toLowerCase();
      return (
        student.fullName.toLowerCase().includes(search) ||
        student.studentCode?.toLowerCase().includes(search) ||
        student.email?.toLowerCase().includes(search)
      );
    }),
  );

  // Helper lấy tên viết tắt cho Avatar
  const getInitials = (name: string): string => {
    const words = name.trim().split(" ");
    if (words.length >= 2) {
      return (words[0][0] + words[words.length - 1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  };

  // Xử lý xóa sinh viên
  async function removeStudent(studentId: string) {
    if (confirm("Bạn có chắc chắn muốn xóa sinh viên này khỏi lớp?")) {
      try {
        await classStore.removeStudent(classroomId, studentId);
        // Cập nhật lại danh sách sinh viên sau khi xóa
        students = students.filter(
          (student: UserInfo) => student.userId !== studentId,
        );
        alert("Xóa sinh viên thành công!");
      } catch (error) {
        console.error("Lỗi khi xóa sinh viên:", error);
        alert("Không thể xóa sinh viên. Vui lòng thử lại.");
      }
    }
  }
</script>

<div class="bg-white rounded-lg shadow-sm p-6">
  <div
    class="mb-6 flex flex-col md:flex-row md:items-center justify-between gap-4"
  >
    <div>
      <h2 class="text-xl font-bold text-gray-800 flex items-center gap-2">
        <Users class="w-6 h-6 text-blue-600" />
        Thành viên lớp học
      </h2>
      <p class="text-sm text-gray-500">Tổng số {students.length} sinh viên</p>
    </div>

    <div class="relative w-full md:w-96">
      <Search
        class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
      />
      <input
        type="text"
        bind:value={searchTerm}
        placeholder="Tìm kiếm sinh viên..."
        class="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
      />
    </div>
  </div>

  {#if filteredStudents.length === 0}
    <div
      class="text-center py-20 border-2 border-dashed border-gray-100 rounded-xl"
    >
      <Users class="w-16 h-16 text-gray-200 mx-auto mb-4" />
      <p class="text-gray-500 font-medium">
        Không tìm thấy sinh viên nào phù hợp
      </p>
      <button
        onclick={() => (searchTerm = "")}
        class="mt-2 text-blue-600 hover:underline text-sm"
      >
        Xóa tìm kiếm
      </button>
    </div>
  {:else}
    <div class="overflow-x-auto border border-gray-100 rounded-xl">
      <table class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th
              class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
              >STT</th
            >
            <th
              class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
              >Họ và Tên</th
            >
            <th
              class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
              >Mã số SV</th
            >
            <th
              class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
              >Thông tin liên hệ</th
            >
            <th
              class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
              >Trạng thái</th
            >
            <th
              class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
              >Hành động</th
            >
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          {#each filteredStudents as student, index (student.userId)}
            <tr class="hover:bg-blue-50/30 transition-colors group">
              <td class="px-6 py-4 text-sm text-gray-500">{index + 1}</td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div
                    class="w-9 h-9 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold shadow-sm"
                  >
                    {getInitials(student.fullName)}
                  </div>
                  <span class="text-sm font-medium text-gray-900"
                    >{student.fullName}</span
                  >
                </div>
              </td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800"
                >
                  {student.studentCode}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="flex flex-col gap-1">
                  <div
                    class="flex items-center gap-2 text-gray-600 group-hover:text-blue-600 transition-colors"
                  >
                    <Mail class="w-3.5 h-3.5" />
                    <span class="text-xs">{student.email}</span>
                  </div>
                  {#if student.phone}
                    <div class="flex items-center gap-2 text-gray-500">
                      <Phone class="w-3.5 h-3.5" />
                      <span class="text-xs">{student.phone}</span>
                    </div>
                  {/if}
                </div>
              </td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center gap-1.5 py-1 px-2.5 rounded-full text-xs font-medium bg-green-100 text-green-700"
                >
                  <span class="w-1.5 h-1.5 rounded-full bg-green-500"></span>
                  Đang học
                </span>
              </td>
              <td class="px-6 py-4">
                <button
                  class="text-red-600 hover:text-red-800 transition-colors flex items-center gap-1"
                  onclick={() => removeStudent(student.userId)}
                >
                  <Trash2 class="w-4 h-4" />
                  Xóa
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
