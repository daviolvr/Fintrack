package dto

type UserUpdateParam struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserChangePasswordParam struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type TransactionCreateParam struct {
	CategoryID  int64   `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionUpdateParam struct {
	CategoryID  int64   `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type BalanceUpdateParam struct {
	Balance float64 `json:"balance"`
}
