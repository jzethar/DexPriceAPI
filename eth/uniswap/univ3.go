// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package eth

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	ecommon "dexprices.io/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/holiman/uint256"
)

type ContractParams struct {
	To   string `json:"to,omitempty"`
	Data string `json:"data,omitempty"`
}

// type TokenPrice struct {
// 	Token0   string `json:"token0,omitempty"`
// 	Token1   string `json:"token1,omitempty"`
// 	PriceT0  string `json:"price_of_token0,omitempty"`
// 	PriceT1  string `json:"price_of_token1,omitempty"`
// 	Decimals string `json:"decimals,omitempty"`
// }

const (
	two96Str             = "79228162514264337593543950336"
	ethCall              = "eth_call"
	slot0Str             = "slot0()"
	depthStr             = "1000000000000"
	uniV3FactoryContract = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	depth                = 9
)

// P = y / x (Price of X in terms of Y)
// sqrtPriceX96 = sqrt(P) * 2^96
// 1eth = 1 / P = 1 / (sqrtPriceX96 / 2^96)^2 -- then calculate the 10^18 / 10^6 (for usdc)
// X = token0
// Y = token1

type UniV3 struct {
	Version string
	client  *rpc.Client
	Error   string
}

func (uni *UniV3) Init(rpcProvider string) error {
	var err error
	uni.client, err = rpc.DialHTTP(rpcProvider)
	if err != nil {
		return errors.New("could not connect to the node")
	}
	defer uni.client.Close()
	return nil
}

func (uni *UniV3) GetPrice(token0 string, token1 string, block string, tokenInfo *ecommon.TokenPrice) error {
	pool, err := uni.getPool(token0, token1, block)
	if err != nil {
		return err
	}
	if err := uni.GetPoolPrice(pool, block, tokenInfo); err != nil {
		return err
	}
	return nil
}

func (uni *UniV3) GetPoolPrice(poolContract string, block string, tokenInfo *ecommon.TokenPrice) error {
	var uniSlot0 ecommon.UniSlot0
	var token0, token1 string
	var err error
	slot0 := []byte(slot0Str)
	hashSlot0 := crypto.Keccak256Hash(slot0) // don't forget to take 10 symbols, because 0x

	var result string
	req := ContractParams{poolContract, hashSlot0.String()[:10]}
	if err := uni.client.Call(&result, ethCall, req, block); err != nil {
		return err
	}

	uniSlot0.Unmarshal(result)
	token0, token1, err = uni.getTokensFromPool(poolContract, block) // evn if we have tokens set -- we need to ask a contract 4 order
	if err != nil {
		return err
	}
	tokenInfo.Token0 = token0
	tokenInfo.Token1 = token1

	decimals0, err := uni.getTokenDecimals(tokenInfo.Token0, block)
	if err != nil {
		return err
	}
	decimals1, err := uni.getTokenDecimals(tokenInfo.Token1, block)
	if err != nil {
		return err
	}

	uni.getPrices(tokenInfo, decimals0, decimals1, uniSlot0)
	return nil
}

func (uni *UniV3) getPool(token0 string, token1 string, block string) (string, error) {
	var pool string
	token0 = token0[2:]
	token1 = token1[2:]
	token0Padded := fmt.Sprintf("%064s", token0)
	token1Padded := fmt.Sprintf("%064s", token1)
	feePadded := fmt.Sprintf("%064s", strconv.FormatInt(3000, 16))
	hashDecimals0 := crypto.Keccak256Hash([]byte("getPool(address,address,uint24)"))
	hashRequest := hashDecimals0.String()[:10] + token0Padded + token1Padded + feePadded

	reqPool := ContractParams{uniV3FactoryContract, hashRequest}
	if err := uni.client.Call(&pool, ethCall, reqPool, block); err != nil {
		return "", err
	}

	return "0x" + pool[26:], nil
}

// OK
func (uni *UniV3) getTokenDecimals(token string, block string) (uint8, error) {
	var decimalsUint string
	decimals := []byte("decimals()")
	hashDecimals0 := crypto.Keccak256Hash(decimals) // don't forget to take 10 symbols, because 0x

	reqDecimals := ContractParams{token, hashDecimals0.String()[:10]}
	if err := uni.client.Call(&decimalsUint, ethCall, reqDecimals, block); err != nil {
		return 0, err
	}

	decimalsUint = decimalsUint[2:]
	return uint8(uni.hex2int(decimalsUint, 8)), nil
}

