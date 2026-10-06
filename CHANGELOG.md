# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

Since the project is just a personal project, with no necessary releases and versions, dates will be used instead of versions.

## 2026-10-5

### Changed
- Added `FibonizeRecursiveV1` and `FibonizeRecursiveV2` to keep a track of what was done.
- Updated `README.md` with new information.
- Added changes to `CHANGELOG.md`

## 2026-09-28

## 2026-10-5

### Added
- `frontend/` application served through nginx, with its own Dockerfile.
- `docker-compose.yml` to run the services together.
- `CHANGELOG.md` (added with Claude's help, to track the previous commits and generate the CHANGELOG).

### Changed
- Moved the Go service (`main.go`, `main_test.go`, `go.mod`, `go.sum`, Dockerfiles) into `backend/`.
- Split `.gitignore` so the backend has its own.
- Improved `FibonizeRecursive` function performance in `main.go`.

## 2026-09-28

### Added
- Dynamic parameter to request a specific number: `/recursive/:num` and `/loop/:num` (5de5a88).
- README documentation for the new endpoints.

### Changed
- Fibonacci functions now return `int64` instead of `int`, so larger results fit.

## 2026-09-24

### Added
- README documentation expanded (2be4ec1).

## 2026-09-23

### Added
- `FibonizeRecursive` and `FibonizeLoop` functions, exposed on `/recursive` and `/loop` with a fixed input of 8 (0f66eb1).
- Logging of the time each implementation takes.
- Welcome page listing the available endpoints.

### Removed
- The template's `IntMin` helper.

## 2026-09-22

### Added
- Initial project setup from Docker's Go template ([docker-gs-ping](https://github.com/docker/docker-gs-ping)): Echo server, `Dockerfile`, `Dockerfile.multistage`, Go module and tests (d5cedcf).

## 2026-09-21

### Added
- Initial commit with `.gitignore` and `README.md` (b5b219c).
