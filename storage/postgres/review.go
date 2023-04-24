package postgres

import (
	"context"
	"fmt"
	"monorepo/src/libs/constants"
	"monorepo/src/review_service/entity"
	"time"

	"gorm.io/gorm"
)

type reviewRepo struct {
	db *gorm.DB
}

// New ...
func New(db *gorm.DB) *reviewRepo {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) Create(ctx context.Context, review entity.Review) error {

	res := r.db.Create(review)
	if res.Error != nil {
		return fmt.Errorf("error in Create: %w", res.Error)
	}

	return nil
}

func (r *reviewRepo) GetAll(ctx context.Context) ([]entity.ReviewOut, error) {

	//TODO: add pagination
	var reviews []entity.ReviewOut
	res := r.db.
		Model(entity.Review{}).
		Scan(&reviews)
	if res.Error != nil {
		return nil, fmt.Errorf("error in GetAll: %w", res.Error)
	}

	return reviews, nil
}

func (r *reviewRepo) Update(ctx context.Context, review entity.Review) error {

	review.UpdatedAt = time.Now()

	res := r.db.Model(entity.Review{}).
		Where("id=?", review.ID).
		Updates(review)
	if res.Error != nil {
		return fmt.Errorf("error in Update: %w", res.Error)
	} else if res.RowsAffected == 0 {
		return fmt.Errorf("error in Update: %w", constants.ErrRowsAffectedIsZero)
	}

	return nil
}
