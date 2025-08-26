package services

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/daviolvr/Fintrack/internal/cache"
	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
)

type CategoryService struct {
	DB    *sql.DB
	cache *cache.Cache
}

// Construtor
func NewCategoryService(db *sql.DB, cache *cache.Cache) *CategoryService {
	return &CategoryService{DB: db, cache: cache}
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

	// Invalida cache do usuário
	s.cache.InvalidateUserCategories(userID)

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

	// Monta a chave do cache
	cacheKey := fmt.Sprintf("categories:%d:%s:%d:%d", userID, search, page, limit)

	// Verifica se existe no cache
	var cached cache.CategoryCacheData
	found, err := s.cache.Get(cacheKey, &cached)
	if err == nil && found {
		fmt.Println("Pegando do cache:", cacheKey)
		return cached.Categories, cached.Total, nil
	}

	// Busca no banco
	categories, total, err := repository.FindCategoriesByUser(s.DB, userID, search, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Salva no cache
	s.cache.Set(cacheKey, cache.CategoryCacheData{
		Categories: categories,
		Total:      total,
	}, time.Minute*5)

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

	// Invalida cache do usuário
	s.cache.InvalidateUserCategories(userID)

	return category, nil
}

// Deleta uma categoria
func (s *CategoryService) DeleteCategory(id, userID uint) error {
	if err := repository.DeleteCategory(s.DB, id, userID); err != nil {
		return err
	}

	// Invalida cache do usuário
	s.cache.InvalidateUserCategories(userID)

	return nil
}

// Calcula total de páginas
func (s *CategoryService) TotalPages(total, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
