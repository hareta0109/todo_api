package response

type AuthLogin struct {
	Token     string
	CompanyID uint64
	UserID    uint64
}
