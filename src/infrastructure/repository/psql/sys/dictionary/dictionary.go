package dictionary

import (
	"time"

	"gorm.io/gorm"
)

// SysDictionary represents the sys_dictionaries table in the database
type SysDictionary struct {
	ID        int64          `gorm:"primaryKey;column:id;type:numeric(20,0)"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`  // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime:milli"`  // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`                 // 软删除标记

	Name   string `gorm:"column:name;type:varchar(191)"` // 字典名（中）
	Type   string `gorm:"column:type;type:varchar(191)"` // 字典名（英）
	Status int16  `gorm:"column:status;type:smallint"`   // 状态
	Desc   string `gorm:"column:desc;type:varchar(191)"` // 描述
}

// TableName returns the name of the database table for this model
func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}
