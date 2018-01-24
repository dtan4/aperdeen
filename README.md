# Aperdeen

[![Build Status](https://travis-ci.org/dtan4/aperdeen.svg?branch=master)](https://travis-ci.org/dtan4/aperdeen)
[![codecov](https://codecov.io/gh/dtan4/aperdeen/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/aperdeen)

Amazon API Gateway client and local proxy

## Commands

### `aperdeen endopoints APINAME`

List API endpoints

```sh-session
$ aperdeen endpoints dtan4
PATH                              ENDPOINT
/dtan4/*  https://terraforming.dtan4.net/*
/foobar/* arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:012345678912:function:api-backend/invocations
```

### `aperdeen export APINAME`

Export API endpoints to YAML

```sh-session
$ aperdeen export dtan4
name: dtan4
endpoints:
  /dtan4/*:
    url: http://terraforming.dtan4.net/*
  /foobar/*:
    url: arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:012345678912:function:api-backend/invocations
```

### `aperdeen local -f YAML`

Start local API Gateway

```sh-session
$ aperdeen local -f api.yaml
server started at :8080 ...

# on another session
$ curl -sL localhost:8080/dtan4/ | head
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; img-src data:; connect-src 'self'">
    <title>Site not found &middot; GitHub Pages</title>
    <style type="text/css" media="screen">
      body {
        background-color: #f1f1f1;
        margin: 0;
```

## Aperdeen YAML

```yaml
name: dtan4
endpoints:
  /dtan4/*:
    url: http://terraforming.dtan4.net/*
  /foobar/*:
    url: arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:012345678912:function:api-backend/invocations
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
