# 📚 Hướng dẫn sử dụng Classroom Updates cho Student

## 🎯 Các thay đổi đã được cập nhật

### 1. **classroom-dto.ts** ✅

- Thêm `WhitelistEntry` interface mới
- Cập nhật `ClassroomResponse` với:
  - `coLecturers?` - Danh sách trợ giảng
  - `allowedEmailDomains?` - Mảng domain email được phép
  - `enableWhitelist?` - Có bật whitelist MSSV không
  - `enableEmailRestriction?` - Có bật giới hạn domain email không
  - `whitelistStudentCode?` - Đổi từ `string[]` sang `WhitelistEntry[]`

### 2. **classroom-join-dto.ts** ✅

- Cập nhật `ClassroomPreviewResponse` thêm:
  - `enableWhitelist?`
  - `enableEmailRestriction?`
  - `allowedEmailDomains?`

### 3. **StudentClassList.svelte** ✅

- Xóa `requiresProjectRegistration` (không còn tồn tại)
- Sửa `project_rounds` → `projectRounds`

### 4. **classroom-error-codes.ts** ✅ (MỚI)

- Error codes cho join classroom
- Error messages tiếng Việt
- Helper function `getClassroomErrorMessage()`

### 5. **ClassroomRequirements.svelte** ✅ (MỚI)

- Component hiển thị yêu cầu tham gia lớp
- Hiển thị whitelist, email restriction, số lượng SV, trợ giảng

---

## 💡 Cách sử dụng

### 1. Hiển thị thông tin lớp học (có trợ giảng)

```typescript
import type { ClassroomResponse } from "../dtos/classroom-dto";

const classroom = await getClassroom(classroomId);

// Hiển thị trợ giảng
if (classroom.coLecturers && classroom.coLecturers.length > 0) {
  classroom.coLecturers.forEach((coLecturer) => {
    console.log(`Trợ giảng: ${coLecturer.fullName}`);
  });
}
```

### 2. Sử dụng ClassroomRequirements component

```svelte
<script lang="ts">
  import ClassroomRequirements from "../components/ClassroomRequirements.svelte";
  import type { ClassroomResponse } from "../dtos/classroom-dto";

  let classroom = $state<ClassroomResponse | null>(null);
</script>

{#if classroom}
  <ClassroomRequirements {classroom} />
{/if}
```

Component này sẽ tự động hiển thị:

- ✅ Thông báo nếu có whitelist MSSV
- ✅ Domain email được phép
- ✅ Số lượng sinh viên / tối đa
- ✅ Danh sách trợ giảng

### 3. Handle error khi join classroom

```typescript
import { joinClassroom } from "../services/classroom-service";
import {
  getClassroomErrorMessage,
  CLASSROOM_ERROR_CODES,
} from "../constants/classroom-error-codes";

try {
  const result = await joinClassroom({ invitationCode: "ABC123" });

  if (result.status === "pending") {
    alert("Yêu cầu tham gia đã được gửi, chờ giảng viên duyệt");
  } else {
    alert("Tham gia lớp học thành công!");
  }
} catch (error: any) {
  const errorCode = error.error_code || error.message;

  // Hiển thị message tiếng Việt tùy theo error code
  const message = getClassroomErrorMessage(errorCode);
  alert(message);

  // Hoặc handle từng error cụ thể:
  switch (errorCode) {
    case CLASSROOM_ERROR_CODES.STUDENT_CODE_ALREADY_USED:
      // MSSV đã được claim bởi user khác
      alert("Mã sinh viên này đã được sử dụng bởi người khác");
      break;

    case CLASSROOM_ERROR_CODES.INVALID_EMAIL_DOMAIN:
      // Email không đúng domain
      alert("Email của bạn không thuộc domain được phép");
      break;

    case CLASSROOM_ERROR_CODES.STUDENT_CODE_NOT_IN_WHITELIST:
      // MSSV không trong whitelist
      alert("MSSV của bạn không có trong danh sách");
      break;

    case CLASSROOM_ERROR_CODES.CLASSROOM_FULL:
      alert("Lớp đã đầy");
      break;
  }
}
```

### 4. Kiểm tra whitelist entry

```typescript
import type { WhitelistEntry } from "../dtos/classroom-dto";

const classroom = await getClassroom(classroomId);

if (classroom.whitelistStudentCode) {
  classroom.whitelistStudentCode.forEach((entry: WhitelistEntry) => {
    console.log(`MSSV: ${entry.studentCode}`);

    if (entry.joinedBy) {
      console.log(`Đã được claim bởi: ${entry.joinedBy}`);
      console.log(`Thời gian: ${entry.joinedAt}`);
    } else {
      console.log("Chưa có ai claim");
    }
  });
}
```

---

## 📝 Error Codes mới cần handle

| Error Code                      | Mô tả                      | Khi nào xảy ra                                       |
| ------------------------------- | -------------------------- | ---------------------------------------------------- |
| `STUDENT_CODE_ALREADY_USED`     | MSSV đã được claim         | Khi join với MSSV đã có user khác dùng               |
| `INVALID_EMAIL_DOMAIN`          | Email không đúng domain    | Khi email không thuộc allowedEmailDomains            |
| `STUDENT_CODE_NOT_IN_WHITELIST` | MSSV không trong whitelist | Khi enableWhitelist=true và MSSV không có trong list |
| `CLASSROOM_FULL`                | Lớp đã đầy                 | Khi số SV >= maxStudents                             |
| `CODE_INVALID`                  | Mã mời không hợp lệ        | Mã mời sai hoặc hết hạn                              |
| `ALREADY_IN_CLASSROOM`          | Đã là thành viên           | User đã join lớp rồi                                 |

---

## 🚀 Các file đã được cập nhật

1. ✅ `frontend/src/dtos/classroom-dto.ts`
2. ✅ `frontend/src/dtos/classroom-join-dto.ts`
3. ✅ `frontend/src/pages/student/StudentClassList.svelte`
4. ✅ `frontend/src/constants/classroom-error-codes.ts` (mới)
5. ✅ `frontend/src/components/ClassroomRequirements.svelte` (mới)

---

## ✨ Next Steps (Optional)

### Nếu có UI join classroom:

1. Import error constants
2. Thêm try-catch với error handling
3. Hiển thị message tiếng Việt cho user

### Nếu có preview classroom page:

1. Import `ClassroomRequirements` component
2. Hiển thị yêu cầu tham gia trước khi join

### Nếu có classroom detail page:

1. Hiển thị danh sách trợ giảng (nếu có)
2. Sử dụng `ClassroomRequirements` để show requirements

---

## ⚠️ Breaking Changes

- `whitelistStudentCode` giờ là `WhitelistEntry[]` thay vì `string[]`
- `requireEmailDomain` (string) đã bị xóa, thay bằng `allowedEmailDomains` (string[])
- `requiresProjectRegistration` không còn tồn tại trong DTO

Kiểm tra code nếu có sử dụng các field này!
