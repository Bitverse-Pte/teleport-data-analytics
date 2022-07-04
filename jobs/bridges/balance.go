package bridges

//func (r *Bridges) QueryPoolBalance(chainName, address string) (chains.TokenAmount, error) {
//	var amount chains.TokenAmount
//	if r.TokenMap[chainName] == nil {
//		return amount, fmt.Errorf("invalid chainName")
//	}
//	tokenQuery := r.TokenMap[chainName][address]
//	if tokenQuery == nil {
//		return amount, fmt.Errorf("invalid token")
//	}
//	decimals, err := tokenQuery.GetDecimals()
//	if err != nil {
//		return amount, fmt.Errorf("tokenQuery.GetDecimals error:%+v", err)
//	}
//	amount.Amount = big.NewInt(-1)
//	amount.Decimals = decimals
//	if chainName == chains.TeleportChain {
//		if address == chains.ZeroAddress {
//			return amount, nil
//		} else {
//			return tokenQuery.GetInTokenAmount(address, chainName, nil)
//		}
//	}
//	if r.AddressToName[fmt.Sprintf("%v/%v", chainName, address)] == "tele" {
//		return tokenQuery.GetInTokenAmount(address, chainName, nil)
//	}
//	return amount, nil
//}
