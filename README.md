> **⚠️ Fork Notice**: This is a fork of the original [Homebox](https://github.com/sysadminsmedia/homebox) project. I created this fork primarily to add a builds system for tracking skateboard and bike parts (with sport-tagged parts, position-based mounting, and per build part tracking), plus a few handy additions like a light/dark theme toggle and a kill-switch to disable the Builds feature. I'll try to keep this fork in sync with the original repo & continue to implement feats I find useful for my use case, but for full support and the most up-to-date version, please use the original repository.
>

<div align="center">
  <img src="/docs/src/assets/lilbox.svg" height="200"/>
</div>

<h1 align="center" style="margin-top: -10px;"> HomeBox </h1>
<p align="center" style="width: 100%;">
   <a href="https://homebox.software/en/">Docs</a>
   |
   <a href="https://demo.homebox.software">Demo</a>
   |
   <a href="https://discord.gg/aY4DCkpNA9">Discord</a>
</p>
<p align="center" style="width: 100%;">
    <img src="https://img.shields.io/github/check-runs/sysadminsmedia/homebox/main" alt="Github Checks"/>
    <img src="https://img.shields.io/github/license/sysadminsmedia/homebox"/>
    <img src="https://img.shields.io/github/v/release/sysadminsmedia/homebox?sort=semver&display_name=release"/>
    <img src="https://img.shields.io/weblate/progress/homebox?server=https%3A%2F%2Ftranslate.sysadminsmedia.com"/>
</p>
<p align="center" style="width: 100%;">
    <img src="https://img.shields.io/reddit/subreddit-subscribers/homebox"/>
    <img src="https://img.shields.io/mastodon/follow/110749314839831923?domain=infosec.exchange"/>
    <img src="https://img.shields.io/lemmy/homebox%40lemmy.world?label=lemmy"/>
</p>
<p align="center" style="width: 100%;">
  <a href="https://www.pikapods.com/pods?run=homebox"><img src="https://www.pikapods.com/static/run-button.svg"/></a>
</p>

## What is HomeBox

HomeBox is the inventory and organization system built for the Home User! With a focus on simplicity and ease of use, Homebox is the perfect solution for your home inventory, organization, and management needs. While developing this project, We've tried to keep the following principles in mind:

- 🧘 _Simple but Expandable_ - Homebox is designed to be simple and easy to use. No complicated setup or configuration required. But expandable to whatever level of infrastructure you want to put into it.
- 🚀 _Blazingly Fast_ - Homebox is written in Go, which makes it extremely fast and requires minimal resources to deploy. In general, idle memory usage is less than 50MB for the whole container.
- 📦 _Portable_ - Homebox is designed to be portable and run on anywhere. We use SQLite and an embedded Web UI to make it easy to deploy, use, and backup.

### Key Features
- 📇 Rich Organization - Organize your items into categories, locations, and tags. You can also create custom fields to store additional information about your items.
- 🔍 Powerful Search - Quickly find items in your inventory using the powerful search feature.
- 📸 Image Upload - Upload images of your items to make it easy to identify them.
- 📄 Document and Warranty Tracking - Keep track of important documents and warranties for your items.
- 💰 Purchase & Maintenance Tracking - Track purchase dates, prices, and maintenance schedules for your items.
- 📱 Responsive Design - Homebox is designed to work on any device, including desktops, tablets, and smartphones.
- 🛹 Builds System - Track complete builds for skateboards and bikes, mounting sport-tagged parts to specific positions (deck, trucks, wheels, bearings).
- 🏷️ Sport Tagging - Tag parts as skateboard, bike, or generic so they only appear in the relevant builds.

## Screenshots
![Login Screen](.github/screenshots/1.png)
![Dashboard](.github/screenshots/2.png)
![Item View](.github/screenshots/3.png)
![Create Item](.github/screenshots/9.png)
![Search](.github/screenshots/8.png)

You can also try the demo instances of Homebox:
- [Demo](https://demo.homebox.software)
- [Nightly](https://nightly.homebox.software)

## Quick Start (Podman)

```bash
# Save pepper to local file
openssl rand -base64 48 > hbox.pepper
chmod 400 hbox.pepper

# If using the rootless or hardened image, ensure the data
# folder is owned by the MAPPED host UID (not 65532).
# Under rootless Podman, container UID 65532 = host UID 165531
# (offset +99999 from the default rootless uid_map).
mkdir -p /path/to/data/folder
sudo chown 165531:165531 -R /path/to/data/folder

# Verify the container will see it as its own UID:
podman unshare ls -ld /path/to/data/folder   # must show 65532 65532

podman run -d \
  --name homebox \
  --restart unless-stopped \
  --publish 8085:7745 \
  --env TZ=Europe/London \
  --env HBOX_AUTH_API_KEY_PEPPER="$(tr -d '\n' < hbox.pepper)" \
  --volume /path/to/data/folder/:/data \
  ghcr.io/antebudimir/homebox:latest
```

**Note (rootless Podman):** The hardened image runs as UID `65532`. Under rootless Podman this maps to a different host UID — with the default `uid_map` the offset is `+99999`, so the host UID is `165531`. Chowning to `65532` directly will fail with `unable to open database file (14)`. Always verify with `podman unshare ls -ld <dir>` — it must report `65532`. If your `subuid` config differs, compute the correct host UID from `podman unshare cat /proc/self/uid_map`.

**Do not add `:Z` to the volume** — it relabels the directory and breaks the UID mapping.

---

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

To get started with code based contributions, please see our [contributing guide](https://homebox.software/en/contribute/).

If you are not a coder and can't help translate, you can still contribute financially. Financial contributions help us maintain the project and keep demos running.

## Help us Translate
We want to make sure that Homebox is available in as many languages as possible. If you are interested in helping us translate Homebox, please help us via our [Weblate instance](https://translate.sysadminsmedia.com/projects/homebox/).

[![Translation status](https://translate.sysadminsmedia.com/widget/homebox/multi-auto.svg)](https://translate.sysadminsmedia.com/engage/homebox/)

## Credits
- Original project by [@hay-kot](https://github.com/hay-kot)
- Logo by [@lakotelman](https://github.com/lakotelman)

### Contributors
<a href="https://github.com/sysadminsmedia/homebox/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=sysadminsmedia/homebox" />
</a>
