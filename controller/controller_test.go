package controller

import (
	"context"
	"errors"
	"monorepo/src/libs/constants"
	"monorepo/src/libs/log"
	"monorepo/src/review_service/entity"
	"monorepo/src/review_service/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCreateReview(t *testing.T) {

	assert := assert.New(t)
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	zapLogger := logger.With(zap.String("service", "review_service"))
	l := log.NewFactory(zapLogger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockReviewRepo(ctrl)

	id := uuid.NewString()

	realUuidGenerator := uuidGenerator
	uuidGenerator = func() string {
		return id
	}
	defer func() {
		uuidGenerator = realUuidGenerator
	}()

	t.Run("success", func(t *testing.T) {

		reviewReq := entity.CreateReviewReq{
			CreatorID:          id,
			CreatorPhoneNumber: "+998901234567",
			CreatorName:        "Aziz",
			Message:            "Juda kech yetib keldi",
		}
		ctx := context.Background()
		review := entity.Review{
			ID:                 id,
			CreatorID:          reviewReq.CreatorID,
			CreatorPhoneNumber: reviewReq.CreatorPhoneNumber,
			CreatorName:        reviewReq.CreatorName,
			Message:            reviewReq.Message,
		}
		repo.EXPECT().Create(ctx, review).Return(nil)

		c := New(repo, l, nil)

		err := c.CreateReview(ctx, reviewReq)

		assert.NoError(err)
	})

	t.Run("error internal", func(t *testing.T) {

		reviewReq := entity.CreateReviewReq{
			CreatorID:          id,
			CreatorPhoneNumber: "+998901234567",
			CreatorName:        "Aziz",
			Message:            "Juda kech yetib keldi",
		}
		ctx := context.Background()

		review := entity.Review{
			ID:                 id,
			CreatorID:          reviewReq.CreatorID,
			CreatorPhoneNumber: reviewReq.CreatorPhoneNumber,
			CreatorName:        reviewReq.CreatorName,
			Message:            reviewReq.Message,
		}
		repo.EXPECT().Create(ctx, review).Return(errors.New("internal error"))

		c := New(repo, l, nil)

		err := c.CreateReview(ctx, reviewReq)

		assert.Error(err)
	})

}

func TestAssignReview(t *testing.T) {

	assert := assert.New(t)
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	zapLogger := logger.With(zap.String("service", "review_service"))
	l := log.NewFactory(zapLogger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockReviewRepo(ctrl)

	id := uuid.NewString()

	realUuidGenerator := uuidGenerator
	uuidGenerator = func() string {
		return id
	}
	defer func() {
		uuidGenerator = realUuidGenerator
	}()

	t.Run("success", func(t *testing.T) {

		reviewReq := entity.AssignReviewReq{
			ID:           id,
			AssignedToID: id,
		}
		ctx := context.Background()
		review := entity.Review{
			ID:           id,
			AssignedToID: id,
		}
		repo.EXPECT().Update(ctx, review).Return(nil)

		c := New(repo, l, nil)

		err := c.AssignReview(ctx, reviewReq)

		assert.NoError(err)
	})

	t.Run("error review not found", func(t *testing.T) {

		reviewReq := entity.AssignReviewReq{
			ID:           id,
			AssignedToID: id,
		}
		ctx := context.Background()
		review := entity.Review{
			ID:           id,
			AssignedToID: id,
		}
		repo.EXPECT().Update(ctx, review).Return(constants.ErrRowsAffectedIsZero)

		c := New(repo, l, nil)

		err := c.AssignReview(ctx, reviewReq)

		assert.Error(err)
		assert.ErrorIs(err, constants.ErrRowsAffectedIsZero)
	})

	t.Run("error internal", func(t *testing.T) {

		reviewReq := entity.AssignReviewReq{
			ID:           id,
			AssignedToID: id,
		}
		ctx := context.Background()
		review := entity.Review{
			ID:           id,
			AssignedToID: id,
		}
		repo.EXPECT().Update(ctx, review).Return(errors.New("internal error"))

		c := New(repo, l, nil)

		err := c.AssignReview(ctx, reviewReq)

		assert.Error(err)
	})

}
