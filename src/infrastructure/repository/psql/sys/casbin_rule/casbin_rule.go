package casbinrule

// CasbinRule represents the casbin_rule table in the database
type CasbinRule struct {
	ID    int64  `gorm:"primaryKey;column:id;type:numeric(20,0)"` // 主键ID
	PType string `gorm:"column:ptype;type:varchar(100)"`          // 策略类型
	V0    string `gorm:"column:v0;type:varchar(100)"`             // 策略字段 v0
	V1    string `gorm:"column:v1;type:varchar(100)"`             // 策略字段 v1
	V2    string `gorm:"column:v2;type:varchar(100)"`             // 策略字段 v2
	V3    string `gorm:"column:v3;type:varchar(100)"`             // 策略字段 v3
	V4    string `gorm:"column:v4;type:varchar(100)"`             // 策略字段 v4
	V5    string `gorm:"column:v5;type:varchar(100)"`             // 策略字段 v5
}

// TableName returns the name of the database table for this model
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
