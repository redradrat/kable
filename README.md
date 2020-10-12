<p align="center">
	<img src="./assets/kable.png" width="50%" align="center" alt="kable-logo">
</p>

# Kable
![Go](https://github.com/redradrat/kable/workflows/Go/badge.svg?branch=master)
![Release](https://img.shields.io/github/v/release/redradrat/kable)
![License](https://img.shields.io/github/license/redradrat/kable)

A tool to manage kubernetes specs in a GitOps fashion.

It reads so-called "concepts" (see [terminology](#terminology)), and renders them into various deployable formats. (YAML Manifests, GitOps controller instructions, helm charts, etc.) 

## Usage

```
~ 
❯ kable    
Usage:
   [command]

Available Commands:
  help        Help about any command
  init        Initialize a concept in the current folder
  list        List all available concepts
  render      Render a concept
  repo        Add/List/Remove concept repositories for kable
  serve       Run kable as a server
  version     Show version information

Flags:
  -c, --config string   config file (default is $HOME/.config/kable/settings.json)
  -h, --help            help for this command
  -t, --toggle          Help message for toggle

Use " [command] --help" for more information about a command.
```

## Install 

**macOS**

```
brew tap redradrat/kable
brew install kable
```

**Others**

Download the binaries from *Releases*, and install to a directory on your path.

### Terminology

**Concept**

A *concept* is a blueprint of an app. It is written in a specific language can be rendered to various outputs.

Supported Types:
* Jsonnet
* JavaScript/Typescript (upcoming)

Each concept needs to build on its own, that's why there is no dependency concept in kable. If let's say a Jsonnet concept depends on another Jsonnet concept, this should be realized via the Jsonnet-specific package management.

**Repo** 

*Repos* are git repositories that contain multiple concepts. They are used as a platform for exchange of concepts, and to render concepts from.

**Rendering**

*Rendering*, means to instantiate a concept. It's "Application" so to say. Multiple output targets supported.

Supported Targets:
* YAML
* FluxCD Application (upcoming)
* Kable Application (upcoming)

## Concepts

A *concept* is a blueprint of an app. It is written in a specific language can be rendered to various outputs.

Supported Types:
* Jsonnet
* JavaScript/Typescript (upcoming)

Each concept needs to build on its own, that's why there is no dependency concept in kable. If let's say a Jsonnet concept depends on another Jsonnet concept, this should be realized via the Jsonnet-specific package management.

### Inputs 
