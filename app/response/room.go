package response

type ListRoom struct {
	ID   string `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
	HasPassword bool `bun:"has_password" json:"has_password"`
}
