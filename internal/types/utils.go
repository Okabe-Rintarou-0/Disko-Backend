package types

import "disko/model"

func FromFile(f *model.File) *FileDTO {
	if f != nil {
		return &FileDTO{
			ID:      f.ID,
			Name:    f.Name,
			Ext:     f.Ext,
			Size:    f.Size,
			UUID:    f.UUID,
			Owner:   f.Owner,
			IsDir:   f.IsDir,
			Private: f.Private,
			Parent:  f.ParentID,
		}
	}
	return nil
}
