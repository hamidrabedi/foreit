import React from 'react';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import SEOHead from '@site/src/components/SEOHead';
import styles from './index.module.css';

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <>
      <SEOHead
        title="forge - Django-like Go Framework"
        description="Build web applications in Go with Django's developer experience. Type-safe ORM, auto-generated admin, REST API framework, and code generation."
        keywords={[
          'go framework',
          'golang framework',
          'django go',
          'type-safe orm',
          'go web framework',
          'forge framework',
        ]}
        url="/"
      />
      <Layout
        title={`${siteConfig.title}`}
        description="Django-like Go framework with type safety">
        <div className={styles.hero}>
          <div className={styles.heroContainer}>
            <div className={styles.heroContent}>
              <h1 className={styles.heroTitle}>
                Django-like Go Framework
              </h1>
              <p className={styles.heroSubtitle}>
                Build web applications in Go with Django's developer experience.
                Type-safe ORM, auto-generated admin, REST API framework, and code generation.
              </p>
              <div className={styles.heroButtons}>
                <Link
                  className={styles.buttonPrimary}
                  to="/docs/getting-started/installation">
                  Get Started →
                </Link>
                <Link
                  className={styles.buttonSecondary}
                  to="/docs/getting-started/hello-world">
                  Quick Start
                </Link>
              </div>
            </div>
            <div className={styles.heroCode}>
              <div className={styles.codeBlock}>
                <div className={styles.codeHeader}>
                  <span className={styles.codeDot}></span>
                  <span className={styles.codeDot}></span>
                  <span className={styles.codeDot}></span>
                  <span className={styles.codeTitle}>models/post.go</span>
                </div>
                <pre className={styles.codeContent}>
{`type Post struct {
    schema.BaseSchema
}

func (Post) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").
            Primary().AutoIncrement().Build(),
        schema.String("title").
            Required().MaxLength(200).Build(),
        schema.Text("content").Required().Build(),
        schema.Bool("published").
            Default(false).Build(),
    }
}`}
                </pre>
              </div>
            </div>
          </div>
        </div>

        <div className={styles.features}>
          <div className={styles.container}>
            <h2 className={styles.sectionTitle}>Why forge?</h2>
            <p className={styles.sectionSubtitle}>
              Everything you need to build production-ready web applications
            </p>
            <div className={styles.featuresGrid}>
              <div className={styles.featureCard}>
                <div className={styles.featureIcon}>🔒</div>
                <h3 className={styles.featureTitle}>Type-Safe ORM</h3>
                <p className={styles.featureDescription}>
                  Full Django ORM features with compile-time type checking. Write queries
                  that your compiler validates before runtime.
                </p>
              </div>
              <div className={styles.featureCard}>
                <div className={styles.featureIcon}>⚡</div>
                <h3 className={styles.featureTitle}>Code Generation</h3>
                <p className={styles.featureDescription}>
                  AST-based code generation for models, managers, and querysets.
                  Generate type-safe code automatically from your schemas.
                </p>
              </div>
              <div className={styles.featureCard}>
                <div className={styles.featureIcon}>🎛️</div>
                <h3 className={styles.featureTitle}>Auto Admin</h3>
                <p className={styles.featureDescription}>
                  Django-style admin interface auto-generated from your models.
                  Just register and get a full-featured admin panel.
                </p>
              </div>
              <div className={styles.featureCard}>
                <div className={styles.featureIcon}>🌐</div>
                <h3 className={styles.featureTitle}>REST API</h3>
                <p className={styles.featureDescription}>
                  Built-in REST API system like Django REST Framework.
                  Build APIs for React, Vue, or any frontend.
                </p>
              </div>
              <div className={styles.featureCard}>
                <div className={styles.featureIcon}>🛡️</div>
                <h3 className={styles.featureTitle}>Security First</h3>
                <p className={styles.featureDescription}>
                  Built-in CSRF, XSS, and SQL injection protection.
                  Security features enabled by default.
                </p>
              </div>
              <div className={styles.featureCard}>
                <div className={styles.featureIcon}>🔌</div>
                <h3 className={styles.featureTitle}>Extensible</h3>
                <p className={styles.featureDescription}>
                  Everything is extendable via plugins. Customize the
                  framework to fit your needs without modifying core code.
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className={styles.quickStart}>
          <div className={styles.container}>
            <div className={styles.quickStartContent}>
              <h2 className={styles.sectionTitle}>Get Started in 5 Minutes</h2>
              <p className={styles.sectionSubtitle}>
                Create your first forge application and see it in action
              </p>
              <div className={styles.quickStartSteps}>
                <div className={styles.step}>
                  <div className={styles.stepNumber}>1</div>
                  <div className={styles.stepContent}>
                    <h3>Create Project</h3>
                    <code>forge new myapp</code>
                  </div>
                </div>
                <div className={styles.step}>
                  <div className={styles.stepNumber}>2</div>
                  <div className={styles.stepContent}>
                    <h3>Define Models</h3>
                    <code>Define your schema</code>
                  </div>
                </div>
                <div className={styles.step}>
                  <div className={styles.stepNumber}>3</div>
                  <div className={styles.stepContent}>
                    <h3>Generate Code</h3>
                    <code>forge generate</code>
                  </div>
                </div>
                <div className={styles.step}>
                  <div className={styles.stepNumber}>4</div>
                  <div className={styles.stepContent}>
                    <h3>Run Server</h3>
                    <code>forge runserver</code>
                  </div>
                </div>
              </div>
              <div className={styles.quickStartButton}>
                <Link
                  className={styles.buttonPrimary}
                  to="/docs/getting-started/installation">
                  Start Building →
                </Link>
              </div>
            </div>
          </div>
        </div>

        <div className={styles.stats}>
          <div className={styles.container}>
            <div className={styles.statsGrid}>
              <div className={styles.stat}>
                <div className={styles.statNumber}>100%</div>
                <div className={styles.statLabel}>Type Safe</div>
              </div>
              <div className={styles.stat}>
                <div className={styles.statNumber}>0</div>
                <div className={styles.statLabel}>Runtime Errors</div>
              </div>
              <div className={styles.stat}>
                <div className={styles.statNumber}>∞</div>
                <div className={styles.statLabel}>Extensible</div>
              </div>
            </div>
          </div>
        </div>
      </Layout>
    </>
  );
}
