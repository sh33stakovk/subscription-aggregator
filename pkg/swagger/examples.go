package swagger

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionExample struct {
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       uint      `json:"price"        example:"999"`
	UserID      uuid.UUID `json:"user_id"      example:"11111111-1111-1111-1111-111111111111"`
	StartDate   string    `json:"start_date"   example:"07-2025"`
}

type UpdateSubscriptionExample struct {
	ServiceName string `json:"service_name" example:"Yandex"`
	Price       uint   `json:"price"        example:"100"`
}

type SubscriptionResponse struct {
	ID          uint      `json:"id"           example:"1"`
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       uint      `json:"price"        example:"999"`
	UserID      uuid.UUID `json:"user_id"      example:"11111111-1111-1111-1111-111111111111"`
	StartDate   string    `json:"start_date"   example:"07-2025"`
	CreatedAt   time.Time `json:"created_at"   example:"2025-07-26T12:34:56Z"`
	UpdatedAt   time.Time `json:"updated_at"   example:"2025-07-26T12:34:56Z"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"failed to bind json"`
}

type MessageResponse struct {
	Message string `json:"message"      example:"created"`
	ID      uint   `json:"id,omitempty" example:"1"`
}

type SumResponse struct {
	SumPrice int `json:"sum_price" example:"2997"`
}
