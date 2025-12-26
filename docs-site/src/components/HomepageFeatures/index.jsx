import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Type-Safe ORM',
    emoji: '🔒',
    description: (
      <>
        Full Django ORM features with compile-time type checking. Write queries
        with confidence knowing the compiler will catch errors before runtime.
      </>
    ),
  },
  {
    title: 'Code Generation',
    emoji: '⚡',
    description: (
      <>
        AST-based code generation for models, managers, and querysets. Generate
        type-safe code automatically from your schema definitions.
      </>
    ),
  },
  {
    title: 'Auto-Generated Admin',
    emoji: '🎛️',
    description: (
      <>
        Django-style admin interface auto-generated from your model registry.
        Just register your models and get a full-featured admin panel.
      </>
    ),
  },
  {
    title: 'REST API',
    emoji: '🌐',
    description: (
      <>
        Built-in REST API system similar to Django REST Framework. Build APIs
        for React, Vue, or any frontend that consumes JSON.
      </>
    ),
  },
  {
    title: 'Security First',
    emoji: '🛡️',
    description: (
      <>
        Built-in CSRF, XSS, and SQL injection protection. Security features are
        enabled by default to keep your application safe.
      </>
    ),
  },
  {
    title: 'Fully Extensible',
    emoji: '🔌',
    description: (
      <>
        Everything is extendable and overridable via plugins. Customize the
        framework to fit your needs without modifying core code.
      </>
    ),
  },
];

function Feature({emoji, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <div className={styles.featureEmoji}>{emoji}</div>
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

