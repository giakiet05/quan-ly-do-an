// Chạy script này trong MongoDB Compass shell hoặc mongosh
// để tạo classroom test cho student

// Thay YOUR_USER_ID bằng user ID thực của bạn
const STUDENT_ID = "YOUR_USER_ID"; // Lấy từ users collection
const LECTURER_ID = "YOUR_USER_ID"; // Có thể dùng cùng ID hoặc tạo user lecturer riêng

// Tạo classroom
db.classrooms.insertOne({
  name: "Đồ án Phát triển ứng dụng Web - K18",
  description: "Môn học về phát triển ứng dụng web full-stack với React và Node.js",
  avatar: "",
  semester: "HK1",
  year: 2024,
  status: "active",
  lecturer_id: ObjectId(LECTURER_ID),
  students: [
    {
      user_id: ObjectId(STUDENT_ID),
      full_name: "Nguyễn Văn A",
      avatar: ""
    }
  ],
  invitation_code: "WEB2024K18",
  require_approval: false,
  require_email_domain: false,
  whitelist_student_codes: [],
  general_channel_id: null,
  can_student_delete_group: true,
  project_rounds: [
    {
      _id: ObjectId(),
      name: "Đợt 1 - Đồ án giữa kỳ",
      start_date: new Date("2024-09-01"),
      end_date: new Date("2024-11-30"),
      description: "Đồ án giữa kỳ - Xây dựng ứng dụng web cơ bản",
      projects: [
        {
          _id: ObjectId(),
          classroom_id: null, // Sẽ update sau
          project_round_id: null, // Sẽ update sau
          title: "Hệ thống quản lý thư viện trực tuyến",
          amount: 5,
          description: "Xây dựng hệ thống quản lý thư viện với các tính năng mượn/trả sách, tìm kiếm, đặt chỗ",
          min_member: 2,
          max_member: 3,
          status: "approved"
        },
        {
          _id: ObjectId(),
          classroom_id: null,
          project_round_id: null,
          title: "Website bán hàng trực tuyến",
          amount: 3,
          description: "Xây dựng website thương mại điện tử với giỏ hàng, thanh toán, quản lý đơn hàng",
          min_member: 2,
          max_member: 4,
          status: "approved"
        }
      ],
      report_periods: [
        {
          _id: ObjectId(),
          title: "Báo cáo tiến độ tuần 1",
          description: "Nộp báo cáo tiến độ và kết quả đạt được trong tuần 1",
          file_type: ["pdf", "docx"],
          start_date: new Date("2024-09-08"),
          end_date: new Date("2024-09-15")
        }
      ],
      created_at: new Date(),
      is_deleted: false
    }
  ],
  created_at: new Date(),
  updated_at: new Date()
});

print("✅ Classroom created successfully!");
print("📝 Remember to:");
print("1. Replace STUDENT_ID and LECTURER_ID with actual ObjectId values");
print("2. Update classroom_id and project_round_id in projects after classroom is created");
