'use client'

import { useState, useEffect } from 'react'

/* ─────────────────────────────────────────
   ALL STYLES
───────────────────────────────────────── */
const STYLES = `
@import url('https://fonts.googleapis.com/css2?family=Fraunces:ital,opsz,wght@0,9..144,300;0,9..144,600;0,9..144,700;1,9..144,300;1,9..144,600&family=Syne:wght@400;500;600;700;800&display=swap');

.kl {
  font-family: 'Syne', system-ui, sans-serif;
  background: #07070f;
  color: #ededf7;
  min-height: 100vh;
  --bg:           #07070f;
  --bg2:          #0c0c1a;
  --bg3:          #101020;
  --border:       rgba(255,255,255,0.065);
  --border2:      rgba(255,255,255,0.11);
  --border3:      rgba(255,255,255,0.22);
  --text:         #ededf7;
  --text2:        #7e7ea0;
  --text3:        #3e3e58;
  --accent:       #7c3aed;
  --accent-lt:    #a78bfa;
  --accent-faint: rgba(124,58,237,0.1);
  --accent-glow:  rgba(124,58,237,0.28);
}
.kl *, .kl *::before, .kl *::after { box-sizing: border-box; }

/* GRAIN */
.kl .grain {
  position: fixed; inset: 0; pointer-events: none; z-index: 999; opacity: 0.024;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='300' height='300'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='300' height='300' filter='url(%23n)'/%3E%3C/svg%3E");
}

/* LAYOUT */
.kl .container    { max-width: 1120px; margin: 0 auto; padding: 0 32px; }
.kl .container-sm { max-width:  760px; margin: 0 auto; padding: 0 32px; }
.kl .display      { font-family: 'Fraunces', Georgia, serif; }

/* ═══ ANIMATIONS ═══ */
@keyframes kl-up   { from { opacity:0; transform: translateY(22px); } to { opacity:1; transform: translateY(0); } }
@keyframes kl-blink{ 0%,100%{opacity:1} 50%{opacity:.35} }
@keyframes kl-float-a { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-10px)} }
@keyframes kl-float-b { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-14px)} }
@keyframes kl-float-c { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-8px)} }
@keyframes kl-persp  {
  from { opacity:0; transform: perspective(1200px) rotateY(-18deg) rotateX(8deg) translateY(32px) scale(0.97); }
  to   { opacity:1; transform: perspective(1200px) rotateY(-8deg)  rotateX(3deg) translateY(0)    scale(1);    }
}

.kl .au { animation: kl-up .65s cubic-bezier(.16,1,.3,1) forwards; }
.kl .d1 { animation-delay:.05s; opacity:0; }
.kl .d2 { animation-delay:.15s; opacity:0; }
.kl .d3 { animation-delay:.25s; opacity:0; }
.kl .d4 { animation-delay:.35s; opacity:0; }
.kl .d5 { animation-delay:.45s; opacity:0; }
.kl .d6 { animation-delay:.60s; opacity:0; }
.kl .d7 { animation-delay:.75s; opacity:0; }

/* ═══ SCROLL REVEAL ═══ */
[data-reveal] {
  opacity: 0;
  transform: translateY(30px);
  transition: opacity .7s cubic-bezier(.16,1,.3,1), transform .7s cubic-bezier(.16,1,.3,1);
}
[data-reveal].in-view { opacity:1; transform:translateY(0); }
[data-reveal][data-delay="1"].in-view { transition-delay:.08s; }
[data-reveal][data-delay="2"].in-view { transition-delay:.16s; }
[data-reveal][data-delay="3"].in-view { transition-delay:.24s; }
[data-reveal][data-delay="4"].in-view { transition-delay:.32s; }
[data-reveal][data-delay="5"].in-view { transition-delay:.40s; }

/* ═══ BUTTONS ═══ */
.kl .btn-primary {
  display: inline-flex; align-items: center; gap: 7px;
  background: var(--accent); color: #fff;
  font-family: 'Syne', sans-serif; font-weight: 600; font-size: 14px;
  padding: 10px 20px; border-radius: 8px; border: none; cursor: pointer;
  text-decoration: none; letter-spacing: .01em;
  transition: background .2s, transform .15s, box-shadow .2s;
}
.kl .btn-primary:hover { background:#6d28d9; transform:translateY(-1px); box-shadow:0 8px 28px var(--accent-glow); }
.kl .btn-primary.lg { font-size:15px; padding:13px 26px; border-radius:9px; }
.kl .btn-ghost {
  display: inline-flex; align-items: center; gap: 7px;
  background: transparent; color: var(--text2);
  font-family: 'Syne', sans-serif; font-weight: 500; font-size: 14px;
  padding: 10px 16px; border-radius: 8px; border: 1px solid var(--border2);
  cursor: pointer; text-decoration: none; transition: all .2s;
}
.kl .btn-ghost:hover { color:var(--text); border-color:var(--border3); background:rgba(255,255,255,.04); }
.kl .btn-ghost.lg { font-size:15px; padding:13px 20px; }

/* ═══ NAVBAR ═══ */
.kl nav {
  position: fixed; top: 0; left: 0; right: 0; z-index: 100;
  border-bottom: 1px solid var(--border);
  background: rgba(7,7,15,.82);
  backdrop-filter: blur(18px); -webkit-backdrop-filter: blur(18px);
}
.kl .nav-inner { display:flex; align-items:center; justify-content:space-between; height:60px; max-width:1120px; margin:0 auto; padding:0 32px; }
.kl .nav-logo  { display:flex; align-items:center; gap:9px; font-weight:800; font-size:17px; letter-spacing:-.02em; color:var(--text); text-decoration:none; }
.kl .nav-logo-mark { width:28px; height:28px; border-radius:7px; background:var(--accent); display:flex; align-items:center; justify-content:center; font-size:13px; }
.kl .nav-links { display:flex; align-items:center; gap:30px; list-style:none; }
.kl .nav-links a { font-size:13.5px; font-weight:500; color:var(--text2); text-decoration:none; transition:color .2s; }
.kl .nav-links a:hover { color:var(--text); }
.kl .nav-right { display:flex; align-items:center; gap:10px; }

/* ═══ HERO — SPLIT LAYOUT ═══ */
.kl .hero {
  padding-top: 60px;
  min-height: 100vh;
  display: flex; align-items: center;
  position: relative; overflow: hidden;
}
.kl .hero-dots {
  position:absolute; inset:0; pointer-events:none;
  background-image: radial-gradient(circle at 1px 1px, rgba(255,255,255,.05) 1px, transparent 0);
  background-size: 36px 36px;
  mask-image: radial-gradient(ellipse 90% 90% at 25% 50%, black 20%, transparent 75%);
  -webkit-mask-image: radial-gradient(ellipse 90% 90% at 25% 50%, black 20%, transparent 75%);
}
.kl .hero-glow-tr {
  position:absolute; top:-100px; right:-100px; width:700px; height:700px; pointer-events:none;
  background: radial-gradient(ellipse at top right, rgba(124,58,237,.15) 0%, transparent 55%);
}
.kl .hero-split {
  display: grid; grid-template-columns: 5fr 7fr;
  gap: 56px; align-items: center;
  padding: 80px 0;
  position: relative; z-index: 1; width: 100%;
}

/* Hero — left text */
.kl .hero-text {}
.kl .badge {
  display: inline-flex; align-items: center; gap: 8px;
  border: 1px solid rgba(167,139,250,.22); border-radius:100px;
  padding: 5px 15px; font-size: 11.5px; font-weight: 600; letter-spacing: .08em;
  text-transform: uppercase; color: var(--accent-lt);
  background: rgba(124,58,237,.08); margin-bottom: 24px;
}
.kl .badge-dot { width:5px; height:5px; border-radius:50%; background:var(--accent-lt); animation: kl-blink 2s ease-in-out infinite; }
.kl .hero-h1 {
  font-family: 'Fraunces', Georgia, serif;
  font-size: clamp(40px, 4.8vw, 66px); font-weight: 700; line-height: 1.04;
  letter-spacing: -.03em; margin-bottom: 20px;
  background: linear-gradient(180deg, #ffffff 30%, rgba(200,200,220,.65) 100%);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
}
.kl .hero-h1 em { font-style:italic; -webkit-text-fill-color: var(--accent-lt); }
.kl .hero-sub { font-size:17px; color:var(--text2); line-height:1.7; max-width:420px; margin-bottom:32px; }
.kl .hero-actions { display:flex; align-items:center; gap:10px; flex-wrap:wrap; margin-bottom:12px; }
.kl .hero-disclaimer { font-size:12px; color:var(--text3); margin-bottom:36px; }
.kl .hero-stats { display:flex; align-items:center; gap:24px; padding-top:28px; border-top:1px solid var(--border); }
.kl .hero-stat-n { display:block; font-family:'Fraunces',serif; font-size:28px; font-weight:700; letter-spacing:-.02em; }
.kl .hero-stat-l { display:block; font-size:12px; color:var(--text3); margin-top:2px; }
.kl .hero-stat-sep { width:1px; height:36px; background:var(--border2); }

/* Hero — right visual */
.kl .hero-visual { position:relative; }
.kl .hero-visual-glow {
  position:absolute; inset:-80px; pointer-events:none;
  background: radial-gradient(ellipse at 55% 40%, rgba(124,58,237,.2) 0%, transparent 55%);
}
.kl .hero-mockup-tilt {
  position: relative; z-index: 1;
  animation: kl-persp .9s cubic-bezier(.16,1,.3,1) .5s both;
  transform-origin: left center;
  transition: transform .5s ease;
}
.kl .hero-mockup-tilt:hover {
  transform: perspective(1200px) rotateY(-4deg) rotateX(1deg) !important;
}

/* Floating badges */
.kl .fbadge {
  position: absolute; z-index: 10;
  display: flex; align-items: center; gap: 10px;
  background: rgba(8,8,20,.94); border: 1px solid var(--border2);
  border-radius: 12px; padding: 10px 14px;
  backdrop-filter: blur(16px); -webkit-backdrop-filter: blur(16px);
  box-shadow: 0 8px 32px rgba(0,0,0,.55);
  white-space: nowrap;
}
.kl .fbadge-icon { font-size:19px; }
.kl .fbadge-label { font-size:10.5px; font-weight:600; color:var(--text3); text-transform:uppercase; letter-spacing:.06em; margin-bottom:2px; }
.kl .fbadge-val   { font-family:'Fraunces',serif; font-size:18px; font-weight:700; line-height:1.1; }
.kl .fb1 { animation: kl-float-a 4s   ease-in-out infinite; }
.kl .fb2 { animation: kl-float-b 5.2s ease-in-out infinite 1.1s; }
.kl .fb3 { animation: kl-float-c 4.6s ease-in-out infinite .6s; }

/* ═══ MINI BROWSER WINDOW ═══ */
.kl .mwin {
  border: 1px solid var(--border2); border-radius:14px; overflow:hidden;
  background: #0b0b1c;
  box-shadow: 0 0 0 1px rgba(255,255,255,.03), 0 40px 80px rgba(0,0,0,.7), 0 0 60px rgba(124,58,237,.07);
  position: relative;
}
.kl .mwin-bar { display:flex; align-items:center; gap:7px; padding:11px 16px; border-bottom:1px solid var(--border); background:rgba(0,0,0,.3); }
.kl .mwin-dot  { width:10px; height:10px; border-radius:50%; }
.kl .mwin-url  { flex:1; margin:0 14px; height:20px; border-radius:4px; background:rgba(255,255,255,.04); border:1px solid var(--border); display:flex; align-items:center; padding:0 10px; font-size:10px; color:var(--text3); }
.kl .mwin-body { display:grid; grid-template-columns:175px 1fr; }
.kl .mwin-sb   { border-right:1px solid var(--border); padding:14px 12px; background:rgba(0,0,0,.2); }
.kl .mwin-sb-sec { font-size:9.5px; font-weight:700; letter-spacing:.1em; text-transform:uppercase; color:var(--text3); margin:12px 0 5px 8px; }
.kl .mwin-sb-sec:first-child { margin-top:0; }
.kl .mwin-sb-item { display:flex; align-items:center; gap:8px; padding:6px 8px; border-radius:5px; font-size:11px; color:var(--text3); margin-bottom:1px; }
.kl .mwin-sb-item.on { background:rgba(124,58,237,.14); color:var(--accent-lt); border:1px solid rgba(124,58,237,.18); }
.kl .mwin-main { padding:18px 22px; }
.kl .mwin-title { font-size:16px; font-weight:700; margin-bottom:2px; }
.kl .mwin-sub   { font-size:10px; color:var(--text3); margin-bottom:14px; }
.kl .mwin-pill  { display:inline-flex; align-items:center; gap:12px; border:1px solid var(--border); border-radius:6px; padding:6px 11px; font-size:10px; color:var(--text2); background:rgba(255,255,255,.02); margin-bottom:14px; }
.kl .mwin-pill b { color:#4ade80; }
.kl .mwin-pill span { color:var(--accent-lt); font-weight:700; }
.kl .mwin-grid  { display:grid; grid-template-columns:1fr 1fr; gap:10px; }
.kl .mwin-card  { border:1px solid var(--border); border-radius:8px; padding:10px; background:rgba(255,255,255,.015); }
.kl .mwin-clabel{ font-size:10px; font-weight:600; margin-bottom:8px; }
.kl .mwin-inp   { width:100%; background:rgba(255,255,255,.06); border:1px solid var(--border); border-radius:5px; padding:6px 8px; font-size:10px; color:var(--text2); font-family:'Syne',sans-serif; margin-bottom:7px; display:block; }
.kl .mwin-gbtn  { width:100%; background:var(--accent); border:none; border-radius:5px; padding:7px; font-size:10px; font-weight:700; color:#fff; font-family:'Syne',sans-serif; }
.kl .mwin-ok    { font-size:9px; color:#4ade80; margin-top:5px; }
.kl .mwin-modes { display:flex; gap:3px; margin-bottom:8px; }
.kl .mwin-mode  { flex:1; padding:5px 2px; border-radius:4px; text-align:center; font-size:9.5px; font-weight:600; background:rgba(255,255,255,.04); color:var(--text3); }
.kl .mwin-mode.on { background:var(--accent); color:#fff; }
.kl .mwin-tgt   { font-size:9px; color:var(--text3); margin-bottom:6px; }
.kl .mwin-pbtn  { width:100%; background:var(--accent); border:none; border-radius:5px; padding:7px; font-size:9.5px; font-weight:700; color:#fff; font-family:'Syne',sans-serif; opacity:.9; }
.kl .mwin-fade  { position:absolute; bottom:0; left:0; right:0; height:64px; background:linear-gradient(to top, #07070f, transparent); pointer-events:none; }

/* ═══ COMPAT BAR ═══ */
.kl .compat { padding:24px 0; border-top:1px solid var(--border); border-bottom:1px solid var(--border); }
.kl .compat-inner { max-width:1120px; margin:0 auto; padding:0 32px; display:flex; align-items:center; justify-content:center; gap:10px; flex-wrap:wrap; }
.kl .compat-label { font-size:12px; color:var(--text3); font-weight:500; margin-right:4px; }
.kl .compat-chip  { padding:5px 14px; border:1px solid var(--border); border-radius:7px; font-size:12px; font-weight:600; color:var(--text3); }

/* ═══ SECTIONS ═══ */
.kl section { padding:96px 0; }
.kl .section-label { display:inline-flex; align-items:center; gap:10px; font-size:11px; font-weight:700; letter-spacing:.1em; text-transform:uppercase; color:var(--accent-lt); margin-bottom:18px; }
.kl .section-label::before { content:''; width:16px; height:1px; background:var(--accent-lt); }
.kl .section-title { font-family:'Fraunces',Georgia,serif; font-size:clamp(30px,4vw,48px); font-weight:700; letter-spacing:-.025em; line-height:1.1; margin-bottom:16px; }
.kl .section-desc  { font-size:17px; color:var(--text2); line-height:1.65; }

/* ═══ PROBLEM ═══ */
.kl .problem-section { border-top:1px solid var(--border); }
.kl .problem-grid { display:grid; grid-template-columns:repeat(3,1fr); border:1px solid var(--border); border-radius:14px; overflow:hidden; gap:1px; background:var(--border); margin-top:56px; }
.kl .problem-card { padding:38px 32px; background:var(--bg2); transition:background .2s; }
.kl .problem-card:hover { background:var(--bg3); }
.kl .problem-n { font-family:'Fraunces',Georgia,serif; font-size:52px; font-weight:700; color:var(--text3); line-height:1; margin-bottom:18px; letter-spacing:-.04em; }
.kl .problem-title { font-size:16px; font-weight:700; margin-bottom:10px; }
.kl .problem-desc  { font-size:14px; color:var(--text2); line-height:1.7; }

/* ═══ FEATURES — TABBED ═══ */
.kl .features-section { border-top:1px solid var(--border); }
.kl .ft-wrap  { display:grid; grid-template-columns:2fr 3fr; gap:20px; margin-top:56px; align-items:start; }
.kl .ft-tabs  { display:flex; flex-direction:column; gap:4px; }
.kl .ft-tab {
  padding:20px 22px; border-radius:12px; cursor:pointer; text-align:left;
  background:transparent; font-family:'Syne',sans-serif; color:var(--text);
  border:1px solid transparent; border-left:2px solid transparent;
  transition:background .2s, border-color .2s;
}
.kl .ft-tab:hover { background:rgba(255,255,255,.025); }
.kl .ft-tab.on   { background:rgba(255,255,255,.03); border-color:rgba(255,255,255,.07); border-left-color:var(--accent-lt); }
.kl .ft-tab-icon { font-size:22px; margin-bottom:8px; display:block; }
.kl .ft-tab-name { font-size:15px; font-weight:700; margin-bottom:4px; }
.kl .ft-tab-hint { font-size:13px; color:var(--text3); line-height:1.5; }
.kl .ft-tab.on .ft-tab-hint { color:var(--text2); }

.kl .ft-preview {
  border:1px solid var(--border2); border-radius:16px; overflow:hidden;
  background:var(--bg2); position:relative; min-height:420px;
}
.kl .ft-bar { height:2px; background:linear-gradient(90deg, var(--accent), var(--accent-lt), transparent); }
.kl .ft-panel {
  position:absolute; inset:0; padding:30px;
  opacity:0; transform:translateY(10px);
  transition:opacity .3s ease, transform .3s ease;
  pointer-events:none;
  overflow:hidden;
}
.kl .ft-panel.on { opacity:1; transform:translateY(0); pointer-events:auto; }

/* Panel: Generate */
.kl .fp-toprow { display:flex; align-items:center; justify-content:space-between; margin-bottom:20px; }
.kl .fp-title  { font-size:12px; font-weight:700; letter-spacing:.08em; text-transform:uppercase; color:var(--text3); }
.kl .fp-model  { font-size:12px; font-weight:600; padding:5px 10px; border-radius:6px; background:rgba(255,255,255,.06); border:1px solid var(--border); color:var(--text); }
.kl .fp-inp    { width:100%; background:rgba(255,255,255,.05); border:1px solid var(--border2); border-radius:8px; padding:11px 14px; font-size:13px; color:var(--text); font-family:'Syne',sans-serif; margin-bottom:10px; display:block; }
.kl .fp-gbtn   { width:100%; background:var(--accent); border:none; border-radius:8px; padding:12px; font-size:13.5px; font-weight:700; color:#fff; font-family:'Syne',sans-serif; cursor:default; letter-spacing:.01em; }
.kl .fp-ok     { font-size:12px; color:#4ade80; margin-top:10px; display:flex; align-items:center; gap:7px; }
.kl .fp-divider{ border:none; border-top:1px solid var(--border); margin:20px 0; }
.kl .fp-metarow{ display:flex; align-items:center; justify-content:space-between; }
.kl .fp-meta-k { font-size:11px; color:var(--text3); font-weight:600; }
.kl .fp-meta-v { font-size:12px; font-weight:700; padding:4px 10px; border-radius:6px; background:rgba(255,255,255,.06); border:1px solid var(--border); }

/* Panel: Schedule */
.kl .fp-day-label { font-size:11px; font-weight:700; color:var(--text3); text-transform:uppercase; letter-spacing:.08em; margin:0 0 8px; }
.kl .fp-day-label + .fp-day-label { margin-top:18px; }
.kl .fp-sched-list { display:flex; flex-direction:column; gap:6px; }
.kl .fp-sched-item { display:flex; align-items:center; gap:12px; border:1px solid var(--border); border-radius:9px; padding:10px 14px; background:rgba(255,255,255,.02); }
.kl .fp-sched-time { font-size:12px; font-weight:700; color:var(--accent-lt); min-width:40px; }
.kl .fp-sched-dot  { width:7px; height:7px; border-radius:50%; background:#4ade80; flex-shrink:0; }
.kl .fp-sched-name { font-size:13px; font-weight:600; flex:1; }
.kl .fp-sched-site { font-size:11px; color:var(--text3); }
.kl .fp-sched-footer { margin-top:16px; padding:10px 14px; border-radius:8px; background:var(--accent-faint); border:1px solid rgba(167,139,250,.18); font-size:12px; color:var(--accent-lt); font-weight:600; }

/* Panel: Multi-site */
.kl .fp-site-count { font-size:14px; font-weight:600; color:var(--text2); margin-bottom:14px; }
.kl .fp-site-count b { color:var(--text); }
.kl .fp-sites { display:flex; flex-direction:column; gap:6px; margin-bottom:20px; }
.kl .fp-site-row { display:flex; align-items:center; gap:12px; border:1px solid var(--border); border-radius:9px; padding:10px 14px; background:rgba(255,255,255,.02); transition:background .15s; }
.kl .fp-site-row.sel { border-color:rgba(167,139,250,.25); background:rgba(124,58,237,.06); }
.kl .fp-site-chk { width:16px; height:16px; border-radius:4px; border:1px solid var(--border2); display:flex; align-items:center; justify-content:center; font-size:9px; flex-shrink:0; }
.kl .fp-site-row.sel .fp-site-chk { background:var(--accent); border-color:var(--accent); color:#fff; }
.kl .fp-site-name { font-size:13px; font-weight:600; flex:1; }
.kl .fp-site-stat { font-size:11px; }
.kl .fp-post-btn  { width:100%; background:var(--accent); border:none; border-radius:9px; padding:12px; font-size:13.5px; font-weight:700; color:#fff; font-family:'Syne',sans-serif; cursor:default; }

/* Panel: SEO */
.kl .fp-seo-header { display:flex; align-items:center; gap:20px; padding:18px; border:1px solid var(--border); border-radius:12px; background:rgba(255,255,255,.02); margin-bottom:22px; }
.kl .fp-seo-big { font-family:'Fraunces',serif; font-size:62px; font-weight:700; letter-spacing:-.04em; line-height:1; color:#4ade80; }
.kl .fp-seo-sub { font-size:13px; color:var(--text3); margin-top:2px; }
.kl .fp-seo-track-wrap { flex:1; }
.kl .fp-seo-track { height:6px; background:rgba(255,255,255,.07); border-radius:100px; overflow:hidden; margin-bottom:6px; }
.kl .fp-seo-fill  { height:100%; border-radius:100px; background:linear-gradient(90deg,#4ade80,#22d3ee); }
.kl .fp-seo-status{ font-size:12px; color:var(--text2); font-weight:600; }
.kl .fp-metrics { display:flex; flex-direction:column; gap:10px; }
.kl .fp-metric   { display:flex; align-items:center; gap:12px; }
.kl .fp-metric-name { font-size:13px; color:var(--text2); width:160px; flex-shrink:0; }
.kl .fp-metric-bar  { flex:1; height:5px; background:rgba(255,255,255,.06); border-radius:100px; overflow:hidden; }
.kl .fp-metric-fill { height:100%; border-radius:100px; }
.kl .fp-metric-val  { font-size:12px; font-weight:700; min-width:30px; text-align:right; }

/* ═══ HOW IT WORKS ═══ */
.kl .how-section { border-top:1px solid var(--border); }
.kl .how-center  { text-align:center; }
.kl .how-grid { display:grid; grid-template-columns:repeat(3,1fr); gap:48px; margin-top:60px; position:relative; }
.kl .how-grid::before { content:''; position:absolute; top:27px; left:15%; right:15%; height:1px; background:linear-gradient(to right, transparent, var(--border2), var(--border2), transparent); }
.kl .how-step { text-align:center; }
.kl .how-circle { width:56px; height:56px; border-radius:50%; border:1px solid var(--border2); background:var(--bg2); display:flex; align-items:center; justify-content:center; margin:0 auto 22px; position:relative; z-index:1; }
.kl .how-n  { font-family:'Fraunces',serif; font-size:20px; font-weight:700; color:var(--accent-lt); }
.kl .how-title { font-size:16px; font-weight:700; margin-bottom:10px; }
.kl .how-desc  { font-size:14px; color:var(--text2); line-height:1.65; }

/* ═══ FOR WHO ═══ */
.kl .who-section { border-top:1px solid var(--border); }
.kl .who-grid { display:grid; grid-template-columns:repeat(2,1fr); gap:24px; margin-top:56px; }
.kl .who-card { border:1px solid var(--border); border-radius:14px; padding:36px; position:relative; overflow:hidden; transition:border-color .25s, transform .25s; }
.kl .who-card:hover { border-color:rgba(167,139,250,.28); transform:translateY(-2px); }
.kl .who-stripe { position:absolute; top:0; left:0; right:0; height:2px; background:linear-gradient(90deg, var(--accent) 0%, transparent 70%); }
.kl .who-icon  { font-size:34px; margin-bottom:18px; }
.kl .who-title { font-size:20px; font-weight:700; margin-bottom:10px; }
.kl .who-desc  { font-size:14px; color:var(--text2); line-height:1.65; margin-bottom:20px; }
.kl .who-tags  { display:flex; flex-wrap:wrap; gap:8px; }
.kl .who-tag   { padding:4px 12px; border-radius:100px; font-size:12px; font-weight:600; color:var(--accent-lt); border:1px solid rgba(167,139,250,.2); background:var(--accent-faint); }

/* ═══ PRICING ═══ */
.kl .pricing-section { border-top:1px solid var(--border); }
.kl .pricing-center  { text-align:center; }
.kl .pricing-grid { display:grid; grid-template-columns:repeat(3,1fr); gap:22px; margin-top:56px; }
.kl .price-card { border:1px solid var(--border); border-radius:14px; padding:34px; background:var(--bg2); position:relative; transition:transform .25s; }
.kl .price-card:hover { transform:translateY(-3px); }
.kl .price-card.pop { border-color:rgba(167,139,250,.35); background:linear-gradient(180deg,rgba(124,58,237,.07) 0%,var(--bg2) 50%); }
.kl .pop-badge { position:absolute; top:-12px; left:50%; transform:translateX(-50%); background:var(--accent); color:#fff; font-size:10.5px; font-weight:700; letter-spacing:.07em; text-transform:uppercase; padding:3px 14px; border-radius:100px; white-space:nowrap; }
.kl .price-tier { font-size:12px; font-weight:700; letter-spacing:.09em; text-transform:uppercase; color:var(--text3); margin-bottom:14px; }
.kl .price-val  { font-family:'Fraunces',serif; font-size:42px; font-weight:700; letter-spacing:-.03em; line-height:1; margin-bottom:4px; }
.kl .price-val sub { font-size:17px; font-weight:400; color:var(--text2); vertical-align:middle; }
.kl .price-per  { font-size:12px; color:var(--text3); margin-bottom:26px; }
.kl .price-sep  { border:none; border-top:1px solid var(--border); margin:22px 0; }
.kl .price-list { list-style:none; display:flex; flex-direction:column; gap:11px; margin-bottom:26px; }
.kl .price-li   { display:flex; align-items:flex-start; gap:10px; font-size:13.5px; color:var(--text2); }
.kl .price-chk  { color:#4ade80; font-size:12px; margin-top:2px; flex-shrink:0; }
.kl .price-cta  { width:100%; padding:12px; border-radius:9px; font-size:13.5px; font-weight:600; font-family:'Syne',sans-serif; cursor:pointer; border:none; transition:all .2s; }
.kl .price-cta.ghost { background:transparent; color:var(--text); border:1px solid var(--border2); }
.kl .price-cta.ghost:hover { border-color:var(--border3); background:rgba(255,255,255,.04); }
.kl .price-cta.solid { background:var(--accent); color:#fff; }
.kl .price-cta.solid:hover { background:#6d28d9; box-shadow:0 6px 24px var(--accent-glow); }

/* ═══ FAQ ═══ */
.kl .faq-section { border-top:1px solid var(--border); }
.kl .faq-list { margin-top:52px; }
.kl .faq-row  { border-bottom:1px solid var(--border); }
.kl .faq-btn  { width:100%; display:flex; align-items:center; justify-content:space-between; gap:20px; padding:22px 0; background:none; border:none; cursor:pointer; text-align:left; font-family:'Syne',sans-serif; font-size:15.5px; font-weight:600; color:var(--text); transition:color .2s; }
.kl .faq-btn:hover { color:var(--accent-lt); }
.kl .faq-icon { font-size:18px; color:var(--text3); flex-shrink:0; transition:transform .22s, color .2s; }
.kl .faq-icon.open { transform:rotate(45deg); color:var(--accent-lt); }
.kl .faq-body { font-size:15px; color:var(--text2); line-height:1.72; overflow:hidden; max-height:0; transition:max-height .32s ease, padding-bottom .32s ease; }
.kl .faq-body.open { max-height:280px; padding-bottom:22px; }

/* ═══ CTA ═══ */
.kl .cta-section { border-top:1px solid var(--border); padding:96px 0; }
.kl .cta-box { border:1px solid var(--border2); border-radius:20px; padding:88px 60px; text-align:center; position:relative; overflow:hidden; background:var(--bg2); }
.kl .cta-halo { position:absolute; top:0; left:50%; transform:translateX(-50%); width:500px; height:260px; pointer-events:none; background:radial-gradient(ellipse at top, rgba(124,58,237,.18) 0%, transparent 65%); }
.kl .cta-box h2 { font-family:'Fraunces',serif; font-size:clamp(30px,5vw,54px); font-weight:700; letter-spacing:-.025em; line-height:1.1; margin-bottom:16px; position:relative; }
.kl .cta-box p  { font-size:17px; color:var(--text2); margin-bottom:36px; position:relative; }
.kl .cta-row    { display:flex; align-items:center; justify-content:center; gap:10px; flex-wrap:wrap; position:relative; }
.kl .cta-note   { font-size:12px; color:var(--text3); margin-top:14px; }

/* ═══ FOOTER ═══ */
.kl footer { border-top:1px solid var(--border); padding:44px 0; }
.kl .footer-inner { max-width:1120px; margin:0 auto; padding:0 32px; display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:20px; }
.kl .footer-logo { display:flex; align-items:center; gap:8px; font-weight:800; font-size:15px; color:var(--text); text-decoration:none; }
.kl .footer-logo-mark { width:24px; height:24px; border-radius:6px; background:var(--accent); display:flex; align-items:center; justify-content:center; font-size:11px; }
.kl .footer-links { display:flex; gap:26px; flex-wrap:wrap; }
.kl .footer-links a { font-size:13px; color:var(--text3); text-decoration:none; transition:color .2s; }
.kl .footer-links a:hover { color:var(--text2); }
.kl .footer-copy { font-size:12px; color:var(--text3); }

/* ═══ RESPONSIVE ═══ */
@media (max-width: 960px) {
  .kl .hero-split { grid-template-columns:1fr; }
  .kl .hero-text  { text-align:center; }
  .kl .hero-sub   { margin-left:auto; margin-right:auto; }
  .kl .hero-actions { justify-content:center; }
  .kl .hero-stats { justify-content:center; }
  .kl .hero-visual { display:none; }
}
@media (max-width: 860px) {
  .kl .container, .kl .container-sm { padding:0 20px; }
  .kl .nav-links  { display:none; }
  .kl .problem-grid  { grid-template-columns:1fr; }
  .kl .ft-wrap    { grid-template-columns:1fr; }
  .kl .who-grid   { grid-template-columns:1fr; }
  .kl .pricing-grid  { grid-template-columns:1fr; }
  .kl .how-grid   { grid-template-columns:1fr; gap:32px; }
  .kl .how-grid::before { display:none; }
  .kl .cta-box    { padding:52px 28px; }
  .kl .footer-inner { flex-direction:column; align-items:flex-start; }
}
`

