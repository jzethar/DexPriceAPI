// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package common

import (
	"math/big"
	"strconv"
)

// Golang doesn’t check for overflows implicitly and
// so this may lead to unexpected results when a number
// larger than 64 bits are stored in a int64.
type UniSlot0 struct {
	SqrtPriceX96               big.Int // uint160
	tick                       int32
	observationIndex           uint16
	observationCardinality     uint16
	observationCardinalityNext uint16
	feeProtocol                uint8
	unlocked                   bool
}

func (uns *UniSlot0) hex2int(data string, bitsize int) int64 {
	result, _ := strconv.ParseInt(data, 16, bitsize)
	return result
}

func (uns *UniSlot0) Unmarshal(data string) error {
	data = data[2:]

	SqrtPriceX96 := data[:64]
	uns.SqrtPriceX96.SetString(SqrtPriceX96, 16)
	data = data[64:]

	tick := data[:64]
	uns.tick = int32(uns.hex2int(tick, 32))
	data = data[64:]

	observationIndex := data[:64]
	uns.observationIndex = uint16(uns.hex2int(observationIndex, 16))
	data = data[64:]

	observationCardinality := data[:64]
	uns.observationCardinality = uint16(uns.hex2int(observationCardinality, 16))
	data = data[64:]

	observationCardinalityNext := data[:64]
	uns.observationCardinalityNext = uint16(uns.hex2int(observationCardinalityNext, 16))
	data = data[64:]

	feeProtocol := data[:64]
	uns.feeProtocol = uint8(uns.hex2int(feeProtocol, 8))
	data = data[64:]

	if string(data[len(data)-1]) == "1" {
		uns.unlocked = true
	} else {
		uns.unlocked = false
	}
	return nil
}
