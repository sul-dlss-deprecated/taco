# TACO Developer Guide

## Adding Envirment Variables for Configuration

```shell
export YOUR_VAR='whatever'
```

### Add config and default value

1. Add to the config struct definition
```go
type Config struct {
    AWS_Region	string
		DB_Endpoint string
		Disable_SSL bool
    Your_Var    string
}
```

2. Load the config and default
```go
func your_var() string {
	var yourvar string
	yourvar = os.Getenv("YOUR_VAR")
	if endpoint == "" {
		yourvar = "Some Default Value"
		log.Printf("YOUR_VAR: Using default [Some Default Value].")
	}
	log.Printf("YOUR_VAR: Found setting [%s]", yourvar)
	return yourvar
}
```

3. Load var into config struct
```go
func NewConfig() *Config {
	return &Config{
		AWS_Region: aws_region(),
		DB_Endpoint: db_endpoint(),
		Disable_SSL: disable_ssl(),
    Your_Var: your_var(),
	}
}
```
