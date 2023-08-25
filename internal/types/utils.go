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

func FromShare(s *model.Share, passwordIncluded bool) *ShareDTO {
	var expireAt *int64
	if s.ExpireAt.Valid {
		tm := s.ExpireAt.Time.UnixMilli()
		expireAt = &tm
	}
	if s != nil {
		dto := &ShareDTO{
			ID:        s.ID,
			UUID:      s.UUID,
			ExpireAt:  expireAt,
			File:      *FromFile(&s.File),
			Username:  s.User.Name,
			CreatedAt: s.CreatedAt.UnixMilli(),
			UpdatedAt: s.UpdatedAt.UnixMilli(),
		}
		if passwordIncluded {
			dto.Password = s.Password
		}
		return dto
	}
	return nil
}
