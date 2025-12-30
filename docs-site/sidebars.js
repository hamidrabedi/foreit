/**
 * Sidebar configuration for forge documentation
 * Simplified structure with Getting Started and Full Guides
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // Getting Started sidebar - Concepts, Quick Start, and Core Logic
  gettingStarted: [
    {
      type: 'doc',
      id: 'introduction',
      label: 'Introduction',
    },
    {
      type: 'category',
      label: 'Concepts',
      items: [
        'learn/what-is-forge',
        'learn/architecture',
        'learn/request-lifecycle',
        'learn/code-generation',
      ],
    },
    {
      type: 'category',
      label: 'Quick Start',
      items: [
        'getting-started/installation',
        'getting-started/hello-world',
        'getting-started/first-api',
        'getting-started/project-structure',
        'getting-started/quickstart',
      ],
    },
    {
      type: 'category',
      label: 'Core Logic',
      items: [
        'learn/extending-forge',
      ],
    },
  ],

  // Full Guides sidebar - All guides, API reference, examples, advanced, and contributing
  fullGuides: [
    {
      type: 'category',
      label: 'Tutorials',
      items: [
        'tutorials/01-overview',
        'tutorials/02-getting-started',
        'tutorials/03-admin-interface',
      ],
    },
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
        'guides/common-patterns',
        'guides/best-practices',
      ],
    },
    {
      type: 'category',
      label: 'API Reference',
      items: [
        'api-reference/schema',
        'api-reference/queryset',
        'api-reference/manager',
        'api-reference/fields',
        'api-reference/relations',
        'api-reference/hooks',
      ],
    },
    {
      type: 'category',
      label: 'Examples',
      items: [
        'examples/blog',
        'examples/ecommerce',
        'examples/library',
      ],
    },
    {
      type: 'category',
      label: 'Advanced',
      items: [
        'advanced/performance',
        'advanced/plugins',
        'advanced/custom-fields',
        'advanced/code-generation',
      ],
    },
    {
      type: 'category',
      label: 'Contributing',
      items: [
        'contributing/architecture',
        'contributing/development',
        'contributing/changelog',
        'contributing/roadmap',
      ],
    },
  ],
};

module.exports = sidebars;
