# gitlabctl
[![pipeline status](https://gitlab.com/isca/gitlabctl/badges/master/pipeline.svg)](https://gitlab.com/isca/gitlabctl/commits/master)
[![coverage report](https://gitlab.com/isca/gitlabctl/badges/master/coverage.svg)](https://gitlab.com/isca/gitlabctl/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/isca/gitlabctl)](https://goreportcard.com/report/gitlab.com/isca/gitlabctl)

<center>
<img src="img/gitlabctl.svg"/>
</center>

This is a CLI to manipulate the gitlab API trough the command line. You can list,copy,delete or create projects and groups with a single command.

### Prerequisites

All you'll need is the `go.1.12` or superior and API token to talk with your gitlab sessions. Moreover on that in the [Getting Started](https://gitlab.com/isca/gitlabctl#getting-started).
This package also use cobra and viper packages. You can just run go get to acquire than.

### Installing

What things you need to install the software and how to install them from the source:

Get the package:

```
go get gitlab.com/isca/gitlabctl
```

Install it:

```
go install
```

You can also build and install using docker, all you'll need is run the `Makefile`:

```
make install
```

## Getting Started

This cli was inspired by the unixes command line, so you'll find almost the same syntax when using it.

#### Listing

_listing all the projects and groups on your gitlab session:_
```
export SESSIONA="myApiToken"
gitlabctl ls 
```

_listing all the content of an specific group:_
```
export MYOTHERSESSION="myApiToken"
gitlabctl ls mygroup/mysubgroup -t myothersession
```

#### Copying 
If you have two gitlab sessions for example an on-premisse omnibus session and a gitlab.com session, you'll
will find this command useful to migrate groups from one session to another.

_copying an entire group with all subgroups from one session to another:_
```
export FROMUSER="myuserFromTheSourceCopy"
export TOUSER="myUserFromTheDestinationCopy"
export SESSIONA="apiTokenOfSourceCopy"
export SESSIONB="apiTokenOfDestinationCopy"
gitlabctl cp group --to=sessionA:mygroup/mysubgroup --to=sessionB:someGroup
```
_If the `someGroup` doesn't exist in `sessionB`, it will be automatically created._

_copying a project from one session to another:_
```
export FROMUSER="myuserFromTheSourceCopy"
export TOUSER="myUserFromTheDestinationCopy"
export SESSIONA="apiTokenOfSourceCopy"
export SESSIONB="apiTokenOfDestinationCopy"
gitlabctl cp proj --to=sessionA:mygroup/mysubgroup/myproject --to=sessionB:someGroup
```
_this will copy `myproject` from `sessionA` into `someGroup` on sessionB._


#### Deleting

You can delete a group like this:
```
export SESSIONA="apiTokenOfSourceCopy"
gitlabctl rm group --from=sessionA:mygroup/mysubgroup
```
_This command will delete the group `mysubgroup` and all the inside content._

## Environments

This project is cloud-native by design so you can run this code with environments instead of the config.toml file.
Here is the list of support environments:

  *	SESSION(\*N) `(string)` specify the token for a gitlab session. _You can use any variable name for the API token, by default it will look to a variable named `SESSIONA`._
  * FROMUSER `(string)` set the login of the user of sessionA (used during the copy command)
  * TOSUER `(string)` set the login of the user of sessionB (used during the copy command)

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

## Built With

* [go](http://golang.org/) - The Go Programming Language

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://git.pgd.to/tools/email/mailhook_consumer/tags).

## Authors

* **Igor Brandao** - *Initial work* - [isca](https://gitlab.com/isca)

See also the list of [contributors](https://gitlab.com/isca/gitlabctl/project_members) who can be participated in this project.
_Isca disclaims all copyright interest in the program “gitlabctl” (which is a gitlab CLI tool) is written by Igor Brandao_  

