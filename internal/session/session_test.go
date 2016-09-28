package session

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/golib/assert"
	"github.com/golib/aws/internal"
	"github.com/golib/aws/internal/credentials"
	"github.com/golib/aws/internal/defaults"
)

func TestNewDefaultSession(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	s := New(&internal.Config{Region: internal.String("region")})

	assert.Equal(t, "region", *s.Config.Region)
	assert.Equal(t, http.DefaultClient, s.Config.HTTPClient)
	assert.NotNil(t, s.Config.Logger)
	assert.Equal(t, internal.LogOff, *s.Config.LogLevel)
}

func TestNew_WithCustomCreds(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	customCreds := credentials.NewStaticCredentials("AKID", "SECRET", "TOKEN")
	s := New(&internal.Config{Credentials: customCreds})

	assert.Equal(t, customCreds, s.Config.Credentials)
}

type mockLogger struct {
	*bytes.Buffer
}

func (w mockLogger) Log(args ...interface{}) {
	fmt.Fprintln(w, args...)
}

func TestSessionCopy(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	os.Setenv("AWS_REGION", "orig_region")

	s := Session{
		Config:   defaults.Config(),
		Handlers: defaults.Handlers(),
	}

	newSess := s.Copy(&internal.Config{Region: internal.String("new_region")})

	assert.Equal(t, "orig_region", *s.Config.Region)
	assert.Equal(t, "new_region", *newSess.Config.Region)
}

func TestSessionClientConfig(t *testing.T) {
	s, err := NewSession(&internal.Config{Region: internal.String("orig_region")})
	assert.NoError(t, err)

	cfg := s.ClientConfig("s3", &internal.Config{Region: internal.String("us-west-2")})

	assert.Equal(t, "https://s3-us-west-2.amazonaws.com", cfg.Endpoint)
	assert.Empty(t, cfg.SigningRegion)
	assert.Equal(t, "us-west-2", *cfg.Config.Region)
}

func TestNewSession_NoCredentials(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	s, err := NewSession()
	assert.NoError(t, err)

	assert.NotNil(t, s.Config.Credentials)
	assert.NotEqual(t, credentials.AnonymousCredentials, s.Config.Credentials)
}

func TestNewSessionWithOptions_OverrideProfile(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "other_profile")

	s, err := NewSessionWithOptions(Options{
		Profile: "full_profile",
	})
	assert.NoError(t, err)

	assert.Equal(t, "full_profile_region", *s.Config.Region)

	creds, err := s.Config.Credentials.Get()
	assert.NoError(t, err)
	assert.Equal(t, "full_profile_akid", creds.AccessKeyID)
	assert.Equal(t, "full_profile_secret", creds.SecretAccessKey)
	assert.Empty(t, creds.SessionToken)
	assert.Contains(t, creds.ProviderName, "SharedConfigCredentials")
}

func TestNewSessionWithOptions_OverrideSharedConfigEnable(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "0")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "full_profile")

	s, err := NewSessionWithOptions(Options{
		SharedConfigState: SharedConfigEnable,
	})
	assert.NoError(t, err)

	assert.Equal(t, "full_profile_region", *s.Config.Region)

	creds, err := s.Config.Credentials.Get()
	assert.NoError(t, err)
	assert.Equal(t, "full_profile_akid", creds.AccessKeyID)
	assert.Equal(t, "full_profile_secret", creds.SecretAccessKey)
	assert.Empty(t, creds.SessionToken)
	assert.Contains(t, creds.ProviderName, "SharedConfigCredentials")
}

func TestNewSessionWithOptions_OverrideSharedConfigDisable(t *testing.T) {
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "full_profile")

	s, err := NewSessionWithOptions(Options{
		SharedConfigState: SharedConfigDisable,
	})
	assert.NoError(t, err)

	assert.Empty(t, *s.Config.Region)

	creds, err := s.Config.Credentials.Get()
	assert.NoError(t, err)
	assert.Equal(t, "full_profile_akid", creds.AccessKeyID)
	assert.Equal(t, "full_profile_secret", creds.SecretAccessKey)
	assert.Empty(t, creds.SessionToken)
	assert.Contains(t, creds.ProviderName, "SharedConfigCredentials")
}

