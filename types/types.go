package types

// Department struct
type Department struct {
	DepartmentID   int    `gorm:"primaryKey;autoIncrement" json:"department_id,omitempty"` // PK
	DepartmentName string `json:"department_name"`                                         // ชื่อแผนก
}

// User struct
type User struct {
	UserID       int    `gorm:"primaryKey;autoIncrement" json:"user_id,omitempty"` // PK
	FirstName    string `json:"first_name"`                                        // ชื่อจริง
	LastName     string `json:"last_name"`                                         // นามสกุล
	Email        string `json:"email"`                                             // อีเมล
	Phone        string `json:"phone"`                                             // เบอร์โทร
	DepartmentID int    `json:"department_id"`                                     // FK → Department
}
