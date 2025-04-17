package response

type ListUser struct {
	ID        string `bun:"id" json:"id"`
	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	Email     string `bun:"email" json:"email"`
	CreatedAt int64  `bun:"created_at" json:"created_at"`
	UpdatedAt int64  `bun:"updated_at" json:"updated_at"`
}
