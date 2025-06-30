package edgegrid

import (
	"context"
	"fmt"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/urfave/cli"
)

type ctxKey string

var sessionKey ctxKey = "session"

func InitializeSession(c *cli.Context) (session.Session, error) {
	edgerc, err := GetEdgegridConfig(c)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve zone configuration: %s", err)
	}

	retryConfig, err := getRetryConfig()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve retry configuration: %w", err)
	}

	options := []session.Option{
		session.WithSigner(edgerc),
		session.WithHTTPTracing(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED") == "true"),
		session.WithLog(log.Default()),
	}

	if retryConfig != nil {
		options = append(options, session.WithRetries(*retryConfig))
	}

	sess, err := session.New(options...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize edgegrid session: %s", err)
	}

	return sess, nil
}

func WithSession(ctx context.Context, sess session.Session) context.Context {
	return context.WithValue(ctx, sessionKey, sess)
}

func GetSession(ctx context.Context) session.Session {
	sess, ok := ctx.Value(sessionKey).(session.Session)
	if !ok {
		panic("edgegrid session not found in context")
	}
	return sess
}