/* ─────────────────────────────────────────
   DATA
───────────────────────────────────────── */
const FEATURES = [
  {
    icon: '✦',
    name: 'Generate Artikel + Gambar',
    hint: 'AI tulis artikel SEO lengkap dari satu keyword.',
  },
  {
    icon: '📅',
    name: 'Penjadwalan Otomatis',
    hint: 'Set jadwal sekali, konten posting sendiri tiap hari.',
  },
  {
    icon: '🌐',
    name: 'Multi-Website Sekaligus',
    hint: 'Satu generate, posting ke banyak website sekaligus.',
  },
  {
    icon: '📊',
    name: 'SEO Score Bawaan',
    hint: 'Nilai SEO otomatis sebelum artikel dipublish.',
  },
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
  {
    icon: '🔍', title: 'Blogger & Pemain SEO',
    desc: 'Punya banyak niche site tapi kewalahan ngisi konten? Kabar bantu kamu scale ke semua website sekaligus dengan konten yang konsisten tiap hari.',
    tags: ['Niche Site', 'Affiliate', 'SEO Content', 'PBN'],
  },
  {
    icon: '📰', title: 'Content Agency & Media',
    desc: 'Kelola konten untuk banyak klien sekaligus? Kabar mempersingkat produksi dari jam menjadi menit, dengan kualitas yang bisa dikontrol penuh.',
    tags: ['Digital Agency', 'Content Marketing', 'Media Online'],
  },
]

