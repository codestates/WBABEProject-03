package main

import (
	"context"
	"fmt"
	ctl "lecture/WBABEProject-03/controller"
	"lecture/WBABEProject-03/model"
	rt "lecture/WBABEProject-03/router"
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
	} else if controller, err := ctl.NewCTL(mod); err != nil {
		panic(fmt.Errorf("controller.New > %v", err))
	} else if rt, err := rt.NewRouter(controller); err != nil {
		panic(fmt.Errorf("router.NewRouter > %v", err))
	} else {
		mapi := &http.Server{
			Addr:           ":8080",
			Handler:        rt.Idx(),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		g.Go(func() error {
			return mapi.ListenAndServe()
		})

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := mapi.Shutdown(ctx); err != nil {
			fmt.Println("Server Shutdown:", err)
		}

		select {
		case <-ctx.Done():
			fmt.Println("timeout of 5 seconds.")
		}

		fmt.Println("Server exiting")
	}
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
