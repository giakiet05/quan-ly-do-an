// ==================== COMMON TYPES ====================

export interface Image {
  url: string;
  publicId: string;
  uploadedAt: string;
}

export interface Video {
  url: string;
  publicId: string;
  uploadedAt: string;
}

export interface File {
  fileName: string;
  fileURL: string;
  fileSize: number; // bytes
  mimeType: string;
}

export interface Attachment {
  fileName: string;
  fileURL: string;
  fileSize: number;
  mimeType: string;
}

export interface UserInfo {
  id: string;
  fullName: string;
  avatar?: Image;
  studentCode?: string;
}

export interface UserInfoResponse {
  userId: string;
  fullName: string;
  avatar: string;
}

// ==================== USER ====================

export type AuthProvider = "local" | "google";

export interface User {
  id: string;
  username: string;
  fullName?: string;
  email: string;
  authProvider: AuthProvider;
  avatar?: Image;
  createdAt?: string;
  deletedAt?: string;
}

// ==================== PROJECT ====================

export type ProjectStatus = "pending" | "approved" | "ongoing" | "completed";

export interface Project {
  id: string;
  classroomId: string;
  projectRoundId: string;
  title: string;
  amount: number;
  description: string;
  minMember: number;
  maxMember: number;
  status: ProjectStatus;
}

export interface ReportPeriod {
  id: string;
  title: string;
  description: string;
  fileType: string[];
  startDate: string;
  endDate: string;
}

export interface ProjectRound {
  id: string;
  name: string;
  startDate: string;
  endDate: string;
  description: string;
  projects: Project[];
  reportPeriods: ReportPeriod[];
  createdAt: string;
  isDeleted: boolean;
}

export type ClassroomStatus = "active" | "inactive" | "archived";

export interface Classroom {
  id: string;
  name: string;
  description: string;
  avatar: string;
  semester: string; // "Fall 2026", "Spring 2027"
  year: number;
  status: ClassroomStatus;
  generalChannelId: string;
  lecturer: UserInfo;
  students: UserInfo[];
  projectRounds: ProjectRound[];

  // Settings
  invitationCode: string;
  whitelistStudentCode?: string[];
  requireEmailDomain?: string;
  maxStudents: number;
  autoApprove: boolean;
  canStudentDeleteGroup: boolean;

  createdAt: string;
}

export type JoinRequestStatus = "pending" | "approved" | "rejected";

export interface ClassroomJoinRequest {
  id: string;
  classroomId: string;
  invitationCode: string;
  userId: string;
  status: JoinRequestStatus;
  createdAt: string;
  reviewedAt?: string;
  reviewedBy?: string;
}

export interface ClassPost {
  id: string;
  classroomId: string;
  author: UserInfo;
  title: string;
  content: string;
  attachments?: Attachment[];
  isPinned: boolean;
  createdAt: string;
  updatedAt: string;
}

export type JoinGroupRequestStatus = "pending" | "accepted" | "rejected";

export interface JoinGroupRequest {
  id: string;
  userId: string;
  userInfo: UserInfo;
  status: JoinGroupRequestStatus;
  message: string;
  requestedAt: string;
  updatedAt: string;
}

export interface JoinGroupInvitation {
  id: string;
  groupId: string;
  recipientId: string;
  status: JoinGroupRequestStatus;
  sentAt: string;
  respondedAt?: string;
}

export interface Group {
  id: string;
  classroomId: string;
  projectId: string;
  groupChannelId?: string;
  leaderId: string;
  members: UserInfoResponse[];  // Backend sends UserInfoResponse, not UserInfo
  tasks: Task[];
  taskStatuses: string[];
  reports: Report[];
  setting: GroupSetting;
  minMember: number;
  maxMember: number;
  joinRequests: JoinGroupRequest[];
  joinInvitation: JoinGroupInvitation[];
  createdAt: string;
  updatedAt: string;
}

export interface GroupSetting {
  allowJoinRequest: boolean;
}

// ==================== TASK ====================

export interface Task {
  id?: string;
  title: string;
  details?: string;
  assignToIDs: string[];
  dueDate: string;
  status: string;
}

// ==================== REPORT ====================

export interface Report {
  id?: string;
  reportPeriodId: string;
  title: string;
  content: string;
  files: File[];
  feedback: ReportFeedback;
  createdAt: string;
  updatedAt: string;
}

export interface ReportFeedback {
  content: string;
  grade: string;
  lecturerId: string;
  commentedAt: string;
  updatedAt: string;
}

// ==================== CHANNEL ====================

export type ChannelStatus = "active" | "block";

export interface ChannelUserSetting {
  userId: string;
  notification: boolean;
  typingIndicator: boolean;
  isDeleted: boolean;
}

export interface ChannelSetting {
  allowMemberMessage: boolean;
  allowMedia: boolean;
  allowAttachments: boolean;
  allowPinMessages: boolean;
  allowMentions: boolean;
}

export interface Channel {
  id: string;
  adminIds: string[];
  members: UserInfo[];
  userSettings: ChannelUserSetting[];
  background?: string;
  status: ChannelStatus;
  setting: ChannelSetting;
  createdAt: string;
  updatedAt: string;
}

// ==================== MESSAGE ====================

export type MessageType = "user" | "system";

export interface Message {
  id: string;
  channelId: string;
  senderId?: string;
  senderUsername?: string;
  type: MessageType;
  content: string;
  readBy: string[];
  isSend: boolean;
  createdAt: string;
  isDeleted: boolean;
  deletedAt?: string;
}

// ==================== NOTIFICATION ====================

export type NotificationType =
  | "comment"
  | "like"
  | "follow"
  | "mention"
  | "new_message"
  | "system";

export interface Notification {
  id: string;
  recipientId: string;
  actorId?: string;
  type: NotificationType;
  message: string;
  link?: string;
  isRead: boolean;
  metadata?: Record<string, any>;
  createdAt: string;
}

// ==================== CLASSROOM INVITATION ====================

export type InvitationStatus =
  | "pending"
  | "accepted"
  | "rejected"
  | "expired"
  | "canceled";

export interface ClassroomInvitation {
  id: string;
  email: string;
  classroomId: string;
  invitedBy: string;
  invitedTo?: string;
  status: InvitationStatus;
  expiresAt: string;
  createdAt: string;
  updatedAt?: string;
}
