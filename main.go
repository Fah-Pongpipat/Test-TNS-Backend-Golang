package main

import (
	// import Types
	"backend/types"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB
var GormDB *gorm.DB

func ConnectDB() {
	// Connection Database
	host := "localhost"  // PostgreSQL server host
	port := 5432         // PostgreSQL port
	user := "postgres"   // user ของ PostgreSQL
	password := "123456" // password ของ user
	dbname := "Test-TNS" // database
	sslmode := "disable"

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	// ตรวจสอบว่า connection ใช้งานได้จริง
	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")

	GormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: DB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ GORM initialized with existing connection")
}
func CreatedTable() {
	createDepartment := `
	CREATE TABLE IF NOT EXISTS department (
    department_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    department_name VARCHAR(255) NOT NULL
	);
	`
	_, errDepartment := DB.Exec(createDepartment)
	if errDepartment != nil {
		log.Fatal("❌ Failed to create Department table:", errDepartment)
	}

	createUser := `
		CREATE TABLE IF NOT EXISTS users (
		user_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		email VARCHAR(100),
		phone VARCHAR(20),
		department_id INT REFERENCES department(department_id)
		);
		`
	_, errUser := DB.Exec(createUser)
	if errUser != nil {
		log.Fatal("❌ Failed to create User table:", errUser)
	}

	fmt.Println("✅ Tables created successfully!")
}
func GetAllUsers(c *gin.Context) {
	var users []types.User
	if err := GormDB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
func AddNewUser(c *gin.Context) {
	var user types.User

	// รับค่าจาก Body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// insert ข้อมูลลง Database table users
	if err := GormDB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
func UpdateUserById(c *gin.Context) {
	// รับค่า id และ เช็คค่าว่ามีอยู่มั้ย ?
	id := c.Param("id")
	var user types.User
	if err := GormDB.Find(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// รับค่าที่ถูกส่งมาใน Body
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// update ค่าที่ถูกส่งมาใย Body
	if err := GormDB.Model(&user).Updates(&update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
func RemoveUserById(c *gin.Context) {
	id := c.Param("id")
	var user types.User
	if err := GormDB.Find(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := GormDB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %d deleted successfully", user.UserID),
	})
}
func main() {
	api := gin.Default()
	ConnectDB()
	CreatedTable()
	api.GET("/users", GetAllUsers)
	api.POST("/users/addUser", AddNewUser)
	api.PUT("/users/:id", UpdateUserById)
	api.DELETE("/users/:id", RemoveUserById)
	api.Run(":8080")
}
