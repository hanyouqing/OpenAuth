# GitHub Actions Workflows

This directory contains GitHub Actions workflows for automated testing, building, and deployment.

## Workflows

### 1. Test and Coverage (`test.yml`)

Automatically runs all unit tests and generates coverage reports on every push and pull request.

**Triggers:**
- Push to `main`, `master`, or `develop` branches
- Pull requests targeting `main`, `master`, or `develop` branches
- Manual trigger via GitHub Actions UI (`workflow_dispatch`)

**Features:**
- Runs all Go unit tests with race detection
- Generates detailed coverage reports
- Posts coverage summary as PR comment
- Identifies uncovered code and packages with low coverage
- Uploads coverage artifacts for download

**PR Comments Include:**
- Overall test coverage percentage with status badge
- Coverage breakdown by package (table format)
- List of packages with coverage below 50%
- List of files with uncovered code
- Test pass/fail status

### 2. Docker Build and Push (`docker-build.yml`)

Builds Docker images for backend and frontend, and pushes them to Docker Hub and GitHub Container Registry (GHCR).

**Triggers:**
- Push to `main` or `master` branches
- Tags matching `v*` pattern (e.g., `v1.0.0`)
- Pull requests (builds only, no push)
- Manual trigger via GitHub Actions UI (`workflow_dispatch`)

**Manual Trigger Options:**
- `version`: Custom version tag (default: `latest`)
- `push_to_dockerhub`: Whether to push to Docker Hub (default: `true`)
- `push_to_ghcr`: Whether to push to GitHub Container Registry (default: `true`)
- `push_to_quay`: Whether to push to Quay.io (default: `false`)
- `push_to_aliyun`: Whether to push to Aliyun Container Registry (default: `false`)

**Features:**
- Multi-platform builds (linux/amd64, linux/arm64)
- Build caching for faster builds
- Automatic version tagging based on git tags
- Pushes to both Docker Hub and GHCR
- Embeds build metadata (version, git commit, build time)

**Required Secrets:**
- `DOCKERHUB_USERNAME`: Your Docker Hub username (required for Docker Hub)
- `DOCKERHUB_TOKEN`: Docker Hub access token (create at https://hub.docker.com/settings/security)

**Optional Secrets (for additional registries):**
- `QUAY_USERNAME`: Your Quay.io username (required for Quay.io)
- `QUAY_TOKEN`: Quay.io access token (create at https://quay.io/user/<username>?tab=settings)
- `ALIYUN_USERNAME`: Your Aliyun Container Registry username (required for Aliyun)
- `ALIYUN_PASSWORD`: Aliyun Container Registry password
- `ALIYUN_NAMESPACE`: Aliyun Container Registry namespace

**Image Tags:**
- `latest`: Latest build from main/master branch
- `v1.0.0`: Semantic version tags
- `v1.0`: Major.minor version tags
- `v1`: Major version tags
- Branch names: For feature branches
- PR numbers: For pull requests (build only)

## Setup Instructions

### 1. Configure Registry Secrets

1. Go to your repository Settings → Secrets and variables → Actions
2. Add the following secrets:

**Required (for Docker Hub):**
   - `DOCKERHUB_USERNAME`: Your Docker Hub username
   - `DOCKERHUB_TOKEN`: Docker Hub access token (create at https://hub.docker.com/settings/security)

**Optional (for additional registries):**
   - `QUAY_USERNAME`: Your Quay.io username
   - `QUAY_TOKEN`: Quay.io access token
   - `ALIYUN_USERNAME`: Your Aliyun Container Registry username
   - `ALIYUN_PASSWORD`: Aliyun Container Registry password
   - `ALIYUN_NAMESPACE`: Aliyun Container Registry namespace

### 2. Enable GitHub Actions

Workflows are automatically enabled. They will run on:
- Every push to protected branches
- Every pull request
- Manual triggers via the Actions tab

### 3. View Results

- **Test Results**: Check PR comments or the Actions tab
- **Docker Images**: 
  - Docker Hub: `https://hub.docker.com/r/<username>/openauth-backend`
  - GHCR: `ghcr.io/<username>/openauth-backend`
  - Quay.io: `quay.io/<username>/openauth-backend` (if configured)
  - Aliyun: `registry.cn-hangzhou.aliyuncs.com/<namespace>/openauth-backend` (if configured)

## Usage Examples

### Manual Test Run

1. Go to Actions tab in GitHub
2. Select "Test and Coverage" workflow
3. Click "Run workflow"
4. Select branch and click "Run workflow"

### Manual Docker Build

1. Go to Actions tab in GitHub
2. Select "Build and Push Docker Images" workflow
3. Click "Run workflow"
4. Configure options:
   - Version: `1.0.0` (optional, default: `latest`)
   - Push to Docker Hub: `true`/`false` (default: `true`)
   - Push to GHCR: `true`/`false` (default: `true`)
   - Push to Quay.io: `true`/`false` (default: `false`)
   - Push to Aliyun: `true`/`false` (default: `false`)
5. Click "Run workflow"

### Release a New Version

1. Create and push a git tag:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```
2. The workflow will automatically build and push images with version tags

## Notes

- Test workflow requires PostgreSQL and Redis services (automatically started)
- Docker builds use buildx for multi-platform support
- Coverage reports are uploaded as artifacts and can be downloaded from the Actions tab
- PR comments are automatically updated on each new commit to the PR
