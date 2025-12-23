<script lang="ts">
  import { createEventDispatcher, onMount, onDestroy } from "svelte";
  import Button from "./Button.svelte";
  import { sendEmailVerification } from "../services/auth-service";

  const dispatch = createEventDispatcher();
  export let email: string = "";

  let values = ["", "", "", "", "", ""];
  let inputRefs: HTMLInputElement[] = [];
  let timeLeft = 60;
  let timerInterval: any;
  let loading = false;
<<<<<<< HEAD
  export let error = ""; 
=======
  let error = "";
>>>>>>> dev

  const startTimer = () => {
    clearInterval(timerInterval);
    timeLeft = 60;
    timerInterval = setInterval(() => {
      if (timeLeft > 0) {
        timeLeft--;
      } else {
        clearInterval(timerInterval);
      }
    }, 1000);
  };

  onMount(() => {
    startTimer();
    if (inputRefs[0]) inputRefs[0].focus();
  });

  onDestroy(() => {
    clearInterval(timerInterval);
  });

  // --- HANDLERS ---
  const handleInput = (index: number, e: Event) => {
    error = "";
    const target = e.target as HTMLInputElement;
    const val = target.value;

    // Chỉ cho nhập số
    if (!/^\d*$/.test(val)) {
      values[index] = "";
      return;
    }

    // Lấy số cuối cùng
    values[index] = val.substring(val.length - 1);

    // Tự động focus ô tiếp theo
    if (val && index < 5) {
      // index < 5 vì mảng có 6 phần tử (0-5)
      inputRefs[index + 1].focus();
    }
  };

  const handleKeyDown = (index: number, e: KeyboardEvent) => {
    if (e.key === "Backspace") {
      if (!values[index] && index > 0) {
        inputRefs[index - 1].focus();
      }
    }
  };

  const handlePaste = (e: ClipboardEvent) => {
    e.preventDefault();
    const clipboardData = e.clipboardData;
    if (!clipboardData) return;

    const pasteData = clipboardData.getData("text").trim();
    if (!/^\d+$/.test(pasteData)) return;

    const chars = pasteData.split("").slice(0, 6); // Cắt lấy 6 ký tự
    chars.forEach((char, i) => {
      values[i] = char;
    });

    // Focus vào ô cuối cùng vừa paste
    const lastIndex = chars.length - 1;
    if (lastIndex < 6 && inputRefs[lastIndex]) {
      inputRefs[lastIndex].focus();
    }
  };

  const handleVerify = () => {
    const token = values.join("");
    if (token.length < 6) return;
    
    dispatch("success", { verification_token: token }); 
  };

  const handleResend = async () => {
    if (timeLeft > 0) return; // Chặn nếu chưa hết giờ

    loading = true;
    try {
      await sendEmailVerification({ email });
      startTimer();
      values = ["", "", "", "", "", ""];
      inputRefs[0].focus();
      error = "";
    } catch (e) {
      error = "Không thể gửi lại mã. Vui lòng thử lại.";
    } finally {
      loading = false;
    }
  };

  $: formattedTime = `00:${timeLeft < 10 ? "0" : ""}${timeLeft}`;
  $: isComplete = values.every((v) => v !== "");
</script>

