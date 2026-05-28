'use client'

import { useState, useEffect, useRef } from 'react'

/* ─────────────────────────────────────
   STYLES
───────────────────────────────────── */
const STYLES = `
@import url('https://fonts.googleapis.com/css2?family=Fraunces:ital,opsz,wght@0,9..144,300;0,9..144,600;0,9..144,700;1,9..144,300;1,9..144,600&family=Syne:wght@400;500;600;700;800&display=swap');

.kl {
  font-family:'Syne',system-ui,sans-serif;
  background:#07070f; color:#ededf7; min-height:100vh;
  --bg:#07070f; --bg2:#0d0d1c; --bg3:#111124;
  --border:rgba(255,255,255,.065); --border2:rgba(255,255,255,.11); --border3:rgba(255,255,255,.22);
  --text:#ededf7; --text2:#7e7ea0; --text3:#3e3e58;
  --accent:#6d28e8; --accent-lt:#a78bfa;
  --accent-faint:rgba(109,40,232,.10); --accent-glow:rgba(109,40,232,.30);
  --green:#34d399;
}
.kl *,.kl *::before,.kl *::after{box-sizing:border-box}
.kl .grain{position:fixed;inset:0;pointer-events:none;z-index:999;opacity:.022;
  background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='300' height='300'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='300' height='300' filter='url(%23n)'/%3E%3C/svg%3E")}
.kl .container{max-width:1120px;margin:0 auto;padding:0 32px}
.kl .container-sm{max-width:760px;margin:0 auto;padding:0 32px}
.kl .display{font-family:'Fraunces',Georgia,serif}

/* keyframes */
@keyframes kl-up{from{opacity:0;transform:translateY(22px)}to{opacity:1;transform:translateY(0)}}
@keyframes kl-blink{0%,100%{opacity:1}50%{opacity:.35}}
@keyframes kl-fa{0%,100%{transform:translateY(0)}50%{transform:translateY(-10px)}}
@keyframes kl-fb{0%,100%{transform:translateY(0)}50%{transform:translateY(-14px)}}
@keyframes kl-fc{0%,100%{transform:translateY(0)}50%{transform:translateY(-8px)}}
@keyframes kl-persp{from{opacity:0;transform:perspective(1200px) rotateY(-18deg) rotateX(8deg) translateY(32px) scale(.97)}to{opacity:1;transform:perspective(1200px) rotateY(-8deg) rotateX(3deg) translateY(0) scale(1)}}
@keyframes kl-spin{from{transform:rotate(0deg)}to{transform:rotate(360deg)}}
@keyframes kl-ring{0%,100%{opacity:.18;transform:scale(1)}50%{opacity:.45;transform:scale(1.07)}}
@keyframes kl-ring2{0%,100%{opacity:.1;transform:scale(1)}50%{opacity:.3;transform:scale(1.12)}}
@keyframes kl-b1{0%,100%{height:4px}50%{height:26px}}
@keyframes kl-b2{0%,100%{height:8px}50%{height:22px}}
@keyframes kl-b3{0%,100%{height:12px}50%{height:32px}}
@keyframes kl-b4{0%,100%{height:6px}50%{height:24px}}
@keyframes kl-b5{0%,100%{height:10px}50%{height:18px}}
@keyframes kl-b6{0%,100%{height:4px}50%{height:20px}}
@keyframes kl-b7{0%,100%{height:14px}50%{height:30px}}
@keyframes kl-prog{from{width:0}to{width:72%}}

.kl .au{animation:kl-up .65s cubic-bezier(.16,1,.3,1) forwards}
.kl .d1{animation-delay:.05s;opacity:0}.kl .d2{animation-delay:.15s;opacity:0}
.kl .d3{animation-delay:.25s;opacity:0}.kl .d4{animation-delay:.35s;opacity:0}
.kl .d5{animation-delay:.45s;opacity:0}.kl .d6{animation-delay:.60s;opacity:0}
.kl .d7{animation-delay:.75s;opacity:0}
[data-reveal]{opacity:0;transform:translateY(30px);transition:opacity .7s cubic-bezier(.16,1,.3,1),transform .7s cubic-bezier(.16,1,.3,1)}
[data-reveal].in-view{opacity:1;transform:translateY(0)}
[data-reveal][data-delay="2"].in-view{transition-delay:.16s}
[data-reveal][data-delay="3"].in-view{transition-delay:.24s}

/* buttons */
.kl .btn-primary{display:inline-flex;align-items:center;gap:7px;background:var(--accent);color:#fff;font-family:'Syne',sans-serif;font-weight:600;font-size:14px;padding:10px 20px;border-radius:8px;border:none;cursor:pointer;text-decoration:none;letter-spacing:.01em;transition:background .2s,transform .15s,box-shadow .2s}
.kl .btn-primary:hover{background:#5b21b6;transform:translateY(-1px);box-shadow:0 8px 28px var(--accent-glow)}
.kl .btn-primary.lg{font-size:15px;padding:13px 26px;border-radius:9px}
.kl .btn-ghost{display:inline-flex;align-items:center;gap:7px;background:transparent;color:var(--text2);font-family:'Syne',sans-serif;font-weight:500;font-size:14px;padding:10px 16px;border-radius:8px;border:1px solid var(--border2);cursor:pointer;text-decoration:none;transition:all .2s}
.kl .btn-ghost:hover{color:var(--text);border-color:var(--border3);background:rgba(255,255,255,.04)}
.kl .btn-ghost.lg{font-size:15px;padding:13px 20px}

/* navbar */
.kl nav{position:fixed;top:0;left:0;right:0;z-index:100;border-bottom:1px solid var(--border);background:rgba(7,7,15,.85);backdrop-filter:blur(20px)}
.kl .nav-inner{display:flex;align-items:center;justify-content:space-between;height:60px;max-width:1120px;margin:0 auto;padding:0 32px}
.kl .nav-logo{display:flex;align-items:center;gap:9px;font-weight:800;font-size:17px;letter-spacing:-.02em;color:var(--text);text-decoration:none}
.kl .nav-links{display:flex;align-items:center;gap:30px;list-style:none;margin:0;padding:0}
.kl .nav-links a{font-size:13.5px;font-weight:500;color:var(--text2);text-decoration:none;transition:color .2s}
.kl .nav-links a:hover{color:var(--text)}
.kl .nav-right{display:flex;align-items:center;gap:10px}

/* hero */
.kl .hero{padding-top:60px;min-height:100vh;display:flex;align-items:center;position:relative;overflow:hidden}
.kl .hero-dots{position:absolute;inset:0;pointer-events:none;background-image:radial-gradient(circle at 1px 1px,rgba(255,255,255,.048) 1px,transparent 0);background-size:36px 36px;mask-image:radial-gradient(ellipse 90% 90% at 25% 50%,black 20%,transparent 75%)}
.kl .hero-glow-tr{position:absolute;top:-100px;right:-100px;width:700px;height:700px;pointer-events:none;background:radial-gradient(ellipse at top right,rgba(109,40,232,.16) 0%,transparent 55%)}
.kl .hero-glow-bl{position:absolute;bottom:-80px;left:-80px;width:480px;height:480px;pointer-events:none;background:radial-gradient(ellipse at bottom left,rgba(52,211,153,.07) 0%,transparent 55%)}
.kl .hero-split{display:grid;grid-template-columns:5fr 7fr;gap:56px;align-items:center;padding:80px 0;position:relative;z-index:1;width:100%}
.kl .badge{display:inline-flex;align-items:center;gap:8px;border:1px solid rgba(167,139,250,.22);border-radius:100px;padding:5px 15px;font-size:11.5px;font-weight:600;letter-spacing:.08em;text-transform:uppercase;color:var(--accent-lt);background:rgba(109,40,232,.08);margin-bottom:24px}
.kl .badge-dot{width:5px;height:5px;border-radius:50%;background:var(--green);animation:kl-blink 2s ease-in-out infinite}
.kl .hero-h1{font-family:'Fraunces',Georgia,serif;font-size:clamp(40px,4.8vw,66px);font-weight:700;line-height:1.04;letter-spacing:-.03em;margin-bottom:20px;background:linear-gradient(180deg,#fff 30%,rgba(200,200,220,.65) 100%);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.kl .hero-h1 em{font-style:italic;-webkit-text-fill-color:var(--accent-lt)}
.kl .hero-sub{font-size:17px;color:var(--text2);line-height:1.7;max-width:420px;margin-bottom:32px}
.kl .hero-actions{display:flex;align-items:center;gap:10px;flex-wrap:wrap;margin-bottom:12px}
.kl .hero-disclaimer{font-size:12px;color:var(--text3);margin-bottom:36px}
.kl .hero-stats{display:flex;align-items:center;gap:24px;padding-top:28px;border-top:1px solid var(--border)}
.kl .hero-stat-n{display:block;font-family:'Fraunces',serif;font-size:28px;font-weight:700;letter-spacing:-.02em}
.kl .hero-stat-l{display:block;font-size:12px;color:var(--text3);margin-top:2px}
.kl .hero-stat-sep{width:1px;height:36px;background:var(--border2)}
.kl .hero-visual{position:relative}
.kl .hero-vglow{position:absolute;inset:-80px;pointer-events:none;background:radial-gradient(ellipse at 55% 40%,rgba(109,40,232,.2) 0%,transparent 55%)}
.kl .hero-tilt{position:relative;z-index:1;animation:kl-persp .9s cubic-bezier(.16,1,.3,1) .5s both;transform-origin:left center;transition:transform .5s ease}
.kl .hero-tilt:hover{transform:perspective(1200px) rotateY(-4deg) rotateX(1deg) !important}
.kl .fbadge{position:absolute;z-index:10;display:flex;align-items:center;gap:10px;background:rgba(8,8,22,.94);border:1px solid var(--border2);border-radius:12px;padding:10px 14px;backdrop-filter:blur(16px);box-shadow:0 8px 32px rgba(0,0,0,.55);white-space:nowrap}
.kl .fbadge-icon{font-size:19px}
.kl .fbadge-label{font-size:10.5px;font-weight:600;color:var(--text3);text-transform:uppercase;letter-spacing:.06em;margin-bottom:2px}
.kl .fbadge-val{font-family:'Fraunces',serif;font-size:18px;font-weight:700;line-height:1.1}
.kl .fb1{animation:kl-fa 4s ease-in-out infinite}
.kl .fb2{animation:kl-fb 5.2s ease-in-out infinite 1.1s}
.kl .fb3{animation:kl-fc 4.6s ease-in-out infinite .6s}

/* mini browser window */
.kl .mwin{border:1px solid var(--border2);border-radius:14px;overflow:hidden;background:#0b0b1e;box-shadow:0 0 0 1px rgba(255,255,255,.03),0 40px 80px rgba(0,0,0,.7),0 0 60px rgba(109,40,232,.08);position:relative}
.kl .mwin-bar{display:flex;align-items:center;gap:7px;padding:11px 16px;border-bottom:1px solid var(--border);background:rgba(0,0,0,.3)}
.kl .mwin-dot{width:10px;height:10px;border-radius:50%}
.kl .mwin-url{flex:1;margin:0 14px;height:20px;border-radius:4px;background:rgba(255,255,255,.04);border:1px solid var(--border);display:flex;align-items:center;padding:0 10px;font-size:10px;color:var(--text3)}
.kl .mwin-body{display:grid;grid-template-columns:175px 1fr}
.kl .mwin-sb{border-right:1px solid var(--border);padding:14px 12px;background:rgba(0,0,0,.2)}
.kl .mwin-sb-sec{font-size:9.5px;font-weight:700;letter-spacing:.1em;text-transform:uppercase;color:var(--text3);margin:12px 0 5px 8px}
.kl .mwin-sb-sec:first-child{margin-top:0}
.kl .mwin-sb-item{display:flex;align-items:center;gap:8px;padding:6px 8px;border-radius:5px;font-size:11px;color:var(--text3);margin-bottom:1px}
.kl .mwin-sb-item.on{background:rgba(109,40,232,.14);color:var(--accent-lt);border:1px solid rgba(109,40,232,.2)}
.kl .mwin-main{padding:18px 22px}
.kl .mwin-title{font-size:16px;font-weight:700;margin-bottom:2px}
.kl .mwin-sub{font-size:10px;color:var(--text3);margin-bottom:14px}
.kl .mwin-pill{display:inline-flex;align-items:center;gap:12px;border:1px solid var(--border);border-radius:6px;padding:6px 11px;font-size:10px;color:var(--text2);background:rgba(255,255,255,.02);margin-bottom:14px}
.kl .mwin-pill b{color:var(--green)}.kl .mwin-pill span{color:var(--accent-lt);font-weight:700}
.kl .mwin-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.kl .mwin-card{border:1px solid var(--border);border-radius:8px;padding:10px;background:rgba(255,255,255,.015)}
.kl .mwin-clabel{font-size:10px;font-weight:600;margin-bottom:8px}
.kl .mwin-inp{width:100%;background:rgba(255,255,255,.06);border:1px solid var(--border);border-radius:5px;padding:6px 8px;font-size:10px;color:var(--text2);font-family:'Syne',sans-serif;margin-bottom:7px;display:block}
.kl .mwin-gbtn{width:100%;background:var(--accent);border:none;border-radius:5px;padding:7px;font-size:10px;font-weight:700;color:#fff;font-family:'Syne',sans-serif}
.kl .mwin-ok{font-size:9px;color:var(--green);margin-top:5px}
.kl .mwin-modes{display:flex;gap:3px;margin-bottom:8px}
.kl .mwin-mode{flex:1;padding:5px 2px;border-radius:4px;text-align:center;font-size:9.5px;font-weight:600;background:rgba(255,255,255,.04);color:var(--text3)}
.kl .mwin-mode.on{background:var(--accent);color:#fff}
.kl .mwin-tgt{font-size:9px;color:var(--text3);margin-bottom:6px}
.kl .mwin-pbtn{width:100%;background:var(--accent);border:none;border-radius:5px;padding:7px;font-size:9.5px;font-weight:700;color:#fff;font-family:'Syne',sans-serif}
.kl .mwin-fade{position:absolute;bottom:0;left:0;right:0;height:64px;background:linear-gradient(to top,#07070f,transparent);pointer-events:none}

/* compat bar */
.kl .compat{padding:24px 0;border-top:1px solid var(--border);border-bottom:1px solid var(--border)}
.kl .compat-inner{max-width:1120px;margin:0 auto;padding:0 32px;display:flex;align-items:center;justify-content:center;gap:10px;flex-wrap:wrap}
.kl .compat-label{font-size:12px;color:var(--text3);font-weight:500;margin-right:4px}
.kl .compat-chip{padding:5px 14px;border:1px solid var(--border);border-radius:7px;font-size:12px;font-weight:600;color:var(--text3);transition:all .2s}
.kl .compat-chip:hover{border-color:var(--border2);color:var(--text2)}

/* shared section */
.kl section{padding:96px 0}
.kl .section-label{display:inline-flex;align-items:center;gap:10px;font-size:11px;font-weight:700;letter-spacing:.1em;text-transform:uppercase;color:var(--accent-lt);margin-bottom:18px}
.kl .section-label::before{content:'';width:16px;height:1px;background:var(--accent-lt)}
.kl .section-title{font-family:'Fraunces',Georgia,serif;font-size:clamp(30px,4vw,48px);font-weight:700;letter-spacing:-.025em;line-height:1.1;margin-bottom:16px}
.kl .section-desc{font-size:17px;color:var(--text2);line-height:1.65}

/* ══ 1. CONTAINER SCROLL (Aceternity-inspired) ══ */
.kl .cs-section{border-top:1px solid var(--border);padding:0}
.kl .cs-outer{position:relative}
.kl .cs-sticky{position:sticky;top:0;display:flex;flex-direction:column;align-items:center;justify-content:center;min-height:100vh;padding:56px 32px;overflow:hidden}
.kl .cs-bg-glow{position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);width:900px;height:600px;pointer-events:none;background:radial-gradient(ellipse at center,rgba(109,40,232,.12) 0%,transparent 60%)}
.kl .cs-header{text-align:center;margin-bottom:48px;position:relative;z-index:1}
.kl .cs-label{display:inline-flex;align-items:center;gap:10px;font-size:11px;font-weight:700;letter-spacing:.1em;text-transform:uppercase;color:var(--accent-lt);margin-bottom:18px}
.kl .cs-label::before{content:'';width:16px;height:1px;background:var(--accent-lt)}
.kl .cs-title{font-family:'Fraunces',Georgia,serif;font-size:clamp(28px,3.8vw,46px);font-weight:700;letter-spacing:-.025em;line-height:1.1;margin-bottom:14px}
.kl .cs-desc{font-size:16px;color:var(--text2);max-width:480px;margin:0 auto;line-height:1.65}
.kl .cs-wrap{width:100%;max-width:900px;position:relative;z-index:1;transform-origin:top center;will-change:transform}
.kl .cs-win{border-radius:16px;overflow:hidden;border:1px solid var(--border2);background:#0b0b1e}
.kl .cs-win-bar{display:flex;align-items:center;gap:7px;padding:12px 18px;border-bottom:1px solid var(--border);background:rgba(0,0,0,.3)}
.kl .cs-win-dot{width:10px;height:10px;border-radius:50%}
.kl .cs-win-tabs{display:flex;gap:1px;margin:0 16px}
.kl .cs-win-tab{font-size:11px;padding:4px 14px;border-radius:4px;color:var(--text3);background:rgba(255,255,255,.04);border:1px solid var(--border)}
.kl .cs-win-tab.active{color:var(--text);background:rgba(109,40,232,.12);border-color:rgba(109,40,232,.2)}
.kl .cs-strip{display:flex;align-items:center;gap:14px;padding:12px 40px;border-bottom:1px solid var(--border);background:rgba(0,0,0,.2);font-size:12px;color:var(--text2);flex-wrap:wrap}
.kl .cs-strip-score{font-weight:700;color:var(--green)}
.kl .cs-strip-dot{width:3px;height:3px;border-radius:50%;background:var(--border2);flex-shrink:0}
.kl .cs-article{padding:36px 52px 56px;max-height:380px;overflow:hidden;position:relative}
.kl .cs-article-tag{display:inline-block;font-size:11px;font-weight:700;letter-spacing:.08em;text-transform:uppercase;color:var(--accent-lt);background:var(--accent-faint);border:1px solid rgba(167,139,250,.2);border-radius:5px;padding:4px 12px;margin-bottom:20px}
.kl .cs-article h1{font-family:'Fraunces',Georgia,serif;font-size:clamp(18px,2.6vw,30px);font-weight:700;letter-spacing:-.025em;line-height:1.15;margin-bottom:14px}
.kl .cs-article-meta{display:flex;align-items:center;gap:12px;font-size:12px;color:var(--text3);margin-bottom:22px;padding-bottom:18px;border-bottom:1px solid var(--border)}
.kl .cs-article p{font-size:15px;color:var(--text2);line-height:1.75;margin-bottom:16px}
.kl .cs-article h2{font-family:'Fraunces',Georgia,serif;font-size:18px;font-weight:600;margin-bottom:10px;color:var(--text)}
.kl .cs-fade{position:absolute;bottom:0;left:0;right:0;height:130px;background:linear-gradient(to top,#07070f,rgba(7,7,15,.8),transparent);pointer-events:none}

/* ══ 2. ARTICLE ENGINE (music-player aesthetic) ══ */
.kl .ae-section{border-top:1px solid var(--border)}
.kl .ae-inner{display:grid;grid-template-columns:1fr 1fr;gap:56px;align-items:center}
.kl .ae-card{background:var(--bg2);border:1px solid var(--border2);border-radius:24px;padding:28px;position:relative;overflow:hidden}
.kl .ae-glow{position:absolute;top:-40px;left:50%;transform:translateX(-50%);width:300px;height:200px;pointer-events:none;background:radial-gradient(ellipse at top,rgba(109,40,232,.22) 0%,transparent 65%)}
.kl .ae-top{display:flex;align-items:center;justify-content:space-between;margin-bottom:24px}
.kl .ae-status{display:flex;align-items:center;gap:7px;font-size:11px;font-weight:700;letter-spacing:.07em;text-transform:uppercase;color:var(--green)}
.kl .ae-dot{width:6px;height:6px;border-radius:50%;background:var(--green);animation:kl-blink 1.5s ease-in-out infinite}
.kl .ae-time{font-size:12px;color:var(--text3);font-weight:600}
.kl .ae-main{display:flex;gap:20px;align-items:center;margin-bottom:22px}
.kl .ae-disc-wrap{position:relative;width:80px;height:80px;flex-shrink:0}
.kl .ae-disc{position:absolute;inset:0;border-radius:50%;background:linear-gradient(135deg,var(--accent) 0%,#1e1b4b 55%,var(--bg3) 100%);display:flex;align-items:center;justify-content:center;animation:kl-spin 12s linear infinite;box-shadow:0 0 32px var(--accent-glow)}
.kl .ae-hole{width:24px;height:24px;border-radius:50%;background:var(--bg2);display:flex;align-items:center;justify-content:center;font-size:11px;position:relative;z-index:1}
.kl .ae-ring1{position:absolute;inset:-10px;border-radius:50%;border:1.5px solid rgba(109,40,232,.3);animation:kl-ring 3s ease-in-out infinite;pointer-events:none}
.kl .ae-ring2{position:absolute;inset:-22px;border-radius:50%;border:1px solid rgba(109,40,232,.15);animation:kl-ring2 4s ease-in-out infinite 1s;pointer-events:none}
.kl .ae-info{flex:1;min-width:0}
.kl .ae-atitle{font-family:'Fraunces',serif;font-size:17px;font-weight:700;line-height:1.3;margin-bottom:5px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.kl .ae-site{font-size:12px;color:var(--text3);margin-bottom:14px}
.kl .ae-bars{display:flex;align-items:flex-end;gap:3px;height:36px}
.kl .ae-bar{width:4px;border-radius:2px 2px 0 0;background:linear-gradient(to top,var(--accent),var(--accent-lt));min-height:4px}
.kl .ae-bar:nth-child(1){animation:kl-b1 .9s ease-in-out infinite 0.00s}
.kl .ae-bar:nth-child(2){animation:kl-b2 .9s ease-in-out infinite 0.10s}
.kl .ae-bar:nth-child(3){animation:kl-b3 .9s ease-in-out infinite 0.20s}
.kl .ae-bar:nth-child(4){animation:kl-b4 .9s ease-in-out infinite 0.30s}
.kl .ae-bar:nth-child(5){animation:kl-b5 .9s ease-in-out infinite 0.15s}
.kl .ae-bar:nth-child(6){animation:kl-b6 .9s ease-in-out infinite 0.25s}
.kl .ae-bar:nth-child(7){animation:kl-b7 .9s ease-in-out infinite 0.05s}
.kl .ae-bar:nth-child(8){animation:kl-b1 .9s ease-in-out infinite 0.35s}
.kl .ae-bar:nth-child(9){animation:kl-b3 .9s ease-in-out infinite 0.12s}
.kl .ae-bar:nth-child(10){animation:kl-b2 .9s ease-in-out infinite 0.22s}
.kl .ae-prog-row{display:flex;align-items:center;justify-content:space-between;font-size:11px;color:var(--text3);margin-bottom:7px}
.kl .ae-prog-track{width:100%;height:3px;background:rgba(255,255,255,.07);border-radius:100px;overflow:hidden;margin-bottom:20px}
.kl .ae-prog-fill{height:100%;background:linear-gradient(90deg,var(--accent),var(--green));border-radius:100px;animation:kl-prog 4s ease-out .5s both}
.kl .ae-div{border:none;border-top:1px solid var(--border);margin-bottom:16px}
.kl .ae-qlabel{font-size:10px;font-weight:700;letter-spacing:.09em;text-transform:uppercase;color:var(--text3);margin-bottom:12px}
.kl .ae-queue{display:flex;flex-direction:column;gap:7px}
.kl .ae-qi{display:flex;align-items:center;gap:10px;padding:9px 12px;border-radius:10px;border:1px solid var(--border);background:rgba(255,255,255,.015)}
.kl .ae-qi-n{font-size:11px;font-weight:700;color:var(--text3);min-width:16px}
.kl .ae-qi-info{flex:1;min-width:0}
.kl .ae-qi-title{font-size:12px;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.kl .ae-qi-site{font-size:10.5px;color:var(--text3)}
.kl .ae-qi-badge{font-size:10px;padding:3px 8px;border-radius:100px;font-weight:600;white-space:nowrap}
.kl .ae-qi-badge.soon{color:var(--accent-lt);background:var(--accent-faint);border:1px solid rgba(167,139,250,.18)}
.kl .ae-qi-badge.queued{color:var(--text3);background:rgba(255,255,255,.04);border:1px solid var(--border)}

/* problem */
.kl .problem-section{border-top:1px solid var(--border)}
.kl .problem-grid{display:grid;grid-template-columns:repeat(3,1fr);border:1px solid var(--border);border-radius:14px;overflow:hidden;gap:1px;background:var(--border);margin-top:56px}
.kl .problem-card{padding:38px 32px;background:var(--bg2);transition:background .2s}
.kl .problem-card:hover{background:var(--bg3)}
.kl .problem-n{font-family:'Fraunces',Georgia,serif;font-size:52px;font-weight:700;color:var(--text3);line-height:1;margin-bottom:18px;letter-spacing:-.04em}
.kl .problem-title{font-size:16px;font-weight:700;margin-bottom:10px}
.kl .problem-desc{font-size:14px;color:var(--text2);line-height:1.7}

/* features tabbed */
.kl .features-section{border-top:1px solid var(--border)}
.kl .ft-wrap{display:grid;grid-template-columns:2fr 3fr;gap:20px;margin-top:56px;align-items:start}
.kl .ft-tabs{display:flex;flex-direction:column;gap:4px}
.kl .ft-tab{padding:20px 22px;border-radius:12px;cursor:pointer;text-align:left;background:transparent;font-family:'Syne',sans-serif;color:var(--text);border:1px solid transparent;border-left:2px solid transparent;transition:background .2s,border-color .2s}
.kl .ft-tab:hover{background:rgba(255,255,255,.025)}
.kl .ft-tab.on{background:rgba(109,40,232,.06);border-color:rgba(255,255,255,.07);border-left-color:var(--accent-lt)}
.kl .ft-tab-icon{font-size:22px;margin-bottom:8px;display:block}
.kl .ft-tab-name{font-size:15px;font-weight:700;margin-bottom:4px}
.kl .ft-tab-hint{font-size:13px;color:var(--text3);line-height:1.5}
.kl .ft-tab.on .ft-tab-hint{color:var(--text2)}
.kl .ft-preview{border:1px solid var(--border2);border-radius:16px;overflow:hidden;background:var(--bg2);position:relative;min-height:400px}
.kl .ft-bar{height:2px;background:linear-gradient(90deg,var(--accent),var(--green),transparent)}
.kl .ft-panel{position:absolute;inset:0;padding:28px;opacity:0;transform:translateY(10px);transition:opacity .3s,transform .3s;pointer-events:none;overflow:hidden}
.kl .ft-panel.on{opacity:1;transform:translateY(0);pointer-events:auto}
.kl .fp-toprow{display:flex;align-items:center;justify-content:space-between;margin-bottom:20px}
.kl .fp-tlabel{font-size:12px;font-weight:700;letter-spacing:.08em;text-transform:uppercase;color:var(--text3)}
.kl .fp-model{font-size:12px;font-weight:600;padding:5px 10px;border-radius:6px;background:rgba(255,255,255,.06);border:1px solid var(--border);color:var(--text)}
.kl .fp-inp{width:100%;background:rgba(255,255,255,.05);border:1px solid var(--border2);border-radius:8px;padding:11px 14px;font-size:13px;color:var(--text);font-family:'Syne',sans-serif;margin-bottom:10px;display:block}
.kl .fp-gbtn{width:100%;background:var(--accent);border:none;border-radius:8px;padding:12px;font-size:13.5px;font-weight:700;color:#fff;font-family:'Syne',sans-serif;cursor:default}
.kl .fp-ok{font-size:12px;color:var(--green);margin-top:10px;display:flex;align-items:center;gap:7px}
.kl .fp-divider{border:none;border-top:1px solid var(--border);margin:20px 0}
.kl .fp-metarow{display:flex;align-items:center;justify-content:space-between}
.kl .fp-meta-k{font-size:11px;color:var(--text3);font-weight:600}
.kl .fp-meta-v{font-size:12px;font-weight:700;padding:4px 10px;border-radius:6px;background:rgba(255,255,255,.06);border:1px solid var(--border)}
.kl .fp-day{font-size:11px;font-weight:700;color:var(--text3);text-transform:uppercase;letter-spacing:.08em;margin:0 0 8px}
.kl .fp-sched-list{display:flex;flex-direction:column;gap:6px}
.kl .fp-sched-item{display:flex;align-items:center;gap:12px;border:1px solid var(--border);border-radius:9px;padding:10px 14px;background:rgba(255,255,255,.02)}
.kl .fp-sched-time{font-size:12px;font-weight:700;color:var(--accent-lt);min-width:40px}
.kl .fp-sched-dot{width:7px;height:7px;border-radius:50%;background:var(--green);flex-shrink:0}
.kl .fp-sched-name{font-size:13px;font-weight:600;flex:1}
.kl .fp-sched-site{font-size:11px;color:var(--text3)}
.kl .fp-sched-footer{margin-top:16px;padding:10px 14px;border-radius:8px;background:var(--accent-faint);border:1px solid rgba(167,139,250,.18);font-size:12px;color:var(--accent-lt);font-weight:600}
.kl .fp-site-count{font-size:14px;font-weight:600;color:var(--text2);margin-bottom:14px}
.kl .fp-site-count b{color:var(--text)}
.kl .fp-sites{display:flex;flex-direction:column;gap:6px;margin-bottom:20px}
.kl .fp-site-row{display:flex;align-items:center;gap:12px;border:1px solid var(--border);border-radius:9px;padding:10px 14px;background:rgba(255,255,255,.02)}
.kl .fp-site-row.sel{border-color:rgba(167,139,250,.25);background:rgba(109,40,232,.06)}
.kl .fp-site-chk{width:16px;height:16px;border-radius:4px;border:1px solid var(--border2);display:flex;align-items:center;justify-content:center;font-size:9px;flex-shrink:0}
.kl .fp-site-row.sel .fp-site-chk{background:var(--accent);border-color:var(--accent);color:#fff}
.kl .fp-site-name{font-size:13px;font-weight:600;flex:1}
.kl .fp-post-btn{width:100%;background:var(--accent);border:none;border-radius:9px;padding:12px;font-size:13.5px;font-weight:700;color:#fff;font-family:'Syne',sans-serif;cursor:default}
.kl .fp-seo-header{display:flex;align-items:center;gap:20px;padding:18px;border:1px solid var(--border);border-radius:12px;background:rgba(255,255,255,.02);margin-bottom:22px}
.kl .fp-seo-big{font-family:'Fraunces',serif;font-size:62px;font-weight:700;letter-spacing:-.04em;line-height:1;color:var(--green)}
.kl .fp-seo-sub{font-size:13px;color:var(--text3);margin-top:2px}
.kl .fp-seo-track-wrap{flex:1}
.kl .fp-seo-track{height:6px;background:rgba(255,255,255,.07);border-radius:100px;overflow:hidden;margin-bottom:6px}
.kl .fp-seo-fill{height:100%;border-radius:100px;background:linear-gradient(90deg,var(--green),#22d3ee)}
.kl .fp-seo-status{font-size:12px;color:var(--text2);font-weight:600}
.kl .fp-metrics{display:flex;flex-direction:column;gap:10px}
.kl .fp-metric{display:flex;align-items:center;gap:12px}
.kl .fp-metric-name{font-size:13px;color:var(--text2);width:160px;flex-shrink:0}
.kl .fp-metric-bar{flex:1;height:5px;background:rgba(255,255,255,.06);border-radius:100px;overflow:hidden}
.kl .fp-metric-fill{height:100%;border-radius:100px}
.kl .fp-metric-val{font-size:12px;font-weight:700;min-width:30px;text-align:right}

/* how it works */
.kl .how-section{border-top:1px solid var(--border)}
.kl .how-center{text-align:center}
.kl .how-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:48px;margin-top:60px;position:relative}
.kl .how-grid::before{content:'';position:absolute;top:27px;left:15%;right:15%;height:1px;background:linear-gradient(to right,transparent,var(--border2),var(--border2),transparent)}
.kl .how-step{text-align:center}
.kl .how-circle{width:56px;height:56px;border-radius:50%;border:1px solid var(--border2);background:var(--bg2);display:flex;align-items:center;justify-content:center;margin:0 auto 22px;position:relative;z-index:1}
.kl .how-n{font-family:'Fraunces',serif;font-size:20px;font-weight:700;color:var(--accent-lt)}
.kl .how-title{font-size:16px;font-weight:700;margin-bottom:10px}
.kl .how-desc{font-size:14px;color:var(--text2);line-height:1.65}

/* for who */
.kl .who-section{border-top:1px solid var(--border)}
.kl .who-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:24px;margin-top:56px}
.kl .who-card{border:1px solid var(--border);border-radius:14px;padding:36px;position:relative;overflow:hidden;transition:border-color .25s,transform .25s}
.kl .who-card:hover{border-color:rgba(167,139,250,.28);transform:translateY(-2px)}
.kl .who-stripe{position:absolute;top:0;left:0;right:0;height:2px;background:linear-gradient(90deg,var(--accent) 0%,var(--green) 60%,transparent 100%)}
.kl .who-icon{font-size:34px;margin-bottom:18px}
.kl .who-title{font-size:20px;font-weight:700;margin-bottom:10px}
.kl .who-desc{font-size:14px;color:var(--text2);line-height:1.65;margin-bottom:20px}
.kl .who-tags{display:flex;flex-wrap:wrap;gap:8px}
.kl .who-tag{padding:4px 12px;border-radius:100px;font-size:12px;font-weight:600;color:var(--accent-lt);border:1px solid rgba(167,139,250,.2);background:var(--accent-faint)}

/* ══ 3. TESTIMONIALS (bankkroll-inspired stacked cards) ══ */
.kl .tst-section{border-top:1px solid var(--border)}
.kl .tst-layout{display:grid;grid-template-columns:1fr 1fr;gap:72px;align-items:center;margin-top:56px}
.kl .tst-trusted{margin-top:36px}
.kl .tst-trusted-label{font-size:11px;font-weight:700;letter-spacing:.09em;text-transform:uppercase;color:var(--text3);margin-bottom:14px}
.kl .tst-logos{display:flex;flex-wrap:wrap;gap:8px}
.kl .tst-logo-pill{padding:5px 14px;border:1px solid var(--border);border-radius:7px;font-size:12px;font-weight:600;color:var(--text3)}
.kl .tst-stack{position:relative}
.kl .tst-card{position:absolute;left:0;right:0;border:1px solid var(--border2);border-radius:18px;padding:28px;background:var(--bg2);transition:all .55s cubic-bezier(.16,1,.3,1)}
.kl .tst-card.active{top:0;opacity:1;transform:scale(1);z-index:3;box-shadow:0 24px 60px rgba(0,0,0,.5),0 0 40px rgba(109,40,232,.08)}
.kl .tst-card.behind-1{top:16px;opacity:.5;transform:scale(.96);z-index:2}
.kl .tst-card.behind-2{top:32px;opacity:.22;transform:scale(.92);z-index:1}
.kl .tst-card.hidden{opacity:0;transform:scale(.88) translateY(20px);z-index:0;top:32px;pointer-events:none}
.kl .tst-top{display:flex;align-items:center;gap:14px;margin-bottom:18px}
.kl .tst-avatar{width:44px;height:44px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-family:'Fraunces',serif;font-size:15px;font-weight:700;color:#fff;flex-shrink:0}
.kl .tst-name{font-size:15px;font-weight:700}
.kl .tst-role{font-size:12px;color:var(--text3);margin-top:1px}
.kl .tst-stars{display:flex;gap:3px;margin-bottom:14px}
.kl .tst-star{font-size:13px;color:#fbbf24}
.kl .tst-quote{font-size:14px;color:var(--text2);line-height:1.75}
.kl .tst-quote::before{content:open-quote;font-family:'Fraunces',serif;font-size:46px;color:rgba(109,40,232,.25);line-height:.6;display:block;margin-bottom:10px}
.kl .tst-nav{display:flex;align-items:center;gap:14px;margin-top:28px}
.kl .tst-nav-btn{width:38px;height:38px;border-radius:50%;border:1px solid var(--border2);background:transparent;color:var(--text2);cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:16px;transition:all .2s}
.kl .tst-nav-btn:hover{border-color:var(--border3);color:var(--text);background:rgba(255,255,255,.04)}
.kl .tst-dots{display:flex;gap:7px;align-items:center}
.kl .tst-dot{width:6px;height:6px;border-radius:50%;background:var(--text3);cursor:pointer;transition:all .25s}
.kl .tst-dot.active{width:20px;border-radius:3px;background:var(--accent-lt)}

/* faq */
.kl .faq-section{border-top:1px solid var(--border)}
.kl .faq-list{margin-top:52px}
.kl .faq-row{border-bottom:1px solid var(--border)}
.kl .faq-btn{width:100%;display:flex;align-items:center;justify-content:space-between;gap:20px;padding:22px 0;background:none;border:none;cursor:pointer;text-align:left;font-family:'Syne',sans-serif;font-size:15.5px;font-weight:600;color:var(--text);transition:color .2s}
.kl .faq-btn:hover{color:var(--accent-lt)}
.kl .faq-icon{font-size:18px;color:var(--text3);flex-shrink:0;transition:transform .22s,color .2s}
.kl .faq-icon.open{transform:rotate(45deg);color:var(--accent-lt)}
.kl .faq-body{font-size:15px;color:var(--text2);line-height:1.72;overflow:hidden;max-height:0;transition:max-height .32s ease,padding-bottom .32s ease}
.kl .faq-body.open{max-height:280px;padding-bottom:22px}

/* cta */
.kl .cta-section{border-top:1px solid var(--border);padding:96px 0}
.kl .cta-box{border:1px solid var(--border2);border-radius:20px;padding:88px 60px;text-align:center;position:relative;overflow:hidden;background:var(--bg2)}
.kl .cta-halo{position:absolute;top:0;left:50%;transform:translateX(-50%);width:500px;height:260px;pointer-events:none;background:radial-gradient(ellipse at top,rgba(109,40,232,.18) 0%,transparent 65%)}
.kl .cta-halo2{position:absolute;bottom:0;right:0;width:360px;height:200px;pointer-events:none;background:radial-gradient(ellipse at bottom right,rgba(52,211,153,.07) 0%,transparent 65%)}
.kl .cta-box h2{font-family:'Fraunces',serif;font-size:clamp(30px,5vw,54px);font-weight:700;letter-spacing:-.025em;line-height:1.1;margin-bottom:16px;position:relative}
.kl .cta-box p{font-size:17px;color:var(--text2);margin-bottom:36px;position:relative}
.kl .cta-row{display:flex;align-items:center;justify-content:center;gap:10px;flex-wrap:wrap;position:relative}
.kl .cta-note{font-size:12px;color:var(--text3);margin-top:14px}

/* footer */
.kl footer{border-top:1px solid var(--border);padding:44px 0}
.kl .footer-inner{max-width:1120px;margin:0 auto;padding:0 32px;display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:20px}
.kl .footer-logo{display:flex;align-items:center;gap:9px;font-weight:800;font-size:15px;color:var(--text);text-decoration:none}
.kl .footer-links{display:flex;gap:26px;flex-wrap:wrap}
.kl .footer-links a{font-size:13px;color:var(--text3);text-decoration:none;transition:color .2s}
.kl .footer-links a:hover{color:var(--text2)}
.kl .footer-copy{font-size:12px;color:var(--text3)}

@media(max-width:960px){
  .kl .hero-split{grid-template-columns:1fr}.kl .hero-text{text-align:center}
  .kl .hero-sub{margin-left:auto;margin-right:auto}.kl .hero-actions{justify-content:center}
  .kl .hero-stats{justify-content:center}.kl .hero-visual{display:none}
  .kl .ae-inner{grid-template-columns:1fr}.kl .tst-layout{grid-template-columns:1fr}
}
@media(max-width:860px){
  .kl .container,.kl .container-sm{padding:0 20px}.kl .nav-links{display:none}
  .kl .problem-grid{grid-template-columns:1fr}.kl .ft-wrap{grid-template-columns:1fr}
  .kl .who-grid{grid-template-columns:1fr}.kl .how-grid{grid-template-columns:1fr;gap:32px}
  .kl .how-grid::before{display:none}.kl .cta-box{padding:52px 28px}
  .kl .footer-inner{flex-direction:column;align-items:flex-start}
  .kl .cs-article{padding:24px 24px 52px}
}
`

