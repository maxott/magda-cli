// Program to create, update & delete aspect schemas in Magda
package cmd

import (
	"github.com/maxott/magda-cli/pkg/adapter"
	"gopkg.in/alecthomas/kingpin.v2"
	lgrus "github.com/sirupsen/logrus"

	"github.com/maxott/magda-cli/pkg/log"
	"github.com/maxott/magda-cli/pkg/log/logrus"
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

	logger = logrus.NewSimpleLogger(lgrus.WarnLevel)
)

func App() *kingpin.Application {
	return app
}

func Adapter() *adapter.Adapter {
	jwtToken := createJwtToken(Logger())
	adapter := adapter.RestAdapter(adapter.ConnectionCtxt{
		Host: *host, TenantID: *tenantID, AuthID: *authID, AuthKey: *authKey, UseTLS: *useTLS,
		SkipGateway: *skipGateway, JwtToken: jwtToken,
	})
	return &adapter
}

func Logger() log.Logger {
	return logger
}

func SetLogger(l log.Logger)  {
	logger = l
}

func createJwtToken(logger log.Logger) string {
	if *skipGateway {
		if jwtSecret == nil || jwtUser == nil {
			logger.Fatal("When skipping gateway, 'jwt-secret' and 'jwt-user-id' are also required")
		}
		token, err := adapter.CreateJwtToken(jwtUser, jwtSecret)
		if err != nil {
			logger.With("error", err).Fatal("While signing JWT token")
		}
		logger.Debugf("JWT Token - ", token)
		return token
	} else {
		return ""
	}
}
