package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Location struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//RoomModel specifies the structure for a single BIM model
type RoomModel struct {
	ID       uint     `json:"id"`
	Sensors  []Sensor `json:"sensors"`
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	ImageUrl string   `json:"image_url"`
	Type     string   `json:"type"`
	Location Location `json:"location" gorm:"embedded"`
	Floors   int      `json:"floors"`
}

//Sensor specifies the structure for a single sensor which is located inside a RoomModel
type Sensor struct {
	ID              uint     `json:"id"`
	RoomModelID     uint     `json:"room_model_id"`
	LatestData      Data     `json:"latest_data" gorm:"-"`
	Data            []Data   `json:"-"`
	ImportName      string   `json:"import_name,omitempty" gorm:"-"`
	MeshID          *int64   `json:"mesh_id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	MeasurementUnit string   `json:"measurement_unit"`
	Range           string   `json:"range"`
	UpperBound      *float64 `json:"upper_bound"`
	LowerBound      *float64 `json:"lower_bound"`
	GradientBound   *float64 `json:"gradient_bound"`
}

//Data specifies the structure for a single measured value w/ timestamp which was recorded by a sensor
type Data struct {
	ID       uint    `json:"id" csv:"-"`
	SensorID uint    `json:"sensor_id" csv:"-"`
	Value    float64 `json:"value" csv:"value"`
	Gradient float64 `json:"gradient" csv:"-"`
	Date     Date    `json:"date" csv:"date" gorm:"embedded"`
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
		DB.DropTableIfExists(&RoomModel{}, &Sensor{}, &Data{}, &Location{})
		fmt.Println("[✓] all data successfully dropped")
	}

	// Migrate the Schema
	DB.AutoMigrate(&RoomModel{}, &Sensor{}, &Data{}, &Location{})
	fmt.Println("[✓] schemes migrated")
}

func LoadModels(dataPath string, modelFolders []string, dataLimit int) {
	fmt.Printf("[i] loading sensor data %s\n", time.Now().String())

	var wg sync.WaitGroup
	wg.Add(len(modelFolders))

	for _, folder := range modelFolders {
		go func(folder string) {
			defer wg.Done()
			f, _ := ioutil.ReadFile(fmt.Sprintf("%s/sensors/%s/model.json", dataPath, folder))
			var m *RoomModel

			if err := json.Unmarshal(f, &m); err != nil {
				fmt.Println("error loading model: ", folder)
				panic(err)
			}

			for i := range m.Sensors {
				m.Sensors[i].Data = loadData(fmt.Sprintf("%s/sensors/%s/%s", dataPath, folder, m.Sensors[i].ImportName), dataLimit)
			}

			DB.Create(&m)
			fmt.Printf("[✓] model %s loaded %s\n", folder, time.Now().String())
		}(folder)
	}

	wg.Wait()
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

func loadData(path string, dataLimit int) []Data {
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
