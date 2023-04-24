package entity

import (
	"time"

	"github.com/lib/pq"
)

type CreateReviewReq struct {
	CreatorID          string
	CreatorPhoneNumber string
	CreatorName        string
	Message            string
}

type ReviewOut struct {
	ID string

	CreatorID          string
	CreatorPhoneNumber string
	CreatorName        string
	Message            string

	AssignedToID string
	IsResolved   bool
	ResolvedAt   time.Time
}

type Review struct {
	ID string `gorm:"id;->;<-:create"`

	CreatorID          string `gorm:"creator_id;->;<-:create"`
	CreatorPhoneNumber string `gorm:"creator_phone_number;->;<-:create"`
	CreatorName        string `gorm:"creator_name;->;<-:create"`
	Message            string `gorm:"message;->;<-:create"`

	AssignedToID string `gorm:"assigned_to_id;<-:update"`

	TakenAction    string      `gorm:"taken_action;<-:update"`
	CustomerRating uint8       `gorm:"customer_rating;<-:update"`
	IsResolved     bool        `gorm:"is_resolved;<-:update"`
	WithSuccess    bool        `gorm:"with_success;<-:update"`
	ResolvedAt     pq.NullTime `gorm:"resolved_at;<-:update"`

	UpdatedAt time.Time `gorm:"updated_at;<-:update"`
}

type AssignReviewReq struct {
	ID string

	AssignedToID string
}

type ResolveReviewReq struct {
	ID string

	AssignedToID string

	TakenAction    string
	CustomerRating uint8
	IsResolved     bool
	WithSuccess    bool
	ResolvedAt     pq.NullTime
}
