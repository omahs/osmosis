package twap_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/osmosis-labs/osmosis/v12/app/apptesting/osmoassert"
	"github.com/osmosis-labs/osmosis/v12/x/twap"
	"github.com/osmosis-labs/osmosis/v12/x/twap/types"
	"github.com/osmosis-labs/osmosis/v12/x/twap/types/twapmock"
)

var zeroDec = sdk.ZeroDec()
var oneDec = sdk.OneDec()
var twoDec = oneDec.Add(oneDec)
var OneSec = sdk.MustNewDecFromStr("1000.000000000000000000")

func newRecord(poolId uint64, t time.Time, sp0, accum0, accum1 sdk.Dec) []types.TwapRecord {
	return []types.TwapRecord{{
		PoolId:          poolId,
		Asset0Denom:     defaultTwoAssetCoins[0].Denom,
		Asset1Denom:     defaultTwoAssetCoins[1].Denom,
		Time:            t,
		P0LastSpotPrice: sp0,
		P1LastSpotPrice: sdk.OneDec().Quo(sp0),
		// make new copies
		P0ArithmeticTwapAccumulator: accum0.Add(sdk.ZeroDec()),
		P1ArithmeticTwapAccumulator: accum1.Add(sdk.ZeroDec()),
	}}
}

// make an expected record for math tests, we adjust other values in the test runner.
func newExpRecord(poolId uint64, accum0, accum1 sdk.Dec) []types.TwapRecord {
	return []types.TwapRecord{{
		PoolId:      poolId,
		Asset0Denom: defaultTwoAssetCoins[0].Denom,
		Asset1Denom: defaultTwoAssetCoins[1].Denom,
		// make new copies
		P0ArithmeticTwapAccumulator: accum0.Add(sdk.ZeroDec()),
		P1ArithmeticTwapAccumulator: accum1.Add(sdk.ZeroDec()),
	}}
}

func newTapRecord(poolId uint64, t time.Time, sp0, accum0, accum1 sdk.Dec) []types.TwapRecord {
	twapAB := types.TwapRecord{
		PoolId:          poolId,
		Asset0Denom:     defaultThreeAssetCoins[0].Denom,
		Asset1Denom:     defaultThreeAssetCoins[1].Denom,
		Time:            t,
		P0LastSpotPrice: sp0,
		P1LastSpotPrice: sdk.OneDec().Quo(sp0),
		// make new copies
		P0ArithmeticTwapAccumulator: accum0.Add(sdk.ZeroDec()),
		P1ArithmeticTwapAccumulator: accum1.Add(sdk.ZeroDec()),
	}
	twapAC := twapAB
	twapAC.Asset1Denom = denom2
	twapBC := twapAC
	twapBC.Asset0Denom = denom1
	return []types.TwapRecord{twapAB, twapAC, twapBC}
}

// make an expected record for math tests, we adjust other values in the test runner.
func newTapExpRecord(poolId uint64, accum0, accum1 sdk.Dec) []types.TwapRecord {
	twapAB := types.TwapRecord{
		PoolId:      poolId,
		Asset0Denom: defaultThreeAssetCoins[0].Denom,
		Asset1Denom: defaultThreeAssetCoins[1].Denom,
		// make new copies
		P0ArithmeticTwapAccumulator: accum0.Add(sdk.ZeroDec()),
		P1ArithmeticTwapAccumulator: accum1.Add(sdk.ZeroDec()),
	}
	twapAC := twapAB
	twapAC.Asset1Denom = denom2
	twapBC := twapAC
	twapBC.Asset0Denom = denom1
	return []types.TwapRecord{twapAB, twapAC, twapBC}
}

