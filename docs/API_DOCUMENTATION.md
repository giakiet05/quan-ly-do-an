# API Documentation

Base URL: `/api`

Tất cả endpoints trả về response theo format:

```json
{
  "success": true|false,
  "message": "Thông báo thành công/lỗi",
  "data": {...},           // Chỉ có khi success: true
  "error_code": "CODE"     // Chỉ có khi success: false
}
```

**Ký hiệu:**
-  = Yêu cầu authentication (gửi `Authorization: Bearer <access_token>` trong header)
-  = Chỉ giảng viên mới có quyền truy cập
-  = Giảng viên hoặc trợ giảng có quyền truy cập
-  = Cả sinh viên và giảng viên đều có quyền truy cập

---

## Authentication

### Local Authentication Flow

#### 1. Gửi mã xác minh email
**Endpoint:** `POST /api/auth/local/send-verification`

**Mô tả:** Gửi mã OTP 6 số đến email để xác minh quyền sở hữu.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Mã xác minh đã được gửi thành công",
  "data": null
}
```

**Lỗi:**
- `400 BAD_REQUEST` - Email không hợp lệ
- `409 EMAIL_EXISTS` - Email đã được đăng ký

---

#### 2. Xác minh mã OTP
**Endpoint:** `POST /api/auth/local/verify-email`

**Mô tả:** Xác minh mã OTP đã gửi tới email và trả về verification token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Email đã được xác minh thành công",
  "data": {
    "verification_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Lỗi:**
- `400 INVALID_OTP` - Mã OTP không đúng
- `400 OTP_EXPIRED` - Mã OTP đã hết hạn (15 phút)

---

#### 3. Hoàn tất đăng ký
**Endpoint:** `POST /api/auth/local/complete-registration`

**Mô tả:** Hoàn tất đăng ký tài khoản với thông tin bổ sung.

**Request Body:**
```json
{
  "verification_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "full_name": "Nguyễn Văn A",
  "student_code": "2021600001",  // Để trống nếu là giảng viên
  "password": "securePassword123"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Đăng ký tài khoản thành công",
  "data": {
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "email": "user@example.com",
      "full_name": "Nguyễn Văn A",
      "student_code": "2021600001",
      "provider": "local",
      "avatar": null,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Lỗi:**
- `400 INVALID_TOKEN` - Token xác minh không hợp lệ
- `400 BAD_REQUEST` - Dữ liệu không hợp lệ (tên quá ngắn, mật khẩu yếu)

---

#### 4. Đăng nhập
**Endpoint:** `POST /api/auth/local/login`

**Mô tả:** Đăng nhập bằng email/username và mật khẩu.

**Request Body:**
```json
{
  "identifier": "user@example.com",  // Email hoặc username
  "password": "securePassword123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Đăng nhập thành công",
  "data": {
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "email": "user@example.com",
      "full_name": "Nguyễn Văn A",
      "student_code": "2021600001",
      "provider": "local",
      "avatar": null,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Lỗi:**
- `401 INVALID_CREDENTIALS` - Email/mật khẩu không đúng
- `403 EMAIL_NOT_VERIFIED` - Email chưa được xác minh
- `403 USER_INACTIVE` - Tài khoản đã bị vô hiệu hóa

---

#### 5. Gửi lại mã OTP
**Endpoint:** `POST /api/auth/local/resend-otp`

**Mô tả:** Gửi lại mã OTP xác minh đến email.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Mã xác minh đã được gửi lại thành công",
  "data": null
}
```

---

### Quên mật khẩu

#### 1. Yêu cầu đặt lại mật khẩu
**Endpoint:** `POST /api/auth/local/forgot-password`

**Mô tả:** Gửi mã OTP đặt lại mật khẩu đến email người dùng.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Mã đặt lại mật khẩu đã được gửi thành công",
  "data": null
}
```

**Lỗi:**
- `404 EMAIL_NOT_REGISTERED` - Email chưa được đăng ký
- `409 LOGIN_METHOD_MISMATCH` - Tài khoản đăng nhập bằng Google

---

#### 2. Xác minh OTP đặt lại mật khẩu
**Endpoint:** `POST /api/auth/local/verify-reset-otp`

**Mô tả:** Xác minh mã OTP và trả về reset token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "OTP đã được xác minh thành công",
  "data": {
    "reset_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Lỗi:**
- `400 INVALID_OTP` - Mã OTP không đúng
- `400 OTP_EXPIRED` - Mã OTP đã hết hạn

---

#### 3. Đặt lại mật khẩu
**Endpoint:** `POST /api/auth/local/reset-password`

**Mô tả:** Đặt lại mật khẩu bằng reset token.

**Request Body:**
```json
{
  "reset_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_password": "newSecurePassword123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Mật khẩu đã được đặt lại thành công",
  "data": null
}
```

**Lỗi:**
- `400 INVALID_TOKEN` - Reset token không hợp lệ
- `400 BAD_REQUEST` - Mật khẩu quá yếu

---

### Google OAuth Flow

#### 1. Khởi tạo đăng nhập Google
**Endpoint:** `GET /api/auth/google/login`

**Mô tả:** Chuyển hướng đến trang đồng ý OAuth của Google.

**Response:** 307 Redirect đến Google OAuth URL

---

#### 2. Google Callback
**Endpoint:** `GET /api/auth/google/callback`

**Mô tả:** Google chuyển hướng về đây sau khi xác thực. Xử lý 2 trường hợp:
- **Người dùng mới:** Trả về `setup_token` để hoàn tất thông tin
- **Người dùng cũ:** Trả về `access_token` và `refresh_token`

**Query Parameters:**
- `code` - Mã authorization từ Google

**Response:** Trang HTML chuyển hướng với tokens trong URL fragment

---

#### 3. Hoàn tất thiết lập Google
**Endpoint:** `POST /api/auth/google/complete-setup`

**Mô tả:** Hoàn tất thông tin cho người dùng Google mới.

**Request Body:**
```json
{
  "setup_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "full_name": "Nguyễn Văn A",
  "student_code": "2021600001"  // Để trống nếu là giảng viên
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Thiết lập tài khoản thành công",
  "data": {
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "email": "user@gmail.com",
      "full_name": "Nguyễn Văn A",
      "student_code": "2021600001",
      "provider": "google",
      "avatar": {
        "url": "https://lh3.googleusercontent.com/...",
        "public_id": null
      },
      "created_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### Quản lý Token

#### Làm mới Access Token
**Endpoint:** `POST /api/auth/refresh`

**Mô tả:** Tạo access token mới bằng refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Token đã được làm mới thành công",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Lỗi:**
- `401 INVALID_TOKEN` - Refresh token không hợp lệ hoặc đã hết hạn

---

#### Đăng xuất
**Endpoint:** `POST /api/auth/logout`

**Mô tả:** Vô hiệu hóa cả access token và refresh token.

**Request Body:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Đăng xuất thành công",
  "data": null
}
```

---

## Notes

### Authentication Headers
For protected endpoints (marked with ), include the access token:
```
Authorization: Bearer <access_token>
```

### Token Expiration
- **Access Token:** 15 minutes
- **Refresh Token:** 7 days
- **OTP Codes:** 15 minutes

### User Types
- **Student:** Has `student_code` field populated
- **Lecturer:** `student_code` is `null`

---

---

## Users

### Lấy thông tin cá nhân 
**Endpoint:** `GET /api/users/me`

**Mô tả:** Lấy thông tin profile của user đang đăng nhập.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy thông tin thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "full_name": "Nguyễn Văn A",
    "student_code": "2021600001",
    "provider": "local",
    "avatar": {
      "url": "https://res.cloudinary.com/...",
      "public_id": "avatars/xyz123"
    },
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### Cập nhật thông tin cá nhân 
**Endpoint:** `PUT /api/users/me`

**Mô tả:** Cập nhật thông tin profile (tên, mã sinh viên).

**Request Body:**
```json
{
  "full_name": "Nguyễn Văn B",
  "student_code": "2021600002"
}
```

**Lưu ý:** Tất cả các field đều optional, chỉ gửi field muốn cập nhật.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật thông tin thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "full_name": "Nguyễn Văn B",
    "student_code": "2021600002",
    "provider": "local",
    "avatar": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### Đổi mật khẩu 
**Endpoint:** `PUT /api/users/me/password`

**Mô tả:** Đổi mật khẩu (chỉ cho tài khoản local, không dùng cho Google OAuth).

**Request Body:**
```json
{
  "old_password": "oldPassword123",
  "new_password": "newPassword456"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Đổi mật khẩu thành công",
  "data": null
}
```

**Lỗi:**
- `401 INVALID_CREDENTIALS` - Mật khẩu cũ không đúng
- `409 LOGIN_METHOD_MISMATCH` - Tài khoản đăng nhập bằng Google

---

### Upload avatar 
**Endpoint:** `POST /api/users/me/avatar`

**Mô tả:** Upload ảnh đại diện lên Cloudinary.

**Request:** `multipart/form-data`
- `file`: File ảnh (jpg, png, gif, ...)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Upload avatar thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "full_name": "Nguyễn Văn A",
    "student_code": "2021600001",
    "provider": "local",
    "avatar": {
      "url": "https://res.cloudinary.com/.../avatar.jpg",
      "public_id": "avatars/xyz123"
    },
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### Xóa avatar 
**Endpoint:** `DELETE /api/users/me/avatar`

**Mô tả:** Xóa ảnh đại diện khỏi Cloudinary và đặt về null.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa avatar thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "full_name": "Nguyễn Văn A",
    "student_code": "2021600001",
    "provider": "local",
    "avatar": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Classrooms

### Tạo lớp học (Chỉ giảng viên)
**Endpoint:** `POST /api/classrooms`

**Mô tả:** Tạo lớp học mới (chỉ giảng viên).

**Request Body:**
```json
{
  "name": "Lập trình Web",
  "description": "Môn lập trình web cơ bản",
  "avatar": "https://res.cloudinary.com/.../class-avatar.jpg",
  "semester": "Spring",
  "year": 2024,
  "max_students": 50,
  "auto_approve": false,
  "allowed_email_domains": ["@hcmut.edu.vn", "@student.hcmut.edu.vn"],
  "enable_email_restriction": true,
  "enable_whitelist": false
}
```

**Lưu ý:**
- `max_students` mặc định: 100
- `auto_approve` mặc định: false (yêu cầu phê duyệt thủ công)
- `allowed_email_domains`: Mảng các domain email được phép (VD: ["@hcmut.edu.vn"])
- `enable_email_restriction`: Bật/tắt giới hạn domain email (mặc định: false)
- `enable_whitelist`: Bật/tắt whitelist mã sinh viên (mặc định: false)

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo lớp học thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Lập trình Web",
    "description": "Môn lập trình web cơ bản",
    "avatar": "https://res.cloudinary.com/.../class-avatar.jpg",
    "semester": "Spring",
    "year": 2024,
    "status": "active",
    "lecturer": {
      "user_id": "507f1f77bcf86cd799439012",
      "full_name": "GV. Nguyễn Văn A",
      "avatar": null
    },
    "co_lecturers": [],
    "students": [],
    "invitation_code": "ABC123XYZ",
    "max_students": 50,
    "auto_approve": false,
    "allowed_email_domains": ["@hcmut.edu.vn", "@student.hcmut.edu.vn"],
    "enable_email_restriction": true,
    "enable_whitelist": false,
    "created_at": "2024-01-01T00:00:00Z",
    "whitelist_student_code": []
  }
}
```

---

### Lấy danh sách lớp do mình tạo
**Endpoint:** `GET /api/classrooms/my`

**Mô tả:** Lấy danh sách các lớp học mà user là giảng viên.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách lớp học thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Lập trình Web",
      "description": "Môn lập trình web cơ bản",
      "semester": "Spring",
      "year": 2024,
      "status": "active",
      "lecturer": {
        "user_id": "507f1f77bcf86cd799439012",
        "full_name": "GV. Nguyễn Văn A",
        "avatar": null
      },
      "students": [],
      "invitation_code": "ABC123XYZ",
      "max_students": 50,
      "auto_approve": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### Lấy danh sách lớp đã tham gia
**Endpoint:** `GET /api/classrooms/joined`

**Mô tả:** Lấy danh sách các lớp học mà user là học sinh.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách lớp học thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Lập trình Web",
      "description": "Môn lập trình web cơ bản",
      "semester": "Spring",
      "year": 2024,
      "status": "active",
      "lecturer": {
        "user_id": "507f1f77bcf86cd799439012",
        "full_name": "GV. Nguyễn Văn A",
        "avatar": null
      },
      "students": [...],
      "max_students": 50,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### Lấy chi tiết lớp học
**Endpoint:** `GET /api/classrooms/:id`

**Mô tả:** Lấy thông tin chi tiết của một lớp học.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy thông tin lớp học thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Lập trình Web",
    "description": "Môn lập trình web cơ bản",
    "avatar": "https://res.cloudinary.com/.../class-avatar.jpg",
    "semester": "Spring",
    "year": 2024,
    "status": "active",
    "lecturer": {
      "user_id": "507f1f77bcf86cd799439012",
      "full_name": "GV. Nguyễn Văn A",
      "avatar": null
    },
    "co_lecturers": [
      {
        "user_id": "507f1f77bcf86cd799439015",
        "full_name": "Trợ giảng Nguyễn Văn C",
        "avatar": null
      }
    ],
    "students": [
      {
        "user_id": "507f1f77bcf86cd799439013",
        "full_name": "Nguyễn Văn B",
        "avatar": null,
        "student_code": "2021600001"
      }
    ],
    "invitation_code": "ABC123XYZ",
    "max_students": 50,
    "auto_approve": false,
    "allowed_email_domains": ["@hcmut.edu.vn"],
    "enable_email_restriction": true,
    "enable_whitelist": true,
    "created_at": "2024-01-01T00:00:00Z",
    "whitelist_student_code": [
      {
        "student_code": "2021600001",
        "joined_by": "507f1f77bcf86cd799439013",
        "joined_at": "2024-01-05T10:00:00Z"
      },
      {
        "student_code": "2021600002",
        "joined_by": null,
        "joined_at": null
      }
    ]
  }
}
```

**Lưu ý:**
- `whitelist_student_code`: Mảng các object chứa mã sinh viên và trạng thái join
  - `joined_by`: ID của user đã claim mã này (null nếu chưa có ai)
  - `joined_at`: Thời điểm user claim mã (null nếu chưa có ai)

**Lỗi:**
- `404 CLASSROOM_NOT_FOUND` - Lớp học không tồn tại
- `403 FORBIDDEN` - Không có quyền truy cập

---

### Cập nhật lớp học (Giảng viên hoặc trợ giảng)
**Endpoint:** `PUT /api/classrooms/:id`

**Mô tả:** Cập nhật thông tin lớp học (giảng viên hoặc trợ giảng).

**Request Body:**
```json
{
  "name": "Lập trình Web Nâng cao",
  "description": "Mô tả mới",
  "max_students": 60,
  "auto_approve": true,
  "allowed_email_domains": ["@hcmut.edu.vn", "@student.hcmut.edu.vn"],
  "enable_email_restriction": true,
  "enable_whitelist": true,
  "can_student_delete_group": false
}
```

**Lưu ý:** Tất cả các field đều optional.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật lớp học thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Lập trình Web Nâng cao",
    ...
  }
}
```

