<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=0:0d1117,50:1a1a2e,100:16213e&height=220&section=header&text=Chaturbate%20DVR&fontSize=80&fontColor=e94560&fontAlignY=35&desc=Record%20live%20streams%20automatically%20%E2%80%94%20Docker%20%7C%20Cloud%20%7C%20Local&descSize=18&descAlignY=55&animation=fadeIn" width="100%"/>

[![Stars](https://img.shields.io/github/stars/vasud3v/chaturbate-recorder?style=flat&color=e94560&logo=github)](https://github.com/vasud3v/chaturbate-recorder/stargazers)
[![Forks](https://img.shields.io/github/forks/vasud3v/chaturbate-recorder?style=flat&color=0f3460&logo=github)](https://github.com/vasud3v/chaturbate-recorder/network/members)
[![Issues](https://img.shields.io/github/issues/vasud3v/chaturbate-recorder?style=flat&color=533483&logo=github)](https://github.com/vasud3v/chaturbate-recorder/issues)
[![License](https://img.shields.io/github/license/vasud3v/chaturbate-recorder?style=flat&color=00b894&logo=opensourceinitiative&logoColor=white)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-2496ED?style=flat&logo=docker&logoColor=white)](docker-compose.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/vasud3v/chaturbate-recorder?style=flat&color=00ADD8&logo=go&logoColor=white)](go.mod)

**An actively maintained, feature-packed DVR** — multi-channel recording, Cloudflare bypass, auto-uploads, and a slick web dashboard. Run it anywhere.

[Quick Start](#rocket-quick-start) &nbsp;·&nbsp; [Features](#zap-features) &nbsp;·&nbsp; [Docs](#book-documentation) &nbsp;·&nbsp; [Support](#heart-support)

</div>

---

## :zap: Features

<table>
<tr>
<td width="50%">

### :movie_camera: Recording
- Multi-channel simultaneous capture
- HLS `.ts` + LL-HLS `.m4s` support
- Auto-split by duration or file size
- ffmpeg compression to `.mkv`

</td>
<td width="50%">

### :globe_with_meridians: Deployment
- **Docker Compose** — one command setup
- **GitHub Actions** — free cloud recording
- **Binary** — single portable executable
- **Web UI** — manage everything from browser

</td>
</tr>
<tr>
<td>

### :shield: Cloudflare Bypass
- Byparr integration with load balancing
- Auto cookie refresh daemon
- Proxy support (SOCKS5/HTTP)

</td>
<td>

### :cloud: Uploads & Storage
- 6+ hosting providers in parallel
- Thumbnail & sprite generation
- Supabase metadata storage
- Browse everything in the dashboard

</td>
</tr>
</table>

---

## :rocket: Quick Start

```bash
git clone https://github.com/vasud3v/chaturbate-recorder.git
cd chaturbate-recorder
cp .env.example .env
docker compose up -d --build
```

Open **http://localhost:8080** — add channels, hit record. That's it.

---

## :camera: Dashboard

<p align="center">
  <img src="docs/images/dashboard.png" alt="Dashboard" width="85%" />
</p>

---

## :book: Documentation

<details>
<summary><b>GitHub Actions Setup</b></summary>
<br>

Run the recorder on GitHub-hosted runners — no server needed.

1. Fork this repo
2. Add repository secrets (`Settings` → `Secrets` → `Actions`)
3. Configure channels in `conf/channels.json`
4. Run workflow: `Actions` → `24/7 Recorder` → `Run workflow`

The dashboard is exposed via a **Cloudflare Tunnel** URL shown in the run summary.

</details>

<details>
<summary><b>Docker Services</b></summary>
<br>

| Service | Port | Purpose |
|---------|------|---------|
| `chaturbate-dvr` | `8080` | Recorder + Web UI |
| `byparr-lb` | `8191` | Cloudflare bypass |
| `cookie-refresher` | — | Auto cookie renewal |
| `uploader` | — | Background uploads |

</details>

<details>
<summary><b>CLI Mode</b></summary>
<br>

```bash
./chaturbate-dvr -u CHANNEL_USERNAME
./chaturbate-dvr -u CHANNEL -resolution 1080 -framerate 30 -max-duration 30
```

</details>

---

## :gear: Tech Stack

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img src="https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white" />
  <img src="https://img.shields.io/badge/ffmpeg-007808?style=for-the-badge&logo=ffmpeg&logoColor=white" />
  <img src="https://img.shields.io/badge/Supabase-3FCF8E?style=for-the-badge&logo=supabase&logoColor=white" />
  <img src="https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white" />
</p>

---

## :bar_chart: Repo Pulse

<p align="center">
  <img src="https://repobeats.axiom.co/api/embed/e8ec122a9f6217881b46ffee305942fc99b8c008.svg" alt="Repobeats" />
</p>

---

## :star: Star History

<a href="https://star-history.com/#vasud3v/chaturbate-recorder&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=vasud3v/chaturbate-recorder&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=vasud3v/chaturbate-recorder&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=vasud3v/chaturbate-recorder&type=Date" width="70%" />
 </picture>
</a>

---

## :heart: Support

If this project helps you, consider:

<p align="center">
  <a href="https://github.com/vasud3v/chaturbate-recorder/stargazers">
    <img src="https://img.shields.io/badge/⭐_Star_This_Repo-e94560?style=for-the-badge&logo=github&logoColor=white" alt="Star" height="40" />
  </a>
</p>

---

## :scroll: License

[MIT](LICENSE) — free to use, modify, and distribute.

<div align="center">
  <sub>Built with :heart: in 2026</sub>
</div>
