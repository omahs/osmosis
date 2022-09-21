package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/osmosis-labs/osmosis/v12/x/validator-preference/keeper"
	"github.com/osmosis-labs/osmosis/v12/x/validator-preference/types"
)

func (suite *KeeperTestSuite) TestCreateValidatorSetPreference() {

	// setup 3 validators
	valAddrs := []string{}
	for i := 0; i < 3; i++ {
		valAddr := suite.SetupValidators([]stakingtypes.BondStatus{stakingtypes.Bonded})
		valAddrs = append(valAddrs, valAddr[0].String())
	}

	type param struct {
		owner       sdk.AccAddress
		preferences []types.ValidatorPreference
	}

	tests := []struct {
		name       string
		param      param
		expectPass bool
	}{
		{
			name: "creation of new validator set",
			param: param{
				owner: sdk.AccAddress([]byte("addr1---------------")),
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: valAddrs[0],
						Weight:         sdk.NewDecWithPrec(5, 1),
					},
					{
						ValOperAddress: valAddrs[1],
						Weight:         sdk.NewDecWithPrec(3, 1),
					},
					{
						ValOperAddress: valAddrs[2],
						Weight:         sdk.NewDecWithPrec(2, 1),
					},
				},
			},
			expectPass: true,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			msgServer := keeper.NewMsgServerImpl(suite.App.ValidatorPreferenceKeeper)
			c := sdk.WrapSDKContext(suite.Ctx)

			_, err := msgServer.CreateValidatorSetPreference(c, types.NewMsgCreateValidatorSetPreference(test.param.owner, test.param.preferences))

			if test.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

		})
	}
}

func (suite *KeeperTestSuite) TestStakeToValidatorSet() {

}

func (suite *KeeperTestSuite) TestUnStakeFromValidatorSet() {

}

/*

	{
			name: "create validator set with unknown validator address",
			param: param{
				owner: "",
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: "",
						Weight:         sdk.NewDec(1),
					},
				},
			},
			expectPass: false,
		},
		{
			name: "create validator set with weights != 1",
			param: param{
				owner: "",
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: "",
						Weight:         sdk.NewDec(1),
					},
				},
			},
			expectPass: false,
		},
		{
			name: "creator validator set with invalid Owner",
			param: param{
				owner: "",
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: "",
						Weight:         sdk.NewDec(1),
					},
				},
			},
			expectPass: false,
		},

*/
