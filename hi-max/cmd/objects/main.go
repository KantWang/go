package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	objectsapi "objects"
	"objects/dbcon"
	objects "objects/gen/objects"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	// fmt.Println(*hostF, *domainF, *httpPortF, *secureF, *dbgF)
	// flag.Parse()
	// fmt.Println(*hostF, *domainF, *httpPortF, *secureF, *dbgF)
	// fmt.Println(log.Ldate, log.Ltime, log.Lmicroseconds, log.Llongfile, log.Lshortfile, log.LUTC, log.Lmsgprefix, log.LstdFlags)
	// log.SetFlags(log.Lshortfile)
	// log.Println("logging")

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[objectsapi] ", log.Ltime)
	}

	// Setup db.
	var (
		db *sql.DB
	)
	{
		db = dbcon.GetDB()
		defer db.Close() // main 루틴이 바로 종료되지 않고 서버 꺼질 때까지 Wait 걸려있어서 여기서 defer 해줘도 문제 없다
	}

	// Initialize the services.
	var (
		objectsSvc objects.Service
	)
	{
		objectsSvc = objectsapi.NewObjects(logger, db)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		objectsEndpoints *objects.Endpoints
	)
	{
		objectsEndpoints = objects.NewEndpoints(objectsSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) // 시스템콜 SIGINT, SIGTERM 발생 시 chan c로 relay
		// fmt.Println("errc!")
		tmp := fmt.Errorf("%s", <-c) // 대기
		errc <- tmp                  // 보냄
		// fmt.Println("errc: ", tmp.Error()) >> 인터럽트
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://localhost:8080"
			u, err := url.Parse(addr)
			if err != nil {
				logger.Fatalf("invalid URL %#v: %s\n", addr, err)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					logger.Fatalf("invalid URL %#v: %s\n", u.Host, err)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, objectsEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		logger.Fatalf("invalid host argument: %q (valid hosts: localhost)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc) // 대기 중. 에러(인터럽트 or ...)가 발생하면 수신되어 다음으로 진행

	// Send cancellation signal to the goroutines.
	cancel() // http.go의 <-ctx.Done() 진행시킴

	wg.Wait()
	logger.Println("exited")
}
