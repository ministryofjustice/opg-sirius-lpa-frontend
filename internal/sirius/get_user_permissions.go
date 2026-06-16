package sirius

import "slices"

type PermissionType struct {
	Permissions  []string `json:"permissions"`
}

type Permissions map[string]PermissionType

func (c *Client) GetUserPermissions(ctx Context) (Permissions, error) {
	var v Permissions
	err := c.get(ctx, "/lpa-api/v1/permissions", &v)

	return v, err
}

func (p Permissions) Includes(permissionType string, method string) bool {
	return slices.Contains(p[permissionType].Permissions, method)
}
