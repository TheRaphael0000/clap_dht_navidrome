plugin.wasm: *.go *.json *.mod
	tinygo build -o plugin.wasm -target wasip1 -gc=leaking -buildmode=c-shared .

clap_dht.ndp: plugin.wasm
	zip -j clap_dht.ndp manifest.json plugin.wasm


dev: clap_dht.ndp
	cp clap_dht.ndp /opt/navidrome_test/data/plugins -f

all: clap_dht.ndp

clean:
	rm -r plugin.wasm clap_dht.ndp