# SheepsTor
Utility for updating static websites (served via Hugo) from Github

## Configuration
Sheepstor is configured from two places:

1. Some environment variables (see the `.env` file). These will need to ne set and exported in the runtime environment.
2. The config file `./config/config.yaml`


## Run as web service
### Using defaults
```bash
./sheepstor server
```

### With debugging
```bash
./sheepstor server --debug=true
```


## Run as command line utility
### Using defaults - single website
```bash
./sheepstor update --sites=www.antleaf.com
```

### Using defaults - update all websites
(also useful for InitContainer to set up web service)
```bash
./sheepstor update --sites=all
```

### With debugging
```bash
./sheepstor update --sites=all --debug=true
```
