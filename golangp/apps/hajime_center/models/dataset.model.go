package models

import "hajime/golangp/apps/hajime_center/initializers"

type Dataset struct {
	ID                string `gorm:"type:varchar(255);primaryKey" json:"id"`
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
	return db.Create(dataset).Error
}

// GetAllDatasets retrieves all Dataset instances from the database
func GetAllDatasets() ([]Dataset, error) {
	var datasets []Dataset
	db := initializers.DB
	if err := db.Find(&datasets).Error; err != nil {
		return nil, err
	}
	return datasets, nil
}

func GetDatasetByID(id string) (*Dataset, error) {
	var dataset Dataset
	db := initializers.DB
	if err := db.First(&dataset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dataset, nil
}
