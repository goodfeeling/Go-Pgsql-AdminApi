package config

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainDictionaryDetail "github.com/gbrayhan/microservices-go/src/domain/sys/dictionary_detail"
)

type Config struct {
	ID            int64                                      `json:"id"`
	ConfigKey     string                                     `json:"config_key"`
	ConfigValue   string                                     `json:"config_value"`
	ConfigType    string                                     `json:"config_type"`
	Module        string                                     `json:"module"`
	EnvType       string                                     `json:"env_type"`
	Sort          int                                        `json:"sort"`
	CreatedAt     domain.CustomTime                          `json:"created_at"`
	UpdatedAt     domain.CustomTime                          `json:"updated_at"`
	SelectOptions *[]domainDictionaryDetail.DictionaryDetail `json:"select_options"`
}

type GroupConfig struct {
	Name    string   `json:"name"`
	Configs []Config `json:"configs"`
}

type IConfigService interface {
	GetConfigByGroup() (*[]GroupConfig, error)
	Update(module string, dataMap map[string]interface{}) error
	GetConfigByModule(module string) (*[]Config, error)
}
