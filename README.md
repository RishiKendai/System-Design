# system-design

This repository is a workspace for learning **system design** from two angles that usually show up together in practice.

**High-level design (HLD)** is about the system as a whole: major components, how they communicate, and how the design behaves under load, failure, and change. That material is still being captured here.

**Low-level design (LLD)** is about turning requirements into a concrete program shape: entities, responsibilities, and the interactions between them. Those exercises live under [`lld/`](lld/) and are implemented in **Go** as small, runnable modules you can read and extend.

## LLD — topics covered

- [Car parking](lld/car-parking/)
- [URL shortener](lld/url-shortner/)
- [Chess CLI](lld/chess/)

    You can download it, try playing chess in the terminal, and have fun. For download and install instructions, see the [Chess CLI instructions](lld/chess/README.md).

## HLD — topics covered

- [Distributed rate limiter](hld/distributed-rate-limiter/)
