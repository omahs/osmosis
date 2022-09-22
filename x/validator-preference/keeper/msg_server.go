package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/osmosis-labs/osmosis/v12/x/validator-preference/types"
)

type msgServer struct {
	keeper *Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) CreateValidatorSetPreference(goCtx context.Context, msg *types.MsgCreateValidatorSetPreference) (*types.MsgCreateValidatorSetPreferenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if a user already have a validator-set created
	_, found := server.keeper.GetValidatorSetPreference(ctx, msg.Owner)
	if found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("user %s already has a validator set", msg.Owner))
	}

	total_weight := sdk.NewDec(0)
	for _, val := range msg.Preferences {
		// validation checks making sure the weights add up to 1 and also the validator given is correct
		_, err := server.keeper.ValidateValidator(ctx, val.ValOperAddress)
		if err != nil {
			return nil, err
		}

		total_weight = total_weight.Add(val.Weight)
	}

	if total_weight != sdk.NewDec(1) {
		return nil, fmt.Errorf("The weights allocated to the validators do not add up to 1")
	}

	setMsg := types.ValidatorSetPreferences{
		Owner:       msg.Owner,
		Preferences: msg.Preferences,
	}

	server.keeper.SetValidatorSetPreferences(ctx, setMsg)
	return &types.MsgCreateValidatorSetPreferenceResponse{}, nil
}

func (server msgServer) StakeToValidatorSet(goCtx context.Context, msg *types.MsgStakeToValidatorSet) (*types.MsgStakeToValidatorSetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the existing validator set preference
	existingSet, found := server.keeper.GetValidatorSetPreference(ctx, msg.Owner)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("user %s doesn't have validator set", msg.Owner))
	}

	// loop through the validatorSetPreference and delegate the proportion of the tokens based on weights
	// user account address
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	tokenAmt := sdk.NewDec(msg.Coin.Amount.Int64())

	for _, val := range existingSet.Preferences {
		validator, err := server.keeper.ValidateValidator(ctx, val.ValOperAddress)
		if err != nil {
			return nil, err
		}

		// NOTE: it'd be nice if this value was decimal
		amountToStake := val.Weight.Mul(tokenAmt).RoundInt()

		_, err = server.keeper.stakingKeeper.Delegate(ctx, owner, amountToStake, validator.Status, validator, true)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgStakeToValidatorSetResponse{}, nil
}

// UnStakeFromValidatorSet unstakes all the tokens from the validator set.
// For ex: UnStake 10osmo with validator-set {ValA -> 0.5, ValB -> 0.3, ValC -> 0.2}
// our unstake logic would attempt to unstake 5osmo from A , 2osmo from B, 3osmo from C
func (server msgServer) UnStakeFromValidatorSet(goCtx context.Context, msg *types.MsgUnStakeFromValidatorSet) (*types.MsgUnStakeFromValidatorSetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the existing validator set preference
	existingSet, found := server.keeper.GetValidatorSetPreference(ctx, msg.Owner)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("user %s doesn't have validator set", msg.Owner))
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	// the total amount the user wants to unstake
	tokenAmt := sdk.NewDec(msg.Coin.Amount.Int64())

	totalAmountFromWeights := sdk.NewDec(0)
	for _, val := range existingSet.Preferences {
		// Calculate the amount to unstake based on the existing weights
		amountToUnStake := val.Weight.Mul(tokenAmt)

		// ValidateValidator gurantees that this exist
		valAddr, err := sdk.ValAddressFromBech32(val.ValOperAddress)
		if err != nil {
			return nil, err
		}

		_, err = server.keeper.stakingKeeper.Undelegate(ctx, owner, valAddr, amountToUnStake)
		if err != nil {
			return nil, err
		}

		totalAmountFromWeights = totalAmountFromWeights.Add(amountToUnStake)
	}

	if totalAmountFromWeights != tokenAmt {
		return nil, fmt.Errorf("The unstake total donot add up with the amount calculated from weights")
	}

	return &types.MsgUnStakeFromValidatorSetResponse{}, nil
}
