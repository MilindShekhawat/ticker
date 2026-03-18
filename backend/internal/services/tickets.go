package services

import (
	"strings"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

type TicketService interface {
	ListForUser(projectID, userID int) ([]models.Ticket, error)
	GetForUser(ticketID, userID int) (*models.Ticket, error)
	CreateForUser(projectID int, title, description string, statusID, priorityID, userID int, assigneeID *int, tagIDs []int) (*models.Ticket, error)
	UpdateForUser(ticketID, userID int, title, description *string, statusID, priorityID, assigneeID *int, tagIDs *[]int) (*models.Ticket, error)
	DeleteForUser(ticketID, userID int) error
}

type ticketService struct {
	tickets    store.TicketStore
	projects   store.ProjectStore
	members    store.ProjectMemberStore
	statuses   store.StatusStore
	priorities store.PriorityStore
	users      store.UserStore
	tags       store.TagStore
}

func NewTicketService(
	tickets store.TicketStore,
	projects store.ProjectStore,
	members store.ProjectMemberStore,
	statuses store.StatusStore,
	priorities store.PriorityStore,
	users store.UserStore,
	tags store.TagStore,
) TicketService {
	return &ticketService{
		tickets:    tickets,
		projects:   projects,
		members:    members,
		statuses:   statuses,
		priorities: priorities,
		users:      users,
		tags:       tags,
	}
}

func (s *ticketService) ListForUser(projectID, userID int) ([]models.Ticket, error) {
	if _, err := s.projects.GetByID(projectID); err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return nil, err
	}

	return s.tickets.ListByProject(projectID)
}

func (s *ticketService) GetForUser(ticketID, userID int) (*models.Ticket, error) {
	ticket, err := s.tickets.GetByID(ticketID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(ticket.ProjectID, userID); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *ticketService) CreateForUser(projectID int, title, description string, statusID, priorityID, userID int, assigneeID *int, tagIDs []int) (*models.Ticket, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTicketTitleRequired
	}
	if len(title) > 200 {
		return nil, ErrTicketTitleTooLong
	}
	if statusID == 0 {
		return nil, ErrInvalidStatusID
	}
	if priorityID == 0 {
		return nil, ErrInvalidPriorityID
	}
	if assigneeID != nil && *assigneeID == 0 {
		return nil, ErrInvalidAssigneeID
	}

	if _, err := s.projects.GetByID(projectID); err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return nil, err
	}
	if err := s.validateStatusProject(statusID, projectID); err != nil {
		return nil, err
	}
	if err := s.validatePriorityProject(priorityID, projectID); err != nil {
		return nil, err
	}
	if err := s.validateAssignee(assigneeID); err != nil {
		return nil, err
	}
	if err := s.validateTagsProject(tagIDs, projectID); err != nil {
		return nil, err
	}

	ticket, err := s.tickets.Create(projectID, title, description, statusID, priorityID, userID, assigneeID)
	if err != nil {
		return nil, err
	}

	if err := s.tickets.ReplaceTags(ticket.ID, tagIDs); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *ticketService) UpdateForUser(ticketID, userID int, title, description *string, statusID, priorityID, assigneeID *int, tagIDs *[]int) (*models.Ticket, error) {
	ticket, err := s.tickets.GetByID(ticketID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(ticket.ProjectID, userID); err != nil {
		return nil, err
	}

	if title != nil {
		trimmed := strings.TrimSpace(*title)
		if trimmed == "" {
			return nil, ErrTicketTitleRequired
		}
		if len(trimmed) > 200 {
			return nil, ErrTicketTitleTooLong
		}
		title = &trimmed
	}

	if statusID != nil {
		if *statusID == 0 {
			return nil, ErrInvalidStatusID
		}
		if err := s.validateStatusProject(*statusID, ticket.ProjectID); err != nil {
			return nil, err
		}
	}
	if priorityID != nil {
		if *priorityID == 0 {
			return nil, ErrInvalidPriorityID
		}
		if err := s.validatePriorityProject(*priorityID, ticket.ProjectID); err != nil {
			return nil, err
		}
	}
	if assigneeID != nil && *assigneeID == 0 {
		return nil, ErrInvalidAssigneeID
	}
	if err := s.validateAssignee(assigneeID); err != nil {
		return nil, err
	}

	if tagIDs != nil {
		if err := s.validateTagsProject(*tagIDs, ticket.ProjectID); err != nil {
			return nil, err
		}
	}

	updatedTicket, err := s.tickets.Update(ticketID, title, description, statusID, priorityID, assigneeID)
	if err != nil {
		return nil, err
	}

	if tagIDs != nil {
		if err := s.tickets.ReplaceTags(ticketID, *tagIDs); err != nil {
			return nil, err
		}
	}

	return updatedTicket, nil
}

func (s *ticketService) DeleteForUser(ticketID, userID int) error {
	ticket, err := s.tickets.GetByID(ticketID)
	if err != nil {
		return err
	}
	if _, err := s.members.GetUserRole(ticket.ProjectID, userID); err != nil {
		return err
	}

	return s.tickets.Delete(ticketID)
}

func (s *ticketService) validateStatusProject(statusID, projectID int) error {
	status, err := s.statuses.GetByID(statusID)
	if err != nil {
		return err
	}
	if status.ProjectID != projectID {
		return store.ErrStatusNotFound
	}

	return nil
}

func (s *ticketService) validatePriorityProject(priorityID, projectID int) error {
	priority, err := s.priorities.GetByID(priorityID)
	if err != nil {
		return err
	}
	if priority.ProjectID != projectID {
		return store.ErrPriorityNotFound
	}

	return nil
}

func (s *ticketService) validateAssignee(assigneeID *int) error {
	if assigneeID == nil {
		return nil
	}
	_, err := s.users.GetByID(*assigneeID)
	return err
}

func (s *ticketService) validateTagsProject(tagIDs []int, projectID int) error {
	for _, tagID := range tagIDs {
		if err := s.tags.BelongsToProject(tagID, projectID); err != nil {
			return err
		}
	}

	return nil
}