/* ─────────────────────────────────────
   LOGO
───────────────────────────────────── */
function KabarLogo({ size = 28 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
      <ellipse cx="14" cy="14" rx="12" ry="5.5" stroke="url(#kl-og)" strokeWidth="1.5"
        strokeDasharray="4 3" transform="rotate(-35 14 14)" opacity="0.7"
        style={{ transformOrigin: '14px 14px', animation: 'kl-spin 8s linear infinite' }} />
      <path d="M14 5.5C14 5.5 18.5 8 18.5 13.5C18.5 17 16.5 19.5 14 21C11.5 19.5 9.5 17 9.5 13.5C9.5 8 14 5.5 14 5.5Z" fill="url(#kl-rg)" />
      <circle cx="14" cy="13" r="2" fill="#0d0d1c" fillOpacity="0.7" />
      <circle cx="14" cy="13" r="1.1" fill="url(#kl-wg)" />
      <path d="M9.5 16.5L7 19.5L10.5 18Z" fill="url(#kl-og)" opacity="0.8" />
      <path d="M18.5 16.5L21 19.5L17.5 18Z" fill="url(#kl-og)" opacity="0.8" />
      <path d="M12.5 21C12.5 21 13 23.5 14 24C15 23.5 15.5 21 15.5 21Z" fill="url(#kl-fg)" opacity="0.9" />
      <circle r="1.8" fill="#34d399">
        <animateMotion dur="4s" repeatCount="indefinite">
          <mpath xlinkHref="#kl-op" />
        </animateMotion>
      </circle>
      <path id="kl-op" d="M 2 14 A 12 5.5 -35 1 1 26 14" fill="none" />
      <defs>
        <linearGradient id="kl-rg" x1="14" y1="5.5" x2="14" y2="21" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#a78bfa"/><stop offset="100%" stopColor="#6d28e8"/>
        </linearGradient>
        <linearGradient id="kl-og" x1="0" y1="0" x2="28" y2="0" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#34d399"/><stop offset="100%" stopColor="#6d28e8"/>
        </linearGradient>
        <linearGradient id="kl-wg" x1="13" y1="12" x2="15" y2="14" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#e0d4ff"/><stop offset="100%" stopColor="#a78bfa"/>
        </linearGradient>
        <linearGradient id="kl-fg" x1="14" y1="21" x2="14" y2="24" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#fbbf24"/><stop offset="100%" stopColor="#f97316" stopOpacity="0"/>
        </linearGradient>
      </defs>
    </svg>
  )
}

