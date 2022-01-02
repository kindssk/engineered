package main

import (
	"context"
	"encoding/json"
	s "engineered/api"
	"engineered/configs/mysql"
	config_mysql_v1 "engineered/configs/mysql/config.mysql.v1"
	"engineered/internel/service/biz"
	"engineered/internel/service/data"
	ser "engineered/internel/service/service"
	"fmt"
	"github.com/ghodss/yaml"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"

	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ShutdownSignals = []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM,
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	group := errgroup.Group{}

	group.Go(func() error {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, ShutdownSignals...)
		select {
		case <-ctx.Done():
			fmt.Println("收到ctx cancel信号")
		case <-signals:
			fmt.Println("收到关闭服务信号")
		}
		return stop(ctx)
	})

	group.Go(func() error {
		err := start(ctx)
		if err != nil {
			cancel()
		}
		return err
	})
	err := group.Wait()
	if err != nil {
		log.Printf("sss %s\n", err)
	}
	fmt.Println("333")
	os.Exit(1)
}

func start(parent context.Context) error {
	//写log
	log.Printf("%s", "开始启动")
	var myError error
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	group := errgroup.Group{}
	group.Go(func() error {
		return initMysql(ctx)
	})
	//启动redis
	//group.Go(func() error {
	//
	//})

	//启动gRPC
	group.Go(func() error {
		return initGRPC(ctx)
	})

	err := group.Wait()
	if err != nil {
		cancel()
	}
	return myError
}

func stop(parent context.Context) error {
	stopSuccess := make(chan int, 1)
	ctx, cancel := context.WithTimeout(parent, time.Second*5)
	defer cancel()

	if mysql.MySqlClient() != nil {
		mysql.MySqlClient().Close()
	}
	stopSuccess <- 1
	select {
	case <-ctx.Done():
		log.Printf("%s\n", "收到超时信号")
	case <-stopSuccess:
		log.Printf("%s\n", "关闭成功")
	}
	return nil
}

func initMysql(ctx context.Context) error {
	log.Printf("%s\n", "初始化mysql")
	mc := new(config_mysql_v1.Config)
	applyYAML(mc)
	err := mysql.Dial(mc.Database, mc.Host, mc.Port, mc.Username, mc.Password,
		mysql.DialEncoding(mc.Encoding),
		mysql.DialTimeout(mc.Timeout),
	)
	if err != nil {
		return err
	}
	log.Printf("%s\n", "初始化mysql成功")
	return nil
}

func initGRPC(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":8099")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	s.RegisterServiceServer(server, ser.NewUserService(biz.NewUserBiz(data.NewDateUser(data.NewMysqlDB()))))
	//
	//
	//// Register reflection service on gRPC server.
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
		return err
	}
	log.Printf("%s\n", "gRPC启动成功")
	return nil
}

func applyYAML(m *config_mysql_v1.Config) {
	js, err := yaml.YAMLToJSON(loadMysqlConfig())
	if err != nil {
		panic(err)
	}
	applyJSON(m, js)
}

func applyJSON(m *config_mysql_v1.Config, js []byte) {
	err := json.Unmarshal(js, m)
	if err != nil {
		panic(err)
	}
}

func loadMysqlConfig() []byte {
	yml, err := ioutil.ReadFile("configs/mysql/mysql.yml")
	if err != nil {
		panic(err)
	}
	return yml
}
