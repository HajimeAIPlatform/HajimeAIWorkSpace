package models

import (
    "errors"
    "github.com/google/uuid"
    "time"

    "gorm.io/gorm"
)

type UploadFiles struct {
    ID            uuid.UUID  `json:"id"`
    TenantID      uuid.UUID  `json:"tenant_id"`
    StorageType   string     `json:"storage_type"`
    Key           string     `json:"key"`
    Name          string     `json:"name"`
    Size          int        `json:"size"`
    Extension     string     `json:"extension"`
    MimeType      string     `json:"mime_type,omitempty"`
    CreatedBy     uuid.UUID  `json:"created_by"`
    CreatedAt     time.Time  `json:"created_at"`
    Used          bool       `json:"used"`
    UsedBy        *uuid.UUID `json:"used_by,omitempty"`
    UsedAt        *time.Time `json:"used_at,omitempty"`
    Hash          string     `json:"hash,omitempty"`
    CreatedByRole string     `json:"created_by_role"`
}

// QueryStorageByID queries a single storage record by ID
func QueryStorageByID(db *gorm.DB, ID string) (*UploadFiles, error) {
    var difyFiles UploadFiles
    result := db.Where("id = ?", ID).First(&difyFiles)

    if result.Error != nil {
        return nil, result.Error
    }

    if result.RowsAffected == 0 {
        return nil, errors.New("no record found with the given ID")
    }

    return &difyFiles, nil
}
