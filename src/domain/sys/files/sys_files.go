package files

import "time"

type SysFiles struct {
	ID        int64
	FileName  string
	FileMD5   string
	FilePath  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
type ISysFilesService interface {
}
