package edgegrid

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/urfave/cli"
)

// Loads Edgegrid config from env, file or CLI flags
func GetEdgegridConfig(c *cli.Context) (*edgegrid.Config, error) {
	edgercOps := []edgegrid.Option{
		edgegrid.WithEnv(true),
		edgegrid.WithFile(GetEdgercPath(c)),
		edgegrid.WithSection(GetEdgercSection(c)),
	}
	config, err := edgegrid.New(edgercOps...)
	if err != nil {
		return nil, err
	}

	if c.IsSet("accountkey") {
		config.AccountKey = c.String("accountkey")
	}

	return config, nil
}

// Get path to .edgerc file
func GetEdgercPath(c *cli.Context) string {
	if path := c.String("edgerc"); path != "" {
		return path
	}
	return edgegrid.DefaultConfigFile
}

// Get .edgerc section
func GetEdgercSection(c *cli.Context) string {
	if section := c.String("section"); section != "" {
		return section
	}
	return edgegrid.DefaultSection
}

// Build Retry config from env variables
func getRetryConfig() (*session.RetryConfig, error) {
	if disabledStr, ok := os.LookupEnv("AKAMAI_RETRY_DISABLED"); ok {
		disabled, err := strconv.ParseBool(disabledStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AKAMAI_RETRY_DISABLED: %w", err)
		}
		if disabled {
			return nil, nil
		}
	}

	conf := session.NewRetryConfig()

	if maxStr, ok := os.LookupEnv("AKAMAI_RETRY_MAX"); ok {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AKAMAI_RETRY_MAX: %w", err)
		}
		conf.RetryMax = max
	}

	if minStr, ok := os.LookupEnv("AKAMAI_RETRY_WAIT_MIN"); ok {
		sec, err := strconv.Atoi(minStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AKAMAI_RETRY_WAIT_MIN: %w", err)
		}
		conf.RetryWaitMin = time.Duration(sec) * time.Second
	}

	if maxStr, ok := os.LookupEnv("AKAMAI_RETRY_WAIT_MAX"); ok {
		sec, err := strconv.Atoi(maxStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AKAMAI_RETRY_WAIT_MAX: %w", err)
		}
		conf.RetryWaitMin = time.Duration(sec) * time.Second
	}

	if excluded, ok := os.LookupEnv("AKAMAI_RETRY_EXCLUDED_ENDPOINTS"); ok {
		conf.ExcludedEndpoints = strings.Split(excluded, ",")
	} else {
		conf.ExcludedEndpoints = []string{"/identity-management/v3/user-admin/ui-identities/*"}
	}

	return &conf, nil
}
