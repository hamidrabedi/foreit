// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'forge Framework',
  tagline: 'Django-like Go framework with type safety',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://hamidrabedi.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/foreit/',

  // GitHub pages deployment config.
  organizationName: 'hamidrabedi',
  projectName: 'foreit',

  onBrokenLinks: 'throw',

  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/forgego/forge/tree/main/newforge/docs-site/',
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/forge-social-card.jpg',
      navbar: {
        title: 'forge',
        logo: {
          alt: 'forge Framework Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'gettingStarted',
            position: 'left',
            label: 'Getting Started',
          },
          {
            type: 'docSidebar',
            sidebarId: 'guides',
            position: 'left',
            label: 'Guides',
          },
          {
            type: 'docSidebar',
            sidebarId: 'reference',
            position: 'left',
            label: 'Reference',
          },
          {
            type: 'docSidebar',
            sidebarId: 'examples',
            position: 'left',
            label: 'Examples',
          },
          {
            type: 'docSidebar',
            sidebarId: 'advanced',
            position: 'left',
            label: 'Advanced',
          },
          {
            type: 'docSidebar',
            sidebarId: 'contributing',
            position: 'left',
            label: 'Contributing',
          },
          {
            href: 'https://github.com/forgego/forge',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Documentation',
            items: [
              {
                label: 'Getting Started',
                to: '/docs/getting-started/installation',
              },
              {
                label: 'Guides',
                to: '/docs/guides/models',
              },
              {
                label: 'API Reference',
                to: '/docs/reference/schema',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/forgego/forge',
              },
              {
                label: 'Issues',
                href: 'https://github.com/forgego/forge/issues',
              },
              {
                label: 'Discussions',
                href: 'https://github.com/forgego/forge/discussions',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Examples',
                to: '/docs/examples/blog',
              },
              {
                label: 'Roadmap',
                to: '/docs/contributing/roadmap',
              },
              {
                label: 'Changelog',
                to: '/docs/contributing/changelog',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} forge Framework. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ['go', 'bash', 'yaml', 'sql'],
      },
      colorMode: {
        defaultMode: 'light',
        disableSwitch: false,
        respectPrefersColorScheme: true,
      },
    }),
};

module.exports = config;

