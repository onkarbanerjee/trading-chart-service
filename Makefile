prepare:
	curl https://pre-commit.com/install-local.py | python3 -
	pre-commit install

generate:
	echo "Not implemented"

migrate:
	echo "Not implemented"

test:
	go test -race ./...

test_and_coverage:
	go test -race -covermode=atomic ./...
