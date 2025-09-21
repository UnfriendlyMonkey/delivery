package main

import (
	"database/sql"
	"delivery/cmd"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/pkg/errs"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := getConfigs()

	connectionString, err := makeConnectionString(
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbSslMode,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	createDBIfNotExists(
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbSslMode,
	)

	gormDB := mustGormOpen(connectionString)
	mustAutoMigrate(gormDB)

	compositionRoot := cmd.NewCompositionRoot(
		config,
		gormDB,
	)
	defer compositionRoot.CloseAll()

	startWebServer(compositionRoot, config.HttpPort)
}

func getConfigs() cmd.Config {
	config := cmd.Config{
		HttpPort:                  goDotEnvVariable("HTTP_PORT"),
		DbHost:                    goDotEnvVariable("DB_HOST"),
		DbPort:                    goDotEnvVariable("DB_PORT"),
		DbUser:                    goDotEnvVariable("DB_USER"),
		DbPassword:                goDotEnvVariable("DB_PASSWORD"),
		DbName:                    goDotEnvVariable("DB_NAME"),
		DbSslMode:                 goDotEnvVariable("DB_SSLMODE"),
		GeoServiceGrpcHost:        goDotEnvVariable("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:                 goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:        goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaBasketConfirmedTopic: goDotEnvVariable("KAFKA_BASKET_CONFIRMED_TOPIC"),
		KafkaOrderChangedTopic:    goDotEnvVariable("KAFKA_ORDER_CHANGED_TOPIC"),
	}
	return config
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func makeConnectionString(host, port, user, password, dbName, sslMode string) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError("host")
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError("port")
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError("user")
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError("password")
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError("dbName")
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError("sslMode")
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbName,
		sslMode), nil
}

func createDBIfNotExists(host, port, user, password, dbName, sslMode string) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", sslMode)
	if err != nil {
		log.Fatalf("Error connecting to Database: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to Database: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing DB: %v", err)
		}
	}()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Printf("Error creating DB (possibly exists already): %v", err)
	}
}

func mustGormOpen(connectionString string) *gorm.DB {
	pgGorm, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN: connectionString,
			PreferSimpleProtocol: true,
		},
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to DB via Gorm: %s", err)
	}
	return pgGorm
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&orderrepo.OrderDTO{})
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
}

func startWebServer(_ *cmd.CompositionRoot, port string) {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
