# bot-translator

Bot provides an interface to the [libretranslate](https://github.com/LibreTranslate/LibreTranslate) (self-hosted)

### Usage
- /start - create new user in db
- /supported - list of supported languages
- /primary \<lang\> - set primary language
- /secondary \<lang\> - set secondary language
- /t \<lang1> \<lang2\> - transalte replied message from lang1 to lang2
- \<text message\> - translate the text into the primary language, from the primary language translates into the secondary language

## Instalation (docker-compose)
1) create docker-compose.yml from examples/docker-compose-example.yml
2) set env args in docker-compose.yml
3) exec "docker-compose up -d"
4) wait until bot start (libretranslate, may be launched within 30 minutes)
### Env args
- BOT_TOKEN - Telegram bot token
- LIBRETRANSLATE_URL - url to libretranslate
- DB_CONNSTR - db connetion string (postgresql://<db_username>:<db_password>@<db_ip>:<db_port>/<db_name>?sslmode=disable)
- DEFAULT_SECONDARY_LANG - default secondary language
- DEFAULT_PRIMARY_LANG - default primary language
