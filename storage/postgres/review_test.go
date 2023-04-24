package postgres

import (
	"context"
	"database/sql"
	"errors"
	"monorepo/src/libs/constants"
	"monorepo/src/review_service/entity"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/test-go/testify/require"
	"github.com/test-go/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *reviewRepo
	// review     entity.Review
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	require.NoError(s.T(), err)

	s.repository = New(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) Test_Create() {

	ctx := context.Background()
	id := uuid.NewString()
	phone := "9090909"
	name := "aziz"
	msg := "some"
	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "reviews" ("id","creator_id","creator_phone_number","creator_name","message") VALUES ($1,$2,$3,$4,$5)`)).
		WithArgs(id, id, phone, name, msg).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	review := entity.Review{
		ID:                 id,
		CreatorID:          id,
		CreatorPhoneNumber: phone,
		CreatorName:        name,
		Message:            msg,
	}
	err := s.repository.Create(ctx, review)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_CreateError() {

	ctx := context.Background()
	id := uuid.NewString()
	phone := "9090909"
	name := "aziz"
	msg := "some"
	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "reviews" ("id","creator_id","creator_phone_number","creator_name","message") VALUES ($1,$2,$3,$4,$5)`)).
		WithArgs(id, id, phone, name, msg).
		WillReturnError(errors.New("internal error"))
	s.mock.ExpectRollback()

	review := entity.Review{
		ID:                 id,
		CreatorID:          id,
		CreatorPhoneNumber: phone,
		CreatorName:        name,
		Message:            msg,
	}
	err := s.repository.Create(ctx, review)

	require.Error(s.T(), err)
}

func (s *Suite) Test_Update_Success() {

	ctx := context.Background()

	id := uuid.NewString()

	updatedAt := time.Now()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		`^UPDATE "reviews" SET .* WHERE .*`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	review := entity.Review{
		ID:           id,
		AssignedToID: id,
		UpdatedAt:    updatedAt,
	}
	err := s.repository.Update(ctx, review)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_Update_ErrorNoRowsAffected() {

	ctx := context.Background()

	id := uuid.NewString()

	updatedAt := time.Now()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		`^UPDATE "reviews" SET .* WHERE .*`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 0))
	s.mock.ExpectCommit()

	review := entity.Review{
		ID:           id,
		AssignedToID: id,
		UpdatedAt:    updatedAt,
	}
	err := s.repository.Update(ctx, review)

	require.Error(s.T(), err)
	require.True(s.T(), errors.Is(err, constants.ErrRowsAffectedIsZero))
}

func (s *Suite) Test_Update_ErrorInternal() {

	ctx := context.Background()

	id := uuid.NewString()

	updatedAt := time.Now()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		`^UPDATE "reviews" SET .* WHERE .*`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("internal error"))
	s.mock.ExpectRollback()

	review := entity.Review{
		ID:           id,
		AssignedToID: id,
		UpdatedAt:    updatedAt,
	}
	err := s.repository.Update(ctx, review)

	require.Error(s.T(), err)
}

func TestCreate(t *testing.T) {
	suite.Run(t, new(Suite))
}
