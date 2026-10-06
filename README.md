# Fibonizer

A simple Fibonacci sequence calculator, that returns the value for the number that is requested.

This is a project made to learn a bit about the Go language, remember some concepts about cache and code optimization.

The final product should have all the various finished implementations that I came up with (and their performance improvements).

## Project structure

- `backend/` - Go API (Echo) with the Fibonacci implementations
- `frontend/` - React + Vite + TypeScript page to request F(N) and compare implementation times
- `docker-compose.yml` - template to wire the services together (work in progress)

## Commands

All commands are run from the repository root.

### Backend

Run the tests (inside Docker):

```bash
docker build --target run-test-stage -f backend/Dockerfile.multistage ./backend
```

Create the Docker image and run a container:

```bash
docker build --tag fibonizer-backend ./backend
docker run --publish 8080:8080 fibonizer-backend
```

Then request the Fibonizer, where `<method>` is either `loop` or `recursive`:

```bash
curl http://localhost:8080/<method>/<n>
```

The response is JSON. `result` is a string so large numbers keep their precision in JavaScript:

```json
{"method":"loop","n":10,"result":"55","durationNs":270}
```

### Frontend

Create the Docker image and run a container (with the backend running on port 8080):

```bash
docker build --tag fibonizer-frontend ./frontend
docker run --publish 3000:80 fibonizer-frontend
```

Then open <http://localhost:3000>. To point it to another backend, pass `--build-arg VITE_API_URL=<url>` to `docker build`.

For local development with hot reload, run `npm install` and `npm run dev` inside `frontend/`.

## Journey

This chapter documents the dificulties I had throughout this project. 

### Go

The first dificulty I had was Go. I decided I wanted to try this language for this project, so I did. It's not like it was hard to learn, but I was used to Typescript/Javascript and coming to Go was somewhat different.

Overall I think that's the only "problem" I had with Go. Every time I wanted to do something new, I had to go to the docs in order to check how it was done.

Along with the language, I didn't know anything about its libraries, so the process was somewhat the same.

#### Libaries

Right now we are using an outdated version of the package **Echo**, which is used to create a server. We are using the `Echo v4` but the most recent is the `Echo v5`. The upgrade will be done soon.

#### 2026-10-06

I am growing a bit more accustomed to the new syntax of Go. Coming from a language with a lot of freedom like JavaScript and TypeScript, it is a bit "weird" dealing with Go syntax, especially the arrays.

I decided to maintain previous versions of the algorithms to keep a track record of all the thought processes I had throughout this "adventure".

The next step is to explore a little bit about the Fibonacci sequence and how to optimize it.

### Docker

When I started, I did not remember anything about Docker. I did remember that it had images and then we would build images into containers but that was it.

At first I tried to work with Docker Desktop, thinking it was a "softer" approach to Docker, but I decided to stick to the command line after a short while, since it was easier to rebuild the images that way.

Docker was actually the first think I set up, with the help of the [Creating a Go image with Docker](https://docs.docker.com/guides/golang/). Since I was going to have a Frontend, API and a Cache using Redis, I decided this was going to be the first step I was going to take.

### Fibonacci

Fibonacci sequence is rather simple when we look at it. F(0) = 0; F(1) = 1; F(n) = F(n-2) + F(n-1). Its simplicity is actually the problem. We implement the first time, and then it works, great! But then we look at the implementation and think: "This is not right...". We try with a big-ish number, lets say, 1000, and everything comes crumbling down. 

At the time of writing, I still have not implemented the part where we can ask for a specific number N to calculate, but with the first implementations of the algorithm, testing F(8) and it had the following results:

- Recursive: around 400ms
  - The problem? I was calculating F(n-2) and F(n-1), which is an O(xM) operation, unnecessarily.
- Loop: around 200ms
  - Still not great. It also makes sense beacuse it's half of the time that the recursive takes.

After seeing this, I thought that, of course, I had to improve performance, especially if a lot of people wanted to Fibonize a lot of numbers (especially large ones).

My first plan, before researching anything about Fibonacci:

- Fix the recursive implementation;
- Implement a cache for Fibonized numbers in Redis to cut calculation times for future calculations.

This is my first plan to make the Fibonizer a bit more performant. Since we cannot have a cache for every number (otherwise we'll run out of space), the cache will have a caveat. Instead of saving every single number, we'll save a number between checkpoints.
The checkpoint threshold will have to be tested for both performance and memory management, but it should be something like the following example

First time Fibonizing a number, with a checkpoint threshold of 50

- Calculate F(200)
- Check if there's any checkpoint < 200 saved
- There's none
- Calculate everything
- Every 50 "steps", save the number

Second time Fibonizing a number

- Calculate F(332)
- Check if there's any checkpoint
- We know 200, which is  < 332
- Calculate from 332 until 201
- Every 50, we save

For this particular implementation, a cap in memory should be implemented, and, whenever that cap is reached, we clear the lower checkpoints, which are always quicker to calculate.

#### After research

This area is destined to be filled once some research about the Fibonacci sequence is conducted.

### Frontend

For the initial frontend, I decided to ask for some help, in this case, Claude's help, and generated a simple frontend using React (with TypeScript) that is served via nginx.

The frontend shows an input that asks for a number. That number is then sent to the backend to Fibonize. It also shows the different algorithms that were used and how long they took to Fibonize the number that was requested.
