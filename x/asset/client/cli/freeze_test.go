package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/CoreumFoundation/coreum/app"
	"github.com/CoreumFoundation/coreum/pkg/config"
	"github.com/CoreumFoundation/coreum/testutil/network"
	"github.com/CoreumFoundation/coreum/x/asset/client/cli"
	"github.com/CoreumFoundation/coreum/x/asset/types"
)

func TestFreezeFungibleToken(t *testing.T) {
	requireT := require.New(t)
	networkCfg, err := config.NetworkByChainID(config.ChainIDDev)
	requireT.NoError(err)
	app.ChosenNetwork = networkCfg
	testNetwork := network.New(t)

	// the denom must start from the letter
	symbol := "l" + uuid.NewString()[:4]
	ctx := testNetwork.Validators[0].ClientCtx
	issuer := testNetwork.Validators[0].Address
	denom := types.BuildFungibleTokenDenom(symbol, issuer)

	// Issue token
	args := []string{symbol, testNetwork.Validators[0].Address.String(), "777", `"My Token"`,
		"--features", types.FungibleTokenFeature_freezable.String(), //nolint:nosnakecase
	}
	args = append(args, txValidator1Args(testNetwork)...)
	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdTxIssueFungibleToken(), args)
	requireT.NoError(err)

	// Freeze part of the token
	token := "100" + denom
	args = append([]string{issuer.String(), token, "--output", "json"}, txValidator1Args(testNetwork)...)
	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdTxFreezeFungibleToken(), args)
	requireT.NoError(err)

	var resp types.QueryFrozenBalanceResponse
	buf, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryFungibleTokenFrozenBalance(), []string{issuer.String(), denom, "--output", "json"})
	requireT.NoError(err)
	requireT.NoError(ctx.Codec.UnmarshalJSON(buf.Bytes(), &resp))
	requireT.Equal(token, resp.Balance.String())

	// test pagination
	for i := 0; i < 2; i++ {
		symbol := "l" + uuid.NewString()[:4]
		denom := types.BuildFungibleTokenDenom(symbol, issuer)
		args := []string{symbol, testNetwork.Validators[0].Address.String(), "777", `"My Token"`,
			"--features", types.FungibleTokenFeature_freezable.String(), //nolint:nosnakecase
		}
		args = append(args, txValidator1Args(testNetwork)...)
		_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdTxIssueFungibleToken(), args)
		requireT.NoError(err)

		// Freeze part of the token
		tokens := "100" + denom
		args = append([]string{issuer.String(), tokens, "--output", "json"}, txValidator1Args(testNetwork)...)
		_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdTxFreezeFungibleToken(), args)
		requireT.NoError(err)
	}

	var balancesResp types.QueryFrozenBalancesResponse
	buf, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryFungibleTokenFrozenBalances(), []string{issuer.String(), "--output", "json"})
	requireT.NoError(err)
	requireT.NoError(ctx.Codec.UnmarshalJSON(buf.Bytes(), &balancesResp))
	requireT.Len(balancesResp.Balances, 3)

	buf, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryFungibleTokenFrozenBalances(), []string{issuer.String(), "--output", "json", "--limit", "1"})
	requireT.NoError(err)
	requireT.NoError(ctx.Codec.UnmarshalJSON(buf.Bytes(), &balancesResp))
	requireT.Len(balancesResp.Balances, 1)

	// Unfreeze part of the frozen token
	unfreezeTokens := "75" + denom
	args = append([]string{issuer.String(), unfreezeTokens, "--output", "json"}, txValidator1Args(testNetwork)...)
	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdTxUnfreezeFungibleToken(), args)
	requireT.NoError(err)

	buf, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryFungibleTokenFrozenBalance(), []string{issuer.String(), denom, "--output", "json"})
	requireT.NoError(err)
	requireT.NoError(ctx.Codec.UnmarshalJSON(buf.Bytes(), &resp))

	requireT.Equal("25"+denom, resp.Balance.String())
}