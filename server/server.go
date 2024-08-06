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

// TODO build on interfaces
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

/*
0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8
0x1F98431c8aD98523631AE4a59f267346ea31F984
0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35
0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
0xCBCdF9626bC03E24f779434178A73a0B4bad62eD
0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599
*/
