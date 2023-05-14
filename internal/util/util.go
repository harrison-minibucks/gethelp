package util

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/params"
)

func ConvertAmount(amount string) (*big.Int, error) {
	bigIntAmount := new(big.Int)
	parts := strings.Split(amount, " ")
	bigIntAmount.SetString(parts[0], 10)
	if len(parts) > 1 {
		unit := strings.ToLower(parts[1])
		switch unit {
		case "eth", "ether":
			bigIntAmount.Mul(bigIntAmount, big.NewInt(params.Ether))
		case "gwei":
			bigIntAmount.Mul(bigIntAmount, big.NewInt(params.GWei))
		case "wei":
		default:
			return nil, errors.New("invalid unit, use ether, gwei or wei")
		}
	}
	return bigIntAmount, nil
}

func FormatEth(wei *big.Int) string {
	// Represent with Wei if smaller than Gwei
	if wei.Cmp(big.NewInt(params.GWei)) == -1 {
		return fmt.Sprintf("%d Wei", wei)
	}
	// Represent with Gwei if smaller than Ether
	if wei.Cmp(big.NewInt(params.Ether)) == -1 {
		return formatDecimalPoints(wei, params.GWei, 5) + " Gwei"
	}
	// Show up to 5 decimals in Ether
	return formatDecimalPoints(wei, params.Ether, 5) + " Ether"
}

// Format up to a specific decimal place
func formatDecimalPoints(n *big.Int, unit float64, places int) string {
	number := new(big.Float).SetInt(n)
	number.Quo(number, big.NewFloat(unit))
	parts := strings.Split(number.Text('f', len(big.NewFloat(unit).String())), ".")
	// Hide depending on the decimal place
	if len(parts) < 2 || strings.HasPrefix(parts[1], strings.Repeat("0", places)) {
		return parts[0]
	} else {
		return fmt.Sprintf("%s.%s", parts[0], parts[1][:5])
	}
}