const PLANS = [
  { tier: 'Starter',  price: 'Gratis',      per: 'Selamanya', pop: false, features: ['5 artikel per bulan', '1 website', 'Generate gambar AI', 'SEO Score', 'Support komunitas'], cta: 'Mulai Gratis',        style: 'ghost' },
  { tier: 'Pro',      price: 'Rp 99.000',   per: 'per bulan', pop: true,  features: ['50 artikel per bulan', '5 website', 'Penjadwalan otomatis', 'Semua AI model', 'Prioritas generate', 'Email support'], cta: 'Coba 7 Hari Gratis', style: 'solid' },
  { tier: 'Agency',   price: 'Rp 299.000',  per: 'per bulan', pop: false, features: ['Artikel tak terbatas', 'Website tak terbatas', 'Penjadwalan otomatis', 'Semua AI model', 'API access', 'Priority support'], cta: 'Hubungi Kami', style: 'ghost' },
]

const FAQS = [
  { q: 'Apakah konten AI akan dihukum Google?', a: 'Google tidak menghukum konten AI — yang dinilai adalah kualitas dan relevansi. Kabar mengoptimalkan artikel untuk SEO dan readability sehingga memenuhi standar Google. Ribuan website menggunakan konten AI tanpa masalah selama kualitasnya terjaga.' },
  { q: 'Website atau platform apa yang bisa diconnect?', a: 'Kabar mendukung berbagai platform blog dan CMS. Kamu bisa menambahkan dan mengelola beberapa "Produk" (website) langsung dari dashboard, lalu memilih satu atau banyak website sekaligus sebagai target posting.' },
  { q: 'Apakah bisa generate konten Bahasa Indonesia?', a: 'Ya, Kabar mendukung multi-bahasa. Kamu bisa generate konten dalam Bahasa Indonesia, Inggris, atau bahasa lain. Kabar dioptimalkan untuk konteks dan gaya penulisan Indonesia.' },
  { q: 'Bagaimana cara kerja penjadwalan posting?', a: 'Setelah generate artikel, pilih mode "Terjadwal" di Konfigurasi Posting, tentukan tanggal dan waktu, lalu pilih website target. Kabar akan otomatis memposting pada waktu yang kamu tentukan.' },
  { q: 'Apa bedanya model AI yang tersedia?', a: 'Kabar mengintegrasikan berbagai model via OpenRouter — termasuk Gemini, GPT-4, Claude, dan lainnya. Kamu bisa memilih model yang paling sesuai dengan kebutuhan kualitas dan anggaran.' },
]