**Lỗi:**
- `403 FORBIDDEN` - Không phải giảng viên hoặc trợ giảng của lớp
- `404 CLASSROOM_NOT_FOUND` - Lớp học không tồn tại

---

### Xóa lớp học (Chỉ giảng viên)
**Endpoint:** `DELETE /api/classrooms/:id`

**Mô tả:** Xóa lớp học (chỉ giảng viên).

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa lớp học thành công",
  "data": {
    "classroom_id": "507f1f77bcf86cd799439011"
  }
}
```

---

### Cập nhật trạng thái lớp học (Giảng viên hoặc trợ giảng)
**Endpoint:** `PATCH /api/classrooms/:id/status`

**Mô tả:** Cập nhật trạng thái lớp học (active, inactive, archived).

**Request Body:**
```json
{
  "status": "archived"
}
```

**Giá trị hợp lệ:** `active`, `inactive`, `archived`

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật trạng thái thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "status": "archived"
  }
}
```

---

### Xóa sinh viên khỏi lớp (Giảng viên hoặc trợ giảng)
**Endpoint:** `DELETE /api/classrooms/:id/students/:student_id`

**Mô tả:** Xóa một sinh viên khỏi lớp học.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa sinh viên khỏi lớp thành công",
  "data": null
}
```

---

### Rời lớp học
**Endpoint:** `POST /api/classrooms/:id/leave`

**Mô tả:** Sinh viên hoặc trợ giảng tự rời khỏi lớp học. Giảng viên chính không thể rời lớp.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Rời lớp học thành công",
  "data": null
}
```

