package api

const (
	// API 分组常量
	GroupAuth             = "鉴权"
	GroupApi              = "api"
	GroupRole             = "角色"
	GroupUpload           = "文件上传与下载"
	GroupMenu             = "菜单"
	GroupDictionary       = "系统字典"
	GroupDictionaryDetail = "系统字典详情"
	GroupOperation        = "操作记录"
	GroupUser             = "系统用户"
	GroupOther            = "其它"
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
	Other            string `json:"other"`
}

func GetApiGroup() ApiGroupType {
	return ApiGroupType{
		Auth:             GroupAuth,
		Api:              GroupApi,
		Role:             GroupRole,
		Upload:           GroupUpload,
		Menu:             GroupMenu,
		Dictionary:       GroupDictionary,
		DictionaryDetail: GroupDictionaryDetail,
		Operation:        GroupOperation,
		User:             GroupUser,
		Other:            GroupOther,
	}
}

func GetApiGroupNames() []string {
	return []string{
		GroupAuth,
		GroupApi,
		GroupRole,
		GroupUpload,
		GroupMenu,
		GroupDictionary,
		GroupDictionaryDetail,
		GroupOperation,
		GroupUser,
		GroupOther,
	}
}
