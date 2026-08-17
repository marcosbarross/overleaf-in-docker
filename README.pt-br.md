<div align="right">

[![English](https://img.shields.io/badge/Language-English-blue.svg)](./README.md)

</div>

# Overleaf in Docker

Ambiente self-hosted do Overleaf (Community Edition) containerizado com Docker Compose, incluindo **TeX Live Scheme-Full**, fontes Microsoft, suporte a `minted` (Pygments) e script automatizado de backup em Go.

---

## 🚀 Como Executar

### 1. Criar os volumes externos

O `docker-compose.yml` utiliza volumes externos para garantir persistência dos dados:

```bash
docker volume create overleaf_mongo_data
docker volume create overleaf_redis_data
docker volume create overleaf_overleaf_data
```

### 2. Iniciar os serviços

Suba os containers em segundo plano:

```bash
docker compose up -d
```

O Overleaf estará disponível em `http://localhost`.

---

## 👤 Primeiro Acesso (Admin)

Para criar o primeiro usuário administrador:

1. Acesse: `http://localhost/launchpad`
2. Cadastre o e-mail e a senha do administrador.

---

## 📦 Imagem Customizada

A imagem utilizada (`marcosbarross/overleaf-on-docker:latest`) é baseada no `sharelatex/sharelatex:6.2.2` e adiciona:

- `scheme-full` (TeX Live completo com todos os pacotes LaTeX)
- `ttf-mscorefonts-installer` (Times New Roman, Arial, etc.)
- `python3-pygments` (para uso do pacote `minted`)

Caso queira reconstruir a imagem localmente:

```bash
docker build -t marcosbarross/overleaf-on-docker:latest .
```

---

## 💾 Backup dos Dados

O projeto inclui um script em Go ([main.go](./main.go)) para realizar o backup completo:

- Dump comprimido do MongoDB (`.gz`)
- Compactação do volume de dados e uploads do Overleaf (`.tar.gz`)
- Limpeza automática de backups com mais de 7 dias de retenção

### Executar backup:

```bash
go run main.go
```

Os arquivos gerados serão salvos no diretório `./backups/`.

---

## 🤖 CI/CD (GitHub Actions)

O repositório possui dois fluxos automatizados:

1. **Build do Binário Go** ([.github/workflows/build-go.yml](./.github/workflows/build-go.yml)):
   - Disparado em alterações de arquivos Go (`**.go`, `go.mod`, etc.).
   - Compila o binário de backup e o disponibiliza como artefato para download.

2. **Build e Push da Imagem Docker** ([.github/workflows/docker-publish.yml](./.github/workflows/docker-publish.yml)):
   - Disparado em alterações no [Dockerfile](./Dockerfile).
   - Realiza o build e publica a imagem no Docker Hub (`marcosbarross/overleaf-on-docker:latest`).

> [!NOTE]
> Para o push no Docker Hub funcionar, adicione os seguintes Secrets nas configurações do repositório GitHub (`Settings > Secrets and variables > Actions`):
> - `DOCKERHUB_USERNAME`: seu usuário do Docker Hub.
> - `DOCKERHUB_TOKEN`: Personal Access Token gerado no Docker Hub.
