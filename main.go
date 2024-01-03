package main

import (
	"fmt"
	"net"
	"net/http"
	"io"
	"context"
	"errors"
	"github.com/gorilla/mux"
)

const keyServerAddr = "serverAddr"

func startEndpoint(port string, c chan int) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	h := mux.NewRouter()
	h.HandleFunc("/", getController).Methods("GET")
	h.HandleFunc("/", postController).Methods("POST")
	server := &http.Server{
		Addr: port,
		Handler: h,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	fmt.Printf("started %s\n", port)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("%s closed\n", port)
	} else if err != nil {
		fmt.Printf("error on %s\n", port)
	}
	cancelCtx()

	c <- 1
}

func getController(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Printf("%s: GET\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "testing")
}

func postController(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Printf("%s: POST\n", ctx.Value(keyServerAddr))
}

func main() {
	fmt.Println("start")

	c1 := make(chan int)
	go startEndpoint(":3333", c1)
	c2 := make(chan int)
	go startEndpoint(":4444", c2)

	<- c1
	<- c2

	fmt.Println("end")
}
