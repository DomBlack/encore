//go:build encore_shell

package auth

import (
	"context"
	"encoding/json"
	goerrs "errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"golang.org/x/oauth2"
)

var ErrInvalidRefreshToken = goerrs.New("invalid refresh token")
var ErrNotLoggedIn = goerrs.New("not logged in: run 'encore auth login' first")

// NewClient returns an HTTP client that is authenticated with the Encore API.
func NewClient() (*http.Client, error) {
	tokenSource, err := newTokenSource()
	if err != nil {
		return nil, err
	}

	return oauth2.NewClient(nil, tokenSource), nil
}

func newTokenSource() (*TokenSource, error) {
	ts := &TokenSource{}

	token, err := ts.readTokenFromConfig()
	if err != nil {
		return nil, err
	}

	ts.token = token
	ts.cfg = &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: APIBaseURL + "/login/oauth:refresh-token",
		}}
	return ts, nil
}

// Config represents the stored Encore configuration.
type Config struct {
	oauth2.Token
	Actor     string `json:"actor,omitempty"`    // The ID of either the user or app authenticated
	Email     string `json:"email,omitempty"`    // non-zero if logged in as a user
	AppSlug   string `json:"app_slug,omitempty"` // non-zero if logged in as an app
	WireGuard struct {
		PublicKey  string `json:"pub,omitempty"`
		PrivateKey string `json:"priv,omitempty"`
	} `json:"wg,omitempty"`
}

// Write persists the configuration for the user.
func writeConfig(cfg *Config) (err error) {
	dir, err := dir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, ".auth_token")
	if data, err := json.Marshal(cfg); err != nil {
		return errors.Wrap(err, "failed to marhsal auth config")
	} else if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return errors.Wrap(err, "failed to create auth config directory")
	} else if err := os.WriteFile(path, data, 0600); err != nil {
		return errors.Wrap(err, "failed to write auth config file")
	}
	return nil
}

func currentUser() (*Config, error) {
	dir, err := dir()
	if err != nil {
		return nil, err
	}
	conf, err := readConf(dir)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func readConf(configDir string) (*Config, error) {
	path := filepath.Join(configDir, ".auth_token")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.WithStack(ErrNotLoggedIn)
		}
		return nil, errors.Wrap(err, "failed to read auth token file")
	}
	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, errors.CombineErrors(err, ErrNotLoggedIn)
	}
	return &conf, nil
}

// TokenSource implements oauth2.TokenSource by looking up the
// current logged in user's API Token.
type TokenSource struct {
	token *oauth2.Token
	cfg   *oauth2.Config
}

// Token implements oauth2.TokenSource.
func (ts *TokenSource) Token() (*oauth2.Token, error) {

	// Use the built-in token source to simplify the logic of
	// refreshing the token as necessary.
	fetch := ts.cfg.TokenSource(context.Background(), ts.token)
	token, err := fetch.Token()
	if err != nil {
		var re *oauth2.RetrieveError
		if errors.As(err, &re) && re.Response.StatusCode == 422 {
			// The refresh token is invalid. Log the user out to reset the token.
			return nil, errors.WithStack(ErrInvalidRefreshToken)
		}
	} else if token.AccessToken != ts.token.AccessToken {
		// The token has changed, so update the config.
		cfg, err := currentUser()
		if err != nil {
			return nil, err
		}
		cfg.Token = *token
		if err := writeConfig(cfg); err != nil {
			return nil, err
		}
	}
	return token, err
}

// readTokenFromConfig reads the oauth token from the config file.
func (ts *TokenSource) readTokenFromConfig() (*oauth2.Token, error) {
	cfg, err := currentUser()
	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.WithStack(ErrNotLoggedIn)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg.Token, nil
}

// Dir reports the directory where Encore's configuration is stored.
func dir() (string, error) {
	dir := os.Getenv("ENCORE_CONFIG_DIR")
	if dir == "" {
		d, err := os.UserConfigDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to get user config dir")
		}
		dir = filepath.Join(d, "encore")
	}
	return dir, nil
}

// APIBaseURL is the base URL for communicating with the Encore Platform.
var APIBaseURL = (func() string {
	if u := os.Getenv("ENCORE_PLATFORM_API_URL"); u != "" {
		return u
	}
	return "https://api.encore.dev"
})()
