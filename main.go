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

	db := initDatabase()
	reconcileRepository := repositories.NewReconcileRepositoryDB(db)

	consumerHandler := consumer.NewConsumer(&wg, pool, consumer.New(reconcileRepository))
	reconcileJobHandler := notification.NewReconcileJob(reconcileRepository)

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

	// test func SaveReconcile
	// var testData = repositories.Reconcile{
	// 		Model:                        gorm.Model{},
	// 		TransactionRefID:             "2testref2oo",
	// 		Status:                       "Success",
	// 		TransactionCreatedTimestamp:  time.Now().Add(-time.Second * 1),
	// 		PaymentInfoAmount:            10000,
	// 		PaymentInfoWebAdditionalInfo: "test",
	// 		PartnerInfoName:              "KIA",
	// 		PartnerInfoDeeplinkUrl:       "www.example.com",
	// 		PaymentPlatform:              "K-NEXT",
	// 		IdCard:                       "426xxxxxxxxxxxxxxxx",
	// 		PlanCode:                     "TA30",
	// 		PlanName:                     "TEST",
	// 		EffectiveDate:                time.Now(),
	// 		ExpireDate:                   time.Now().AddDate(1, 0, 0),
	// 		IssueDate:                    time.Now(),
	// 		InsuranceStatus:              "Success",
	// 		TotalSumInsured:              10000,
	// 		ProductOwner:                 "KIA",
	// 		PlanType:                     "TA",
	// }
	// err = reconcileRepository.SaveReconcile(testData)
	// if err != nil {
	// 	panic(err)
	// }

	// test func CheckNullReconcile true = มี  flase = ไม่มี
	// boo, err := reconcileRepository.CheckNullReconcile("lkjkllk")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(boo)

	// test func UpdateReconcile
	// var testData = repositories.Reconcile{
	// 	Model:                        gorm.Model{},
	// 	TransactionRefID:             "2testref2oo",
	// 	Status:                       "Successss",
	// 	TransactionCreatedTimestamp:  time.Now().Add(-time.Second * 1),
	// 	PaymentInfoAmount:            10000,
	// 	PaymentInfoWebAdditionalInfo: "test",
	// 	PartnerInfoName:              "KIA",
	// 	PartnerInfoDeeplinkUrl:       "www.example.com",
	// 	PaymentPlatform:              "K-NEXT",
	// 	IdCard:                       "426xxxxxxxxxxxxxxxx",
	// 	PlanCode:                     "TA30",
	// 	PlanName:                     "TEST",
	// 	EffectiveDate:                time.Now(),
	// 	ExpireDate:                   time.Now().AddDate(1, 0, 0),
	// 	IssueDate:                    time.Now(),
	// 	InsuranceStatus:              "Success",
	// 	TotalSumInsured:              10000,
	// 	ProductOwner:                 "KIA",
	// 	PlanType:                     "TA",
	// }
	// err = reconcileRepository.UpdateReconcile(testData)
	// if err != nil {
	// 	panic(err)
	// }

	// test func GetReconcileFail
	// recon, err := reconcileRepository.GetReconcileFail()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(recon)

	// test func SaveAlert
	// var testData = []repositories.Alert{
	// 	{
	// 		Model:     gorm.Model{},
	// 		Messages:  "success",
	// 		Count:     2,
	// 		Status:    "success",
	// 		NextAlert: time.Now(),
	// 		RefId:     "klhjllkj",
	// 		Missing:   "Next",
	// 	}, {
	// 		Model:     gorm.Model{},
	// 		Messages:  "success",
	// 		Count:     2,
	// 		Status:    "success",
	// 		NextAlert: time.Now(),
	// 		RefId:     "klhjllkj",
	// 		Missing:   "Next",
	// 	},
	// }
	// err = reconcileRepository.SaveAlert(testData)
	// if err != nil {
	// 	panic(err)
	// }

	//test func UpdateAlert
	// var testData = repositories.Alert{
	// 		Model:     gorm.Model{},
	// 		Messages:  "fail",
	// 		Count:     2,
	// 		Status:    "successss",
	// 		NextAlert: time.Now(),
	// 		RefId:     "klhjllkj",
	// 		Missing:   "Next",
	// 	}
	// err = reconcileRepository.UpdateAlert(testData)
	// if err != nil {
	// 	panic(err)
	// }

	//test func GetAlertFail
	// a, err := reconcileRepository.GetAlertFail()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print(a)

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
