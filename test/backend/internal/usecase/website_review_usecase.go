package usecase

import (
	"context"
	"time"

	"github.com/etemademan/backend/internal/domain"
	"github.com/google/uuid"
)

type websiteUsecase struct {
	websiteRepo    domain.WebsiteRepository
	reviewRepo     domain.ReviewRepository
	contextTimeout time.Duration
}

func NewWebsiteUsecase(wRepo domain.WebsiteRepository, rRepo domain.ReviewRepository, timeout time.Duration) domain.WebsiteUsecase {
	return &websiteUsecase{
		websiteRepo:    wRepo,
		reviewRepo:     rRepo,
		contextTimeout: timeout,
	}
}

func (u *websiteUsecase) Search(ctx context.Context, query string) ([]domain.Website, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.websiteRepo.Search(ctx, query)
}

func (u *websiteUsecase) GetDetails(ctx context.Context, domainName string) (*domain.Website, []domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	website, err := u.websiteRepo.GetByDomain(ctx, domainName)
	if err != nil {
		return nil, nil, err
	}

	reviews, err := u.reviewRepo.GetByWebsiteID(ctx, website.ID)
	if err != nil {
		return website, nil, err
	}

	return website, reviews, nil
}

type reviewUsecase struct {
	reviewRepo     domain.ReviewRepository
	contextTimeout time.Duration
}

func NewReviewUsecase(rRepo domain.ReviewRepository, timeout time.Duration) domain.ReviewUsecase {
	return &reviewUsecase{
		reviewRepo:     rRepo,
		contextTimeout: timeout,
	}
}

func (u *reviewUsecase) SubmitReview(ctx context.Context, userID, websiteID uuid.UUID, rating int, comment string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	review := &domain.Review{
		ID:         uuid.New(),
		UserID:     userID,
		WebsiteID:  websiteID,
		Rating:     rating,
		Comment:    comment,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	return u.reviewRepo.Create(ctx, review)
}
