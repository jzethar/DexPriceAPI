// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package ton

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/contract/jetton"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	dcommon "dexprices.io/common"
)

const (
	depthStr = "1000000000000"
	depth    = 9
)

type GetExpectedOutputs struct {
	Out              uint64
	Protocol_fee_out uint64
	Ref_fee_out      uint64
}

type GetPoolDataOutput struct {
	Reserve0                      uint64
	Reserve1                      uint64
	Token0_address                tlb.MsgAddress
	Token1_address                tlb.MsgAddress
	Lp_fee                        uint64
	Protocol_fee                  uint64
	Ref_fee                       uint64
	Protocol_fee_address          tlb.MsgAddress
	Collected_token0_protocol_fee uint64
	Collected_token1_protocol_fee uint64
}

type GetWalletDataOutput struct {
	Balance          tlb.Int257
	Owner            tlb.MsgAddress
	Jetton           tlb.MsgAddress
	JettonWalletCode tlb.Any
}

type GetPoolOutput struct {
	Pool tlb.MsgAddress
}

type STonFi struct {
	Version string
	client  *liteapi.Client
	Error   string
}

var MainnetRouter = ton.MustParseAccountID("EQB3ncyBUTjZUA5EnFKR5_EnOMI9V1tTEAAPaiU71gc4TiUt")

// TODO do we need another connectors???

func (uni *STonFi) Init(rpcProvider string) error {
	var err error
	uni.client, err = liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		fmt.Printf("Unable to create tongo client: %v", err)
	}
	return nil
}

/*
TODO
	1. Get metadata of each token
	2. Get `get_expected_outputs` 4 each wallet in account with decimals
	3. Count with decimals
	4. Get pool by addresses abi/get_methods.go:2061
*/

func (sTonFi *STonFi) GetPrice(token0 string, token1 string, block string, tokenInfo *dcommon.TokenPrice) error {
	pool, err := sTonFi.getPool(token0, token1)
	if err != nil {
		return err
	}
	poolStr, err := ton.AccountIDFromTlb(pool)
	if err != nil {
		return err
	}
	err = sTonFi.GetPoolPrice(poolStr.ToRaw(), block, tokenInfo)
	if err != nil {
		return err
	}
	return nil
}

func (sTonFi *STonFi) GetPoolPrice(poolContract string, block string, tokenInfo *dcommon.TokenPrice) error {
	accountId := tongo.MustParseAddress(poolContract)
	// Get data from Pool
	rs0, rs1, ta0, ta1, err := sTonFi.getPoolData(accountId)
	if err != nil {
		return err
	}
	ta0ID, err := ton.AccountIDFromTlb(ta0)
	if err != nil {
		return err
	}
	ta1ID, err := ton.AccountIDFromTlb(ta1)
	if err != nil {
		return err
	}
	// Get token decimals
	d0, err := sTonFi.getTokenDecimals(ta0)
	if err != nil {
		return err
	}
	d1, err := sTonFi.getTokenDecimals(ta1)
	if err != nil {
		return err
	}
	// Count prices
	err = sTonFi.getPrices(uint8(d0), uint8(d1), rs0, rs1, tokenInfo)
	if err != nil {
		return err
	}
	// Get tokens masters
	/*
		Unfortunately in STon.fi they keep only wallets that provides liquidity for pools
		So for knowing what is a real token is -- we need to ask these wallets for masters
	*/
	ta0Master, err := sTonFi.getJettonMaster(*ta0ID)
	if err != nil {
		return err
	}
	ta1Master, err := sTonFi.getJettonMaster(*ta1ID)
	if err != nil {
		return err
	}
	ta0MasterID, err := ton.AccountIDFromTlb(ta0Master)
	if err != nil {
		return err
	}
	ta1MasterID, err := ton.AccountIDFromTlb(ta1Master)
	if err != nil {
		return err
	}
	tokenInfo.Token0 = ta0MasterID.String()
	tokenInfo.Token1 = ta1MasterID.String()
	return nil
}

/*
	In a case dc0 < dc1:
	 1. d0d1 = dc0 - dc1
	 2. nx = x * 10^d0d1
	 3. P0 = nx / y = (nx * 10^d) / (y * 10^d) => A = (nx * 10^d) / y => P0 = A / 10^d
	 4. P1 = y / nx = (y * 10^d) / (nx * 10^d) = (y * 10^d) / (x * 10^(d+d0d1)) => B = (y * 10^d) / (x) => P1 = B / 10^(d+d0d1)

Perhaps not the best way to count, maybe we need to make decimals flow
*/
func (uni *STonFi) getPrices(decimals0 uint8, decimals1 uint8, rs0, rs1 big.Int, tokenPrice *dcommon.TokenPrice) error {
	var nrs0, nrs1 big.Int
	var p0, p1 big.Int
	d0d1 := int64(math.Abs(float64(decimals1) - float64(decimals0)))
	if decimals0 < decimals1 {
		nrs0.Mul(&rs0, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(depth+d0d1), big.NewInt(0)))
		nrs1.Mul(&rs1, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(depth-d0d1), big.NewInt(0)))
	} else {
		nrs0.Mul(&rs0, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(depth-d0d1), big.NewInt(0)))
		nrs1.Mul(&rs1, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(depth+d0d1), big.NewInt(0)))
	}
	log.Print(rs0.String())
	log.Print(rs1.String())
	p1.Div(&nrs0, &rs1)
	p0.Div(&nrs1, &rs0)
	log.Print(p0.String())
	log.Print(p1.String())
	tokenPrice.PriceT0 = p0.String()
	tokenPrice.PriceT1 = p1.String()
	tokenPrice.Decimals = strconv.FormatInt(9, 10)
	return nil
}

