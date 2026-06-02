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

// acmeMu serialises all writes to certmagic's package-level globals
// (certmagic.Default, certmagic.DefaultACME). Without this, concurrent calls
// from the web server and the sub server goroutines would race on Email, Storage
// and Logger, with the last writer winning unpredictably.
var acmeMu sync.Mutex

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
func ACMETLSConfig(domain, email, storageDir string) (*tls.Config, error) {
	acmeMu.Lock()
	defer acmeMu.Unlock()

	certmagic.Default.Logger = logger.NewZapForwarder("acme")
	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = email
	certmagic.DefaultACME.DisableTLSALPNChallenge = true // use HTTP-01 on port 80
	certmagic.Default.Storage = &certmagic.FileStorage{Path: storageDir}

	cfg := certmagic.NewDefault()

	// Derive from context.Background(), NOT the server's context: a panel
	// restart (SIGHUP) cancels the server context before Start() runs again,
	// which would abort issuance immediately with "context canceled". The
	// obtain is still bounded by acmeObtainTimeout.
	obtainCtx, cancel := context.WithTimeout(context.Background(), acmeObtainTimeout)
	defer cancel()
	if err := cfg.ManageSync(obtainCtx, []string{domain}); err != nil {
		return nil, err
	}

	tlsConfig := cfg.TLSConfig()
	tlsConfig.NextProtos = append([]string{"h2", "http/1.1"}, tlsConfig.NextProtos...)
	return tlsConfig, nil
}
