package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"projectlayout/app/cmd/album/internal/api"
	"projectlayout/app/cmd/album/internal/data"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}

type handerFunc func(http.ResponseWriter, *http.Request)

func (f handerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func ServerListen(ctx context.Context, server *http.Server) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	errChan := make(chan error, 1)

	go func() {
		url := server.Addr
		if err := server.ListenAndServe(); err != nil {
			errChan <- errors.Wrapf(err, "httpServer,url:%s,%s", url, err)
		} else {
			errChan <- nil
		}

	}()

	select {
	case <-ctx.Done():
		server.Shutdown(context.Background())
		return ctx.Err()
	case e := <-errChan:
		return e
	}

}

func kill(signal string, pid string) error {
	cmd := exec.Command("kill", signal, pid) // kill -SIGINT 144234
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "kill %s %s, %s", signal, pid, err)
	}
	return nil
}

func signalListen(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	pid := os.Getpid()
	fmt.Errorf("pid", pid)

	breakChan := make(chan struct{}, 1)
	go func() {
		interrupt := make(chan os.Signal, 1)
		reload := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT) // SIGINT 中断信息，相当于 Ctrl+c
		signal.Notify(reload, syscall.SIGHUP)
	OUTOUT:
		for {
			select {
			case <-interrupt:
				breakChan <- struct{}{}
				break OUTOUT
			case <-reload:
			default:
			}
		}
	}()

	select {
	case <-ctx.Done():
		if err := kill("-SIGINT", strconv.Itoa(pid)); err != nil { // 这里主要目的是为了 中断上面 goroutine 里信号的监听
			log.Print(err)
		}
		return ctx.Err()
	case <-breakChan:
		return errors.New("signalListen catch SIGINT")
	}
}

func main() {

	_ = InitSetting()

	defer data.Close()

	root, _ := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(root)

	g.Go(func() error {
		var server http.Server
		mux := http.NewServeMux()
		mux.Handle("/album", handerFunc(api.Album))
		mux.Handle("/shutdown", handerFunc(func(w http.ResponseWriter, r *http.Request) {
			server.Shutdown(context.Background())
		}))

		server = http.Server{Addr: viper.GetString("port"), Handler: mux}
		return ServerListen(ctx, &server)
	})

	g.Go(func() error {
		return signalListen(ctx)
	})

	fmt.Println("服务已启动")
	if err := g.Wait(); err != nil {

		fmt.Printf("%+v", err)
	}

}
