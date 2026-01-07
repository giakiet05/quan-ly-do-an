import type { File, GroupSetting } from "../models";

// ==================== REQUEST DTOs ====================

export interface CreateGroupRequest {
  classroomId: string;
  projectId: string;
  projectRoundId: string;
  memberIds: string[];
}

export interface UpdateGroupRequest {
  groupId: string;
  projectId?: string;
  leaderId?: string;
  setting?: GroupSetting;
}

export interface UpdateGroupMembersRequest {
  groupId: string;
  memberIds: string[];
  action: "add" | "remove";
}

export interface CreateJoinGroupRequest {
  groupId: string;
  message: string;
}

export interface UpdateJoinGroupRequest {
  groupId: string;
  requestId: string;
  status: "pending" | "accepted" | "rejected";
}

export interface CreateGroupInvitationRequest {
  groupId: string;
  recipientId: string;
}

export interface UpdateGroupInvitationRequest {
  groupId: string;
  invitationId: string;
  status: "pending" | "accepted" | "rejected";
}

export interface CreateTaskRequest {
  groupId: string;
  title: string;
  details?: string;
  assignToIds: string[];
  dueDate: string; // ISO date string
  status: string;
}

export interface UpdateTaskRequest {
  groupId: string;
  taskId: string;
  title?: string;
  details?: string;
  assignToIds?: string[];
  dueDate?: string;
  status?: string;
}

export interface CreateReportRequest {
  groupId: string;
  title: string;
  content: string;
  files: File[];
}

export interface UpdateReportRequest {
  groupId: string;
  reportId: string;
  title?: string;
  content?: string;
  files?: File[];
}

export interface CreateReportFeedbackRequest {
  groupId: string;
  reportId: string;
  content: string;
  grade: string;
}

export interface UpdateReportFeedbackRequest {
  groupId: string;
  feedbackId: string;
  content?: string;
  grade?: string;
}

export interface GetGroupsFilterQuery {
  classroomId: string;
  projectId?: string;
  memberId?: string;
}

// ==================== RESPONSE DTOs ====================

export interface GroupResponse {
  id: string;
  classroomId: string;
  projectId: string;
  groupChannelId: string;
  leaderId: string;
  members: UserInfoResponse[];
  tasks: TaskResponse[];
  taskStatuses: string[];
  reports: ReportResponse[];
  setting: GroupSetting;
}

export interface TaskResponse {
  id: string;
  title: string;
  assignToIds: string[];
  dueDate: string;
  status: string;
}

export interface ReportResponse {
  id: string;
  title: string;
  content: string;
  files: File[];
  feedback: ReportFeedbackResponse;
}

export interface ReportFeedbackResponse {
  content: string;
  grade: string;
  lecturerId: string;
  commentedAt: string;
}

export interface UserInfoResponse {
  userId: string;
  username: string;
  avatar?: Image;
}

// Import Image from models
import type { Image } from "../models";
