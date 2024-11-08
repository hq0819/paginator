package paginator

import (
	"gorm.io/gorm"
	"math"
)

type BaseModel struct {
	PageNum  int64  `json:"pageNum"`
	PageSize int64  `json:"PageSize"`
	OrderBy  string `json:"OrderBy"`
}

type PageInfo[T any] struct {
	PageNum  int64 `json:"pageNum"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
	Rows     []T   `json:"rows"`
	LastPage bool  `json:"lastPage"`
}

type Paginator[T any] struct {
}

func (p *Paginator[T]) StartPage(db *gorm.DB, baseModel BaseModel) (*PageInfo[T], error) {
	total := new(int64)
	pageInfo := new(PageInfo[T])
	pageInfo.PageNum = baseModel.PageNum
	pageInfo.PageSize = baseModel.PageSize
	pageInfo.LastPage = true
	err := db.Count(total).Error
	pageInfo.Total = *total
	if err != nil {
		return pageInfo, err
	}

	if baseModel.PageNum == 0 {
		baseModel.PageNum = 1
	}

	if baseModel.PageNum*baseModel.PageSize < *total {
		pageInfo.LastPage = false
	}
	n := math.Ceil(float64(pageInfo.Total) / float64(baseModel.PageSize))

	if float64(pageInfo.PageNum) <= n {
		rows := make([]T, 0, baseModel.PageSize)
		db.Order(baseModel.OrderBy)
		err := db.Offset(int((baseModel.PageNum - 1) * baseModel.PageSize)).Limit(int(baseModel.PageSize)).Find(&rows).Error
		if err != nil {
			return pageInfo, err
		}
		pageInfo.Rows = rows

	}
	return pageInfo, nil

}