/* ─────────────────────────────────────────
   FEATURE PREVIEW PANELS
───────────────────────────────────────── */
function PanelGenerate() {
  return (
    <>
      <div className="fp-toprow">
        <span className="fp-title">Generate Konten</span>
        <span className="fp-model">Gemini 2.5 Flash ▾</span>
      </div>
      <div className="fp-inp" style={{ display: 'block' }}>
        10 Tips SEO Terbaik untuk Meningkatkan Traffic 2025
      </div>
      <button className="fp-gbtn">✦&nbsp;&nbsp;Generate Artikel</button>
      <div className="fp-ok"><span>●</span> Artikel siap · 1.240 kata</div>
      <hr className="fp-divider" />
      <div className="fp-metarow">
        <span className="fp-meta-k">SEO Score</span>
        <span className="fp-meta-v" style={{ color: '#4ade80' }}>87 / 100</span>
      </div>
      <div className="fp-metarow" style={{ marginTop: 10 }}>
        <span className="fp-meta-k">Mode Posting</span>
        <span className="fp-meta-v">Terjadwal · Besok 09:00</span>
      </div>
    </>
  )
}

function PanelSchedule() {
  return (
    <>
      <p className="fp-day-label">Hari ini</p>
      <div className="fp-sched-list">
        {[
          { t: '09:00', name: '10 Tips SEO Terbaik 2025', site: 'techblog.id' },
          { t: '15:00', name: 'Review Laptop Gaming Terbaik', site: 'reviewgadget.com' },
        ].map(i => (
          <div key={i.t} className="fp-sched-item">
            <span className="fp-sched-time">{i.t}</span>
            <span className="fp-sched-dot" />
            <span className="fp-sched-name">{i.name}</span>
            <span className="fp-sched-site">{i.site}</span>
          </div>
        ))}
      </div>
      <p className="fp-day-label" style={{ marginTop: 18 }}>Besok</p>
      <div className="fp-sched-list">
        {[
          { t: '10:00', name: 'Strategi Affiliate Marketing 2025', site: 'affiliateku.net' },
          { t: '14:00', name: 'Cara Monetisasi Blog dari Nol', site: 'tutorialseo.id' },
        ].map(i => (
          <div key={i.t} className="fp-sched-item">
            <span className="fp-sched-time">{i.t}</span>
            <span className="fp-sched-dot" />
            <span className="fp-sched-name">{i.name}</span>
            <span className="fp-sched-site">{i.site}</span>
          </div>
        ))}
      </div>
      <div className="fp-sched-footer">✓&nbsp; 8 artikel terjadwal bulan ini</div>
    </>
  )
}