**Lỗi:**
- `403 FORBIDDEN` - Giảng viên chính không thể rời lớp / Không phải thành viên lớp
- `404 CLASSROOM_NOT_FOUND` - Lớp học không tồn tại

---

### Xóa trợ giảng khỏi lớp (Chỉ giảng viên chính)
**Endpoint:** `DELETE /api/classrooms/:id/co-lecturers/:co_lecturer_id`

**Mô tả:** Giảng viên chính xóa trợ giảng khỏi lớp học.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa trợ giảng thành công",
  "data": null
}
```

**Lỗi:**
- `403 FORBIDDEN` - Không phải giảng viên chính
- `404 CLASSROOM_NOT_FOUND` - Lớp học không tồn tại

---

### Tải template whitelist
**Endpoint:** `GET /api/classrooms/whitelist-template`

**Mô tả:** Tải file Excel template cho danh sách mã sinh viên được phép.

**Response:** File Excel (.xlsx)

---

### Lấy whitelist mã sinh viên (Giảng viên hoặc trợ giảng)
**Endpoint:** `GET /api/classrooms/:id/whitelist-student-code`

**Mô tả:** Lấy danh sách mã sinh viên được phép tham gia lớp, kèm trạng thái đã join.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy whitelist thành công",
  "data": {
    "whitelist": [
      {
        "student_code": "2021600001",
        "joined_by": "507f1f77bcf86cd799439013",
        "joined_at": "2024-01-05T10:00:00Z"
      },
      {
        "student_code": "2021600002",
        "joined_by": null,
        "joined_at": null
      },
      {
        "student_code": "2021600003",
        "joined_by": null,
        "joined_at": null
      }
    ]
  }
}
```

**Lưu ý:**
- `joined_by`: ID user đã claim mã sinh viên này (null nếu chưa ai claim)
- `joined_at`: Thời điểm claim (null nếu chưa ai claim)

---

### Upload whitelist (JSON) (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/:id/whitelist-student-code`

**Mô tả:** Upload danh sách mã sinh viên bằng JSON (thay thế toàn bộ whitelist cũ).

**Request Body:**
```json
{
  "student_codes": ["2021600001", "2021600002", "2021600003"]
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Upload whitelist thành công",
  "data": {
    "student_codes": ["2021600001", "2021600002", "2021600003"]
  }
}
```

---

### Upload whitelist (Excel) (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/:id/whitelist-student-code/upload`

**Mô tả:** Upload danh sách mã sinh viên bằng file Excel.

**Request:** `multipart/form-data`
- `file`: File Excel (.xlsx)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Upload whitelist thành công",
  "data": {
    "student_codes": ["2021600001", "2021600002", "2021600003"]
  }
}
```

---

### Cập nhật whitelist (Giảng viên hoặc trợ giảng)
**Endpoint:** `PATCH /api/classrooms/:id/whitelist-student-code`

**Mô tả:** Thêm hoặc xóa mã sinh viên khỏi whitelist.

**Request Body:**
```json
{
  "add_codes": ["2021600004", "2021600005"],
  "remove_codes": ["2021600001"]
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật whitelist thành công",
  "data": {
    "student_codes": ["2021600002", "2021600003", "2021600004", "2021600005"]
  }
}
```

---

### Xóa whitelist (Giảng viên hoặc trợ giảng)
**Endpoint:** `DELETE /api/classrooms/:id/whitelist-student-code`

**Mô tả:** Xóa toàn bộ whitelist mã sinh viên.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa whitelist thành công",
  "data": null
}
```

---

### Tạo lại mã mời (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/:id/regenerate-code`

**Mô tả:** Tạo lại invitation code mới cho lớp học.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tạo mã mời mới thành công",
  "data": {
    "new_code": "XYZ789ABC",
    "code_expires_at": "2024-12-31T23:59:59Z"
  }
}
```

---

## Classroom Join

### Xem trước lớp học (Không cần đăng nhập)
**Endpoint:** `GET /api/classrooms/preview?code=ABC123XYZ`

**Mô tả:** Xem thông tin cơ bản của lớp học trước khi tham gia (không cần auth).

**Query Parameters:**
- `code`: Mã mời lớp học

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy thông tin lớp học thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Lập trình Web",
    "lecturer": "GV. Nguyễn Văn A",
    "student_count": 25,
    "max_students": 50
  }
}
```

**Lỗi:**
- `400 CODE_INVALID` - Mã mời không hợp lệ hoặc đã hết hạn

---

### Tham gia lớp học
**Endpoint:** `POST /api/classrooms/join`

**Mô tả:** Gửi yêu cầu tham gia lớp học bằng mã mời.

**Request Body:**
```json
{
  "invitation_code": "ABC123XYZ"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Tham gia lớp học thành công",
  "data": {
    "classroom_id": "507f1f77bcf86cd799439011",
    "status": "approved",
    "message": "Bạn đã tham gia lớp học thành công"
  }
}
```

**Lưu ý:**
- Nếu `auto_approve = true`: Tự động duyệt, `status = "approved"`
- Nếu `auto_approve = false`: Chờ phê duyệt, `status = "pending"`

**Lỗi:**
- `400 CODE_INVALID` - Mã mời không hợp lệ
- `409 ALREADY_IN_CLASSROOM` - Đã là thành viên của lớp
- `409 STUDENT_CODE_NOT_IN_WHITELIST` - Mã SV không trong whitelist
- `400 INVALID_EMAIL_DOMAIN` - Email không đúng domain yêu cầu
- `409 CLASSROOM_FULL` - Lớp học đã đầy

---

### Lấy danh sách yêu cầu tham gia (Giảng viên hoặc trợ giảng)
**Endpoint:** `GET /api/classrooms/:id/join-requests`

**Mô tả:** Lấy danh sách các yêu cầu tham gia đang chờ duyệt.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách yêu cầu thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439014",
      "student_code": "2021600001",
      "full_name": "Nguyễn Văn B",
      "email": "student@hcmut.edu.vn",
      "status": "pending",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### Phê duyệt yêu cầu tham gia (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/join-requests/:id/approve`

**Mô tả:** Chấp nhận yêu cầu tham gia lớp.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Phê duyệt yêu cầu thành công",
  "data": null
}
```

**Lỗi:**
- `404 JOIN_REQUEST_NOT_FOUND` - Yêu cầu không tồn tại
- `409 JOIN_REQUEST_ALREADY_PROCESSED` - Yêu cầu đã được xử lý

---

### Từ chối yêu cầu tham gia (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/join-requests/:id/reject`

**Mô tả:** Từ chối yêu cầu tham gia lớp.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Từ chối yêu cầu thành công",
  "data": null
}
```

---

## Classroom Invitations (Mời trợ giảng)

### Mời trợ giảng vào lớp (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/:id/invitations`

**Mô tả:** Mời người khác làm trợ giảng cho lớp học thông qua email.

