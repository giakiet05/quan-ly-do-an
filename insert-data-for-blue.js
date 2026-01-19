// ============================================
// SCRIPT INSERT CLASSROOM CHO USER: Blue
// User ID: 696bc941be3b4acc6f83ed1b (NEW - Google OAuth)
// ============================================
// Chạy trong MongoDB Compass Shell (_MONGOSH tab)
// Đảm bảo đã chọn database ProjectGate trước khi chạy

// Generate ObjectIds cho các documents
const classroomId = ObjectId();
const projectRoundId = ObjectId();
const project1Id = ObjectId();
const project2Id = ObjectId();
const reportPeriod1Id = ObjectId();
const reportPeriod2Id = ObjectId();
const channelId = ObjectId();
const userId = ObjectId("696bc941be3b4acc6f83ed1b");

print("📝 Generated IDs:");
print("Classroom ID:", classroomId.toString());
print("Project Round ID:", projectRoundId.toString());
print("Project 1 ID:", project1Id.toString());
print("Channel ID:", channelId.toString());

// 1. Tạo Classroom
db.classrooms.insertOne({
  _id: classroomId,
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
      _id: userId,
      full_name: "Blue",
      avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg"
    }
  ],
  invitation_code: "WEB2024K18",
  auto_approve: true,
  require_email_domain: null,
  whitelist_student_codes: [],
  max_students: 100,
  general_channel_id: channelId,
  can_student_delete_group: true,
  project_rounds: [
    {
      _id: projectRoundId,
      name: "Đợt 1 - Đồ án giữa kỳ 2024",
      start_date: new Date("2024-09-01T00:00:00.000Z"),
      end_date: new Date("2024-11-30T23:59:59.000Z"),
      description: "Đợt đăng ký đồ án giữa kỳ - Sinh viên thực hiện xây dựng ứng dụng web theo đề tài đã chọn",
      report_periods: [
        {
          _id: reportPeriod1Id,
          title: "Báo cáo tiến độ tuần 1-2",
          description: "Nộp báo cáo phân tích yêu cầu, thiết kế database và giao diện prototype. Yêu cầu: file PDF hoặc Word, tối đa 10 trang.",
          file_type: ["pdf", "docx"],
          start_date: new Date("2024-09-08T00:00:00.000Z"),
          end_date: new Date("2024-09-22T23:59:59.000Z")
        },
        {
          _id: reportPeriod2Id,
          title: "Báo cáo tiến độ tuần 3-4",
          description: "Nộp báo cáo tiến độ implementation, demo video sản phẩm (nếu có). Trình bày các chức năng đã hoàn thành và kế hoạch tiếp theo.",
          file_type: ["pdf", "docx", "mp4"],
          start_date: new Date("2024-09-23T00:00:00.000Z"),
          end_date: new Date("2024-10-06T23:59:59.000Z")
        }
      ],
      created_at: new Date(),
      is_deleted: false
    }
  ],
  created_at: new Date(),
  updated_at: new Date()
});

print("✅ Classroom created!");

// 2. Tạo Projects trong collection riêng
db.projects.insertMany([
  {
    _id: project1Id,
    classroom_id: classroomId,
    project_round_id: projectRoundId,
    title: "Hệ thống quản lý thư viện trực tuyến",
    amount: 5,
    description: "Xây dựng hệ thống quản lý thư viện với các tính năng: mượn/trả sách, tìm kiếm nâng cao, đặt chỗ, quản lý thành viên, báo cáo thống kê. Hệ thống cho phép người dùng tra cứu sách, đặt mượn online và nhận thông báo.",
    min_member: 2,
    max_member: 3,
    status: "approved",
    created_at: new Date()
  },
  {
    _id: project2Id,
    classroom_id: classroomId,
    project_round_id: projectRoundId,
    title: "Website thương mại điện tử bán hàng trực tuyến",
    amount: 3,
    description: "Xây dựng website bán hàng trực tuyến với giỏ hàng, thanh toán VNPay/MoMo, quản lý đơn hàng, sản phẩm, khách hàng. Tích hợp chatbot hỗ trợ khách hàng và hệ thống đánh giá sản phẩm.",
    min_member: 2,
    max_member: 4,
    status: "approved",
    created_at: new Date()
  }
]);

