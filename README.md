# Unconditional API

Unconditional API is a Go project designed to provide a web server, a feed repository, and a source client for [Unconditional](https://unconditional.day/about).

![Project Status](https://img.shields.io/badge/status-active-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)

## :file_folder: Project Structure

The project is structured into several key components:

- **Web Server**: Located in the `webserver` package, the server uses the Echo framework to define and handle HTTP routes. It's configured with parameters such as the server address, port, and allowed origins. The server is defined in the 'cmd/serve.go' file.

- **Feed Repository**: Defined in the `bleveRepo` package, the repository uses the Bleve library to create and manage a full-text index for feeds.

- **Source Client**: Located in the `github` package, the client uses the Github library to interact with the Github API.

## :wrench: Installation

To install the project, you'll need to have Go installed on your machine. You'll also need to install the Echo, Bleve, and Github libraries. Once you have these prerequisites, you can clone the repository and build the project using the `make build` command in the 'Makefile'. You can also use the 'Dockerfile' to build a Docker image of the project.

## :computer: Running the Project

To run the project, use the `go run cmd/serve.go` command. If the project is running successfully, you should see output indicating that the server is running and listening on the configured port.

## :handshake: Contributing

We welcome contributions to the project. To contribute, you can submit an issue describing the bug you found or the feature you want to add. You can also create a pull request with your proposed changes. Please follow our code of conduct when contributing. For more information on how to contribute, refer to the 'Makefile', 'Dockerfile', and 'cmd/serve.go'.

## License and Acknowledgments

This project is licensed under the MIT License. We'd like to acknowledge the creators and contributors of the Echo, Bleve, and Github libraries, which this project uses extensively.