**Request Body:**
```json
{
  "email": "ta@hcmut.edu.vn"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Gửi lời mời thành công",
  "data": {
    "id": "507f1f77bcf86cd799439030",
    "classroom_id": "507f1f77bcf86cd799439011",
    "inviter_id": "507f1f77bcf86cd799439012",
    "invitee_id": "507f1f77bcf86cd799439016",
    "status": "pending",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

**Lỗi:**
- `400 BAD_REQUEST` - Email không hợp lệ
- `404 USER_NOT_FOUND` - Không tìm thấy user với email này
- `409 ALREADY_LECTURER` - User đã là giảng viên chính của lớp
- `409 ALREADY_CO_LECTURER` - User đã là trợ giảng của lớp
- `409 INVITATION_ALREADY_EXISTS` - Đã có lời mời pending cho user này
- `403 FORBIDDEN` - Không thể mời chính mình

---

### Lấy danh sách lời mời của lớp (Giảng viên hoặc trợ giảng)
**Endpoint:** `GET /api/classrooms/:id/invitations?page=1&pageSize=20`

**Mô tả:** Lấy danh sách các lời mời trợ giảng của lớp.

**Query Parameters:**
- `page`: Số trang (mặc định: 1)
- `pageSize`: Số lời mời/trang (mặc định: 20)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách lời mời thành công",
  "data": {
    "invitations": [
      {
        "id": "507f1f77bcf86cd799439030",
        "classroom": {
          "id": "507f1f77bcf86cd799439011",
          "name": "Lập trình Web"
        },
        "inviter": {
          "id": "507f1f77bcf86cd799439012",
          "full_name": "GV. Nguyễn Văn A",
          "email": "lecturer@hcmut.edu.vn"
        },
        "invitee": {
          "id": "507f1f77bcf86cd799439016",
          "full_name": "Trợ giảng B",
          "email": "ta@hcmut.edu.vn"
        },
        "status": "pending",
        "created_at": "2024-01-01T10:00:00Z",
        "responded_at": null
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 1,
      "page_size": 20
    }
  }
}
```

---

### Lấy danh sách lời mời của tôi
**Endpoint:** `GET /api/classroom-invitations/my?page=1&pageSize=20`

**Mô tả:** Lấy danh sách các lời mời trợ giảng mà mình nhận được.

**Query Parameters:**
- `page`: Số trang (mặc định: 1)
- `pageSize`: Số lời mời/trang (mặc định: 20)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách lời mời thành công",
  "data": {
    "invitations": [
      {
        "id": "507f1f77bcf86cd799439030",
        "classroom": {
          "id": "507f1f77bcf86cd799439011",
          "name": "Lập trình Web"
        },
        "inviter": {
          "id": "507f1f77bcf86cd799439012",
          "full_name": "GV. Nguyễn Văn A",
          "email": "lecturer@hcmut.edu.vn"
        },
        "invitee": {
          "id": "507f1f77bcf86cd799439016",
          "full_name": "Trợ giảng B",
          "email": "ta@hcmut.edu.vn"
        },
        "status": "pending",
        "created_at": "2024-01-01T10:00:00Z",
        "responded_at": null
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 1,
      "page_size": 20
    }
  }
}
```

---

### Chấp nhận lời mời trợ giảng
**Endpoint:** `POST /api/classroom-invitations/:id/accept`

**Mô tả:** Chấp nhận lời mời làm trợ giảng.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Chấp nhận lời mời thành công",
  "data": null
}
```

**Lỗi:**
- `404 INVITATION_NOT_FOUND` - Lời mời không tồn tại
- `403 FORBIDDEN` - Không phải người được mời
- `409 INVITATION_ALREADY_PROCESSED` - Lời mời đã được xử lý

---

### Từ chối lời mời trợ giảng
**Endpoint:** `POST /api/classroom-invitations/:id/reject`

**Mô tả:** Từ chối lời mời làm trợ giảng.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Từ chối lời mời thành công",
  "data": null
}
```

**Lỗi:**
- `404 INVITATION_NOT_FOUND` - Lời mời không tồn tại
- `403 FORBIDDEN` - Không phải người được mời
- `409 INVITATION_ALREADY_PROCESSED` - Lời mời đã được xử lý

---

### Hủy lời mời trợ giảng (Người mời)
**Endpoint:** `DELETE /api/classroom-invitations/:id`

**Mô tả:** Hủy lời mời đã gửi (chỉ người mời hoặc giảng viên/trợ giảng của lớp).

**Success Response (200):**
```json
{
  "success": true,
  "message": "Hủy lời mời thành công",
  "data": null
}
```

**Lỗi:**
- `404 INVITATION_NOT_FOUND` - Lời mời không tồn tại
- `403 FORBIDDEN` - Không có quyền hủy lời mời
- `409 INVITATION_ALREADY_PROCESSED` - Lời mời đã được xử lý

---

## Class Posts

### Tạo bài đăng (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/classrooms/:id/posts`

**Mô tả:** Tạo bài đăng mới trong lớp học.

**Workflow:**
1. Upload files trước: `POST /api/classrooms/posts/upload-attachments`
2. Lấy URLs từ bước 1
3. Tạo post với URLs trong `attachments`

**Request Body:**
```json
{
  "title": "Bài tập tuần 1",
  "content": "Nội dung bài tập...",
  "attachments": [
    {
      "file_name": "bt-tuan-1.pdf",
      "file_url": "https://res.cloudinary.com/.../bt-tuan-1.pdf",
      "file_size": 2048000,
      "mime_type": "application/pdf"
    }
  ]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo bài đăng thành công",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "classroom_id": "507f1f77bcf86cd799439011",
    "author": {
      "user_id": "507f1f77bcf86cd799439012",
      "full_name": "GV. Nguyễn Văn A",
      "avatar": null
    },
    "title": "Bài tập tuần 1",
    "content": "Nội dung bài tập...",
    "attachments": [...],
    "is_pinned": false,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### Lấy danh sách bài đăng
**Endpoint:** `GET /api/classrooms/:id/posts?page=1&pageSize=20`

**Mô tả:** Lấy danh sách bài đăng trong lớp (có phân trang).

**Query Parameters:**
- `page`: Số trang (mặc định: 1)
- `pageSize`: Số bài/trang (mặc định: 20)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách bài đăng thành công",
  "data": {
    "posts": [...],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

**Lưu ý:** Bài đăng được sắp xếp: Bài ghim trước, sau đó theo thời gian tạo giảm dần.

---

### Lấy chi tiết bài đăng
**Endpoint:** `GET /api/classrooms/posts/:post_id`

**Mô tả:** Lấy thông tin chi tiết một bài đăng.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy bài đăng thành công",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "classroom_id": "507f1f77bcf86cd799439011",
    "author": {...},
    "title": "Bài tập tuần 1",
    "content": "Nội dung bài tập...",
    "attachments": [...],
    "is_pinned": false,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### Cập nhật bài đăng
**Endpoint:** `PUT /api/classrooms/posts/:post_id`

**Mô tả:** Cập nhật bài đăng (giảng viên, trợ giảng hoặc tác giả).

**Request Body:**
```json
{
  "title": "Bài tập tuần 1 (Cập nhật)",
  "content": "Nội dung mới...",
  "attachments_to_add": [
    {
      "file_name": "file-moi.pdf",
      "file_url": "https://res.cloudinary.com/.../file-moi.pdf",
      "file_size": 1024000,
      "mime_type": "application/pdf"
    }
  ],
  "attachments_to_remove": [
    "https://res.cloudinary.com/.../file-cu.pdf"
  ]
}
```

**Lưu ý:**
- Tất cả field đều optional
- `attachments_to_add`: Mảng file mới cần thêm
- `attachments_to_remove`: Mảng URLs file cần xóa

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật bài đăng thành công",
  "data": {...}
}
```

---

### Xóa bài đăng (Giảng viên hoặc trợ giảng)
**Endpoint:** `DELETE /api/classrooms/posts/:post_id`

