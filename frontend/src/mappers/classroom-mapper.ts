import type { ClassroomResponse } from "../dtos/classroom-dto";
import type { ClassItem } from "../types/class";

export function mapClassroomToClassItem(
    cls: ClassroomResponse
): ClassItem {
    return {
        id: cls.id,
        name: cls.name,
        school: "—",
        description: cls.description,
        studentCount: cls.students.length,
        semester: `${cls.semester} ${cls.year}`,
        status: cls.status === "active" ? "active" : "inactive",
    };
}
