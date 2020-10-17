<p align="center">
	<img src="./assets/kable.png" width="50%" align="center" alt="kable-logo">
</p>

# Kable
![Go](https://github.com/redradrat/kable/workflows/Go/badge.svg?branch=master)
![Release](https://img.shields.io/github/v/release/redradrat/kable)
![License](https://img.shields.io/github/license/redradrat/kable)

A tool to manage kubernetes resources in a GitOps fashion.

It reads so-called "concepts" (see [terminology](#terminology)), and renders them into 
various deployable formats. (YAML Manifests, GitOps controller instructions, helm charts, etc.) 

## Install 

Download the binaries from *Releases*, and install to a directory on your path.

**macOS**

```
brew tap redradrat/kable
brew install kable
```

## Getting Started

You want to manage your kubernetes resources with kable? Great! 

Here is a quick primer:
* Kable essentially is a tool to manage so-called [concepts](#concept) and the
[repos](#repo) they reside in.
* A concept is a directory, containing a `concept.json` and "source", defining Kubernetes resources.
* A repo is a git repository, containing concepts.
* [Rendering](#render) is the process of "instantiating" a concept, 
resulting in various deployable formats. 
 
First we want to add an existing concept [repository](#repository):

```
~/Development/getting-started 
❯ kable repo add demo https://github.com/redradrat/demo-concepts.git
Fetching repository...
? Does this repository require basic authentication? N
✔ Successfully added repository!
```

> There you go! We've added our first repository! What can we do with that?

A repository contains various concepts that are ready to be rendered by us. To list
our available concepts, we can use the `kable list` command.

```
~/Development/getting-started 
❯ kable list
       ID      | REPOSITORY |                MAINTAINER                 
---------------+------------+-------------------------------------------
  apps/grafana | demo       | Name <email>  
  apps/sentry  | demo       | Name <email>  
```
As you can see, there are some concepts already available in the repository. 
Why not head over there and check it out?

Now let's actually [render](#render) one of these into a format that we can actually
apply to our k8s cluster. With `kable render` we can do so! This command has quite a few
tricks up its sleeve. So be sure to check out what you can do with `kable render --help`.

For now let's use `kable render apps/grafana@demo -o out/`: 

```
~/Development/getting-started 
❯ kable render apps/grafana@demo -o out/
Fetching Concept 'apps/grafana@demo'...
Mandatory Values
? instanceName (string) [? for help] 
``` 

As this is our first time rendering this, a dialog will pop up, asking us for values.
So for now we will comply with what this pesky dialog wants.

```
~/Development/getting-started 
❯ kable render apps/grafana@demo -o out/ -t yaml
Fetching Concept 'apps/grafana@demo'...
Mandatory Values
? instanceName (string) test
? nameSelection Option 1
Rendering concept...
✔ Successfully created concept!
```

Alright! Let's see what we got, shall we? 
The `-o` flag we used, defined an output directory `out/` for kable to write to.

```
~/Development/getting-started 
❯ tree out 
out
├── apps-v1_Deployment_test.yaml
├── renderinfo.json
└── v1_Service_test.yaml
```

Apparently our concept consists of multiple k8s resources. A separate manifest has 
been created for each. You can change this behavior by using `-s`, which will render
all resources into a single manifest, `manifest.yaml`.

Notice the `renderinfo.json` file? This file contains the information of how this
rendering has been created. On subsequent render runs, the values we initially provided
will be reused, if this file is detected. You can even pass a specific file with `-r`!

```
~/Development/getting-started 
❯ bat out/renderinfo.json 
───────┬────────────────────────────────────────────────────────────────────
       │ File: out/renderinfo.json
───────┼────────────────────────────────────────────────────────────────────
   1   │ {
   2   │     "version": 1,
   3   │     "meta": {
   4   │         "date": "17 Oct 20 17:42 CEST"
   5   │     },
   6   │     "origin": {
   7   │         "repository": "https://github.com/redradrat/demo-concepts",
   8   │         "ref": "refs/heads/master"
   9   │     },
  10   │     "values": {
  11   │         "instanceName": "test",
  12   │         "nameSelection": "Option 1"
  13   │     }
  14   │ }
───────┴────────────────────────────────────────────────────────────────────
```

When we now execute our render command again, we will see that kable automatically
detects the `renderinfo.json` file in the output path, and reuses it's values.

```
~/Development/getting-started 
❯ kable render apps/grafana@demo -o out/ -t yaml 
Fetching Concept 'apps/grafana@demo'...
Rendering concept...
✔ Successfully created concept!
```

See? No pesky dialog this time!

Whenever kable auto-detects a `renderinfo.json` in the output path, or if we pass an 
existing one via `-r path/to/renderinfo.json`, the dialog will not appear.

That's it! You're now able to render and use concepts!
Make sure to check out the [development](#development) section to take a deeper 
look at how to write concepts.

## Usage

```
~ 
❯ kable    
Usage:
   [command]

Available Commands:
  helm        Tools to interact with helm
  help        Help about any command
  init        Initialize a concept in the current folder
  list        List all available concepts
  render      Render a concept
  repo        Add/List/Remove concept repositories for kable
  serve       Run kable as a server
  version     Show version information

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```

## Terminology

### Concept

A *concept* is a blueprint of an app. It is written in a specific language can be rendered to various outputs.

Supported Types:
* Jsonnet
* JavaScript/Typescript (upcoming)

A concept defines a specific set of [inputs](#inputs), that are passed on to the underlying type. (jsonnet, javascript, etc.)

Each concept needs to build on its own, that's why there is no dependency concept in kable. If let's say a Jsonnet 
concept depends on another Jsonnet concept, this should be realized via the Jsonnet-specific package management.

**Example:**

```
demo-concepts/apps/grafana 
❯ tree
.
├── Makefile
├── concept.json
├── jsonnetfile.json
├── lib
│   ├── k.libsonnet
│   └── main.libsonnet
├── main.jsonnet
└── vendor
```

#### Concept File

What makes a regular directory a concept, is a `concept.json` file at its root. It describes:
* metadata - Name and maintainer of the Concept
* type - The concept type tells kable what the actual content is. Jsonnet? Javascript?
* inputs - See [inputs](#inputs).

#### Inputs

Inputs are a core aspect of any concept. They define a set of required or optional values, that are needed to render
the underlying definiton.

When rendering via `kable render` a dialog will ask the user about the defined inputs.

**Types:**
* string
* int
* bool
* map
* select

**concept.json**

Examples of how to define aforementioned inputs in your concept.json file:

```
{
    "apiVersion": 1,
    "type": "jsonnet",
    "metadata": {...},
    "inputs": {
        "mandatory": {

            // HERE WE CAN DEFINE INPUTS

            "string": {
                "type": "string",
            },

            "int": {
                "type": "int",
            },

            "bool": {
                "type": "bool",
            },

            "map": {
                "type": "map",
            },

            "selection": {
                "type": "select",
                "options": [
                    "Option 1",
                    "Option 2"
                ]
            },

        },
        "optional": {...},
    }
}
```

#### Repo

*Repos* are git repositories that contain multiple concepts. They are used as a 
platform for exchange of concepts, and to render concepts from.

Crucially, they contain a `kable.json` file at their root, listing all the concepts 
contained within.

Example kable.json:
```
{
  "version": 1,
  "concepts": [
    "apps/grafana",
    "apps/sentry"
  ]
}
```

A local kable installation can configure multiple repositories at the same time.

**Demo Repository**

A demo repository can be found at https://github.com/redradrat/demo-concepts

```
kable repo add demo https://github.com/redradrat/demo-concepts.git
```

### Render

*Rendering*, means to instantiate a concept. It's "Application" so to say. Multiple output targets supported.

Supported Targets:
* YAML
* FluxCD Application (upcoming)
* Kable Application (upcoming)

Rendering a concept will give the user a dialog, helping users to define their input values. These values will be
stored in the `renderinfo.json` file. On consecutive render interactions, and pointing kable to this file, those 
values will be reused. 

## Development

*TBD*

## Thanks

* [grafana/tanka](https://github.com/grafana/tanka) - For sticking to jsonnet and maintaining 
an amazing project. Kable relies on their yaml rendering "engine" including the amazing *Helmraiser* 😎🔥.