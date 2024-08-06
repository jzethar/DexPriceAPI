// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package eth

import (
	"strings"
	"testing"

	ecommon "dexprices.io/common"
)

func TestGetPoolPrice(t *testing.T) {
	var uni UniV3
	uni.Init("https://eth.llamarpc.com")
	got := uni.GetPoolPrice("0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35", "latest", &ecommon.TokenPrice{})

	if got != nil {
		t.Errorf("got %q, wanted nil", got)
	}
}

func TestGetTokensFromPool(t *testing.T) {
	var uni UniV3
	uni.Init("https://eth.llamarpc.com")
	t0, t1, err := uni.getTokensFromPool("0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35", "latest")

	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if t0 != "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599" {
		t.Errorf("got %q, wanted 0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", t0)
	}
	if t1 != "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" {
		t.Errorf("got %q, wanted ", t1)
	}
}

func TestGetTokenDecimals(t *testing.T) {
	var uni UniV3
	uni.Init("https://eth.llamarpc.com")
	d0, err := uni.getTokenDecimals("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", "latest")

	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if d0 != 8 {
		t.Errorf("got %d, wanted 8", d0)
	}
}

func TestHex2Int(t *testing.T) {
	var uni UniV3
	d0 := uni.hex2int("0000000000000000000000000000000000000000000000000000000000000008", 8)

	if d0 != 8 {
		t.Errorf("got %d, wanted 8", d0)
	}
}

func TestGetPool(t *testing.T) {
	var uni UniV3
	uni.Init("https://eth.llamarpc.com")
	pool, err := uni.getPool("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", "latest")

	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if pool != strings.ToLower("0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35") {
		t.Errorf("got %q, wanted 0x99ac8cA7087fA4A2A1FB6357269965A2014ABc35", pool)
	}
}