print("✅ Projects created!");

// 3. Tạo Channel cho classroom
db.channels.insertOne({
  _id: channelId,
  name: "general",
  type: "group",
  classroom_id: classroomId,
  members: [
    {
      user_id: userId,
      full_name: "Blue",
      avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg"
    }
  ],
  created_at: new Date(),
  updated_at: new Date()
});

print("✅ Channel created!");

// 4. Tạo Group để user đăng ký project
const groupId = ObjectId();
const groupChannelId = ObjectId();

db.groups.insertOne({
  _id: groupId,
  classroom_id: classroomId,
  project_id: project1Id,
  group_channel_id: groupChannelId,
  leader_id: userId,
  members: [
    {
      user_id: userId,
      full_name: "Blue",
      avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg"
    }
  ],
  tasks: [
    {
      _id: ObjectId(),
      title: "Phân tích yêu cầu hệ thống",
      details: "Thu thập và phân tích yêu cầu chức năng, phi chức năng của hệ thống quản lý thư viện",
      assign_to_ids: [userId],
      due_date: new Date("2024-09-15T23:59:59.000Z"),
      status: "in_progress"
    },
    {
      _id: ObjectId(),
      title: "Thiết kế database",
      details: "Thiết kế ERD, schema MongoDB cho hệ thống",
      assign_to_ids: [userId],
      due_date: new Date("2024-09-20T23:59:59.000Z"),
      status: "pending"
    }
  ],
  task_statuses: ["pending", "in_progress", "completed"],
  reports: [],
  setting: {
    allow_join_request: true
  },
  min_member: 2,
  max_member: 3,
  join_requests: [],
  join_invitations: []
});

print("✅ Group created!");

// 4. Tạo Channel cho group
db.channels.insertOne({
  _id: groupChannelId,
  name: "Nhóm - Hệ thống quản lý thư viện",
  type: "group",
  group_id: groupId,
  members: [
    {
      user_id: userId,
      full_name: "Blue",
      avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg"
    }
  ],
  created_at: new Date(),
  updated_at: new Date()
});

print("✅ Group channel created!");

// 5. Tạo vài tin nhắn mẫu trong channel classroom
db.messages.insertMany([
  {
    _id: ObjectId(),
    channel_id: channelId,
    sender_id: userId,
    content: "Chào mọi người! Đây là kênh thảo luận chung của lớp học.",
    created_at: new Date("2024-09-01T08:00:00.000Z"),
    updated_at: new Date("2024-09-01T08:00:00.000Z")
  },
  {
    _id: ObjectId(),
    channel_id: channelId,
    sender_id: userId,
    content: "Các bạn đã xem qua danh sách đề tài chưa? Hãy lựa chọn đề tài phù hợp nhé!",
    created_at: new Date("2024-09-01T10:30:00.000Z"),
    updated_at: new Date("2024-09-01T10:30:00.000Z")
  }
]);

print("✅ Sample messages created!");

// 6. Tạo notification
db.notifications.insertOne({
  _id: ObjectId(),
  user_id: userId,
  type: "classroom_invitation",
  title: "Chào mừng đến lớp học!",
  content: "Bạn đã tham gia lớp 'Đồ án Phát triển ứng dụng Web - K18'. Hãy bắt đầu chọn đề tài nhé!",
  is_read: false,
  created_at: new Date()
});

print("✅ Notification created!");

// Summary
print("\n🎉 ========================================");
print("🎉 TẠO DỮ LIỆU THÀNH CÔNG!");
print("🎉 ========================================");
print("\n📊 Thông tin đã tạo:");
print("✓ 1 Classroom: Đồ án Phát triển ứng dụng Web - K18");
print("✓ 1 Project Round với 2 projects");
print("✓ 2 Report periods");
print("✓ 1 General channel cho classroom");
print("✓ 1 Group đã đăng ký project");
print("✓ 1 Group channel");
print("✓ 2 Tasks trong group");
print("✓ 2 Messages mẫu");
print("✓ 1 Notification");
print("\n🚀 Bây giờ refresh frontend để xem kết quả!");
print("📍 Vào sidebar Student → Lớp học của tôi");
print("📍 Hoặc sidebar Student → Đề tài của tôi");
