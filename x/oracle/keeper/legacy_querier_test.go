package keeper

import (
	"bytes"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/oracle/types"
)

func TestLegacyNewLegacyQuerier(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.Ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)
}

func TestLegacyQueryParams(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryParameters}, req)
	require.NoError(t, err)

	var params types.Params
	err = input.Cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.OracleKeeper.GetParams(input.Ctx), params)
}

func TestLegacyQueryExchangeRate(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	rate := sdk.NewDec(1700)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, rate)

	// denom query params
	queryParams := types.NewQueryExchangeRateParams(core.MicroBSDRDenom)
	bz, err := input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryExchangeRate}, req)
	require.NoError(t, err)

	var queriedRate sdk.Dec
	err = input.Cdc.UnmarshalJSON(res, &queriedRate)
	require.NoError(t, err)
	require.Equal(t, rate, queriedRate)
}

func TestLegacyQueryExchangeRates(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	rate := sdk.NewDec(1700)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, rate)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBUSDDenom, rate)

	res, err := querier(input.Ctx, []string{types.QueryExchangeRates}, abci.RequestQuery{})
	require.NoError(t, err)

	var queriedRate sdk.DecCoins
	err2 := input.Cdc.UnmarshalJSON(res, &queriedRate)
	require.NoError(t, err2)
	require.Equal(t, sdk.DecCoins{
		sdk.NewDecCoinFromDec(core.MicroBSDRDenom, rate),
		sdk.NewDecCoinFromDec(core.MicroBUSDDenom, rate),
	}, queriedRate)
}

func TestLegacyQueryActives(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	rate := sdk.NewDec(1700)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, rate)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBKRWDenom, rate)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBUSDDenom, rate)

	res, err := querier(input.Ctx, []string{types.QueryActives}, abci.RequestQuery{})
	require.NoError(t, err)

	targetDenoms := []string{
		core.MicroBKRWDenom,
		core.MicroBSDRDenom,
		core.MicroBUSDDenom,
	}

	var denoms []string
	err2 := input.Cdc.UnmarshalJSON(res, &denoms)
	require.NoError(t, err2)
	require.Equal(t, targetDenoms, denoms)
}

func TestLegacyQueryFeederDelegation(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	input.OracleKeeper.SetFeederDelegation(input.Ctx, ValAddrs[0], Addrs[1])

	queryParams := types.NewQueryFeederDelegationParams(ValAddrs[0])
	bz, err := input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryFeederDelegation}, req)
	require.NoError(t, err)

	var delegate sdk.AccAddress
	input.Cdc.UnmarshalJSON(res, &delegate)
	require.Equal(t, Addrs[1], delegate)
}

func TestLegacyQueryAggregatePrevote(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	prevote1 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[0], 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[0], prevote1)
	prevote2 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[1], 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[1], prevote2)
	prevote3 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[2], 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[2], prevote3)

	// validator 0 address params
	queryParams := types.NewQueryAggregatePrevoteParams(ValAddrs[0])
	bz, err := input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryAggregatePrevote}, req)
	require.NoError(t, err)

	var prevote types.AggregateExchangeRatePrevote
	err = input.Cdc.UnmarshalJSON(res, &prevote)
	require.NoError(t, err)
	require.Equal(t, prevote1, prevote)

	// validator 1 address params
	queryParams = types.NewQueryAggregatePrevoteParams(ValAddrs[1])
	bz, err = input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryAggregatePrevote}, req)
	require.NoError(t, err)

	err = input.Cdc.UnmarshalJSON(res, &prevote)
	require.NoError(t, err)
	require.Equal(t, prevote2, prevote)
}

func TestLegacyQueryAggregatePrevotes(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	prevote1 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[0], 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[0], prevote1)
	prevote2 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[1], 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[1], prevote2)
	prevote3 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[2], 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[2], prevote3)

	expectedPrevotes := []types.AggregateExchangeRatePrevote{prevote1, prevote2, prevote3}
	sort.SliceStable(expectedPrevotes, func(i, j int) bool {
		addr1, _ := sdk.ValAddressFromBech32(expectedPrevotes[i].Voter)
		addr2, _ := sdk.ValAddressFromBech32(expectedPrevotes[j].Voter)
		return bytes.Compare(addr1, addr2) == -1
	})

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryAggregatePrevotes}, req)
	require.NoError(t, err)

	var prevotes []types.AggregateExchangeRatePrevote
	err = input.Cdc.UnmarshalJSON(res, &prevotes)
	require.NoError(t, err)
	require.Equal(t, expectedPrevotes, prevotes)
}

