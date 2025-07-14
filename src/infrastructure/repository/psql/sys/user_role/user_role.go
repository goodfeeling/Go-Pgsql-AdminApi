package user_role

type SysUserRole struct {
	SysUserID uint64 `gorm:"column:sys_user_id;primaryKey" json:"sysUserId"`
	SysRoleID uint64 `gorm:"column:sys_role_id;primaryKey" json:"sysRoleId"`
}

func (SysUserRole) TableName() string {
	return "public.sys_user_role"
}
