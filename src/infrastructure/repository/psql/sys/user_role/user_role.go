package user_role

type SysUserRole struct {
	SysUserID int64 `gorm:"column:sys_user_id;primaryKey" json:"sysUserId"`
	SysRoleID int64 `gorm:"column:sys_role_id;primaryKey" json:"sysRoleId"`
}

func (SysUserRole) TableName() string {
	return "sys_user_roles"
}
