package models

// LoginRequest > รับคำร้องล็อกอิน
type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}
