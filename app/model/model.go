package model

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
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
	ID       uint      `json:"id" csv:"-"`
	SensorID uint      `json:"sensor_id" csv:"-"`
	Value    float64   `json:"value" csv:"value"`
	Gradient float64   `json:"gradient" csv:"-"`
	Date     Date      `json:"date" csv:"date" gorm:"embedded"`
}

type Date struct {
	time.Time `gorm:"column:date"`
}

// Convert the CSV string as internal date
func (date *Date) UnmarshalCSV(csv string) (err error) {
	t, err := strconv.ParseInt(csv, 10, 64)
	date.Time = time.Unix(t, 0).UTC()
	return err
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

	sampleData := [8][]Data{
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/outdoor_air_temperature.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/flow.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/return_flow.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/room_temperature.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/heating_target_temperature.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/pressure.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/tap_water_temperature_boiler_room.csv", sampleDataPath), dataLimit),
		loadSampleData(fmt.Sprintf("%s/sensors/sample_model/tap_water_temperature_inflow.csv", sampleDataPath), dataLimit),
	}

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

func loadSampleData(path string, dataLimit int) []Data {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[!] Error loading file")
		panic(err)
	}
	defer file.Close()

	data := make([]Data, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	headerText := scanner.Text()

	i := 0
	for scanner.Scan() {
		if i == dataLimit {
			break
		}

		var d []Data
		combined := fmt.Sprintf("%s\n%s", headerText, scanner.Text())
		if err := gocsv.UnmarshalString(combined, &d); err != nil {
			panic(err)
		}
		data = append(data, d[0])

		if i > 0 {
			data[i].Gradient, _, _ = calculateGradient(&data[i-1], &data[i])
		}

		i++
	}

	return data
}

func calculateGradient(d1 *Data, d2 *Data) (float64, float64, int64) {
	d := d2.Value - d1.Value
	dTime := d2.Date.Unix() - d1.Date.Unix()
	grad := d / float64(dTime)
	grad = math.Round(grad*100000) / 100000
	return grad, d, dTime
}
