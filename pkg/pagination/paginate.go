package pagination

import (
	"math"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int    `json:"limit,omitempty" query:"limit"`
	Page       int    `json:"page,omitempty" query:"page"`
	TotalRows  int64  `json:"total_rows" query:"total_rows"`
	TotalPages int    `json:"total_pages" query:"total_pages"`
	SortBy     string `json:"sort_by,omitempty" query:"sort_by"`
	SortOrder  string `json:"sort_order,omitempty" query:"sort_order"`
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
