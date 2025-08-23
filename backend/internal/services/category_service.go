package services

import (
	"database/sql"
	"math"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
)

type CategoryService struct {
	DB *sql.DB
}

// Construtor
func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{DB: db}
}

// Cria uma categoria
func (s *CategoryService) CreateCategory(userID uint, name string) (*models.Category, error) {
	category := &models.Category{
		UserID: userID,
		Name:   name,
	}
	if err := repository.CreateCategory(s.DB, category); err != nil {
		return nil, err
	}
	return category, nil
}

// Lista categorias com paginação e filtro
func (s *CategoryService) ListCategories(
	userID uint,
	search string,
	page, limit int,
) ([]models.Category, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	categories, total, err := repository.FindCategoriesByUser(s.DB, userID, search, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

// Atualiza uma categoria
func (s *CategoryService) UpdateCategory(userID, id uint, name string) (*models.Category, error) {
	category := &models.Category{
		ID:     id,
		UserID: userID,
		Name:   name,
	}

	if err := repository.UpdateCategory(s.DB, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Deleta uma categoria
func (s *CategoryService) DeleteCategory(id, userID uint) error {
	if err := repository.DeleteCategory(s.DB, id, userID); err != nil {
		return err
	}
	return nil
}

// Calcula total de páginas
func (s *CategoryService) TotalPages(total, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
