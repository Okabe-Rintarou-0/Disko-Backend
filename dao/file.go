package dao

import (
	"disko/model"
	"disko/repository/query"
	"fmt"
)

type IFileDAO interface {
	FindById(id uint) (*model.File, error)
	FindByOwnerAndParentAndName(owner uint, parent *uint, name string) (*model.File, error)
	FindByUUID(uuid string) (*model.File, error)
	Save(files ...*model.File) error
	FindByOwner(ownerID uint) ([]*model.File, error)
	Search(owner *uint, parent *uint, keyword string, extensions []string) ([]*model.File, error)
	DeleteById(id uint) error
}

func NewFileDAO() IFileDAO {
	return &FileDAO{}
}

type FileDAO struct{}

func (fd *FileDAO) DeleteById(id uint) error {
	f := query.Use(db).File
	_, err := f.WithContext(ctx).Where(f.ID.Eq(id)).Delete()
	return err
}

func (fd *FileDAO) FindByOwnerAndParentAndName(owner uint, parent *uint, name string) (*model.File, error) {
	f := query.Use(db).File
	q := f.WithContext(ctx).Where(f.Owner.Eq(owner), f.Name.Eq(name))
	if parent != nil {
		q = q.Where(f.ParentID.Eq(*parent))
	}
	return q.Take()
}

func (fd *FileDAO) Search(owner *uint, parent *uint, keyword string, extensions []string) ([]*model.File, error) {
	f := query.Use(db).File
	q := f.WithContext(ctx)

	if owner != nil {
		q = q.Where(f.Owner.Eq(*owner))
	}

	if parent != nil {
		q = q.Where(f.ParentID.Eq(*parent))
	} else {
		q = q.Where(f.ParentID.IsNull())
	}

	if len(keyword) > 0 {
		// %<keyword>%
		q = q.Where(f.Name.Like(fmt.Sprintf("%%%s%%", keyword)))
	}

	if len(extensions) > 0 {
		q = q.Where(f.Ext.In(extensions...))
	}

	return q.Find()
}

func (fd *FileDAO) FindById(id uint) (*model.File, error) {
	f := query.Use(db).File
	return f.WithContext(ctx).Where(f.ID.Eq(id)).Take()
}

func (fd *FileDAO) FindByUUID(uuid string) (*model.File, error) {
	f := query.Use(db).File
	return f.WithContext(ctx).Where(f.UUID.Eq(uuid)).Take()
}

func (fd *FileDAO) Save(files ...*model.File) error {
	f := query.Use(db).File
	return f.WithContext(ctx).Save(files...)
}

func (fd *FileDAO) FindByOwner(ownerID uint) ([]*model.File, error) {
	f := query.Use(db).File
	return f.WithContext(ctx).Where(f.Owner.Eq(ownerID)).Find()
}
