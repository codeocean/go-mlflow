# go-mlflow

go-mlflow is a Go client library for the [MLflow API](https://mlflow.org/docs/latest/rest-api.html), including user & permission management as described in the [MLFlow Authentication API](https://mlflow.org/docs/latest/auth/rest-api.html).

Currently, most APIs are covered, missing mainly registered models APIs.

## Installation

go-mlflow is compatible with modern Go releases in module mode, with Go installed:
```
go get github.com/codeocean/go-mlflow
```
will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:
```
import "github.com/codeocean/go-mlflow/mlflow"
```
and run go get without parameters.

## Usage

```
import "github.com/codeocean/go-mlflow/mlflow"
```

Construct a new MLflow client:
```
client, err := mlflow.NewClient(nil, "http://localhost:5000")
```
or if needed, with admin credentials:
```
client, err := mlflow.NewClient(nil, "http://admin-user:admin-password@localhost:5000")
```
then use the various services on the client to access different parts of the MLflow API. For example:
```
// Get an experiment
experiment, err := client.Experiments.Get(context.Background(), "1")
```

The services of a client divide the API into logical chunks and correspond to the structure of the MLflow API documentation at https://mlflow.org/docs/latest/rest-api.html .

NOTE: Using the context package, one can easily pass cancellation signals and deadlines to various services of the client for handling a request. In case there is no context available, then context.Background() can be used as a starting point.
