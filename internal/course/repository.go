package course

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DanyJDuque/gocourse_domain/domain"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, course *domain.Course) error
		Get(ctx context.Context, id string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, course *domain.Course) error {
	if err := r.db.Create(course).Error; err != nil {
		r.log.Panicf("error: %v", err)
		return err
	}

	r.log.Println("course created with id: ", course.ID)
	return nil
}

func (repo *repo) Get(ctx context.Context, id string) (*domain.Course, error) {
	course := domain.Course{ID: id}

	if err := repo.db.WithContext(ctx).First(&course).Error; err != nil {
		repo.log.Println(err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}

		return nil, err
	}
	return &course, nil

}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var c []domain.Course

	tx := repo.db.WithContext(ctx).Model(&c)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&c)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	return c, nil
}

func (r *repo) Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error {
	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}
	if endDate != nil {
		values["end_date"] = *endDate
	}

	result := r.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.log.Printf("course %s doesn't exist", id)
		return ErrNotFound{id}
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	course := domain.Course{ID: id}
	result := r.db.WithContext(ctx).Delete(&course)

	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.log.Printf("course %s doesn't exist", id)
		return ErrNotFound{id}
	}
	return nil
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, nil
	}
	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)

	}

	return tx
}
