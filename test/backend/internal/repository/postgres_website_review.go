package repository

import (
	"context"
	"errors"

	"github.com/etemademan/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresWebsiteRepository struct {
	db *pgxpool.Pool
}

func NewPostgresWebsiteRepository(db *pgxpool.Pool) domain.WebsiteRepository {
	return &postgresWebsiteRepository{db: db}
}

func (r *postgresWebsiteRepository) Search(ctx context.Context, query string) ([]domain.Website, error) {
	sql := `SELECT id, domain, name, description, average_rating, review_count, created_at 
	        FROM websites WHERE name ILIKE $1 OR domain ILIKE $1`
	rows, err := r.db.Query(ctx, sql, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var websites []domain.Website
	for rows.Next() {
		var w domain.Website
		if err := rows.Scan(&w.ID, &w.Domain, &w.Name, &w.Description, &w.AverageRating, &w.ReviewCount, &w.CreatedAt); err != nil {
			return nil, err
		}
		websites = append(websites, w)
	}
	return websites, nil
}

func (r *postgresWebsiteRepository) GetByDomain(ctx context.Context, domainName string) (*domain.Website, error) {
	sql := `SELECT id, domain, name, description, average_rating, review_count, created_at FROM websites WHERE domain = $1`
	w := &domain.Website{}
	err := r.db.QueryRow(ctx, sql, domainName).Scan(&w.ID, &w.Domain, &w.Name, &w.Description, &w.AverageRating, &w.ReviewCount, &w.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("website not found")
		}
		return nil, err
	}
	return w, nil
}

func (r *postgresWebsiteRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Website, error) {
	sql := `SELECT id, domain, name, description, average_rating, review_count, created_at FROM websites WHERE id = $1`
	w := &domain.Website{}
	err := r.db.QueryRow(ctx, sql, id).Scan(&w.ID, &w.Domain, &w.Name, &w.Description, &w.AverageRating, &w.ReviewCount, &w.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("website not found")
		}
		return nil, err
	}
	return w, nil
}

type postgresReviewRepository struct {
	db *pgxpool.Pool
}

func NewPostgresReviewRepository(db *pgxpool.Pool) domain.ReviewRepository {
	return &postgresReviewRepository{db: db}
}

func (r *postgresReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert review
	sqlReview := `INSERT INTO reviews (id, user_id, website_id, rating, comment, is_verified, created_at) 
	              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(ctx, sqlReview, review.ID, review.UserID, review.WebsiteID, review.Rating, review.Comment, review.IsVerified, review.CreatedAt)
	if err != nil {
		return err
	}

	// Update website stats
	sqlWebsite := `UPDATE websites SET 
	               average_rating = (average_rating * review_count + $1) / (review_count + 1),
	               review_count = review_count + 1
	               WHERE id = $2`
	_, err = tx.Exec(ctx, sqlWebsite, review.Rating, review.WebsiteID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *postgresReviewRepository) GetByWebsiteID(ctx context.Context, websiteID uuid.UUID) ([]domain.Review, error) {
	sql := `SELECT id, user_id, website_id, rating, comment, is_verified, created_at FROM reviews WHERE website_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, sql, websiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rev domain.Review
		if err := rows.Scan(&rev.ID, &rev.UserID, &rev.WebsiteID, &rev.Rating, &rev.Comment, &rev.IsVerified, &rev.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}
