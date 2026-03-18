package services

import (
	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

type TicketActivityService interface {
	ListForUser(ticketID, userID int) ([]models.TicketActivity, error)
}

type ticketActivityService struct {
	activities store.TicketActivityStore
	members    store.ProjectMemberStore
}

func NewTicketActivityService(activities store.TicketActivityStore, members store.ProjectMemberStore) TicketActivityService {
	return &ticketActivityService{activities: activities, members: members}
}

func (s *ticketActivityService) ListForUser(ticketID, userID int) ([]models.TicketActivity, error) {
	projectID, err := s.activities.GetTicketProjectID(ticketID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return nil, err
	}

	return s.activities.ListByTicket(ticketID)
}
