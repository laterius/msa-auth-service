package domain

type Db struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	DbName   string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Extras   string `env:"DB_EXTRAS"`
}

type Config struct {
	Db   Db
	Http struct {
		Port string `env:"HTTP_PORT"`
	}
	App struct {
		Env string `env:"APP_ENV"`
	}
}
