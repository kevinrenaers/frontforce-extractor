package configuration

type Configuration struct {
	FrontforceURL                 string `mapstructure:"frontforce_url"`
	FrontforceInitialRefreshToken string `mapstructure:"frontforce_initial_refresh_token"`
	FrontforceDashboardID         string `mapstructure:"frontforce_dashboard_id"`
	FrontforceUserID              int    `mapstructure:"frontforce_user_id"`

	HaUrl                          string       `mapstructure:"ha_url"`
	HaToken                        string       `mapstructure:"ha_token"`
	HaFrontforceStatusEntity       entityConfig `mapstructure:"ha_frontforce_status_entity"`
	HaFrontforceInterventionEntity entityConfig `mapstructure:"ha_frontforce_intervention_entity"`
	HaFrontforceVehicleEntity      entityConfig `mapstructure:"ha_frontforce_vehicle_entity"`
}

type entityConfig struct {
	EntityID     string `mapstructure:"id"`
	FriendlyName string `mapstructure:"friendly_name"`
}
