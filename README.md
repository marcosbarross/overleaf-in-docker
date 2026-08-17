<div align="right">

[![Português](https://img.shields.io/badge/Language-Portugu%C3%AAs-green.svg)](./README.pt-br.md)

</div>

# Overleaf in Docker

A self-hosted Overleaf (Community Edition) environment containerized with Docker Compose, featuring **TeX Live Scheme-Full**, Microsoft fonts, `minted` (Pygments) syntax highlighting support, and an automated backup script written in Go.

---

## 🚀 Getting Started

### 1. Create external volumes

`docker-compose.yml` uses external volumes to ensure data persistence:

```bash
docker volume create overleaf_mongo_data
docker volume create overleaf_redis_data
docker volume create overleaf_overleaf_data
```

### 2. Start the services

Launch the containers in detached mode:

```bash
docker compose up -d
```

Overleaf will be accessible at `http://localhost`.

---

## 👤 First Access (Admin)

To create the initial admin user:

1. Navigate to: `http://localhost/launchpad`
2. Set up the administrator email and password.

---

## 📦 Custom Image

The image used (`marcosbarross/overleaf-on-docker:latest`) is based on `sharelatex/sharelatex:6.2.2` and includes:

- `scheme-full` (Full TeX Live distribution with all LaTeX packages)
- `ttf-mscorefonts-installer` (Times New Roman, Arial, etc.)
- `python3-pygments` (required by the `minted` package)

To rebuild the image locally:

```bash
docker build -t marcosbarross/overleaf-on-docker:latest .
```

---

## 💾 Data Backup

The repository includes a Go script ([main.go](./main.go)) to perform full backups:

- Compressed MongoDB dump (`.gz`)
- Archive of Overleaf data volume and uploads (`.tar.gz`)
- Automatic cleanup of backups older than 7 days

### Run backup:

```bash
go run main.go
```

Generated backup archives will be stored in the `./backups/` directory.

---

## 🤖 CI/CD (GitHub Actions)

This repository includes two automated workflows:

1. **Build Go Binary** ([.github/workflows/build-go.yml](./.github/workflows/build-go.yml)):
   - Triggered on changes to Go files (`**.go`, `go.mod`, etc.).
   - Compiles the backup binary and uploads it as a workflow artifact for easy download.

2. **Build and Push Docker Image** ([.github/workflows/docker-publish.yml](./.github/workflows/docker-publish.yml)):
   - Triggered on changes to the [Dockerfile](./Dockerfile).
   - Builds and publishes the image to Docker Hub (`marcosbarross/overleaf-on-docker:latest`).

> [!NOTE]
> For Docker Hub automated push to work, configure the following secrets in your GitHub repository (`Settings > Secrets and variables > Actions`):
> - `DOCKERHUB_USERNAME`: your Docker Hub username.
> - `DOCKERHUB_TOKEN`: Personal Access Token generated from Docker Hub.
