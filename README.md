# SheepsTor
Utility for updating static websites (served via Hugo) from Github

## Run as web service
### Using defaults
```bash
./sheepstor
```

### With configuration options
```bash
./sheepstor --debug=true --config=./config/config.yaml
```


## Run as command line utility
### Using defaults - single website
```bash
./sheepstor --update=www.antleaf.com
```

### Using defaults - update all websites
(also useful for InitContainer to set up web service)
```bash
./sheepstor --update=all
```

### With configuration options
```bash
./sheepstor --update=all --debug=true --config=./config/config.yaml
```
