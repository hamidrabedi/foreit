package models

import (
	"context"
	"database/sql"
	"errors"
)

// Session represents a database session (SQLAlchemy-inspired)
type Session interface {
	// Begin starts a transaction
	Begin(ctx context.Context) (Transaction, error)
	
	// Commit commits the current transaction
	Commit() error
	
	// Rollback rolls back the current transaction
	Rollback() error
	
	// Query executes a query
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	
	// Exec executes a statement
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	
	// Close closes the session
	Close() error
}

// Transaction represents a database transaction
type Transaction interface {
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
	
	// Query executes a query within the transaction
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	
	// Exec executes a statement within the transaction
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// SessionManager manages database sessions
type SessionManager interface {
	// GetSession returns a session for the given context
	GetSession(ctx context.Context) (Session, error)
	
	// WithTransaction executes a function within a transaction
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

// ContextSession extracts session from context
func ContextSession(ctx context.Context) (Session, error) {
	session, ok := ctx.Value(sessionKey).(Session)
	if !ok {
		return nil, ErrNoSession
	}
	return session, nil
}

// WithSession adds a session to context
func WithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// WithTransaction executes a function within a transaction (type-safe context manager)
func WithTransaction(ctx context.Context, session Session, fn func(context.Context) error) error {
	tx, err := session.Begin(ctx)
	if err != nil {
		return err
	}
	
	txCtx := WithSession(ctx, &transactionSession{tx: tx})
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}
	
	return tx.Commit()
}

// Transactional executes operations within a transaction (SQLAlchemy-style)
func Transactional[T any](
	ctx context.Context,
	session Session,
	fn func(context.Context) (T, error),
) (T, error) {
	var zero T
	
	tx, err := session.Begin(ctx)
	if err != nil {
		return zero, err
	}
	
	txCtx := WithSession(ctx, &transactionSession{tx: tx})
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	result, err := fn(txCtx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return zero, errors.Join(err, rollbackErr)
		}
		return zero, err
	}
	
	if err := tx.Commit(); err != nil {
		return zero, err
	}
	
	return result, nil
}

// transactionSession wraps a transaction as a session
type transactionSession struct {
	tx Transaction
}

func (s *transactionSession) Begin(ctx context.Context) (Transaction, error) {
	return nil, errors.New("cannot begin transaction within transaction")
}

func (s *transactionSession) Commit() error {
	return s.tx.Commit()
}

func (s *transactionSession) Rollback() error {
	return s.tx.Rollback()
}

func (s *transactionSession) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.tx.Query(ctx, query, args...)
}

func (s *transactionSession) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.tx.Exec(ctx, query, args...)
}

func (s *transactionSession) Close() error {
	return nil
}

type sessionKeyType struct{}

var sessionKey = sessionKeyType{}

var (
	ErrNoSession = errors.New("no session in context")
)

