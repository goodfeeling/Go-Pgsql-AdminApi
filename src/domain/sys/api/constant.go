package api

import "sync"

var (
	ApiGroup map[string]string
	once     sync.Once
)

type ApiGroupType struct {
	ApiGroup map[string]string `json:"api_group"`
	Names    []string          `json:"groups"`
}

func initApiGroup() {
	ApiGroup = map[string]string{
		"auth":             "鉴权",
		"api":              "api",
		"role":             "角色",
		"upload":           "文件上传与下载",
		"menu":             "菜单",
		"dictionary":       "系统字典",
		"dictionaryDetail": "系统字典详情",
		"operation":        "操作记录",
		"user":             "系统用户",
	}
}

func GetApiGroup() map[string]string {
	once.Do(initApiGroup)
	return ApiGroup
}

func GetApiGroupNames() []string {
	once.Do(initApiGroup)
	var names []string
	for _, name := range ApiGroup {
		names = append(names, name)
	}
	return names
}
