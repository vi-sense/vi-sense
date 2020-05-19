package model

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"encoding/csv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

//RoomModel specifies the structure for a single BIM model
type RoomModel struct {
	ID          uint		`json:"id"`
	Sensors     []Sensor	`json:"sensors"`
	Name        string		`json:"name"`
	Description string		`json:"description"`
	Url         string		`json:"url"`
	ImageUrl    string		`json:"image_url"`
}

//Sensor specifies the structure for a single sensor which is located inside a RoomModel
type Sensor struct {
	ID              uint	`json:"id"`
	RoomModelID     uint	`json:"room_model_id"`
	Data            []Data	`json:"-"`
	MeshID          string	`json:"mesh_id"`
	Name            string	`json:"name"`
	Description     string	`json:"description"`
	MeasurementUnit string	`json:"measurement_unit"`
}

//Data specifies the structure for a single measured value w/ timestamp which was recorded by a sensor
type Data struct {
	ID       uint			`json:"id"`
	SensorID uint			`json:"sensor_id"`
	Value    float64		`json:"value"`
	Date     time.Time		`json:"date"`
}

// default mnemonic time
// https://pauladamsmith.com/blog/2011/05/go_time.html
const Layout = "2006-01-02 15:04:05"

var DB *gorm.DB

//SetupDatabase initializes the database w/ the orm mapping and postgres as the dialect;
//drop defines whether or not the currently active scheme should be dropped
func SetupDatabase(drop bool) {
	// TODO check if ssl mode can be enabled later
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))
	var err error
	DB, err = gorm.Open("postgres", dbInfo)

	if err != nil {
		fmt.Println(err)
		fmt.Println(dbInfo)
		panic("[!] failed to connect to db")
	} else {
		fmt.Println("[✓] successfully connected to db")
	}

	if drop {
		DB.DropTableIfExists(&RoomModel{}, &Sensor{}, &Data{})
		fmt.Println("[✓] all data successfully dropped")
	}

	// Migrate the Schema
	DB.AutoMigrate(&RoomModel{}, &Sensor{}, &Data{})
	fmt.Println("[✓] schemes migrated")
}

func CreateMockData(sampleDataPath string, dataLimit int) {
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
		DB.Create(&m)
		s1 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Flow Sensor",
			Description:     "A basic flow sensor.",
			MeshID:          meshIds[i][0],
			MeasurementUnit: "°C",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_004_vorlauf_deg-celcius.csv", sampleDataPath), dataLimit),
		}
		DB.Create(&s1)

		s2 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Return Flow Sensor",
			Description:     "A basic return flow sensor with a longer description. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.",
			MeshID:          meshIds[i][1],
			MeasurementUnit: "°C",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_003_ruecklauf_deg-celcius.csv", sampleDataPath), dataLimit),
		}
		DB.Create(&s2)

		s3 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Fuel Sensor",
			Description:     "A basic thermal sensor",
			MeshID:          meshIds[i][2],
			MeasurementUnit: "l",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_002_fuel_litres.csv", sampleDataPath), dataLimit),
		}
		DB.Create(&s3)

		s4 := &Sensor{
			RoomModelID:     m.ID,
			Name:            "Pressure Sensor",
			Description:     "A basic thermal sensor",
			MeshID:          meshIds[i][3],
			MeasurementUnit: "bar",
			Data:            loadSampleData(fmt.Sprintf("%s/sensors/sensor_001_pressure_bar.csv", sampleDataPath), dataLimit),
		}
		DB.Create(&s4)
	}

	fmt.Println("[✓] finished loading sensor data")
}

//SetupTestDatabase creates local sqlite db for testing
func SetupTestDatabase() {
	var err error
	DB, err = gorm.Open("sqlite3", "./gorm_test.db")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	DB.DB().SetMaxIdleConns(3)
	// Migrate the Schema
	DB.AutoMigrate(&RoomModel{}, &Sensor{}, &Data{})
}

//DeleteTestDatabase deletes local sqlite db for testing
func DeleteTestDatabase() {
	err := DB.Close()
	err = os.Remove("./gorm_test.db")
	if err != nil {
		fmt.Println("unable to delete test db", err)
	}
}

func loadSampleData(path string, dataLimit int) []Data {
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Println("[!] Error ", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var data []Data

	i := 0

	for {
		if i == dataLimit {
			break
		}

		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			panic("[!] Error loading sample data")
		}

		if len(line) < 3 {
			fmt.Println("[!] Line skipped because there is not enough data prepared.", err)
			fmt.Printf("[!] Line: %s", line)
			continue
		}

		s := fmt.Sprintf("%s %s", line[0], line[1])
		t, _ := time.Parse(Layout, s)
		v, _ := strconv.ParseFloat(line[2], 64)

		data = append(data, Data{
			Date:  t,
			Value: v,
		})

		i++
	}
	return data
}
