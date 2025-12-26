<script lang="ts">
    // Mock data
    let profile = {
        name: "TS. Nguyễn Văn A",
        mssv: "20181234",
        email: "an.nv181234@sis.hust.edu.vn",
        phone: "0987654321",
        class: "CNTT 2018.01",
        dob: "2000-01-15",
        department: "Khoa CNTT",
    };

    const role = "teacher";

    function saveChanges() {
        alert("Lưu thay đổi thành công (mock)!");
    }

    function changePassword() {
        alert("Đổi mật khẩu (mock)");
    }

    function handleCancel() {
        alert("Huỷ thay đổi (mock)");
    }
</script>

<div class="space-y-6">
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
                <div
                    class="w-[114px] h-[114px] rounded-full overflow-hidden bg-blue-600 flex items-center justify-center text-white text-3xl"
                >
                    {role === "teacher" ? "TS" : profile.name.charAt(0)}
                </div>
            </div>

            <!-- Name & Info -->
            <h3 class="text-xl text-slate-900 mb-2">{profile.name}</h3>
            <div class="text-center text-slate-500 mb-6">
                <p class="text-sm">MSSV: {profile.mssv}</p>
                <p class="text-sm">{profile.department}</p>
            </div>

            <!-- Change Photo Button -->
            <button
                class="w-full bg-slate-100 hover:bg-slate-200 rounded-lg py-2.5 px-4 text-sm text-slate-800 transition-colors mb-6"
            >
                Thay đổi ảnh
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
                            <label class="block text-sm text-slate-600 mb-2"
                                >Họ và Tên</label
                            >
                            <input
                                type="text"
                                bind:value={profile.name}
                                disabled
                                class="w-full px-3 py-2 bg-slate-100 border border-slate-200 rounded-lg text-slate-500"
                            />
                        </div>

                        <div>
                            <label class="block text-sm text-slate-600 mb-2"
                                >Mã số sinh viên</label
                            >
                            <input
                                type="text"
                                value={profile.mssv}
                                disabled
                                class="w-full px-3 py-2 bg-slate-100 border border-slate-200 rounded-lg text-slate-500"
                            />
                        </div>
                    </div>

                    <!-- Row 2 -->
                    <div class="grid grid-cols-2 gap-6">
                        <div>
                            <label class="block text-sm text-slate-600 mb-2"
                                >Email trường</label
                            >
                            <input
                                type="email"
                                value={profile.email}
                                disabled
                                class="w-full px-3 py-2 bg-slate-100 border border-slate-200 rounded-lg text-slate-500"
                            />
                        </div>

                        <div>
                            <label class="block text-sm text-slate-600 mb-2"
                                >Số điện thoại</label
                            >
                            <input
                                type="tel"
                                bind:value={profile.phone}
                                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>

                    <!-- Row 3 -->
                    <div class="grid grid-cols-2 gap-6">
                        <div>
                            <label class="block text-sm text-slate-600 mb-2"
                                >Lớp sinh hoạt</label
                            >
                            <input
                                type="text"
                                value={profile.class}
                                disabled
                                class="w-full px-3 py-2 bg-slate-100 border border-slate-200 rounded-lg text-slate-500"
                            />
                        </div>

                        <div>
                            <label class="block text-sm text-slate-600 mb-2"
                                >Ngày sinh</label
                            >
                            <input
                                type="date"
                                bind:value={profile.dob}
                                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>
                </div>
            </div>

            <!-- Account UserSettings Card -->
            <div class="bg-white border border-slate-200 rounded-xl shadow-sm">
                <div class="border-b border-slate-200 px-6 py-5">
                    <h4 class="text-lg text-slate-900">Cài đặt tài khoản</h4>
                </div>

                <div class="px-6 py-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <h5 class="text-slate-800 mb-1">Đổi mật khẩu</h5>
                            <p class="text-sm text-slate-500">
                                Thay đổi mật khẩu đăng nhập của bạn.
                            </p>
                        </div>
                        <button
                            class="bg-slate-100 hover:bg-slate-200 rounded-lg px-6 py-2 text-sm text-slate-800 transition-colors"
                            on:click={changePassword}
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
        <button
            on:click={handleCancel}
            class="bg-slate-200 hover:bg-slate-300 rounded-lg px-6 py-2.5 text-sm text-slate-800 transition-colors"
        >
            Hủy
        </button>
        <button
            on:click={saveChanges}
            class="bg-[#2b8cee] hover:bg-blue-600 rounded-lg px-6 py-2.5 text-sm text-white transition-colors"
        >
            Lưu thay đổi
        </button>
    </div>
</div>
