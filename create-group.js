// Tạo group cho user trong classroom
const userId = ObjectId("696bc941be3b4acc6f83ed1b");
const classroomId = ObjectId("696bde961efa5a3b7ed8b75f");

// Lấy project từ classroom
const classroom = db.classrooms.findOne({ _id: classroomId });
const projectId = classroom.project_rounds[0].projects[0]._id;

print("Creating group...");

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

// Tạo channel cho group
db.channels.insertOne({
  _id: groupChannelId,
  name: "Group Chat",
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

print("✅ Channel created!");
print("\n🎉 Done! Refresh the page now.");
