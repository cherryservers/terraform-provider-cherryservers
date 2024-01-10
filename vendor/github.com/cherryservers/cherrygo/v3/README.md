# cherrygo

Cherry Servers golang API client library for Cherry Servers RESTful API.

You can view the client API docs here: [https://pkg.go.dev/github.com/cherryservers/cherrygo](https://pkg.go.dev/github.com/cherryservers/cherrygo)

You can view Cherry Servers API docs here: [https://api.cherryservers.com/doc](https://api.cherryservers.com/doc)

## Table of Contents

- [Installation](#installation)
- [Authentication](#authentication)
- [Examples](#examples)
  - [Get teams](#get-teams)
  - [Get projects](#get-projects)
  - [Get plans](#get-plans)
  - [Get images](#get-images)
  - [Order new server](#order-new-server)

## Installation

Download the library to you GOPATH:
```
go get github.com/cherryservers/cherrygo
```

Then import the library in your Go code:
```
import "github.com/cherryservers/cherrygo"
```

### Authentication

To authenticate to the Cherry Servers API, you must have API token, you can create authentication tokens in the [Cherry Server client portal](https://portal.cherryservers.com). Token must be exported in env var `CHERRY_AUTH_TOKEN` or passed to client directly.

```
export CHERRY_AUTH_TOKEN="4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
```
Use exported CHERRY_AUTH_TOKEN env variable:
```go
func main() {
    c, err := cherrygo.NewClient()
}
```
Pass token directly to client:
```go
func main() {
    c, err := cherrygo.NewClient(cherrygo.WithAuthToken("your-api-token"))
}
```

### Examples ###

#### Get teams
You will need team ID for later calls, for example to get projects for specified team, you will need to provide team ID.
```go
teams, _, err := c.Teams.List(nil)
if err != nil {
    log.Fatal("Error", err)
}

for _, t := range teams {
    log.Println(t.ID, t.Name, t.Credit.Promo.Remaining, t.Credit.Promo.Usage, t.Credit.Resources.Pricing.Price)
}
```

#### Get projects
After you have your team ID, you can list your projects. You will need your project ID to list your servers or order new ones.
```go
projects, _, err := c.Projects.List(teamID, nil)
if err != nil {
    log.Fatal("Error", err)
}

for _, p := range projects {
    log.Println(p.ID, p.Name, p.Href)
}
```

#### Get plans
Next thing in order to get new server is to choose one, we call it plans

```go
plans, _, err := c.Plans.List(teamID, nil)
if err != nil {
    log.Fatalf("Plans error: %v", err)
}

for _, p := range plans {
    log.Println(p.Name, p.Slug)
}
```

#### Get images
After you manage to know desired plan, you need to get available images for that plan
```go
images, _, err := c.Images.List(planSlug, nil)
if err != nil {
    log.Fatal("Error", err)
}

for _, i := range images {
    log.Println(i.Name, i.Slug)
}
```

#### Order new server
Now you are ready to order new server
```go
addServerRequest := cherrygo.CreateServer{
    ProjectID:   projectID,
    Hostname:    hostname,
    Image:       imageSlug,
    Region:      regionSlug,
    Plan:        planSlug,
}

server, _, err := c.Server.Create(&addServerRequest)
if err != nil {
    log.Fatal("Error while creating new server: ", err)
}

log.Println(server.ID, server.Name, server.Hostname)
```

## Debug

In case you want to debug this library and get requests and responses from API you need to export CHERRY_DEBUG variable
```
export CHERRY_DEBUG="true"
```

When you done, just unset that variable:
```
unset CHERRY_DEBUG
```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations.