function PanelMultiSite() {
  const sites = [
    { name: 'techblog.id',       sel: true  },
    { name: 'reviewgadget.com',  sel: true  },
    { name: 'tutorialseo.id',    sel: true  },
    { name: 'affiliateku.net',   sel: false },
    { name: 'produkfavorit.com', sel: false },
  ]
  return (
    <>
      <p className="fp-site-count"><b>3</b> dari 5 website dipilih</p>
      <div className="fp-sites">
        {sites.map(s => (
          <div key={s.name} className={`fp-site-row ${s.sel ? 'sel' : ''}`}>
            <div className="fp-site-chk">{s.sel ? '✓' : ''}</div>
            <span className="fp-site-name">{s.name}</span>
            <span className="fp-site-stat" style={{ color: s.sel ? '#4ade80' : 'var(--text3)' }}>
              {s.sel ? '● aktif' : '○'}
            </span>
          </div>
        ))}
      </div>
      <button className="fp-post-btn">📅&nbsp;&nbsp;Post ke 3 Website →</button>
    </>
  )
}

function PanelSEO() {
  const metrics = [
    { label: 'Keyword Density',   pct: 80,  color: '#4ade80' },
    { label: 'Struktur Heading',  pct: 100, color: '#4ade80' },
    { label: 'Meta Description',  pct: 65,  color: '#fbbf24' },
    { label: 'Jumlah Kata',       pct: 100, color: '#4ade80' },
  ]
  return (
    <>
      <div className="fp-seo-header">
        <div>
          <div className="fp-seo-big">87</div>
          <div className="fp-seo-sub">dari 100</div>
        </div>
        <div className="fp-seo-track-wrap">
          <div className="fp-seo-track">
            <div className="fp-seo-fill" style={{ width: '87%' }} />
          </div>
          <div className="fp-seo-status">Skor Bagus — siap dipublish</div>
        </div>
      </div>
      <div className="fp-metrics">
        {metrics.map(m => (
          <div key={m.label} className="fp-metric">
            <span className="fp-metric-name">{m.label}</span>
            <div className="fp-metric-bar">
              <div className="fp-metric-fill" style={{ width: `${m.pct}%`, background: m.color }} />
            </div>
            <span className="fp-metric-val" style={{ color: m.color }}>{m.pct}</span>
          </div>
        ))}
      </div>
    </>
  )
}

