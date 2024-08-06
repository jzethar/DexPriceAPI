// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package servers

import (
	"net/http"

	common "dexprices.io/common"
	uniSwap "dexprices.io/eth/uniswap"
)

type Dex struct{}

func (t *Dex) GetPoolPrice(r *http.Request, args *Args, reply *common.TokenPrice) error {
	var uniV3 uniSwap.UniV3
	if err := uniV3.Init("https://eth.llamarpc.com"); err != nil {
		r.Response.StatusCode = 404
		return err
	}
	if err := uniV3.GetPoolPrice(args.Pool, "latest", reply); err != nil {
		// r.Response.StatusCode = 200
		return err
	}
	return nil
}

func (t *Dex) GetTokensPrices(r *http.Request, args *Args, reply *common.TokenPrice) error {
	var uniV3 uniSwap.UniV3
	if err := uniV3.Init("https://eth.llamarpc.com"); err != nil {
		r.Response.StatusCode = 404
		return err
	}
	if err := uniV3.GetPrice(args.Token0, args.Token1, "latest", reply); err != nil {
		// r.Response.StatusCode = 200
		return err
	}
	return nil
}
