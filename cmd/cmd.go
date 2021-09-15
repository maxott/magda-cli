// Program to create, update & delete aspect schemas in Magda
package cmd

import (
	"os"
	"github.com/maxott/magda-cli/pkg/adapter"
	"gopkg.in/alecthomas/kingpin.v2"

	log "go.uber.org/zap"
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

	useYaml        = app.Flag("use-yaml", "Use and assume data formated in YAML [MAGDA_USE_YAML]").Short('y').Envar("MAGDA_USE_YAML").Bool()

	logger *log.Logger
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

func Logger() *log.Logger {
	return logger
}

func SetLogger(l *log.Logger)  {
	logger = l
}

func createJwtToken(logger *log.Logger) string {
	if *skipGateway {
		if jwtSecret == nil || jwtUser == nil {
			logger.Fatal("When skipping gateway, 'jwt-secret' and 'jwt-user-id' are also required")
		}
		token, err := adapter.CreateJwtToken(jwtUser, jwtSecret)
		if err != nil {
			logger.Fatal("While signing JWT token", log.Error(err))
		}
		logger.Debug("JWT Token", log.String("token", token))
		return token
	} else {
		return ""
	}
}

func loadObjFromFile(fileName string) map[string]interface{} {
	if fileName != "-" {
		if s, err := os.Stat(fileName); os.IsNotExist(err) {
			App().Fatalf("file '%s' does not exist", fileName)
		} else if err != nil {
			App().Fatalf("failed to verify existence of file '%s' - %s", fileName, err)
		} else {
			if s.IsDir() {
				App().Fatalf("path '%s' is not a file", fileName)
			}
		}
	}
	adata, err := adapter.LoadPayloadFromFile(fileName, *useYaml)
	if err != nil {
		App().Fatalf("failed to load '%s' - %s", fileName, err)
	}
	obj, err := adata.AsObject()
	if err != nil {
		App().Fatalf("failed to verify '%s' - %s", fileName, err)
	}
	return obj
}

func loadObjFromStdin() map[string]interface{} {
	adata, err := adapter.LoadPayloadFromStdin(*useYaml)
	if err != nil {
		App().Fatalf("failed to load data from stdin - %s", err)
	}
	obj, err := adata.AsObject()
	if err != nil {
		App().Fatalf("failed to verify data from stdin - %s", err)
	}
	return obj
}
