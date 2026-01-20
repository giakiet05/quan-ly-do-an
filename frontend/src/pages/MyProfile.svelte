<script lang="ts">
  import { onMount } from "svelte";
  import {
    getUserProfile,
    updateUserProfile,
    uploadAvatar,
    deleteAvatar,
    changePassword as changePasswordAPI,
    type UserProfile,
  } from "../services/user-service";
  import { authStore } from "../stores/auth-store";

  // State using Svelte 5 runes
  let profile = $state<UserProfile | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let saving = $state(false);

  // Editable fields
  let editedFullName = $state("");
  let editedStudentCode = $state("");

  // Avatar upload
  let avatarInput: HTMLInputElement;
  let uploadingAvatar = $state(false);

  // Password change modal
  let showPasswordModal = $state(false);
  let oldPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");
  let changingPassword = $state(false);

  // Load profile on mount
  onMount(async () => {
    try {
      loading = true;
      error = null;
      const data = await getUserProfile();
      profile = data;
      editedFullName = data.fullName;
      editedStudentCode = data.studentCode || "";

      // Update authStore with latest profile
      authStore.update((state) => ({
        ...state,
        user: {
          id: data.id,
          avatar: data.avatar?.url || "",
          email: data.email,
          fullname: data.fullName,
          role: "user",
          is_verified: true,
        },
      }));
    } catch (err) {
      error =
        err instanceof Error ? err.message : "Không thể tải thông tin hồ sơ";
      console.error("Error loading profile:", err);
    } finally {
      loading = false;
    }
  });

  async function saveChanges() {
    if (!profile) return;

    try {
      saving = true;
      error = null;

      // Only send changed fields
      const updates: { fullName?: string; studentCode?: string } = {};
      if (editedFullName !== profile.fullName) {
        updates.fullName = editedFullName;
      }
      if (editedStudentCode !== (profile.studentCode || "")) {
        updates.studentCode = editedStudentCode;
      }

      if (Object.keys(updates).length === 0) {
        alert("Không có thay đổi nào để lưu");
        return;
      }

      const updatedProfile = await updateUserProfile(updates);
      profile = updatedProfile;
      editedFullName = updatedProfile.fullName;
      editedStudentCode = updatedProfile.studentCode || "";

      // Update authStore
      authStore.update((state) => ({
        ...state,
        user: state.user
          ? {
              ...state.user,
              fullname: updatedProfile.fullName,
            }
          : null,
      }));

      alert("Lưu thay đổi thành công!");
    } catch (err) {
      error = err instanceof Error ? err.message : "Không thể lưu thay đổi";
      alert(`Lỗi: ${error}`);
      console.error("Error saving profile:", err);
    } finally {
      saving = false;
    }
  }

  function openPasswordModal() {
    showPasswordModal = true;
    oldPassword = "";
    newPassword = "";
    confirmPassword = "";
  }

  function closePasswordModal() {
    showPasswordModal = false;
    oldPassword = "";
    newPassword = "";
    confirmPassword = "";
  }

  async function handleChangePassword() {
    if (!newPassword || newPassword.length < 6) {
      alert("Mật khẩu mới phải có ít nhất 6 ký tự");
      return;
    }

    if (newPassword !== confirmPassword) {
      alert("Mật khẩu xác nhận không khớp");
      return;
    }

    try {
      changingPassword = true;
      error = null;
      await changePasswordAPI(oldPassword, newPassword);
      alert("Đổi mật khẩu thành công!");
      closePasswordModal();
    } catch (err) {
      error = err instanceof Error ? err.message : "Không thể đổi mật khẩu";
      alert(`Lỗi: ${error}`);
      console.error("Error changing password:", err);
    } finally {
      changingPassword = false;
    }
  }

  function openAvatarPicker() {
    avatarInput?.click();
  }

  async function handleAvatarChange(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];

    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      alert("Vui lòng chọn file ảnh");
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert("Kích thước ảnh không được vượt quá 5MB");
      return;
    }

    try {
      uploadingAvatar = true;
      error = null;
      const updatedProfile = await uploadAvatar(file);
      profile = updatedProfile;

      // Update authStore with new avatar
      authStore.update((state) => ({
        ...state,
        user: state.user
          ? {
              ...state.user,
              avatar: updatedProfile.avatar?.url || "",
            }
          : null,
      }));

      alert("Cập nhật ảnh đại diện thành công!");
    } catch (err) {
      error = err instanceof Error ? err.message : "Không thể tải ảnh lên";
      alert(`Lỗi: ${error}`);
      console.error("Error uploading avatar:", err);
    } finally {
      uploadingAvatar = false;
      // Reset input
      target.value = "";
    }
  }

  async function handleDeleteAvatar() {
    if (!confirm("Bạn có chắc muốn xóa ảnh đại diện?")) return;

    try {
      uploadingAvatar = true;
      error = null;
      const updatedProfile = await deleteAvatar();
      profile = updatedProfile;

      // Update authStore
      authStore.update((state) => ({
        ...state,
        user: state.user
          ? {
              ...state.user,
              avatar: "",
            }
          : null,
      }));

      alert("Đã xóa ảnh đại diện");
    } catch (err) {
      error = err instanceof Error ? err.message : "Không thể xóa ảnh";
      alert(`Lỗi: ${error}`);
      console.error("Error deleting avatar:", err);
    } finally {
      uploadingAvatar = false;
    }
  }

  function handleCancel() {
    if (!profile) return;
    // Reset to original values
    editedFullName = profile.fullName;
    editedStudentCode = profile.studentCode || "";
    alert("Đã hủy thay đổi");
  }

  // Derived values
  const avatarLetter = $derived(
    profile?.fullName?.charAt(0)?.toUpperCase() || "U",
  );
  const displayRole = $derived(
    profile?.provider === "google" ? "Google" : "Local",
  );
