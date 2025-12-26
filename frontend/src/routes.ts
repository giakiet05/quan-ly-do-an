import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import ForgotPassword from './pages/ForgotPassword.svelte';
import MyProfile from './pages/MyProfile.svelte';
import Chat from './pages/Chat.svelte';
import MyProjects from './pages/MyProjects.svelte';
import MyProjectDetail from './pages/MyProjectDetail.svelte';
import ClassLecture from './pages/lecture/class/ClassLecture.svelte';
import ClassDetail from './pages/lecture/class/ClassDetail.svelte';
import Notification from './pages/Notification.svelte';
import CreateCategoryModal from './pages/lecture/class/category/CreateCategoryModal.svelte';
import CategoryDetail from './pages/lecture/class/category/CategoryDetail.svelte';

// Archived pages available in src/_archived/pages for reference:
// Home, Profile, Settings, PostDetail, Community, ManageCommunities, ModTools, CreateCommunity

const routes = {
    '/': Login,
    '/auth': Login,
    '/auth/': Login,
    '/auth/login': Login,
    '/auth/register': Register,
    '/auth/forgot-password': ForgotPassword,
    '/home': MyProfile, // Trang chủ sau khi login
    '/my-profile': MyProfile,
    '/chats': Chat,
    '/my-projects': MyProjects,
    '/my-projects/:id': MyProjectDetail,
    '/lecture/my-classes': ClassLecture,
    '/lecture/my-classes/:id': ClassDetail,
    '/lecture/my-classes/:id/categories/create': CreateCategoryModal,
    '/notifications': Notification,
    '/lecture/my-classes/:id/categories/:categoryId': CategoryDetail,
};

export default routes;
