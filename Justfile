# Run all tests.
test:
  find . -name go.mod | grep -v /_ | xargs -n1 dirname | xargs -n1 -I{} sh -c 'cd {} && go test -v ./...'

# Run all tests with XML reporting.
test-xml:
  find . -name go.mod | grep -v /_ | xargs -n1 dirname | xargs -n1 -I{} sh -c 'cd {} && go test -v 2>&1 ./...| go-junit-report -set-exit-code > report.xml'

# Lint everything.
lint:
  find . -name go.mod | grep -v /_ | xargs -n1 dirname | xargs -n1 -I{} sh -c 'cd {} && golangci-lint run && staticcheck ./...'
