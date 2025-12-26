---
sidebar_position: 3
---

# Library Management Example

A library management system demonstrating forge with SQLite.

## Overview

This example shows:

- Library management models
- SQLite database support
- Book borrowing system
- Category organization

## Models

### Author

```go
type Author struct {
    schema.BaseSchema
}

func (Author) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("name").Required().MaxLength(100).Build(),
        schema.String("email").Unique().MaxLength(255).Build(),
        schema.Text("bio").Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
    }
}

func (Author) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToMany("books", "Book").
            RelatedName("author"),
    }
}
```

### Book

```go
type Book struct {
    schema.BaseSchema
}

func (Book) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("title").Required().MaxLength(200).Build(),
        schema.String("isbn").Unique().MaxLength(20).Build(),
        schema.Int32("published_year").Build(),
        schema.Text("description").Build(),
        schema.Bool("available").Default(true).Build(),
    }
}

func (Book) Relations() []schema.Relation {
    return []schema.Relation{
        relations.ForeignKey("author", "Author").
            Required().
            OnDelete(schema.Cascade),
        relations.ForeignKey("category", "Category").
            OnDelete(schema.SetNull),
        relations.OneToMany("loans", "Loan").
            RelatedName("book"),
    }
}
```

### Borrower

```go
type Borrower struct {
    schema.BaseSchema
}

func (Borrower) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("name").Required().MaxLength(100).Build(),
        schema.String("email").Unique().Required().MaxLength(255).Build(),
        schema.String("phone").MaxLength(20).Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
    }
}

func (Borrower) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToMany("loans", "Loan").
            RelatedName("borrower"),
    }
}
```

### Loan

```go
type Loan struct {
    schema.BaseSchema
}

func (Loan) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.Time("borrowed_at").AutoNowAdd().Build(),
        schema.Time("due_date").Required().Build(),
        schema.Time("returned_at").Build(),
        schema.String("status").
            Choices("active", "returned", "overdue").
            Default("active").
            Build(),
    }
}

func (Loan) Relations() []schema.Relation {
    return []schema.Relation{
        relations.ForeignKey("book", "Book").
            Required().
            OnDelete(schema.Cascade),
        relations.ForeignKey("borrower", "Borrower").
            Required().
            OnDelete(schema.Cascade),
    }
}
```

## Usage

### Borrow a Book

```go
loan := &models.Loan{
    Book:     book,
    Borrower: borrower,
    DueDate: time.Now().AddDate(0, 0, 14), // 2 weeks
    Status:   "active",
}

err := models.Loan.Objects.Create(ctx, loan)

// Mark book as unavailable
book.Available = false
err = models.Book.Objects.Update(ctx, book)
```

### Return a Book

```go
loan.ReturnedAt = time.Now()
loan.Status = "returned"
err := models.Loan.Objects.Update(ctx, loan)

// Mark book as available
book.Available = true
err = models.Book.Objects.Update(ctx, book)
```

### Find Overdue Loans

```go
overdueLoans, err := models.Loan.Objects.
    Filter(models.Loan.Fields.Status.Equals("active")).
    Filter(models.Loan.Fields.DueDate.Less(time.Now())).
    PrefetchRelated("book", "borrower").
    All(ctx)
```

### Get Available Books

```go
availableBooks, err := models.Book.Objects.
    Filter(models.Book.Fields.Available.Equals(true)).
    SelectRelated("author", "category").
    OrderBy("title").
    All(ctx)
```

## Configuration

For SQLite, configure in `config/config.yaml`:

```yaml
database:
  driver: sqlite3
  dsn: "./library.db"
```

Note: SQLite requires CGO. Build with:

```bash
CGO_ENABLED=1 go build
```

## See Also

- [Models Guide](/docs/guides/models) - Model definitions
- [Queries Guide](/docs/guides/queries) - Query examples
- [E-commerce Example](/docs/examples/ecommerce) - More complex example

