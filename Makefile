GOCMD=go
GOBUILD=$(GOCMD) build
TARGET_FILE=*.go
OUTPUT_GO_FILE=main

build:
	$(GOBUILD) -o $(OUTPUT_GO_FILE) $(TARGET_FILE)

test: build
	./test.sh

clean:
	rm -f main *.o *~ tmp*

.PHONY: test clean