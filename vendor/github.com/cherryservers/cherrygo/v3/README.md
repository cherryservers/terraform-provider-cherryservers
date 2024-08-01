# cherrygo

Cherry Servers golang API client library for Cherry Servers RESTful API.

You can view the client API docs here: [https://pkg.go.dev/github.com/cherryservers/cherrygo/v3](https://pkg.go.dev/github.com/cherryservers/cherrygo/v3)

You can view Cherry Servers API docs here: [https://api.cherryservers.com/doc](https://api.cherryservers.com/doc)

## Table of Contents

- [Installation](#installation)
- [Authentication](#authentication)
- [Examples](#examples)
  - [Get teams](#get-teams)
  - [Get projects](#get-projects)
  - [Get plans](#get-plans)
  - [Get images](#get-images)
  - [Request new server](#request-new-server)

## Installation

Download the library to your GOPATH:
```
go get github.com/cherryservers/cherrygo/v3
```

Then import the library in your Go code:
```
import "github.com/cherryservers/cherrygo/v3
```

### Authentication

To authenticate to the Cherry Servers API, you must have an API token. You can create API tokens in the [Cherry Servers client portal](https://portal.cherryservers.com). Tokens must be exported in the `CHERRY_AUTH_TOKEN` environment variable or passed to the client directly.

Use an exported CHERRY_AUTH_TOKEN environment variable:
```
export CHERRY_AUTH_TOKEN="4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
```
```go
func main() {
    c, err := cherrygo.NewClient()
}
```
Pass a token directly to the client:
```go
func main() {
    c, err := cherrygo.NewClient(cherrygo.WithAuthToken("your-api-token"))
}
```

### Examples

#### Get teams
You will need a team ID for subsequent function calls, for example, to get projects for a specified team, you will need to provide a team ID.
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
View available server plans.

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
View OS images available for a specific plan.

```go
images, _, err := c.Images.List(planSlug, nil)
if err != nil {
    log.Fatal("Error", err)
}

for _, i := range images {
    log.Println(i.Name, i.Slug)
}
```

#### Request new server
```go
addServerRequest := cherrygo.CreateServer{
    ProjectID:   projectID,
    Image:       imageSlug,
    Region:      regionSlug,
    Plan:        planSlug,
}

server, _, err := c.Servers.Create(&addServerRequest)
if err != nil {
    log.Fatal("Error while creating new server: ", err)
}

log.Println(server.ID, server.Name, server.Hostname)
```

## Debug

If you want to debug this library, set the CHERRY_DEBUG environment variable to true, which enable full API request and response logging.
```
export CHERRY_DEBUG="true"
```

Unset the variable to stop debugging.
```
unset CHERRY_DEBUG
```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations.
