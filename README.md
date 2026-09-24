# Fibonizer

A simple Fibonacci sequence calculator, that returns the value for the number that is requested.

This is a project made to learn a bit about the Go language, remember some concepts about cache and code optimization.

The final product should have all the various finished implementations that I came up with (and their performance improvements).

## Journey

This chapter documents the dificulties I had throughout this project. 

### Go

The first dificulty I had was Go. I decided I wanted to try this language for this project, so I did. It's not like it was hard to learn, but I was used to Typescript/Javascript and coming to Go was somewhat different.

Overall I think that's the only "problem" I had with Go. Every time I wanted to do something new, I had to go to the docs in order to check how it was done.

Along with the language, I didn't know anything about its libraries, so the process was somewhat the same.

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