<div class="otp-container">
  <div class="auth-header">
    <div class="icon-wrapper">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="currentColor"
        class="shield-icon"
      >
        <path
          fill-rule="evenodd"
          d="M12.516 2.17a.75.75 0 0 0-1.032 0 11.209 11.209 0 0 1-7.877 3.08.75.75 0 0 0-.722.515A12.74 12.74 0 0 0 2.25 9.75c0 5.942 4.064 10.933 9.563 12.348a.749.749 0 0 0 .374 0c5.499-1.415 9.563-6.406 9.563-12.348 0-1.352-.272-2.636-.759-3.803a.75.75 0 0 0-.722-.515 11.208 11.208 0 0 1-7.877-3.08ZM12 13.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3Z"
          clip-rule="evenodd"
        />
      </svg>
    </div>
    <h2 class="title">Xác thực OTP</h2>
    <p class="subtitle">Nhập mã 6 số vừa được gửi đến email của bạn.</p>
  </div>

  <div class="otp-inputs">
    {#each values as val, i}
      <input
        type="text"
        inputmode="numeric"
        maxlength="1"
        class="otp-box"
        bind:value={values[i]}
        bind:this={inputRefs[i]}
        on:input={(e) => handleInput(i, e)}
        on:keydown={(e) => handleKeyDown(i, e)}
        on:paste={i === 0 ? handlePaste : null}
        disabled={loading}
      />
    {/each}
  </div>

  {#if error}
    <div class="error-message">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="error-icon"
      >
        <path
          fill-rule="evenodd"
          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
          clip-rule="evenodd"
        />
      </svg>
      <span>{error}</span>
    </div>
  {/if}

  <div class="timer-wrapper">
    {#if timeLeft > 0}
      <span class="timer-text"
        >Mã hết hạn sau: <span class="time">{formattedTime}</span></span
      >
    {:else}
      <span class="timer-expired">Mã đã hết hạn</span>
    {/if}
  </div>

  <Button
    variant="primary"
    onclick={handleVerify}
    disabled={loading || !isComplete}
    class="w-full"
  >
    {#if loading}
      <span class="loader"></span> Đang kiểm tra...
    {:else}
      Xác nhận
    {/if}
  </Button>

  <div class="footer-text">
    Bạn không nhận được mã?
    <button
      class="resend-btn"
      on:click={handleResend}
      disabled={timeLeft > 0 || loading}
    >
      Gửi lại
    </button>
  </div>
</div>

<style>
  :root {
    --primary-blue: #1a56db;
    --text-dark: #111827;
    --text-gray: #6b7280;
  }

  /* Header style giống Register */
  .auth-header {
    margin-bottom: 28px;
    text-align: center;
  }
  .icon-wrapper {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 16px;
  }
  .shield-icon {
    width: 48px;
    height: 48px;
    color: var(--primary-blue);
  }
  .title {
    font-size: 28px;
    font-weight: 700;
    color: var(--text-dark);
    margin: 0 0 8px 0;
  }
  .subtitle {
    font-size: 16px;
    color: var(--text-gray);
    margin: 0;
  }

  /* Input Style */
  .otp-inputs {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-bottom: 24px;
    width: 100%;
  }

  .otp-box {
    width: 48px;
    height: 56px;
    border: 2px solid #e5e7eb;
    border-radius: 10px;
    text-align: center;
    font-size: 22px;
    font-weight: 700;
    color: var(--text-dark);
    background-color: #f9fafb;
    transition: all 0.2s;
    outline: none;
    padding: 0;
  }

  .otp-box:focus {
    border-color: var(--primary-blue);
    background-color: #fff;
    box-shadow: 0 0 0 4px rgba(26, 86, 219, 0.1);
  }

  /* Timer */
  .timer-wrapper {
    text-align: center;
    margin-bottom: 24px;
    font-size: 14px;
  }
  .time {
    color: #ef4444;
    font-weight: 600;
  }
  .timer-expired {
    color: #ef4444;
    font-weight: 600;
  }

  /* Utilities */
  .w-full {
    width: 100% !important;
    display: flex !important;
  }

  .loader {
    width: 16px;
    height: 16px;
    border: 2px solid #fff;
    border-bottom-color: transparent;
    border-radius: 50%;
    display: inline-block;
    animation: rotation 1s linear infinite;
    margin-right: 8px;
  }
  @keyframes rotation {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  /* Footer */
  .footer-text {
    margin-top: 24px;
    font-size: 14px;
    color: var(--text-dark);
    text-align: center;
  }

  .resend-btn {
    background: none;
    border: none;
    color: var(--primary-blue);
    font-weight: 700;
    cursor: pointer;
    padding: 0;
    margin-left: 4px;
    font-size: 14px;
  }

  .resend-btn:disabled {
    color: #9ca3af;
    cursor: not-allowed;
    text-decoration: none;
  }
  .resend-btn:hover:not(:disabled) {
    text-decoration: underline;
  }

  .error-message {
    background-color: #fef2f2;
    color: #991b1b;
    font-size: 14px;
    padding: 10px 12px;
    border-radius: 8px;
    margin-bottom: 16px; /* Cách timer ra một chút */
    display: flex;
    align-items: center;
    justify-content: center; /* Căn giữa nội dung */
    gap: 8px;
    border: 1px solid #fee2e2;
    animation: slideDown 0.3s ease-out;
  }

  .error-icon {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
