import type { Role } from "./role";


export type SidebarItem = {
    id: string;
    label: string;
    route?: string;
    icon?: string;
    badge?: number;
    children?: SidebarItem[];
    roles: Role[];
};
