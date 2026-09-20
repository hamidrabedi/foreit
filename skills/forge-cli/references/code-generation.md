---
sidebar_position: 1
---

# Code Generation

forge uses AST-based code generation to create type-safe code from your schema definitions.

## How It Works

1. **Parse Models** - AST parser reads your Go model files
2. **Extract Schema** - Extracts field definitions, relations, meta, hooks
3. **Generate Code** - Creates type-safe managers, querysets, and field expressions
4. **Write Files** - Writes all generated code for a package to one `gen.go` in the output directory (plus `api_gen.go` with `--api`)

## Generated Files

### Generated struct

Generation does not write your model struct — you author it. `gen.go` holds a
`<Model>Generated` struct carrying the schema fields and their tags:

```go
type PostGenerated struct {
    schema.BaseSchema
    ID        int64     `json:"id" db:"id" validate:""`
    Title     string    `json:"title" db:"title" validate:"required,max=200"`
    Published bool      `json:"published" db:"published" validate:""`
    CreatedAt time.Time `json:"created_at" db:"created_at" validate:""`
}
```

It also emits a `Validate()` method on your `Post` type.

### Manager

`<Model>Objects` is a package-level `orm.Manager[Model]`, bound to the table from
`Meta().TableName` (or the snake_case plural of the model name):

```go
var PostObjects = orm.MustNewManager[Post]("posts")
```

Use it directly: `PostObjects.Filter(...)`, `.Get(ctx, id)`, `.Create(ctx, &post)`.

### Field expressions

```go
type PostFields struct {
    ID        orm.Field[int64]
    Title     orm.Field[string]
    Published orm.Field[bool]
}

var PostFieldsInstance = PostFields{
    ID:        orm.NewField[int64]("id", "posts"),
    Title:     orm.NewField[string]("title", "posts"),
    Published: orm.NewField[bool]("published", "posts"),
}
```

### Relations

For each model with relations, generation adds `<Model>RelationExpr` (typed accessors
returning `*orm.RelationField[...]`) and `<Model>RelationHelper` with `Load<Relation>`
methods for ForeignKey, OneToOne, and ManyToMany.

## Running Code Generation

```bash
forge generate
```

This will:
1. Scan `models/` directory for model definitions
2. Parse each model file
3. Generate code for each model
4. Write generated files

## Customizing Generation

### Custom Templates

You can customize generation templates, though this is advanced and not recommended for most users.

### Generation Options

```bash
# Read schemas from another directory (default ./models)
forge generate --models ./app/blog

# Output to different directory (default ./models)
forge generate --output ./generated

# Also generate REST API ViewSets, serializers, and routes (api_gen.go)
forge generate --api

# Fail when a model expression cannot be evaluated
forge generate --strict
```

`--api` requires each model to have one auto-increment `int64` primary key named `id`,
and a concrete `int64` `ID`/`Id` field or an `orm.ModelWithID` implementation.
`gen.go` and `api_gen.go` are replaced together; if a write fails, both keep their
previous contents.

## Best Practices

1. **Don't Edit Generated Files** - They will be overwritten
2. **Regenerate After Model Changes** - Always regenerate after modifying models
3. **Commit Generated Files** - Include `gen.go` and `api_gen.go` in version control
4. **Use Type-Safe APIs** - Use generated field expressions and querysets

## Troubleshooting

### Generation Fails

- Check model syntax
- Ensure all imports are correct
- Verify schema interface implementation

### Generated Code Errors

- Regenerate code
- Check for circular dependencies
- Verify model definitions

## See Also

- [Models Guide](/docs/guides/models) - Model definitions
- [Plugins](/docs/advanced/plugins) - Extend code generation

