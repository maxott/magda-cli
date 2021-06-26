// Program to create, update & delete aspect schemas in Magda
package cmd

import (
	"github.com/maxott/magda-cli/pkg/adapter"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("magda-cli", "Managing records & schemas in Magda.")

	host        = app.Flag("host", "DNS name/IP of Magda host [MAGDA_HOST]").Short('H').Envar("MAGDA_HOST").String()
	tenantID    = app.Flag("tenantID", "Tenant ID [MAGDA_TENANT_ID]").Envar("MAGDA_TENANT_ID").String()
	authID      = app.Flag("authID", "Authorization Key ID [MAGDA_AUTH_ID]").Envar("MAGDA_AUTH_ID").String()
	authKey     = app.Flag("authKey", "Authorization Key [MAGDA_AUTH_KEY]").Envar("MAGDA_AUTH_KEY").String()
	useTLS      = app.Flag("useTLS", "Use https").Default("false").Bool()
	skipGateway = app.Flag("skipGateway", "Skip gateway server and call registry server directly [MAGDA_SKIP_GATEWAY]]").
			Default("false").Envar("MAGDA_SKIP_GATEWAY").Bool()
)

func App() *kingpin.Application {
	return app
}

func Adapter() *adapter.Adapter {
	adapter := adapter.RestAdapter(adapter.ConnectionCtxt{
		Host: *host, TenantID: *tenantID, AuthID: *authID, AuthKey: *authKey, UseTLS: *useTLS, SkipGateway: *skipGateway,
	})
	return &adapter
}
