package sensor

import "github.com/example/components/contracts"

// SensorService implements contracts.TemperatureContractHandler.
type SensorService struct{}

func New() *SensorService { return &SensorService{} }

func (s *SensorService) HandleGetReading(req contracts.GetReadingRequest) (contracts.GetReadingResponse, error) {
	return contracts.GetReadingResponse{Value: 22.5, Unit: "celsius"}, nil
}

func (s *SensorService) TriggerReadingUpdate() contracts.ReadingUpdateEvent {
	return contracts.ReadingUpdateEvent{SensorId: "sensor-1", Value: 22.5, Unit: "celsius", Timestamp: 1234567890}
}
