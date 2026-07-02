package service

import (
	"time"

	"reserveflow-v1/dao"
)

// ReservationService defines the interface for reservation-related business logic.
type ReservationService interface {
	GetExpiredReservationIDs(currentTime time.Time) ([]uint, error)
	ExpireReservation(resID uint) error
}

type reservationService struct{}

// NewReservationService returns a new instance of ReservationService.
func NewReservationService() ReservationService {
	return &reservationService{}
}

func (s *reservationService) GetExpiredReservationIDs(currentTime time.Time) ([]uint, error) {
	return dao.GetExpiredReservationIDs(currentTime)
}

func (s *reservationService) ExpireReservation(resID uint) error {
	return dao.MarkReservationAsExpired(resID)
}
