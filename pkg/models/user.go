package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt string    `json:"created_at"`
}

type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	Description string    `json:"description"`
	PublishedAt string    `json:"published_at"`
	ImageURL    string    `json:"image_url"`
	UserID      uuid.UUID `json:"user_id"`
	Reviews     []Review  `json:"reviews"`
}

type Borrow struct {
	ID        uuid.UUID `json:"id"`         // Unique borrowing ID
	BookID    uuid.UUID `json:"book_id"`    // ID of the borrowed book
	UserID    uuid.UUID `json:"user_id"`    // ID of the borrowing user
	Status    string    `json:"status"`     // 'pending', 'approved', 'rejected'
	DueDate   string    `json:"due_date"`   // Return due date
	CreatedAt string    `json:"created_at"` // Request creation time
}

type Review struct {
	ID        string `json:"id"`      // Assuming UUID or string for ID
	BookID    string `json:"book_id"` // Assuming UUID or string for BookID
	UserID    string `json:"user_id"` // Assuming UUID or string for UserID
	Rating    int    `json:"rating"`  // Assuming rating is an integer
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}
