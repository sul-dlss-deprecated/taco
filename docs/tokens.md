# Token Authorization

## Installing

On a mac:
```
brew tap mike-engel/jwt-cli
brew install jwt-cli
```

## Create a token

*note* the secret must be the same secret that TACO is configured with
(in the `TACO_SECRET_KEY` variable)
```
jwt encode --secret superSeckeret --iss hydrus.stanford.edu --sub "lmcrae@stanford.edu"
```

## Use the token

```
curl -H "Content-Type: application/json" \
 -H "Authorization: bearer <token>" \
 -d@examples/request.json http://localhost:8080/v1/resource
```
