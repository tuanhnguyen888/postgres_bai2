package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tuanhnguyen888/postgres_Bai2/models"
	"github.com/tuanhnguyen888/postgres_Bai2/storage"
	"gorm.io/gorm"
)

type Alert struct {
	AlertId     *string   `json:"alertID"`
	Category    *string   `json:"category"`
	CloseBy     *string   `json:"closeBy"`
	Create      time.Time `json:"create"`
	Datetime    time.Time `json:"dateTime"`
	Description *string   `json:"description"`
	LastUpdate  time.Time `json:"lastUpdate"`
	Message     *string   `json:"message"`
	Object      *string   `json:"object"`
	ObjectType  *string   `json:"objectType"`
	Owner       *string   `json:"owner"`
	Status      *string   `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Type        *string   `json:"type"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateAlert(context *fiber.Ctx) error {
	alert := Alert{}
	err := context.BodyParser(&alert)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "not request"})
		return err
	}

	err = r.DB.Create(&alert).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create alert"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "alert has been added",
		})
	return nil

}

func (r *Repository) GetAlert(context *fiber.Ctx) error {
	q := `
		SELECT 
			concat( date_part('year', timestamp) ,'/', date_part('month', timestamp)) AS _interval,
			type,
			COUNT(type)
		FROM alerts 
		GROUP by timestamp, type
		ORDER BY timestamp, count DESC
	`
	rows, err := r.DB.Raw(q).Rows()

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could get alert"})
		return err
	}

	defer rows.Close()

	type AlertCountbyMonth struct {
		Interval *string `json:"interval"`
		Type     *string `json:"type"`
		Count    int     `json:"count"`
	}

	rowDatas := []AlertCountbyMonth{}
	for rows.Next() {
		var Interval *string
		var Type *string
		var Count int

		rows.Scan(&Interval, &Type, &Count)
		rowData := AlertCountbyMonth{Interval, Type, Count}
		rowDatas = append(rowDatas, rowData)
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Group fetched",
			"data":    rowDatas,
		})
	return nil
}

func (r *Repository) DeleteAlert(context *fiber.Ctx) error {
	groupAlert := models.Alert{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "ID EMPTY"})
		return nil
	}

	err := r.DB.Delete(groupAlert, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete alert"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "deleted"})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	app.Post("/create_alert", r.CreateAlert)

	app.Get("/Alerts", r.GetAlert)

	app.Delete("/DeleteAlert/:id", r.DeleteAlert)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSQLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewInit(config)
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	err = models.MigrateAlert(db)
	if err != nil {
		panic(err)
	}
	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":1012")

}
