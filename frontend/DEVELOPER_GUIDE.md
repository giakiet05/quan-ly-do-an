# 📘 Hướng dẫn Phát triển Frontend - Quản lý Đồ án

> Tài liệu này hướng dẫn chi tiết cho developer mới tham gia dự án về cách làm việc với Authentication và API.

## 📂 Cấu trúc Project quan trọng

```
frontend/src/
├── dtos/                    # Type definitions cho API
│   ├── auth-dto.ts         # Auth request/response types
│   ├── user-dto.ts         # User types
│   ├── classroom-dto.ts    # Classroom types
│   └── api-response-dto.ts # Generic API response wrapper
├── services/               # Business logic & API calls
│   ├── auth-service.ts     # ⭐ Auth functions (login, register, etc.)
│   ├── user-service.ts     # User-related API calls
│   └── classroom-service.ts
├── utils/
│   └── api-fetch.ts        # ⭐⭐ Core API wrapper (XEM ĐẦU TIÊN)
├── stores/                 # Svelte stores (global state)
│   └── auth-store.ts       # Auth state management
└── pages/                  # UI components
```

---

## 🔐 PHẦN 1: Authentication - Bắt đầu từ đây

### 1.1. File quan trọng nhất: `api-fetch.ts`

**Đường dẫn:** `src/utils/api-fetch.ts`

Đây là **wrapper function** bao bọc `fetch` API, tự động xử lý:

- ✅ Thêm Bearer token vào header
- ✅ Tự động refresh token khi hết hạn (401)
- ✅ Retry request sau khi refresh thành công
- ✅ Xử lý lỗi network

#### Cách sử dụng cơ bản:

```typescript
import { apiFetch } from "../utils/api-fetch";

// GET request với auth tự động
const data = await apiFetch<UserResponse>("/api/users/me");

// POST request
const result = await apiFetch<CreateResponse>("/api/classrooms", {
  method: "POST",
  body: JSON.stringify({ name: "Lớp mới", description: "..." }),
});

// Request KHÔNG cần auth (login, register)
const loginData = await apiFetch<LoginResponse>("/api/auth/login", {
  method: "POST",
  skipAuth: true, // ⭐ Quan trọng: bỏ qua Bearer token
  body: JSON.stringify({ email, password }),
});
```

#### Token Flow tự động:

```
Request → apiFetch
  ├─ Có token? → Thêm Authorization: Bearer {token}
  ├─ Call API
  ├─ 401 Unauthorized?
  │   ├─ Có refresh_token?
  │   ├─ Call /api/auth/refresh
  │   ├─ Lưu token mới
  │   └─ Retry request ban đầu
  └─ Trả về data
```

---

### 1.2. File thứ 2: `auth-service.ts`

**Đường dẫn:** `src/services/auth-service.ts`

Chứa các function auth cấp cao, sử dụng `apiFetch` bên trong.

#### Các function chính:

```typescript
// 1. Login
export async function login(email: string, password: string) {
  const response = await apiFetch<ApiResponse<LoginResponseDto>>(
    "/api/auth/login",
    {
      method: "POST",
      skipAuth: true, // Login không cần token
      body: JSON.stringify({ email, password })
    }
  );

  // Lưu token vào localStorage
  if (response.success && response.data) {
    setAccessToken(response.data.access_token);
    setRefreshToken(response.data.refresh_token);
    setUser(response.data.user);
  }

  return response;
}

// 2. Register (3 bước)
// Bước 1: Gửi email verification
export async function sendVerificationEmail(email: string, role: string) {
  return apiFetch<ApiResponse<void>>("/api/auth/send-verification", {
    method: "POST",
    skipAuth: true,
    body: JSON.stringify({ email, role })
  });
}

// Bước 2: Verify OTP
export async function verifyEmailOtp(email: string, otp: string) {
  return apiFetch<ApiResponse<VerifyEmailResponseDto>>(
    "/api/auth/verify-email",
    {
      method: "POST",
      skipAuth: true,
      body: JSON.stringify({ email, otp })
    }
  );
}

// Bước 3: Hoàn thành đăng ký
export async function completeRegistration(
  verification_token: string,
  full_name: string,
  password: string,
  // ... các field khác
) {
  return apiFetch<ApiResponse<RegisterResponseDto>>(
    "/api/auth/register",
    {
      method: "POST",
      skipAuth: true,
      body: JSON.stringify({ verification_token, full_name, password, ... })
    }
  );
}

// 3. Logout
export async function logout() {
  const refreshToken = getRefreshToken();
  await apiFetch("/api/auth/logout", {
    method: "POST",
    body: JSON.stringify({ refresh_token: refreshToken })
  });

  // Xóa token khỏi localStorage
  clearAccessToken();
  clearRefreshToken();
  clearUser();
}

// 4. Forgot Password (2 bước)
export async function forgotPassword(email: string) { ... }
export async function resetPassword(token: string, newPassword: string) { ... }

// 5. Google OAuth
export async function loginWithGoogle() { ... }
```

---

### 1.3. LocalStorage Keys