/*
	TODO:
		1. Find permalink on github
		2. Deal with uint or bigInt

		Source: https://tonscan.org/address/EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE#source
		In pool/get.func

		```
		(int, int, slice, slice, int, int, int, slice, int, int) get_pool_data() method_id {
		    load_storage();
		    return (
		        storage::reserve0,
		        storage::reserve1,
		        storage::token0_address,
		        storage::token1_address,
		        storage::lp_fee,
		        storage::protocol_fee,
		        storage::ref_fee,
		        storage::protocol_fee_address,
		        storage::collected_token0_protocol_fee,
		        storage::collected_token1_protocol_fee
		    );
		}
		```
*/

func (uni *STonFi) getPoolData(pool ton.Address) (big.Int, big.Int, tlb.MsgAddress, tlb.MsgAddress, error) {
	stack := tlb.VmStack{}
	rc, vmStack, err := uni.client.RunSmcMethod(context.Background(), pool.ID, "get_pool_data", stack)
	if err != nil {
		log.Print("Rc: " + strconv.FormatUint(uint64(rc), 10))
		return *big.NewInt(0), *big.NewInt(0), tlb.MsgAddress{}, tlb.MsgAddress{}, err
	}
	if len(vmStack) != 10 ||
		vmStack[0].SumType != "VmStkTinyInt" ||
		vmStack[1].SumType != "VmStkTinyInt" ||
		vmStack[2].SumType != "VmStkSlice" ||
		vmStack[3].SumType != "VmStkSlice" {
		return *big.NewInt(0), *big.NewInt(0), tlb.MsgAddress{}, tlb.MsgAddress{}, errors.New("stack corrupted")
	}
	var result GetPoolDataOutput
	_ = vmStack.Unmarshal(&result)
	return *big.NewInt(int64(result.Reserve0)), *big.NewInt(int64(result.Reserve1)), result.Token0_address, result.Token1_address, nil
}

// Do golang has overloading?
func (uni *STonFi) getPool(master0 string, master1 string) (tlb.MsgAddress, error) {
	t0 := tongo.MustParseAddress(master0)
	j0 := jetton.New(t0.ID, uni.client)
	token0S, err := j0.GetJettonWallet(context.Background(), MainnetRouter)
	if err != nil {
		return tlb.MsgAddress{}, nil
	}

	t1 := tongo.MustParseAddress(master1)
	j1 := jetton.New(t1.ID, uni.client)
	token1S, err := j1.GetJettonWallet(context.Background(), MainnetRouter)
	if err != nil {
		return tlb.MsgAddress{}, nil
	}

	stack := tlb.VmStack{}
	val0, err := tlb.TlbStructToVmCellSlice(token0S.ToMsgAddress())
	if err != nil {
		return tlb.MsgAddress{}, err
	}
	stack.Put(val0)
	val1, err := tlb.TlbStructToVmCellSlice(token1S.ToMsgAddress())
	if err != nil {
		return tlb.MsgAddress{}, err
	}
	stack.Put(val1)
	rc, vmStack, err := uni.client.RunSmcMethod(context.Background(), MainnetRouter, "get_pool_address", stack)
	if err != nil {
		log.Print("Rc: " + strconv.FormatUint(uint64(rc), 10))
		return tlb.MsgAddress{}, err
	}
	if len(vmStack) != 1 || vmStack[0].SumType != "VmStkSlice" {
		return tlb.MsgAddress{}, errors.New("stack corrupted")
	}
	var result GetPoolOutput
	_ = vmStack.Unmarshal(&result)

	return result.Pool, nil
}

func (uni *STonFi) getTokenDecimals(token tlb.MsgAddress) (int, error) {
	tokenTLB, err := ton.AccountIDFromTlb(token)
	if err != nil {
		return 0, err
	}
	j := jetton.New(*tokenTLB, uni.client)
	d, err := j.GetDecimals(context.Background())
	if err != nil {
		// log.Fatalf("Get decimals error: %v", err)
		master, err := uni.getJettonMaster(*tokenTLB)
		if err != nil {
			return 0, err
		}
		tokenTLB, err := ton.AccountIDFromTlb(master)
		if err != nil {
			return 0, err
		}
		j := jetton.New(*tokenTLB, uni.client)
		d, err := j.GetDecimals(context.Background())
		if err != nil {
			if err.Error() == "only onchain jetton data supported" {
				return 9, nil // TODO offchain support
			}
			return 0, err
		}
		return d, nil
	}
	return d, nil
}

func (uni *STonFi) getJettonMaster(jettonWallet ton.AccountID) (tlb.MsgAddress, error) {
	errCode, stack, err := uni.client.RunSmcMethod(context.Background(), jettonWallet, "get_wallet_data", tlb.VmStack{})
	if err != nil {
		return tlb.MsgAddress{}, err
	}
	if errCode == 0xFFFFFF00 { // contract not init
		return tlb.MsgAddress{}, nil
	}
	if errCode != 0 && errCode != 1 {
		return tlb.MsgAddress{}, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 4 || (stack[0].SumType != "VmStkTinyInt" && stack[0].SumType != "VmStkInt") ||
		stack[1].SumType != "VmStkSlice" ||
		stack[2].SumType != "VmStkSlice" ||
		stack[3].SumType != "VmStkCell" {
		return tlb.MsgAddress{}, fmt.Errorf("invalid stack")
	}
	var result GetWalletDataOutput
	_ = stack.Unmarshal(&result)
	return result.Jetton, nil
}
