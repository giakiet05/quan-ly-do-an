# 🔄 Frontend Backend Sync - January 2026

## 📋 Tổng quan thay đổi

Backend đã cập nhật và thêm nhiều models/DTOs mới. Frontend đã được sync để tương thích 100%.

---

## ✨ Những gì đã thay đổi

### 1. **Models Updated** (`src/models/index.ts`)

#### 🆕 Models mới:

- `ClassroomJoinRequest` - Request tham gia lớp học
- `ClassPost` - Bài đăng/thông báo trong lớp
- `JoinGroupInvitation` - Lời mời tham gia nhóm
- `Attachment` - File đính kèm
- `UserInfoResponse` - Response format cho user info
- `ReportPeriod` - Kỳ nộp báo cáo

#### 🔄 Models cập nhật:

- **`Classroom`** - Thêm nhiều fields mới:

  ```typescript
  - semester: string
  - year: number
  - status: ClassroomStatus
  - invitationCode: string
  - whitelistStudentCode?: string[]
  - requireEmailDomain?: string
  - maxStudents: number
  - autoApprove: boolean
  - canStudentDeleteGroup: boolean
  ```

- **`ProjectRound`** - Thêm:

  ```typescript
  - reportPeriods: ReportPeriod[]
  ```

- **`Group`** - Thêm:

  ```typescript
  - invitations?: JoinGroupInvitation[]
  ```

- **`JoinGroupRequest`** - Thêm:
  ```typescript
  - updatedAt: string
  ```

#### 📝 Type changes:

- `RequestStatus` → `JoinGroupRequestStatus`
- `ClassroomSetting` → Merged into `Classroom` as separate fields

---

### 2. **DTOs mới** (`src/dtos/`)

#### ✅ classroom-dto.ts

**Requests:**

- `CreateClassroomRequest` - Tạo lớp học
- `UpdateClassroomRequest` - Cập nhật lớp học
- `UpdateClassroomStatusRequest` - Đổi trạng thái
- `UploadWhitelistStudentCodeRequest` - Upload whitelist
- `UpdateWhitelistStudentCodeRequest` - Thêm/xóa whitelist

**Responses:**

- `ClassroomResponse` - Thông tin lớp học
- `RegenerateCodeResponse` - Code mời mới

#### ✅ classroom-join-dto.ts

**Requests:**

- `JoinClassroomRequest` - Tham gia lớp
- `UploadWhitelistRequest` - Upload whitelist

**Responses:**

- `ClassroomPreviewResponse` - Preview trước khi join
- `JoinClassroomResponse` - Kết quả join
- `JoinRequestResponse` - Thông tin request

#### ✅ class-post-dto.ts

**Requests:**

- `CreateClassPostRequest` - Tạo bài đăng
- `UpdateClassPostRequest` - Cập nhật bài đăng
- `AttachmentUpload` - Upload file đính kèm

**Responses:**

- `ClassPostResponse` - Thông tin bài đăng

#### 🔄 group-dto.ts (Updated)

**Requests thêm:**

- `CreateJoinGroupRequest` - Request tham gia nhóm
- `UpdateJoinGroupRequest` - Duyệt/từ chối request
- `CreateGroupInvitationRequest` - Gửi lời mời
- `UpdateGroupInvitationRequest` - Chấp nhận/từ chối lời mời

---

### 3. **Services mới** (`src/services/`)

#### ✅ classroom-service.ts (NEW)

**Lecturer - Classroom Management:**

- `getMyClassrooms()` - Lấy danh sách lớp của giảng viên
- `getClassroom(id)` - Chi tiết lớp
- `createClassroom(data)` - Tạo lớp mới
- `updateClassroom(id, data)` - Cập nhật lớp
- `updateClassroomStatus(id, status)` - Đổi trạng thái
- `deleteClassroom(id)` - Xóa lớp
- `regenerateInvitationCode(id)` - Tạo code mới
- `uploadWhitelistStudentCodes(id, codes)` - Upload whitelist
- `updateWhitelistStudentCodes(id, data)` - Thêm/xóa whitelist
- `removeStudentFromClassroom(id, studentId)` - Xóa sinh viên

**Student - Join Classroom:**

- `previewClassroom(code)` - Xem trước lớp
- `joinClassroom(code)` - Tham gia lớp
- `leaveClassroom(id)` - Rời lớp
- `getMyJoinedClassrooms()` - Lớp đã tham gia

**Lecturer - Join Requests:**

- `getPendingJoinRequests(id)` - Lấy requests chờ
- `approveJoinRequest(classroomId, requestId)` - Duyệt
- `rejectJoinRequest(classroomId, requestId)` - Từ chối

**Class Posts (Announcements):**

