// Setup đầy đủ: Project Round → Project → Group cho user
const userId = ObjectId("696bc941be3b4acc6f83ed1b");
const classroomId = ObjectId("696bde961efa5a3b7ed8b75f");

print("=== SETUP COMPLETE DATA ===\n");

// ============================================
// BƯỚC 1: Thêm Project Round vào Classroom
// ============================================
print("📚 Bước 1: Tạo Project Round...");

const projectRoundId = ObjectId();
const projectId = ObjectId();

db.classrooms.updateOne(
  { _id: classroomId },
  {
    $push: {
      project_rounds: {
        _id: projectRoundId,
        name: "Đợt 1 - Đồ án cuối kỳ HK2",
        description: "Phát triển ứng dụng web full-stack",
        start_date: new Date("2024-02-01T00:00:00.000Z"),
        end_date: new Date("2024-05-31T23:59:59.000Z"),
        projects: [
          {
            _id: projectId,
            classroom_id: classroomId,
            project_round_id: projectRoundId,
            title: "Hệ thống quản lý thư viện trực tuyến",
            amount: 5,
            description: "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt trước",
            min_member: 2,
            max_member: 4,
            status: "approved",
            is_deleted: false
          },
          {
            _id: ObjectId(),
            classroom_id: classroomId,
            project_round_id: projectRoundId,
            title: "Website bán hàng trực tuyến",
            amount: 3,
            description: "Xây dựng website thương mại điện tử với giỏ hàng, thanh toán",
            min_member: 2,
            max_member: 4,
            status: "approved",
            is_deleted: false
          },
          {
            _id: ObjectId(),
            classroom_id: classroomId,
            project_round_id: projectRoundId,
            title: "Ứng dụng quản lý công việc nhóm",
            amount: 4,
            description: "Task management app với real-time collaboration",
            min_member: 2,
            max_member: 3,
            status: "approved",
            is_deleted: false
          }
        ]
      }
    }
  }
);

print("✅ Project Round đã tạo với 3 projects");
print(`   - Project Round ID: ${projectRoundId}`);
print(`   - Project ID (sẽ dùng): ${projectId}\n`);

// ============================================
// BƯỚC 2: Tạo Group cho User
// ============================================
print("👥 Bước 2: Tạo Group...");

const groupId = ObjectId();
const groupChannelId = ObjectId();

db.groups.insertOne({
  _id: groupId,
  classroom_id: classroomId,
  project_id: projectId,
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
      title: "Setup project structure",
      details: "Initialize React app and backend server",
      assign_to_ids: [userId],
      due_date: new Date("2024-09-20T23:59:59.000Z"),
      status: "in_progress"
    },
    {
      _id: ObjectId(),
      title: "Thiết kế database schema",
      details: "Thiết kế ERD và tạo collections MongoDB",
      assign_to_ids: [userId],
      due_date: new Date("2024-09-25T23:59:59.000Z"),
      status: "pending"
    }
  ],
  task_statuses: ["pending", "in_progress", "completed"],
  reports: [
    {
      _id: ObjectId(),
      title: "Báo cáo tuần 1",
      description: "Khảo sát yêu cầu và thiết kế hệ thống",
      files: [],
      status: "submitted",
      created_at: new Date("2024-02-10T10:00:00.000Z"),
      updated_at: new Date("2024-02-10T10:00:00.000Z")
    }
  ],
  setting: {
    allow_join_request: true
  },
  min_member: 2,
  max_member: 4,
  join_requests: [],
  join_invitations: []
});

print("✅ Group đã tạo");
print(`   - Group ID: ${groupId}\n`);

// ============================================
// BƯỚC 3: Tạo Channel cho Group
// ============================================
print("💬 Bước 3: Tạo Group Channel...");

db.channels.insertOne({
  _id: groupChannelId,
  name: "Hệ thống quản lý thư viện - Team Chat",
  type: "group",
  classroom_id: classroomId,
  group_id: groupId,
  members: [
    {
      user_id: userId,
      full_name: "Blue",
      avatar: "https://res.cloudinary.com/dkjxtedna/image/upload/v1768671866/lkforum/g6s4anxuuwjwdp7yss4w.jpg",
      role: "leader"
    }
  ],
  created_at: new Date(),
  updated_at: new Date()
});

print("✅ Group Channel đã tạo");
print(`   - Channel ID: ${groupChannelId}\n`);

// ============================================
// KIỂM TRA KẾT QUẢ
// ============================================
print("=== KIỂM TRA ===");

const classroom = db.classrooms.findOne({ _id: classroomId });
print(`✅ Classroom có ${classroom.project_rounds?.length || 0} project rounds`);
print(`✅ Project round có ${classroom.project_rounds?.[0]?.projects?.length || 0} projects`);

const groupCount = db.groups.countDocuments({ leader_id: userId });
print(`✅ User có ${groupCount} group(s)`);

const channelCount = db.channels.countDocuments({ group_id: groupId });
print(`✅ Group có ${channelCount} channel(s)`);

print("\n🎉 HOÀN TẤT! Bây giờ refresh trang để xem:\n");
print("   1. 'Lớp học của tôi' > Chi tiết lớp > Tab 'Tổng quan' → Hiển thị Project Round");
print("   2. 'Đề tài của tôi' → Hiển thị 'Hệ thống quản lý thư viện trực tuyến'");
print("   3. Click vào đề tài → Xem tasks, reports, team chat");
