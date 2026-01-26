import { apiFetch } from "./api";
import type { Group, Task, Report, JoinGroupRequest } from "../models";

// ==================== GROUP SERVICES ====================

/**
 * Get all groups in a classroom
 */
export async function getClassroomGroups(
  classroomId: string
): Promise<Group[]> {
  return apiFetch<Group[]>(`/api/classrooms/${classroomId}/groups`);
}

/**
 * Get groups with filters (e.g., by project_round_id, user_id, etc.)
 */
export async function getGroupsFilter(
  filters: Record<string, string | number>
): Promise<Group[]> {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, value]) => {
    params.append(key, String(value));
  });

  return apiFetch<Group[]>(`/api/groups?${params.toString()}`);
}

/**
 * Get a single group by ID
 */
export async function getGroup(groupId: string): Promise<Group> {
  return apiFetch<Group>(`/api/groups/${groupId}`);
}

/**
 * Get user's group for a project
 */
export async function getMyProjectGroup(projectId: string): Promise<Group> {
  return apiFetch<Group>(`/api/projects/${projectId}/my-group`);
}

/**
 * Create a new group
 */
export async function createGroup(
  classroomId: string,
  projectId: string,
  projectRoundId: string
): Promise<Group> {
  return apiFetch<Group>(`/api/groups`, {
    method: "POST",
    body: JSON.stringify({
      classroom_id: classroomId,
      project_id: projectId,
      project_round_id: projectRoundId,
    }),
  });
}

/**
 * Update group settings
 */
export async function updateGroupSettings(
  groupId: string,
  settings: Partial<Group["setting"]>
): Promise<Group> {
  return apiFetch<Group>(`/api/groups/${groupId}/settings`, {
    method: "PUT",
    body: JSON.stringify(settings),
  });
}

/**
 * Delete a group
 */
export async function deleteGroup(groupId: string): Promise<void> {
  return apiFetch<void>(`/api/groups/${groupId}`, {
    method: "DELETE",
  });
}

// ==================== TASK SERVICES ====================

/**
 * Get all tasks for a group
 */
export async function getGroupTasks(groupId: string): Promise<Task[]> {
  return apiFetch<Task[]>(`/api/groups/${groupId}/tasks`);
}

/**
 * Create a new task
 */
export async function createTask(
  groupId: string,
  taskData: Partial<Task>
): Promise<Task> {
  return apiFetch<Task>(`/api/groups/${groupId}/tasks`, {
    method: "POST",
    body: JSON.stringify(taskData),
  });
}

/**
 * Update a task
 */
export async function updateTask(
  groupId: string,
  taskId: string,
  updates: Partial<Task>
): Promise<Task> {
  return apiFetch<Task>(`/api/groups/${groupId}/tasks/${taskId}`, {
    method: "PUT",
    body: JSON.stringify(updates),
  });
}

/**
 * Delete a task
 */
export async function deleteTask(groupId: string, taskId: string): Promise<void> {
  return apiFetch<void>(`/api/groups/${groupId}/tasks/${taskId}`, {
    method: "DELETE",
  });
}

// ==================== REPORT SERVICES ====================

/**
 * Get all reports for a group
 */
export async function getGroupReports(groupId: string): Promise<Report[]> {
  return apiFetch<Report[]>(`/api/groups/${groupId}/reports`);
}

/**
 * Submit a new report
 */
export async function submitReport(
  groupId: string,
  reportData: FormData
): Promise<Report> {
  return apiFetch<Report>(`/api/groups/${groupId}/reports`, {
    method: "POST",
    body: reportData,
  });
}

/**
 * Update a report
 */
export async function updateReport(
  groupId: string,
  reportId: string,
  reportData: FormData
): Promise<Report> {
  return apiFetch<Report>(`/api/groups/${groupId}/reports/${reportId}`, {
    method: "PUT",
    body: reportData,
  });
}

/**
 * Add feedback to a report (lecturer only)
 */
export async function addReportFeedback(
  groupId: string,
  reportId: string,
  feedback: {
    content: string;
    grade: string;
  }
): Promise<Report> {
  return apiFetch<Report>(
    `/api/groups/${groupId}/reports/${reportId}/feedback`,
    {
      method: "POST",
      body: JSON.stringify(feedback),
    }
  );
}

// ==================== JOIN REQUEST SERVICES ====================

/**
 * Send a join request to a group
 */
export async function sendJoinRequest(
  groupId: string,
  message: string
): Promise<JoinGroupRequest> {
  return apiFetch<JoinGroupRequest>(`/api/groups/${groupId}/join-requests`, {
    method: "POST",
    body: JSON.stringify({ message }),
  });
}

/**
 * Get all join requests for a group (leader only)
 */
export async function getGroupJoinRequests(
  groupId: string
): Promise<JoinGroupRequest[]> {
  return apiFetch<JoinGroupRequest[]>(`/api/groups/${groupId}/join-requests`);
}

/**
 * Accept a join request (leader only)
 */
export async function acceptJoinRequest(
  groupId: string,
  requestId: string
): Promise<void> {
  return apiFetch<void>(
    `/api/groups/${groupId}/join-requests/${requestId}/accept`,
    {
      method: "POST",
    }
  );
}

/**
 * Reject a join request (leader only)
 */
export async function rejectJoinRequest(
  groupId: string,
  requestId: string
): Promise<void> {
  return apiFetch<void>(
    `/api/groups/${groupId}/join-requests/${requestId}/reject`,
    {
      method: "POST",
    }
  );
}

/**
 * Remove a member from group (leader only)
 */
export async function removeMember(
  groupId: string,
  memberId: string
): Promise<void> {
  return apiFetch<void>(`/api/groups/${groupId}/members/${memberId}`, {
    method: "DELETE",
  });
}

// ==================== GROUP INVITATIONS ====================

/**
 * Send an invitation to join group (leader only)
 */
export async function sendGroupInvitation(
  groupId: string,
  recipientId: string
): Promise<void> {
  const payload = { groupId, recipientId };
  console.log("📤 Gửi lời mời nhóm:", {
    url: `/api/groups/invitations`,
    method: "POST",
    body: JSON.stringify(payload, null, 2), // in đẹp
    groupId,
    recipientId,
  });

  return apiFetch<void>(`/api/groups/invitations`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

/**
 * Get all invitations sent by this group (leader only)
 */
export async function getGroupInvitations(groupId: string): Promise<any[]> {
  return apiFetch<any[]>(`/api/groups/${groupId}/invitations`);
}

/**
 * Get all invitations received by current user
 */
export async function getMyGroupInvitations(): Promise<any[]> {
  return apiFetch<any[]>("/api/groups/my-invitations");
}

/**
 * Accept a group invitation
 */
export async function acceptGroupInvitation(
  groupId: string,
  invitationId: string
): Promise<void> {
  return apiFetch<void>(
    `/api/groups/${groupId}/invitations/${invitationId}/accept`,
    {
      method: "POST",
    }
  );
}

/**
 * Reject a group invitation
 */
export async function rejectGroupInvitation(
  groupId: string,
  invitationId: string
): Promise<void> {
  return apiFetch<void>(
    `/api/groups/${groupId}/invitations/${invitationId}/reject`,
    {
      method: "POST",
    }
  );
}
