package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/x/superfluid/types"
)

func (k Keeper) HandleSetSuperfluidAssetsProposal(ctx sdk.Context, p *types.SetSuperfluidAssetsProposal) error {
	for _, asset := range p.Assets {
		k.SetSuperfluidAsset(ctx, asset)
		event := sdk.NewEvent(
			types.TypeEvtSetSuperfluidAsset,
			sdk.NewAttribute(types.AttributeDenom, asset.Denom),
			sdk.NewAttribute(types.AttributeSuperfluidAssetType, asset.AssetType.String()),
		)
		ctx.EventManager().EmitEvent(event)
	}
	return nil
}

func (k Keeper) HandleRemoveSuperfluidAssetsProposal(ctx sdk.Context, p *types.RemoveSuperfluidAssetsProposal) error {
	for _, denom := range p.SuperfluidAssetDenoms {
		// TODO: Why do we need this SetEpochOsmoEquivalentTWAP, why don't we have a unified delete?
		k.SetEpochOsmoEquivalentTWAP(ctx, 0, denom, sdk.ZeroDec())
		k.DeleteSuperfluidAsset(ctx, denom)
		event := sdk.NewEvent(
			types.TypeEvtRemoveSuperfluidAsset,
			sdk.NewAttribute(types.AttributeDenom, denom),
		)
		ctx.EventManager().EmitEvent(event)
	}
	return nil
}