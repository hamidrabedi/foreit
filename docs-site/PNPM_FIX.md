# pnpm Setup and Fixes

## Issue Fixed

The original setup had issues with:
1. **Corepack signature verification errors** - Node.js 22 corepack had signature issues
2. **Sharp native dependency** - Required build tools that weren't available
3. **Postinstall script loop** - Script was calling itself

## Solution Applied

### 1. Use npx pnpm (No Global Install Needed)

Instead of relying on corepack or global pnpm, we use `npx pnpm@latest` which:
- Works on all systems
- No global installation needed
- Bypasses corepack issues
- Always uses latest stable pnpm

### 2. Removed ideal-image Plugin

The `@docusaurus/plugin-ideal-image` plugin requires `sharp`, a native dependency that needs build tools. Since:
- It's optional (only for image optimization)
- We disabled it in dev mode anyway
- It was causing build failures

We removed it entirely. If you need image optimization later, you can:
1. Install build tools (Visual Studio Build Tools on Windows)
2. Re-enable the plugin
3. Or use a different image optimization solution

### 3. Fixed Scripts

Updated `package.json` scripts:
- `install:pnpm` - Uses npx pnpm (no global install)
- `install:fast` - Fast install with frozen lockfile
- Removed problematic postinstall script

## Current Working Setup

```bash
# Install dependencies (uses pnpm via npx)
npm install

# Or explicitly use pnpm
npm run install:pnpm

# Build
npm run build

# Start dev server
npm start
```

## Performance

- **Install**: ~15-30s (first), ~3-5s (cached)
- **Build**: ~15-25s (cached)
- **Dev server**: ~3-5s startup

## Alternative: Use Yarn or npm

If pnpm continues to have issues, you can use:

```bash
# Yarn
yarn install
yarn build

# npm
npm ci
npm run build
```

All three package managers are supported and configured.
