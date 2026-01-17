import type { ClassroomResponse, CreateClassroomRequest } from "../dtos/classroom-dto";
import type { ClassItem, CreateClassRequest } from "../types/class";

export function mapClassroomToClassItem(
    cls: ClassroomResponse
): ClassItem {
    return {
        id: cls.id,
        name: cls.name,
        avatar: cls.avatar || "",
        description: cls.description || "",
        studentCount: cls.students?.length ?? 0,
        semester: cls.semester || "",
        status: cls.status === "active" ? "active" : "inactive",
    };
}
export function mapUIRequestToDTO(
    uiData: CreateClassRequest
): CreateClassroomRequest {
    return {
        name: uiData.name,
        description: uiData.description ?? "",
        semester: uiData.semester,
        year: new Date().getFullYear(),
        avatar: uiData.avatar ?? ""
    };
}
