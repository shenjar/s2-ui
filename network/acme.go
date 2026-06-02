package network

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/shenjar/s2-ui/logger"

	"github.com/caddyserver/certmagic"
)

// acmeObtainTimeout bounds the synchronous, first-time certificate issuance so
// that server startup cannot hang indefinitely if port 80 is unreachable or the
// ACME server is slow. It only affects the initial obtain; background renewal is
// driven by certmagic's own cache maintenance goroutine and is not bounded by it.
const acmeObtainTimeout = 120 * time.Second

// acmeLogOnce wires certmagic's zap logger into the panel log buffer exactly
// once, before the first certmagic.NewDefault() (which copies Default.Logger
// into the shared cache). This surfaces ACME issuance/renewal diagnostics in
// the panel "Logs" view.
var acmeLogOnce sync.Once

// ACMETLSConfig obtains (and thereafter auto-renews) a Let's Encrypt certificate
// for the given domain using the HTTP-01 challenge, and returns a *tls.Config
// ready to wrap a listener with tls.NewListener.
//
// HTTP-01 requires TCP port 80 to be free and reachable from the internet:
// certmagic's HTTP solver transiently binds :80 for the duration of each
// challenge (issuance and renewal), then releases it. TLS-ALPN is disabled
// because the panel/sub servers do not necessarily listen on :443.
//
// Certificates are persisted under storageDir so they survive restarts and we
// avoid re-issuing (which would otherwise hit Let's Encrypt rate limits).
func ACMETLSConfig(ctx context.Context, domain, email, storageDir string) (*tls.Config, error) {
	acmeLogOnce.Do(func() {
		certmagic.Default.Logger = logger.NewZapForwarder("acme")
	})
	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = email
	certmagic.DefaultACME.DisableTLSALPNChallenge = true // use HTTP-01 on port 80
	certmagic.Default.Storage = &certmagic.FileStorage{Path: storageDir}

	cfg := certmagic.NewDefault()

	obtainCtx, cancel := context.WithTimeout(ctx, acmeObtainTimeout)
	defer cancel()
	if err := cfg.ManageSync(obtainCtx, []string{domain}); err != nil {
		return nil, err
	}

	tlsConfig := cfg.TLSConfig()
	tlsConfig.NextProtos = append([]string{"h2", "http/1.1"}, tlsConfig.NextProtos...)
	return tlsConfig, nil
}