**Mô tả:** Xóa bài đăng.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa bài đăng thành công",
  "data": {
    "post_id": "507f1f77bcf86cd799439015"
  }
}
```

---

### Ghim/bỏ ghim bài đăng (Giảng viên hoặc trợ giảng)
**Endpoint:** `PATCH /api/classrooms/posts/:post_id/pin`

**Mô tả:** Ghim hoặc bỏ ghim bài đăng.

**Request Body:**
```json
{
  "is_pinned": true
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật trạng thái ghim thành công",
  "data": {
    "is_pinned": true
  }
}
```

---

### Upload attachments
**Endpoint:** `POST /api/classrooms/posts/upload-attachments`

**Mô tả:** Upload nhiều file lên Cloudinary trước khi tạo/cập nhật bài đăng.

**Request:** `multipart/form-data`
- `files`: Mảng các file cần upload

**Success Response (200):**
```json
{
  "success": true,
  "message": "Upload files thành công",
  "data": {
    "attachments": [
      {
        "file_name": "document.pdf",
        "file_url": "https://res.cloudinary.com/.../document.pdf",
        "file_size": 2048000,
        "mime_type": "application/pdf",
        "public_id": "class_posts/xyz123"
      }
    ]
  }
}
```

---

### Xóa attachment
**Endpoint:** `DELETE /api/classrooms/posts/attachments/:public_id`

**Mô tả:** Xóa một file khỏi Cloudinary bằng public_id.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa file thành công",
  "data": {
    "public_id": "class_posts/xyz123",
    "result": "ok"
  }
}
```

---

### Thêm attachment vào post
**Endpoint:** `POST /api/classrooms/posts/:post_id/attachments`

**Mô tả:** Thêm một attachment vào bài đăng đã tồn tại.

**Request Body:**
```json
{
  "file_name": "new-file.pdf",
  "file_url": "https://res.cloudinary.com/.../new-file.pdf",
  "file_size": 1024000,
  "mime_type": "application/pdf"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Thêm attachment thành công",
  "data": {...}
}
```

---

### Xóa attachment khỏi post
**Endpoint:** `DELETE /api/classrooms/posts/:post_id/attachments`

**Mô tả:** Xóa một attachment khỏi bài đăng.

**Request Body:**
```json
{
  "file_url": "https://res.cloudinary.com/.../file-to-remove.pdf"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa attachment thành công",
  "data": {...}
}
```

---

## Groups

### Tạo nhóm
**Endpoint:** `POST /api/groups`

**Mô tả:** Tạo nhóm đồ án mới.

**Request Body:**
```json
{
  "classroom_id": "507f1f77bcf86cd799439011",
  "project_id": "507f1f77bcf86cd799439020",
  "project_round_id": "507f1f77bcf86cd799439021",
  "member_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo nhóm thành công",
  "data": {
    "id": "507f1f77bcf86cd799439022",
    "classroom_id": "507f1f77bcf86cd799439011",
    "project_id": "507f1f77bcf86cd799439020",
    "group_channel_id": "507f1f77bcf86cd799439023",
    "leader_id": "507f1f77bcf86cd799439013",
    "members": [...],
    "tasks": [],
    "reports": [],
    "setting": {...}
  }
}
```

---

### Lấy thông tin nhóm
**Endpoint:** `GET /api/groups/:group_id`

**Mô tả:** Lấy chi tiết thông tin một nhóm.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy thông tin nhóm thành công",
  "data": {
    "id": "507f1f77bcf86cd799439022",
    "classroom_id": "507f1f77bcf86cd799439011",
    "project_id": "507f1f77bcf86cd799439020",
    "group_channel_id": "507f1f77bcf86cd799439023",
    "leader_id": "507f1f77bcf86cd799439013",
    "members": [
      {
        "user_id": "507f1f77bcf86cd799439013",
        "full_name": "Nguyễn Văn A",
        "avatar": null
      }
    ],
    "tasks": [],
    "task_statuses": ["todo", "in_progress", "done"],
    "reports": [],
    "setting": {
      "can_members_add_task": true,
      "can_members_edit_task": true
    }
  }
}
```

---

### Lấy danh sách nhóm (Filter)
**Endpoint:** `GET /api/groups?classroom_id=xxx&project_id=yyy&member_id=zzz`

**Mô tả:** Lấy danh sách nhóm theo filter.

**Query Parameters:**
- `classroom_id` (required): ID lớp học
- `project_id` (optional): ID đồ án
- `member_id` (optional): ID thành viên

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách nhóm thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439022",
      "classroom_id": "507f1f77bcf86cd799439011",
      "project_id": "507f1f77bcf86cd799439020",
      "leader_id": "507f1f77bcf86cd799439013",
      "members": [...]
    }
  ]
}
```

---

### Cập nhật nhóm
**Endpoint:** `PUT /api/groups`

**Mô tả:** Cập nhật thông tin nhóm (project, leader, settings).

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "project_id": "507f1f77bcf86cd799439024",
  "leader_id": "507f1f77bcf86cd799439014",
  "setting": {
    "can_members_add_task": false,
    "can_members_edit_task": true
  }
}
```

**Lưu ý:** Tất cả field đều optional trừ `group_id`.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật nhóm thành công",
  "data": {...}
}
```

---

### Xóa nhóm
**Endpoint:** `DELETE /api/groups/:group_id`

**Mô tả:** Xóa nhóm (leader hoặc giảng viên).

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa nhóm thành công",
  "data": null
}
```

---

## Group Join Requests & Invitations

### Gửi yêu cầu tham gia nhóm
**Endpoint:** `POST /api/groups/join-requests` 

**Mô tả:** Sinh viên gửi yêu cầu tham gia nhóm.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "message": "Mình muốn tham gia nhóm này"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Gửi yêu cầu thành công",
  "data": {
    "id": "507f1f77bcf86cd799439030",
    "user_id": "507f1f77bcf86cd799439015",
    "status": "pending",
    "message": "Mình muốn tham gia nhóm này",
    "requested_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Lỗi:**
- `403 FORBIDDEN` - Nhóm không cho phép yêu cầu tham gia (AllowJoinRequest = false)
- `400 BAD_REQUEST` - Người dùng đã là thành viên hoặc đã có yêu cầu pending
- `400 BAD_REQUEST` - Nhóm đã đạt số lượng thành viên tối đa

**Notification:**
- **Người nhận:** Nhóm trưởng (leader)
- **Nội dung:** "Bạn có yêu cầu tham gia nhóm mới"

---

### Chấp nhận yêu cầu tham gia nhóm
**Endpoint:** `PUT /api/groups/join-requests/accept` 

**Mô tả:** Nhóm trưởng chấp nhận yêu cầu tham gia.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "request_id": "507f1f77bcf86cd799439030"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Chấp nhận yêu cầu thành công",
  "data": null
}
```

**Lỗi:**
- `403 FORBIDDEN` - Chỉ nhóm trưởng mới có quyền chấp nhận
- `400 BAD_REQUEST` - Yêu cầu không ở trạng thái pending
- `400 BAD_REQUEST` - Nhóm đã đạt số lượng thành viên tối đa

**Notification:**
- **Người nhận:** Người gửi yêu cầu
- **Nội dung:** "Yêu cầu tham gia nhóm của bạn đã được chấp nhận"

---

### Từ chối yêu cầu tham gia nhóm
**Endpoint:** `PUT /api/groups/join-requests/reject` 

**Mô tả:** Nhóm trưởng từ chối yêu cầu tham gia.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "request_id": "507f1f77bcf86cd799439030"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Từ chối yêu cầu thành công",
  "data": null
}
```

**Lỗi:**
- `403 FORBIDDEN` - Chỉ nhóm trưởng mới có quyền từ chối
- `400 BAD_REQUEST` - Yêu cầu không ở trạng thái pending

**Notification:**
- **Người nhận:** Người gửi yêu cầu
- **Nội dung:** "Yêu cầu tham gia nhóm của bạn đã bị từ chối"

---

### Mời sinh viên vào nhóm
**Endpoint:** `POST /api/groups/invitations` 

**Mô tả:** Nhóm trưởng mời sinh viên tham gia nhóm.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "recipient_id": "507f1f77bcf86cd799439016"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Gửi lời mời thành công",
  "data": {
    "id": "507f1f77bcf86cd799439031",
    "group_id": "507f1f77bcf86cd799439022",
    "recipient_id": "507f1f77bcf86cd799439016",
    "status": "pending",
    "sent_at": "2024-01-15T11:00:00Z"
  }
}
```

