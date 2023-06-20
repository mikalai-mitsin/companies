package entity

type PermissionID string

func (p PermissionID) String() string {
	return string(p)
}