func (s *TestSuite) TestNewTwapRecord() {
	// prepare pool before test
	poolId := s.PrepareBalancerPoolWithCoins(defaultTwoAssetCoins...)

	tests := map[string]struct {
		poolId        uint64
		denom0        string
		denom1        string
		expectedErr   error
		expectedPanic bool
	}{
		"denom with lexicographical order": {
			poolId,
			denom0,
			denom1,
			nil,
			false,
		},
		"denom with non-lexicographical order": {
			poolId,
			denom1,
			denom0,
			nil,
			false,
		},
		"new record with same denom": {
			poolId,
			denom0,
			denom0,
			fmt.Errorf("both assets cannot be of the same denom: assetA: %s, assetB: %s", denom0, denom0),
			false,
		},
		"error in getting spot price": {
			poolId + 1,
			denom1,
			denom0,
			nil,
			true,
		},
	}
	for name, test := range tests {
		s.Run(name, func() {
			twapRecord, err := twap.NewTwapRecord(s.App.GAMMKeeper, s.Ctx, test.poolId, test.denom0, test.denom1)

			if test.expectedPanic {
				s.Require().Equal(twapRecord.LastErrorTime, s.Ctx.BlockTime())
			} else if test.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(test.expectedErr.Error(), err.Error())
			} else {
				s.Require().NoError(err)

				s.Require().Equal(test.poolId, twapRecord.PoolId)
				s.Require().Equal(s.Ctx.BlockHeight(), twapRecord.Height)
				s.Require().Equal(s.Ctx.BlockTime(), twapRecord.Time)
				s.Require().Equal(sdk.ZeroDec(), twapRecord.P0ArithmeticTwapAccumulator)
				s.Require().Equal(sdk.ZeroDec(), twapRecord.P1ArithmeticTwapAccumulator)
			}

		})
	}
}

