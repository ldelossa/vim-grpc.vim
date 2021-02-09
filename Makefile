.PHONY: protoc
protoc:
	protoc --proto_path=./proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
		./proto/*.proto \
		./proto/env/*.proto \
        ./proto/commands/*.proto

.PHONY: test-env
test-env:
	ln -s ~/git/go/vim-grpc.vim ~/.local/share/nvim/plugged  	

.PHONY: test-env-down
test-env-down:
	rm -rf ~/.local/share/nvim/plugged
