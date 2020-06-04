package model

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"encoding/csv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//RoomModel specifies the structure for a single BIM model
type RoomModel struct {
	ID          uint     `json:"id"`
	Sensors     []Sensor `json:"sensors"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	ImageUrl    string   `json:"image_url"`
}

//Sensor specifies the structure for a single sensor which is located inside a RoomModel
type Sensor struct {
	ID              uint     `json:"id"`
	RoomModelID     uint     `json:"room_model_id"`
	LatestData      Data     `json:"latest_data" gorm:"-"`
	Data            []Data   `json:"-"`
	MeshID          string   `json:"mesh_id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	MeasurementUnit string   `json:"measurement_unit"`
	UpperBound      *float64 `json:"upper_bound"`
	LowerBound      *float64 `json:"lower_bound"`
	GradientBound   *float64 `json:"gradient_bound"`
}

//Data specifies the structure for a single measured value w/ timestamp which was recorded by a sensor
type Data struct {
	ID       uint      `json:"id"`
	SensorID uint      `json:"sensor_id"`
	Value    float64   `json:"value"`
	Gradient float64   `json:"gradient"`
	Date     time.Time `json:"date"`
}

//Layout default mnemonic time
// https://pauladamsmith.com/blog/2011/05/go_time.html
const Layout = "2006-01-02 15:04:05"

var DB *gorm.DB

//SetupDatabase initializes the database w/ the orm mapping and postgres as the dialect;
//drop defines whether or not the currently active scheme should be dropped
func SetupDatabase(drop bool) {
	// TODO check if ssl mode can be enabled later
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))

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
		Name: "Overhead MEP Installation",
		Description: "This model shows a overhead MEP installation. To make things look better this model" +
			" has a longer description. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam" +
			" nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.",

		Url:      "files/overhead-mep-installation/model.zip",
		ImageUrl: "files/overhead-mep-installation/thumbnail.png",
	}

	m4 := &RoomModel{
		Name: "PGN Model",
		Description: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam" +
			" nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.",

		Url:      "files/pgn-model/model.zip",
		ImageUrl: "files/pgn-model/thumbnail.png",
	}

	models = append(models, m1, m2, m3, m4)
	meshIds := [][]string{
		{"vent1", "valve7", "valve8", "valve5", "tank1", "pump1", "pump2", "", ""}, // Facility Mechanical Room
		{"node13", "node11", "node14", "node8", "", "", "", "", ""},                // MEP Building Model
		{"1883", "1887", "9673", "10147", "", "", "", "", ""},                      // Overhead MEP Installation
		{"30058", "29608", "29632", "29614", "29660", "5417", "5411", "", ""},      // PGN Model
	}

	fmt.Printf("[i] loading sensor data %s\n", time.Now().String())

	sampleData := loadSampleData(fmt.Sprintf("%s/sensors/sensor_data.csv", sampleDataPath), dataLimit)

	for i, m := range models {

		// deep copy sample data
		d := make([][]Data, len(sampleData))
		for i := range sampleData {
			d[i] = make([]Data, len(sampleData[i]))
			copy(d[i], sampleData[i])
		}

		DB.Create(&m)

		s1 := Sensor{
			RoomModelID: m.ID,
			Name:        "Outdoor Temperature Sensor",
			Description: "This Outdoor Temperature Sensor is a narrow-band, long range, low power" +
				" consumption, high performance and high quality wireless sensor transmitting" +
				" temperature from a NTC probe.",
			MeshID:          meshIds[i][0],
			MeasurementUnit: "°C",
			Data:            d[0],
		}
		DB.Create(&s1)

		s2 := Sensor{
			RoomModelID: m.ID,
			Name:        "Thermal Flow Sensor",
			Description: "Inline Flow-Through Temperature Sensor monitors the temperature of a fluid" +
				" that passes through it where a system control module receives this temperature" +
				" reading and uses a control loop to control the overall system temperature.",

			MeshID:          meshIds[i][1],
			MeasurementUnit: "°C",
			Data:            d[1],
		}
		DB.Create(&s2)

		s3 := Sensor{
			RoomModelID: m.ID,
			Name:        "Thermal Return Flow Sensor",
			Description: "Inline Return Flow Temperature Sensor monitors the temperature of a fluid" +
				" that passes through it. The output values are related to the Thermal Flow Sensor.",

			MeshID:          meshIds[i][2],
			MeasurementUnit: "°C",
			Data:            d[2],
		}
		DB.Create(&s3)

		s4 := Sensor{
			RoomModelID: m.ID,
			Name:        "Room Temperature Sensor",
			Description: "Calibratable room temperature measuring transducer with Modbus connection," +
				" in an impact-resistant plastic housing.",
			MeshID:          meshIds[i][3],
			MeasurementUnit: "°C",
			Data:            d[3],
		}
		DB.Create(&s4)

		s5 := Sensor{
			RoomModelID: m.ID,
			Name:        "Heating Room Temperature Sensor",
			Description: "Calibratable room temperature measuring transducer with Modbus connection," +
				" in an impact-resistant plastic housing. Designed for higher temperatures.",
			MeshID:          meshIds[i][4],
			MeasurementUnit: "°C",
			Data:            d[4],
		}
		DB.Create(&s5)

		s6 := Sensor{
			RoomModelID:     m.ID,
			Name:            "Pressure Sensor",
			Description:     "Pressure Sensor for water pipe pressure measurement at water distribution utilities.",
			MeshID:          meshIds[i][5],
			MeasurementUnit: "bar",
			Data:            d[5],
		}
		DB.Create(&s6)

		s7 := Sensor{
			RoomModelID: m.ID,
			Name:        "Tap Water Temperature Sensor Heating Room",
			Description: "This Tap Water Temperature Sensor is a probe that measures water temperature from" +
				" -40° to +70°C. It consists of a thermistor encased in a sheath made from grade 316L" +
				" stainless steel.",
			MeshID:          meshIds[i][6],
			MeasurementUnit: "°C",
			Data:            d[6],
		}
		DB.Create(&s7)

		s8 := Sensor{
			RoomModelID: m.ID,
			Name:        "Tap Water Temperature Sensor Inflow",
			Description: "This Tap Water Temperature Sensor is a probe that measures water temperature from" +
				" -40° to +70°C. It consists of a thermistor encased in a sheath made from grade 316L" +
				" stainless steel.",
			MeshID:          meshIds[i][7],
			MeasurementUnit: "°C",
			Data:            d[7],
		}
		DB.Create(&s8)
	}

	fmt.Printf("[✓] finished loading sensor data %s\n", time.Now().String())
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

