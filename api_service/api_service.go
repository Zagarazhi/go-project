package apiservice

import (
	"context"

	"github.com/Zagarazhi/go-project/generated"
)

type ApiServiceServer struct {
	generated.UnimplementedApiServiceServer
}

func (s *ApiServiceServer) CallApiService(ctx context.Context, req *generated.ApiServiceRequest) (*generated.ApiServiceResponce, error) {
	res := &generated.ApiServiceResponce{
		Symbol:     req.Symbol,
		Interval:   req.Interval,
		OpenTime:   123456789,
		OpenPrice:  "123.45",
		HighPrice:  "130.00",
		LowPrice:   "120.00",
		ClosePrice: "125.00",
		Volume:     "1000000",
		CloseTime:  987654321,
	}
	return res, nil
}
