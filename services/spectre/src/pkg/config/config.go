package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_PATH_TO_ENV = "../../.env"
)

type Config struct {
	LogLevel string
	JWTexp   int // in hours

	Server ServerCfg
	Routes RoutesCfg
}

type ServerCfg struct {
	SpectreHost string
	SpectrePort string

	ProxyHost string
	ProxyPort string
}

type RoutesCfg struct {
	APILettersPoint string
	APIUsersPoint   string

	ECDHPoint      string
	AuthLoginPoint string

	ProxyEncryptPoing string
	ProxyDecryptPoing string
	ProxyECDHPoing    string
}

func MustLoad() *Config {
	if os.Getenv("START") == "manually" {
		if err := godotenv.Load(DEFAULT_PATH_TO_ENV); err != nil {
			log.Fatalf("No .env file found (and ENV not set): %v", err)
		}
	}

	llevel := getEnv("SPECTRE_LOG_LEVEL", "debug")
	sjwtexp := getEnv("SPECTRE_JWT_EXP", "24")
	jwtexp, err := strconv.Atoi(sjwtexp)
	if err != nil {
		log.Fatalf("Cannot parse jwt exp as int: %v", err)
	}
	shost := getEnv("SPECTRE_HOST", "0.0.0.0")
	sport := getEnv("SPECTRE_PORT", "5000")
	chost := getEnv("CRYPTO_HOST", "0.0.0.0")
	cport := getEnv("CRYPTO_PORT", "7654")

	apiLP := getEnv("SPECTRE_API_LTS_POINT", "/api/letters")
	apiUP := getEnv("SPECTRE_API_USRS_POINT", "/api/users")
	ecdhP := getEnv("SPECTRE_ECDH_POINT", "/ecdh")
	loginP := getEnv("SPECTRE_LOGIN_POINT", "/login")

	pep := getEnv("CRYPTO_ENCR_POINT", "/encrypt")
	pdp := getEnv("CRYPTO_DECR_POINT", "/decrypt")
	pecdh := getEnv("CRYPTO_ECDH_POINT", "/ecdh")

	sc := ServerCfg{
		SpectreHost: shost,
		SpectrePort: sport,

		ProxyHost: chost,
		ProxyPort: cport,
	}
	rc := RoutesCfg{
		APILettersPoint: apiLP,
		APIUsersPoint:   apiUP,

		ECDHPoint:      ecdhP,
		AuthLoginPoint: loginP,

		ProxyEncryptPoing: pep,
		ProxyDecryptPoing: pdp,
		ProxyECDHPoing:    pecdh,
	}

	cfg := Config{
		LogLevel: llevel,
		JWTexp:   jwtexp,

		Server: sc,
		Routes: rc,
	}

	return &cfg
}

func getEnv(key string, defVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defVal
}
