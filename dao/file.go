package dao

import (
	"disko/model"
	"disko/repository/query"
	"fmt"
)

type IFileDAO interface {
	FindById(id uint) (*model.File, error)
	FindByUUID(uuid string) (*model.File, error)
	Save(files ...*model.File) error
	FindByOwner(ownerID uint) ([]*model.File, error)
	Search(owner *uint, parent *uint, keyword string, extensions []string) ([]*model.File, error)
}

func NewFileDAO() IFileDAO {
	return &FileDAO{}
}

type FileDAO struct{}

func (fd *FileDAO) Search(owner *uint, parent *uint, keyword string, extensions []string) ([]*model.File, error) {
	f := query.Use(db).File
	q := f.WithContext(ctx)

	if owner != nil {
		q = q.Where(f.Owner.Eq(*owner))
	}

	if parent != nil {
		q = q.Where(f.ParentID.Eq(*parent))
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
