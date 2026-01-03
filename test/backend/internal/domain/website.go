package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Website struct {
	ID            uuid.UUID `json:"id"`
	Domain        string    `json:"domain"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	AverageRating float64   `json:"average_rating"`
	ReviewCount   int       `json:"review_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type Review struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	WebsiteID  uuid.UUID `json:"website_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type WebsiteRepository interface {
	Search(ctx context.Context, query string) ([]Website, error)
	GetByDomain(ctx context.Context, domain string) (*Website, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Website, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	GetByWebsiteID(ctx context.Context, websiteID uuid.UUID) ([]Review, error)
}

type WebsiteUsecase interface {
	Search(ctx context.Context, query string) ([]Website, error)
	GetDetails(ctx context.Context, domain string) (*Website, []Review, error)
}

type ReviewUsecase interface {
	SubmitReview(ctx context.Context, userID, websiteID uuid.UUID, rating int, comment string) error
}