/* ─────────────────────────────────────
   DATA
───────────────────────────────────── */
const FEATURES = [
  { icon: '✦', name: 'Generate Artikel + Gambar', hint: 'AI tulis artikel SEO lengkap dari satu keyword.' },
  { icon: '📅', name: 'Penjadwalan Otomatis', hint: 'Set jadwal sekali, konten posting sendiri tiap hari.' },
  { icon: '🌐', name: 'Multi-Website Sekaligus', hint: 'Satu generate, posting ke banyak website sekaligus.' },
  { icon: '📊', name: 'SEO Score Bawaan', hint: 'Nilai SEO otomatis sebelum artikel dipublish.' },
]

const PROBLEMS = [
  { n: '01', title: 'Nulis konten makan waktu', desc: 'Satu artikel bisa habiskan 2–4 jam. Kalau punya 5 niche site, itu 10–20 jam per minggu — hanya untuk konten.' },
  { n: '02', title: 'Posting tidak konsisten', desc: 'Algoritma Google menyukai website yang update rutin. Tapi konsistensi itu sulit ketika ada banyak hal lain yang harus dikerjakan.' },
  { n: '03', title: 'Biaya writer membengkak', desc: 'Hire penulis lepas per artikel bisa Rp 50–200 ribu. Untuk 30 artikel/bulan per website, biayanya langsung tak terkendali.' },
]

