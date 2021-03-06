package commands_test

import (
	"context"
	"testing"

	"github.com/filecoin-project/specs-actors/actors/abi"

	"github.com/filecoin-project/go-filecoin/internal/pkg/types"

	"github.com/filecoin-project/go-filecoin/internal/app/go-filecoin/node/test"

	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-filecoin/fixtures/fortest"
	tf "github.com/filecoin-project/go-filecoin/internal/pkg/testhelpers/testflags"
	"github.com/stretchr/testify/assert"
)

func TestListAsks(t *testing.T) {
	tf.IntegrationTest(t)
	ctx := context.TODO()

	seed, cfg, fakeClk, chainClk := test.CreateBootstrapSetup(t)
	n := test.CreateBootstrapMiner(ctx, t, seed, chainClk, cfg)

	minerDaemon, apiDone := test.RunNodeAPI(ctx, n, t)
	defer apiDone()

	minerDaemon.RunSuccess(ctx, "miner", "set-price", "20", "10")

	test.RequireMineOnce(ctx, t, fakeClk, n)

	var asks []*storagemarket.SignedStorageAsk
	minerDaemon.RunMarshaledJSON(ctx, &asks, "client", "list-asks")
	assert.Len(t, asks, 1)
	ask := asks[0].Ask
	assert.Equal(t, fortest.TestMiners[0], ask.Miner)
	assert.Equal(t, uint64(1), ask.SeqNo)
	assert.Equal(t, types.NewAttoFILFromFIL(20), ask.Price)
	assert.Equal(t, abi.ChainEpoch(10), ask.Expiry)
}
