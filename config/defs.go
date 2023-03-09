package config

type ConfigPostamark struct {
	Token        string `mapstructure:"token"`
	SlackWebhook string `mapstructure:"slack_webhook"`
	ApiURL       string `mapstructure:"api_url"`
}

type ConfigGitlab struct {
	Token         string `mapstructure:"token"`
	GroupID       int    `mapstructure:"group_id"`
	ApiURL        string `mapstructure:"api_url"`
	DefaultBranch string `mapstructure:"default_branch"`
}

type ConfigOnePassword struct {
	Address string `mapstructure:"address"`
}

type Config struct {
	Postmark    ConfigPostamark   `mapstructure:"postmark"`
	Gitlab      ConfigGitlab      `mapstructure:"gitlab"`
	OnePassword ConfigOnePassword `mapstructure:"onepassword"`
}
