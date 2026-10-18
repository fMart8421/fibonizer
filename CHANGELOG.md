# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

Since the project is just a personal project, with no necessary releases and versions, dates will be used instead of versions.

## 2026-10-18

### Added
- Binet's formula implementation (`FibonizeBinet`, using `math/big`) exposed on `/binet/:num`, so results beyond `int64` fit (dd894d1).
- `getCumulativeSum` helper (not used yet) (dd894d1).

### Changed
- Frontend now offers the `recursive/v1`, `recursive/v2` and `binet` methods (5c86ee3).
- Frontend requests now time out after 5 seconds, to avoid waiting indefinitely (21c6f7b).
- Reformatted `frontend/src/api.ts` (double quotes and semicolons).

### Fixed
- `FibonizeBinet` returned wrong results: the loop skipped the last halving step, the sort compared indexes instead of values, the even/odd branch used the loop index instead of the halved value, and the result pair was returned the wrong way round (54a49e4).

## 2026-10-06

### Added
- `/recursive/v1/:num` and `/recursive/v2/:num` endpoints, each reporting its own method name (e1287f8).
- `FiboError` response: the endpoints now return `400 Bad Request` for negative numbers (e1287f8).
- `GetParsedNumber` helper to share the `:num` parsing between handlers (e1287f8).
- `README.md` notes about keeping previous algorithm versions (e6d4465).

### Changed
- Optimized the recursive Fibonizer: `FibonizeRecursive` now returns the pair `[F(n-1), F(n)]` instead of recomputing both branches (7296af7).
- Kept the previous implementation as `FibonizeRecursiveV1` and the optimized one as `FibonizeRecursiveV2`, to keep a track of what was done (e1287f8).
- Removed the `/recursive/:num` endpoint in favour of the versioned ones above (e1287f8).
- Handlers now log the result along with the time taken, in a more visible format (e1287f8).
- Updated the tests to cover `FibonizeRecursiveV2` (e1287f8).

## 2026-10-05

### Added
- `frontend/` application served through nginx, with its own Dockerfile.
- `docker-compose.yml` to run the services together.
- `CHANGELOG.md` (added with Claude's help, to track the previous commits and generate the CHANGELOG).

### Changed
- Moved the Go service (`main.go`, `main_test.go`, `go.mod`, `go.sum`, Dockerfiles) into `backend/`.
- Split `.gitignore` so the backend has its own.

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