func (s *TestSuite) TestUpdateTwap() {
	s.PrepareBalancerPoolWithCoins(defaultTwoAssetCoins...)
	s.PrepareBalancerPoolWithCoins(defaultThreeAssetCoins...)
	programmableAmmInterface := twapmock.NewProgrammedAmmInterface(s.App.TwapKeeper.GetAmmInterface())
	s.App.TwapKeeper.SetAmmInterface(programmableAmmInterface)

	spotPriceResOne := twapmock.SpotPriceResult{Sp: sdk.OneDec(), Err: nil}
	spotPriceResOneErr := twapmock.SpotPriceResult{Sp: sdk.OneDec(), Err: errors.New("dummy err")}
	spotPriceResOneErrNilDec := twapmock.SpotPriceResult{Sp: sdk.Dec{}, Err: errors.New("dummy err")}
	baseTime := time.Unix(2, 0).UTC()
	updateTime := time.Unix(3, 0).UTC()
	baseTimeMinusOne := time.Unix(1, 0).UTC()

	zeroAccumNoErrSp10Record := newRecord(1, baseTime, sdk.NewDec(10), zeroDec, zeroDec)
	tapZeroAccumNoErrSp10Record := newTapRecord(2, baseTime, sdk.NewDec(10), zeroDec, zeroDec)
	// all tests occur with updateTime = base time + time.Unix(1, 0)
	tests := map[string]struct {
		record           []types.TwapRecord
		spotPriceResult0 twapmock.SpotPriceResult
		spotPriceResult1 twapmock.SpotPriceResult
		expRecord        []types.TwapRecord
	}{
		"0 accum start, sp change": {
			record:           zeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOne,
			expRecord:        newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)),
		},
		"three asset, 0 accum start, sp change": {
			record:           tapZeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOne,
			expRecord:        newTapExpRecord(2, OneSec.MulInt64(10), OneSec.QuoInt64(10)),
		},
		"0 accum start, sp0 err at update": {
			record:           zeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOneErr,
			spotPriceResult1: spotPriceResOne,
			expRecord:        withLastErrTime(newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime),
		},
		"three asset, 0 accum start, sp0 err at update": {
			record:           tapZeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOneErr,
			spotPriceResult1: spotPriceResOne,
			expRecord:        withLastErrTime(newTapExpRecord(2, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime),
		},
		"0 accum start, sp0 err at update with nil dec": {
			record:           zeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOneErrNilDec,
			spotPriceResult1: spotPriceResOne,
			expRecord:        withSp0(withLastErrTime(newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime), sdk.ZeroDec()),
		},
		"three asset, 0 accum start, sp0 err at update with nil dec": {
			record:           tapZeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOneErrNilDec,
			spotPriceResult1: spotPriceResOne,
			expRecord:        withSp0(withLastErrTime(newTapExpRecord(2, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime), sdk.ZeroDec()),
		},
		"0 accum start, sp1 err at update with nil dec": {
			record:           zeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOneErrNilDec,
			expRecord:        withSp1(withLastErrTime(newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime), sdk.ZeroDec()),
		},
		"three asset, 0 accum start, sp1 err at update with nil dec": {
			record:           tapZeroAccumNoErrSp10Record,
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOneErrNilDec,
			expRecord:        withSp1(withLastErrTime(newTapExpRecord(2, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime), sdk.ZeroDec()),
		},
		"startRecord err time preserved": {
			record:           withLastErrTime(newRecord(1, baseTime, sdk.NewDec(10), zeroDec, zeroDec), baseTimeMinusOne),
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOne,
			expRecord:        withLastErrTime(newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)), baseTimeMinusOne),
		},
		"three asset, startRecord err time preserved": {
			record:           withLastErrTime(newTapRecord(2, baseTime, sdk.NewDec(10), zeroDec, zeroDec), baseTimeMinusOne),
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOne,
			expRecord:        withLastErrTime(newTapExpRecord(2, OneSec.MulInt64(10), OneSec.QuoInt64(10)), baseTimeMinusOne),
		},
		"err time bumped with start": {
			record:           withLastErrTime(newRecord(1, baseTime, sdk.NewDec(10), zeroDec, zeroDec), baseTimeMinusOne),
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOneErr,
			expRecord:        withLastErrTime(newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime),
		},
		"three asset, err time bumped with start": {
			record:           withLastErrTime(newTapRecord(2, baseTime, sdk.NewDec(10), zeroDec, zeroDec), baseTimeMinusOne),
			spotPriceResult0: spotPriceResOne,
			spotPriceResult1: spotPriceResOneErr,
			expRecord:        withLastErrTime(newTapExpRecord(2, OneSec.MulInt64(10), OneSec.QuoInt64(10)), updateTime),
		},
	}
	for name, test := range tests {
		s.Run(name, func() {
			s.Ctx = s.Ctx.WithBlockTime(updateTime)
			for i := range test.record {
				// setup common, block time, pool Id, expected spot prices
				if (test.expRecord[i].P0LastSpotPrice == sdk.Dec{}) {
					test.expRecord[i].P0LastSpotPrice = test.spotPriceResult0.Sp
				}
				if (test.expRecord[i].P1LastSpotPrice == sdk.Dec{}) {
					test.expRecord[i].P1LastSpotPrice = test.spotPriceResult1.Sp
				}
				test.expRecord[i].Height = s.Ctx.BlockHeight()
				test.expRecord[i].Time = s.Ctx.BlockTime()

				programmableAmmInterface.ProgramPoolSpotPriceOverride(test.record[i].PoolId,
					test.record[i].Asset0Denom, test.record[i].Asset1Denom,
					test.spotPriceResult0.Sp, test.spotPriceResult0.Err)
				programmableAmmInterface.ProgramPoolSpotPriceOverride(test.record[i].PoolId,
					test.record[i].Asset1Denom, test.record[i].Asset0Denom,
					test.spotPriceResult1.Sp, test.spotPriceResult1.Err)

				newRecord := s.twapkeeper.UpdateRecord(s.Ctx, test.record[i])
				s.Equal(test.expRecord[i], newRecord)
			}
		})
	}
}

