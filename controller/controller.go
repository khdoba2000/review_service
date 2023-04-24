package controller

import (
	"context"
	"fmt"
	"monorepo/src/libs/log"
	"monorepo/src/review_service/entity"
	"monorepo/src/review_service/pkg/utils"
	"monorepo/src/review_service/storage/repo"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

var uuidGenerator = func() string {
	return uuid.NewString()
}

// Controller defines a metadata service controller.
type ControllerImpl struct {
	repo   repo.ReviewRepo
	logger log.Factory
	tracer opentracing.Tracer
}

type Controller interface {
	CreateReview(ctx context.Context, req entity.CreateReviewReq) error
	GetAll(ctx context.Context) ([]entity.ReviewOut, error)
	AssignReview(ctx context.Context, req entity.AssignReviewReq) error
	ResolveReview(ctx context.Context, req entity.ResolveReviewReq) error
}

// New creates a metadata service controller.
func New(repo repo.ReviewRepo, logger log.Factory, tracer opentracing.Tracer) Controller {
	return &ControllerImpl{
		repo:   repo,
		logger: logger,
		tracer: tracer,
	}
}

func (c *ControllerImpl) CreateReview(ctx context.Context, req entity.CreateReviewReq) error {

	id := uuidGenerator()
	review := entity.Review{
		ID:                 id,
		CreatorID:          req.CreatorID,
		CreatorPhoneNumber: req.CreatorPhoneNumber,
		CreatorName:        req.CreatorName,
		Message:            req.Message,
	}

	err := c.repo.Create(ctx, review)
	if err != nil {
		c.logger.For(ctx).Error("Controller: error creating review", zap.Error(err))
		return err
	}
	return nil
}

func (c *ControllerImpl) GetAll(ctx context.Context) ([]entity.ReviewOut, error) {

	res, err := c.repo.GetAll(ctx)
	if err != nil {
		c.logger.For(ctx).Error("Controller: error GetAll", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (c *ControllerImpl) AssignReview(ctx context.Context, req entity.AssignReviewReq) error {

	review := entity.Review{
		ID:           req.ID,
		AssignedToID: req.AssignedToID,
	}
	err := c.repo.Update(ctx, review)
	if err != nil {
		c.logger.For(ctx).Error("Controller: error in Update", zap.Error(err))
		return err
	}
	return nil
}

func (c *ControllerImpl) ResolveReview(ctx context.Context, req entity.ResolveReviewReq) error {

	if !utils.IsValidRating(req.CustomerRating) {
		return fmt.Errorf("not a valid rating: %d", req.CustomerRating)
	}

	t := pq.NullTime{}
	t.Scan(time.Now())
	review := entity.Review{
		IsResolved:     true,
		ID:             req.ID,
		AssignedToID:   req.AssignedToID,
		TakenAction:    req.TakenAction,
		WithSuccess:    req.WithSuccess,
		CustomerRating: req.CustomerRating,
		ResolvedAt:     t,
	}
	err := c.repo.Update(ctx, review)
	if err != nil {
		c.logger.For(ctx).Error("Controller: error in Update", zap.Error(err))
		return err
	}
	return nil
}
