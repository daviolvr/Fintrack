package repository

import (
	"github.com/daviolvr/Fintrack/internal/models"
	"gorm.io/gorm"
)

// Cria categoria para o usuário
func CreateCategory(db *gorm.DB, category *models.Category) error {
	if err := db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

// Busca categorias do usuário
func FindCategoriesByUser(
	db *gorm.DB,
	userID uint,
	search string,
	page, limit int,
) ([]models.Category, int, error) {
	var categories []models.Category
	var total int64

	query := db.Model(&models.Category{}).Where("user_id = ?", userID)

	// Filtro de busca
	if search != "" {
		query = query.Where("name ILIKIE ?", "%"+search+"%")
	}

	// Conta o total de registros antes da paginação
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginação e ordenação
	offset := (page - 1) * limit
	if err := query.Order("name").Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, int(total), nil
}

// Atualiza categoria pelo ID e user_id
func UpdateCategory(db *gorm.DB, category *models.Category) error {
	result := db.Model(&models.Category{}).
		Where("id = ? AND user_id = ?", category.ID, category.UserID).
		Update("name", category.Name)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Deleta categoria pelo ID e user_id
func DeleteCategory(db *gorm.DB, id, userID uint) error {
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Category{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
