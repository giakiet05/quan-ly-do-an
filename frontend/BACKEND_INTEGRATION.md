# Backend Integration Guide

## Tổng quan

Toàn bộ code frontend đã được thiết kế để dễ dàng tích hợp với backend. Hiện tại đang dùng **mock data** và chỉ cần thay bằng **API calls**.

## Đã chuẩn bị

### 1. **Models/Types** - `src/models/index.ts`

- ✅ Tất cả models đã được port từ backend Go models
- ✅ TypeScript interfaces khớp với cấu trúc BE
- ✅ Enum types (ProjectStatus, MessageType, NotificationType, etc.)

### 2. **API Services** - `src/services/`

- ✅ `api.ts` - Base API client với authentication
- ✅ `project-service.ts` - Project & ProjectRound APIs
- ✅ `group-service.ts` - Group, Task, Report, Join Request APIs
- ✅ `message-service.ts` - Message/Chat APIs
- ✅ `notification-service.ts` - Notification APIs
- ✅ `classroom-service.ts` - Classroom management APIs

### 3. **Components với Props Types**

Tất cả components đều có TypeScript interfaces rõ ràng:

- `StudentProjectReportsTab` - nhận `projectName`, `teamName`
- `StudentProjectTeamChatTab` - nhận `team`, `currentStudentId`
- `StudentProjectMembersTab` - nhận `team`, `isLeader`
- `ReportDetailModal` - nhận `reportTitle`, `submittedFile`, `grade`, `feedback`

## Cách tích hợp Backend

### Bước 1: Cấu hình API Base URL

Trong file `.env`:

```
VITE_API_BASE_URL=http://localhost:8080
```

### Bước 2: Thay Mock Data bằng API Calls

**VÍ DỤ: Trang My Projects**

**Trước (Mock data):**

```typescript
// MyProjects.svelte
let projects = $state<Project[]>([
  {
    id: "1",
    name: "Hệ thống quản lý thư viện",
    // ... mock data
  },
]);
```

**Sau (API integration):**

```typescript
// MyProjects.svelte
import { getMyProjectGroup } from "../services/group-service";
import { onMount } from "svelte";

let projects = $state<Project[]>([]);
let loading = $state(true);
let error = $state<string | null>(null);

onMount(async () => {
  try {
    loading = true;
    // Gọi API thay vì dùng mock data
    const myGroups = await getClassroomGroups(classroomId);
    projects = myGroups.map((group) => ({
      // Map BE data to FE structure
      id: group.projectId,
      // ...
    }));
  } catch (err: any) {
    error = err.message;
  } finally {
    loading = false;
  }
});
```

### Bước 3: Thêm Loading & Error States

```typescript
{#if loading}
  <div class="loading">Đang tải...</div>
{:else if error}
  <div class="error">{error}</div>
{:else}
  <!-- Render nội dung -->
{/if}
```

### Bước 4: WebSocket cho Real-time Features

Cho chat và notifications:

```typescript
// src/services/websocket-service.ts
let ws: WebSocket | null = null;

export function connectWebSocket(token: string) {
  ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    // Handle incoming messages
  };
}
```

## Mapping: Pages → API Services

| Page/Component              | API Service               | Functions                               |
| --------------------------- | ------------------------- | --------------------------------------- |
| `MyProjects.svelte`         | `group-service.ts`        | `getClassroomGroups()`                  |
| `MyProjectDetail.svelte`    | `group-service.ts`        | `getGroup()`, `getMyProjectGroup()`     |
| `StudentProjectReportsTab`  | `group-service.ts`        | `getGroupReports()`, `submitReport()`   |
| `StudentProjectTeamChatTab` | `message-service.ts`      | `getChannelMessages()`, `sendMessage()` |
| `StudentProjectMembersTab`  | `group-service.ts`        | `getGroup()` (members included)         |
| `Chat.svelte`               | `message-service.ts`      | `getChannelMessages()`, `sendMessage()` |
| `Notification.svelte`       | `notification-service.ts` | `getNotifications()`, `markAsRead()`    |

## Example: Tích hợp Report Submission

```typescript
// StudentProjectReportsTab.svelte
import { submitReport } from "../services/group-service";

async function handleSubmitReport(file: File, reportTitle: string) {
  try {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("title", reportTitle);
    formData.append("content", "Report content");

    const newReport = await submitReport(groupId, formData);

    // Update UI with new report
    reports = [...reports, newReport];
  } catch (err: any) {
    alert(`Lỗi: ${err.message}`);
  }
}
```

## Example: Tích hợp Chat

```typescript
// StudentProjectTeamChatTab.svelte
import { getChannelMessages, sendMessage } from "../services/message-service";

let messages = $state<Message[]>([]);

onMount(async () => {
  messages = await getChannelMessages(channelId);
});

async function handleSendMessage(content: string) {
  const newMessage = await sendMessage(channelId, content);
  messages = [...messages, newMessage];
}
```

## WebSocket Events

Khi nhận message từ WebSocket:

```typescript
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  switch (data.type) {
    case "new_message":
      // Thêm message vào chat
      break;
    case "notification":
      // Hiển thị notification
      break;
    case "group_update":
      // Update group info
      break;
  }
};
```

## Testing

### 1. Test với Mock Backend

Dùng Mock Service Worker (MSW) hoặc json-server

### 2. Test với Real Backend

```bash
# Start BE
cd backend
go run main.go

# Start FE
cd frontend
pnpm dev
```

## Checklist Tích Hợp

- [ ] Cấu hình `VITE_API_BASE_URL` trong `.env`
- [ ] Test authentication flow (login/logout)
- [ ] Thay mock data bằng API calls ở từng page
- [ ] Thêm loading states
- [ ] Thêm error handling
- [ ] Implement WebSocket cho real-time
- [ ] Test file upload (reports)
- [ ] Test pagination (nếu có)
- [ ] Test permission checks (student vs lecturer)

## Notes

- Tất cả API functions đã có **type safety** với TypeScript
- Authentication tự động handle qua `apiFetch()` trong `api.ts`
- Token refresh tự động trong `auth-service.ts`
- Logout tự động khi 401 Unauthorized
- Tất cả dates từ BE (ISO strings) có thể dùng `new Date()`

## Liên hệ

Nếu có vấn đề khi tích hợp, check:

1. Network tab trong DevTools
2. Console logs
3. Backend logs
4. API response structure khớp với models chưa
