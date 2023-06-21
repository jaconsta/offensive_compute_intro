package main

/**
 * RAT -> Remote Administrtion Tool
 *
 * Install protobufs (`protoc`)
 * (Installer guide)[https://grpc.io/docs/protoc-installation/]
 * ```sh
 * $ brew install protobufs
 * $ protoc --version  # Expect `libprotoc 23.3` or greater
 * ```
 * (Install the Go's side compiler)[https://grpc.io/docs/languages/go/quickstart/]
 * ```sh
 * $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
 * $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
 * ```
 *
 * Generate / update the proto files
 * ```sh
 * $ cd c2grpcapi # before was: embed
 * $ ls # embed.proto
 * $ protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative embed.proto
 * $ ls # embed.pb.go      embed.proto      embed_grpc.pb.go
 * ```
 * Then download complimentary packages
 * ```sh
 * $ go mod tidy
 * ```
**/
import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"simple_proxy/remote_admin_tool/c2grpcapi"

	"google.golang.org/grpc"
)

type emberServer struct {
	work, output chan *c2grpcapi.Command
	c2grpcapi.UnsafeEmbedServer
}
type adminServer struct {
	work, output chan *c2grpcapi.Command
	c2grpcapi.UnimplementedAdminServer
}

func NewEmbededServer(work, output chan *c2grpcapi.Command) *emberServer {
	s := new(emberServer)
	s.work = work
	s.output = output
	return s
}
func NewAdminServer(work, output chan *c2grpcapi.Command) *adminServer {
	s := new(adminServer)
	s.work = work
	s.output = output
	return s
}

func (s *emberServer) GetCommand(ctx context.Context, empty *c2grpcapi.Empty) (*c2grpcapi.Command, error) {
	var cmd = new(c2grpcapi.Command)
	select {
	case cmd, ok := <-s.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New(("Channel error"))
	default:
		return cmd, nil
	}
}
func (s *emberServer) SendResult(ctx context.Context, result *c2grpcapi.Command) (*c2grpcapi.Empty, error) {
	s.output <- result
	return &c2grpcapi.Empty{}, nil
}

func (s *adminServer) GetCommand(ctx context.Context, cmd *c2grpcapi.Command) (*c2grpcapi.Command, error) {
	var res *c2grpcapi.Command
	go func() {
		s.work <- cmd
	}()
	res = <-s.output
	return res, nil
}

func main() {
	var (
		embedListener, adminListener net.Listener
		err                          error
		opts                         []grpc.ServerOption
		work, output                 chan *c2grpcapi.Command
	)
	work, output = make(chan *c2grpcapi.Command), make(chan *c2grpcapi.Command)
	embed := NewEmbededServer(work, output)
	admin := NewAdminServer(work, output)

	if embedListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 4445)); err != nil {
		log.Fatal(err)
	}
	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9995)); err != nil {
		log.Fatal(err)
	}
	grpcAdminServer, grpcEmbedServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	c2grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	c2grpcapi.RegisterEmbedServer(grpcEmbedServer, embed)
	fmt.Println("Started the c2 server for admin and embed clients")
	go func() {
		grpcEmbedServer.Serve(embedListener)
	}()
	grpcAdminServer.Serve(adminListener)
}
