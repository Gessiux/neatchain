BUILD_FLAGS = -tags "$(BUILD_TAGS)" -ldflags "

build:
	@ echo "start building......"
	@ go build -o $(GOPATH)/bin/neatchain ./cmd/neatchain/
	@ echo "Done building."
#.PHONY: neatchain
neatchain:
	@ echo "start building......"
	@ go build -o $(GOPATH)/bin/neatchain ./cmd/neatchain/
	@ echo "Done building."
	@ echo "Run neatchain to launch neatchain network."

install:
	@ echo "start install......"
	@ go install -mod=readonly $(BUILD_FLAGS) ./cmd/neatchain
	@ echo "install success......"