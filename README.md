# gin-scalar

A gin middleware that helps automatically generate RESTful API documentation with Scalar.

> [!NOTE]
> Unlike Swagger implementations, this middleware bundles Scalar, and it does not need a separate files package.

## Usage

### Start using it

1. Add comments to your API source code, [See Declarative Comments Format](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format).
2. Download [Swag](https://github.com/swaggo/swag) for Go with:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

3. Run [Swag](https://github.com/swaggo/swag) at your Go project root path, [Swag](https://github.com/swaggo/swag) will parse comments and generate the required files (`docs` folder and `docs/doc.go`).

```sh
swag init
```

4. Download [gin-scalar](https://github.com/tinyauthapp/gin-scalar) with:

```sh
go get -u github.com/tinyauthapp/gin-scalar
```

Import the middleware:

```go
import "github.com/tinyauthapp/gin-scalar"
```

### Canonical example

Now assume you have implemented a simple API as follows:

```go
// A get function that returns a hello world string by JSON
type HelloWorldResponse struct {
    Message string `json:"message"`
}

func HelloWorld(ctx *gin.Context)  {
   ctx.JSON(http.StatusOK, HelloWorldResponse{
   Message: "Hello, World!",
   })
}
```

Now, add the scalar middleware to your gin router:

1. Add comments in the API and in the main function:

```go
type HelloWorldResponse struct {
	Message string `json:"message"`
}
// @BasePath /api/v1

// HelloWorldExample godoc
// @Summary Hello World Example
// @Description Just return a hello world string
// @Tags example
// @Produce json
// @Success 200 {object} HelloWordResponse
// @Router /example/helloworld [get]
func HelloWorld(ctx *gin.Context)  {
    ctx.JSON(http.StatusOK, HelloWorldResponse{
		Message: "Hello, World!",
    })
}
```

2. Use the `swag init` command to generate the docs, the generated docs will be stored at `docs/`.

3. Add the scalar middleware to your gin router:

```go
package main

import (
   "github.com/gin-gonic/gin"
   docs "github.com/go-project-name/docs"
   ginScalar "github.com/tinyauthapp/gin-scalar"
   "net/http"
)
type HelloWorldResponse struct {
   Message string `json:"message"`
}
// @BasePath /api/v1

// HelloWorldExample godoc
// @Summary Hello World Example
// @Description Just return a hello world string
// @Tags example
// @Produce json
// @Success 200 {object} HelloWorldResponse
// @Router /example/helloworld [get]
func HelloWorld(ctx *gin.Context)  {
   ctx.JSON(http.StatusOK, HelloWorldResponse{
      Message: "Hello, World!",
   })
}

func main()  {
   r := gin.Default()
   docs.SwaggerInfo.BasePath = "/api/v1"
   v1 := r.Group("/api/v1")
   {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld", HelloWorld)
      }
   }
   r.GET("/scalar/*any", ginScalar.WrapHandler(nil))
   r.Run(":8080")
}
```

Demo project tree, `swag init` is run at the root path of the project.

```
.
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
└── main.go
```

## Configuration

You can configure Scalar using different configuration options. For example:

```go
ginScalar.WrapHandler(nil, ginScalar.URL("http://localhost:8080/swagger/doc.json"))
```

| Option       | Type   | Default    | Description                                                                                                                                                                                                                                                |
| ------------ | ------ | ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| URL          | string | "doc.json" | URL pointing to API definition (normally swagger.json or swagger.yaml)                                                                                                                                                                                     |
| InstanceName | string | "swagger"  | The instance name of the swagger document. If multiple different swagger instances should be deployed on one gin router, ensure that each instance has a unique name (use the _--instanceName_ parameter to generate swagger documents with _swag init_).  |
| BasePath     | string | "/scalar"  | The base path in which the Scalar UI will be served                                                                                                                                                                                                        |
| ProjectName  | string | "Scalar"   | Project name displayed as the page title of the Scalar UI                                                                                                                                                                                                  |

## License

Licensed under [MIT](./LICENSE).