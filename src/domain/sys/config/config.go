package config

import (
	"github.com/gbrayhan/microservices-go/src/domain"
)

type Config struct {
	ID          int64             `json:"id"`
	ConfigKey   string            `json:"config_key"`
	ConfigValue string            `json:"config_vValue"`
	ConfigType  string            `json:"config_type"`
	Description string            `json:"description"`
	Module      string            `json:"module"`
	EnvType     string            `json:"env_type"`
	IsEnabled   bool              `json:"isEnabled"`
	CreatedAt   domain.CustomTime `json:"created_at"`
	UpdatedAt   domain.CustomTime `json:"updated_at"`
}

type GroupConfig struct {
	Name    string   `json:"name"`
	Configs []Config `json:"configs"`
}

type IConfigService interface {
	GetConfigByGroup() (*[]GroupConfig, error)
	Update(dataMap map[string]interface{}) (*Config, error)
	GetConfigByModule(module string) (*[]Config, error)
}
