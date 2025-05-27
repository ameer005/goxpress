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

## 📦 Installation

```bash
go get github.com/ameer005/goxpress

```
