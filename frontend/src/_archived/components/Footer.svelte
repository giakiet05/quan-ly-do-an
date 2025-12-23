<script lang="ts">
  import { link } from "svelte-spa-router";

  // Dữ liệu Links giúp HTML gọn gàng hơn
  const footerSections = [
    {
      title: "Product",
      links: [
        { name: "Overview", path: "/overview", highlight: true },
        { name: "Pricing", path: "/pricing" },
        { name: "Customer stories", path: "/stories" },
      ]
    },
    {
      title: "Resources",
      links: [
        { name: "Blog", path: "/blog" },
        { name: "Guides & tutorials", path: "/guides" },
        { name: "Help center", path: "/help" },
      ]
    },
    {
      title: "Company",
      links: [
        { name: "About us", path: "/about" },
        { name: "Careers", path: "/careers" },
        { name: "Media kit", path: "/media" },
      ]
    }
  ];
</script>

<footer class="site-footer">
  <div class="container">
    
    <div class="footer-top">
      
      <div class="footer-col brand-col">
        <div class="logo-wrapper">
          <svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="white"/>
            <path d="M2 17L12 22L22 17" stroke="white"/>
            <path d="M2 12L12 17L22 12" stroke="white"/>
          </svg>
          <span class="brand-name">DoAnHub</span>
        </div>
        <p class="brand-desc">
          whitespace was created for the new ways we live and work. We make a better workspace around the world.
        </p>
      </div>

      {#each footerSections as section}
        <div class="footer-col">
          <h4 class="col-heading">{section.title}</h4>
          <ul class="footer-links">
            {#each section.links as item}
              <li>
                <a 
                  href={item.path} 
                  use:link 
                  class:highlight-link={item.highlight}
                >
                  {item.name}
                </a>
              </li>
            {/each}
          </ul>
        </div>
      {/each}

      <div class="footer-col cta-col">
        <h4 class="col-heading">Try It Today</h4>
        <p class="cta-desc">
          Get started for free. Add your whole team as your needs grow.
        </p>
        <a href="/register" use:link class="btn-start">
          Start today 
          <span class="arrow-icon">&rarr;</span>
        </a>
      </div>
    </div>

    <div class="footer-divider"></div>

    <div class="footer-bottom">
      
      <div class="bottom-left">
        <button class="language-selector">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
          <span>English</span>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 9l6 6 6-6"/></svg>
        </button>
        
        <div class="legal-links">
          <a href="/terms" use:link>Terms & privacy</a>
          <a href="/security" use:link>Security</a>
          <a href="/status" use:link>Status</a>
        </div>
        <span class="copyright">©2021 Whitepace LLC.</span>
      </div>

      <div class="bottom-right social-icons">
        <a href="javascript:void(0)" aria-label="Facebook"><svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"></path></svg></a>
        <a href="javascript:void(0)" aria-label="Twitter"><svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M23 3a10.9 10.9 0 0 1-3.14 1.53 4.48 4.48 0 0 0-7.86 3v1A10.66 10.66 0 0 1 3 4s-4 9 5 13a11.64 11.64 0 0 1-7 2c9 5 20 0 20-11.5a4.5 4.5 0 0 0-.08-.83A7.72 7.72 0 0 0 23 3z"></path></svg></a>
        <a href="javascript:void(0)" aria-label="LinkedIn"><svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"></path><rect x="2" y="9" width="4" height="12"></rect><circle cx="4" cy="4" r="2"></circle></svg></a>
      </div>

    </div>
  </div>
</footer>

<style>
  /* --- 1. Base Styles --- */
  .site-footer {
    /* CHỈNH SỬA Ở ĐÂY: */
    /* Đổi từ var(--primary-color) sang var(--darkblue--) để lấy màu xanh #1a56db */
    background-color: var(--darkblue--); 
    
    color: var(--surface-color); /* Màu chữ trắng */
    font-family: var(--font-primary, 'Inter', sans-serif);
    padding: 60px 20px 32px;
    width: 100%;
    box-sizing: border-box;
    margin-top: auto;
  }

  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  a { text-decoration: none; transition: 0.2s; color: inherit; }

  /* --- 2. Top Section (Grid Layout) --- */
  .footer-top {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 1fr 1.5fr; 
    gap: 40px;
    margin-bottom: 60px;
  }

  /* Logo */
  .logo-wrapper { display: flex; align-items: center; gap: 10px; margin-bottom: 24px; }
  .brand-name { font-size: 24px; font-weight: 700; font-family: var(--font-secondary); }
  .brand-desc { font-size: 14px; line-height: 1.6; opacity: 0.9; max-width: 280px; }

  /* Links */
  .col-heading { font-size: 16px; font-weight: 700; margin-bottom: 20px; font-family: var(--font-secondary); }
  .footer-links { list-style: none; padding: 0; margin: 0; }
  .footer-links li { margin-bottom: 14px; }
  .footer-links a { font-size: 14px; opacity: 0.8; }
  .footer-links a:hover { opacity: 1; text-decoration: underline; }
  
  /* Link highlight: Dùng biến từ app.css */
  .highlight-link { 
    color: var(--text-highlight) !important; 
    opacity: 1 !important; 
  }

  /* CTA Button */
  .cta-desc { font-size: 14px; line-height: 1.5; margin-bottom: 20px; opacity: 0.9; }
  
  .btn-start {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    /* Nút con bên trong: Dùng màu sáng hơn để nổi bật trên nền xanh */
    background-color: var(--cta-blue); 
    color: white;
    padding: 12px 24px;
    border-radius: 6px;
    font-weight: 600;
    font-size: 14px;
    transition: background-color 0.2s;
  }
  .btn-start:hover { 
    background-color: var(--cta-blue-hover); 
  }

  /* Divider */
  .footer-divider {
    height: 1px;
    /* Dùng biến divider từ app.css */
    background-color: var(--footer-divider); 
    margin-bottom: 30px;
  }

  /* --- 3. Bottom Section --- */
  .footer-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    flex-wrap: wrap;
    gap: 20px;
  }

  .bottom-left { display: flex; align-items: center; gap: 30px; flex-wrap: wrap; }
  .legal-links { display: flex; gap: 30px; }
  .bottom-left a, .copyright { opacity: 0.8; }
  .bottom-left a:hover { opacity: 1; }

  /* Language Selector */
  .language-selector {
    display: flex;
    align-items: center;
    gap: 6px;
    background: none;
    border: none;
    color: inherit;
    font: inherit;
    cursor: pointer;
    opacity: 0.9;
    padding: 0;
  }
  .language-selector:hover { opacity: 1; }

  .social-icons { display: flex; gap: 24px; }
  .social-icons a { color: white; opacity: 0.9; }
  .social-icons a:hover { opacity: 1; }

  /* --- 4. Responsive --- */
  @media (max-width: 1024px) {
    .footer-top {
      grid-template-columns: repeat(3, 1fr); 
      gap: 30px;
    }
    .brand-col { grid-column: span 3; max-width: 100%; } 
    .brand-desc { max-width: 600px; }
    .cta-col { grid-column: span 3; margin-top: 10px; } 
  }

  @media (max-width: 600px) {
    .site-footer { padding: 40px 20px 24px; }
    .footer-top {
      grid-template-columns: 1fr;
      text-align: center;
      gap: 32px;
    }
    .brand-col, .cta-col { grid-column: auto; }
    .logo-wrapper, .bottom-left, .social-icons, .legal-links {
      justify-content: center;
    }
    .brand-desc { margin: 0 auto; }
    .footer-bottom {
      flex-direction: column-reverse;
      gap: 24px;
    }
    .bottom-left { flex-direction: column; gap: 16px; }
    .legal-links { flex-direction: column; gap: 12px; }
  }
</style>