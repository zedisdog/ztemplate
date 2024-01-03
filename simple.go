package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"simple/internal/api/handler"
	"simple/internal/config"
	exampleServer "simple/internal/server/exampleservice"
	"simple/internal/svc"
	"simple/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var configFile = flag.String("f", "etc/simple.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterExampleServiceServer(grpcServer, exampleServer.NewExampleServiceServer(ctx))

		reflection.Register(grpcServer)
	})

	gg, err := grpcHandler(
		context.TODO(),
		"127.0.0.1:8080",
		[]func(http.HandlerFunc) http.HandlerFunc{}, //middleware
		pb.RegisterExampleServiceHandlerFromEndpoint,
	)
	if err != nil {
		panic(err)
	}
	server := rest.MustNewServer(
		c.RestConf,
		rest.WithNotFoundHandler(gg),
	)
	handler.RegisterHandlers(server, ctx)

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(server)
	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway at %s:%d...\n", c.RestConf.Host, c.RestConf.Port)
	group.Start()
}

func withMiddlewares(handler http.Handler, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}
	var hf http.HandlerFunc
	hf = func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	}
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}

func grpcHandler(
	ctx context.Context,
	grpcEndpoint string,
	middlewares []func(http.HandlerFunc) http.HandlerFunc,
	registerFuncs ...func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error),
) (http.Handler, error) {
	mux := runtime.NewServeMux(
		// runtime.WithXxx方法可对请求的解析做一些设置
		// runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		// runtime.WithErrorHandler(func(ctx context.Context, serveMux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		// 	if err == nil {
		// 		return
		// 	}
		// 	w.WriteHeader(500)
		// }),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, reg := range registerFuncs {
		err := reg(ctx, mux, grpcEndpoint, opts)
		if err != nil {
			return nil, err
		}
	}
	return withMiddlewares(mux, middlewares...), nil
}
