package userservice

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Zagarazhi/go-project/generated"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

// Структура модели для таблицы klines
type Kline struct {
	gorm.Model               `gorm:"-"`
	Symbol                   string
	Interval                 string
	OpenTime                 int64
	OpenPrice                string
	HighPrice                string
	LowPrice                 string
	ClosePrice               string
	Volume                   string
	CloseTime                int64
	QuoteAssetVolume         string
	NumberOfTrades           int
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
}

type UserServiceServer struct {
	generated.UnimplementedUserServiceServer
	ApiServiceClient generated.ApiServiceClient
	DB               *gorm.DB
}

// Метод передачи данных от пользователя в Api Service
func (s *UserServiceServer) ProcessData(ctx context.Context, req *generated.UserData) (*generated.ApiServiceResponce, error) {
	apiServiceResponse, err := s.ApiServiceClient.CallApiService(ctx, &generated.ApiServiceRequest{
		Symbol:   req.Symbol,
		Interval: req.Interval,
		Limit:    1,
	})

	if err != nil {
		log.Printf("Error calling ApiService: %v", err)
		return nil, err
	}

	// Преобразование apiServiceResponse в структуру Kline
	kline := Kline{
		Symbol:                   apiServiceResponse.Symbol,
		Interval:                 apiServiceResponse.Interval.String(),
		OpenTime:                 apiServiceResponse.OpenTime,
		OpenPrice:                apiServiceResponse.OpenPrice,
		HighPrice:                apiServiceResponse.HighPrice,
		LowPrice:                 apiServiceResponse.LowPrice,
		ClosePrice:               apiServiceResponse.ClosePrice,
		Volume:                   apiServiceResponse.Volume,
		CloseTime:                apiServiceResponse.CloseTime,
		QuoteAssetVolume:         apiServiceResponse.QuoteAssetVolume,
		NumberOfTrades:           int(apiServiceResponse.NumberOfTrades),
		TakerBuyBaseAssetVolume:  apiServiceResponse.TakerBuyBaseAssetVolume,
		TakerBuyQuoteAssetVolume: apiServiceResponse.TakerBuyQuoteAssetVolume,
	}

	// Запрос с обновлением или созданием записи
	result := s.DB.Where(&Kline{Symbol: kline.Symbol, Interval: kline.Interval}).Assign(kline).FirstOrCreate(&kline)
	if result.Error != nil {
		log.Printf("Error creating or updating kline record: %v", result.Error)
		return nil, result.Error
	}

	return apiServiceResponse, nil
}

// Метод запуска HTTP-сервера
func StartHTTPServer(s *UserServiceServer, addr string) {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/process":
			s.handleProcessRequest(ctx)
		default:
			ctx.Error("Not found", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(addr, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

// Обработка запросов
func (s *UserServiceServer) handleProcessRequest(ctx *fasthttp.RequestCtx) {
	var userData generated.UserData

	if err := json.Unmarshal(ctx.PostBody(), &userData); err != nil {
		log.Printf("Error unmarshaling user data: %v", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	result, err := s.ProcessData(ctx, &userData)
	if err != nil {
		ctx.Error("BadRequest", fasthttp.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(responseJSON)
}
