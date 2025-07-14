package role_menu

import (
	"strconv"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysRoleMenu struct {
	SysBaseMenuID uint64 `gorm:"column:sys_base_menu_id;primaryKey" json:"sysBaseMenuId"`
	SysRoleID     uint64 `gorm:"column:sys_role_id;primaryKey" json:"sysRoleId"`
}

func (SysRoleMenu) TableName() string {
	return "public.sys_role_menus"
}

var ColumnsRoleMapping = map[string]string{
	"id":        "id",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

type ISysRoleMenuRepository interface {
	Insert(roleId int, UpdateMap map[string]any) error
	GetByRoleId(roleId int) ([]int, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

// GetByRoleId implements ISysRoleMenuRepository.
func (r *Repository) GetByRoleId(roleId int) ([]int, error) {
	var roleMenus []SysRoleMenu
	err := r.DB.Where("sys_role_id = ?", roleId).Find(&roleMenus).Error
	if err != nil {
		r.Logger.Error("Error getting all roles", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	var menuIds []int
	for _, roleMenu := range roleMenus {
		menuIds = append(menuIds, int(roleMenu.SysBaseMenuID))
	}
	return menuIds, nil
}

// Insert implements ISysRoleMenuRepository.
func (r *Repository) Insert(roleId int, UpdateMap map[string]any) error {
	menuIdsInterface, ok := UpdateMap["menuIds"]
	if !ok {
		r.Logger.Error("menuIds not found in update map")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	menuIdsInterfaceSlice, ok := menuIdsInterface.([]interface{})
	if !ok {
		r.Logger.Error("menuIds is not an array")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	roleMenus := make([]SysRoleMenu, 0, len(menuIdsInterfaceSlice))
	for _, item := range menuIdsInterfaceSlice {
		menuIdString, ok := item.(string)
		if !ok {
			r.Logger.Error("menuId is not a number")
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		menuId, err := strconv.ParseUint(menuIdString, 10, 64)
		if err != nil {
			r.Logger.Error("menuId is not a number")
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		roleMenu := SysRoleMenu{
			SysBaseMenuID: uint64(menuId),
			SysRoleID:     uint64(roleId),
		}
		roleMenus = append(roleMenus, roleMenu)
	}
	if err := r.DB.Model(&SysRoleMenu{}).
		Where("sys_role_id = ?", roleId).
		Delete(&SysRoleMenu{}).Error; err != nil {
		return err
	}
	if err := r.DB.Model(&SysRoleMenu{}).Create(&roleMenus).Error; err != nil {
		return err
	}
	return nil
}

func NewSysRoleMenuRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysRoleMenuRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
