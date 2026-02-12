# Releasing ocmgr

## How GitHub Releases Work

```
You tag a commit  →  Pipeline triggers  →  Binaries built  →  Release published
    (v1.0.0)         (GitHub Actions)      (3 platforms)      (download page)
```

A **release** is a GitHub page tied to a **git tag** where you can attach downloadable files (your binaries). Users see it at `github.com/acchapm1/ocmgr-app/releases`.

## What's Included

### Pipeline — `.github/workflows/release.yml`

Triggers automatically when you push a tag starting with `v`. It:

- Builds binaries for **macOS arm64**, **Linux arm64**, and **Linux amd64**
- Packages each as a `.tar.gz` with SHA256 checksums
- Creates a GitHub Release with auto-generated release notes
- Tags containing `-rc`, `-beta`, or `-alpha` are marked as pre-releases

### Makefile — `make dist`

Cross-compiles all 3 platforms locally (useful for testing before tagging):

```bash
make dist
```

Output:

```
dist/ocmgr_v1.0.0_darwin_arm64.tar.gz
dist/ocmgr_v1.0.0_linux_arm64.tar.gz
dist/ocmgr_v1.0.0_linux_amd64.tar.gz
```

## How to Publish a Release

It's just two commands:

```bash
# 1. Tag the current commit
git tag v1.0.0

# 2. Push the tag to GitHub
git push origin v1.0.0
```

That's it! The pipeline runs automatically. After ~2 minutes, you'll see the release at:
`https://github.com/acchapm1/ocmgr-app/releases`

### Future releases

```bash
git tag v1.1.0
git push origin v1.1.0
```

### Pre-releases (for testing)

```bash
git tag v1.1.0-beta.1
git push origin v1.1.0-beta.1
```

### Install script integration

Once a release exists, `curl ... | bash` will automatically download the pre-built binary instead of building from source — the asset naming pattern (`ocmgr_v1.0.0_darwin_arm64.tar.gz`) matches what `install.sh` looks for.

## Setup Todo

- [ ] Commit the workflow file (`.github/workflows/release.yml`) and push to `main`
- [ ] Verify GitHub Actions is enabled for the repo
  - Go to **Settings → Actions → General** in the repo
  - Ensure "Allow all actions and reusable workflows" is selected
- [ ] Verify workflow permissions
  - Go to **Settings → Actions → General → Workflow permissions**
  - Select "Read and write permissions"
  - This allows the pipeline to create releases and upload assets
- [ ] Test with a pre-release tag
  - `git tag v0.1.0-beta.1 && git push origin v0.1.0-beta.1`
  - Check the **Actions** tab in GitHub to watch the pipeline run
  - Verify the release appears at `github.com/acchapm1/ocmgr-app/releases`
- [ ] Verify the install script picks up the release
  - `curl -sSL https://raw.githubusercontent.com/acchapm1/ocmgr-app/main/install.sh | bash`
  - It should download the pre-built binary instead of building from source
- [ ] Tag your first stable release
  - `git tag v1.0.0 && git push origin v1.0.0`
