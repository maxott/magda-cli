
build:
	go mod tidy
	go build .

create-schemas: build
	./magda-cli schema create --id cse-order --name cse-order --schema-file example/schema/order.json
	./magda-cli schema create --id cse-service --name cse-service --schemaFile example/schema/service.json

load-services:
	./magda-cli record update --id ffdi --name ffdi -a cse-service -f example/record/ffdi_service.json