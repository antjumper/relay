/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package util

import (
	"errors"
	"fmt"
	"github.com/Loopring/relay/dao"
	"github.com/Loopring/relay/log"
	"github.com/Loopring/relay/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

const weiToEther = 1e18

type TokenPair struct {
	TokenS common.Address
	TokenB common.Address
}

func ByteToFloat(amount []byte) float64 {
	var rst big.Int
	rst.UnmarshalText(amount)
	return float64(rst.Int64()) / weiToEther
}

func FloatToByte(amount float64) []byte {
	rst, _ := big.NewInt(int64(amount * weiToEther)).MarshalText()
	return rst
}

var (
	SupportTokens  map[string]types.Token // token symbol to entity
	AllTokens      map[string]types.Token
	SupportMarkets map[string]types.Token // token symbol to contract hex address
	AllMarkets     []string
	AllTokenPairs  []TokenPair
)

var ContractVersionConfig = map[string]string{
	"v1.0": "0x39kdjfskdfjsdfj",
	"v1.2": "0x39kdjfskdfjsdfj",
}

func Initialize(rds dao.RdsService) {
	SupportTokens = make(map[string]types.Token)
	SupportMarkets = make(map[string]types.Token)
	AllTokens = make(map[string]types.Token)

	tokens, err := rds.FindUnDeniedTokens()
	if err != nil {
		log.Fatalf("market util cann't find any token!")
	}
	markets, err := rds.FindUnDeniedMarkets()
	if err != nil {
		log.Fatalf("market util cann't find any base market!")
	}

	// set support tokens
	for _, v := range tokens {
		var token types.Token
		v.ConvertUp(&token)
		SupportTokens[v.Symbol] = token
		log.Infof("supported token %s->%s", token.Symbol, token.Protocol.Hex())
	}

	// set all tokens
	for k, v := range SupportTokens {
		AllTokens[k] = v
	}

	// set support markets
	for _, v := range markets {
		var token types.Token
		v.ConvertUp(&token)
		SupportMarkets[token.Symbol] = token
	}

	// set all markets
	for _, k := range SupportTokens { // lrc,omg
		for _, kk := range SupportMarkets { //eth
			symbol := k.Symbol + "-" + kk.Symbol
			AllMarkets = append(AllMarkets, symbol)
			log.Infof("supported market:%s", symbol)
		}
	}

	// set all token pairs
	pairsMap := make(map[string]TokenPair, 0)
	for _, v := range SupportMarkets {
		for _, vv := range SupportTokens {
			pairsMap[v.Symbol+"-"+vv.Symbol] = TokenPair{v.Protocol, vv.Protocol}
			pairsMap[vv.Symbol+"-"+v.Symbol] = TokenPair{vv.Protocol, v.Protocol}
		}
	}
	for _, v := range pairsMap {
		AllTokenPairs = append(AllTokenPairs, v)
	}
}

// todo: add token, delete token ...

func WrapMarket(s, b string) (market string, err error) {

	s, b = strings.ToLower(s), strings.ToLower(b)

	if IsSupportedMarket(s) && IsSupportedToken(b) {
		market = fmt.Sprintf("%s-%s", b, s)
	} else if IsSupportedMarket(b) && IsSupportedToken(s) {
		market = fmt.Sprintf("%s-%s", s, b)
	} else {
		err = errors.New(fmt.Sprintf("not supported market type : %s-%s", s, b))
	}
	return
}

func WrapMarketByAddress(s, b string) (market string, err error) {
	return WrapMarket(AddressToAlias(s), AddressToAlias(b))
}

func UnWrap(market string) (s, b string) {
	mkt := strings.Split(strings.TrimSpace(market), "-")
	if len(mkt) != 2 {
		return "", ""
	}

	s, b = strings.ToLower(mkt[0]), strings.ToLower(mkt[1])
	return
}

func UnWrapToAddress(market string) (s, b common.Address) {
	sa, sb := UnWrap(market)
	return common.StringToAddress(sa), common.StringToAddress(sb)
}

func IsSupportedMarket(market string) bool {
	_, ok := SupportMarkets[market]
	return ok
}

func IsSupportedToken(token string) bool {
	_, ok := SupportTokens[token]
	return ok
}

func AliasToAddress(t string) common.Address {
	return AllTokens[t].Protocol
}

func AddressToAlias(t string) string {
	for k, v := range AllTokens {
		if t == v.Protocol.Hex() {
			return k
		}
	}
	return ""
}

func CalculatePrice(amountS, amountB []byte, s, b string) float64 {

	as := ByteToFloat(amountS)
	ab := ByteToFloat(amountB)

	if as == 0 || ab == 0 {
		return 0
	}

	if IsBuy(s) {
		return ab / as
	}

	return as / ab

}

func IsBuy(s string) bool {
	if IsAddress(s) {
		s = AddressToAlias(s)
	}
	if _, ok := SupportTokens[s]; !ok {
		return false
	}
	return true
}

func IsAddress(token string) bool {
	return strings.HasPrefix(token, "0x")
}

func getContractVersion(address string) string {
	for k, v := range ContractVersionConfig {
		if v == address {
			return k
		}
	}
	return ""
}

func IsSupportedContract(address string) bool {
	return getContractVersion(address) != ""
}