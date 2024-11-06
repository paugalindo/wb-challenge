package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	natsconn "github.com/nats-io/nats.go"

	"wb-challenge/internal/http"
	"wb-challenge/internal/nats"
	"wb-challenge/internal/postgres"
	"wb-challenge/internal/query"
)

func Run(ctx context.Context, cfg Config, logger *log.Logger) {
	ctx, cancelSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	dbSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser,
		cfg.DatabasePassword, cfg.DatabaseName,
	)

	sql, err := sql.Open("postgres", dbSource)
	if err != nil {
		logger.Fatal(err)
	}

	var natsConn *natsconn.Conn
	count := 0
Loop:
	for {
		natsConn, err = natsconn.Connect(cfg.NatsURL)
		switch {
		case err == nil:
			break Loop
		case count == 10:
			logger.Fatal(err)
		}
		time.Sleep(time.Second)
		count++
	}

	publisher := nats.NewPublisher(natsConn)

	vehiclesRepository := postgres.NewVehiclesRepository(sql)
	groupRepository := postgres.NewGroupsRepository(sql)
	cmdBus := InitCommandBus(&vehiclesRepository, &vehiclesRepository, &groupRepository, &groupRepository, &publisher)

	if err := InitNATSConsumer(natsConn, &cmdBus, logger); err != nil {
		logger.Fatal(err)
	}

	groupQS := query.NewGroupQS(&groupRepository)

	srv := http.New(cfg.ServicePort, &cmdBus, &groupQS)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal(err)
		}
	}()

	defer func() {
		cancelSignal()
		if err := natsConn.Drain(); err != nil {
			logger.Fatal(err)
		}
		if err := srv.Stop(ctx); err != nil {
			logger.Fatal(err)
		}
	}()

	<-ctx.Done()
}