const STEPS = [
  { n: '1', title: 'Masukkan Keyword', desc: 'Ketik topik atau keyword. Kabar generate artikel lengkap dengan struktur heading, body, dan meta yang SEO-ready.' },
  { n: '2', title: 'Review & Finalisasi', desc: 'Cek preview artikel, lihat SEO score, tambahkan gambar. Edit seperlunya atau langsung lanjut ke posting.' },
  { n: '3', title: 'Jadwalkan & Tinggalkan', desc: 'Pilih website target, atur jadwal, selesai. Kabar handle sisanya — otomatis, tiap hari, tanpa intervensi.' },
]

const WHOS = [
  { icon: '🔍', title: 'Blogger & Pemain SEO', desc: 'Punya banyak niche site tapi kewalahan ngisi konten? Kabar bantu kamu scale ke semua website sekaligus dengan konten yang konsisten tiap hari.', tags: ['Niche Site', 'Affiliate', 'SEO Content', 'PBN'] },
  { icon: '📰', title: 'Content Agency & Media', desc: 'Kelola konten untuk banyak klien sekaligus? Kabar mempersingkat produksi dari jam menjadi menit, dengan kualitas yang bisa dikontrol penuh.', tags: ['Digital Agency', 'Content Marketing', 'Media Online'] },
]

