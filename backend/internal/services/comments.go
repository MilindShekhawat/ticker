package services

import (
	"strings"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

type CommentService interface {
	ListForUser(ticketID, userID int) ([]models.Comment, error)
	CreateForUser(ticketID, userID int, body string) (*models.Comment, error)
	UpdateForUser(commentID, userID int, body string) (*models.Comment, error)
	DeleteForUser(commentID, userID int) error
}

type commentService struct {
	comments store.CommentStore
	members  store.ProjectMemberStore
}

func NewCommentService(comments store.CommentStore, members store.ProjectMemberStore) CommentService {
	return &commentService{comments: comments, members: members}
}

func (s *commentService) ListForUser(ticketID, userID int) ([]models.Comment, error) {
	projectID, err := s.comments.GetTicketProjectID(ticketID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return nil, err
	}

	return s.comments.ListByTicket(ticketID)
}

func (s *commentService) CreateForUser(ticketID, userID int, body string) (*models.Comment, error) {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return nil, ErrInvalidCommentBody
	}

	projectID, err := s.comments.GetTicketProjectID(ticketID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return nil, err
	}

	return s.comments.Create(ticketID, userID, trimmedBody)
}

func (s *commentService) UpdateForUser(commentID, userID int, body string) (*models.Comment, error) {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return nil, ErrInvalidCommentBody
	}

	comment, err := s.comments.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	projectID, err := s.comments.GetTicketProjectID(comment.TicketID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return nil, err
	}

	return s.comments.Update(commentID, trimmedBody)
}

func (s *commentService) DeleteForUser(commentID, userID int) error {
	comment, err := s.comments.GetByID(commentID)
	if err != nil {
		return err
	}

	projectID, err := s.comments.GetTicketProjectID(comment.TicketID)
	if err != nil {
		return err
	}
	if _, err := s.members.GetUserRole(projectID, userID); err != nil {
		return err
	}

	return s.comments.Delete(commentID)
}