func TestRecordWithUpdatedAccumulators(t *testing.T) {

	tests := map[string]struct {
		record          []types.TwapRecord
		interpolateTime time.Time
		expRecord       []types.TwapRecord
	}{
		"0accum": {
			record:          newRecord(1, time.Unix(1, 0), sdk.NewDec(10), zeroDec, zeroDec),
			interpolateTime: time.Unix(2, 0),
			expRecord:       newExpRecord(1, OneSec.MulInt64(10), OneSec.QuoInt64(10)),
		},
		"small starting accumulators": {
			record:          newRecord(1, time.Unix(1, 0), sdk.NewDec(10), oneDec, twoDec),
			interpolateTime: time.Unix(2, 0),
			expRecord:       newExpRecord(1, oneDec.Add(OneSec.MulInt64(10)), twoDec.Add(OneSec.QuoInt64(10))),
		},
		"larger time interval": {
			record:          newRecord(1, time.Unix(11, 0), sdk.NewDec(10), oneDec, twoDec),
			interpolateTime: time.Unix(55, 0),
			expRecord:       newExpRecord(1, oneDec.Add(OneSec.MulInt64(44*10)), twoDec.Add(OneSec.MulInt64(44).QuoInt64(10))),
		},
		"same time": {
			record:          newRecord(1, time.Unix(1, 0), sdk.NewDec(10), oneDec, twoDec),
			interpolateTime: time.Unix(1, 0),
			expRecord:       newExpRecord(1, oneDec, twoDec),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for i := range test.record {
				// correct expected record based off copy/paste values
				test.expRecord[i].Time = test.interpolateTime
				test.expRecord[i].P0LastSpotPrice = test.record[i].P0LastSpotPrice
				test.expRecord[i].P1LastSpotPrice = test.record[i].P1LastSpotPrice

				gotRecord := twap.RecordWithUpdatedAccumulators(test.record[i], test.interpolateTime)
				require.Equal(t, test.expRecord[i], gotRecord)
			}
		},
		)
	}
}

func newOneSidedRecord(time time.Time, accum sdk.Dec, useP0 bool) []types.TwapRecord {
	record := types.TwapRecord{Time: time, Asset0Denom: denom0, Asset1Denom: denom1}
	if useP0 {
		record.P0ArithmeticTwapAccumulator = accum
	} else {
		record.P1ArithmeticTwapAccumulator = accum
	}
	record.P0LastSpotPrice = sdk.ZeroDec()
	record.P1LastSpotPrice = sdk.OneDec()
	records := []types.TwapRecord{record}
	return records
}

func newTapOneSidedRecord(time time.Time, accum sdk.Dec, useP0 bool) []types.TwapRecord {
	record := types.TwapRecord{Time: time, Asset0Denom: denom0, Asset1Denom: denom1}
	if useP0 {
		record.P0ArithmeticTwapAccumulator = accum
	} else {
		record.P1ArithmeticTwapAccumulator = accum
	}
	record.P0LastSpotPrice = sdk.ZeroDec()
	record.P1LastSpotPrice = sdk.OneDec()
	records := []types.TwapRecord{record, record, record}
	records[1].Asset1Denom = denom2
	records[2].Asset0Denom = denom1
	records[2].Asset1Denom = denom2
	return records
}

type computeArithmeticTwapTestCase struct {
	startRecord []types.TwapRecord
	endRecord   []types.TwapRecord
	quoteAsset  []string
	expTwap     sdk.Dec
	expErr      bool
}

