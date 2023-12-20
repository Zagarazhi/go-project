package main

// Импорты
import (
	"fmt"
	"log"
	"net"

	apiservice "github.com/Zagarazhi/go-project/api_service"
	"github.com/Zagarazhi/go-project/generated"
	userservice "github.com/Zagarazhi/go-project/user_service"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Kline структура представляет модель для таблицы klines.
type Kline struct {
	gorm.Model               `gorm:"-"`
	Symbol                   string `gorm:"column:symbol"`
	Interval                 string `gorm:"column:interval"`
	OpenTime                 int64  `gorm:"column:open_time"`
	OpenPrice                string `gorm:"column:open_price"`
	HighPrice                string `gorm:"column:high_price"`
	LowPrice                 string `gorm:"column:low_price"`
	ClosePrice               string `gorm:"column:close_price"`
	Volume                   string `gorm:"column:volume"`
	CloseTime                int64  `gorm:"column:close_time"`
	QuoteAssetVolume         string `gorm:"column:quote_asset_volume"`
	NumberOfTrades           int    `gorm:"column:number_of_trades"`
	TakerBuyBaseAssetVolume  string `gorm:"column:taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume string `gorm:"column:taker_buy_quote_asset_volume"`
}

func main() {
	// Инициализация viper
	viper.SetConfigName("config") // имя файла конфигурации без расширения
	viper.SetConfigType("yaml")   // или "json"
	viper.AddConfigPath(".")      // путь к директории с файлом конфигурации

	// Чтение файла конфигурации
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		log.Fatalf("Failed to read config file: %v", viperErr)
	}

	// Переменные из конфигурации
	serverAdress := viper.GetString("server.adress")
	serverPort := viper.GetString("server.port")
	grpcAdress := viper.GetString("api.grpcAdress")
	grpcPort := viper.GetString("api.grpcPort")

	// Запуск gRPC-сервера для ApiService
	runApiServiceServer(grpcAdress, grpcPort)

	// Создаем клиент gRPC для ApiService
	apiServiceConn, err := grpc.Dial(fmt.Sprintf("%s:%s", grpcAdress, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ApiService: %v", err)
	}
	defer apiServiceConn.Close()

	apiServiceClient := generated.NewApiServiceClient(apiServiceConn)

	// Инициализация БД
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создаем gRPC-сервер для UserService
	s := grpc.NewServer()
	userServiceServer := &userservice.UserServiceServer{
		ApiServiceClient: apiServiceClient,
		DB:               db,
	}
	generated.RegisterUserServiceServer(s, userServiceServer)

	// Запускаем HTTP-сервер
	fmt.Printf("HTTP Server is running on %s\n", fmt.Sprintf("%s:%s", serverAdress, serverPort))
	userservice.StartHTTPServer(userServiceServer, fmt.Sprintf("%s:%s", serverAdress, serverPort))
}

// initViper инициализирует viper для работы с конфигурацией.
func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		log.Fatalf("Failed to read config file: %v", viperErr)
	}
}

// InitDB инициализирует и возвращает соединение с базой данных.
func InitDB() (*gorm.DB, error) {
	dsn := viper.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if viper.GetBool("database.shouldMigrate") {
		Migrate(db)
	}
	return db, nil
}

// runApiServiceServer создает и запускает gRPC-сервер для ApiService.
func runApiServiceServer(grpcAdres string, grpcPort string) {
	apiServiceServer := &apiservice.ApiServiceServer{}
	grpcServer := grpc.NewServer()
	generated.RegisterApiServiceServer(grpcServer, apiServiceServer)

	apiListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcAdres, grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", fmt.Sprintf("%s:%s", grpcAdres, grpcPort), err)
	}
	fmt.Printf("gRPC Server is running on %s\n", fmt.Sprintf("%s:%s", grpcAdres, grpcPort))
	go func() {
		if err := grpcServer.Serve(apiListener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
}

// Создает таблицу klines и индекс idx_symbol_interval.
func Migrate(db *gorm.DB) error {
	// Миграция для создания таблицы klines.
	if err := db.AutoMigrate(&Kline{}); err != nil {
		return err
	}

	// Миграция для создания индекса idx_symbol_interval.
	if err := db.Exec("CREATE INDEX idx_symbol_interval ON klines (symbol, interval)").Error; err != nil {
		return err
	}

	fmt.Println("Migration successful.")
	return nil
}
