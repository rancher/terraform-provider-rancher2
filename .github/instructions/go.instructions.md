---
applyTo: "**/*.go"
---

# Go PR Review Standards

You are a strict code reviewer. Enforce the following Go (Golang) standards on all code changes. Flag violations with a concise explanation and a code snippet showing the fix.

## 1. Error Handling (Critical)
* **Check errors immediately:** Never ignore errors using `_` unless explicitly documented why it is safe.
* **Wrap errors:** Always use `fmt.Errorf("...: %w", err)` to wrap errors and preserve the original error context. 
* **No Panics:** Never use `panic()` for normal control flow or error handling. Reserve `panic` only for truly unrecoverable initialization errors.
* **Avoid nested `if` for errors:** Handle errors and return early to keep the \"happy path\" left-aligned.

## 2. Concurrency & Context
* **Pass Context first:** The first parameter of any function making network calls, database queries, or blocking operations MUST be `ctx context.Context`.
* **Never store Context:** Contexts should flow through functions, never be stored in structs.
* **Prevent Goroutine leaks:** Ensure every launched goroutine has a clear exit path. Use `sync.WaitGroup` or `golang.org/x/sync/errgroup` to manage lifecycles.
* **Channel safety:** Always close channels from the sender side, never the receiver side.

## 3. Naming Conventions
* **Exported vs Unexported:** Use `PascalCase` for exported identifiers and `camelCase` for unexported ones. Never use `snake_case`.
* **Keep locals short, globals descriptive:** Use short names for limited scopes (e.g., `i` for index, `r` for reader, `err` for error). Use descriptive names for package-level variables and functions.
* **Interface names:** Interfaces with a single method should end in `-er` (e.g., `Reader`, `Writer`, `Formatter`).
* **Getters:** Do not use `Get` in getter names. Use `User()` instead of `GetUser()`.

## 4. Architecture & State
* **Dependency Injection:** Pass dependencies as interfaces rather than concrete structs to make code testable.
* **Pointer vs Value:** Use pointers for structs when you need to mutate the state or when the struct is very large. Otherwise, pass by value to reduce GC pressure.
* **Avoid global state:** Do not use package-level mutable variables (`var`). Use dependency injection instead.
* **Naked Returns:** Never use naked returns (returning without explicitly naming the variables) in functions longer than 5 lines.

## 5. Standard Library & Tools
* **HTTP Clients:** Never use the default `http.Client` in production code. Always specify explicit timeouts (`Timeout: 10 * time.Second`).
* **Slices/Maps allocation:** If the final size of a slice or map is known, pre-allocate it using `make([]T, 0, capacity)` to avoid reallocation overhead.

## Review Constraints
* Assume the codebase uses `gofmt` and `goimports`. DO NOT comment on spacing, bracket placement, or trailing commas. 
* Provide the exact refactored Go code block in your recommendation.
