/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  gettingStarted: [
    'intro',
    {
      type: 'category',
      label: 'Getting Started',
      items: [
        'getting-started/installation',
        'getting-started/quickstart',
        'getting-started/first-app',
      ],
    },
  ],

  guides: [
    {
      type: 'category',
      label: 'Guides',
      items: [
        'guides/models',
        'guides/queries',
        'guides/admin',
        'guides/rest-api',
        'guides/migrations',
        'guides/security',
        'guides/deployment',
      ],
    },
  ],

  reference: [
    {
      type: 'category',
      label: 'API Reference',
      items: [
        'reference/schema',
        'reference/queryset',
        'reference/manager',
        'reference/fields',
        'reference/relations',
        'reference/hooks',
      ],
    },
  ],

  examples: [
    {
      type: 'category',
      label: 'Examples',
      items: [
        'examples/blog',
        'examples/ecommerce',
        'examples/library',
      ],
    },
  ],

  advanced: [
    {
      type: 'category',
      label: 'Advanced Topics',
      items: [
        'advanced/code-generation',
        'advanced/plugins',
        'advanced/custom-fields',
        'advanced/performance',
      ],
    },
  ],

  contributing: [
    {
      type: 'category',
      label: 'Contributing',
      items: [
        'contributing/development',
        'contributing/architecture',
        'contributing/roadmap',
        'contributing/changelog',
      ],
    },
  ],
};

module.exports = sidebars;

