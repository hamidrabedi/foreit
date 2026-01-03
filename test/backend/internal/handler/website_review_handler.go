package handler

import (
	"github.com/etemademan/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WebsiteHandler struct {
	WebsiteUsecase domain.WebsiteUsecase
}

func NewWebsiteHandler(app *fiber.App, websiteUsecase domain.WebsiteUsecase) {
	handler := &WebsiteHandler{
		WebsiteUsecase: websiteUsecase,
	}

	api := app.Group("/api/v1/websites")
	api.Get("/", handler.Search)
	api.Get("/:domain", handler.GetDetails)
}

func (h *WebsiteHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	websites, err := h.WebsiteUsecase.Search(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(websites)
}

func (h *WebsiteHandler) GetDetails(c *fiber.Ctx) error {
	domainName := c.Params("domain")
	website, reviews, err := h.WebsiteUsecase.GetDetails(c.Context(), domainName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"website": website,
		"reviews": reviews,
	})
}

type ReviewHandler struct {
	ReviewUsecase domain.ReviewUsecase
}

func NewReviewHandler(app *fiber.App, reviewUsecase domain.ReviewUsecase) {
	handler := &ReviewHandler{
		ReviewUsecase: reviewUsecase,
	}

	api := app.Group("/api/v1/reviews")
	// In a real app, we would add JWT middleware here
	api.Post("/", handler.SubmitReview)
}

type submitReviewRequest struct {
	UserID    string `json:"user_id"`
	WebsiteID string `json:"website_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

func (h *ReviewHandler) SubmitReview(c *fiber.Ctx) error {
	var req submitReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	uID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}

	wID, err := uuid.Parse(req.WebsiteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid website_id"})
	}

	err = h.ReviewUsecase.SubmitReview(c.Context(), uID, wID, req.Rating, req.Comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "review submitted successfully"})
}
