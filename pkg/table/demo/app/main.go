package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hoorayui/core-framework/pkg/cap/database/mysql"
	"github.com/hoorayui/core-framework/pkg/table/demo/tables"
	"github.com/hoorayui/core-framework/pkg/table/doc"
	cap "github.com/hoorayui/core-framework/pkg/table/proto"
	"github.com/hoorayui/core-framework/pkg/table/registry"
	"github.com/hoorayui/core-framework/pkg/table/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	//"google.golang.org/grpc"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var (
	dbRead  *mysql.DB
	dbWrite *mysql.DB
)

func main() {
	// 0. Read configuration file
	{
		fmt.Println(os.Getwd())
		f := flag.String("f", "table/demo/app/etc/app.yml", "config file path")
		flag.Parse()
		log.Println("Configuration file:", *f)
		viper.SetConfigFile(*f)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
	}
	/******************************************************************************/
	// 1. 数据库连接
	{
		mysql.DebugLogOn = true
		db, err := mysql.NewDatabase(&mysql.ConnConfig{
			Driver:          "mysql",
			Host:            viper.GetString("mysql.host"),
			Port:            viper.GetString("mysql.port"),
			Database:        viper.GetString("mysql.database"),
			User:            viper.GetString("mysql.xuser"),
			Password:        viper.GetString("mysql.password"),
			MaxOpenConns:    200,
			MaxIdelConns:    15,
			ConnMaxLifeTime: 5 * time.Second,
			ConnMaxIdelTime: 5 * time.Second,
		})
		if err != nil {
			panic(err)
		}
		// 不用读写分离
		dbRead, dbWrite = db, db
	}
	/******************************************************************************/
	// 2. 注册表格数据
	{
		tables.Init()
	}
	// 启动表格文档服务
	{
		docPort := ":" + viper.GetString("app.docport")
		go doc.ServeHTTP(docPort, dbWrite, registry.GlobalTableRegistry())
		log.Println("Table doc server listen at", docPort)
	}
	/******************************************************************************/

	/******************************************************************************/
	// 4. 启动表格服务
	{
		gRPCPort := ":" + viper.GetString("app.grpcport")
		lis, err := net.Listen("tcp", gRPCPort)
		if err != nil {
			log.Fatalf("failed to listen at [%s]: %s", gRPCPort, err.Error())
		}

		server := grpc.NewServer()
		cap.RegisterTableWServiceServer(server,
			service.NewTableWService(dbWrite, dbRead, &TestUserInfoProvider{})) // 启动grpc服务

		go func() {
			log.Println("gRPC listen at", gRPCPort)
			if err := server.Serve(lis); err != nil {
				log.Fatal(err)
			}
		}()
	}
	/******************************************************************************/
	{
		// 启动gateway
		conn, err := grpc.DialContext(context.Background(),
			"0.0.0.0:8002",
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalln("Failed to dial server:", err)
		}

		gwmux := runtime.NewServeMux()
		_ = cap.RegisterTableWServiceHandler(context.Background(), gwmux, conn)
		gwServer := &http.Server{ // http服务端口
			Addr:    ":8090",
			Handler: gwmux,
		}
		//
		//8090端口提供gRPC-Gateway服务
		//log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
		log.Fatalln(gwServer.ListenAndServe())

		fmt.Println("*********")
	}
	// 5. hold
	{
		fmt.Println("Demo is running")
		ch := make(chan int, 1)
		<-ch
	}
}