func loadSampleData(path string, dataLimit int) [][]Data {
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Println("[!] Error ", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	data := make([][]Data, 8)

	i := 0

	for {
		if dataLimit != -1 && i == dataLimit + 1 {
			break
		}

		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			panic("[!] Error loading sample data")
		}

		if len(line) != 9 {
			fmt.Println("[!] Line skipped because there is not enough data prepared.", err)
			fmt.Printf("[!] Line: %s", line)
			continue
		}

		if i == 0 {
			i++
			continue
		}

		t, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			fmt.Println(err)
			panic("[!] Error parsing date")
		}

		dt := time.Unix(t, 0).UTC()

		for j := 0; j < len(data); j++ {
			v, err := strconv.ParseFloat(line[j+1], 64)
			if err != nil {
				fmt.Println(err)
				panic("[!] Error parsing value")
			}

			d := Data{
				Date:     dt,
				Value:    v,
				Gradient: 0.0,
			}

			if len(data[j]) > 0 {
				d.Gradient, _, _ = calculateGradient(data[j][len(data[j])-1], d)
			} else {
				data[j] = make([]Data, 0)
			}

			data[j] = append(data[j], d)
		}

		i++
	}

	for _, d := range data {
		sort.Slice(d, func(i, j int) bool {
			return d[i].Date.Before(d[j].Date)
		})
	}

	return data
}

func calculateGradient(d1 Data, d2 Data) (float64, float64, int64) {
	d := d2.Value - d1.Value
	dTime := d2.Date.Unix() - d1.Date.Unix()
	grad := d / float64(dTime)
	grad = math.Round(grad*100000) / 100000
	return grad, d, dTime
}
