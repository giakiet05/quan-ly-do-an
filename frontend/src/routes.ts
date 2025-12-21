import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import ForgotPassword from './pages/ForgotPassword.svelte';
import MyProfile from './pages/MyProfile.svelte';
import Chat from './pages/Chat.svelte';

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

};

export default routes;
