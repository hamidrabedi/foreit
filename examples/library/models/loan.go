package models

import (
	"github.com/forgego/forge/schema"
)

// Loan represents a book loan/borrowing record
type Loan struct {
	schema.BaseSchema
}

// Fields returns all field definitions for Loan
func (Loan) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("book_id").Required().VerboseName("Book ID").Build(),
		schema.Int64("borrower_id").Required().VerboseName("Borrower ID").Build(),
		schema.Time("borrowed_at").AutoNowAdd().VerboseName("Borrowed Date").Build(),
		schema.Time("due_at").Required().VerboseName("Due Date").Build(),
		schema.Time("returned_at").VerboseName("Returned Date").Build(),
		schema.Bool("is_returned").Default(false).VerboseName("Returned").Build(),
		schema.Float64("fine_amount").Default(0.0).VerboseName("Fine Amount").Build(),
		schema.String("notes").VerboseName("Notes").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (Loan) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "loans",
		VerboseName:       "Loan",
		VerboseNamePlural: "Loans",
		OrderBy:           []string{"-borrowed_at"},
		Indexes: []schema.Index{
			{Name: "idx_loan_book", Fields: []string{"book_id"}, Unique: false},
			{Name: "idx_loan_borrower", Fields: []string{"borrower_id"}, Unique: false},
			{Name: "idx_loan_returned", Fields: []string{"is_returned"}, Unique: false},
		},
	}
}

// Relations returns all relationship definitions
func (Loan) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Hooks returns model lifecycle hooks
func (Loan) Hooks() *schema.ModelHooks {
	return nil
}