const TESTIMONIALS = [
  { name: 'Rizki Pratama', role: 'SEO Specialist', company: 'Niche Site Pro', initials: 'RP', color: 'linear-gradient(135deg,#6d28e8,#a78bfa)', quote: 'Dulu ngisi 8 niche site sendiri hampir impossible. Sekarang dengan Kabar, semua site konsisten update tiap hari. Traffic organik naik 3x dalam 2 bulan.', stars: 5 },
  { name: 'Sari Dewi', role: 'Content Director', company: 'DigitalKu Agency', initials: 'SD', color: 'linear-gradient(135deg,#0d9488,#34d399)', quote: 'Klien kami butuh 20+ artikel per bulan per site. Kabar memangkas waktu produksi dari 3 minggu jadi 3 hari. Tim sekarang fokus ke strategi, bukan nulis.', stars: 5 },
  { name: 'Budi Santoso', role: 'Affiliate Marketer', company: 'Mandiri Digital', initials: 'BS', color: 'linear-gradient(135deg,#9333ea,#c084fc)', quote: 'Skeptis awalnya soal AI content. Tapi SEO score rata-rata 85+ dan artikel-artikelnya readable banget. Google belum pernah flagging satupun site saya.', stars: 5 },
]

const FAQS = [
  { q: 'Apakah konten AI akan dihukum Google?', a: 'Google tidak menghukum konten AI — yang dinilai adalah kualitas dan relevansi. Kabar mengoptimalkan artikel untuk SEO dan readability sehingga memenuhi standar Google. Ribuan website menggunakan konten AI tanpa masalah selama kualitasnya terjaga.' },
  { q: 'Website atau platform apa yang bisa diconnect?', a: 'Kabar mendukung berbagai platform blog dan CMS. Kamu bisa menambahkan dan mengelola beberapa "Produk" (website) langsung dari dashboard, lalu memilih satu atau banyak website sekaligus sebagai target posting.' },
  { q: 'Apakah bisa generate konten Bahasa Indonesia?', a: 'Ya, Kabar mendukung multi-bahasa. Kamu bisa generate konten dalam Bahasa Indonesia, Inggris, atau bahasa lain. Kabar dioptimalkan untuk konteks dan gaya penulisan Indonesia.' },
  { q: 'Bagaimana cara kerja penjadwalan posting?', a: 'Setelah generate artikel, pilih mode "Terjadwal" di Konfigurasi Posting, tentukan tanggal dan waktu, lalu pilih website target. Kabar akan otomatis memposting pada waktu yang kamu tentukan.' },
  { q: 'Apa bedanya model AI yang tersedia?', a: 'Kabar mengintegrasikan berbagai model via OpenRouter — termasuk Gemini, GPT-4, Claude, dan lainnya. Kamu bisa memilih model yang paling sesuai dengan kebutuhan kualitas dan anggaran.' },
]

