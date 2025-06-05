# Goxpress

**Goxpress** is an Express-like web framework built entirely in Go — including its own custom HTTP server implementation.  
It does not rely on Go's `net/http` package. Instead, it parses raw TCP connections and builds HTTP handling, routing, middleware, and more from scratch.

---

## 🎯 Purpose

Goxpress is built for learning and experimentation.  
Use it in personal or hobby projects — but it's **not production-ready** and not recommended for professional deployments.

This project is ideal if you want to:

- Understand how web frameworks like Express.js work under the hood
- Learn low-level HTTP, TCP, and middleware patterns in Go
- Experiment with building your own server infrastructure

---

## ✨ Features

- 🔧 **Custom HTTP parser** — handles raw TCP connections
- 🛣️ **Routing** with support for:
  - Static routes: `/about`
  - Dynamic routes: `/user/:id`
- 📦 **Body parsers**:
  - `application/json`
  - `application/x-www-form-urlencoded`
  - `text/plain`
  - `multipart/form-data`
- 🔐 **Middleware support** (like Express)
  - Global middleware
  - Early return (e.g., auth middleware)
- 🔐 **Built-in middlewares**
  - Cors
- ⚙️ **Typed query and body parsing** using Go generics
- 🧠 **Custom context**: `ctx.Req`, `ctx.Res`, `ctx.Data` — shared across middleware and routes
- 📤 Easy JSON/string responses

---

## 📦 Installation

```bash
go get github.com/ameer005/goxpress
```

---

## 🚀 Quick Start

### Basic Server

```go
package main

import "github.com/ameer005/goxpress"

func main() {
    var app *goxpress.Server = goxpress.NewServer(":8080")
    router := app.Router

    router.Route(httpmethod.GET, "/", func(ctx *goxpress.Context) {
        ctx.Res.Status(200).JSON(map[string]any{"status": "success"})
    })

    app.Listen()
}
```

---

### Complete Example with Middleware & Routes

```go
package main

import (
    "fmt"
    "log"
    "github.com/ameer005/goxpress"
)

func main() {
    app := goxpress.NewServer(":8080")
    router := app.Router

    // Global middleware - runs for all requests
    router.Use(func(ctx *goxpress.Context) {
        fmt.Printf("Request: %s %s\n", ctx.Req.RequestMethod())
    })

    // Auth middleware
    authMiddleware := func(ctx *goxpress.Context) {
        token := ctx.Req.Headers("Authorization")
        if token == "" {
            ctx.Res.Status(401).JSON(map[string]string{
                "error": "Authorization required",
            })
            return
        }
        ctx.Data["user"] = "authenticated_user"
    }

    // Routes
    router.Route(httpmethod.GET, "/", func(ctx *goxpress.Context) {
        ctx.Res.Status(200).JSON(map[string]string{
            "message": "Welcome to Goxpress!",
            "version": "1.0.0",
        })
    })

    router.Route(httpmethod.GET, "/user/:id", func(ctx *goxpress.Context) {
        userID := ctx.Req.Params["id"]
        ctx.Res.Status(200).JSON(map[string]any{
            "user_id": userID,
            "message": fmt.Sprintf("Hello user %s!", userID),
        })
    })

    router.Route(httpmethod.GET, "/profile", authMiddleware, func(ctx *goxpress.Context) {
        user := ctx.Data["user"]
        ctx.Res.Status(200).JSON(map[string]any{
            "profile": user,
            "message": "This is a protected route",
        })
    })

    router.Route(httpmethod.POST, "/users", func(ctx *goxpress.Context) {
        type CreateUser struct {
            Name  string `json:"name"`
            Email string `json:"email"`
        }

        user, err := JSONBody[CreateUser](ctx.Req)


        ctx.Res.Status(201).JSON(map[string]any{
            "message": "User created successfully",
            "user":    user,
        })
    })

    router.Route(httpmethod.GET, "/search", func(ctx *goxpress.Context) {

		q = ctx.Req.UntypedQuery()

        query := q["q"]
        page := q["page"]

        if query == "" {
            ctx.Res.Status(400).JSON(map[string]string{
                "error": "Query parameter 'q' is required",
            })
            return
        }

        ctx.Res.Status(200).JSON(map[string]any{
            "query":   query,
            "page":    page,
            "results": []string{"result1", "result2", "result3"},
        })
    })

        type UploadForm struct {
        Title       string `form:"title"`
        Description string `form:"description"`
    }

        func uploadHandler(ctx *goxpress.Context) {
            // Parse form fields
            form, err := goxpress.FormBody[UploadForm](ctx.Req)
            if err != nil {
                ctx.Res.Status(400).JSON(map[string]string{
                    "error": "Invalid form fields",
                })
                return
            }

            // Get uploaded file metadata
            file, err := ctx.Req.FormFile("upload")
            if err != nil {
                ctx.Res.Status(400).JSON(map[string]string{
                    "error": "File not found",
                })
                return
            }

            ctx.Res.Status(200).JSON(map[string]any{
                "title":       form.Title,
                "description": form.Description,
                "filename":    file.Filename,
                "size":        file.Size,
                "mime_type":   file.Mime,
            })
        }

        app.Router.Route(httpmethod.POST, "/upload", uploadHandler)


        log.Printf("🚀 Goxpress server starting on http://localhost:8080")
        app.Listen()
    }
```

---

## 🧪 Testing Your Server

```bash
# Basic GET request
curl http://localhost:8080/

# Dynamic route
curl http://localhost:8080/user/123

# Protected route (without auth)
curl http://localhost:8080/profile

# Protected route (with auth)
curl -H "Authorization: Bearer token123" http://localhost:8080/profile

# POST with JSON
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Query parameters
curl "http://localhost:8080/search?q=golang&page=1"

# Form data
curl -X POST http://localhost:8080/contact \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=Jane&email=jane@example.com&message=Hello"

curl -X POST http://localhost:8080/upload \
  -F "title=Test File" \
  -F "description=This is a test upload" \
  -F "upload=@/path/to/file.jpg"

```

---

## 📝 TODO

Here’s what’s coming next (or what you can contribute):

- 🍪 Cookie handling — `ctx.Cookies.Get()`, `ctx.Cookies.Set()`, etc.
- 🛡️ Rate limiting middleware — basic protection against abuse
- 📄 Static file serving — serve HTML, CSS, JS, etc.
- 🔄 Custom error handler middleware — centralized error handling
- 🧰 Request validators — validate input data using tags or schemas

---

💡 Suggestions welcome!  
Open an issue or PR if you have an idea you'd like to contribute or discuss.
