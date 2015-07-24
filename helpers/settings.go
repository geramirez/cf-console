package helpers

import (
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"encoding/gob"
	"errors"
)

// Settings is the object to hold global values and objects for the service.
type Settings struct {
	// OAuthConfig is the OAuth client with all the paramters to talk with CF's UAA OAuth Provider.
	OAuthConfig *oauth2.Config
	// Console API
	ConsoleAPI string
	// Sessions is the session store for all connected users.
	Sessions sessions.Store
}

// InitSettings attempts to populate all the fields of the Settings struct. It will return an error if it fails,
// otherwise it returns nil for success.
func (s *Settings) InitSettings(envVars EnvVars) error {
	if len(envVars.ClientID) == 0 {
		return errors.New("Unable to find '" + ClientIDEnvVar + "' in environment. Exiting.\n")
	}
	if len(envVars.ClientSecret) == 0 {
		return errors.New("Unable to find '" + ClientSecretEnvVar + "' in environment. Exiting.\n")
	}
	if len(envVars.Hostname) == 0 {
		return errors.New("Unable to find '" + HostnameEnvVar + "' in environment. Exiting.\n")
	}
	if len(envVars.AuthURL) == 0 {
		return errors.New("Unable to find '" + AuthURLEnvVar + "' in environment. Exiting.\n")
	}
	if len(envVars.TokenURL) == 0 {
		return errors.New("Unable to find '" + TokenURLEnvVar + "' in environment. Exiting.\n")
	}
	if len(envVars.APIURL) == 0 {
		return errors.New("Unable to find '" + APIEnvVar + "' in environment. Exiting.\n")
	}
	s.ConsoleAPI = envVars.APIURL

	// Setup OAuth2 Client Service.
	s.OAuthConfig = &oauth2.Config{
		ClientID:     envVars.ClientID,
		ClientSecret: envVars.ClientSecret,
		RedirectURL:  envVars.Hostname + "/oauth2callback",
		Scopes:       []string{"cloud_controller.read", "cloud_controller.write", "cloud_controller.admin", "scim.read", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  envVars.AuthURL,
			TokenURL: envVars.TokenURL,
		},
	}

	// Initialize Sessions.
	// Temp FIXME that fixes the problem of using a cookie store which would cause the secure encoding
	// of the oauth 2.0 token struct in production to exceed the max size of 4096 bytes.
	filesystemStore := sessions.NewFilesystemStore("", []byte("some key"))
	filesystemStore.MaxLength(4096 * 4)
	s.Sessions = filesystemStore
	// Want to save a struct into the session. Have to register it.
	gob.Register(oauth2.Token{})

	return nil
}