func TestLegacyQueryAggregateVote(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	vote1 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "", ExchangeRate: sdk.OneDec()}}, ValAddrs[0])
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[0], vote1)
	vote2 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "", ExchangeRate: sdk.OneDec()}}, ValAddrs[1])
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[1], vote2)
	vote3 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "", ExchangeRate: sdk.OneDec()}}, ValAddrs[2])
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[2], vote3)

	// validator 0 address params
	queryParams := types.NewQueryAggregateVoteParams(ValAddrs[0])
	bz, err := input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryAggregateVote}, req)
	require.NoError(t, err)

	var vote types.AggregateExchangeRateVote
	err = input.Cdc.UnmarshalJSON(res, &vote)
	require.NoError(t, err)
	require.Equal(t, vote1, vote)

	// validator 1 address params
	queryParams = types.NewQueryAggregateVoteParams(ValAddrs[1])
	bz, err = input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryAggregateVote}, req)
	require.NoError(t, err)

	err = input.Cdc.UnmarshalJSON(res, &vote)
	require.NoError(t, err)
	require.Equal(t, vote2, vote)
}

func TestLegacyQueryAggregateVotes(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	vote1 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "", ExchangeRate: sdk.OneDec()}}, ValAddrs[0])
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[0], vote1)
	vote2 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "", ExchangeRate: sdk.OneDec()}}, ValAddrs[1])
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[1], vote2)
	vote3 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "", ExchangeRate: sdk.OneDec()}}, ValAddrs[2])
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[2], vote3)

	expectedVotes := []types.AggregateExchangeRateVote{vote1, vote2, vote3}
	sort.SliceStable(expectedVotes, func(i, j int) bool {
		addr1, _ := sdk.ValAddressFromBech32(expectedVotes[i].Voter)
		addr2, _ := sdk.ValAddressFromBech32(expectedVotes[j].Voter)
		return bytes.Compare(addr1, addr2) == -1
	})

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryAggregateVotes}, req)
	require.NoError(t, err)

	var votes []types.AggregateExchangeRateVote
	err = input.Cdc.UnmarshalJSON(res, &votes)
	require.NoError(t, err)
	require.Equal(t, expectedVotes, votes)
}

func TestLegacyQueryVoteTargets(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	// clear tobin taxes
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)

	voteTargets := []string{"denom", "denom2", "denom3"}
	for _, target := range voteTargets {
		input.OracleKeeper.SetTobinTax(input.Ctx, target, sdk.OneDec())
	}

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryVoteTargets}, req)
	require.NoError(t, err)

	var voteTargetsRes []string
	err2 := input.Cdc.UnmarshalJSON(res, &voteTargetsRes)
	require.NoError(t, err2)
	require.Equal(t, voteTargets, voteTargetsRes)
}

func TestLegacyQueryTobinTaxes(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	// clear tobin taxes
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)

	tobinTaxes := types.DenomList{{
		Name:     core.MicroBKRWDenom,
		TobinTax: sdk.OneDec(),
	}, {
		Name:     core.MicroBSDRDenom,
		TobinTax: sdk.NewDecWithPrec(123, 2),
	}}
	for _, item := range tobinTaxes {
		input.OracleKeeper.SetTobinTax(input.Ctx, item.Name, item.TobinTax)
	}

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryTobinTaxes}, req)
	require.NoError(t, err)

	var tobinTaxesRes types.DenomList
	err2 := input.Cdc.UnmarshalJSON(res, &tobinTaxesRes)
	require.NoError(t, err2)
	require.Equal(t, tobinTaxes, tobinTaxesRes)
}

func TestLegacyQueryTobinTax(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.OracleKeeper, input.Cdc)

	denom := types.Denom{Name: core.MicroBKRWDenom, TobinTax: sdk.OneDec()}
	input.OracleKeeper.SetTobinTax(input.Ctx, denom.Name, denom.TobinTax)

	queryParams := types.NewQueryTobinTaxParams(core.MicroBKRWDenom)
	bz, err := input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryTobinTax}, req)
	require.NoError(t, err)

	var tobinTaxRes sdk.Dec
	input.Cdc.UnmarshalJSON(res, &tobinTaxRes)
	require.Equal(t, denom.TobinTax, tobinTaxRes)
}