**Lỗi:**
- `403 FORBIDDEN` - Chỉ nhóm trưởng mới có quyền mời
- `400 BAD_REQUEST` - Người được mời đã là thành viên hoặc đã có lời mời pending
- `400 BAD_REQUEST` - Nhóm đã đạt số lượng thành viên tối đa

**Notification:**
- **Người nhận:** Người được mời (recipient)
- **Nội dung:** "Bạn có thư mời tham gia nhóm"

---

### Chấp nhận lời mời tham gia nhóm
**Endpoint:** `PUT /api/groups/invitations/accept` 

**Mô tả:** Sinh viên chấp nhận lời mời tham gia nhóm.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "invitation_id": "507f1f77bcf86cd799439031"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Chấp nhận lời mời thành công",
  "data": null
}
```

**Lỗi:**
- `403 FORBIDDEN` - Chỉ người được mời mới có quyền chấp nhận
- `400 BAD_REQUEST` - Lời mời không ở trạng thái pending
- `400 BAD_REQUEST` - Nhóm đã đạt số lượng thành viên tối đa

**Notification:**
- **Người nhận:** Nhóm trưởng (inviter/leader)
- **Nội dung:** "Lời mời tham gia nhóm đã được chấp nhận"

---

### Từ chối lời mời tham gia nhóm
**Endpoint:** `PUT /api/groups/invitations/reject` 

**Mô tả:** Sinh viên từ chối lời mời tham gia nhóm.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "invitation_id": "507f1f77bcf86cd799439031"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Từ chối lời mời thành công",
  "data": null
}
```

**Lỗi:**
- `403 FORBIDDEN` - Chỉ người được mời mới có quyền từ chối
- `400 BAD_REQUEST` - Lời mời không ở trạng thái pending

**Notification:**
- **Người nhận:** Nhóm trưởng (inviter/leader)
- **Nội dung:** "Lời mời tham gia nhóm đã bị từ chối"

---

## Group Tasks

### Tạo task
**Endpoint:** `POST /api/groups/tasks`

**Mô tả:** Tạo task mới cho nhóm.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "title": "Thiết kế UI",
  "details": "Thiết kế giao diện màn hình đăng nhập",
  "assign_to_ids": ["507f1f77bcf86cd799439013"],
  "due_date": "2024-12-31T23:59:59Z",
  "status": "todo"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo task thành công",
  "data": {
    "id": "507f1f77bcf86cd799439025",
    "title": "Thiết kế UI",
    "assign_to_ids": ["507f1f77bcf86cd799439013"],
    "due_date": "2024-12-31T23:59:59Z",
    "status": "todo"
  }
}
```

---

### Cập nhật task
**Endpoint:** `PUT /api/groups/tasks`

**Mô tả:** Cập nhật thông tin task.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "task_id": "507f1f77bcf86cd799439025",
  "title": "Thiết kế UI (Updated)",
  "status": "in_progress"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật task thành công",
  "data": {...}
}
```

---

### Xóa task
**Endpoint:** `DELETE /api/groups/:group_id/tasks/:task_id`

**Mô tả:** Xóa task khỏi nhóm.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa task thành công",
  "data": null
}
```

---

### Tạo báo cáo
**Endpoint:** `POST /api/groups/reports`

**Mô tả:** Tạo báo cáo tiến độ cho nhóm.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "title": "Báo cáo tuần 1",
  "content": "Nội dung báo cáo...",
  "files": [
    {
      "name": "report.pdf",
      "url": "https://res.cloudinary.com/.../report.pdf",
      "size": 1024000
    }
  ]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo báo cáo thành công",
  "data": {
    "id": "507f1f77bcf86cd799439026",
    "title": "Báo cáo tuần 1",
    "content": "Nội dung báo cáo...",
    "files": [...],
    "feedback": null
  }
}
```

---

### Cập nhật báo cáo
**Endpoint:** `PUT /api/groups/reports`

**Mô tả:** Cập nhật nội dung báo cáo.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "report_id": "507f1f77bcf86cd799439026",
  "title": "Báo cáo tuần 1 (Cập nhật)",
  "content": "Nội dung mới..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật báo cáo thành công",
  "data": {...}
}
```

---

### Xóa báo cáo
**Endpoint:** `DELETE /api/groups/:group_id/reports/:report_id`

**Mô tả:** Xóa báo cáo.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa báo cáo thành công",
  "data": null
}
```

---

### Tạo feedback cho báo cáo (Giảng viên hoặc trợ giảng)
**Endpoint:** `POST /api/groups/reports/feedback`

**Mô tả:** Giảng viên hoặc trợ giảng đưa feedback và điểm cho báo cáo.

**Request Body:**
```json
{
  "group_id": "507f1f77bcf86cd799439022",
  "report_id": "507f1f77bcf86cd799439026",
  "content": "Báo cáo tốt, cần bổ sung phần...",
  "grade": "8.5"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo feedback thành công",
  "data": {
    "content": "Báo cáo tốt, cần bổ sung phần...",
    "grade": "8.5",
    "lecturer_id": "507f1f77bcf86cd799439012",
    "commented_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### Cập nhật feedback
**Endpoint:** `PUT /api/groups/:group_id/reports/:report_id/feedback`

**Mô tả:** Cập nhật feedback đã tạo.

**Request Body:**
```json
{
  "content": "Feedback mới...",
  "grade": "9.0"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật feedback thành công",
  "data": {...}
}
```

---

### Xóa feedback
**Endpoint:** `DELETE /api/groups/:group_id/reports/:report_id/feedback`

**Mô tả:** Xóa feedback.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa feedback thành công",
  "data": null
}
```

---

## Channels

### Tạo channel
**Endpoint:** `POST /api/channels`

**Mô tả:** Tạo channel chat mới giữa các user.

**Request Body:**
```json
{
  "members": [
    {
      "id": "507f1f77bcf86cd799439013",
      "full_name": "Nguyễn Văn A",
      "avatar": null
    },
    {
      "id": "507f1f77bcf86cd799439014",
      "full_name": "Nguyễn Văn B",
      "avatar": null
    }
  ]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo channel thành công",
  "data": {
    "id": "507f1f77bcf86cd799439027",
    "members": [...],
    "settings": [...],
    "background": null,
    "status": "active",
    "unread_message_count": 0,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### Lấy thông tin channel
**Endpoint:** `GET /api/channels/:channel_id`

**Mô tả:** Lấy chi tiết một channel.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy thông tin channel thành công",
  "data": {
    "id": "507f1f77bcf86cd799439027",
    "members": [...],
    "settings": [
      {
        "user_id": "507f1f77bcf86cd799439013",
        "nickname": "User A",
        "notification": true,
        "typing_indicator": true
      }
    ],
    "background": null,
    "status": "active",
    "unread_message_count": 5,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### Lấy danh sách channel của user
**Endpoint:** `GET /api/channels/user?user_id=xxx&page=1&pageSize=20`

**Mô tả:** Lấy danh sách channel mà user tham gia.

**Query Parameters:**
- `user_id`: ID của user
- `page`: Số trang (mặc định: 1)
- `pageSize`: Số channel/trang (mặc định: 20)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách channel thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439027",
      "members": [...],
      "unread_message_count": 5,
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### Lấy channel giữa 2 user
**Endpoint:** `GET /api/channels/between/:user1/:user2`

**Mô tả:** Tìm channel giữa 2 user cụ thể.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy channel thành công",
  "data": {
    "id": "507f1f77bcf86cd799439027",
    "members": [...],
    "settings": [...]
  }
}
```

**Lỗi:**
- `404 CHANNEL_NOT_FOUND` - Không tìm thấy channel

---

### Cập nhật channel
**Endpoint:** `PUT /api/channels`

**Mô tả:** Cập nhật settings của channel (nickname, background, notification, ...).

