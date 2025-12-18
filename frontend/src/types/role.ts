export const Roles = {
    LECTURER: "LECTURER",
    STUDENT: "STUDENT",
} as const;

export type Role = (typeof Roles)[keyof typeof Roles];
