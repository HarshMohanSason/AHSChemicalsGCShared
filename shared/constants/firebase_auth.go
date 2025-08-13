package constants

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

// Map of valid roles.
var (
	Roles = map[string]struct{}{
		RoleSuperAdmin: {},
		RoleAdmin:      {},
		RoleUser:       {},
	}
)
