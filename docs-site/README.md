# forge Documentation Site

This is the Docusaurus documentation site for the forge framework.

## Getting Started

### Prerequisites

- Node.js 18 or later
- npm or yarn

### Installation

```bash
cd docs-site
npm install
```

### Development

Start the development server:

```bash
npm start
```

This starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

### Build

Build the site for production:

```bash
npm run build
```

This generates static content into the `build` directory and can be served using any static contents hosting service.

### Deployment

The site is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The deployment is handled by GitHub Actions (see `.github/workflows/deploy.yml`).

## Project Structure

```
docs-site/
├── docs/                    # Documentation markdown files
│   ├── intro.md            # Landing page
│   ├── getting-started/    # Getting started guides
│   ├── guides/             # Tutorial guides
│   ├── reference/          # API reference
│   ├── examples/           # Example applications
│   ├── advanced/            # Advanced topics
│   └── contributing/       # Contributing guides
├── src/
│   ├── components/         # React components
│   ├── css/                # Custom CSS
│   └── pages/              # Custom pages
├── static/                  # Static assets
├── docusaurus.config.js     # Docusaurus configuration
├── sidebars.js             # Sidebar navigation
└── package.json            # Dependencies
```

## Adding Documentation

1. Create a new markdown file in the appropriate directory under `docs/`
2. Add frontmatter with metadata:

```markdown
---
sidebar_position: 1
---

# Your Document Title

Content here...
```

3. Update `sidebars.js` to include the new document in navigation

## Customization

### Theme

Customize the theme in `docusaurus.config.js` and `src/css/custom.css`.

### Homepage

Edit `src/pages/index.jsx` to customize the homepage.

## See Also

- [Docusaurus Documentation](https://docusaurus.io/docs)
- [forge Framework](../README.md)

