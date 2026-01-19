# 🗺️ Classroom Creation & Setup Flow (Backend Guide for FE)

Tài liệu này tổng hợp quy trình tạo và thiết lập Lớp học (Classroom) từ phía Backend để Frontend nắm bắt và tích hợp.

---

## 🚀 Giai đoạn 1: Tạo Lớp (The Core)

Đây là bước khởi tạo "khung sườn" cho một lớp học mới.

### Endpoint
`POST /api/classrooms`

### Authentication
*   **Header:** `Authorization: Bearer <token>`
*   **Role:** User phải là **Lecturer** (Giảng viên).

### Request Body (`CreateClassroomRequest`)
```json
{
  "name": "Nhập môn Công nghệ Phần mềm",
  "description": "Lớp học kỳ 1 năm 2024",
  "avatar": "https://example.com/avatar.jpg",
  "semester": 1,
  "year": 2024,
  "max_students": 50,
  "auto_approve": true,
  "require_email_domain": "student.hcmute.edu.vn"
}
```

*   `max_students`: Mặc định là 100 nếu không gửi.
*   `auto_approve`: `true` (sinh viên vào là duyệt ngay), `false` (cần giảng viên duyệt).

### Backend Logic
1.  Tạo lớp với trạng thái mặc định là `active`.
2.  **Tự động sinh Invitation Code** (Mã mời).
3.  **Tạo sẵn GeneralChannelID**: ID của kênh chat chung cho lớp.
4.  Gán User hiện tại làm chủ sở hữu (Lecturer).

### Response Success (201 Created)
Trả về object `Classroom` đầy đủ. **FE cần lưu lại `id` của lớp để dùng cho các bước sau.**

---

## 🛡️ Giai đoạn 2: Thiết lập Whitelist (Bảo mật sinh viên)

Sau khi tạo lớp, giảng viên có thể import danh sách mã số sinh viên (MSSV) được phép tham gia.

### 1. Tải file mẫu (Excel Template)
Giảng viên cần file mẫu để điền MSSV trước khi upload.
*   **Endpoint:** `GET /api/classrooms/whitelist-template`
*   **Response:** File `.xlsx` (binary).

### 2. Upload danh sách bằng Excel
*   **Endpoint:** `POST /api/classrooms/:id/whitelist-student-code/upload`
*   **Body (Multipart/Form-data):** `file: <file_excel.xlsx>`

### 3. Upload danh sách bằng JSON (Optional)
Nếu UI có chức năng copy-paste hoặc nhập tay danh sách MSSV.
*   **Endpoint:** `POST /api/classrooms/:id/whitelist-student-code`
*   **Body:**
    ```json
    {
      "student_codes": ["20110398", "20110400", "20110401"]
    }
    ```

### 4. Quản lý Whitelist (Xem/Sửa/Xóa)
*   **Xem:** `GET /api/classrooms/:id/whitelist-student-code`
*   **Thêm/Bớt:** `PATCH /api/classrooms/:id/whitelist-student-code`
    ```json
    {
      "add_codes": ["20119999"],
      "remove_codes": ["20110398"]
    }
    ```
*   **Xóa sạch:** `DELETE /api/classrooms/:id/whitelist-student-code`

---

## ⚙️ Giai đoạn 3: Vận hành & Quản trị

Các thao tác cấu hình sau khi lớp đã hoạt động.

### 1. Đổi Mã Mời (Regenerate Invitation Code)
Dùng khi mã cũ bị lộ ra ngoài.
*   **Endpoint:** `POST /api/classrooms/:id/regenerate-code`
*   **Response:** Trả về `new_code` mới. FE cần cập nhật ngay lên UI.

### 2. Cập nhật thông tin lớp (Settings)
*   **Endpoint:** `PUT /api/classrooms/:id`
*   **Body:** Gửi các trường cần sửa (name, description, settings...).

### 3. Đóng/Mở lớp (Status)
*   **Endpoint:** `PATCH /api/classrooms/:id/status`
*   **Body:**
    ```json
    { "status": "archived" } // hoặc "active"
    ```

### 4. Quản lý thành viên (Kích sinh viên)
*   **Endpoint:** `DELETE /api/classrooms/:id/students/:student_id`

---

## 💡 Notes for Frontend Devs

1.  **Invitation Code:** Luôn hiển thị mã này ở trang chi tiết lớp (cho giảng viên) để họ gửi cho sinh viên.
2.  **General Channel:** Khi user vào màn hình lớp học, có thể dùng `general_channel_id` trong response để connect vào khung chat chung ngay lập tức.
3.  **Error Handling:** Chú ý handle lỗi `401 Unauthorized` (hết hạn token) và `403 Forbidden` (sinh viên cố tình gọi API của giảng viên).
