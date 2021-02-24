package xyzpb

import (
	context "context"
	"strconv"
	"time"

	"github.com/pedidopago/ms-grpcd/pkg/grpcd"
	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app/config"
	"github.com/rs/zerolog/log"
)

// Client returns a service client
func Client(cl grpcd.Client) XYZServiceClient {
	conn, err := cl.Conn("xyzservice")
	if err != nil {
		log.Error().Err(err).Msg("helper: failed to obtain client")
		return nil
	}
	return NewXYZServiceClient(conn)
}

// ServerPublish publishes the address of the service every 2 minutes
func ServerPublish(ctx context.Context, cfg config.Config) {
	sv, err := grpcd.NewServer(cfg.GRPCDAddress())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to grpcd (publish intent)")
	}
	addr := cfg.Host() + ":" + strconv.Itoa(cfg.Port())
	if cfg.GRPCTLS() != nil && cfg.GRPCTLS().IsAlts() {
		sv.PublishAltsInterval(ctx, time.Minute*2, "xyzservice", "xyzservice", addr)
	} else {
		sv.PublishInterval(ctx, time.Minute*2, "xyzservice", "xyzservice", addr)
	}
}
