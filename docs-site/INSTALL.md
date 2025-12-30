# Installation Guide

## Recommended: pnpm (Fastest) ⚡

pnpm is the fastest package manager, offering:
- **2-3x faster** installs than npm
- **1.5x faster** than yarn
- **50% less disk space**
- Better dependency resolution

### Install Dependencies (No Global Install Needed)

```bash
# Fast install with pnpm (via npx)
npm install
# or
npm run install

# Or use frozen lockfile (CI/CD)
npm run install:fast
```

**Note**: pnpm is used automatically via npx, so no global installation is required. If you have pnpm installed globally, you can use `pnpm` commands directly.

### Build

```bash
# Development
npm start

# Production build
npm run build

# Fast build (no minification)
npm run build:fast
```

## Alternative: Yarn

If you prefer yarn:

```bash
# Install yarn (if not installed)
npm install -g yarn

# Install dependencies
yarn install

# Build
yarn build
```

## Fallback: npm

If neither pnpm nor yarn is available:

```bash
# Install dependencies
npm ci

# Build
npm run build
```

## Performance Comparison

| Package Manager | First Install | Cached Install | Speed |
|----------------|---------------|----------------|-------|
| **pnpm** | ~15-30s | ~3-5s | ⚡⚡⚡ Fastest |
| **yarn** | ~30-60s | ~5-10s | ⚡⚡ Fast |
| **npm** | ~60-120s | ~10-20s | ⚡ Standard |

## Troubleshooting

### Clear cache

```bash
# pnpm
pnpm store prune

# yarn
yarn cache clean

# npm
npm cache clean --force
```

### Reinstall

```bash
# Remove node_modules and reinstall
rm -rf node_modules
pnpm install
```

## Why pnpm?

1. **Fastest installs** - Parallel downloads, smart caching
2. **Disk efficient** - Symlinks instead of copies
3. **Strict** - Prevents phantom dependencies
4. **Industry standard** - Used by Vue, Svelte, and others