func TestNewSessionWithOptions_Overrides(t *testing.T) {
	cases := []struct {
		InEnvs    map[string]string
		InProfile string
		OutRegion string
		OutCreds  credentials.Value
	}{
		{
			InEnvs: map[string]string{
				"AWS_SDK_LOAD_CONFIG":         "0",
				"AWS_SHARED_CREDENTIALS_FILE": testConfigFilename,
				"AWS_PROFILE":                 "other_profile",
			},
			InProfile: "full_profile",
			OutRegion: "full_profile_region",
			OutCreds: credentials.Value{
				AccessKeyID:     "full_profile_akid",
				SecretAccessKey: "full_profile_secret",
				ProviderName:    "SharedConfigCredentials",
			},
		},
		{
			InEnvs: map[string]string{
				"AWS_SDK_LOAD_CONFIG":         "0",
				"AWS_SHARED_CREDENTIALS_FILE": testConfigFilename,
				"AWS_REGION":                  "env_region",
				"AWS_ACCESS_KEY":              "env_akid",
				"AWS_SECRET_ACCESS_KEY":       "env_secret",
				"AWS_PROFILE":                 "other_profile",
			},
			InProfile: "full_profile",
			OutRegion: "env_region",
			OutCreds: credentials.Value{
				AccessKeyID:     "env_akid",
				SecretAccessKey: "env_secret",
				ProviderName:    "EnvConfigCredentials",
			},
		},
		{
			InEnvs: map[string]string{
				"AWS_SDK_LOAD_CONFIG":         "0",
				"AWS_SHARED_CREDENTIALS_FILE": testConfigFilename,
				"AWS_CONFIG_FILE":             testConfigOtherFilename,
				"AWS_PROFILE":                 "shared_profile",
			},
			InProfile: "config_file_load_order",
			OutRegion: "shared_config_region",
			OutCreds: credentials.Value{
				AccessKeyID:     "shared_config_akid",
				SecretAccessKey: "shared_config_secret",
				ProviderName:    "SharedConfigCredentials",
			},
		},
	}

	for _, c := range cases {
		oldEnv := initSessionTestEnv()
		defer popEnv(oldEnv)

		for k, v := range c.InEnvs {
			os.Setenv(k, v)
		}

		s, err := NewSessionWithOptions(Options{
			Profile:           c.InProfile,
			SharedConfigState: SharedConfigEnable,
		})
		assert.NoError(t, err)

		creds, err := s.Config.Credentials.Get()
		assert.NoError(t, err)
		assert.Equal(t, c.OutRegion, *s.Config.Region)
		assert.Equal(t, c.OutCreds.AccessKeyID, creds.AccessKeyID)
		assert.Equal(t, c.OutCreds.SecretAccessKey, creds.SecretAccessKey)
		assert.Equal(t, c.OutCreds.SessionToken, creds.SessionToken)
		assert.Contains(t, creds.ProviderName, c.OutCreds.ProviderName)
	}
}

func TestSessionAssumeRole_DisableSharedConfig(t *testing.T) {
	// Backwards compatibility with Shared config disabled
	// assume role should not be built into the config.
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "0")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	s, err := NewSession()
	assert.NoError(t, err)

	creds, err := s.Config.Credentials.Get()
	assert.NoError(t, err)
	assert.Equal(t, "assume_role_w_creds_akid", creds.AccessKeyID)
	assert.Equal(t, "assume_role_w_creds_secret", creds.SecretAccessKey)
	assert.Contains(t, creds.ProviderName, "SharedConfigCredentials")
}

func TestSessionAssumeRole_InvalidSourceProfile(t *testing.T) {
	// Backwards compatibility with Shared config disabled
	// assume role should not be built into the config.
	oldEnv := initSessionTestEnv()
	defer popEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_invalid_source_profile")

	s, err := NewSession()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SharedConfigAssumeRoleError: failed to load assume role")
	assert.Nil(t, s)
}

func initSessionTestEnv() (oldEnv []string) {
	oldEnv = stashEnv()
	os.Setenv("AWS_CONFIG_FILE", "file_not_exists")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "file_not_exists")

	return oldEnv
}
