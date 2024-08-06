// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package common

type TokenPrice struct {
	Token0   string `json:"token0,omitempty"`
	Token1   string `json:"token1,omitempty"`
	PriceT0  string `json:"price_of_token0,omitempty"`
	PriceT1  string `json:"price_of_token1,omitempty"`
	Decimals string `json:"decimals,omitempty"`
}
