# Environment configuration
SHELL := /bin/bash
export PATH := $(PATH):/usr/local/bin:/opt/homebrew/bin

# Tooling
NPM := npm
NPX := npx
GO := go

.PHONY: test-all test-backend test-frontend test-e2e lint lint-backend lint-frontend help

# Default target
help:
	@echo "Available commands:"
	@echo "  make test-all       - Run all tests (backend, frontend, E2E)"
	@echo "  make test-backend   - Run backend Go tests"
	@echo "  make test-frontend  - Run frontend Jest tests"
	@echo "  make test-e2e       - Run frontend Playwright E2E tests"
	@echo "  make test-report    - Show Playwright E2E test report"
	@echo "  make lint           - Run all linters"
	@echo "  make lint-backend   - Run backend linter"
	@echo "  make lint-frontend  - Run frontend linter"

test-all: test-backend test-frontend test-e2e

test-backend:
	@echo "Running Backend Tests..."
	cd backend && $(GO) test ./...

test-frontend:
	@echo "Running Frontend Unit Tests..."
	cd frontend && $(NPM) test

test-e2e:
	@echo "Running E2E Tests (Production Mode)..."
	cd frontend && $(NPM) run build
	cd frontend && $(NPX) playwright test

test-report:
	@echo "Cleaning up existing report servers..."
	-@lsof -ti :9323 | xargs kill -9 2>/dev/null || true
	@echo "Opening report..."
	cd frontend && $(NPX) playwright show-report $(if $(PORT),--port $(PORT))

lint: lint-backend lint-frontend

lint-backend:
	@echo "Linting Backend..."
	cd backend && golangci-lint run ./...

lint-frontend:
	@echo "Linting Frontend..."
	cd frontend && $(NPM) run lint
