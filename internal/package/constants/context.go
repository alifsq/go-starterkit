package constants

// ContextKey adalah custom string type agar type-safe saat dipakai di gin.Context
type ContextKey string

const (
	CtxUserID   ContextKey = "user_id"
	CtxUserRole ContextKey = "user_role"
	CtxEmail    ContextKey = "user_email"
)

// String mengembalikan representasi string bawaan agar kompatibel dengan c.Get/c.Set
func (c ContextKey) String() string {
	return string(c)
}