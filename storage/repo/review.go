package repo

import (
	"context"
	"fmt"
	"monorepo/src/review_service/configs"
	"monorepo/src/review_service/entity"
	"monorepo/src/review_service/pkg/db"
	"monorepo/src/review_service/storage/postgres"
)

// ReviewRepo defines base interface for Review Storage
type ReviewRepo interface {
	Create(context.Context, entity.Review) error
	Update(ctx context.Context, req entity.Review) error
	GetAll(ctx context.Context) ([]entity.ReviewOut, error)
}

// New ...
func New(config *configs.Configuration) ReviewRepo {
	fmt.Println("storage ")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDatabase)

	dbIns, err := db.Init(dbURL)
	if err != nil {
		panic(err)
	}

	storage := postgres.New(dbIns)

	return storage
}
