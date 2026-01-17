<script lang="ts">
  import { onMount } from "svelte";
  import { replace } from "svelte-spa-router";
  import { authService } from "../services/auth-service";

  onMount(() => {
    const hash = window.location.hash;
    const queryString = hash.includes("?") ? hash.split("?")[1] : "";
    const params = new URLSearchParams(queryString);

    const accessToken = params.get("access_token");
    const setupToken = params.get("setup_token");

    console.log("Callback nhận được:", { accessToken, setupToken });

    if (setupToken) {
      replace(`/auth/google-setup?setup_token=${setupToken}`);
      return;
    }

    if (accessToken) {
      const result = authService.handleLoginCallback();
      console.log("Login callback result:", result);

      if (result.success) {
        setTimeout(() => {
          replace("/home");
        }, 100);
      } else {
        console.error("Login failed:", result);
        replace("/auth/login");
      }
      return;
    }
    replace("/auth/login");
  });
</script>

<div style="text-align: center; margin-top: 50px; font-family: sans-serif;">
  <p>Đang xử lý đăng nhập Google...</p>
</div>
