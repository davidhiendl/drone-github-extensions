package shared

type AppConfig struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET" required:"true"`

	DroneConfigIncludeMax      int  `envconfig:"DRONE_CONFIG_INCLUDE_MAX" default:"20"`
	EmulateCIPrefixedVariables bool `envconfig:"EMULATE_CI_PREFIXED_ENV_VARS" default:"true"`
	EnvAddTagSemver            bool `envconfig:"ENV_ADD_TAG_SEMVER" default:"true"`
}
