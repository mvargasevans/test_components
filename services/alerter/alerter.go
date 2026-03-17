package alerter

import (
	"fmt"
	"sync"

	"github.com/example/components/contracts"
)

// AlerterService implements contracts.AlertsContractHandler.
type AlerterService struct {
	mu sync.Mutex
	n  int
}

func New() *AlerterService { return &AlerterService{} }

func (a *AlerterService) HandleRaiseAlert(req contracts.RaiseAlertRequest) (contracts.RaiseAlertResponse, error) {
	a.mu.Lock()
	a.n++
	id := fmt.Sprintf("alert-%d", a.n)
	a.mu.Unlock()
	return contracts.RaiseAlertResponse{Id: id, Ok: true}, nil
}

func (a *AlerterService) TriggerAlertFired() contracts.AlertFiredEvent {
	a.mu.Lock()
	a.n++
	id := fmt.Sprintf("alert-%d", a.n)
	a.mu.Unlock()
	return contracts.AlertFiredEvent{Id: id, Level: "warn", Message: "threshold exceeded"}
}
