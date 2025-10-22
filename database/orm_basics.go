package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORMUser model for GORM (extends the base User struct)
type GORMUser struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Age       int            `json:"age"`
	Profile   Profile        `gorm:"foreignKey:UserID" json:"profile"`
	Posts     []Post         `gorm:"foreignKey:UserID" json:"posts"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Profile model for one-to-one relationship
type Profile struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"not null" json:"user_id"`
	Bio      string `gorm:"type:text" json:"bio"`
	Website  string `gorm:"size:255" json:"website"`
	Location string `gorm:"size:100" json:"location"`
}

// Post model for one-to-many relationship
type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Published bool      `gorm:"default:false" json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ORMBasics demonstrates GORM operations
type ORMBasics struct {
	db *gorm.DB
}

// NewORMBasics creates a new ORMBasics instance
func NewORMBasics(db *gorm.DB) *ORMBasics {
	return &ORMBasics{db: db}
}

// ConnectPostgreSQL connects to PostgreSQL using GORM
func ConnectPostgreSQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Println("GORM PostgreSQL connection established")
	return db, nil
}

// ConnectMySQL connects to MySQL using GORM
func ConnectMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	log.Println("GORM MySQL connection established")
	return db, nil
}

// ConnectSQLite connects to SQLite using GORM
func ConnectSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	log.Println("GORM SQLite connection established")
	return db, nil
}

// AutoMigrate runs database migrations
func (o *ORMBasics) AutoMigrate() error {
	err := o.db.AutoMigrate(&GORMUser{}, &Profile{}, &Post{})
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// CreateUser demonstrates GORM Create operation
func (o *ORMBasics) CreateUser(name, email string, age int) (*GORMUser, error) {
	user := &GORMUser{
		Name:  name,
		Email: email,
		Age:   age,
	}

	result := o.db.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %w", result.Error)
	}

	log.Printf("User created with ID: %d", user.ID)
	return user, nil
}

// GetUserByID demonstrates GORM First operation
func (o *ORMBasics) GetUserByID(id uint) (*GORMUser, error) {
	var user GORMUser
	result := o.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// GetAllUsers demonstrates GORM Find operation
func (o *ORMBasics) GetAllUsers() ([]GORMUser, error) {
	var users []GORMUser
	result := o.db.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %w", result.Error)
	}

	log.Printf("Retrieved %d users", len(users))
	return users, nil
}

// UpdateUser demonstrates GORM Update operation
func (o *ORMBasics) UpdateUser(id uint, name, email string, age int) (*GORMUser, error) {
	var user GORMUser
	result := o.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	user.Name = name
	user.Email = email
	user.Age = age

	result = o.db.Save(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", result.Error)
	}

	log.Printf("User updated: %+v", user)
	return &user, nil
}

// DeleteUser demonstrates GORM Delete operation
func (o *ORMBasics) DeleteUser(id uint) error {
	result := o.db.Delete(&GORMUser{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	log.Printf("User with id %d deleted successfully", id)
	return nil
}

// SearchUsers demonstrates GORM Where operation
func (o *ORMBasics) SearchUsers(searchTerm string) ([]GORMUser, error) {
	var users []GORMUser
	result := o.db.Where("name ILIKE ? OR email ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search users: %w", result.Error)
	}

	log.Printf("Found %d users matching '%s'", len(users), searchTerm)
	return users, nil
}

// GetUsersByAgeRange demonstrates GORM Where with range
func (o *ORMBasics) GetUsersByAgeRange(minAge, maxAge int) ([]GORMUser, error) {
	var users []GORMUser
	result := o.db.Where("age BETWEEN ? AND ?", minAge, maxAge).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users by age range: %w", result.Error)
	}

	log.Printf("Found %d users between ages %d and %d", len(users), minAge, maxAge)
	return users, nil
}

// CreateUserWithProfile demonstrates GORM associations
func (o *ORMBasics) CreateUserWithProfile(name, email string, age int, bio, website, location string) (*GORMUser, error) {
	user := &GORMUser{
		Name:  name,
		Email: email,
		Age:   age,
		Profile: Profile{
			Bio:      bio,
			Website:  website,
			Location: location,
		},
	}

	result := o.db.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user with profile: %w", result.Error)
	}

	log.Printf("User with profile created with ID: %d", user.ID)
	return user, nil
}

// GetUserWithProfile demonstrates GORM Preload
func (o *ORMBasics) GetUserWithProfile(id uint) (*GORMUser, error) {
	var user GORMUser
	result := o.db.Preload("Profile").First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user with profile: %w", result.Error)
	}

	return &user, nil
}

// CreatePost demonstrates GORM association creation
func (o *ORMBasics) CreatePost(userID uint, title, content string) (*Post, error) {
	post := &Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Published: false,
	}

	result := o.db.Create(post)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create post: %w", result.Error)
	}

	log.Printf("Post created with ID: %d", post.ID)
	return post, nil
}

// GetUserWithPosts demonstrates GORM Preload with associations
func (o *ORMBasics) GetUserWithPosts(id uint) (*GORMUser, error) {
	var user GORMUser
	result := o.db.Preload("Posts").First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user with posts: %w", result.Error)
	}

	return &user, nil
}

// GetUserCount demonstrates GORM Count
func (o *ORMBasics) GetUserCount() (int64, error) {
	var count int64
	result := o.db.Model(&GORMUser{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count users: %w", result.Error)
	}

	log.Printf("Total users: %d", count)
	return count, nil
}

// GetUsersWithPagination demonstrates GORM Limit and Offset
func (o *ORMBasics) GetUsersWithPagination(page, pageSize int) ([]GORMUser, error) {
	var users []GORMUser
	offset := (page - 1) * pageSize

	result := o.db.Limit(pageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users with pagination: %w", result.Error)
	}

	log.Printf("Retrieved %d users for page %d", len(users), page)
	return users, nil
}

// SoftDeleteUser demonstrates GORM soft delete
func (o *ORMBasics) SoftDeleteUser(id uint) error {
	result := o.db.Delete(&GORMUser{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	log.Printf("User with id %d soft deleted successfully", id)
	return nil
}

// GetDeletedUsers demonstrates GORM Unscoped
func (o *ORMBasics) GetDeletedUsers() ([]GORMUser, error) {
	var users []GORMUser
	result := o.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get deleted users: %w", result.Error)
	}

	log.Printf("Found %d deleted users", len(users))
	return users, nil
}

// CleanupDatabase removes all records
func (o *ORMBasics) CleanupDatabase() error {
	// Delete in reverse order of dependencies
	if err := o.db.Unscoped().Delete(&Post{}, "1=1").Error; err != nil {
		return fmt.Errorf("failed to cleanup posts: %w", err)
	}

	if err := o.db.Unscoped().Delete(&Profile{}, "1=1").Error; err != nil {
		return fmt.Errorf("failed to cleanup profiles: %w", err)
	}

	if err := o.db.Unscoped().Delete(&GORMUser{}, "1=1").Error; err != nil {
		return fmt.Errorf("failed to cleanup users: %w", err)
	}

	log.Println("Database cleaned up successfully")
	return nil
}
