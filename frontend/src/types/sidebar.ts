import type { Icon } from "../assets/icons/icons";
import type { Role } from "./role";


export type SidebarItem = {
    id: string;
    label: string;
    route?: string;
    icon?: Icon;
    badge?: number;
    children?: SidebarItem[];
    roles: Role[];
};
