package main

import (
	"context"
	"fmt"
	ctl "lecture/go-project/controller"
	model "lecture/go-project/model"
	rt "lecture/go-project/router"
	"log"
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
			Handler:        rt.Index(),
			ReadTimeout:    3 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		g.Go(func() error {
			if err := mapi.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(fmt.Errorf("listen: %s\n", err))
			}
			return nil
		})
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mapi.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		select {
		case <-ctx.Done():
			log.Println("timeout of 5 seconds.")
		}

		log.Println("Server exiting")
	}
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
