package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type ReadDeadliner interface {
	SetReadDeadline(t time.Time) error
}

func SetReadDeadlineOnCancel(ctx context.Context, d ReadDeadliner) {
	done := ctx.Done()
	go func() {
		select {
		case <-done:
			d.SetReadDeadline(time.Now())
		}
	}()
}

func handleConnection(ctx context.Context, client net.Conn) {
	remote, err := net.Dial("tcp", "1.1.1.1:443") // 替换为远端地址(vps ip)
	if err != nil {
		fmt.Println("Failed to connect to remote server:", err)
		return
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		SetReadDeadlineOnCancel(ctx, remote)

		_, err := io.Copy(remote, client)
		if err != nil {
			fmt.Println("Error copying data to remote server:", err)
		}
	}()

	_, err = io.Copy(client, remote)
	if err != nil {
		fmt.Println("Error copying data from remote server:", err)
	}

	time.AfterFunc(time.Second*3, func() {
		defer remote.Close()
		defer client.Close()
	})
}

func main() {
	listener, err := net.Listen("tcp", ":20801") // Replace with the desired local port 中转vps端口
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}

	fmt.Println("Proxy server is running on :20801")

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(context.Background(), client)
	}
}
