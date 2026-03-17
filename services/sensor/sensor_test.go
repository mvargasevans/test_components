package sensor_test

import (
	"testing"

	"github.com/example/components/contracttests"
	"github.com/example/components/contracts"
	"github.com/example/components/services/sensor"
)

func TestSensorProviderContract(t *testing.T) {
	contracttests.RunTemperatureContractProviderTests(t, func() contracts.TemperatureContractHandler {
		return sensor.New()
	})
}
