# Goxpress

**Goxpress** is an Express-like web framework built entirely in **Go** — including its own custom **HTTP server implementation**.  
It does **not** rely on Go’s `net/http` package. Instead, it parses raw TCP connections and builds HTTP handling, routing, middleware, and more from scratch.

---

## 🎯 Purpose

> Goxpress is built **for learning and experimentation**.  
> Use it in personal or hobby projects — but it's **not production-ready** and not recommended for professional deployments.

This project is ideal if you want to:

- Understand how web frameworks like Express.js work under the hood
- Learn low-level HTTP, TCP, and middleware patterns in Go
- Experiment with building your own server infrastructure

---

## ✨ Features

- 🔧 Custom HTTP parser — handles raw TCP connections
- 🛣️ Routing with support for:
  - Static routes: `/about`
  - Dynamic routes: `/user/:id`
- 📦 Body parsers:
  - `application/json`
  - `application/x-www-form-urlencoded`
  - `text/plain`
- 🔐 Middleware support (like Express)
  - Global middleware
  - Early return (e.g., auth middleware)
- ⚙️ Typed query and body parsing using Go generics
- 🧠 Custom context: `ctx.Req`, `ctx.Res`, `ctx.Data` — shared across middleware and routes
- 📤 Easy JSON/string responses

---

## 📝 TODO

Here’s what’s coming next (or what you can contribute):

- 🗂️ **Multipart form-data parsing** — support file uploads via `multipart/form-data`
- 🍪 **Cookie handling** — `ctx.Cookies.Get()`, `ctx.Cookies.Set()`, etc.
- 🌐 **CORS middleware** — handle cross-origin resource sharing
- 🛡️ **Rate limiting middleware** — basic protection against abuse
- 📄 **Static file serving** — serve HTML, CSS, JS, etc.
- 🔄 **Custom error handler middleware** — centralized error handling
- 🧰 **Request validators** — validate input data using tags or schemas

> 💡 **Suggestions welcome!** Open an issue or PR if you have an idea you'd like to contribute or discuss.

---

## 📦 Installation

```bash
go get github.com/ameer005/goxpress

```
