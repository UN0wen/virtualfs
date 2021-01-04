JSCMD=yarn
JSBUILD=$(JSCMD) build
JSINSTALL=$(JSCMD) install
all: test build
build: frontend backend
frontend:
		cd app; $(JSINSTALL); $(JSBUILD)

vendor:
		$(MAKE) -C ./server get
backend: vendor
		$(MAKE) -C ./server build

test: 
		$(MAKE) -C ./server test
clean: 
		$(MAKE) -C ./server clean
run: frontend vendor
		$(MAKE) -C ./server run