# prints available just commands. all you need to do is run `just`
_help:
  @just -l

# Run all tests.
test:
  @echo "Running tests..."
  @go test -cover ./...

# Run all tests with XML reporting.
test-xml:
  @echo "Running tests with XML reporting..."
  @find . -name go.mod | grep -v /_ | xargs -n1 dirname | xargs -n1 -I{} sh -c 'cd {} && go test -v 2>&1 ./... | go-junit-report -set-exit-code > report.xml'
  @echo "Test results can be found in report.xml"

lint:
  @echo "Running linter..."
  @golangci-lint run

# Copies JSON schemas from the tbdex submodule repo into the validator dir.
schemas:
  @git submodule update --init --recursive
  @cp -r spec/hosted/json-schemas tbdex/validator/
