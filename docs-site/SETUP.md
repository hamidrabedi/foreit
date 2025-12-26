# Documentation Site Setup Guide

## Quick Start

1. **Install Dependencies**
   ```bash
   cd docs-site
   npm install
   ```

2. **Start Development Server**
   ```bash
   npm start
   ```
   This will start a local server at `http://localhost:3000`

3. **Build for Production**
   ```bash
   npm run build
   ```

## GitHub Pages Setup

### Initial Setup

1. Go to your repository settings on GitHub
2. Navigate to "Pages" in the left sidebar
3. Under "Source", select "GitHub Actions"
4. The site will automatically deploy when you push to the `main` branch

### Manual Deployment

If you need to deploy manually:

```bash
npm run build
npm run deploy
```

## Adding New Documentation

1. Create a new markdown file in the appropriate directory:
   - `docs/getting-started/` - Getting started guides
   - `docs/guides/` - Tutorial guides
   - `docs/reference/` - API reference
   - `docs/examples/` - Example applications
   - `docs/advanced/` - Advanced topics
   - `docs/contributing/` - Contributing guides

2. Add frontmatter:
   ```markdown
   ---
   sidebar_position: 1
   ---
   
   # Your Title
   ```

3. Update `sidebars.js` to include the new document

## Customization

### Changing Colors

Edit `src/css/custom.css` to change the color scheme.

### Changing Homepage

Edit `src/pages/index.jsx` and `src/components/HomepageFeatures/index.jsx`

### Adding Images

Place images in `static/img/` and reference them as `/img/filename.png`

## Troubleshooting

### Build Errors

- Make sure Node.js 18+ is installed
- Run `npm install` to ensure all dependencies are installed
- Check for syntax errors in markdown files

### Deployment Issues

- Ensure GitHub Pages is enabled in repository settings
- Check GitHub Actions workflow for errors
- Verify the workflow file path matches your repository structure

## Next Steps

- Read the [Docusaurus Documentation](https://docusaurus.io/docs)
- Check out [Docusaurus Examples](https://docusaurus.io/docs/showcase)
- Customize the theme and branding

