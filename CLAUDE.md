# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Regenerate all *_gen.go files from YAML contracts
go generate ./...

# Run all tests (real Unix socket IPC exercised)
go test -count=1 ./...

# Run a single service's contract tests
go test -count=1 ./services/sensor/...
go test -count=1 ./services/alerter/...

# Vet
go vet ./...
```

## Architecture

This repo implements **YAML-driven consumer-driven contract testing over Unix IPC**. The consumer team writes a YAML contract; the provider team runs a generated test harness that exercises real Unix socket communication.

### Flow

```
contracts/*.yaml
      │
      ▼
cmd/contractgen  ──► contracts/*_gen.go          (message structs + handler interface)
                 ──► contracttests/*_contract_gen.go  (RunXxxProviderTests — real IPC)
                              ▲
                    ipc/ipc.go  (shared Unix socket primitives)
                              │
                    services/*/  ← provider implements interface + wires test
```

### Key concepts

- **`contracts/*.yaml`** — written by the consumer. Defines interactions (`request_reply` or `publish`), typed fields, and **example values**. Example values are used in the generated test harness to assert the provider returns non-zero, meaningful data — not just that a field exists.
- **`cmd/contractgen`** — reads every `*.yaml` in `contracts/`, renders two templates per contract: one for structs/interfaces (`gen_contracts.go`) and one for the IPC test harness (`gen_tests.go`). Helper functions (`pascal`, `goType`, `exampleLiteral`, etc.) live in `gen_contracts.go`.
- **`ipc/ipc.go`** — hand-written transport: `TempSocket`, `Envelope`, `Conn`, `Dial`, `Listen`, `Accept`. All messages are newline-delimited JSON.
- **`contracttests/`** — generated `Run<Name>ProviderTests` functions + hand-written `helpers.go` (`dialWithRetry`, `assertField`, `assertFieldValue`). Helper functions are hand-written here specifically to avoid duplication across generated files in the same package.
- **`services/<name>/`** — provider implementation (implements the generated handler interface) + a ~4-line test that calls `Run<Name>ProviderTests`.

### Wire protocol

Request: `{"interaction": "GetReading", "body": {...}}`
Response: `{"interaction": "GetReading", "ok": true, "body": {...}}`
Publish: `{"interaction": "ReadingUpdate", "body": {...}}`

### Adding a new contract

1. Add a new `contracts/<name>.yaml`
2. Run `go generate ./...`
3. Implement the generated `<Name>Handler` interface in a new `services/<name>/` package
4. Wire `Run<Name>ProviderTests` in `services/<name>/<name>_test.go`
