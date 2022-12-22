# go-everypay

[![Go Reference](https://pkg.go.dev/badge/github.com/cesbo/go-everypay.svg)](https://pkg.go.dev/github.com/cesbo/go-everypay)

Simple library for [EveryPay](https://every-pay.com) payment gateway

## Installation

To install the library use the following command in the project directory:

```
go get github.com/cesbo/go-everypay
```

## Quick Start

```go
e := everypay.NewEverypay(
    "api username",
    "api secret",
    "EUR3D1",
    false,
)

link, err := e.InitialPayment(&everypay.OneOff{
    Amount: 100.33,
    CustomerUrl: "https://example.com/thank-you",
    OrderReference: "order-id",
    Description: "example order",
    CustomerEmail: "name@example.com",
    CustomerIp: "127.0.0.1",
})
```
