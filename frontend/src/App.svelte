<script>
  import Router from "svelte-spa-router";
  import { location, push } from "svelte-spa-router";
  import { onMount } from "svelte";

  import MainLayout from "./layouts/MainLayout.svelte";
  import AuthLayout from "./layouts/AuthLayout.svelte";
  import Login from "./pages/auth/Login.svelte";
  import routes from "./routes";

  onMount(() => {
    if ($location === "/") {
      push("/home");
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
