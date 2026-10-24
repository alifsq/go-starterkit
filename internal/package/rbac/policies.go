package rbac

type Permission struct {
	Path   string
	Method string
}

var RolePermission = map[string][]Permission{
	"admin": {
		{Path: "/api/v1/user", Method: "GET"},
	},
	"user": {
		{Path: "/api/v1/user/me", Method: "GET"},
	},
}
