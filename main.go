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

	"github.com/labstack/echo"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kafka-reconsign/consumer"
	scramkafka "kafka-reconsign/internal/screamkafka"
	"kafka-reconsign/notification"
	"kafka-reconsign/repositories"
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

func initDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=%v",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
		viper.GetBool("db.parseTime"))

	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
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

	//consumer 2
	consumerGroupClient2, err := scramkafka.NewConsumerClient(viper.GetBool("kafka.auth"))
	if err != nil {
		log.Fatalf("cannot new sarama consumer client:%s", err)
	}
	defer consumerGroupClient2.Close()
	//consumer 2

	db := initDatabase()
	reconcileRepository := repositories.NewReconcileRepositoryDB(db)

	consumerHandler := consumer.NewConsumer(&wg, pool, consumer.NewServicePayment(reconcileRepository))
	consumer2Handler := consumer.NewConsumer2(&wg, pool, consumer.NewServiceInsurance(reconcileRepository))
	reconcileJobHandler := notification.NewReconcileJob(reconcileRepository)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "success",
		})
	})

	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroupClient.Consume(ctx, viper.GetStringSlice("kafka.topic.consumeTopic.Next"), consumerHandler); err != nil {
				log.Printf("Consume error[%s]: %v", viper.GetStringSlice("kafka.topic.consumeTopic.Next"), err)
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
	//consumer 2
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroupClient2.Consume(ctx, viper.GetStringSlice("kafka.topic.consumeTopic.insourer"), consumer2Handler); err != nil {
				log.Printf("Consume error[%s]: %v", viper.GetStringSlice("kafka.topic.consumeTopic.insourer"), err)
			}
			if ctx.Err() != nil {
				log.Printf("Context Error %v", ctx.Err())
				return
			}

			consumer2Handler.Ready = make(chan struct{})
		}
	}()
	<-consumer2Handler.Ready
	log.Println("Consumer2 Up and Running!")
	// consumer 2
	go func() {
		if err := reconcileJobHandler.CheckReconcileStatus(); err != nil {
			log.Println(err)
			//return
		}
	}()
	go func() {
		if err := reconcileJobHandler.CheckAlertStatus(); err != nil {
			log.Println(err)
			//return
		}
	}()

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
