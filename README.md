# Dev

Dev is a tool to enable local development of modern applications.

* Upstream dependencies can be run from standard images using Docker.
* Every process gets a hostname entry that other processes can access.
* Processes can be knitted together, downstream process only starting when upstream processes are ready,
* Logs are captured and saved so they can be searched later on.
* Do this all rootless (hard to impossible?).

## References

- [Containers from scratch](https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909)