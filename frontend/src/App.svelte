<script>
  import Router from "svelte-spa-router";
  import { location, push } from "svelte-spa-router";
  import { onMount } from "svelte";

  import MainLayout from "./layouts/MainLayout.svelte";
  import AuthLayout from "./layouts/AuthLayout.svelte";
  import Login from "./pages/Login.svelte";
  import routes from "./routes";

  // Redirect root to auth/login
  onMount(() => {
    if ($location === "/") {
      push("/auth/login");
    }
  });
</script>

{#if $location.startsWith("/auth") || $location === "/"}
  <AuthLayout>
    <Router {routes} />
  </AuthLayout>
{:else}
  <MainLayout>
    <Router {routes} />
  </MainLayout>
{/if}