const PANELS = [PanelGenerate, PanelSchedule, PanelMultiSite, PanelSEO]

/* ─────────────────────────────────────────
   COMPONENT
───────────────────────────────────────── */
export default function Landing() {
  const [activeFeature, setActiveFeature] = useState(0)
  const [openFaq, setOpenFaq] = useState<number | null>(null)

  /* Scroll reveal */
  useEffect(() => {
    const obs = new IntersectionObserver(
      (entries) => entries.forEach(e => { if (e.isIntersecting) e.target.classList.add('in-view') }),
      { threshold: 0.1, rootMargin: '-40px 0px' }
    )
    document.querySelectorAll('[data-reveal]').forEach(el => obs.observe(el))
    return () => obs.disconnect()
  }, [])

  return (
    <div className="kl">
      <style dangerouslySetInnerHTML={{ __html: STYLES }} />
      <div className="grain" />

      {/* ── NAVBAR ── */}
      <nav>
        <div className="nav-inner">
          <a href="#" className="nav-logo">
            <div className="nav-logo-mark">✦</div>
            KABAR
          </a>
          <ul className="nav-links">
            <li><a href="#fitur">Fitur</a></li>
            <li><a href="#cara-kerja">Cara Kerja</a></li>
            <li><a href="#harga">Harga</a></li>
            <li><a href="#faq">Pertanyaan yang Sering Diajukan (FAQ)</a></li>
          </ul>
          <div className="nav-right">
            <a href="#" className="btn-ghost">Masuk</a>
            <a href="#" className="btn-primary">Mulai Gratis</a>
          </div>
        </div>
      </nav>

      {/* ── HERO ── */}
      <section className="hero">
        <div className="hero-dots" />
        <div className="hero-glow-tr" />
        <div className="container">
          <div className="hero-split">

            {/* LEFT — text */}
            <div className="hero-text">
              <div className="badge au d1">
                <span className="badge-dot" />
                Platform Konten AI · Dibuat untuk Indonesia
              </div>
              <h1 className="hero-h1 display au d2">
                Isi Semua<br />Website-mu.<br /><em>Otomatis.</em>
              </h1>
              <p className="hero-sub au d3">
                Generate artikel SEO berkualitas, buat gambar, dan posting terjadwal ke semua website-mu — tanpa menulis satu kata pun.
              </p>
              <div className="hero-actions au d4">
                <a href="#" className="btn-primary lg">✦&nbsp; Mulai Gratis</a>
                <a href="#cara-kerja" className="btn-ghost lg">Cara Kerja →</a>
              </div>
              <p className="hero-disclaimer au d4">
                Tidak perlu kartu kredit · Setup dalam 2 menit
              </p>
              <div className="hero-stats au d5">
                <div>
                  <span className="hero-stat-n">500+</span>
                  <span className="hero-stat-l">artikel/hari</span>
                </div>
                <div className="hero-stat-sep" />
                <div>
                  <span className="hero-stat-n">87</span>
                  <span className="hero-stat-l">rata‑rata SEO score</span>
                </div>
                <div className="hero-stat-sep" />
                <div>
                  <span className="hero-stat-n">5+</span>
                  <span className="hero-stat-l">model AI tersedia</span>
                </div>
              </div>
            </div>

            {/* RIGHT — visual */}
            <div className="hero-visual">
              <div className="hero-visual-glow" />

              {/* Floating badge: SEO score */}
              <div className="fbadge fb1 au d6"
                style={{ position: 'absolute', top: '10%', right: '-20px', zIndex: 10 }}>
                <span className="fbadge-icon">📋</span>
                <div>
                  <div className="fbadge-label">SEO Score</div>
                  <div className="fbadge-val" style={{ color: '#4ade80' }}>87</div>
                </div>
              </div>

              {/* Floating badge: article ready */}
              <div className="fbadge fb2 au d7"
                style={{ position: 'absolute', bottom: '28%', left: '-22px', zIndex: 10 }}>
                <span className="fbadge-icon">✅</span>
                <div>
                  <div className="fbadge-label">Artikel Siap</div>
                  <div className="fbadge-val">1.240 kata</div>
                </div>
              </div>

              {/* Floating badge: posted */}
              <div className="fbadge fb3 au d7"
                style={{ position: 'absolute', bottom: '10%', right: '8%', zIndex: 10 }}>
                <span className="fbadge-icon">🌐</span>
                <div>
                  <div className="fbadge-label">Posted ke</div>
                  <div className="fbadge-val">3 Website</div>
                </div>
              </div>

              {/* Browser mockup */}
              <div className="hero-mockup-tilt au d6">
                <div className="mwin">
                  <div className="mwin-bar">
                    <div className="mwin-dot" style={{ background: '#ff5f57' }} />
                    <div className="mwin-dot" style={{ background: '#ffbd2e' }} />
                    <div className="mwin-dot" style={{ background: '#28c840' }} />
                    <div className="mwin-url">app.kabar.id/generate</div>
                  </div>
                  <div className="mwin-body">
                    {/* sidebar */}
                    <div className="mwin-sb">
                      <div className="mwin-sb-sec">Menu</div>
                      {[['⊞','Dashboard'],['✦','Buat Konten',true],['📦','Produk'],['📄','Draft']].map(([ic,lb,on]) => (
                        <div key={String(lb)} className={`mwin-sb-item ${on ? 'on' : ''}`}>
                          <span>{ic}</span><span>{lb}</span>
                        </div>
                      ))}
                      <div className="mwin-sb-sec">Manajemen</div>
                      {[['📅','Jadwal'],['🕐','Riwayat'],['⚙️','Pengaturan']].map(([ic,lb]) => (
                        <div key={String(lb)} className="mwin-sb-item">
                          <span>{ic}</span><span>{lb}</span>
                        </div>
                      ))}
                    </div>
                    {/* main */}
                    <div className="mwin-main">
                      <div className="mwin-title">Buat Konten</div>
                      <div className="mwin-sub">Buat artikel + gambar dengan AI</div>
                      <div className="mwin-pill">
                        <span>📋 Skor SEO: <b>87</b></span>
                        <span>📝 Jumlah Kata: <span>1.240</span></span>
                      </div>
                      <div className="mwin-grid">
                        <div className="mwin-card">
                          <div className="mwin-clabel">✏️ Topik Artikel</div>
                          <div className="mwin-inp" style={{ display: 'block' }}>10 Tips SEO Terbaik untuk 2025</div>
                          <button className="mwin-gbtn">✦ Buat Artikel</button>
                          <div className="mwin-ok">● Artikel siap. Tambahkan gambar!</div>
                        </div>
                        <div className="mwin-card">
                          <div className="mwin-clabel">⚙️ Konfigurasi Posting</div>
                          <div className="mwin-modes">
                            <div className="mwin-mode">Langsung</div>
                            <div className="mwin-mode on">Jadwal</div>
                            <div className="mwin-mode">Draf</div>
                          </div>
                          <div className="mwin-tgt">Target: 3 website yang dipilih</div>
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

      {/* ── PROBLEM ── */}
      <section className="problem-section">
        <div className="container">
          <div data-reveal>
            <div className="section-label">Masalah yang Kita Selesaikan</div>
            <h2 className="section-title display">
              Kenapa content creator selalu<br />kehabisan waktu?
            </h2>
            <p className="section-desc" style={{ maxWidth: 480 }}>
              Bukan karena malas. Tapi karena cara lama memang tidak skalabel.
            </p>
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

      {/* ── FEATURES — TABBED ── */}
      <section className="features-section" id="fitur">
        <div className="container">
          <div data-reveal>
            <div className="section-label">Fitur Utama</div>
            <h2 className="section-title display">
              Satu platform untuk<br />semua kebutuhan kontenmu
            </h2>
          </div>

          <div className="ft-wrap" data-reveal data-delay="2">
            {/* Tabs */}
            <div className="ft-tabs">
              {FEATURES.map((f, i) => (
                <button
                  key={i}
                  className={`ft-tab ${activeFeature === i ? 'on' : ''}`}
                  onClick={() => setActiveFeature(i)}
                >
                  <span className="ft-tab-icon">{f.icon}</span>
                  <div className="ft-tab-name">{f.name}</div>
                  <div className="ft-tab-hint">{f.hint}</div>
                </button>
              ))}
            </div>

            {/* Preview */}
            <div className="ft-preview">
              <div className="ft-bar" />
              {PANELS.map((Panel, i) => (
                <div key={i} className={`ft-panel ${activeFeature === i ? 'on' : ''}`}>
                  <Panel />
                </div>
              ))}
            </div>
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
      <section className="who-section" id="untuk-siapa">
        <div className="container">
          <div data-reveal>
            <div className="section-label">Untuk Siapa</div>
            <h2 className="section-title display">
              Kabar dibuat untuk kamu<br />yang serius di konten
            </h2>
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

      {/* ── PRICING ── */}
      <section className="pricing-section" id="harga">
        <div className="container">
          <div className="pricing-center" data-reveal>
            <div className="section-label" style={{ justifyContent: 'center' }}>Harga</div>
            <h2 className="section-title display">Mulai gratis, scale<br />sesuai kebutuhanmu</h2>
            <p className="section-desc">Tidak ada kontrak panjang. Upgrade atau downgrade kapan saja.</p>
          </div>
          <div className="pricing-grid" data-reveal data-delay="2">
            {PLANS.map(plan => (
              <div key={plan.tier} className={`price-card ${plan.pop ? 'pop' : ''}`}>
                {plan.pop && <div className="pop-badge">Paling Populer</div>}
                <div className="price-tier">{plan.tier}</div>
                <div className="price-val display">
                  {plan.price}{plan.per === 'per bulan' && <sub>/bln</sub>}
                </div>
                <div className="price-per">{plan.per}</div>
                <hr className="price-sep" />
                <ul className="price-list">
                  {plan.features.map(f => (
                    <li key={f} className="price-li"><span className="price-chk">✓</span>{f}</li>
                  ))}
                </ul>
                <button className={`price-cta ${plan.style}`}>{plan.cta}</button>
              </div>
            ))}
          </div>
        </div>
      </section>

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
              <div className="cta-halo" />
              <h2 className="display">Siap isi semua website-mu<br />tanpa kerja extra?</h2>
              <p>Mulai hari ini. Gratis. Tanpa kartu kredit.</p>
              <div className="cta-row">
                <a href="#" className="btn-primary lg">✦&nbsp; Mulai Gratis Sekarang</a>
                <a href="#" className="btn-ghost lg">Lihat Demo →</a>
              </div>
              <p className="cta-note">Bergabung dengan ratusan content creator Indonesia 🇮🇩</p>
            </div>
          </div>
        </div>
      </section>

      {/* ── FOOTER ── */}
      <footer>
        <div className="footer-inner">
          <a href="#" className="footer-logo">
            <div className="footer-logo-mark">✦</div>
            KABAR
          </a>
          <div className="footer-links">
            <a href="#">Fitur</a>
            <a href="#">Harga</a>
            <a href="#">Blog</a>
            <a href="#">Kebijakan Privasi</a>
            <a href="#">Syarat &amp; Ketentuan</a>
          </div>
          <p className="footer-copy">© 2025 Kabar · Dibuat di Indonesia 🇮🇩</p>
        </div>
      </footer>
    </div>
  )
}