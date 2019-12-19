# gitlabctl
[![pipeline status](https://gitlab.com/isca/gitlabctl/badges/master/pipeline.svg)](https://gitlab.com/isca/gitlabctl/commits/master)
[![coverage report](https://gitlab.com/isca/gitlabctl/badges/master/coverage.svg)](https://gitlab.com/isca/gitlabctl/commits/master)

This is a CLI to manipulate the gitlab API trough command line. You can list,copy,delete or create projects and groups in a single line.

### Prerequisites

This package is built with go1.12 and use only the packages from stdlibrary so there are no external dependencies  
to compile this code.
All you'll need is the `go.1.12` or superior and a kubernetes resource template. Moreover on that in the template section of this manual.
_If will use installation with Makefile using the `make install` command you'll also need some version of docker installed._

### Installing

What things you need to install the software and how to install them from the source:

```
go install
```

You can also build and install using docker, all you'll need is run the `Makefile`:

```
make install
```
_This process creates a docker container to build and install the go binary, so it needs a docker installed._

To uninstall run:

```
make uninstall
```

## Running the tests

Until I finish this README there are not so many Unit tests written.  
But I will try to cover at least 80% of unit tests for this code as soon as possible.  

You can run tests like this:

```
go test ./...
```

To run this code locally for test purposes use:

```
go run main.go
```

## Getting Started

You'll need to create a valid k8s template and use `kploy` to render it. For e.g:

```
---
# Deployment resource template for project {{ .ProjectName }} on k8s 1.14.6
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .ProjectName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .ProjectName }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: {{ .ProjectName }}
  template:
    metadata:
      labels:
        app: {{ .ProjectName }}
    spec:
      containers:
        - name: {{ .ProjectName }}
          image: {{ .ContainerName }}
          ports:
          - containerPort: {{ .ContainerPort }}
          resources:
            requests:
              memory: "{{ .Memory }}"
              cpu: "{{ .Cpu }}"
            limits:
              memory: "{{ .LimitMem }}"
              cpu: "{{ .LimitCpu }}"
```

_Assuming this template as deploy-template.yml, to render this file you'll need to export the template variables,
and after that, you'll be ready to use kploy._

```
export ProjectName=MyProject;\
export NameSpace=MyNameSpace;\
export Replicas=1;\
export ContainerName=MyContainer:Version;\
export ContainerPort=8000;\
export Memory=10Mi;\
export Cpu=100m;\
export LimitMem=12Mi;\
export LimitCpu=120m;\
kploy -c deploy-template.yml -o > deploy.yml
```
_If `-c` was not specified this script will search by default a template in `/etc/kploy/deploy.yml`, which is provided with the command `make install`._
_By default with `-o` flag will only output the rendered template._

_You can render and apply the rendered template direct to the default cluster with:_

```
export ProjectName=MyProject;\
export NameSpace=MyNameSpace;\
export Replicas=1;\
export ContainerName=MyContainer:Version;\
export ContainerPort=8000;\
export Memory=10Mi;\
export Cpu=100m;\
export LimitMem=12Mi;\
export LimitCpu=120m;\
kploy -c deploy-template.yml -apply
```

_You can specify a different cluster using the flag -target yourKubeConfig._


That's the part where you automate everything on your CI and forget.

There are a few more flags on this command you can list them with `kploy -h`.

## Custom Variables

You can inject your own variables into containers with kploy. All you'll need is  
declare your variables starting with `KP_`. e.g:
```
...
export KP_MyVar="myValue";\
...
kploy -c deploy-template -apply
```

The command above will create a template like:
```
...
spec:
   containers:
   - env:
     - name: "MyVar"
     - value: "myValue"
...
```
The variable `MyVar` can be accessed from inside the containers. 

## Host Aliases

You can add custom host aliases with kploy by exporting a variable like this:
```
...
export KHOST="hostA.com,hostB.com>1.1.1.1:hostC.io>2.2.2.2"
...
kploy -c deploy-template -apply
```

You can also append variables into `KHOST` with this syntax:
```
export HOSTA="hostA.com>1.1.1.1"
export HOSTBC="hostB.com,hostC.com>2.2.2.2"
export KHOST=$HOSTA:$HOSTBC
...
kploy -c template -o
```
Both syntaxes are supported and the result will be the same.

As a result of the above commands, you should have a deployment with these entries:
```
...
spec:
   hostAliases:
   - hostnames:
    - "hostA.com"
    - "hostB.com"
    ip: "1.1.1.1"
   - hostnames:
    - "hostC.io"
    ip: "2.2.2.2"
   containers:
   ...
...
```

## Environments

This project is cloud-native by design so you can run this code with environments instead of the config.toml file.
Here is the list of support environments:

  *	SESSION(\*N) `(string)` specify the token for a gitlab session.
  * FROMUSER `(string)` set the login of the user of sessionA (used during the copy command)
  * TOSUER `(string)` set the login of the user of sessionB (used during the copy command)

## Built With

* [go](http://golang.org/) - The Go Programming Language

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://git.pgd.to/tools/email/mailhook_consumer/tags).

## Authors

* **Igor Brandao** - *Initial work* - [isca](https://gitlab.com/isca)

See also the list of [contributors](https://gitlab.com/isca/gitlabctl/project_members) who can be participated in this project.
_Isca disclaims all copyright interest in the program “gitlabctl” (which render kubernetes templates from environment variables) is written by Igor Brandao_  

