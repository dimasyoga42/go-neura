package supabase

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Config struct {
	Url string `env:"SUPABASE_URL,required"`
	Key string `env:"SUPABASE_KEY,required"`
}

func Supabase() *supabase.Client {
	// Load .env file kalau ada, kalau tidak ada tetap lanjut pakai env sistem
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from system environment")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse env: %+v", err)
	}

	supa, err := supabase.NewClient(cfg.Url, cfg.Key, nil)
	if err != nil {
		log.Fatalf("failed to create supabase client: %+v", err)
	}

	return supa
}
