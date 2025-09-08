package main

import (
	// import Types
	"backend/types"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
    department_id SERIAL PRIMARY KEY,
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
func InsertDataDepartment() {
	// ข้อมูลเริ่มต้นพร้อม ID
	departments := []types.Department{
		{DepartmentID: 1, DepartmentName: "Research & Development"},
		{DepartmentID: 2, DepartmentName: "Software Engineering"},
		{DepartmentID: 3, DepartmentName: "IT Support"},
		{DepartmentID: 4, DepartmentName: "Cybersecurity"},
		{DepartmentID: 5, DepartmentName: "Cloud Services"},
		{DepartmentID: 6, DepartmentName: "Human Resources"},
		{DepartmentID: 7, DepartmentName: "Marketing & Sales"},
		{DepartmentID: 8, DepartmentName: "Finance & Accounting"},
		{DepartmentID: 9, DepartmentName: "Legal & Compliance"},
		{DepartmentID: 10, DepartmentName: "Customer Success"},
	}

	for _, dept := range departments {
		// เช็คว่า ID มีอยู่แล้วหรือยัง ถ้าไม่มีก็ insert
		var existing types.Department
		if err := GormDB.First(&existing, dept.DepartmentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := GormDB.Create(&dept).Error; err != nil {
					fmt.Println("❌ Failed to insert department:", dept.DepartmentName, err)
				} else {
					fmt.Println("✅ Inserted department:", dept.DepartmentName)
				}
			} else {
				fmt.Println("❌ Error checking department:", err)
			}
		} else {
			fmt.Println("ℹ️ Department already exists:", dept.DepartmentName)
		}
	}
}

func GetAllUsers(c *gin.Context) {
	var users []types.User
	if err := GormDB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
func GetAllDepartment(c *gin.Context) {
	var department []types.Department
	if err := GormDB.Model(&types.Department{}).
		Select("department_id, department_name").
		Find(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, department)
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
func FindDepartmentById(c *gin.Context) {
	id := c.Param("id")
	var department types.Department

	// ใช้ First เพื่อหา 1 record ตาม id
	if err := GormDB.First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	c.JSON(http.StatusOK, department)
}

func FindUserById(c *gin.Context) {
	id := c.Param("id")
	var user types.User

	// ใช้ First เช่นกัน
	if err := GormDB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func main() {
	api := gin.Default()

	// ตั้งค่า CORS
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Angular frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	ConnectDB()
	CreatedTable()
	InsertDataDepartment()
	api.GET("/users", GetAllUsers)
	api.GET("/users/:id", FindUserById)
	api.POST("/users/addUser", AddNewUser)
	api.PUT("/users/:id", UpdateUserById)
	api.DELETE("/users/:id", RemoveUserById)

	api.GET("/department", GetAllDepartment)
	api.GET("/department/:id", FindDepartmentById)
	api.Run(":8080")
}
