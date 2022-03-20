package manage

import (
	"errors"
	"gorm.io/gorm"
	"main.go/global"
	"main.go/model/common"
	"main.go/model/common/request"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"main.go/utils"
	"strconv"
	"time"
)

type GoodsCategoryService struct {
}

// GoodsCategoryServiceApp 提供全局调用
var GoodsCategoryServiceApp = new(GoodsCategoryService)

// AddCategory 添加商品分类
func (goodsCategoryService *GoodsCategoryService) AddCategory(req manageReq.MallGoodsCategoryReq) (err error) {
	if !errors.Is(global.GVA_DB.Where("category_level=? AND category_name=? AND is_deleted=0",
		req.CategoryLevel, req.CategoryName).First(&manage.MallGoodsCategory{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同分类")
	}

	rank, _ := strconv.Atoi(req.CategoryRank)
	category := manage.MallGoodsCategory{
		CategoryLevel: req.CategoryLevel,
		CategoryName:  req.CategoryName,
		CategoryRank:  rank,
		IsDeleted:     0,
		CreateTime:    common.JSONTime{Time: time.Now()},
		UpdateTime:    common.JSONTime{Time: time.Now()},
	}
	return global.GVA_DB.Create(&category).Error
}

// UpdateCategory 更新商品分类
func (goodsCategoryService *GoodsCategoryService) UpdateCategory(category manage.MallGoodsCategory) (err error) {
	if !errors.Is(global.GVA_DB.Where("category_level=? AND category_name=? AND is_deleted=0",
		category.CategoryLevel, category.CategoryName).First(&manage.MallGoodsCategory{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同分类")
	}
	return global.GVA_DB.Updates(&category).Error

}

// SelectCategoryPage 获取分类分页数据
func (goodsCategoryService *GoodsCategoryService) SelectCategoryPage(info manageReq.SearchCategoryParams) (err error, list interface{}, total int64) {
	limit := info.PageSize
	if limit > 1000 {
		limit = 1000
	}
	offset := info.PageSize * (info.PageNumber - 1)
	db := global.GVA_DB.Model(&manage.MallGoodsCategory{})
	var categoryList []manage.MallGoodsCategory

	if utils.NumsInList(info.CategoryLevel, []int{1, 2, 3}) {
		db.Where("category_level=?", info.CategoryLevel)
	}
	if info.ParentId >= 0 {
		db.Where("parent_id=?", info.ParentId)
	}
	err = db.Where("is_deleted=0").Count(&total).Error

	if err != nil {
		return err, categoryList, total

	} else {
		db = db.Where("is_deleted=0").Order("category_rank desc").Limit(limit).Offset(offset)
		err = db.Find(&categoryList).Error
	}
	return err, categoryList, total
}

// SelectCategoryById 获取单个分类数据
func (goodsCategoryService *GoodsCategoryService) SelectCategoryById(categoryId int) (err error, goodsCategory manage.MallGoodsCategory) {
	err = global.GVA_DB.Where("category_id=?", categoryId).First(&goodsCategory).Error
	return err, goodsCategory
}

// DeleteGoodsCategoriesByIds 批量设置失效
func (goodsCategoryService *GoodsCategoryService) DeleteGoodsCategoriesByIds(ids request.IdsReq) (err error, goodsCategory manage.MallGoodsCategory) {
	err = global.GVA_DB.Where("category_id in ?", ids.Ids).UpdateColumns(manage.MallGoodsCategory{IsDeleted: 1}).Error
	return err, goodsCategory
}

func (goodsCategoryService *GoodsCategoryService) SelectByLevelAndParentIdsAndNumber(parentId int, categoryLevel int) (err error, goodsCategories []manage.MallGoodsCategory) {
	err = global.GVA_DB.Where("category_id in ?", parentId).Where("category_level=?", categoryLevel).Where("is_deleted=0").Error
	return err, goodsCategories

}
