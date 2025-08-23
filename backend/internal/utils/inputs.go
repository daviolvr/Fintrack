package utils

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required,jwt"`
}

type CategoryInput struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

type TransactionInput struct {
	CategoryID  uint    `json:"category_id" binding:"required,min=1"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"max=255"`
	Date        string  `json:"date" binding:"required,datetime=2006-01-02"`
}

type UserUpdateInput struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Email     string `json:"email" binding:"omitempty,email"`
}

type UserUpdateBalanceInput struct {
	Balance float64 `json:"balance" binding:"required"`
}

type UserDeleteInput struct {
	Password string `json:"password" binding:"required"`
}

type UserUpdatePasswordInput struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=72,nefield=Password"`
}