// TestComputeArithmeticTwap tests ComputeArithmeticTwap on various inputs.
// The test vectors are structured by setting up different start and records,
// based on time interval, and their accumulator values.
// Then an expected TWAP is provided in each test case, to compare against computed.
func TestComputeArithmeticTwap(t *testing.T) {
	testCaseFromDeltas := func(startAccum, accumDiff sdk.Dec, timeDelta time.Duration, expectedTwap sdk.Dec) computeArithmeticTwapTestCase {
		return computeArithmeticTwapTestCase{
			newOneSidedRecord(baseTime, startAccum, true),
			newOneSidedRecord(baseTime.Add(timeDelta), startAccum.Add(accumDiff), true),
			[]string{denom0},
			expectedTwap,
			false,
		}
	}

	testTapCaseFromDeltas := func(startAccum, accumDiff sdk.Dec, timeDelta time.Duration, expectedTwap sdk.Dec) computeArithmeticTwapTestCase {
		return computeArithmeticTwapTestCase{
			newTapOneSidedRecord(baseTime, startAccum, true),
			newTapOneSidedRecord(baseTime.Add(timeDelta), startAccum.Add(accumDiff), true),
			[]string{denom0, denom0, denom1},
			expectedTwap,
			false,
		}
	}

	testCaseFromDeltasAsset1 := func(startAccum, accumDiff sdk.Dec, timeDelta time.Duration, expectedTwap sdk.Dec) computeArithmeticTwapTestCase {
		return computeArithmeticTwapTestCase{
			newOneSidedRecord(baseTime, startAccum, false),
			newOneSidedRecord(baseTime.Add(timeDelta), startAccum.Add(accumDiff), false),
			[]string{denom1},
			expectedTwap,
			false,
		}
	}

	testTapCaseFromDeltasAsset1 := func(startAccum, accumDiff sdk.Dec, timeDelta time.Duration, expectedTwap sdk.Dec) computeArithmeticTwapTestCase {
		return computeArithmeticTwapTestCase{
			newTapOneSidedRecord(baseTime, startAccum, false),
			newTapOneSidedRecord(baseTime.Add(timeDelta), startAccum.Add(accumDiff), false),
			[]string{denom1, denom2, denom2},
			expectedTwap,
			false,
		}
	}

	tenSecAccum := OneSec.MulInt64(10)
	pointOneAccum := OneSec.QuoInt64(10)
	tests := map[string]computeArithmeticTwapTestCase{
		"basic: spot price = 1 for one second, 0 init accumulator": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecord(tPlusOne, OneSec, true),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.OneDec(),
		},
		"three asset basic: spot price = 1 for one second, 0 init accumulator": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecord(tPlusOne, OneSec, true),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.OneDec(),
		},
		// this test just shows what happens in case the records are reversed.
		// It should return the correct result, even though this is incorrect internal API usage
		"invalid call: reversed records of above": {
			startRecord: newOneSidedRecord(tPlusOne, OneSec, true),
			endRecord:   newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.OneDec(),
		},
		"three asset invalid call: reversed records of above": {
			startRecord: newTapOneSidedRecord(tPlusOne, OneSec, true),
			endRecord:   newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.OneDec(),
		},
		"same record: denom0, end spot price = 0": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.ZeroDec(),
		},
		"three asset same record: asset0, end spot price = 0": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.ZeroDec(),
		},
		"same record: denom1, end spot price = 1": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			quoteAsset:  []string{denom1},
			expTwap:     sdk.OneDec(),
		},
		"three asset same record: asset1, end spot price = 1": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			quoteAsset:  []string{denom1, denom2, denom2},
			expTwap:     sdk.OneDec(),
		},
		"accumulator = 10*OneSec, t=5s. 0 base accum": testCaseFromDeltas(
			sdk.ZeroDec(), tenSecAccum, 5*time.Second, sdk.NewDec(2)),
		"three asset. accumulator = 10*OneSec, t=5s. 0 base accum": testTapCaseFromDeltas(
			sdk.ZeroDec(), tenSecAccum, 5*time.Second, sdk.NewDec(2)),
		"accumulator = 10*OneSec, t=3s. 0 base accum": testCaseFromDeltas(
			sdk.ZeroDec(), tenSecAccum, 3*time.Second, ThreePlusOneThird),
		"three asset. accumulator = 10*OneSec, t=3s. 0 base accum": testTapCaseFromDeltas(
			sdk.ZeroDec(), tenSecAccum, 3*time.Second, ThreePlusOneThird),
		"accumulator = 10*OneSec, t=100s. 0 base accum": testCaseFromDeltas(
			sdk.ZeroDec(), tenSecAccum, 100*time.Second, sdk.NewDecWithPrec(1, 1)),
		"three asset. accumulator = 10*OneSec, t=100s. 0 base accum": testTapCaseFromDeltas(
			sdk.ZeroDec(), tenSecAccum, 100*time.Second, sdk.NewDecWithPrec(1, 1)),

		// test that base accum has no impact
		"accumulator = 10*OneSec, t=5s. 10 base accum": testCaseFromDeltas(
			sdk.NewDec(10), tenSecAccum, 5*time.Second, sdk.NewDec(2)),
		"three asset. accumulator = 10*OneSec, t=5s. 10 base accum": testTapCaseFromDeltas(
			sdk.NewDec(10), tenSecAccum, 5*time.Second, sdk.NewDec(2)),
		"accumulator = 10*OneSec, t=3s. 10*second base accum": testCaseFromDeltas(
			tenSecAccum, tenSecAccum, 3*time.Second, ThreePlusOneThird),
		"three asset. accumulator = 10*OneSec, t=3s. 10*second base accum": testTapCaseFromDeltas(
			tenSecAccum, tenSecAccum, 3*time.Second, ThreePlusOneThird),
		"accumulator = 10*OneSec, t=100s. .1*second base accum": testCaseFromDeltas(
			pointOneAccum, tenSecAccum, 100*time.Second, sdk.NewDecWithPrec(1, 1)),
		"three asset. accumulator = 10*OneSec, t=100s. .1*second base accum": testTapCaseFromDeltas(
			pointOneAccum, tenSecAccum, 100*time.Second, sdk.NewDecWithPrec(1, 1)),

		"accumulator = 10*OneSec, t=100s. 0 base accum (asset 1)":              testCaseFromDeltasAsset1(sdk.ZeroDec(), OneSec.MulInt64(10), 100*time.Second, sdk.NewDecWithPrec(1, 1)),
		"three asset. accumulator = 10*OneSec, t=100s. 0 base accum (asset 1)": testTapCaseFromDeltasAsset1(sdk.ZeroDec(), OneSec.MulInt64(10), 100*time.Second, sdk.NewDecWithPrec(1, 1)),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for i, startRec := range test.startRecord {
				actualTwap, err := twap.ComputeArithmeticTwap(startRec, test.endRecord[i], test.quoteAsset[i])
				require.Equal(t, test.expTwap, actualTwap)
				require.NoError(t, err)
			}
		})
	}
}

