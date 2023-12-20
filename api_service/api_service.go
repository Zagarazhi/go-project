package apiservice

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Zagarazhi/go-project/generated"
	"github.com/valyala/fasthttp"
)

type ApiServiceServer struct {
	generated.UnimplementedApiServiceServer
}

// Обработка запроса от пользователя к API Binance
func (s *ApiServiceServer) CallApiService(ctx context.Context, req *generated.ApiServiceRequest) (*generated.ApiServiceResponce, error) {
	// Формируем URL для отправки HTTP-запроса
	interval, stringErr := convertIntervalFormat(req.Interval.String())
	if stringErr != nil {
		return nil, stringErr
	}
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", req.Symbol, interval, req.Limit)

	// Отправляем GET-запрос
	statusCode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}

	// Проверяем код статуса ответа
	if statusCode != fasthttp.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var data [][]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON: %v\n", err)
	}

	res := &generated.ApiServiceResponce{
		Symbol:     req.Symbol,
		Interval:   req.Interval,
		OpenTime:   int64(data[0][0].(float64)),
		OpenPrice:  data[0][1].(string),
		HighPrice:  data[0][2].(string),
		LowPrice:   data[0][3].(string),
		ClosePrice: data[0][4].(string),
		Volume:     data[0][5].(string),
		CloseTime:  int64(data[0][6].(float64)),
	}

	return res, nil
}

// Обработка строки
func convertIntervalFormat(input string) (string, error) {
	re := regexp.MustCompile(`^INTERVAL_(\d+)([SMHDW]|MONTH)$`)

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return "", fmt.Errorf("invalid input format")
	}

	quantity := matches[1]
	unit := strings.ToLower(matches[2])

	// Если единица измерения "month", заменяем на "M"
	if unit == "month" {
		unit = "M"
	}

	return quantity + unit, nil
}
