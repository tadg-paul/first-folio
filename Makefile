INSTALL_DIR := $(HOME)/.local/bin
SCRIPTS := org-play-to-pdf org-play-to-markdown
PROJECT_DIR := $(shell cd "$(dir $(lastword $(MAKEFILE_LIST)))" && pwd)

.PHONY: install uninstall test test-one-off

install:
	@for script in $(SCRIPTS); do \
		ln -sf "$(PROJECT_DIR)/$$script" "$(INSTALL_DIR)/$$script"; \
		echo "Linked $$script -> $(INSTALL_DIR)/$$script"; \
	done

uninstall:
	@for script in $(SCRIPTS); do \
		rm -f "$(INSTALL_DIR)/$$script"; \
		echo "Removed $(INSTALL_DIR)/$$script"; \
	done

test:
	@bash tests/regression/test_org_play_export.sh

ifdef ISSUE
test-one-off:
	@bash tests/one_off/test_*$(ISSUE)*.sh
else
test-one-off:
	@echo "No one-off tests defined yet"
endif
