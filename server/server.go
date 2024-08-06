// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package main

import (
	"flag"
	"fmt"
	"net/http"

	servers "dexprices.io/server/servers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/spf13/viper"
)

// TODO interfaces
//

func main() {
	var config string
	var server_ip, server_port string

	flag.StringVar(&config, "c", "/home/dexprices", "Enter your config")
	flag.Parse()
	viper.SetConfigType("yaml")
	viper.AddConfigPath(config)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	server_ip = viper.GetString("server.ip")
	server_port = viper.GetString("server.port")
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(servers.Dex), "")
	t := rpc.NewServer()
	t.RegisterCodec(json.NewCodec(), "application/json")
	t.RegisterService(new(servers.TONDex), "")
	r := mux.NewRouter()
	r.Handle("/dex", s)
	r.Handle("/ton", t)
	http.ListenAndServe(server_ip+":"+server_port, r)
}
