module example.com/kafka_ingester

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/evanphx/json-patch v0.5.2
	github.com/segmentio/kafka-go v0.4.16
	gopkg.in/alecthomas/kingpin.v2 v2.2.6

	github.com/maxott/magda-cli v0.0.0

)
replace "github.com/maxott/magda-cli" => "../.."