// OK
func (uni *UniV3) getTokensFromPool(pool string, block string) (string, string, error) {
	var resultToken0 string
	var resultToken1 string
	token0 := []byte("token0()")
	token1 := []byte("token1()")
	hashToken0 := crypto.Keccak256Hash(token0) // don't forget to take 10 symbols, because 0x
	hashToken1 := crypto.Keccak256Hash(token1) // don't forget to take 10 symbols, because 0x

	reqToken0 := ContractParams{pool, hashToken0.String()[:10]}
	if err := uni.client.Call(&resultToken0, ethCall, reqToken0, block); err != nil {
		return "0x", "0x", err
	}

	reqToken1 := ContractParams{pool, hashToken1.String()[:10]}
	if err := uni.client.Call(&resultToken1, ethCall, reqToken1, block); err != nil {
		return "0x", "0x", err
	}

	return "0x" + resultToken0[26:], "0x" + resultToken1[26:], nil
}

// OK
func (uni UniV3) hex2int(data string, bitsize int) int64 {
	result, _ := strconv.ParseInt(data, 16, bitsize)
	return result
}

func (uni *UniV3) getPrices(tokenInfo *ecommon.TokenPrice, decimals0 uint8, decimals1 uint8, uniSlot0 ecommon.UniSlot0) {
	var X uint256.Int
	var P uint256.Int
	var depthStr uint256.Int
	var d0d1 uint64
	var newDepth uint64
	var divSqrtTwo96 uint256.Int
	two96BI, _ := new(big.Int).SetString(two96Str, 10)
	two96, _ := uint256.FromBig(two96BI)
	depthStr.Exp(uint256.NewInt(10), uint256.NewInt(depth))
	// WE NEED TO CMP two96 and sqrtPriceX96
	// WE NEED DEPTH due to the dividing 1 to (A) we have to provide depth that gonna be 10^9 (I think it should be enough) to calculate properly it
	// But
	// Maybe even it could be set by user
	// BUT now let use 10^9
	// SHIT

	// ((10^dp / (SqrtPriceX96 / two96))^2) * 10^(d1-d0)
	log.Print("SqrtPriceX96: " + uniSlot0.SqrtPriceX96.String())
	log.Print("Dec0: " + strconv.FormatUint(uint64(decimals0), 10))
	log.Print("Dec1: " + strconv.FormatUint(uint64(decimals1), 10))
	if len(uniSlot0.SqrtPriceX96.String())-len(two96Str) < 8 {
		diff := int64(len(uniSlot0.SqrtPriceX96.String()) - len(two96Str))
		uniSlot0.SqrtPriceX96.Mul(&uniSlot0.SqrtPriceX96,
			new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil))
		newDepth = uint64(int64(depth) + diff)
		depthStr.Exp(uint256.NewInt(10), uint256.NewInt(newDepth))
	}
	if int8(decimals1)-int8(decimals0) < 0 {
		d0d1 = uint64(math.Abs(float64(decimals1) - float64(decimals0)))
	} // TODO errors d0d1 is not initialized
	divSqrtTwo96.Div(uint256.MustFromBig(&uniSlot0.SqrtPriceX96), two96)
	X.Mul(
		X.Exp(X.Div(&depthStr,
			&divSqrtTwo96), uint256.NewInt(2)),
		new(uint256.Int).Exp(uint256.NewInt(10), uint256.NewInt(d0d1)))
	log.Print(X.String())

	// (SqrtPriceX96 / two96) ^ 2
	P.Mul(P.Exp(&divSqrtTwo96, uint256.NewInt(2)),
		new(uint256.Int).Exp(uint256.NewInt(10), uint256.NewInt(uint64((newDepth+newDepth)-d0d1))))
	// AFTER ALL OPERATIONS we have depth ^ 2 or 10^18 be AWARE!!!
	tokenInfo.Decimals = strconv.FormatInt(int64(newDepth+newDepth), 10)
	tokenInfo.PriceT1 = X.String()
	tokenInfo.PriceT0 = P.String()
}
