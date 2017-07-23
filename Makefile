BIN = cryptz

install:
	go install .

container:
	docker build -t cryptz:latest .

clean:
	rm -f $(BIN)
	rm -f $(GOPATH)/bin/$(BIN)
	find $(GOPATH)/pkg -maxdepth 2 -type d -name "cryptz" -exec rm -rf {} \;

