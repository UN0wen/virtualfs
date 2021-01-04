JSCMD=yarn
JSBUILD=$(JSCMD) build
all: test build
build: frontend backend
frontend:
		cd web; $(JSBUILD)

backend:
		$(MAKE) -C ./server build

test: 
		$(MAKE) -C ./server test
clean: 
		$(MAKE) -C ./server clean
run: frontend
		$(MAKE) -C ./server run