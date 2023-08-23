# fleet.api-gateway
## Repo contains everything related to the Tyk.io api gateway/mgmt platform including:
1. IaC code for api definitions management and deployment
2. Custom Go plug-in for custom rate limiting for phase 1 legacy implementation ([see README for details](src/custom-go-plugin/README.md))
   - everything that is necessary to develop and test a custom Go plugin that can then be used by Tyk.io api gateway 
   - see [Makefile](https://github.com/Omnitracs/fleet.api-gateway/blob/develop/src/custom-go-plugin/Makefile) and above README for details regarding excuting an actual build and bringing of the docker containers for tesing  of plugin code
   is in this src directory
   - this dev environment is provided by Tyk for just this purpose -- the environment itself is Linux basedf and there should be used within a Linux environment (it has been tested and confirmed using a Windows Hyper-V Ubuntu VM)
   - the basic idea is that a number of containers are provided that provide an out-of-the-box runtime environment which can be used for development
   - therefore, a container runtime environment such as docker and docker-compose are required along with Git, Go, Make and VSCode with plugins:
     - Go language support
     - Gitlens
     - Go Test Explorer (not evaluated yet)
## Go Dev resources:
- [Effective Go](https://go.dev/doc/effective_go)
- [Go for beginers video](https://youtu.be/yyUHQIec83I) -- 3:24 basic Go tutorial

## Go Custom Plugin Dev Env Setup
1. configure Ubuntu VM using Windows Hyper-V
- create and configure: [sample instructions](https://cloudbytes.dev/snippets/install-ubuntu-in-a-vm-on-windows-using-hyper-v) (quick create option is typically sufficient)
- you will likely need to expand the VM disk size: [sample instructions](http://linguist.is/2020/08/12/expand-ubuntu-disk-after-hyper-v-quick-create/)
- can also change ubuntu video resolution: [sample instructions](https://cloudbytes.dev/snippets/make-ubuntu-fullscreen-on-windows-hyper-v)

2. install Git
```shell
$ sudo apt update
$ sudo apt install git-all

# check install
$ git --version

# configure git
$ git config --global user.name "replace with your github user name here"
$ git config --global user.email youremail@solera.com
```
3. install Go
```shell
# note: the make script is setup to use the snap install location
$ sudo snap install go --classic
```
4. install Make 
```shell
$ sudo apt-get install -y build-essential
```
5. install Docker
- follow instructions [HERE](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04)

6. install Docker Compose
- folow instrctions [HERE](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-22-04)

7. run docker commands without sudo
- https://docs.docker.com/engine/install/linux-postinstall/ 
- NOTE: after adding your user to the new "docker" group you will MOST LIKELY need to shutdown and then restart your VM

8. install VS Code
```shell
$ sudo snap install --classic code
```
- add VS Code plugins for Go and GitLens
- [configure git access](https://code.visualstudio.com/docs/sourcecontrol/github) in VS Code
- since Go was installed with snap there is a little Code configuration that likely needs be updated in order for Code debugging to work as expected for Go code with F5
  - Go to File | Preferences | Settings
  - Search for "go.alternateTools"
  - Select the "User" tab towards the top left
  - Select "Edit in settings.json"
  - Add "go": "/snap/go/current/bin/go" as a property to the empty value object of "go.alternateTools" as follows:
    ```json
    "go.alternateTools": {
        "go": "/snap/go/current/bin/go"
    }
    ```
  - Save the settings file
  - Restart Code
  - Install any Go related tool that a notification is received for
  - F5 should now work for debugging Go code
  - see this [thread](https://github.com/golang/vscode-go/issues/1411#issuecomment-816972221) for some additional info

## Example Rate Limiting Custom Plugin Configs
Request Headers strategy
  ```json
  {
  "rateLimiting": {
    "LogLevel": 0,
    "active": true,
    "overrides": [
      {
        "method": "GET",
        "requests": -1,
        "resource": "/hello/ability/",
        "seconds": -1
      },
      {
        "method": "GET",
        "requests": 5,
        "resource": "/hello/location/",
        "seconds": 60
      }
    ],
    "requests": 2,
    "seconds": 10,
    "sessionTtlMin": 120,
    "strategy": {
      "config": {
        "headerNames": [
          "x-tenant-id",
          "Host",
          "Content-Type",
          "Authorization"
        ],
        "separator": "::"
      },
      "name": "requestHeaders"
    }
  }
}
```
Session Guid strategy
```json
{
  "rateLimiting": {
    "active": true,
    "LogLevel": 0,
    "overrides": [
      {
        "method": "GET",
        "requests": -1,
        "resource": "/rna-soap/ability/",
        "seconds": -1
      },
      {
        "method": "GET",
        "requests": 5,
        "resource": "/rna-soap/location/",
        "seconds": 60
      }
    ],
    "requests": 2,
    "seconds": 10,
    "sessionTtlMin": 120,
    "strategy": {
      "config": {
        "headerNames": [
          "Authorization"
        ],
        "separator": "::"
      },
      "name": "sessionGuid"
    }
  }
}
```
Request headers XRS strategy
```json
{
  "rateLimiting": {
    "active": true,
    "LogLevel": 0,
    "overrides": [
      {
        "method": "GET",
        "requests": -1,
        "resource": "/xrs/ability/",
        "seconds": -1
      },
      {
        "method": "GET",
        "requests": 5,
        "resource": "/xrs/location/",
        "seconds": 60
      }
    ],
    "requests": 2,
    "seconds": 10,
    "sessionTtlMin": 120,
    "strategy": {
      "config": {
        "headerNames": [
          "Authorization"
        ],
        "separator": "::"
      },
      "name": "requestHeadersXRS"
    }
  }
}
```