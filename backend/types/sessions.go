package types

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int    `json:"userID"`
	Token     string `json:"token"`
	TokenHash string `json:"tokenHash"`
}
