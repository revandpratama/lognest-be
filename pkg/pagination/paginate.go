package pagination

import (
	"math"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
	SortBy     string `json:"sort_by,omitempty"`
	SortOrder  string `json:"sort_order,omitempty"`
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func Paginate(db *gorm.DB, pagination *Pagination, model any, allowedSortColumns []string) *gorm.DB {
	var totalRows int64
	db.Model(model).Count(&totalRows)

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	isColumnAllowed := false
	if slices.Contains(allowedSortColumns, pagination.SortBy) {
		isColumnAllowed = true
	}

	if !isColumnAllowed {
		db = db.Order("created_at desc")
	} else {

		sortOrder := "DESC" // Default to a safe value
		if strings.ToUpper(pagination.SortOrder) == "ASC" {
			sortOrder = "ASC"
		}
		
		db = db.Order(pagination.SortBy + " " + sortOrder)
	}

	return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
}
