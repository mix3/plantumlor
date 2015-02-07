plantumlor: bindata_assetfs.go controller/controller.go main.go middleware/logger.go middleware/recover.go plantuml/plantuml.go
	go build

bindata_assetfs.go: assets/css/style.css assets/index.html assets/js/app.js
	go-bindata-assetfs -prefix=assets assets/...

clean:
	rm -rf plantumlor