// This tests the behavior of computeArithmeticTwap, around error returning
// when there has been an intermediate spot price error.
func TestComputeArithmeticTwapWithSpotPriceError(t *testing.T) {
	newOneSidedRecordWErrorTime := func(time time.Time, accum sdk.Dec, useP0 bool, errTime time.Time) []types.TwapRecord {
		record := newOneSidedRecord(time, accum, useP0)
		record[0].LastErrorTime = errTime
		return record
	}

	newTapOneSidedRecordWErrorTime := func(time time.Time, accum sdk.Dec, useP0 bool, errTime time.Time) []types.TwapRecord {
		records := newTapOneSidedRecord(time, accum, useP0)
		for i := range records {
			records[i].LastErrorTime = errTime
		}
		return records
	}

	tests := map[string]computeArithmeticTwapTestCase{
		// should error, since end time may have been used to interpolate this value
		"errAtEndTime from end record": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecordWErrorTime(tPlusOne, OneSec, true, tPlusOne),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.OneDec(),
			expErr:      true,
		},
		"three asset, errAtEndTime from end record": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecordWErrorTime(tPlusOne, OneSec, true, tPlusOne),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.OneDec(),
			expErr:      true,
		},
		// should error, since start time may have been used to interpolate this value
		"err at StartTime exactly": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecordWErrorTime(tPlusOne, OneSec, true, baseTime),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.OneDec(),
			expErr:      true,
		},
		"three asset, err at StartTime exactly": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecordWErrorTime(tPlusOne, OneSec, true, baseTime),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.OneDec(),
			expErr:      true,
		},
		"err before StartTime": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecordWErrorTime(tPlusOne, OneSec, true, tMinOne),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.OneDec(),
			expErr:      false,
		},
		"three asset, err before StartTime": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecordWErrorTime(tPlusOne, OneSec, true, tMinOne),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.OneDec(),
			expErr:      false,
		},
		// Should not happen, but if it did would error
		"err after EndTime": {
			startRecord: newOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newOneSidedRecordWErrorTime(tPlusOne, OneSec.MulInt64(2), true, baseTime.Add(20*time.Second)),
			quoteAsset:  []string{denom0},
			expTwap:     sdk.OneDec().MulInt64(2),
			expErr:      true,
		},
		"three asset, err after EndTime": {
			startRecord: newTapOneSidedRecord(baseTime, sdk.ZeroDec(), true),
			endRecord:   newTapOneSidedRecordWErrorTime(tPlusOne, OneSec.MulInt64(2), true, baseTime.Add(20*time.Second)),
			quoteAsset:  []string{denom0, denom0, denom1},
			expTwap:     sdk.OneDec().MulInt64(2),
			expErr:      true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for i, startRec := range test.startRecord {
				actualTwap, err := twap.ComputeArithmeticTwap(startRec, test.endRecord[i], test.quoteAsset[i])
				require.Equal(t, test.expTwap, actualTwap)
				osmoassert.ConditionalError(t, test.expErr, err)
			}
		})
	}
}

