import { User, Settings, Search, LogOut, LayoutDashboard, BookOpen, Bell, Users, MessageCircle, ArrowRightLeft } from 'lucide-svelte';


export type Icon = typeof User;

export const Icons = {
    user: User,
    settings: Settings,
    search: Search,
    logout: LogOut,
    dashboard: LayoutDashboard,
    class: BookOpen,
    notification: Bell,
    students: Users,
    message: MessageCircle,
    resize: ArrowRightLeft
};
