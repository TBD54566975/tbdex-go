# Run all tests.
test:
  @echo "Running tests..."
  @go test -cover ./...

# Run all tests with XML reporting.
test-xml:
  find . -name go.mod | grep -v /_ | xargs -n1 dirname | xargs -n1 -I{} sh -c 'cd {} && go test -v 2>&1 ./...| go-junit-report -set-exit-code > report.xml'

lint:
  @echo "Running linter..."
  @golangci-lint run

# Copies JSON schemas from the tbdex submodule repo into the protocol dir.
schemas:
  @git submodule update --init --recursive
  @cp -r spec/hosted/json-schemas tbdex/
