package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/initializers"
)

type Dataset struct {
	ID                string `gorm:"type:uuid;primaryKey" json:"id"`
	Name              string `gorm:"type:varchar(255)" json:"name"`
	Description       string `gorm:"type:text" json:"description"`
	CreatedAt         int64  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy         string `gorm:"type:varchar(255)" json:"created_by"`
	DataSourceType    string `gorm:"type:varchar(255)" json:"data_source_type"`
	IndexingTechnique string `gorm:"type:varchar(255)" json:"indexing_technique"`
	Permission        string `gorm:"type:varchar(255)" json:"permission"`
	Owner             string `gorm:"type:varchar(255)" json:"owner"`
}

// SaveDataset saves a given Dataset instance to the database
func SaveDataset(dataset *Dataset) error {
	db := initializers.DB

	var user User
	ownerUUID, err := uuid.Parse(dataset.Owner)
	if err != nil {
		return fmt.Errorf("invalid UUID format for owner: %v", err)
	}
	if err := db.First(&user, ownerUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}
	err = user.UpdateBalance(constants.UploadKnowledgePoints, "UploadKnowledgePoints")
	if err != nil {
		return err
	}

	return db.Create(dataset).Error
}

func GetDatasetByID(id string) (*Dataset, error) {
	var dataset Dataset
	db := initializers.DB
	if err := db.First(&dataset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dataset, nil
}

func DeleteDatasetByID(id string) error {
	db := initializers.DB

	// Check if the dataset exists
	var dataset Dataset
	if err := db.First(&dataset, "id = ?", id).Error; err != nil {
		return errors.New("dataset not found")
	}

	// Delete the dataset
	if err := db.Delete(&dataset).Error; err != nil {
		return err
	}

	return nil
}
