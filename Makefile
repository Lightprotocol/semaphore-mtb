tiny-wasm: main.go
	tinygo build -o ./tiny.wasm -target wasi -gc=leaking -opt=2 -panic=trap -no-debug -scheduler asyncify ./main.go

tiny-run: main.go
	tinygo run main.go