/*
── PRICING (dinonaktifkan — harga belum ditentukan) ──
const PLANS = [
  { tier:'Starter', price:'Gratis',     per:'Selamanya', pop:false, features:['5 artikel/bulan','1 website','Generate gambar AI','SEO Score','Support komunitas'], cta:'Mulai Gratis', style:'ghost' },
  { tier:'Pro',     price:'Rp 99.000',  per:'per bulan', pop:true,  features:['50 artikel/bulan','5 website','Penjadwalan otomatis','Semua AI model','Prioritas generate','Email support'], cta:'Coba 7 Hari Gratis', style:'solid' },
  { tier:'Agency',  price:'Rp 299.000', per:'per bulan', pop:false, features:['Artikel tak terbatas','Website tak terbatas','Penjadwalan otomatis','Semua AI model','API access','Priority support'], cta:'Hubungi Kami', style:'ghost' },
]
*/

/* ─────────────────────────────────────
   FEATURE PANELS
───────────────────────────────────── */
function PanelGenerate() {
  return <>
    <div className="fp-toprow">
      <span className="fp-tlabel">Generate Konten</span>
      <span className="fp-model">Gemini 2.5 Flash ▾</span>
    </div>
    <div className="fp-inp" style={{ display: 'block' }}>10 Tips SEO Terbaik untuk Meningkatkan Traffic 2025</div>
    <button className="fp-gbtn">✦&nbsp;&nbsp;Generate Artikel</button>
    <div className="fp-ok"><span>●</span> Artikel siap · 1.240 kata</div>
    <hr className="fp-divider" />
    <div className="fp-metarow">
      <span className="fp-meta-k">SEO Score</span>
      <span className="fp-meta-v" style={{ color: 'var(--green)' }}>87 / 100</span>
    </div>
    <div className="fp-metarow" style={{ marginTop: 10 }}>
      <span className="fp-meta-k">Mode Posting</span>
      <span className="fp-meta-v">Terjadwal · Besok 09:00</span>
    </div>
  </>
}

function PanelSchedule() {
  const rows1 = [{ t: '09:00', name: '10 Tips SEO Terbaik 2025', site: 'techblog.id' }, { t: '15:00', name: 'Review Laptop Gaming Terbaik', site: 'reviewgadget.com' }]
  const rows2 = [{ t: '10:00', name: 'Strategi Affiliate Marketing 2025', site: 'affiliateku.net' }, { t: '14:00', name: 'Cara Monetisasi Blog dari Nol', site: 'tutorialseo.id' }]
  return <>
    <p className="fp-day">Hari ini</p>
    <div className="fp-sched-list">{rows1.map(r => <div key={r.t} className="fp-sched-item"><span className="fp-sched-time">{r.t}</span><span className="fp-sched-dot" /><span className="fp-sched-name">{r.name}</span><span className="fp-sched-site">{r.site}</span></div>)}</div>
    <p className="fp-day" style={{ marginTop: 18 }}>Besok</p>
    <div className="fp-sched-list">{rows2.map(r => <div key={r.t} className="fp-sched-item"><span className="fp-sched-time">{r.t}</span><span className="fp-sched-dot" /><span className="fp-sched-name">{r.name}</span><span className="fp-sched-site">{r.site}</span></div>)}</div>
    <div className="fp-sched-footer">✓&nbsp; 8 artikel terjadwal bulan ini</div>
  </>
}

function PanelMultiSite() {
  const sites = [{ n: 'techblog.id', s: true }, { n: 'reviewgadget.com', s: true }, { n: 'tutorialseo.id', s: true }, { n: 'affiliateku.net', s: false }, { n: 'produkfavorit.com', s: false }]
  return <>
    <p className="fp-site-count"><b>3</b> dari 5 website dipilih</p>
    <div className="fp-sites">{sites.map(site => (
      <div key={site.n} className={`fp-site-row ${site.s ? 'sel' : ''}`}>
        <div className="fp-site-chk">{site.s ? '✓' : ''}</div>
        <span className="fp-site-name">{site.n}</span>
        <span style={{ fontSize: 10, color: site.s ? 'var(--green)' : 'var(--text3)' }}>{site.s ? '● aktif' : '○'}</span>
      </div>
    ))}</div>
    <button className="fp-post-btn">📅&nbsp;&nbsp;Post ke 3 Website →</button>
  </>
}

function PanelSEO() {
  const metrics = [{ l: 'Keyword Density', p: 80, c: 'var(--green)' }, { l: 'Struktur Heading', p: 100, c: 'var(--green)' }, { l: 'Meta Description', p: 65, c: '#fbbf24' }, { l: 'Jumlah Kata', p: 100, c: 'var(--green)' }]
  return <>
    <div className="fp-seo-header">
      <div><div className="fp-seo-big">87</div><div className="fp-seo-sub">dari 100</div></div>
      <div className="fp-seo-track-wrap">
        <div className="fp-seo-track"><div className="fp-seo-fill" style={{ width: '87%' }} /></div>
        <div className="fp-seo-status">Skor Bagus — siap dipublish</div>
      </div>
    </div>
    <div className="fp-metrics">{metrics.map(m => (
      <div key={m.l} className="fp-metric">
        <span className="fp-metric-name">{m.l}</span>
        <div className="fp-metric-bar"><div className="fp-metric-fill" style={{ width: `${m.p}%`, background: m.c }} /></div>
        <span className="fp-metric-val" style={{ color: m.c }}>{m.p}</span>
      </div>
    ))}</div>
  </>
}

const PANELS = [PanelGenerate, PanelSchedule, PanelMultiSite, PanelSEO]

/* ─────────────────────────────────────
   ARTICLE ENGINE  (music-player theme)
───────────────────────────────────── */
function ArticleEngine() {
  const queue = [
    { title: 'Review Laptop Gaming Terbaik 2025', site: 'reviewgadget.com', badge: 'soon', label: 'Segera' },
    { title: 'Cara Monetisasi Blog dari Nol', site: 'tutorialseo.id', badge: 'queued', label: 'Antrian' },
    { title: 'Panduan Affiliate Marketing 2025', site: 'affiliateku.net', badge: 'queued', label: 'Antrian' },
  ]
  return (
    <div className="ae-card">
      <div className="ae-glow" />
      <div className="ae-top">
        <div className="ae-status"><span className="ae-dot" />Sedang Diproses</div>
        <div className="ae-time">09:41 WIB</div>
      </div>

      <div className="ae-main">
        <div className="ae-disc-wrap">
          <div className="ae-ring2" /><div className="ae-ring1" />
          <div className="ae-disc"><div className="ae-hole">✦</div></div>
        </div>
        <div className="ae-info">
          <div className="ae-atitle">10 Strategi SEO yang Wajib Dicoba di 2025</div>
          <div className="ae-site">techblog.id</div>
          <div className="ae-bars">{Array.from({ length: 10 }).map((_, i) => <div key={i} className="ae-bar" />)}</div>
        </div>
      </div>

      <div className="ae-prog-row">
        <span>Generate artikel</span>
        <span style={{ color: 'var(--accent-lt)', fontWeight: 700 }}>72%</span>
      </div>
      <div className="ae-prog-track"><div className="ae-prog-fill" /></div>

      <hr className="ae-div" />
      <div className="ae-qlabel">Selanjutnya · 3 artikel</div>
      <div className="ae-queue">{queue.map((item, i) => (
        <div key={i} className="ae-qi">
          <span className="ae-qi-n">{i + 2}</span>
          <div className="ae-qi-info">
            <div className="ae-qi-title">{item.title}</div>
            <div className="ae-qi-site">{item.site}</div>
          </div>
          <span className={`ae-qi-badge ${item.badge}`}>{item.label}</span>
        </div>
      ))}</div>
    </div>
  )
}

