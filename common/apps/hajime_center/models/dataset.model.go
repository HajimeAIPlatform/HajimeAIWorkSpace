package models

import (
	"time"

	"github.com/google/uuid"
)

type Dataset struct {
	ID                     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name                   string    `gorm:"type:varchar(255);not null" json:"name"`
	Description            string    `gorm:"type:text" json:"description,omitempty"`
	Provider               string    `gorm:"type:varchar(255);not null;default:'vendor'" json:"provider"`
	Permission             string    `gorm:"type:varchar(255);not null;default:'only_me'" json:"permission"`
	DataSourceType         string    `gorm:"type:varchar(255);default:'upload_file'" json:"data_source_type"`
	CreatedBy              string    `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt              time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy              string    `gorm:"type:uuid" json:"updated_by,omitempty"`                                         // 使用指针类型，允许 NULL
	UpdatedAt              time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"` // 使用指针类型，允许 NULL
	EmbeddingModel         string    `gorm:"type:varchar(255)" json:"embedding_model,omitempty"`
	EmbeddingModelProvider string    `gorm:"type:varchar(255)" json:"embedding_model_provider,omitempty"`
}

type Document struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	DatasetID           uuid.UUID `gorm:"type:uuid;not null;index:idx_document_dataset" json:"dataset_id"`
	Position            int       `gorm:"type:int;not null" json:"position"`
	Name                string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedBy           string    `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt           time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ProcessingStartedAt time.Time `gorm:"type:timestamp" json:"processing_started_at,omitempty"` // 使用指针类型，允许 NULL
	FileID              uuid.UUID `gorm:"type:uuid" json:"file_id,omitempty"`                    // 使用指针类型，允许 NULL
	ParsingCompletedAt  time.Time `gorm:"type:timestamp" json:"parsing_completed_at,omitempty"`  // 使用指针类型，允许 NULL
	CompletedAt         time.Time `gorm:"type:timestamp" json:"completed_at,omitempty"`          // 使用指针类型，允许 NULL
	UpdatedAt           time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DocType             string    `gorm:"type:varchar(40)" json:"doc_type,omitempty"`
	IsUploadToServer    bool      `gorm:"type:boolean;default:false" json:"is_upload_to_server"`
	IndexingStatus      string    `gorm:"type:varchar(40);default:'indexing" json:"indexing_status"`
}

type UploadFile struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	StorageType   string     `gorm:"type:varchar(255);not null" json:"storage_type"`
	Key           string     `gorm:"type:varchar(255);not null" json:"key"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name"`
	Size          int64      `gorm:"type:int;not null" json:"size"`
	Extension     string     `gorm:"type:varchar(255);not null" json:"extension"`
	MimeType      string     `gorm:"type:varchar(255)" json:"mime_type"`
	CreatedByRole string     `gorm:"type:varchar(255);not null;default:'account'" json:"created_by_role"`
	CreatedBy     string     `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Used          bool       `gorm:"type:boolean;not null;default:false" json:"used"`
	UsedBy        *string    `gorm:"type:uuid" json:"used_by,omitempty"`      // 使用指针类型，允许 NULL
	UsedAt        *time.Time `gorm:"type:timestamp" json:"used_at,omitempty"` // 使用指针类型，允许 NULL
}

type DatasetInfoResult struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Permission     string    `json:"permission"`
	DataSourceType string    `json:"data_source_type"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      int       `json:"created_at"`
}

type DocumentInfoResult struct {
	ID             uuid.UUID `json:"id"`
	Position       int       `json:"position"`
	DataSourceInfo struct {
		UploadFileID string `json:"upload_file_id"`
	} `json:"data_source_info"`
	Name        string `json:"name"`
	CreatedFrom string `json:"created_from"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   int    `json:"created_at"`
}

type FileUploadResult struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Size      int64              `json:"size"`
	Extension string             `json:"extension"`
	MimeType  string             `json:"mime_type"`
	CreatedBy string             `json:"created_by"`
	CreatedAt int                `json:"created_at"`
	Dataset   DatasetInfoResult  `json:"dataset"`
	Documents DocumentInfoResult `json:"documents"`
}

type FileUploadForChatResult struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
	MimeType  string    `json:"mime_type"`
	CreatedBy string    `json:"created_by"`
	CreatedAt int       `json:"created_at"`
}

type RenameDocumentRequest struct {
	Name string `json:"name" binding:"required"`
}
