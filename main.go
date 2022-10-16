package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"opera/util"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	// 암호 생성시만 사용
	if len(os.Args) > 1 {
		args := os.Args[1:]
		if args[0] == "--encbibop" {
			fmt.Println(args[1])
			fmt.Println(util.Encrypt(args[1]))
			return
		}
	}

	// goroutine 중지용.
	ctx, cancelFunc := context.WithCancel(context.Background())
	// goroutine 이 모두 중지될때까지 기다리는 용도
	var wg sync.WaitGroup
	// ctrl-c, process kill
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 환경변수를 3초에 한번씩 다시 로드.
	go util.LoadEnv(ctx, &wg)
	err := godotenv.Overload()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logLevel := util.GetLogLevel()
	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	fmt.Println("db bibop is ->", util.Conf("dbbibop"))
	fmt.Println("db bibop is ->", util.Conf("dbbibop"))
	go func() {
		wg.Add(1)
		defer wg.Done()
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				log.Info("helo info")
				log.Warn("this is warn")
				log.Debug("this is debug")
			case <-ctx.Done():
				ticker.Stop()
				fmt.Println("goroutine1 exit")
				return
			}
		}
	}()

	fmt.Println("awaiting signal")
	// wg.Wait()
	<-sigs
	cancelFunc()
	wg.Wait()
	fmt.Println("exiting")
	// time.Sleep(4 * time.Second)

}
