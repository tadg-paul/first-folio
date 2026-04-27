INSTALL_DIR := $(HOME)/.local/bin
PROJECT_DIR := $(shell cd "$(dir $(lastword $(MAKEFILE_LIST)))" && pwd)

.PHONY: install uninstall test test-one-off

install:
	@ln -sf "$(PROJECT_DIR)/bin/folio" "$(INSTALL_DIR)/folio"
	@echo "Linked folio -> $(INSTALL_DIR)/folio"

uninstall:
	@rm -f "$(INSTALL_DIR)/folio"
	@echo "Removed $(INSTALL_DIR)/folio"

test:
	@for t in tests/regression/test_*.sh; do \
		echo ""; \
		echo ">>> Running $$t"; \
		bash "$$t" || exit 1; \
	done

ifdef ISSUE
test-one-off:
	@bash tests/one_off/test_*$(ISSUE)*.sh
else
test-one-off:
	@echo "No one-off tests defined yet"
endif
