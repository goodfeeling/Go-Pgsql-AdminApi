package files

type SysFiles struct {
	ID       int64  `json:"id"`
	FileName string `json:"file_name"`
	FileMD5  string `json:"file_md5"`
	FilePath string `json:"file_path"`
	FileUrl  string `json:"file_url"`
}
type ISysFilesService interface {
	Create(data *SysFiles) (*SysFiles, error)
}
