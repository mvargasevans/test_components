# test_components

A lightweight, YAML-driven **consumer-driven contract testing** framework for services that communicate over Unix IPC — inspired by [Pact](https://docs.pact.io/).

## What is consumer-driven contract testing?

In a system where multiple services communicate over a protocol, the **consumer** (the service that calls or reads) writes a contract describing exactly what it needs from the **provider** (the service that serves or publishes). The provider then runs a generated test harness against their own service to verify they satisfy that contract.

If the provider's test passes, the consumer is guaranteed the protocol is satisfied — without the two teams needing to coordinate a shared test environment.

## How it works here

Contracts are declared in YAML:

```yaml
name: TemperatureContract
consumer: dashboard
provider: sensor-service
interactions:
  - name: GetReading
    type: request_reply
    request:
      - {name: sensor_id, type: string, example: "sensor-1"}
    response:
      - {name: value, type: float64, example: "22.5"}
      - {name: unit,  type: string,  example: "celsius"}

  - name: ReadingUpdate
    type: publish
    fields:
      - {name: sensor_id, type: string,  example: "sensor-1"}
      - {name: value,     type: float64, example: "22.5"}
      - {name: unit,      type: string,  example: "celsius"}
      - {name: timestamp, type: int,     example: "1234567890"}
```

Running `go generate ./...` produces:

- **`contracts/temperature_gen.go`** — typed Go structs and a `TemperatureContractHandler` interface
- **`contracttests/temperature_contract_gen.go`** — a `RunTemperatureContractProviderTests` function that spins up a real Unix socket, sends messages over it, and validates field types and values against the contract

The provider team implements the interface and wires the test in ~4 lines:

```go
func TestSensorProviderContract(t *testing.T) {
    contracttests.RunTemperatureContractProviderTests(t, func() contracts.TemperatureContractHandler {
        return sensor.New()
    })
}
```

## Configurability

The contract schema supports two interaction types:

| Type | Description |
|---|---|
| `request_reply` | Client sends a request, server responds synchronously |
| `publish` | Server emits an unsolicited event (e.g. a stream update) |

Field types supported: `string`, `int`, `float64`, `bool`.

Each field accepts an optional `example` value. When present, the generated test asserts the provider returns that exact value — not just that the field exists and has the right type. This is what catches a provider that compiles fine but never populates a newly added field.

## Demonstrating the consumer-driven guarantee

1. Consumer edits `contracts/temperature.yaml` — adds a field, changes a type, or adds a new interaction
2. Consumer runs `go generate ./...`
3. Provider runs `go test ./services/sensor/...`
4. Test fails (missing field value) or fails to compile (interface changed) — forcing the provider to update their implementation

## Running

```bash
go generate ./...          # regenerate from YAML
go test -count=1 ./...     # run all contract tests over real sockets
go vet ./...
```
