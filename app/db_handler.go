package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

var db *gorm.DB

func setupDatabase(drop bool) {
	// load env file
	//err := godotenv.Load("../database.env")
	//if err != nil {
	//	fmt.Println(err)
	//	panic("[!] Please create a 'database.env' file and prepare the needed variables.")
	//}

	// TODO check if sslmode can be enabled later
	dbinfo := fmt.Sprintf("host=localhost port=32300 user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))

	db, err := gorm.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println(err)
		panic("[!] failed to connect to db")
	} else {
		fmt.Println("[✓] successfully connected to db")
	}

	if drop {
		db.DropTableIfExists(&RoomModel{}, &Sensor{}, &Data{})
		fmt.Println("[✓] all data successfully dropped")
	}

	// Migrate the Schema
	db.AutoMigrate(&RoomModel{}, &Sensor{}, &Data{})
	fmt.Println("[✓] schemes migrated")
}

// query all room models
func queryRoomModels() string {
	var q []RoomModel
	db.Find(&q)
	return asJson(&q)
}

// query room model by id
// - returns model with sensors
func queryRoomModel(id string) string {
	var q RoomModel
	db.Preload("Sensors").First(&q, id)
	return asJson(&q)
}

// TODO check if needed
// query all sensors
func querySensors() string {
	var q []Sensor
	db.Find(&q)
	return asJson(&q)
}

// query sensor by id
// - returns sensor with sensor data
func querySensor(id string) string {
	var q Sensor
	db.Preload("Data").First(&q, id)
	return asJson(&q)
}

func asJson(obj interface{}) string {
	b, err := json.Marshal(&obj)
	if err != nil {
		fmt.Println("[!]", err)
		return ""
	}

	return string(b)
}

func createMockData() {
	m1 := &RoomModel{
		Name:        "Basic 1",
		Description: "A basic model",
		Path:        "path_to_model",
		ImagePath:   "path_to_image",
	}

	m2 := &RoomModel{
		Name:        "Basic 2",
		Description: "A basic model",
		Path:        "path_to_model",
		ImagePath:   "path_to_image",
	}

	// Create
	db.Create(&m1)
	db.Create(&m2)

	fmt.Println("[✓] mock room model created")

	s1 := &Sensor{
		RoomModelID:     m1.ID,
		Name:            "Thermal Sensor",
		Description:     "A basic thermal sensor",
		MeshID:          425,
		MeasurementUnit: "°C",
		Data: []Data{
			{Value: 3},
		},
	}

	db.Create(&s1)
	fmt.Println("[✓] mock sensor w/ data created")
}
