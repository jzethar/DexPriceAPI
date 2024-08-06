// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package ton

import (
	"math/big"
	"testing"

	common "dexprices.io/common"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/ton"
)

/*
For tests: EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE USTD/pTON pool
*/
func TestGetPoolPrice(t *testing.T) {
	var sTonFi STonFi
	sTonFi.Init("")
	got := sTonFi.GetPoolPrice("EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE", "", &common.TokenPrice{})

	if got != nil {
		t.Errorf("got %q, wanted nil", got)
	}
}

func TestGetPoolData(t *testing.T) {
	var sTonFi STonFi
	sTonFi.Init("")
	accountId := tongo.MustParseAddress("EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE")
	_, _, ta0, ta1, err := sTonFi.getPoolData(accountId)

	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	t0, err := ton.AccountIDFromTlb(ta0)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	t1, err := ton.AccountIDFromTlb(ta1)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if t0.ToRaw() != "0:4eec921b9d4d56a0d94676016d49c26b39b1d7901a7e153e66033f146a0886f4" {
		t.Errorf("got %q, wanted nil", t0.Address)
	}
	if t1.ToRaw() != "0:1150b518b2626ad51899f98887f8824b70065456455f7fe2813f012699a4061f" {
		t.Errorf("got %q, wanted nil", t1.Address)
	}
}

func TestGetTokenDecimals(t *testing.T) {
	var sTonFi STonFi
	sTonFi.Init("")
	accountId := tongo.MustParseAddress("EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE")
	_, _, ta0, ta1, err := sTonFi.getPoolData(accountId)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	d0, err := sTonFi.getTokenDecimals(ta0)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if d0 != 6 {
		t.Errorf("got %q, wanted nil", d0)
	}
	d1, err := sTonFi.getTokenDecimals(ta1)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if d1 != 9 {
		t.Errorf("got %q, wanted nil", d1)
	}
}

func TestGetJettonMaster(t *testing.T) {
	var sTonFi STonFi
	sTonFi.Init("")
	accountId := tongo.MustParseAddress("EQARULUYsmJq1RiZ-YiH-IJLcAZUVkVff-KBPwEmmaQGH6aC")
	master, err := sTonFi.getJettonMaster(accountId.ID)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	ms, err := ton.AccountIDFromTlb(master)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if ms.ToRaw() != "0:8cdc1d7640ad5ee326527fc1ad0514f468b30dc84b0173f0e155f451b4e11f7c" {
		t.Errorf("got %q, wanted nil", ms.Address)
	}
}

func TestGetPool(t *testing.T) {
	var sTonFi STonFi
	sTonFi.Init("")
	pool, err := sTonFi.getPool("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs", "EQCM3B12QK1e4yZSf8GtBRT0aLMNyEsBc_DhVfRRtOEffLez")
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	ms, err := ton.AccountIDFromTlb(pool)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if ms.ToRaw() != "0:fc4c9f311160754a99d113877cf583b78e5d16d048819cd4b820168769499d7e" {
		t.Errorf("got %q, wanted nil", ms.Address)
	}
}

func TestGetPrices(t *testing.T) {
	var sTonFi STonFi
	var tokenPrice common.TokenPrice
	sTonFi.Init("")
	err := sTonFi.getPrices(6, 9, *big.NewInt(150681859999009), *big.NewInt(22133954922754892), &tokenPrice)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if tokenPrice.PriceT0 != "146891967" {
		t.Errorf("got %q, wanted 146891967771704", tokenPrice.PriceT0)
	}
	if tokenPrice.PriceT1 != "6807724174" {
		t.Errorf("got %q, wanted 6807724174", tokenPrice.PriceT1)
	}
}
