import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import ForgotPassword from './pages/ForgotPassword.svelte';
import GoogleCallBack from './pages/GoogleCallBack.svelte';
import GoogleSetUp from './pages/GoogleSetUp.svelte';
import GoogleError from './pages/GoogleError.svelte';
import MyProfile from './pages/MyProfile.svelte';
import Chat from './pages/Chat.svelte';
import MyProjects from './pages/MyProjects.svelte';
import MyProjectDetail from './pages/MyProjectDetail.svelte';
import ClassLecture from './pages/lecture/class/ClassLecture.svelte';
import ClassDetail from './pages/lecture/class/ClassDetail.svelte';
import Notification from './pages/Notification.svelte';
import CreateCategoryModal from './pages/lecture/class/category/CreateCategoryModal.svelte';
import CategoryDetail from './pages/lecture/class/category/CategoryDetail.svelte';
import StudentsManagementView from './pages/lecture/student/StudentsManagementView.svelte';
import ProjectDetailView from './pages/lecture/class/category/project/ProjectDetailView.svelte';

// Student pages
import StudentClassList from './pages/student/StudentClassList.svelte';
import StudentClassDetail from './pages/student/StudentClassDetail.svelte';
import StudentCategoryDetail from './pages/student/StudentCategoryDetail.svelte';
import StudentProjectRegister from './pages/student/StudentProjectRegister.svelte';
import JoinClassroom from './pages/student/JoinClassroom.svelte';

// Archived pages available in src/_archived/pages for reference:
// Home, Profile, Settings, PostDetail, Community, ManageCommunities, ModTools, CreateCommunity

const routes = {
    '/': Login,
    '/auth': Login,
    '/auth/': Login,
    '/auth/login': Login,
    '/auth/register': Register,
    '/auth/forgot-password': ForgotPassword,
    '/auth/callback': GoogleCallBack,
    '/auth/google-setup': GoogleSetUp,
    "/auth/error": GoogleError,
    '/home': MyProfile, // Trang chủ sau khi login
    '/my-profile': MyProfile,
    '/chats': Chat,
    '/my-projects': MyProjects,
    '/my-projects/:id': MyProjectDetail,
    '/notifications': Notification,

    // Lecturer routes
    '/lecture/my-classes': ClassLecture,
    '/lecture/my-classes/:id': ClassDetail,
    '/lecture/my-classes/:id/categories/create': CreateCategoryModal,
    '/lecture/my-classes/:id/categories/:categoryId': CategoryDetail,
    '/lecture/my-classes/:id/categories/:categoryId/projects/:projectId': ProjectDetailView,
    '/lecture/students-management': StudentsManagementView,

    // Student routes
    '/student/classes': StudentClassList,
    '/student/classes/join': JoinClassroom,
    '/student/classes/:id': StudentClassDetail,
    '/student/classes/:id/categories/:categoryId': StudentCategoryDetail,
    '/student/classes/:id/categories/:categoryId/projects/:projectId': StudentProjectRegister,

};

export default routes;
