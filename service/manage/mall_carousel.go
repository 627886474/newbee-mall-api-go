package manage

import (
	"errors"
	"gorm.io/gorm"
	"main.go/global"
	"main.go/model/common"
	"main.go/model/common/request"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"strconv"
	"time"
)

type MallCarouselService struct {
}

// CreateMallCarousel 创建MallCarousel记录
func (m *MallCarouselService) CreateMallCarousel(req manageReq.MallCarouselAddParam) (err error) {
	carouseRank, _ := strconv.Atoi(req.CarouselRank)
	mallCarousel := manage.MallCarousel{
		CarouselUrl:  req.CarouselUrl,
		RedirectUrl:  req.RedirectUrl,
		CarouselRank: carouseRank,
		CreateTime:   common.JSONTime{Time: time.Now()},
		UpdateTime:   common.JSONTime{Time: time.Now()},
	}

	err = global.GVA_DB.Create(&mallCarousel).Error
	return err
}

// DeleteMallCarousel 删除MallCarousel记录
func (m *MallCarouselService) DeleteMallCarousel(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&manage.MallCarousel{}, "carousel_id in ?", ids.Ids).Error
	return err
}

// UpdateMallCarousel 更新MallCarousel记录
func (m *MallCarouselService) UpdateMallCarousel(req manageReq.MallCarouselUpdateParam) (err error) {
	carouseRank, _ := strconv.Atoi(req.CarouselRank)
	if errors.Is(global.GVA_DB.Where("carousel_id = ?", req.CarouselId).First(&manage.MallCarousel{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("未查询到记录！")
	}
	err = global.GVA_DB.Where("carousel_id = ?", req.CarouselId).UpdateColumns(&manage.MallCarousel{
		CarouselUrl:  req.CarouselUrl,
		RedirectUrl:  req.RedirectUrl,
		CarouselRank: carouseRank,
		UpdateTime:   common.JSONTime{Time: time.Now()},
	}).Error
	return err
}

// GetMallCarousel 根据id获取MallCarousel记录
func (m *MallCarouselService) GetMallCarousel(id int) (err error, mallCarousel manage.MallCarousel) {
	err = global.GVA_DB.Where("carousel_id = ?", id).First(&mallCarousel).Error
	return
}

// GetMallCarouselInfoList 分页获取MallCarousel记录
func (m *MallCarouselService) GetMallCarouselInfoList(info manageReq.MallCarouselSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&manage.MallCarousel{})
	var mallCarousels []manage.MallCarousel
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("carousel_rank desc").Find(&mallCarousels).Error
	return err, mallCarousels, total
}
