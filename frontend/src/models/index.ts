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
  // TODO: Define file structure
}

export interface UserInfo {
  id: string;
  username: string;
  avatar?: Image;
}

// ==================== USER ====================

export type AuthProvider = "local" | "google";

export interface User {
  id: string;
  username: string;
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

export interface ProjectRound {
  id: string;
  name: string;
  startDate: string;
  endDate: string;
  description: string;
  projects: Project[];
  createdAt: string;
  isDeleted: boolean;
}

// ==================== CLASSROOM ====================

export interface ClassroomSetting {
  canStudentDeleteGroup: boolean;
}

export interface Classroom {
  id: string;
  universityId: string;
  universityName: string;
  name: string;
  avatar: string;
  generalChannelId: string;
  lecturer: UserInfo;
  students: UserInfo[];
  projectRounds: ProjectRound[];
  setting: ClassroomSetting;
  createdAt: string;
}

// ==================== GROUP ====================

export interface Task {
  id?: string;
  title: string;
  details?: string;
  assignToIds: string[];
  dueDate: string;
  status: string;
}

export interface ReportFeedback {
  content: string;
  grade: string;
  lecturerId: string;
  commentedAt: string;
  updatedAt: string;
}

export interface Report {
  id: string;
  title: string;
  content: string;
  files: File[];
  feedback: ReportFeedback;
  createdAt: string;
  updatedAt: string;
}

export interface GroupSetting {
  allowJoinRequest: boolean;
}

export type RequestStatus = "pending" | "accepted" | "rejected";

export interface JoinGroupRequest {
  id: string;
  userId: string;
  userInfo: UserInfo;
  status: RequestStatus;
  message: string;
  requestedAt: string;
}

export interface Group {
  id: string;
  classroomId: string;
  projectId: string;
  groupChannelId?: string;
  leaderId: string;
  members: UserInfo[];
  tasks: Task[];
  taskStatuses: string[];
  reports: Report[];
  setting: GroupSetting;
  joinRequests: JoinGroupRequest[];
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
