<script>
  import Router from "svelte-spa-router";
  import { location, push } from "svelte-spa-router";
  import { onMount } from "svelte";

  import MainLayout from "./layouts/MainLayout.svelte";
  import AuthLayout from "./layouts/AuthLayout.svelte";
  import routes from "./routes";

  onMount(() => {
    if ($location === "/" || $location === "") {
      push("/home");
    }
  });
</script>

{#if $location.startsWith("/auth") || $location === "/"}
  <AuthLayout>
    {#key $location}
      <Router {routes} />
    {/key}
  </AuthLayout>
{:else}
  <MainLayout>
    {#key $location}
      <Router {routes} />
    {/key}
  </MainLayout>
{/if}
