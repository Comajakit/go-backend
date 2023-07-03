package interfaces

import "github.com/google/uuid"

type CreateTripRequest struct {
	TripCategoryID uint            `json:"tripCategoryId"`
	TripName       string          `json:"tripName"`
	Members        []MemberRequest `json:"members"`
}

type MemberRequest struct {
	ID uuid.UUID `json:"id"`
}
