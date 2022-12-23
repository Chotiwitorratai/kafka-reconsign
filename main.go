package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"

	"kafka-reconsign/consumer"
	scramkafka "kafka-reconsign/internal/screamkafka"
	"kafka-reconsign/repositories"
)

const (
	dbServer = "server_name"
	dbName   = "database_name"
	dbUser   = "user_name"
	dbPass   = "password"
)

// READ CONFIG
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("fatal error config file: %s", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	runtime.GOMAXPROCS(1)
	log.Printf("%v : %v", viper.GetString("log.level"), viper.GetString("log.env"))
}

func InitDatabase() *sqlx.DB {
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", dbServer, dbName, dbUser, dbPass)
	db, err := sqlx.Connect("mssql", connString)
	if err != nil {
		return nil
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db
}

func main() {
	ctx, cancelConsumer := context.WithCancel(context.Background())
	defer cancelConsumer()

	var (
		e       = initEcho()
		pool, _ = ants.NewPool(viper.GetInt("app.pool"))
		wg      sync.WaitGroup
	)

	go run(e)

	consumerGroupClient, err := scramkafka.NewConsumerClient(viper.GetBool("kafka.auth"))
	if err != nil {
		log.Fatalf("cannot new sarama consumer client:%s", err)
	}
	defer consumerGroupClient.Close()

	db := InitDatabase()
	orderRepository := repositories.NewOrderRepositoryDB(db)
	consumerHandler := consumer.NewConsumer(&wg, pool, consumer.New(orderRepository))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "success",
		})
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroupClient.Consume(ctx, viper.GetStringSlice("kafka.topic.consumeTopic"), consumerHandler); err != nil {
				log.Printf("Consume error[%s]: %v", viper.GetStringSlice("kafka.topic.consumeTopic"), err)
			}
			if ctx.Err() != nil {
				log.Printf("Context Error %v", ctx.Err())
				return
			}

			consumerHandler.Ready = make(chan struct{})
		}
	}()
	<-consumerHandler.Ready
	log.Println("Consumer Up and Running!")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-sigterm:
	}

	cancelConsumer()
	wg.Wait()

}

func initEcho() *echo.Echo {
	e := echo.New()
	return e
}

func run(router *echo.Echo) {
	log.Printf("Starting %s", viper.GetString("app.name"))
	log.Printf("Health probe serve at port %s", viper.GetString("app.port"))
	log.Println(router.Start(":" + viper.GetString("app.port")))
	// logx.Infof("Starting %s", viper.GetString("app.name"))
	// logx.Infof("Health probe serve at port %s", viper.GetString("app.port"))
	// logx.Info(router.Start(":" + viper.GetString("app.port")))
}
