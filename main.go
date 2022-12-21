package main

import (
	"context"
	"fmt"
	ctl "lecture/go-mvc/controller"
	"lecture/go-mvc/model"
	rt "lecture/go-mvc/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	//model 모듈 선언
	if mod, err := model.NewModel(); err != nil {
		panic(err)
		//~생략
	} else if controller, err := ctl.NewCTL(mod); err != nil { //controller 모듈 설정
		panic(err)
		//~생략
	} else if rt, err := rt.NewRouter(controller); err != nil { //router 모듈 설정
		panic(err)
		//~생략
	} else {
		mapi := &http.Server{
			Addr:           ":8080",
			Handler:        rt.Idx(),
			ReadTimeout:    3 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		g.Go(func() error {
			return mapi.ListenAndServe()
		})
		quit := make(chan os.Signal)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := mapi.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}

		select {
		case <-ctx.Done():
		}
		fmt.Println("sever exit")
	}
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
