# Kiểm tra User trong MongoDB

## Backend đang kết nối:

- **Cluster**: `lkforumcluster.vugy4ss.mongodb.net`
- **Database**: `LKForum`
- **Collection**: `users`

## Cách kiểm tra trong MongoDB Compass:

### 1. Kết nối:

```
mongodb+srv://23520790_db_user:1234567890@lkforumcluster.vugy4ss.mongodb.net/
```

### 2. Chọn database `LKForum` (quan trọng!)

### 3. Mở collection `users`

### 4. Filter tìm user:

```json
{ "email": "your-email@gmail.com" }
```

Hoặc xem tất cả users:

```json
{}
```

## Nếu không thấy user:

### Kiểm tra logs khi đăng nhập:

```bash
docker-compose logs -f backend
```

Sau đó đăng nhập lại và xem logs có lỗi gì không.

## Có thể bạn của bạn đang dùng database khác?

Trong cluster `lkforumcluster.vugy4ss.mongodb.net` có thể có nhiều database:

- `LKForum` ← Backend production (đang dùng)
- `LKForum_dev`
- `test`
- Hoặc tên khác

Hỏi bạn của bạn xem họ đang check database nào trong Compass.
