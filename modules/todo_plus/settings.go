package todo_plus

import (
	"os"

	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const defaultTitle = "Todo"

type Settings struct {
	common *cfg.Common

	backendType     string
	backendSettings *config.Config
	projects        []interface{}
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {

	backend, _ := ymlConfig.Get("backendSettings")

	settings := Settings{
		common: cfg.NewCommonSettingsFromModule(name, defaultTitle, ymlConfig, globalConfig),

		backendType:     ymlConfig.UString("backendType"),
		backendSettings: backend,
	}

	return &settings
}

func FromTodoist(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	apiKey := ymlConfig.UString("apiKey", ymlConfig.UString("apikey", os.Getenv("WTF_TODOIST_TOKEN")))
	projects := ymlConfig.UList("projects")
	backend, _ := config.ParseYaml("apiKey: " + apiKey)
	backend.Set(".projects", projects)

	settings := Settings{
		common: cfg.NewCommonSettingsFromModule(name, defaultTitle, ymlConfig, globalConfig),

		backendType:     "todoist",
		backendSettings: backend,
	}

	return &settings
}

func FromTrello(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {

	accessToken := ymlConfig.UString("accessToken", ymlConfig.UString("apikey", os.Getenv("WTF_TRELLO_ACCESS_TOKEN")))
	apiKey := ymlConfig.UString("apiKey", os.Getenv("WTF_TRELLO_APP_KEY"))
	board := ymlConfig.UString("board")
	username := ymlConfig.UString("username")
	var lists []interface{}
	list, err := ymlConfig.String("list")
	if err == nil {
		lists = append(lists, list)
	} else {
		lists = ymlConfig.UList("list")
	}
	backend, _ := config.ParseYaml("apiKey: " + apiKey)
	backend.Set(".accessToken", accessToken)
	backend.Set(".board", board)
	backend.Set(".username", username)
	backend.Set(".lists", lists)

	settings := Settings{
		common: cfg.NewCommonSettingsFromModule(name, defaultTitle, ymlConfig, globalConfig),

		backendType:     "trello",
		backendSettings: backend,
	}

	return &settings
}
