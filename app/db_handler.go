package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type RoomModel struct {
	ID          uint
	Sensors     []Sensor
	Name        string
	Description string
	Path        string
	ImagePath   string
}

type Sensor struct {
	ID              uint
	RoomModelID     uint
	Data            []Data
	MeshID          uint
	Name            string
	Description     string
	MeasurementUnit string
}

type Data struct {
	ID       uint
	SensorID uint
	Value    int64
}

func setupDatabase(drop bool) *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("failed to connect to db")
	} else {
		fmt.Println("successfully connected to db")
	}

	if drop {
		db.DropTableIfExists(&RoomModel{}, &Sensor{}, &Data{})
		fmt.Println("all data successfully dropped")
	}

	// Migrate the Schema
	db.AutoMigrate(&RoomModel{}, &Sensor{}, &Data{})
	fmt.Println("all schemas migrated")

	return db
}

func createMockData(db *gorm.DB) {
	model := &RoomModel{
		Name:        "Basic",
		Description: "A basic model",
		Path:        "path_to_model",
		ImagePath:   "path_to_image",
	}

	// Create
	db.Create(&model)

	fmt.Println("model created")

	sensor := &Sensor{
		RoomModelID:     model.ID,
		Name:            "Heat Sensor",
		Description:     "A basic model",
		MeshID:          425,
		MeasurementUnit: "mm",
		Data: []Data{
			{Value: 3},
		},
	}

	db.Create(&sensor)

	fmt.Println("sensor w/ data created")
}
