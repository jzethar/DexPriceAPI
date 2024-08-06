// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package servers

import (
	"net/http"

	common "dexprices.io/common"
	sTonFi "dexprices.io/ton/sTonFi"
)

type TONDex struct{}

func (t *TONDex) GetPoolPrice(r *http.Request, args *Args, reply *common.TokenPrice) error {
	var sTonFi sTonFi.STonFi
	if err := sTonFi.Init(""); err != nil {
		r.Response.StatusCode = 404
		return err
	}
	if err := sTonFi.GetPoolPrice(args.Pool, "latest", reply); err != nil {
		// r.Response.StatusCode = 200
		return err
	}
	return nil
}

func (t *TONDex) GetTokensPrices(r *http.Request, args *Args, reply *common.TokenPrice) error {
	var sTonFi sTonFi.STonFi
	if err := sTonFi.Init(""); err != nil {
		r.Response.StatusCode = 404
		return err
	}
	if err := sTonFi.GetPrice(args.Token0, args.Token1, "reply", reply); err != nil {
		// r.Response.StatusCode = 200
		return err
	}
	return nil
}
