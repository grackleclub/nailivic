# nailivic
Go + htmx rewrite of [nailivic](https://github.com/Turkosaurus/nailivic)

## frontend
`htmx` drives the frontend, utilizing HTML as the engine of application state.

### Air
Live reload for Go apps, to see your changes in realtime without having to restart your server

Via `go install`
```bash
go install github.com/air-verse/air@latest
```
Then add an alias to your .bashrc(or whatever you have) like `alias air='~/usr/go/bin/air'` so you can run the air server with fewer keystrokes.

Github for [air](https://github.com/air-verse/air).

## templates
Templates are handled by go's standard `template/html` package. 

### ordering templates
Templates may be comprised of other templates, but template import order matters.
> [!IMPORTANT]
> Always import **parent** elements first, then any **children** parsed after.

This will work, because index defines subtemplates first, then the subtemplates are loaded.
```go
tmpl, err := template.ParseFS(content,
    "static/html/index.html",
    "static/html/head.html",
    "static/html/footer.html",
)
```

This *will not work*. 😿
```go
tmpl, err := template.ParseFS(content,
    "static/html/head.html",
    "static/html/footer.html",
    "static/html/index.html",
)
```

This will only work if the files are alphabetized. 🙃
```go
tmpl, err := template.ParseFS(content,
    "static/html/*.html",
)
```

## backend
Backend is a bog standard go server, utilizing a filesystem embedded in the binary for ease and speed.
