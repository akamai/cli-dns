package edgegrid

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/urfave/cli"
)

func expandHome(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[1:])
		}
		// If error getting home, just return path as-is
	}
	return path
}

// Loads Edgegrid config from env, file or CLI flags
func GetEdgegridConfig(c *cli.Context) (*edgegrid.Config, error) {
	edgercPath := expandHome(GetEdgercPath(c))
	section := GetEdgercSection(c)
	//accountKey := c.GlobalString("accountkey")

	edgercOps := []edgegrid.Option{
		edgegrid.WithEnv(true),
		edgegrid.WithFile(edgercPath),
		edgegrid.WithSection(section),
	}
	config, err := edgegrid.New(edgercOps...)
	if err != nil {
		return nil, err
	}

	if c.IsSet("accountkey") {
		config.AccountKey = c.GlobalString("accountkey")
	}

	return config, nil
}

// Get path to .edgerc file
func GetEdgercPath(c *cli.Context) string {
	if path := c.GlobalString("edgerc"); path != "" {
		return path
	}
	return edgegrid.DefaultConfigFile
}

// Get .edgerc section
func GetEdgercSection(c *cli.Context) string {
	if section := c.GlobalString("section"); section != "" {
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
		conf.RetryWaitMax = time.Duration(sec) * time.Second
	}

	if excluded, ok := os.LookupEnv("AKAMAI_RETRY_EXCLUDED_ENDPOINTS"); ok {
		conf.ExcludedEndpoints = strings.Split(excluded, ",")
	} else {
		conf.ExcludedEndpoints = []string{"/identity-management/v3/user-admin/ui-identities/*"}
	}

	return &conf, nil
}
