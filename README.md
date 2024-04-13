# bl-forms

## install

    go install github.com/cosmtrek/air@latest
    go install github.com/swaggo/swag/cmd/swag@latest
    go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest


    git config user.name "gitUsername"
    git config user.email "git@email"

    ~/go/bin/swag init -t public,auth && \
    ~/go/bin/swag init -t intern --instanceName intern 

    ~/go/bin/swag init -t public && ~/go/bin/swag init -t intern --instanceName intern 

## DEV

     ~/go/bin/gin  --laddr 0.0.0.0 --all --appPort 4090 --excludeDir webapp run main.go     
     ~/go/bin/air  --laddr 0.0.0.0 --all --appPort 4090 --excludeDir webapp run main.go



     https://github.com/SushritPasupuleti/Go-Chi-Boilerplate/tree/main


goi18n extract --outdir i18n -format yaml 
goi18n merge  --outdir i18n  -format yaml i18n/active.*.yaml 
goi18n merge  --outdir i18n  -format yaml i18n/active.*.yaml i18n/translate.*.yaml 
