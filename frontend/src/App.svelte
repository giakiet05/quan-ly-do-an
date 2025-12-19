<script>
  import Router from "svelte-spa-router";
  import { location } from "svelte-spa-router";

  import MainLayout from "./layouts/MainLayout.svelte";
  import AuthLayout from "./layouts/AuthLayout.svelte";
  import Login from "./pages/Login.svelte";
  import routes from "./routes";
</script>

{#if $location.startsWith("/auth")}
  <AuthLayout>
    <Router {routes} />
  </AuthLayout>
{:else}
  <MainLayout>
    <Router {routes} />
  </MainLayout>
{/if}
