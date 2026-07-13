.PHONY: help tidy clonator db-start db-stop db-create db-drop test test-e2e start proto appsec-install appsec-uninstall support-install-grpc-tools support-install-github-ssh-key appsec-test

help:
	@echo ""
	@echo ""
	@echo "Available targets:"
	@echo ""
	@echo "  help - Show this help"
	@echo ""
	@echo "  clonator - Start the interactive Clonator CLI to generate a new service"
	@echo ""
	@echo "  db-start - Start PostgreSQL database"
	@echo "  db-stop - Stop PostgreSQL database"
	@echo "  db-create - Create database"
	@echo "  db-drop - Drop database"
	@echo ""
	@echo "  test - Run application tests"
	@echo "  test-e2e - Run End-to-End gRPC tests using grpcurl"
	@echo "  start - Start application"
	@echo "  proto - Generate protocol buffers"
	@echo "  tidy - Run code formatting and linting"
	@echo ""
	@echo "  appsec-test      - Run all application security tests"
	@echo "  appsec-test-code - Run application security tests for code only"
	@echo "  appsec-install   - Install application security tools"
	@echo "  appsec-uninstall - Uninstall application security tools"
	@echo ""
	@echo "  support-install-grpc-tools - Install gRPC tools"
	@echo "  support-install-github-ssh-key - Install GitHub SSH key"
	@echo ""
	
db-start:
	bash bin/database/start.sh
db-stop:
	bash bin/database/stop.sh
db-create:
	bash  bin/database/create.sh
db-drop:
	bash  bin/database/drop.sh

clonator:
	@if [ "$(FILE)" != "" ]; then \
		go run cmd/clonator/*.go --from-file=$(FILE); \
	else \
		go run cmd/clonator/*.go; \
	fi
test:
	bash bin/app/test.sh
test-e2e:
	bash bin/app/test_e2e_grpc.sh
start:
	bash bin/app/start.sh
proto:
	bash bin/app/proto.sh
tidy:
	go mod tidy

appsec-test:
	@SNYK_TOKEN=${SNYK_TOKEN} bash bin/appsec/test.sh
appsec-test-code:
	@SNYK_TOKEN=${SNYK_TOKEN} bash bin/appsec/test-code.sh
appsec-install:
	bash bin/appsec/install.sh
appsec-install-tools-with-hook:
	bash bin/appsec/install-tools-with-hook.sh
appsec-uninstall-tools-with-hook:
	bash bin/appsec/uninstall-tools-with-hook.sh
appsec-install-hook:
	bash bin/appsec/install-hook.sh
appsec-uninstall-hook:
	bash bin/appsec/uninstall-hook.sh
	
support-install-grpc-tools:
	bash bin/support/install-grpc-tools.sh
support-install-github-ssh-key:
	bash bin/support/install-github-ssh-key.sh github-for-ssh-key markitos.es.info@gmail.com
support-set-git-globals:
	git config --global user.email "markitos.es.info@gmail.com"
	git config --global user.name "marco antonio - markitos"