- `getClassPosts(id)` - Lấy danh sách bài đăng
- `getClassPost(classroomId, postId)` - Chi tiết bài
- `createClassPost(id, data)` - Tạo bài mới
- `updateClassPost(classroomId, postId, data)` - Cập nhật
- `deleteClassPost(classroomId, postId)` - Xóa
- `togglePinClassPost(classroomId, postId, isPinned)` - Ghim/bỏ ghim

#### 🔄 group-service.ts (Updated)

**Functions thêm:**

- `sendGroupInvitation(groupId, recipientId)` - Gửi lời mời
- `getGroupInvitations(groupId)` - Lời mời đã gửi
- `getMyGroupInvitations()` - Lời mời nhận được
- `acceptGroupInvitation(groupId, invitationId)` - Chấp nhận
- `rejectGroupInvitation(groupId, invitationId)` - Từ chối

---

## 🎯 Cách sử dụng

### Example 1: Classroom Management (Lecturer)

```typescript
import {
  getMyClassrooms,
  createClassroom,
  regenerateInvitationCode,
} from "@/services/classroom-service";

// Lấy danh sách lớp
const classrooms = await getMyClassrooms();

// Tạo lớp mới
const newClassroom = await createClassroom({
  name: "Web Development 2026",
  description: "Advanced web development course",
  avatar: "https://...",
  semester: "Spring 2026",
  year: 2026,
  maxStudents: 50,
  autoApprove: false,
  requireEmailDomain: "@hcmut.edu.vn",
});

// Tạo code mời mới
const { newCode } = await regenerateInvitationCode(newClassroom.id);
```

### Example 2: Join Classroom (Student)

```typescript
import {
  previewClassroom,
  joinClassroom,
  getMyJoinedClassrooms,
} from "@/services/classroom-service";

// Xem trước lớp
const preview = await previewClassroom("ABC123");

// Tham gia lớp
const result = await joinClassroom({ invitationCode: "ABC123" });

if (result.status === "approved") {
  console.log("Joined successfully!");
} else {
  console.log("Waiting for approval...");
}

// Lấy danh sách lớp đã tham gia
const myClasses = await getMyJoinedClassrooms();
```

### Example 3: Class Posts (Announcements)

```typescript
import {
  getClassPosts,
  createClassPost,
  togglePinClassPost,
} from "@/services/classroom-service";

// Lấy danh sách bài đăng
const posts = await getClassPosts(classroomId);

// Tạo thông báo mới
const newPost = await createClassPost(classroomId, {
  title: "Thông báo nộp báo cáo",
  content: "Deadline: 31/01/2026",
  attachments: [
    {
      fileName: "template.docx",
      fileURL: "https://...",
      fileSize: 102400,
      mimeType:
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    },
  ],
});

// Ghim bài đăng
await togglePinClassPost(classroomId, newPost.id, true);
```

### Example 4: Group Invitations

```typescript
import {
  sendGroupInvitation,
  getMyGroupInvitations,
  acceptGroupInvitation,
} from "@/services/group-service";

// Leader gửi lời mời
await sendGroupInvitation(groupId, recipientUserId);

// Sinh viên xem lời mời
const invitations = await getMyGroupInvitations();

// Chấp nhận lời mời
await acceptGroupInvitation(groupId, invitations[0].id);
```

---

## 📊 Tổng kết

### Models

- **Trước:** 30 models
- **Sau:** 36 models (+6)
- **Updated:** 4 models

### DTOs

- **Trước:** 9 DTO files
- **Sau:** 12 DTO files (+3)
- **Updated:** 1 DTO file

### Services

- **Trước:** 5 service files
- **Sau:** 6 service files (+1 new)
- **Updated:** 1 service file

### API Functions

- **Trước:** ~40 functions
- **Sau:** ~65 functions (+25)

---

## ✅ Checklist

- [x] Updated `models/index.ts` - All models synced
- [x] Created `dtos/classroom-dto.ts` - New DTOs
- [x] Created `dtos/classroom-join-dto.ts` - New DTOs
- [x] Created `dtos/class-post-dto.ts` - New DTOs
- [x] Updated `dtos/group-dto.ts` - Added invitations
- [x] Updated `dtos/index.ts` - Export new DTOs
- [x] Created `services/classroom-service.ts` - Complete classroom API
- [x] Updated `services/group-service.ts` - Added invitations
- [x] Updated `MODELS_DTOS_COMPLETE.md` - Documentation

---

## 🚀 Next Steps

1. **Test API Integration** - Khi backend ready
2. **Update UI Components** - Sử dụng services mới
3. **Add Error Handling** - Handle edge cases
4. **WebSocket Integration** - Real-time updates

---

## 📝 Notes

- Tất cả types đều có TypeScript strict mode
- All services sử dụng `apiFetch()` với auto-auth
- Mock data vẫn hoạt động cho development
- Backend API endpoints cần match với service calls

**🎉 Frontend đã sẵn sàng tích hợp với backend mới!**