</script>

<div class="space-y-6">
  <!-- Loading State -->
  {#if loading}
    <div class="flex items-center justify-center h-[400px]">
      <div class="text-center">
        <div
          class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"
        ></div>
        <p class="text-slate-600">Đang tải thông tin...</p>
      </div>
    </div>
  {:else if error && !profile}
    <!-- Error State -->
    <div class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-800">❌ {error}</p>
    </div>
  {:else if profile}
    <!-- Header -->
    <div>
      <h1 class="text-4xl tracking-tight mb-2">Hồ sơ của bạn</h1>
      <p class="text-slate-500">
        Quản lý thông tin cá nhân và cài đặt tài khoản của bạn.
      </p>
    </div>

    <div class="flex gap-6">
      <!-- Left Card - Profile Summary -->
      <div
        class="bg-white border border-slate-200 rounded-xl shadow-sm w-[398px] h-[360px] flex flex-col items-center p-6"
      >
        <!-- Avatar -->
        <div class="relative mt-4 mb-6">
          {#if profile.avatar}
            <img
              src={profile.avatar.url}
              alt="Avatar"
              class="w-[114px] h-[114px] rounded-full object-cover"
            />
          {:else}
            <div
              class="w-[114px] h-[114px] rounded-full overflow-hidden bg-blue-600 flex items-center justify-center text-white text-3xl font-semibold"
            >
              {avatarLetter}
            </div>
          {/if}
        </div>

        <!-- Name & Info -->
        <h3 class="text-xl text-slate-900 mb-2">{profile.fullName}</h3>
        <div class="text-center text-slate-500 mb-6">
          <p class="text-sm">Email: {profile.email}</p>
          {#if profile.studentCode}
            <p class="text-sm">MSSV: {profile.studentCode}</p>
          {/if}
          <p class="text-xs text-slate-400 mt-1">
            Đăng nhập qua: {displayRole}
          </p>
        </div>

        <!-- Change Photo Button -->
        <input
          type="file"
          accept="image/*"
          bind:this={avatarInput}
          onchange={handleAvatarChange}
          class="hidden"
        />
        <button
          onclick={openAvatarPicker}
          disabled={uploadingAvatar}
          class="w-full bg-slate-100 hover:bg-slate-200 rounded-lg py-2.5 px-4 text-sm text-slate-800 transition-colors disabled:opacity-50 mb-6"
        >
          {uploadingAvatar ? "Đang Tải..." : "Thay Đổi Ảnh"}
        </button>
      </div>

      <!-- Right Cards -->
      <div class="flex-1 space-y-6">
        <!-- Personal Information Card -->
        <div class="bg-white border border-slate-200 rounded-xl shadow-sm">
          <div class="border-b border-slate-200 px-6 py-5">
            <h4 class="text-lg text-slate-900">Thông tin cá nhân</h4>
          </div>

          <div class="p-6 space-y-6">
            <!-- Row 1 -->
            <div class="grid grid-cols-2 gap-6">
              <div>
                <label class="block text-sm text-slate-600 mb-2 font-semibold"
                  >Họ Và Tên</label
                >
                <input
                  type="text"
                  bind:value={editedFullName}
                  oninput={(e) =>
                    console.log("Input value:", e.currentTarget.value)}
                  autocomplete="name"
                  spellcheck="false"
                  style="text-transform: none !important;"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm text-slate-600 mb-2 font-semibold"
                  >Mã Số Sinh Viên</label
                >
                <input
                  type="text"
                  bind:value={editedStudentCode}
                  placeholder="Nhập MSSV Của Bạn"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            <!-- Row 2 -->
            <div class="grid grid-cols-2 gap-6">
              <div>
                <label class="block text-sm text-slate-600 mb-2 font-semibold"
                  >Email</label
                >
                <input
                  type="email"
                  value={profile.email}
                  disabled
                  class="w-full px-3 py-2 bg-slate-100 border border-slate-200 rounded-lg text-slate-500"
                />
                <p class="text-xs text-slate-400 mt-1">
                  Email không thể thay đổi
                </p>
              </div>

              <div>
                <label class="block text-sm text-slate-600 mb-2 font-semibold"
                  >Tài Khoản Tạo Lúc</label
                >
                <input
                  type="text"
                  value={new Date(profile.createdAt).toLocaleDateString(
                    "vi-VN",
                  )}
                  disabled
                  class="w-full px-3 py-2 bg-slate-100 border border-slate-200 rounded-lg text-slate-500"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Account Settings Card -->
        <div class="bg-white border border-slate-200 rounded-xl shadow-sm">
          <div class="border-b border-slate-200 px-6 py-5">
            <h4 class="text-lg text-slate-900">Cài đặt tài khoản</h4>
          </div>

          <div class="px-6 py-6">
            <div class="flex items-center justify-between">
              <div>
                <h5 class="text-slate-800 mb-1">Đổi mật khẩu</h5>
                <p class="text-sm text-slate-500">
                  {#if profile.provider === "local"}
                    Thay đổi mật khẩu đăng nhập của bạn.
                  {:else}
                    Bạn đăng nhập qua Google, không thể đổi mật khẩu.
                  {/if}
                </p>
              </div>
              <button
                class="bg-slate-100 hover:bg-slate-200 rounded-lg px-6 py-2 text-sm text-slate-800 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={openPasswordModal}
                disabled={profile.provider !== "local"}
              >
                Đổi
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex justify-end gap-3">
      {#if error}
        <div class="flex-1 text-red-600 text-sm">
          ⚠️ {error}
        </div>
      {/if}
      <button
        onclick={handleCancel}
        class="bg-slate-200 hover:bg-slate-300 rounded-lg px-6 py-2.5 text-sm text-slate-800 transition-colors disabled:opacity-50"
        disabled={saving}
      >
        Hủy
      </button>
      <button
        onclick={saveChanges}
        class="bg-[#2b8cee] hover:bg-blue-600 rounded-lg px-6 py-2.5 text-sm text-white transition-colors disabled:opacity-50"
        disabled={saving}
      >
        {saving ? "Đang lưu..." : "Lưu thay đổi"}
      </button>
    </div>
  {/if}
</div>

<!-- Password Change Modal -->
{#if showPasswordModal}
  <div
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    onclick={closePasswordModal}
  >
    <div
      class="bg-white rounded-xl shadow-xl p-6 w-[450px] max-w-[90vw]"
      onclick={(e) => e.stopPropagation()}
    >
      <h3 class="text-xl text-slate-900 mb-4">Đổi mật khẩu</h3>

      <div class="space-y-4">
        <div>
          <label class="block text-sm text-slate-600 mb-2 font-semibold"
            >Mật Khẩu Hiện Tại</label
          >
          <input
            type="password"
            bind:value={oldPassword}
            placeholder="Nhập mật khẩu hiện tại"
            class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm text-slate-600 mb-2 font-semibold"
            >Mật Khẩu Mới</label
          >
          <input
            type="password"
            bind:value={newPassword}
            placeholder="Nhập mật khẩu mới (tối thiểu 6 ký tự)"
            class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm text-slate-600 mb-2 font-semibold"
            >Xác Nhận Mật Khẩu Mới</label
          >
          <input
            type="password"
            bind:value={confirmPassword}
            placeholder="Nhập lại mật khẩu mới"
            class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-6">
        <button
          onclick={closePasswordModal}
          class="bg-slate-200 hover:bg-slate-300 rounded-lg px-6 py-2.5 text-sm text-slate-800 transition-colors"
          disabled={changingPassword}
        >
          Hủy
        </button>
        <button
          onclick={handleChangePassword}
          class="bg-[#2b8cee] hover:bg-blue-600 rounded-lg px-6 py-2.5 text-sm text-white transition-colors disabled:opacity-50"
          disabled={changingPassword}
        >
          {changingPassword ? "Đang xử lý..." : "Đổi mật khẩu"}
        </button>
      </div>
    </div>
  </div>
{/if}
