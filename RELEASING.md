# Releasing quick-branch

This project uses [GoReleaser](https://goreleaser.com/) to automate releases and publish to GitHub Releases and Homebrew.

## Prerequisites

1. Create a Homebrew tap repository at `github.com/rangoons/homebrew-tap`
   - Go to GitHub and create a new public repository named `homebrew-tap`
   - No need to add any files - GoReleaser will handle it

2. Install GoReleaser locally (for testing):
   ```bash
   brew install goreleaser
   ```

## Creating a Release

### Automated Release (via GitHub Actions)

1. **Tag a version:**
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

2. **GitHub Actions automatically:**
   - Builds binaries for Linux, macOS, and Windows (amd64 + arm64)
   - Creates a GitHub Release with binaries attached
   - Updates your Homebrew tap with the new cask

That's it! Users can now install via:
```bash
brew install rangoons/tap/quick-branch
```

### Manual Release (local testing)

To test the release process locally without publishing:

```bash
# Test the build without publishing
goreleaser release --snapshot --clean

# Check the dist/ folder for generated binaries
ls -la dist/
```

## Version Naming Convention

Follow semantic versioning (semver):
- `v0.1.0` - Initial release
- `v0.2.0` - New features
- `v0.2.1` - Bug fixes
- `v1.0.0` - First stable release

## First Release Checklist

- [ ] Create `homebrew-tap` repository on GitHub
- [ ] Update README.md with installation instructions
- [ ] Ensure LICENSE file exists
- [ ] Test build locally with `goreleaser release --snapshot --clean`
- [ ] Create and push the first tag
- [ ] Verify GitHub Release was created
- [ ] Verify Homebrew cask was updated in tap repo (check `Casks/` directory)
- [ ] Test installation: `brew install rangoons/tap/quick-branch`

## Troubleshooting

### Homebrew tap push fails

If GoReleaser can't push to your Homebrew tap:
1. Ensure the `homebrew-tap` repository exists and is public
2. Check that GitHub Actions has write permissions (should be automatic)

### Build fails

Check the logs in GitHub Actions at:
`https://github.com/rangoons/quick-branch/actions`

### Testing Homebrew cask locally

```bash
# Install from your tap
brew install rangoons/tap/quick-branch

# Or test the cask directly
brew install --cask path/to/homebrew-tap/Casks/quick-branch.rb
```
