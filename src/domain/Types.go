package domain

import "time"

type DateRangeFilter struct {
	Field string     `json:"field"`
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

func (sd SortDirection) IsValid() bool {
	return sd == SortAsc || sd == SortDesc
}

type DataFilters struct {
	LikeFilters      map[string][]string `json:"likeFilters"`
	Matches          map[string][]string `json:"matches"`
	DateRangeFilters []DateRangeFilter   `json:"dataRanges"`
	SortBy           []string            `json:"sortBy"`
	SortDirection    SortDirection       `json:"sortDirection"`
	Page             int                 `json:"page"`
	PageSize         int                 `json:"pageSize"`
}

type CommonResponse[T interface{}] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
type PageList[T interface{}] struct {
	List       T           `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_page"`
	Filters    DataFilters `json:"filters"`
}
