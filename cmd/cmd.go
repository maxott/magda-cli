// Program to create, update & delete aspect schemas in Magda
package cmd

import (
	"github.com/maxott/magda-cli/pkg/adapter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("magda-cli", "Managing records & schemas in Magda.")

	host        = app.Flag("host", "DNS name/IP of Magda host [MAGDA_HOST]").Short('H').Envar("MAGDA_HOST").String()
	tenantID    = app.Flag("tenant-id", "Tenant ID [MAGDA_TENANT_ID]").Envar("MAGDA_TENANT_ID").String()
	authID      = app.Flag("auth-id", "Authorization Key ID [MAGDA_AUTH_ID]").Envar("MAGDA_AUTH_ID").String()
	authKey     = app.Flag("auth-key", "Authorization Key [MAGDA_AUTH_KEY]").Envar("MAGDA_AUTH_KEY").String()
	useTLS      = app.Flag("use-tls", "Use https").Default("false").Bool()
	skipGateway = app.Flag("skip-gateway", "Skip gateway server and call registry server directly [MAGDA_SKIP_GATEWAY]]").
			Default("false").Envar("MAGDA_SKIP_GATEWAY").Bool()
	jwtSecret = app.Flag("jwt-secret", "Secret used for creating JWT token for inernal comms [MAGDA_JWT_SECRET]").Envar("MAGDA_JWT_SECRET").String()
	jwtUser   = app.Flag("jwt-user-id", "User ID for creating JWT token for inernal comms [MAGDA_JWT_USER_ID]").Envar("MAGDA_JWT_USER_ID").String()
	// verbose   = app.Flag("verbose", "Be chatty [MAGDA_VERBOSE]").Short('v').Envar("MAGDA_VERBOSE").Bool()
)

func App() *kingpin.Application {
	return app
}

func Adapter() *adapter.Adapter {
	jwtToken := createJwtToken()
	adapter := adapter.RestAdapter(adapter.ConnectionCtxt{
		Host: *host, TenantID: *tenantID, AuthID: *authID, AuthKey: *authKey, UseTLS: *useTLS,
		SkipGateway: *skipGateway, JwtToken: jwtToken,
	})
	return &adapter
}

func createJwtToken() string {
	if *skipGateway {
		if jwtSecret == nil || jwtUser == nil {
			log.Fatal("When skipping gateway, 'jwt-secret' and 'jwt-user-id' are also required")
		}
		token, err := adapter.CreateJwtToken(jwtUser, jwtSecret)
		if err != nil {
			log.Fatal("Error while signing JWT token - ", err)
		}
		log.Info("JWT Token - ", token)
		return token
	} else {
		return ""
	}
}
