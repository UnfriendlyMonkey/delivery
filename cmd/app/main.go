package main

import (
	"database/sql"
	"delivery/cmd"
	httpadapter "delivery/internal/adapters/in/http"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
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

	startCronJobs(compositionRoot)
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
			DSN:                  connectionString,
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

func startWebServer(compositionRoot *cmd.CompositionRoot, port string) {
	handlers, err := httpadapter.NewServer(
		compositionRoot.NewCreateOrderHandler(),
		compositionRoot.NewCreateCourierHandler(),
		compositionRoot.NewGetAllCouriersHandler(),
		compositionRoot.NewGetIncompleteOrdersHandler(),
	)
	if err != nil {
		log.Fatalf("HTTP Server initialization error: %v", err)
	}
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	// Custom error handler to log errors
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Errorf("HTTP Error: %v", err)
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	e.Pre(middleware.RemoveTrailingSlash())

	// Register Swagger and health check
	registerSwaggerOpenAPI(e)
	registerSwaggerUI(e)
	registerHealthCheck(e)

	// Register API handlers
	servers.RegisterHandlers(e, handlers)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}

func registerSwaggerOpenAPI(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := servers.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
		}

		data, err := swagger.MarshalJSON()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
		}

		return c.Blob(http.StatusOK, "application/json", data)
	})
}

func registerSwaggerUI(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
		return c.HTML(http.StatusOK, html)
	})
}

func registerHealthCheck(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})
}

func startCronJobs(compositionRoot *cmd.CompositionRoot) {
	c := cron.New()
	_, err := c.AddJob("@every 10s", compositionRoot.NewAssignOrderJob())
	if err != nil {
		log.Fatalf("error adding cron job: %v", err)
	}
	_, err = c.AddJob("@every 1s", compositionRoot.NewMoveCouriersJob())
	if err != nil {
		log.Fatalf("error adding cron job: %v", err)
	}
}