// TestPruneRecords tests that twap records earlier than
// current block time - RecordHistoryKeepPeriod are pruned from the store
// while keeping the newest record before the above time threshold.
// Such record is kept for each pool.
func (s *TestSuite) TestPruneRecords() {
	recordHistoryKeepPeriod := s.twapkeeper.RecordHistoryKeepPeriod(s.Ctx)

	pool1OlderMin2MsRecord, // deleted
		pool2OlderMin1MsRecord,  // deleted
		pool3OlderBaseRecord,    // kept as newest under keep period
		pool4OlderPlus1Record := // kept as newest under keep period
		s.createTestRecordsFromTime(baseTime.Add(2 * -recordHistoryKeepPeriod))

	pool1Min2MsRecord, // kept as newest under keep period
		pool2Min1MsRecord,  // kept as newest under keep period
		pool3BaseRecord,    // kept as it is at the keep period boundary
		pool4Plus1Record := // kept as it is above the keep period boundary
		s.createTestRecordsFromTime(baseTime.Add(-recordHistoryKeepPeriod))

	// non-ascending insertion order.
	recordsToPreSet := []types.TwapRecord{
		pool2OlderMin1MsRecord,
		pool4Plus1Record,
		pool4OlderPlus1Record,
		pool3OlderBaseRecord,
		pool2Min1MsRecord,
		pool3BaseRecord,
		pool1Min2MsRecord,
		pool1OlderMin2MsRecord,
	}

	// tMin2Record is before the threshold and is pruned away.
	// tmin1Record is the newest record before current block time - record history keep period.
	// All other records happen after the threshold and are kept.
	expectedKeptRecords := []types.TwapRecord{
		pool3OlderBaseRecord,
		pool4OlderPlus1Record,
		pool1Min2MsRecord,
		pool2Min1MsRecord,
		pool3BaseRecord,
		pool4Plus1Record,
	}
	s.SetupTest()
	s.preSetRecords(recordsToPreSet)

	ctx := s.Ctx
	twapKeeper := s.twapkeeper

	ctx = ctx.WithBlockTime(baseTime)

	err := twapKeeper.PruneRecords(ctx)
	s.Require().NoError(err)

	s.validateExpectedRecords(expectedKeptRecords)
}
