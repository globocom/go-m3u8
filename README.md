# go-m3u8

# Troubleshooting

Erros ao executar o comando `go get`:

```
        remote: 
        remote: ========================================================================
        remote: 
        remote: The project you were looking for could not be found or you don't have permission to view it.
        remote: 
        remote: ========================================================================
        remote: 
        fatal: Could not read from remote repository.

        Please make sure you have the correct access rights
        and the repository exists.
```

Execute os comandos a seguir:

1. Crie o token de acesso no perfil do gitlab. Siga
   a [documentação oficial](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

2. Adiciona token de acesso no arquivo `.netrc`

```
echo "machine gitlab.globoi.com login <USUARIO> password <ACCESS_TOKEN>" > ~/.netrc

```

3. Agora será possível fazer download do repositório

```
GOPRIVATE="*.globoi.com" go get "gitlab.globoi.com/webmedia/media-delivery-advertising/go-m3u8"
```
