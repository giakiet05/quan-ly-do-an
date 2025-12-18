import Home from './pages/Home.svelte';
import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import ForgotPassword from './pages/ForgotPassword.svelte';
import Profile from './pages/Profile.svelte';
import Settings from './pages/Settings.svelte';
import PostDetail from './pages/PostDetail.svelte';
import Community from './pages/Community.svelte';
import ManageCommunities from './pages/ManageCommunities.svelte';
import ModTools from './pages/ModTools.svelte';

const routes = {
    '/': Home,
    '/auth/login': Login,
    '/auth/register': Register,
    '/auth/forgot-password': ForgotPassword,
    '/profile': Profile,
    '/settings': Settings,
    '/post/:id': PostDetail,
    '/community/:id': Community,
    '/manage-communities': ManageCommunities,
    '/mod-tools': ModTools
};

export default routes;
