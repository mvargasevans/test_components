package alerter_test

import (
	"testing"

	"github.com/example/components/contracttests"
	"github.com/example/components/contracts"
	"github.com/example/components/services/alerter"
)

func TestAlerterProviderContract(t *testing.T) {
	contracttests.RunAlertsContractProviderTests(t, func() contracts.AlertsContractHandler {
		return alerter.New()
	})
}
