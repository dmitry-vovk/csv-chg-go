export SERVICE_NAME=go-challenge-tst
export NETWORK=$(SERVICE_NAME)-net
MINIMUM_TEST_COVERAGE=80

default: clean fmt test check-coverage clean

clean:
	docker-compose -f ./.divido/development/compose.yml rm --force
	docker network rm $(NETWORK) || true

test:
	docker network create $(NETWORK) || true
	docker-compose -f ./.divido/development/compose.yml up --build --force-recreate --abort-on-container-exit

check-coverage:
	@docker run --rm \
		-e MINIMUM_COVERAGE=${MINIMUM_TEST_COVERAGE} \
		-e README_FILE=/app/README.md \
		-e COVERAGE_FILE=/app/coverage.out \
		-v $(PWD):/app \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		divido/devtools:go-coverage
	git diff --exit-code README.md

fmt:
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		divido/gobuild:1.14-alpine3.11 sh -c "go mod tidy; go fmt ./..."
