# bl-forms

## install

    go install github.com/wandercn/hotbuild@latest
    go install github.com/cosmtrek/air@latest
    go install github.com/swaggo/swag/cmd/swag@latest

    git config user.name "gitUsername"
    git config user.email "git@email"

    ~/go/bin/swag init -t public && \
    ~/go/bin/swag init -t intern --instanceName intern 

    ~/go/bin/swag init -t public && ~/go/bin/swag init -t intern --instanceName intern 

## DEV

     ~/go/bin/gin  --laddr 0.0.0.0 --all --appPort 4090 --excludeDir webapp run main.go     
     ~/go/bin/air  --laddr 0.0.0.0 --all --appPort 4090 --excludeDir webapp run main.go