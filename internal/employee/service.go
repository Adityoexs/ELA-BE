package employee

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Create(ctx context.Context, input CreateEmployeeInput) (*Employee, error)
	GetByID(ctx context.Context, id uint) (*Employee, error)
	List(ctx context.Context) ([]Employee, error)
	Update(ctx context.Context, id uint, input UpdateEmployeeInput) (*Employee, error)
	Delete(ctx context.Context, id uint) error
}

type CreateEmployeeInput struct {
	Name     string
	Email    string
	Position string
}

type UpdateEmployeeInput struct {
	Name     string
	Email    string
	Position string
}

type service struct {
	repo   Repository
	logger *logrus.Entry
}

func NewService(repo Repository, logger *logrus.Entry) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) Create(ctx context.Context, input CreateEmployeeInput) (*Employee, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return nil, errors.New("email is required")
	}

	employee := &Employee{
		Name:     input.Name,
		Email:    input.Email,
		Position: input.Position,
	}

	if err := s.repo.Create(ctx, employee); err != nil {
		return nil, err
	}

	s.logger.WithField("employee_id", employee.ID).Info("employee created")
	return employee, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (s *service) List(ctx context.Context) ([]Employee, error) {
	employees, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (s *service) Update(ctx context.Context, id uint, input UpdateEmployeeInput) (*Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Name) != "" {
		employee.Name = input.Name
	}
	if strings.TrimSpace(input.Email) != "" {
		employee.Email = input.Email
	}
	if strings.TrimSpace(input.Position) != "" {
		employee.Position = input.Position
	}

	if err := s.repo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("update employee: %w", err)
	}

	s.logger.WithField("employee_id", employee.ID).Info("employee updated")
	return employee, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.WithField("employee_id", id).Info("employee deleted")
	return nil
}
