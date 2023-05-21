\# Installation Guide for development

This guide provides instructions for setting up and running the Xavier repository.

## Prerequisites

Before proceeding with the installation, ensure that the following prerequisites are met:

- Go programming language (version 1.19 or later) is installed. Refer to the \[official Go website\]\(https://golang.org/doc/install\) for installation instructions.
- Buf cli tool must be installed. Refer to the \[official Buf website\]\(https://docs.buf.build/installation\) for installation instructions.
  - The Buf cli tool is used to generate the gRPC code from the proto files.
  - As the project matures, the generated code will be checked in to the repository and this dependency will be removed unless changes are made to the proto files.
  - On a Mac, you can install the Buf cli tool using Homebrew: `brew install buf`

## Clone the Repository

First, clone the Xavier repository to your local machine:

```shell
git clone https://github.com/zencodinglab/xavier.git
```

## Install Dependencies and generate gRPC code
```shell
make buf-generate
```

## Build the Binaries

To build the client and server binaries, navigate to the repository root directory and run the following command:

```shell
make build
```

This command will compile the client and server binaries and place them in the `bin` directory.

## Run the Server

To run the server, execute the following command from the repository root directory:

```shell
make run-server
```

This will start the server application.
## Run the Client

To run the sample client, execute the following command from the repository root directory:

```shell
make run-client
```

This will start the client application.
  - Note: The sample client changes frequently as services are added. 

## Clean Up

To clean up the built binaries, run the following command:

```shell
make clean
```

This command will remove the compiled binaries from the `bin` directory.

## Additional Information

- The `Makefile` provides various build and execution commands. Refer to the file for more details on available commands.
- The repository contains a `docs` directory where you can find additional documentation, including the API reference and usage examples.
