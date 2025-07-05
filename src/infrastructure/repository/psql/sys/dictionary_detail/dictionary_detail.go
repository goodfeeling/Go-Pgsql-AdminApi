package dictionary_detail

import (
	"time"

	"gorm.io/gorm"
)

// SysDictionaryDetail represents the sys_dictionary_details table in the database
type SysDictionaryDetail struct {
	ID        int64          `gorm:"primaryKey;column:id;type:numeric(20,0)"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`  // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime:milli"`  // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`                 // 软删除标记

	Label           string `gorm:"column:label;type:varchar(191)"`              // 展示值
	Value           string `gorm:"column:value;type:varchar(191)"`              // 字典值
	Extend          string `gorm:"column:extend;type:varchar(191)"`             // 扩展值
	Status          int16  `gorm:"column:status;type:smallint"`                 // 启用状态
	Sort            int8   `gorm:"column:sort;type:bigint"`                     // 排序标记
	SysDictionaryID int64  `gorm:"column:sys_dictionary_id;type:numeric(20,0)"` // 关联字典ID
}

// TableName returns the name of the database table for this model
func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_details"
}
