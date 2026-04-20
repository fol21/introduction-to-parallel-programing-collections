# Introduction to Parallel Programming: Example Collections

A collection of source code discussed in **"An Introduction to Parallel Programming"** by Peter Pacheco.

---

## About

This repository gathers the example programs from the book, which covers the three most widely-used parallel programming models in C/C++:

| Technology | Description |
|---|---|
| **MPI** (Message Passing Interface) | Distributed-memory parallelism via message passing between processes |
| **Pthreads** (POSIX Threads) | Shared-memory parallelism using low-level POSIX thread primitives |
| **OpenMP** | Shared-memory parallelism using compiler directives and a runtime library |

### C++ Programs Overview

The examples cover a progressive set of parallel-programming concepts, including:

- **Hello World** — basic process/thread spawning and rank/ID identification
- **Trapezoidal Rule (numerical integration)** — data decomposition and collective communication (`MPI_Reduce`)
- **Odd-Even Transposition Sort** — point-to-point communication patterns
- **Matrix–Vector Multiplication** — partitioned data distribution, scatter/gather
- **Parallel Prefix Sum** — tree-structured communication and reduction trees
- **Producer–Consumer Queue** — mutex locks, condition variables (Pthreads)
- **Read-Write Locks** — concurrent reads with exclusive writes (Pthreads)
- **Barriers and Semaphores** — synchronisation primitives
- **Cache-Coherence & False Sharing** — performance pitfalls in shared-memory programs
- **OpenMP Loop Parallelism** — `#pragma omp parallel for`, schedules, and reductions
- **Task Parallelism** — `#pragma omp task` and the OpenMP tasking model

---

## Goals

The long-term goal of this project is to **translate every C++ example into idiomatic Go and Rust**, making it easy to compare how each language expresses the same parallel algorithms.

### Go (Golang)

Go's concurrency model — goroutines and channels — maps naturally to many of the message-passing and shared-memory patterns in the book.  The translations will use:

- `goroutines` + `channels` as equivalents to MPI processes and message passing
- `sync.Mutex`, `sync.RWMutex`, and `sync.WaitGroup` as equivalents to Pthreads primitives
- The standard `sync/atomic` package for lock-free operations

### Rust

Rust's ownership model provides compile-time safety guarantees for concurrent code.  The translations will use:

- `std::thread` with message-passing via `std::sync::mpsc` channels
- `Arc<Mutex<T>>` and `Arc<RwLock<T>>` for shared-state concurrency
- The [Rayon](https://github.com/rayon-rs/rayon) data-parallelism library as a high-level OpenMP analogue
- The [MPI crate](https://github.com/rsmpi/rsmpi) for distributed-memory examples

---

## Repository Layout (planned)

```
.
├── cpp/          # Original C/C++ examples (MPI, Pthreads, OpenMP)
│   ├── mpi/
│   ├── pthreads/
│   └── openmp/
├── go/           # Go translations
│   ├── mpi/
│   ├── goroutines/
│   └── sync/
└── rust/         # Rust translations
    ├── mpi/
    ├── threads/
    └── rayon/
```

---

## Prerequisites

| Language | Tools needed |
|---|---|
| C++ | GCC/Clang, OpenMPI, libpthread, libomp |
| Go | Go ≥ 1.21 |
| Rust | Rust (via `rustup`) ≥ 1.75, Cargo |

---

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file included in this repository.
