# Hướng dẫn thêm classroom test cho student

## Bước 1: Lấy User ID của bạn

1. Mở MongoDB Compass
2. Kết nối: `mongodb+srv://23520790_db_user:1234567890@lkforumcluster.vugy4ss.mongodb.net/`
3. Database: `LKForum` → Collection: `users`
4. Tìm user của bạn: `{ "email": "your-email@gmail.com" }`
5. Copy `_id` (ví dụ: `677a1234567890abcdef1234`)

## Bước 2: Tạo classroom đơn giản

Vào collection `classrooms` → Insert Document:

```json
{
  "name": "Đồ án Web - K18",
  "description": "Lớp học test",
  "avatar": "",
  "semester": "HK1",
  "year": 2024,
  "status": "active",
  "lecturer_id": { "$oid": "PASTE_YOUR_USER_ID_HERE" },
  "students": [
    {
      "user_id": { "$oid": "PASTE_YOUR_USER_ID_HERE" },
      "full_name": "Your Name",
      "avatar": ""
    }
  ],
  "invitation_code": "TEST2024",
  "require_approval": false,
  "require_email_domain": false,
  "whitelist_student_codes": [],
  "general_channel_id": null,
  "can_student_delete_group": true,
  "project_rounds": [
    {
      "_id": { "$oid": "677a1111111111111111111a" },
      "name": "Đợt 1 - Giữa kỳ",
      "start_date": { "$date": "2024-09-01T00:00:00.000Z" },
      "end_date": { "$date": "2024-11-30T00:00:00.000Z" },
      "description": "Đồ án giữa kỳ",
      "projects": [
        {
          "_id": { "$oid": "677a2222222222222222222b" },
          "title": "Hệ thống quản lý thư viện",
          "amount": 5,
          "description": "Xây dựng hệ thống quản lý thư viện trực tuyến",
          "min_member": 2,
          "max_member": 3,
          "status": "approved"
        }
      ],
      "report_periods": [],
      "created_at": { "$date": "2024-09-01T00:00:00.000Z" },
      "is_deleted": false
    }
  ],
  "created_at": { "$date": "2024-09-01T00:00:00.000Z" },
  "updated_at": { "$date": "2024-09-01T00:00:00.000Z" }
}
```

**Lưu ý:** Thay `PASTE_YOUR_USER_ID_HERE` bằng ID thực của bạn (ví dụ: `677a1234567890abcdef1234`)

## Bước 3: Update classroom_id và project_round_id

Sau khi insert classroom, MongoDB sẽ tạo `_id` cho classroom. Bạn cần update:

1. Mở classroom vừa tạo
2. Copy `_id` của classroom
3. Edit document:
   - Trong `project_rounds[0].projects[0]`:
     - Thêm `"classroom_id": {"$oid": "classroom_id_vừa_copy"}`
     - Thêm `"project_round_id": {"$oid": "677a1111111111111111111a"}`

## Bước 4: Tạo channel cho classroom (optional)

Nếu muốn test chat, vào collection `channels` → Insert:

```json
{
  "name": "general",
  "type": "group",
  "classroom_id": { "$oid": "PASTE_CLASSROOM_ID_HERE" },
  "members": [
    {
      "user_id": { "$oid": "PASTE_YOUR_USER_ID_HERE" },
      "full_name": "Your Name",
      "avatar": ""
    }
  ],
  "created_at": { "$date": "2024-09-01T00:00:00.000Z" }
}
```

Sau đó update classroom:

- Thêm `"general_channel_id": {"$oid": "channel_id_vừa_tạo"}`

## Bước 5: Reload frontend

Refresh trang web → Xem mục "Lớp học của tôi" ở sidebar student.

---

## Nhanh hơn: Dùng API

Nếu có Postman/Thunder Client:

**POST** `http://localhost:8080/api/classrooms`
Headers:

- `Authorization: Bearer YOUR_TOKEN`
- `Content-Type: application/json`

Body:

```json
{
  "name": "Đồ án Web - K18",
  "description": "Lớp test",
  "semester": "HK1",
  "year": 2024
}
```

Nhưng điều này yêu cầu tài khoản có role="lecturer".
