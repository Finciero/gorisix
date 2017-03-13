<img align="right" src="https://cdn-images-1.medium.com/max/800/1*NlSd8_zML3LljxAypRw64w.png">
# Gorisix
**Gorisix** is an open source tool belt to integrate your web service with [46degrees](http://www.46degrees.net/)'s API.


## Overview

* Lightweight
* Extensible
* Zero dependencies

## Instalation

```
go get github.com/Finciero/gorisix
```

## How to use

To use **Gorisix** you only need to create a new client.

```go
import "github.com/Finciero/gorisix"

// ...

gc := gorisix.NewClient("secret", "receiver_id")

```

## Custom request

If you want to create a new payment url using 46degrees's API, you can use **Gorisix** as follow:

```go
import (
    "http"

    "github.com/Finciero/gorisix"
)

// create a new payment
payment := gosirix.Payment{
  ...
}

resp, err := gc,CreatePayment(payment)
if err != nil {
  // do something with this error
}

```
`resp` has all the information to lead your user through 46degrees's APP.
