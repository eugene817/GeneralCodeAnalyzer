# General Code Analyzer 🚀

A high-performance automated source code analysis platform designed for statistical, syntactic, and dynamic analysis of SQL and backend code. Built with **Go (Golang)**, **Echo**, **HTMX**, and **Docker-in-Docker (DinD)** architecture, this platform isolates code execution environments to provide performance benchmarks, static code analysis (linters), and AI-driven optimizations.

---

## 💡 The Idea & Vision

The core objective is to build a unified system capable of executing, profiling, and analyzing user-submitted code in completely isolated sandboxes. By capturing dynamic runtime metrics (CPU, memory footprints, execution time) and static analytics, the platform helps beginner developers understand resource constraints and guides them toward writing optimized, production-grade code. 

Future milestones include training customized LLM/AI models to identify structural patterns and recommend architectural improvements based on collected execution statistics.

---

## 🛠 Tech Stack

- **Backend Core:** Go (Golang) 1.23, Echo Web Framework
- **Dynamic Frontend:** UI rendered serverside with Go Templates, styled via TailwindCSS, and powered by HTMX for smooth, asynchronous SPA-like interactions without heavy JS frameworks.
- **Database:** PostgreSQL (with migration layers and schema models)
- **Infrastructure & Sandboxing:** Docker Engine API, Docker-in-Docker (DinD), Multi-stage Docker optimization.

---

## 🏗 Repository Structure & Architecture

```text
├── api/
│   ├── handlers/   # HTTP request lifecycles, route registration & input validation
│   ├── static/     # Asset pipeline (compiled TailwindCSS styles)
│   └── templates/  # Go html/template files driven by HTMX endpoints
├── config/         # Application bootstrap configurations & environment management
├── database/       # DB initialization, ORM models, and automatic migration layers
├── services/       # Core business logic: container orchestration, LLM interfacing, analytics
├── Dockerfile      # Highly optimized, multi-stage production container definition
└── main.go         # Application bootstrap & high-level infrastructure orchestration
