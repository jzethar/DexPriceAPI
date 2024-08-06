// Copyright 2024, Yulian Volianskyi aka jzethar, All rights reserved.
// This code is a part of DexPriceAPI project
// See the LICENSE file

package servers

// Args holds arguments passed to JSON RPC service
type Args struct {
	Pool   string `json:"pool,omitempty"`
	Token0 string `json:"token0,omitempty"`
	Token1 string `json:"token1,omitempty"`
	Block  string `json:"block,omitempty"`
}
