package cli_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/made-in-block/slash-refund/testutil/network"
)

func TestE2ETestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cfg := network.DefaultConfig()

	// ==== Query tests ====
	// Run query tests. A new network will be created.
	// In order to isolate queries behaviour from transactions behaviour, deposits,
	// deposit pools, unbonding deposits, refund and refund pools will be set inside
	// genesis configuration.
	suite.Run(t, NewE2EQueryTestSuite(cfg))

	// ==== Transaction tests ====
	// Set specific mnemonics for genesis accounts. This is needed to set inside
	// genesis configuration:
	// 1. specific deposits that will be withdrawn, in order to isolate whitdraw tx
	//    tests from deposit tx tests.
	// 2. specific refunds that will be claimed, in order to isolate claim tx tests
	//    from the slashing event.
	cfg.Mnemonics = []string{
		"tiny void main swift patch note shuffle glue amateur acquire walk bulk river toe test master kind minor canal chicken vanish column woman pioneer",
		"marine hat caught work bulb tourist sniff bacon earth home woman kingdom disorder labor elder become giraffe stage plate lawn truly hold eagle sauce",
		"inmate leader marble ready flag wire hurdle still asset kiss category correct hockey ceiling weekend memory tray empower ticket motor sentence best summer supreme",
		"sing wedding notable immense clown thunder stick match chase rack donkey track indicate woman script flee feature wealth truck verb gallery theory bicycle correct"}

	// This parameter must not be less than 4 for the Tx tests to work as expected.
	// In order to isolate behaviours of the different transactions, the following
	// scheme is used.
	// Validator 0 is used in claim tests. Refunds and refund pools are set in genesis.
	// Validator 1 is used in deposit tests. Deposits and deposit pools will be created
	// for this validator during deposit transaction tests.
	// Validator 2 is used in withdraw tests. Deposits and deposit pools are set in
	// genesis.
	// Validator 3 represents a validator left untouched by transactions.
	cfg.NumValidators = 4

	// Run transaction tests. A network will be created. Specific objects will be
	// created in order to isolate each transaction behaviour.
	suite.Run(t, NewE2ETestSuite(cfg))
}
