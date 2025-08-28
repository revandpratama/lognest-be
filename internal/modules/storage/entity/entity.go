package entity

import (
	"mime/multipart"
)

// Storage represents the data structure for a storage.
type Storage struct {
	PathName string                `json:"path_name" form:"path_name"`
	File     *multipart.FileHeader `form:"file"`
	Image    *multipart.FileHeader `form:"image"`
}

// TableName sets the table name for the Storage.
// func (Storage) TableName() string {
// 	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "storages")
// }

// func (p *Storage) BeforeCreate(tx *gorm.DB) error {
// 	if p.ID == uuid.Nil {
// 		uuidGenerated, err := uuid.NewV7()
// 		if err != nil {
// 			return err
// 		}
// 		p.ID = uuidGenerated
// 	}
// 	return nil
// }