Token được lưu trong localStorage với các key sau:

```typescript
localStorage.getItem("access_token"); // JWT access token (15 phút)
localStorage.getItem("refresh_token"); // Refresh token (7 ngày)
localStorage.getItem("user"); // User info (JSON string)
```

**⚠️ Lưu ý:** Không bao giờ lưu password trong localStorage!

---

## 🚀 PHẦN 2: Viết API Call mới

### Bước 1: Tạo DTO Types

File: `src/dtos/your-feature-dto.ts`

```typescript
// REQUEST DTOs
export interface CreateProjectRequest {
  title: string;
  description: string;
  category_id: string;
}

export interface UpdateProjectRequest {
  title?: string;
  description?: string;
}

// RESPONSE DTOs
export interface ProjectResponse {
  id: string;
  title: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

// API Response wrapper (generic)
export type ApiProjectResponse = ApiResponse<ProjectResponse>;
```

### Bước 2: Tạo Service

File: `src/services/project-service.ts`

```typescript
import { apiFetch } from "../utils/api-fetch";
import type { ApiResponse } from "../dtos/api-response-dto";
import type {
  CreateProjectRequest,
  ProjectResponse,
} from "../dtos/project-dto";

// GET: Lấy danh sách
export async function getProjects(classroomId: string) {
  return apiFetch<ApiResponse<ProjectResponse[]>>(
    `/api/classrooms/${classroomId}/projects`
  );
}

// GET: Lấy chi tiết
export async function getProjectById(projectId: string) {
  return apiFetch<ApiResponse<ProjectResponse>>(`/api/projects/${projectId}`);
}

// POST: Tạo mới
export async function createProject(
  classroomId: string,
  data: CreateProjectRequest
) {
  return apiFetch<ApiResponse<ProjectResponse>>(
    `/api/classrooms/${classroomId}/projects`,
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
}

// PUT: Cập nhật
export async function updateProject(
  projectId: string,
  data: UpdateProjectRequest
) {
  return apiFetch<ApiResponse<ProjectResponse>>(`/api/projects/${projectId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

// DELETE: Xóa
export async function deleteProject(projectId: string) {
  return apiFetch<ApiResponse<void>>(`/api/projects/${projectId}`, {
    method: "DELETE",
  });
}
```

### Bước 3: Sử dụng trong Component

File: `src/pages/ProjectList.svelte`

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { getProjects } from "../services/project-service";
  import type { ProjectResponse } from "../dtos/project-dto";

  let projects = $state<ProjectResponse[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      loading = true;
      const response = await getProjects("classroom-id-123");

      if (response.success && response.data) {
        projects = response.data;
      } else {
        error = response.message || "Lỗi không xác định";
      }
    } catch (err: any) {
      error = err.message || "Không thể tải danh sách";
      console.error("Error fetching projects:", err);
    } finally {
      loading = false;
    }
  });

  async function handleCreate() {
    try {
      const newProject = await createProject("classroom-id-123", {
        title: "Dự án mới",
        description: "Mô tả...",
        category_id: "cat-1"
      });

      if (newProject.success && newProject.data) {
        projects = [...projects, newProject.data];
        alert("Tạo thành công!");
      }
    } catch (err: any) {
      alert(err.message || "Tạo thất bại");
    }
  }
</script>

{#if loading}
  <p>Đang tải...</p>
{:else if error}
  <p class="error">{error}</p>
{:else}
  <div>
    <button onclick={handleCreate}>Tạo dự án mới</button>
    {#each projects as project}
      <div class="project-card">
        <h3>{project.title}</h3>
        <p>{project.description}</p>
      </div>
    {/each}
  </div>
{/if}
```

---

## 📋 PHẦN 3: Xử lý lỗi

### Các loại lỗi thường gặp:

```typescript
try {
  const data = await apiFetch("/api/endpoint");
} catch (error: any) {
  // 1. Network error (mất mạng)
  if (error.message === "Failed to fetch") {
    console.error("Không có kết nối internet");
  }

  // 2. Backend error response
  // error structure: { message: string, error_code?: string }
  if (error.message) {
    console.error("Backend error:", error.message);
  }

  // 3. 401 Unauthorized (đã xử lý tự động trong apiFetch)
  // Nếu vẫn lỗi 401 sau retry → token hết hạn hoàn toàn
  if (error.status === 401) {
    // Redirect to login
    window.location.href = "/login";
  }

  // 4. 403 Forbidden (không có quyền)
  if (error.status === 403) {
    alert("Bạn không có quyền thực hiện thao tác này");
  }
}
```

---

## 🔧 PHẦN 4: Upload File

### Upload với FormData:

```typescript
import { apiFetch, uploadFile } from "../utils/api-fetch";

// Sử dụng helper uploadFile
export async function uploadAvatar(file: File) {
  return uploadFile<ApiResponse<{ file_url: string }>>(
    "/api/users/avatar",
    file,
    "avatar" // field name
  );
}

// Hoặc tự viết với FormData
async function uploadDocument(file: File, projectId: string) {
  const formData = new FormData();
  formData.append("document", file);
  formData.append("project_id", projectId);

  return apiFetch<ApiResponse<DocumentResponse>>("/api/projects/documents", {
    method: "POST",
    body: formData,
    // ⚠️ Không set Content-Type, browser tự set với boundary
  });
}
```

### Trong Component:

```svelte
<script lang="ts">
  let fileInput: HTMLInputElement;

  async function handleUpload() {
    const file = fileInput.files?.[0];
    if (!file) return;

    try {
      const response = await uploadAvatar(file);
      if (response.success) {
        alert("Upload thành công!");
      }
    } catch (err) {
      console.error("Upload failed:", err);
    }
  }
</script>

<input type="file" bind:this={fileInput} accept="image/*" />
<button onclick={handleUpload}>Upload Avatar</button>
```

---

## 🔒 PHẦN 5: Protected Routes

### Kiểm tra auth trước khi render:

```svelte
<!-- src/pages/ProtectedPage.svelte -->
<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { isAuthenticated } from "../stores/auth-store";

  onMount(() => {
    // Nếu chưa login, redirect về login
    if (!$isAuthenticated) {
      push("/login");
    }
  });
</script>

{#if $isAuthenticated}
  <div>
    <!-- Nội dung trang -->
  </div>
{:else}
  <p>Đang kiểm tra xác thực...</p>
{/if}
```

---

## 📝 PHẦN 6: Best Practices

### ✅ DO (Nên làm):

1. **Luôn dùng `apiFetch`** thay vì `fetch` trực tiếp
2. **Define types** trong `dtos/` trước khi viết service
3. **Xử lý error** với try-catch trong component
4. **Loading state**: Hiển thị trạng thái đang tải
5. **Validate input** trước khi gửi API
6. **Log errors** để debug: `console.error("Context:", error)`

### ❌ DON'T (Không nên):

1. ❌ Lưu password trong localStorage/state
2. ❌ Hard-code token trong code
3. ❌ Gọi API trực tiếp với `fetch()` (bỏ qua auth flow)
4. ❌ Bỏ qua xử lý lỗi (empty catch block)
5. ❌ Set `Content-Type: multipart/form-data` khi upload (browser tự làm)

---

## 🧪 PHẦN 7: Testing với Mock Data

Khi backend chưa sẵn sàng:

```typescript
// services/project-service.ts

export async function getProjects(classroomId: string) {
  // MOCK DATA - Comment khi backend ready
  await new Promise((resolve) => setTimeout(resolve, 500)); // Giả lập delay
  return {
    success: true,
    message: "OK",
    data: [
      { id: "1", title: "Project 1", description: "..." },
      { id: "2", title: "Project 2", description: "..." },
    ],
  };

  // REAL API - Uncomment khi backend ready
  // return apiFetch<ApiResponse<ProjectResponse[]>>(
  //   `/api/classrooms/${classroomId}/projects`
  // );
}
```

---

## 🆘 PHẦN 8: Troubleshooting

### Problem: Token hết hạn liên tục

**Nguyên nhân:** Access token expire sau 15 phút, refresh token expire sau 7 ngày.

**Giải pháp:**

- `apiFetch` tự động refresh khi 401
- Nếu refresh token cũng hết hạn → User phải login lại

### Problem: CORS error

```
Access to fetch at 'http://localhost:8080/api/...' from origin 'http://localhost:5173'
has been blocked by CORS policy
```

**Giải pháp:**

- Đảm bảo backend đã config CORS cho `http://localhost:5173`
- Check file `backend/internal/middleware/cors.go`

### Problem: 401 dù đã login

**Kiểm tra:**

1. Token có trong localStorage? `localStorage.getItem("access_token")`
2. Token format đúng? Kiểm tra Bearer token trong Network tab
3. Token có expire chưa? Decode JWT tại [jwt.io](https://jwt.io)

### Problem: Request bị "stuck"

**Nguyên nhân:** Backend không response hoặc CORS issue.

**Debug:**

1. Mở DevTools → Network tab
2. Check request có được gửi không
3. Check response status code
4. Check console có lỗi CORS không

---

## 📚 Tài liệu tham khảo

- **Svelte 5 Docs:** https://svelte.dev/docs/svelte/overview
- **TypeScript Handbook:** https://www.typescriptlang.org/docs/
- **Fetch API:** https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API
- **JWT:** https://jwt.io/introduction

---

## 🎯 Quick Reference

### Environment Variables

File: `.env` (tạo từ `.env.example`)

```bash
VITE_API_BASE_URL=http://localhost:8080
```

### Useful Commands

```bash
# Dev server
pnpm dev

# Build production
pnpm build

# Preview production build
pnpm preview

# Type check
pnpm check
```

---

## 👨‍💻 Need Help?

1. **Check existing code:** Xem các service/component đã có để tham khảo pattern
2. **Read error message:** Đọc kỹ error trong console
3. **Check Network tab:** Xem request/response thực tế
4. **Ask team:** Liên hệ team lead nếu còn vướng mắc

---

**Last updated:** January 8, 2026