**Request Body:**
```json
{
  "channel_id": "507f1f77bcf86cd799439027",
  "nickname": "Nhóm dự án",
  "background": "https://res.cloudinary.com/.../bg.jpg",
  "notification": false,
  "typing_indicator": true
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật channel thành công",
  "data": {...}
}
```

---

### Xóa channel
**Endpoint:** `DELETE /api/channels/:channel_id`

**Mô tả:** Xóa channel.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa channel thành công",
  "data": null
}
```

---

## Messages

### Lấy thông tin tin nhắn
**Endpoint:** `GET /api/messages/:message_id`

**Mô tả:** Lấy chi tiết một tin nhắn.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy tin nhắn thành công",
  "data": {
    "id": "507f1f77bcf86cd799439028",
    "channel_id": "507f1f77bcf86cd799439027",
    "sender_id": "507f1f77bcf86cd799439013",
    "sender_username": "user_a",
    "type": "text",
    "content": "Hello world!",
    "read_by": ["507f1f77bcf86cd799439013"],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### Lấy danh sách tin nhắn (Filter)
**Endpoint:** `GET /api/messages/filter?channel_id=xxx&page=1&pageSize=50`

**Mô tả:** Lấy danh sách tin nhắn theo filter.

**Query Parameters:**
- `channel_id` (required): ID channel
- `sender_id` (optional): ID người gửi
- `search_content` (optional): Tìm kiếm nội dung
- `is_read` (optional): Lọc tin đã đọc
- `is_send` (optional): Lọc tin đã gửi
- `is_media` (optional): Lọc tin có media
- `page`: Số trang (mặc định: 1)
- `pageSize`: Số tin/trang (mặc định: 50)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách tin nhắn thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439028",
      "channel_id": "507f1f77bcf86cd799439027",
      "sender_id": "507f1f77bcf86cd799439013",
      "sender_username": "user_a",
      "type": "text",
      "content": "Hello world!",
      "read_by": ["507f1f77bcf86cd799439013"],
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### Xóa tin nhắn
**Endpoint:** `DELETE /api/messages/:message_id`

**Mô tả:** Xóa một tin nhắn.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa tin nhắn thành công",
  "data": null
}
```

**Lưu ý:** Gửi tin nhắn thực hiện qua WebSocket, không qua REST API.

---

## Notifications

### Lấy danh sách thông báo
**Endpoint:** `GET /api/notifications?page=1&pageSize=20`

**Mô tả:** Lấy danh sách thông báo của user.

**Query Parameters:**
- `page`: Số trang (mặc định: 1)
- `pageSize`: Số thông báo/trang (mặc định: 20)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách thông báo thành công",
  "data": {
    "notifications": [
      {
        "id": "507f1f77bcf86cd799439029",
        "type": "classroom_invitation",
        "message": "Bạn được mời vào lớp Lập trình Web",
        "link": "/classrooms/507f1f77bcf86cd799439011",
        "is_read": false,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 100,
      "page_size": 20
    }
  }
}
```

---

### Đánh dấu tất cả đã đọc
**Endpoint:** `PUT /api/notifications/read-all`

**Mô tả:** Đánh dấu tất cả thông báo là đã đọc.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Đánh dấu tất cả thông báo đã đọc thành công",
  "data": null
}
```

---

## WebSocket

### Kết nối WebSocket
**Endpoint:** `ws://localhost:8080/api/ws?token=<access_token>`

**Mô tả:** Kết nối WebSocket để nhận/gửi tin nhắn real-time.

**Query Parameters:**
- `token`: Access token để xác thực

### Gửi tin nhắn
**Message Format (Client → Server):**
```json
{
  "type": "send_message",
  "data": {
    "channel_id": "507f1f77bcf86cd799439027",
    "content": "Hello!",
    "message_type": "text"
  }
}
```

**Message Types:**
- `text`: Tin nhắn văn bản
- `image`: Hình ảnh
- `file`: File đính kèm
- `voice`: Tin nhắn voice

### Nhận tin nhắn
**Message Format (Server → Client):**
```json
{
  "type": "new_message",
  "data": {
    "id": "507f1f77bcf86cd799439028",
    "channel_id": "507f1f77bcf86cd799439027",
    "sender_id": "507f1f77bcf86cd799439013",
    "sender_username": "user_a",
    "type": "text",
    "content": "Hello!",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### Typing indicator
**Client → Server:**
```json
{
  "type": "typing",
  "data": {
    "channel_id": "507f1f77bcf86cd799439027",
    "is_typing": true
  }
}
```

**Server → Client:**
```json
{
  "type": "user_typing",
  "data": {
    "channel_id": "507f1f77bcf86cd799439027",
    "user_id": "507f1f77bcf86cd799439013",
    "username": "user_a",
    "is_typing": true
  }
}
```

### Mark message as read
**Client → Server:**
```json
{
  "type": "mark_read",
  "data": {
    "message_id": "507f1f77bcf86cd799439028"
  }
}
```

---

## Projects

### Tạo report period
**Endpoint:** `POST /api/projects/classrooms/:classroom_id/rounds/:round_id/report-periods`

**Mô tả:** Tạo một report period mới cho project round.

**Request Body:**
```json
{
  "name": "Báo cáo tuần 1",
  "description": "Mô tả báo cáo",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-07T23:59:59Z"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo report period thành công",
  "data": {
    "id": "507f1f77bcf86cd799439030",
    "name": "Báo cáo tuần 1",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-07T23:59:59Z"
  }
}
```

---

### Tạo nhiều report periods
**Endpoint:** `POST /api/projects/classrooms/:classroom_id/report-periods/bulk`

**Mô tả:** Tạo nhiều report periods một lúc.

**Request Body:**
```json
{
  "round_id": "507f1f77bcf86cd799439021",
  "report_periods": [
    {
      "name": "Báo cáo tuần 1",
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-01-07T23:59:59Z"
    },
    {
      "name": "Báo cáo tuần 2",
      "start_date": "2024-01-08T00:00:00Z",
      "end_date": "2024-01-14T23:59:59Z"
    }
  ]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo report periods thành công",
  "data": {
    "created_count": 2,
    "report_periods": [...]
  }
}
```

---

### Lấy thông tin report period
**Endpoint:** `GET /api/projects/classrooms/:classroom_id/rounds/:round_id/report-periods/:period_id`

**Mô tả:** Lấy chi tiết một report period.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy report period thành công",
  "data": {
    "id": "507f1f77bcf86cd799439030",
    "name": "Báo cáo tuần 1",
    "description": "Mô tả báo cáo",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-07T23:59:59Z"
  }
}
```

---

### Lấy danh sách report periods
**Endpoint:** `GET /api/projects/classrooms/:classroom_id/rounds/:round_id/report-periods`

**Mô tả:** Lấy danh sách report periods theo round.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách report periods thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439030",
      "name": "Báo cáo tuần 1",
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-01-07T23:59:59Z"
    }
  ]
}
```

---

### Cập nhật report period
**Endpoint:** `PUT /api/projects/report-periods`

**Mô tả:** Cập nhật thông tin report period.

**Request Body:**
```json
{
  "period_id": "507f1f77bcf86cd799439030",
  "name": "Báo cáo tuần 1 (Updated)",
  "description": "Mô tả mới",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-10T23:59:59Z"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật report period thành công",
  "data": {...}
}
```

---

### Xóa report period
**Endpoint:** `DELETE /api/projects/classrooms/:classroom_id/rounds/:round_id/report-periods/:period_id`

**Mô tả:** Xóa report period.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa report period thành công",
  "data": null
}
```

---

### Tạo project
**Endpoint:** `POST /api/projects`

**Mô tả:** Tạo project mới.

**Request Body:**
```json
{
  "classroom_id": "507f1f77bcf86cd799439011",
  "round_id": "507f1f77bcf86cd799439021",
  "name": "Hệ thống quản lý thư viện",
  "description": "Mô tả đồ án",
  "requirements": "Yêu cầu chức năng..."
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo project thành công",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "name": "Hệ thống quản lý thư viện",
    "description": "Mô tả đồ án",
    "requirements": "Yêu cầu chức năng..."
  }
}
```

