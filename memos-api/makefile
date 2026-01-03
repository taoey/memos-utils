# 参考自：https://gist.github.com/MarkWarneke/2e26d7caef237042e9374ebf564517ad
# Makefile for python code
# 
# > make help
#
# The following commands can be used.
#
# init:  sets up environment and installs requirements
# install:  Installs development requirments
# format:  Formats the code with autopep8
# lint:  Runs flake8 on src, exit if critical rules are broken
# clean:  Remove build and cache files
# env:  Source venv and environment files for testing
# leave:  Cleanup and deactivate venv
# test:  Run pytest
# run:  Executes the logic

VENV_PATH='env/bin/activate'
ENVIRONMENT_VARIABLE_FILE='.env'
DOCKER_NAME=$DOCKER_NAME
DOCKER_TAG=$DOCKER_TAG

define find.functions
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
endef

help:
	@echo 'The following commands can be used.'
	@echo ''
	$(call find.functions)


init: ## sets up environment and installs requirements
init:
	pip3 install -r requirements.txt
	pip3 install yapf
	pip3 install pytest
	pip3 install pipreqs
	pip3 install autoflake


clean: ## remove build and cache files
clean:
	rm -rf *.egg-info
	rm -rf build
	rm -rf dist
	rm -rf .pytest_cache
	# Remove all pycache
	find . | grep -E '(__pycache__|\.pyc|\.pyo$\)' | xargs rm -rf