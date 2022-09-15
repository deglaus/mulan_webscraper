######################################
##### DOCUMENTATION REQUIREMENTS #####
######################################

# Python: 		Sphinx-doc 
# REQUIREMENTS:
#	pip install sphinx
#	pip install autodocsumm

# Go: 			godoc
# Installation:
#	- https://www.bswen.com/2020/07/How-to-install-godoc-tool-for-golang.html
# How to open documentation after executing "make doc":
# 	- Go to: http://localhost:6060/pkg/gosecondhand/ 
# If you have problems with godoc installation then enter this in your terminal:
# 	1. export GOPATH=$HOME/go
#	2. export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

####################
##### COMMANDS #####
####################

# Initialises the enabled scrapers in Main.go
# EXAMPLE:
#	make scrape arg="mumin"
scrape:
	go run ./src/Main.go $(arg)

# 1. Starts API 
# 2. Starts React.js App in browser on: http://localhost:3000
ui:
	npm start --prefix src/server & npm start --prefix src/ui && fg


# Drops all tables in the `test` database
drop:
	go run ./src/droptables.go

# Start and stop XAMPP
startX:
	python3 ./scripts/XamppStarter.py

stopX:
	python3 ./scripts/XamppStopper.py


# Installs ui-related packages
install:
	npm install --prefix ./src/ui; \
	npm install --prefix ./src/server express --save; \
	npm install --prefix -/src/ui react-router-dom@5.2.0

# Generates function and module documentation 
doc:
	sphinx-build -b html src/ doc/Python; \
	godoc -http=localhost:6060

# Deletes junk files from documentation generation 
clean:
	find doc/Python/ -type f -maxdepth 1 -delete
	find src/targets/__pycache__ -type f -maxdepth 1 -delete
	rm -f doc/Python/.doctrees/index.doctree
	rm -f doc/Python/.doctrees/environment.pickle

# Runs all tests
# NO EXTRA REQUIREMENETS! 
# The -v flag in go test module is used to show 
# the live status of every test that executing in the terminal
tests:
	go test -v ./tests/utils_test.go; \
	go test -v ./tests/scrape_test.go

#################
##### OTHER #####
#################

.PHONY: scrape ui doc clean tests install
