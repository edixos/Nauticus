## Getting Started

You must have a working [Go environment](https://golang.org/doc/install) and
then fork the repo and clone it:

```shell  title="Clone Nauticus Locally"
git clone https://github.com/<your-username>/nauticus.git
cd nauticus
```


Do not forget to add the upstream repository to rebase when needed.

```shell  title="Clone Nauticus Locally"
git remote add upstream https://github.com/edixos/nauticus.git
```



## Building & Testing

First you need a local kubernetes cluster if you want to run the controller againt a real API Server, you can create a kind 
cluster for this purpose.


 ```bash title="Create a kind cluster"
 kind create cluster --image kindest/node:v1.24.1 --wait 5m --name nauticus
 ```

The project uses the `make` build system. It'll run code generators, tests and
static code analysis.

Start by displaying the help target of make to discover what we can do with **make**.

```shell  title="Display Make Help"
make help
```

### Initialize your dev environment

1. First Install the `golangci-lint` tool

    ```shell  title="Install golangci-lint tool"
    make golangci-lint
    ```

2. Build the controller locally

    ```shell  title="Build the controller locally"
    make build
    ```

3. Run tests

    ```shell  title="Run tests locally"
    make tests
    ```
4. Run the controller 
    ```shell  title="Run the controller locally"
    make manifests
    kubectl apply -f config/crd/bases/nauticus.io_spaces.yaml
    make run
    ```
5. Lint the code

    ```shell  title="Run the controller locally"
    make golint
    ```
   


## Documentation

We use [mkdocs material](https://squidfunk.github.io/mkdocs-material/). See `/docs` for the source code.

Install mkdocs using pip:

```shell
pip install mkdocs-material
```

When writing documentation it is advised to run the mkdocs server with livereload:

```shell
mkdocs serve
```

Open `http://localhost:8000` in your browser.

We generate CRD documentation using crddoc, Build the api documentation:
```shell  title="Generate Api doc"
make apidoc
```