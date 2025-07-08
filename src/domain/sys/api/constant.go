package api

import "sync"

var (
	ApiGroup ApiGroupType
	once     sync.Once
)

type ApiGroupResponse struct {
	ApiGroup ApiGroupType `json:"api_group"`
	Groups   []string     `json:"groups"`
}

type ApiGroupType struct {
	Auth             string `json:"auth"`
	Api              string `json:"api"`
	Role             string `json:"role"`
	Upload           string `json:"upload"`
	Menu             string `json:"menu"`
	Dictionary       string `json:"dictionary"`
	DictionaryDetail string `json:"dictionary_detail"`
	Operation        string `json:"operation"`
	User             string `json:"user"`
}

func initApiGroup() {
	ApiGroup = ApiGroupType{
		Auth:             "鉴权",
		Api:              "api",
		Role:             "角色",
		Upload:           "文件上传与下载",
		Menu:             "菜单",
		Dictionary:       "系统字典",
		DictionaryDetail: "系统字典详情",
		Operation:        "操作记录",
		User:             "系统用户",
	}
}

func GetApiGroup() ApiGroupType {
	once.Do(initApiGroup)
	return ApiGroup
}

func GetApiGroupNames() []string {
	group := GetApiGroup()
	return []string{
		group.Auth,
		group.Api,
		group.Role,
		group.Upload,
		group.Menu,
		group.Dictionary,
		group.DictionaryDetail,
		group.Operation,
		group.User,
	}
}
