package dto

import (
	"github.com/toheart/functrace"
	"gorm.io/gorm"
)

type Pagination struct {
	PageIndex int `form:"page_index"`
	PageSize  int `form:"page_size"`
}

func (m *Pagination) GetPageIndex() int {
	defer functrace.Trace([]interface {
	}{m})()
	if m.PageIndex <= 0 {
		m.PageIndex = 1
	}
	return m.PageIndex
}

func (m *Pagination) GetPageSize() int {
	defer functrace.Trace([]interface {
	}{m})()
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	defer functrace.Trace([]interface {
	}{pageSize, pageIndex})()
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

type Base struct {
	ID string `json:"id"`
}
