import { apiFetch } from "./api";
import type { Classroom } from "../models";

// ==================== CLASSROOM SERVICES ====================

/**
 * Get all classrooms for current user
 */
export async function getMyClassrooms(): Promise<Classroom[]> {
  return apiFetch<Classroom[]>(`/api/classrooms/my`);
}

/**
 * Get a single classroom by ID
 */
export async function getClassroom(classroomId: string): Promise<Classroom> {
  return apiFetch<Classroom>(`/api/classrooms/${classroomId}`);
}

/**
 * Create a new classroom (lecturer only)
 */
export async function createClassroom(
  classroomData: Partial<Classroom>
): Promise<Classroom> {
  return apiFetch<Classroom>(`/api/classrooms`, {
    method: "POST",
    body: JSON.stringify(classroomData),
  });
}

/**
 * Update classroom settings
 */
export async function updateClassroomSettings(
  classroomId: string,
  settings: Partial<Classroom["setting"]>
): Promise<Classroom> {
  return apiFetch<Classroom>(`/api/classrooms/${classroomId}/settings`, {
    method: "PUT",
    body: JSON.stringify(settings),
  });
}

/**
 * Add students to classroom
 */
export async function addStudentsToClassroom(
  classroomId: string,
  studentIds: string[]
): Promise<Classroom> {
  return apiFetch<Classroom>(`/api/classrooms/${classroomId}/students`, {
    method: "POST",
    body: JSON.stringify({ studentIds }),
  });
}

/**
 * Remove student from classroom
 */
export async function removeStudentFromClassroom(
  classroomId: string,
  studentId: string
): Promise<void> {
  return apiFetch<void>(`/api/classrooms/${classroomId}/students/${studentId}`, {
    method: "DELETE",
  });
}

/**
 * Join classroom with invitation code (student)
 */
export async function joinClassroom(invitationCode: string): Promise<Classroom> {
  return apiFetch<Classroom>(`/api/classrooms/join`, {
    method: "POST",
    body: JSON.stringify({ invitationCode }),
  });
}

/**
 * Leave a classroom
 */
export async function leaveClassroom(classroomId: string): Promise<void> {
  return apiFetch<void>(`/api/classrooms/${classroomId}/leave`, {
    method: "POST",
  });
}

/**
 * Delete a classroom (lecturer only)
 */
export async function deleteClassroom(classroomId: string): Promise<void> {
  return apiFetch<void>(`/api/classrooms/${classroomId}`, {
    method: "DELETE",
  });
}
