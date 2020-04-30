package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"encoding/csv"
	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type RoomModel struct {
	ID          uint
	Sensors     []Sensor
	Name        string
	Description string
	Url         string
	ImageUrl    string
}

type Sensor struct {
	ID              uint
	RoomModelID     uint
	Data            []Data
	MeshID          string
	Name            string
	Description     string
	MeasurementUnit string
}

type Data struct {
	ID       uint
	SensorID uint
	Value    float64
	Date     time.Time
}

// default mnemonic time
// https://pauladamsmith.com/blog/2011/05/go_time.html
const Layout = "2006-01-02 15:04:05"

var db *gorm.DB

func setupDatabase(drop bool) {
	// load env file
	//err := godotenv.Load("../database.env")
	//if err != nil {
	//	fmt.Println(err)
	//	panic("[!] Please create a 'database.env' file and prepare the needed variables.")
	//}

	// TODO check if sslmode can be enabled later
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))
	var err error
	db, err = gorm.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println(err)
		fmt.Println(dbinfo)
		panic("[!] failed to connect to db")
	} else {
		fmt.Println("[✓] successfully connected to db")
	}

	if drop {
		db.DropTableIfExists(&RoomModel{}, &Sensor{}, &Data{})
		fmt.Println("[✓] all data successfully dropped")
	}

	autoMigrate()
	fmt.Println("[✓] schemes migrated")
}

func autoMigrate() {
	// Migrate the Schema
	db.AutoMigrate(&RoomModel{}, &Sensor{}, &Data{})
}
func setupTestDatabase() *gorm.DB {
	testDb, err := gorm.Open("sqlite3", "./gorm_test.db")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	testDb.DB().SetMaxIdleConns(3)
	db = testDb
	autoMigrate()
	return testDb
}

func deleteTestDatabase(testDB *gorm.DB) {
	err := testDB.Close()
	err = os.Remove("./gorm_test.db")
	if err != nil {
		fmt.Println("unable to delete test db", err)
	}
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

func createMockData(sampleDataPath string) {
	var models []*RoomModel

	m1 := &RoomModel{
		Name:        "Facility Mechanical Room",
		Description: "This model shows a facility mechanical room with lots of pipes and stuff.",
		Url:         "files/facility-mechanical-room/model.zip",
		ImageUrl:    "files/facility-mechanical-room/thumbnail.png",
	}

	m2 := &RoomModel{
		Name:        "MEP Building Model",
		Description: "This model shows a MEP building with two floors and pipes.",
		Url:         "files/mep-building-model/model.zip",
		ImageUrl:    "files/mep-building-model/thumbnail.png",
	}

	m3 := &RoomModel{
		Name:        "Overhead MEP Installation",
		Description: "This model shows a overhead MEP installation. To make things look better this model has a longer description. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.",
		Url:         "files/overhead-mep-installation/model.zip",
		ImageUrl:    "files/overhead-mep-installation/thumbnail.png",
	}

	models = append(models, m1, m2, m3)
	meshIds := [][]string{
		{"node358", "node422", "node441", "node505"},
		{"node13", "node11", "node14", "node8"},
		{"1883", "1887", "9673", "10147"},
	}

	fmt.Println("[i] loading sensor data")

	for i, m := range models {
		db.Create(&m)
		s1 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Flow Sensor",
			Description:     "A basic flow sensor.",
			MeshID:          meshIds[i][0],
			MeasurementUnit: "°C",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_004_vorlauf_deg-celcius.csv", sampleDataPath)),
		}
		db.Create(&s1)

		s2 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Return Flow Sensor",
			Description:     "A basic return flow sensor with a longer description. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.",
			MeshID:          meshIds[i][1],
			MeasurementUnit: "°C",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_003_ruecklauf_deg-celcius.csv", sampleDataPath)),
		}
		db.Create(&s2)

		s3 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Fuel Sensor",
			Description:     "A basic thermal sensor",
			MeshID:          meshIds[i][2],
			MeasurementUnit: "l",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_002_fuel_litres.csv", sampleDataPath)),
		}
		db.Create(&s3)

		s4 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Pressure Sensor",
			Description:     "A basic thermal sensor",
			MeshID:          meshIds[i][3],
			MeasurementUnit: "bar",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_001_pressure_bar.csv", sampleDataPath)),
		}
		db.Create(&s4)
	}

	fmt.Println("[✓] finished loading sensor data")
}

func loadSampleData(path string) []Data {
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Println("[!] Error ", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var data []Data

	maxCount := 1000
	i := 0

	for {
		if i == maxCount {
			break
		}

		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			panic("[!] Error loading sample data")
		}

		t, _ := time.Parse(Layout, line[0]+" "+line[1])
		v, _ := strconv.ParseFloat(line[2], 64)

		data = append(data, Data{
			Date:  t,
			Value: v,
		})

		i++
	}
	return data
}
