// ============================================
// SCRIPT: XÓA VÀ TẠO LẠI CLASSROOM
// ============================================
// Chạy trong MongoDB Compass Shell (_MONGOSH tab)

const userId = ObjectId("696bc941be3b4acc6f83ed1b");
const oldClassroomId = ObjectId("696bcb554eea63c26694751d");

print("🗑️ Deleting old classroom...");
db.classrooms.deleteOne({ _id: oldClassroomId });

print("✅ Creating new classroom with correct structure...");

// Tạo classroom mới với structure ĐÚNG
const newClassroomId = ObjectId();
const projectRoundId = ObjectId();
const project1Id = ObjectId();
const channelId = ObjectId();

db.classrooms.insertOne({
  _id: newClassroomId,
  name: "Đồ án Phát triển ứng dụng Web - K18",
  description: "Môn học về phát triển ứng dụng web full-stack với React, Node.js và MongoDB. Sinh viên sẽ thực hiện đồ án xây dựng ứng dụng web hoàn chỉnh.",
  avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1234567890/classroom-web.jpg",
  semester: "HK1",
  year: 2024,
  status: "active",
  lecturer: {
    _id: userId,
    full_name: "Blue",
    avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg"
  },
  students: [
    {
      user_id: userId,
      full_name: "Blue",
      avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg"
    }
  ],
  invitation_code: "WEB2024K18",
  auto_approve: true,
  require_email_domain: null,
  whitelist_student_code: [],
  max_students: 100,
  general_channel_id: channelId,
  can_student_delete_group: true,
  project_rounds: [
    {
      _id: projectRoundId,
      name: "Đợt 1 - Đồ án giữa kỳ 2024",
      start_date: new Date("2024-09-01T00:00:00.000Z"),
      end_date: new Date("2024-11-30T23:59:59.000Z"),
      description: "Đợt đăng ký đồ án giữa kỳ",
      projects: [
        {
          _id: project1Id,
          classroom_id: newClassroomId,
          project_round_id: projectRoundId,
          title: "Hệ thống quản lý thư viện trực tuyến",
          amount: 5,
          description: "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm nâng cao.",
          min_member: 2,
          max_member: 3,
          status: "approved"
        }
      ],
      report_periods: [],
      created_at: new Date(),
      is_deleted: false
    }
  ],
  created_at: new Date()
});

print("✅ Classroom created with ID:", newClassroomId.toString());

// Cập nhật group với classroom_id mới
print("🔧 Updating group with new classroom_id...");
db.groups.updateMany(
  { classroom_id: oldClassroomId },
  { $set: { classroom_id: newClassroomId, project_id: project1Id } }
);

// Cập nhật channel với classroom_id mới
print("🔧 Updating channels...");
db.channels.updateMany(
  { classroom_id: oldClassroomId },
  { $set: { classroom_id: newClassroomId } }
);

print("\n✅ All done! Test the API now.");