---

### Tạo nhiều projects
**Endpoint:** `POST /api/projects/bulk`

**Mô tả:** Tạo nhiều projects một lúc.

**Request Body:**
```json
{
  "classroom_id": "507f1f77bcf86cd799439011",
  "round_id": "507f1f77bcf86cd799439021",
  "projects": [
    {
      "name": "Project 1",
      "description": "Mô tả 1"
    },
    {
      "name": "Project 2",
      "description": "Mô tả 2"
    }
  ]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo projects thành công",
  "data": {
    "created_count": 2,
    "projects": [...]
  }
}
```

---

### Lấy thông tin project
**Endpoint:** `GET /api/projects/classrooms/:classroom_id/:project_id`

**Mô tả:** Lấy chi tiết một project.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy project thành công",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "classroom_id": "507f1f77bcf86cd799439011",
    "round_id": "507f1f77bcf86cd799439021",
    "name": "Hệ thống quản lý thư viện",
    "description": "Mô tả đồ án",
    "requirements": "Yêu cầu chức năng..."
  }
}
```

---

### Lấy danh sách projects
**Endpoint:** `GET /api/projects/classrooms/:classroom_id/rounds/:round_id/projects`

**Mô tả:** Lấy danh sách projects theo round.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách projects thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "name": "Hệ thống quản lý thư viện",
      "description": "Mô tả đồ án"
    }
  ]
}
```

---

### Cập nhật project
**Endpoint:** `PUT /api/projects`

**Mô tả:** Cập nhật thông tin project.

**Request Body:**
```json
{
  "project_id": "507f1f77bcf86cd799439020",
  "name": "Hệ thống quản lý thư viện (Updated)",
  "description": "Mô tả mới",
  "requirements": "Yêu cầu mới..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật project thành công",
  "data": {...}
}
```

---

### Xóa project
**Endpoint:** `DELETE /api/projects/classrooms/:classroom_id/:project_id`

**Mô tả:** Xóa project.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa project thành công",
  "data": null
}
```

---

### Tạo project round
**Endpoint:** `POST /api/projects/rounds`

**Mô tả:** Tạo project round mới cho lớp học.

**Request Body:**
```json
{
  "classroom_id": "507f1f77bcf86cd799439011",
  "name": "Đợt 1",
  "description": "Đợt đồ án học kỳ 1",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-05-31T23:59:59Z",
  "max_members_per_group": 5
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo project round thành công",
  "data": {
    "id": "507f1f77bcf86cd799439021",
    "classroom_id": "507f1f77bcf86cd799439011",
    "name": "Đợt 1",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-05-31T23:59:59Z",
    "max_members_per_group": 5
  }
}
```

---

### Tạo nhiều project rounds
**Endpoint:** `POST /api/projects/rounds/bulk`

**Mô tả:** Tạo nhiều project rounds một lúc.

**Request Body:**
```json
{
  "classroom_id": "507f1f77bcf86cd799439011",
  "rounds": [
    {
      "name": "Đợt 1",
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-05-31T23:59:59Z",
      "max_members_per_group": 5
    },
    {
      "name": "Đợt 2",
      "start_date": "2024-06-01T00:00:00Z",
      "end_date": "2024-10-31T23:59:59Z",
      "max_members_per_group": 5
    }
  ]
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Tạo project rounds thành công",
  "data": {
    "created_count": 2,
    "rounds": [...]
  }
}
```

---

### Lấy thông tin project round
**Endpoint:** `GET /api/projects/classrooms/:classroom_id/rounds/:round_id`

**Mô tả:** Lấy chi tiết một project round.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy project round thành công",
  "data": {
    "id": "507f1f77bcf86cd799439021",
    "classroom_id": "507f1f77bcf86cd799439011",
    "name": "Đợt 1",
    "description": "Đợt đồ án học kỳ 1",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-05-31T23:59:59Z",
    "max_members_per_group": 5
  }
}
```

---

### Lấy danh sách project rounds
**Endpoint:** `GET /api/projects/classrooms/:classroom_id/rounds`

**Mô tả:** Lấy danh sách project rounds theo lớp học.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lấy danh sách project rounds thành công",
  "data": [
    {
      "id": "507f1f77bcf86cd799439021",
      "name": "Đợt 1",
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-05-31T23:59:59Z"
    }
  ]
}
```

---

### Cập nhật project round
**Endpoint:** `PUT /api/projects/rounds`

**Mô tả:** Cập nhật thông tin project round.

**Request Body:**
```json
{
  "round_id": "507f1f77bcf86cd799439021",
  "name": "Đợt 1 (Updated)",
  "description": "Mô tả mới",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-06-30T23:59:59Z",
  "max_members_per_group": 6
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cập nhật project round thành công",
  "data": {...}
}
```

---

### Xóa project round
**Endpoint:** `DELETE /api/projects/classrooms/:classroom_id/rounds/:round_id`

**Mô tả:** Xóa project round.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Xóa project round thành công",
  "data": null
}
```

---

## Common Error Codes

| Error Code | HTTP Status | Mô tả |
|------------|-------------|-------|
| `BAD_REQUEST` | 400 | Dữ liệu request không hợp lệ |
| `INVALID_CREDENTIALS` | 401 | Sai username/password |
| `INVALID_TOKEN` | 401 | Token không hợp lệ hoặc hết hạn |
| `UNAUTHORIZED` | 401 | Chưa đăng nhập |
| `FORBIDDEN` | 403 | Không có quyền truy cập |
| `EMAIL_NOT_VERIFIED` | 403 | Email chưa xác minh |
| `USER_NOT_FOUND` | 404 | Không tìm thấy user |
| `CLASSROOM_NOT_FOUND` | 404 | Không tìm thấy lớp học |
| `POST_NOT_FOUND` | 404 | Không tìm thấy bài đăng |
| `GROUP_NOT_FOUND` | 404 | Không tìm thấy nhóm |
| `CHANNEL_NOT_FOUND` | 404 | Không tìm thấy channel |
| `INVITATION_NOT_FOUND` | 404 | Không tìm thấy lời mời |
| `EMAIL_EXISTS` | 409 | Email đã tồn tại |
| `ALREADY_IN_CLASSROOM` | 409 | Đã là thành viên của lớp |
| `ALREADY_CO_LECTURER` | 409 | Đã là trợ giảng của lớp |
| `ALREADY_LECTURER` | 409 | Đã là giảng viên của lớp |
| `INVITATION_ALREADY_EXISTS` | 409 | Đã có lời mời pending |
| `INVITATION_ALREADY_PROCESSED` | 409 | Lời mời đã được xử lý |
| `CLASSROOM_FULL` | 409 | Lớp học đã đầy |
| `STUDENT_CODE_NOT_IN_WHITELIST` | 409 | Mã SV không trong whitelist |
| `STUDENT_CODE_ALREADY_USED` | 409 | Mã SV đã được sử dụng |
| `INVALID_EMAIL_DOMAIN` | 400 | Email không đúng domain cho phép |
| `LOGIN_METHOD_MISMATCH` | 409 | Sai phương thức đăng nhập |
| `INTERNAL_ERROR` | 500 | Lỗi server |

---

## Notes

### Authentication Headers
Đối với các endpoint yêu cầu xác thực, gửi access token trong header:
```
Authorization: Bearer <access_token>
```

### Token Expiration
- **Access Token:** 15 phút
- **Refresh Token:** 7 ngày
- **OTP Codes:** 15 phút
- **Invitation Codes:** Tùy theo cấu hình lớp học

### User Types
- **Student:** Có field `student_code` khác null
- **Lecturer:** Field `student_code` là null

### Pagination
Hầu hết endpoints có phân trang đều sử dụng:
- `page`: Số trang (bắt đầu từ 1)
- `pageSize`: Số items mỗi trang (max: 100)

### File Upload
- Upload files qua `multipart/form-data`
- Files được lưu trên Cloudinary
- Trả về `public_id` để xóa sau này

---
