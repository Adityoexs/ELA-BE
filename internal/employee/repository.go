package employee

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("employee not found")

type Repository interface {
	Create(ctx context.Context, employee *Employee) error
	GetByID(ctx context.Context, id uint) (*Employee, error)
	List(ctx context.Context) ([]Employee, error)
	Update(ctx context.Context, employee *Employee) error
	Delete(ctx context.Context, id uint) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, employee *Employee) error {
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return fmt.Errorf("create employee: %w", err)
	}
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*Employee, error) {
	var employee Employee
	if err := r.db.WithContext(ctx).First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get employee by id: %w", err)
	}
	return &employee, nil
}

func (r *gormRepository) List(ctx context.Context) ([]Employee, error) {
	var employees []Employee
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("list employees: %w", err)
	}
	return employees, nil
}

func (r *gormRepository) Update(ctx context.Context, employee *Employee) error {
	if err := r.db.WithContext(ctx).Save(employee).Error; err != nil {
		return fmt.Errorf("update employee: %w", err)
	}
	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Employee{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
