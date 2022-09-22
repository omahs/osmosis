package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) ValidateValidator(ctx sdk.Context, valOperAddress string) (stakingtypes.Validator, error) {
	vals, err := sdk.ValAddressFromBech32(valOperAddress)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("validator address not formatted")
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, vals)
	if !found {
		return stakingtypes.Validator{}, fmt.Errorf("validator doesnot exist %s", validator)
	}

	return validator, nil
}
