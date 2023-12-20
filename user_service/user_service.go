package userservice

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Zagarazhi/go-project/generated"
	"github.com/valyala/fasthttp"
)

type UserServiceServer struct {
	generated.UnimplementedUserServiceServer
	ApiServiceClient generated.ApiServiceClient
}

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

	return apiServiceResponse, nil
}

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

func (s *UserServiceServer) handleProcessRequest(ctx *fasthttp.RequestCtx) {
	var userData generated.UserData

	if err := json.Unmarshal(ctx.PostBody(), &userData); err != nil {
		log.Printf("Error unmarshaling user data: %v", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	result, err := s.ProcessData(ctx, &userData)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
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
