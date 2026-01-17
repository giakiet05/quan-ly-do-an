// ============================================
// SCRIPT KIỂM TRA CLASSROOMS TRONG DATABASE
// ============================================
// Chạy trong MongoDB Compass Shell (_MONGOSH tab)
// Đảm bảo đã chọn database ProjectGate trước khi chạy

const userId = ObjectId("696bc941be3b4acc6f83ed1b");

print("🔍 Kiểm tra user ID:", userId.toString());

// 1. Kiểm tra tổng số classrooms
const totalClassrooms = db.classrooms.countDocuments();
print("\n📊 Tổng số classrooms trong database:", totalClassrooms);

// 2. Tìm classrooms mà user là lecturer
const asLecturer = db.classrooms.find({ lecturer_id: userId }).toArray();
print("\n👨‍🏫 Classrooms mà user là lecturer:", asLecturer.length);
if (asLecturer.length > 0) {
  print("IDs:", asLecturer.map(c => c._id.toString()));
}

// 3. Tìm classrooms mà user là student (trong students array)
const asStudent = db.classrooms.find({ 
  "students.user_id": userId 
}).toArray();
print("\n👨‍🎓 Classrooms mà user là student:", asStudent.length);
if (asStudent.length > 0) {
  print("IDs:", asStudent.map(c => c._id.toString()));
  print("\nChi tiết classroom đầu tiên:");
  print("- Name:", asStudent[0].name);
  print("- Status:", asStudent[0].status);
  print("- Students count:", asStudent[0].students?.length || 0);
  print("- Project rounds count:", asStudent[0].project_rounds?.length || 0);
}

// 4. Kiểm tra groups của user
const userGroups = db.groups.find({
  "members.user_id": userId
}).toArray();
print("\n👥 Groups mà user là member:", userGroups.length);
if (userGroups.length > 0) {
  print("Group IDs:", userGroups.map(g => g._id.toString()));
  print("Classroom IDs:", userGroups.map(g => g.classroom_id?.toString()));
  print("Project IDs:", userGroups.map(g => g.project_id?.toString()));
}

print("\n✅ Kiểm tra hoàn tất!");