/* ─────────────────────────────────────
   TESTIMONIALS  (stacked cards)
───────────────────────────────────── */
function TestimonialsSection() {
  const [active, setActive] = useState(0)
  const n = TESTIMONIALS.length

  useEffect(() => {
    const t = setInterval(() => setActive(p => (p + 1) % n), 4500)
    return () => clearInterval(t)
  }, [n])

  const cls = (i: number) => {
    if (i === active) return 'active'
    if (i === (active + 1) % n) return 'behind-1'
    if (i === (active + 2) % n) return 'behind-2'
    return 'hidden'
  }

  return (
    <section className="tst-section">
      <div className="container">
        <div className="tst-layout" data-reveal>
          <div>
            <div className="section-label">Testimoni</div>
            <h2 className="section-title display">Dipercaya oleh<br />content creator Indonesia</h2>
            <p className="section-desc" style={{ maxWidth: 380 }}>
              Dari blogger niche site sampai agency digital, Kabar membantu ribuan kreator memproduksi konten lebih cepat dan konsisten.
            </p>
            <div className="tst-trusted">
              <div className="tst-trusted-label">Digunakan di berbagai niche</div>
              <div className="tst-logos">
                {['Teknologi', 'Kesehatan', 'Finance', 'Lifestyle', 'Otomotif', 'Travel'].map(l => (
                  <div key={l} className="tst-logo-pill">{l}</div>
                ))}
              </div>
            </div>
          </div>

          <div>
            <div className="tst-stack" style={{ height: `${240 + (n - 1) * 32}px` }}>
              {TESTIMONIALS.map((t, i) => (
                <div key={i} className={`tst-card ${cls(i)}`}>
                  <div className="tst-top">
                    <div className="tst-avatar" style={{ background: t.color }}>{t.initials}</div>
                    <div>
                      <div className="tst-name">{t.name}</div>
                      <div className="tst-role">{t.role} · {t.company}</div>
                    </div>
                  </div>
                  <div className="tst-stars">{Array.from({ length: t.stars }).map((_, si) => <span key={si} className="tst-star">★</span>)}</div>
                  <div className="tst-quote">{t.quote}</div>
                </div>
              ))}
            </div>
            <div className="tst-nav">
              <button className="tst-nav-btn" onClick={() => setActive(p => (p - 1 + n) % n)}>←</button>
              <div className="tst-dots">
                {TESTIMONIALS.map((_, i) => (
                  <div key={i} className={`tst-dot ${i === active ? 'active' : ''}`} onClick={() => setActive(i)} />
                ))}
              </div>
              <button className="tst-nav-btn" onClick={() => setActive(p => (p + 1) % n)}>→</button>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

/* ─────────────────────────────────────
   MAIN LANDING PAGE
───────────────────────────────────── */
export default function Landing() {
  const [activeFeature, setActiveFeature] = useState(0)
  const [openFaq, setOpenFaq] = useState<number | null>(null)
  const csSectionRef = useRef<HTMLDivElement>(null)
  const csMockupRef = useRef<HTMLDivElement>(null)

  // Scroll reveal
  useEffect(() => {
    const obs = new IntersectionObserver(
      entries => entries.forEach(e => { if (e.isIntersecting) e.target.classList.add('in-view') }),
      { threshold: 0.1, rootMargin: '-40px 0px' }
    )
    document.querySelectorAll('[data-reveal]').forEach(el => obs.observe(el))
    return () => obs.disconnect()
  }, [])

  // Container scroll — direct DOM mutation for 60fps
  useEffect(() => {
    const section = csSectionRef.current
    const mockup = csMockupRef.current
    if (!section || !mockup) return
    const onScroll = () => {
      const rect = section.getBoundingClientRect()
      const p = Math.max(0, Math.min(1, (window.innerHeight - rect.top) / (window.innerHeight * 0.72)))
      mockup.style.transform = `perspective(1100px) rotateX(${22 * (1 - p)}deg) scale(${0.76 + 0.24 * p})`
      mockup.style.opacity = String(0.3 + 0.7 * p)
      mockup.style.boxShadow = `0 ${20 + p * 80}px ${60 + p * 80}px rgba(0,0,0,.65),0 0 ${p * 100}px rgba(109,40,232,${(p * 0.18).toFixed(2)})`
    }
    window.addEventListener('scroll', onScroll, { passive: true })
    onScroll()
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  return (
    <div className="kl">
      <style dangerouslySetInnerHTML={{ __html: STYLES }} />
      <div className="grain" />

      {/* ── NAVBAR ── */}
      <nav>
        <div className="nav-inner">
          <a href="/" className="nav-logo"><KabarLogo size={28} />KABAR</a>
          <ul className="nav-links">
            <li><a href="#fitur">Fitur</a></li>
            <li><a href="#cara-kerja">Cara Kerja</a></li>
            <li><a href="#faq">FAQ</a></li>
          </ul>
          <div className="nav-right">
            <a href="/login" className="btn-ghost">Masuk</a>
            <a href="/login" className="btn-primary">Mulai Gratis</a>
          </div>
        </div>
      </nav>

      {/* ── HERO ── */}
      <section className="hero">
        <div className="hero-dots" /><div className="hero-glow-tr" /><div className="hero-glow-bl" />
        <div className="container">
          <div className="hero-split">
            {/* Left */}
            <div className="hero-text">
              <div className="badge au d1"><span className="badge-dot" />Platform Konten AI · Dibuat untuk Indonesia</div>
              <h1 className="hero-h1 display au d2">Isi Semua<br />Website-mu.<br /><em>Otomatis.</em></h1>
              <p className="hero-sub au d3">Generate artikel SEO berkualitas, buat gambar, dan posting terjadwal ke semua website-mu — tanpa menulis satu kata pun.</p>
              <div className="hero-actions au d4">
                <a href="/login" className="btn-primary lg">✦&nbsp; Mulai Gratis</a>
                <a href="#cara-kerja" className="btn-ghost lg">Cara Kerja →</a>
              </div>
              <p className="hero-disclaimer au d4">Tidak perlu kartu kredit · Setup dalam 2 menit</p>
              <div className="hero-stats au d5">
                <div><span className="hero-stat-n">500+</span><span className="hero-stat-l">artikel/hari</span></div>
                <div className="hero-stat-sep" />
                <div><span className="hero-stat-n">87</span><span className="hero-stat-l">rata-rata SEO score</span></div>
                <div className="hero-stat-sep" />
                <div><span className="hero-stat-n">5+</span><span className="hero-stat-l">model AI tersedia</span></div>
              </div>
            </div>
            {/* Right — browser mockup */}
            <div className="hero-visual">
              <div className="hero-vglow" />
              <div className="fbadge fb1 au d6" style={{ position: 'absolute', top: '10%', right: '-20px', zIndex: 10 }}>
                <span className="fbadge-icon">📋</span>
                <div><div className="fbadge-label">SEO Score</div><div className="fbadge-val" style={{ color: 'var(--green)' }}>87</div></div>
              </div>
              <div className="fbadge fb2 au d7" style={{ position: 'absolute', bottom: '28%', left: '-22px', zIndex: 10 }}>
                <span className="fbadge-icon">✅</span>
                <div><div className="fbadge-label">Artikel Siap</div><div className="fbadge-val">1.240 kata</div></div>
              </div>
              <div className="fbadge fb3 au d7" style={{ position: 'absolute', bottom: '10%', right: '8%', zIndex: 10 }}>
                <span className="fbadge-icon">🌐</span>
                <div><div className="fbadge-label">Posted ke</div><div className="fbadge-val">3 Website</div></div>
              </div>
              <div className="hero-tilt au d6">
                <div className="mwin">
                  <div className="mwin-bar">
                    <div className="mwin-dot" style={{ background: '#ff5f57' }} />
                    <div className="mwin-dot" style={{ background: '#ffbd2e' }} />
                    <div className="mwin-dot" style={{ background: '#28c840' }} />
                    <div className="mwin-url">app.kabar.id/generate</div>
                  </div>
                  <div className="mwin-body">
                    <div className="mwin-sb">
                      <div className="mwin-sb-sec">Menu</div>
                      {[['⊞', 'Dashboard', false], ['✦', 'Buat Konten', true], ['📦', 'Produk', false], ['📄', 'Draft', false]].map(([ic, lb, on]) => (
                        <div key={String(lb)} className={`mwin-sb-item ${on ? 'on' : ''}`}><span>{ic}</span><span>{lb}</span></div>
                      ))}
                      <div className="mwin-sb-sec">Manajemen</div>
                      {[['📅', 'Jadwal'], ['🕐', 'Riwayat'], ['⚙️', 'Pengaturan']].map(([ic, lb]) => (
                        <div key={String(lb)} className="mwin-sb-item"><span>{ic}</span><span>{lb}</span></div>
                      ))}
                    </div>
                    <div className="mwin-main">
                      <div className="mwin-title">Buat Konten</div>
                      <div className="mwin-sub">Buat artikel + gambar dengan AI</div>
                      <div className="mwin-pill"><span>📋 SEO: <b>87</b></span><span>📝 Kata: <span>1.240</span></span></div>
                      <div className="mwin-grid">
                        <div className="mwin-card">
                          <div className="mwin-clabel">✏️ Topik Artikel</div>
                          <div className="mwin-inp" style={{ display: 'block' }}>10 Tips SEO Terbaik untuk 2025</div>
                          <button className="mwin-gbtn">✦ Buat Artikel</button>
                          <div className="mwin-ok">● Artikel siap. Tambahkan gambar!</div>
                        </div>
                        <div className="mwin-card">
                          <div className="mwin-clabel">⚙️ Konfigurasi</div>
                          <div className="mwin-modes">
                            <div className="mwin-mode">Langsung</div>
                            <div className="mwin-mode on">Jadwal</div>
                            <div className="mwin-mode">Draf</div>
                          </div>
                          <div className="mwin-tgt">Target: 3 website dipilih</div>
                          <button className="mwin-pbtn">📅 Jadwal Postingan →</button>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="mwin-fade" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ── COMPAT BAR ── */}
      <div className="compat">
        <div className="compat-inner">
          <span className="compat-label">Didukung oleh</span>
          {['WordPress', 'OpenRouter', 'Gemini', 'GPT-4o', 'Claude', 'Llama'].map(l => (
            <div key={l} className="compat-chip">{l}</div>
          ))}
        </div>
      </div>

      {/* ══ SECTION 1: CONTAINER SCROLL (Aceternity-inspired) ══ */}
      <section className="cs-section">
        <div className="cs-outer" ref={csSectionRef} style={{ height: '240vh' }}>
          <div className="cs-sticky">
            <div className="cs-bg-glow" />
            <div className="cs-header" data-reveal>
              <div className="cs-label">Hasil Generate</div>
              <h2 className="cs-title display">Artikel yang rapi,<br />terstruktur, dan SEO-ready</h2>
              <p className="cs-desc">Bukan sekadar teks mentah — Kabar menghasilkan artikel yang benar-benar siap publish, lengkap dengan heading, meta, dan skor SEO bawaan.</p>
            </div>
            <div
              className="cs-wrap"
              ref={csMockupRef}
              style={{ transform: 'perspective(1100px) rotateX(22deg) scale(0.76)', opacity: '0.3', willChange: 'transform,opacity,box-shadow' }}
            >
              <div className="cs-win">
                <div className="cs-win-bar">
                  <div className="cs-win-dot" style={{ background: '#ff5f57' }} />
                  <div className="cs-win-dot" style={{ background: '#ffbd2e' }} />
                  <div className="cs-win-dot" style={{ background: '#28c840' }} />
                  <div className="cs-win-tabs">
                    <div className="cs-win-tab active">Preview Artikel</div>
                    <div className="cs-win-tab">Preview Gambar</div>
                    <div className="cs-win-tab">Ringkasan</div>
                  </div>
                </div>
                <div className="cs-strip">
                  <span>📋 SEO Score: <span className="cs-strip-score">87 / 100</span></span>
                  <span className="cs-strip-dot" />
                  <span>📝 1.240 kata</span>
                  <span className="cs-strip-dot" />
                  <span>✅ Siap publish</span>
                  <span style={{ marginLeft: 'auto', color: 'var(--accent-lt)', fontWeight: 600 }}>techblog.id · Besok 09:00</span>
                </div>
                <div className="cs-article">
                  <div className="cs-article-tag">SEO · Tutorial</div>
                  <h1>10 Strategi SEO yang Wajib Dicoba di 2025 untuk Tingkatkan Traffic Organik</h1>
                  <div className="cs-article-meta">
                    <span>Kabar AI</span><span>·</span><span>8 min baca</span><span>·</span><span>24 Mei 2025</span>
                  </div>
                  <p>Di era persaingan digital yang semakin ketat, memiliki strategi SEO yang tepat bukan lagi sekadar pilihan — melainkan keharusan. Artikel ini membahas 10 strategi terkini yang terbukti efektif meningkatkan peringkat di halaman pertama Google.</p>
                  <h2>1. Fokus pada Search Intent, Bukan Hanya Keyword</h2>
                  <p>Google semakin pintar memahami konteks pencarian. Alih-alih mengisi artikel dengan keyword sebanyak mungkin, fokuslah pada niat di balik pencarian pengguna — apakah mereka mencari informasi, ingin membeli, atau sedang membandingkan produk?</p>
                  <h2>2. Optimalkan Core Web Vitals</h2>
                  <p>Sejak pembaruan Page Experience, Google secara eksplisit menggunakan Core Web Vitals sebagai faktor ranking. Pastikan LCP di bawah 2.5 detik dan CLS di bawah 0.1 untuk hasil terbaik.</p>
                  <div className="cs-fade" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ── PROBLEM ── */}
      <section className="problem-section">
        <div className="container">
          <div data-reveal>
            <div className="section-label">Masalah yang Kita Selesaikan</div>
            <h2 className="section-title display">Kenapa content creator selalu<br />kehabisan waktu?</h2>
            <p className="section-desc" style={{ maxWidth: 480 }}>Bukan karena malas. Tapi karena cara lama memang tidak skalabel.</p>
          </div>
          <div className="problem-grid" data-reveal data-delay="2">
            {PROBLEMS.map(p => (
              <div key={p.n} className="problem-card">
                <div className="problem-n">{p.n}</div>
                <div className="problem-title">{p.title}</div>
                <div className="problem-desc">{p.desc}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── FEATURES TABBED ── */}
      <section className="features-section" id="fitur">
        <div className="container">
          <div data-reveal>
            <div className="section-label">Fitur Utama</div>
            <h2 className="section-title display">Satu platform untuk<br />semua kebutuhan kontenmu</h2>
          </div>
          <div className="ft-wrap" data-reveal data-delay="2">
            <div className="ft-tabs">
              {FEATURES.map((f, i) => (
                <button key={i} className={`ft-tab ${activeFeature === i ? 'on' : ''}`} onClick={() => setActiveFeature(i)}>
                  <span className="ft-tab-icon">{f.icon}</span>
                  <div className="ft-tab-name">{f.name}</div>
                  <div className="ft-tab-hint">{f.hint}</div>
                </button>
              ))}
            </div>
            <div className="ft-preview">
              <div className="ft-bar" />
              {PANELS.map((Panel, i) => (
                <div key={i} className={`ft-panel ${activeFeature === i ? 'on' : ''}`}><Panel /></div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* ══ SECTION 2: ARTICLE ENGINE (music-player aesthetic) ══ */}
      <section className="ae-section">
        <div className="container">
          <div className="ae-inner">
            <div data-reveal>
              <div className="section-label">Kabar Bekerja 24/7</div>
              <h2 className="section-title display">AI yang tidak<br />pernah berhenti<br />berkarya</h2>
              <p className="section-desc" style={{ maxWidth: 400, marginBottom: 28 }}>
                Sementara kamu tidur, Kabar terus generate dan menjadwalkan konten untuk semua website-mu. Queue artikel berjalan otomatis, tanpa henti.
              </p>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
                {[
                  { i: '⚡', t: 'Generate artikel dalam hitungan detik' },
                  { i: '🔄', t: 'Queue artikel berjalan tanpa perlu diawasi' },
                  { i: '🌙', t: 'Jadwal tengah malam pun bisa diatur' },
                ].map(x => (
                  <div key={x.t} style={{ display: 'flex', alignItems: 'center', gap: 12, fontSize: 14, color: 'var(--text2)' }}>
                    <span style={{ fontSize: 18 }}>{x.i}</span>{x.t}
                  </div>
                ))}
              </div>
            </div>
            <div data-reveal data-delay="2"><ArticleEngine /></div>
          </div>
        </div>
      </section>

      {/* ── HOW IT WORKS ── */}
      <section className="how-section" id="cara-kerja">
        <div className="container">
          <div className="how-center" data-reveal>
            <div className="section-label" style={{ justifyContent: 'center' }}>Cara Kerja</div>
            <h2 className="section-title display">3 langkah, website-mu<br />penuh konten</h2>
          </div>
          <div className="how-grid" data-reveal data-delay="2">
            {STEPS.map(s => (
              <div key={s.n} className="how-step">
                <div className="how-circle"><span className="how-n">{s.n}</span></div>
                <div className="how-title">{s.title}</div>
                <div className="how-desc">{s.desc}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── FOR WHO ── */}
      <section className="who-section">
        <div className="container">
          <div data-reveal>
            <div className="section-label">Untuk Siapa</div>
            <h2 className="section-title display">Kabar dibuat untuk kamu<br />yang serius di konten</h2>
          </div>
          <div className="who-grid" data-reveal data-delay="2">
            {WHOS.map(w => (
              <div key={w.title} className="who-card">
                <div className="who-stripe" />
                <div className="who-icon">{w.icon}</div>
                <div className="who-title">{w.title}</div>
                <div className="who-desc">{w.desc}</div>
                <div className="who-tags">{w.tags.map(t => <span key={t} className="who-tag">{t}</span>)}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ══ SECTION 3: TESTIMONIALS (bankkroll-inspired stacked cards) ══ */}
      <TestimonialsSection />

      {/* ── FAQ ── */}
      <section className="faq-section" id="faq">
        <div className="container-sm">
          <div data-reveal>
            <div className="section-label">FAQ</div>
            <h2 className="section-title display">Pertanyaan yang<br />sering ditanyakan</h2>
          </div>
          <div className="faq-list" data-reveal data-delay="2">
            {FAQS.map((faq, i) => (
              <div key={i} className="faq-row">
                <button className="faq-btn" onClick={() => setOpenFaq(openFaq === i ? null : i)}>
                  <span>{faq.q}</span>
                  <span className={`faq-icon ${openFaq === i ? 'open' : ''}`}>+</span>
                </button>
                <div className={`faq-body ${openFaq === i ? 'open' : ''}`}>{faq.a}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── FINAL CTA ── */}
      <section className="cta-section">
        <div className="container">
          <div data-reveal>
            <div className="cta-box">
              <div className="cta-halo" /><div className="cta-halo2" />
              <h2 className="display">Siap isi semua website-mu<br />tanpa kerja extra?</h2>
              <p>Mulai hari ini. Gratis. Tanpa kartu kredit.</p>
              <div className="cta-row">
                <a href="/login" className="btn-primary lg">✦&nbsp; Mulai Gratis Sekarang</a>
                <a href="/login" className="btn-ghost lg">Masuk ke Akun →</a>
              </div>
              <p className="cta-note">Bergabung dengan ratusan content creator Indonesia 🇮🇩</p>
            </div>
          </div>
        </div>
      </section>

      {/* ── FOOTER ── */}
      <footer>
        <div className="footer-inner">
          <a href="/" className="footer-logo"><KabarLogo size={22} />KABAR</a>
          <div className="footer-links">
            <a href="#fitur">Fitur</a>
            <a href="#cara-kerja">Cara Kerja</a>
            <a href="#faq">FAQ</a>
            <a href="/privacy">Kebijakan Privasi</a>
            <a href="/terms">Syarat &amp; Ketentuan</a>
          </div>
          <p className="footer-copy">© 2025 Kabar · Dibuat di Indonesia 🇮🇩</p>
        </div>
      </footer>
    </div>
  )
}