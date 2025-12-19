import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import ForgotPassword from './pages/ForgotPassword.svelte';
import MyProfile from './pages/MyProfile.svelte';

// Archived pages available in src/_archived/pages for reference:
// Home, Profile, Settings, PostDetail, Community, ManageCommunities, ModTools, CreateCommunity

const routes = {
    '/': Login, // Temporary: redirecting to login
    '/auth/login': Login,
    '/auth/register': Register,
    '/auth/forgot-password': ForgotPassword,
    '/my-profile': MyProfile,

};

export default routes;
