package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

const path = "/Users/deklerk/workspace/lsprobe"

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err == nil {
		fmt.Println("+++++++++++++ Success")
	} else {
		fmt.Println("+++++++++++++ Fail:", err)
	}

	cmd := exec.Command("bash", "-c", "kill -9 `lsof -i:8081 | grep gopls | awk '{print $2}'`")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("+++++++++++++ Failed to clean up:", err)
	}
}

func run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "gopls", "serve", "-listen=localhost:8081")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go func() {
		fmt.Println("+++++++++++++ Starting server")
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("+++++++++++++ Waiting for server initialisation")
	time.Sleep(1 * time.Second) // TODO(deklerk): Find a better way to wait for server init.

	c := &client{}
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		return err
	}
	defer conn.Close()
	jsonrpc2Conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), c)

	if err := initConn(ctx, jsonrpc2Conn); err != nil {
		return fmt.Errorf("failed to init: %w", err)
	}

	if err := tidy(ctx, jsonrpc2Conn); err != nil {
		return fmt.Errorf("failed to do thing: %w", err)
	}

	fmt.Println("+++++++++++++ Sleeping for an hour")
	time.Sleep(time.Hour)
	return nil
}

func initConn(ctx context.Context, conn *jsonrpc2.Conn) error {
	var initializeRes interface{}
	initializesParams := InitializeParams{
		RootPath: path,
		RootURI:  pathToURI(path),
		WorkspaceFolders: []WorkspaceFolder{
			{
				URI:  "file://" + path,
				Name: "dummy_workspace",
			},
		},
	}
	if err := conn.Call(ctx, "initialize", &initializesParams, &initializeRes); err != nil {
		return fmt.Errorf("error with initialize call: %w", err)
	}
	b, err := json.Marshal(initializeRes)
	if err != nil {
		return fmt.Errorf("error unmarshaling initialize call: %w", err)
	}
	fmt.Println("+++++++++++++ initialize result:", string(b))

	type initializedParams struct{}
	var initializedRes interface{}
	err = conn.Call(context.Background(), "initialized", &initializedParams{}, &initializedRes)
	if err != nil {
		return fmt.Errorf("error with initialized call: %w", err)
	}
	b, err = json.Marshal(initializedRes)
	if err != nil {
		return fmt.Errorf("error unmarshaling initialized call: %w", err)
	}
	fmt.Println("+++++++++++++ initialized result:", string(b))

	type handshakeParams struct{}
	var handshakeRes interface{}
	err = conn.Call(context.Background(), "gopls/handshake", &handshakeParams{}, &handshakeRes)
	if err != nil {
		return fmt.Errorf("error with gopls/handshake call: %w", err)
	}
	b, err = json.Marshal(handshakeRes)
	if err != nil {
		return fmt.Errorf("error unmarshaling gopls/handshake call: %w", err)
	}
	fmt.Println("+++++++++++++ gopls/handshake result:", string(b))

	return nil
}

func tidy(ctx context.Context, conn *jsonrpc2.Conn) error {
	var res interface{}
	type ThingParams struct {
		URIs []string `json:"URIs"`
	}
	if err := conn.Call(ctx, "gopls.tidy", &ThingParams{URIs: []string{path + "main.go"}}, &res); err != nil {
		return err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Println("+++++++++++++ Tidy result:", string(b))
	return nil
}

// pathToURI converts given absolute path to file URI
func pathToURI(path string) string {
	path = filepath.ToSlash(path)
	parts := strings.SplitN(path, "/", 2)

	// If the first segment is a Windows drive letter, prefix with a slash and skip encoding
	head := parts[0]
	if head != "" {
		head = "/" + head
	}

	rest := ""
	if len(parts) > 1 {
		rest = "/" + parts[1]
	}

	return "file://" + head + rest
}

type client struct{}

func (c *client) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	fmt.Printf("Handle! method: %s. remainder: %+v\n", req.Method, req)
}
