package tools

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Param struct {
	DB        *gorm.DB
	PageIndex int
	PageSize  int
	OrderBy   []string
	ShowSQL   bool
}

type Pagination struct {
	CurrentPage int `json:"current_page" form:"current_page"`
	PageSize    int `json:"page_size" form:"page_size"`
	LastPage    int `json:"last_page"`
	Total       int `json:"total" form:"total"`
}

func Paging(p Param, result interface{}) (Pagination, error) {
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			p.DB = p.DB.Order(o)
		}
	}

	//if p.PageIndex == 0 && p.PageSize == 0 {
	//	if err := p.DB.Find(result).Error; err != nil {
	//		logrus.Errorf("Paging db get record err: %v", err.Error())
	//		return Pagination{}, err
	//	}
	//	return Pagination{}, nil
	//}

	pagination := Pagination{}
	if p.PageIndex <= 0 {
		p.PageIndex = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 50
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			p.DB = p.DB.Order(o)
		}
	}

	var totalCount int64
	err := p.DB.Count(&totalCount).Error
	if err != nil {
		logrus.Errorf("Paging db get count err: %v", err.Error())
		return pagination, err
	}
	pagination.Total = int(totalCount)
	pagination.LastPage = int(totalCount)/p.PageSize + 1
	if p.PageIndex > pagination.LastPage {
		p.PageIndex = pagination.LastPage
	}

	if err := p.DB.Limit(p.PageSize).Offset((p.PageIndex - 1) * p.PageSize).Find(result).Error; err != nil {
		logrus.Errorf("Paging db get record err: %v", err.Error())
		return pagination, err
	}

	pagination.CurrentPage = p.PageIndex
	pagination.PageSize = p.PageSize
	return pagination, nil
}
