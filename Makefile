
PRJ_PATH ?= $(PWD)

benchmark: $(PRJ_PATH)/bench
	go test -v $(PRJ_PATH)/bench -bench=$(bench) -run=None -benchmem

go.test: $(PRJ_PATH)/test
	go test -v $(PRJ_PATH)/test -count=1