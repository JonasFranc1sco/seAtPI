# Como usar

Para executar o projeto localmente, siga as etapas abaixo.

## 1. Subir os containers

Na raiz do projeto, execute:

```bash
docker compose up -d --build
```

## 2. Entrar no container

Acesse o container da aplicação.

## 3. Executar a aplicação

Já dentro do container, estando na raiz do projeto, execute:

```bash
go run app/cmd/api/main.go
```
