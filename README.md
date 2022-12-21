# Doctor Doom

![Dr.Doom](./images/drdoom-removebg-preview.png)

## Description

Doctor Doom will destroy files or folder which are:

- Live longer than certain time (default 30 days)
- Have a size bigger than certain size (default 100MB)
- Have a certain name (default "*", all files)
Doctor Doom will alway find file victims in recursive way. It will not destroy the folder itself.
Doctor Doom will only destroy folder if it is empty (of course, match the conditions).

NOTE: With the default config, Doctor Doom will destroy all files.

## Environment

- `doom_path`: The root folder path, where Dr.Doom will look for files to destroy
- `circle`: The time interval (in time unit, integer ) between each Dr.Doom run. Cron tab definition ex: `0 0 * * 0` (every Sunday at midnight)
- `doom_export`: The path where Dr.Doom will export the list of files it destroyed. Default is `./doom_victims.log`

## Rules (These rule will alway use the OR logic)

- `age`: The time (in time unit) a file must be older than to be destroyed. Default is 30d
  - `d`: day
  - `h`: hour
  - `m`: minute
- `size`: The size (in size unit) a file must be bigger than to be destroyed. Default is 100MB
  - `B`: byte
  - `K`: kilobyte
  - `M`: megabyte
  - `G`: gigabyte
  - `T`: terabyte
- `name`: The name of the file to be destroyed. Default is `.*`. Can use regex

### Rule priority

- `Age` > `Size` > `Name`

## Example 1

```yaml
# Example 1
doom_path: /home/user
circle: 0 0 * * 0
doom_export: /home/user/doom_victims.log
rule:
  age: 30d
  size: 100M
  name: "*"
rule_and: false

# Meaning: Dr.Doom will destroy all files in /home/user that are:
# - Live longer than 30 days
# - Have a size bigger than 100MB
# - Have a name that matches regex .*
# - Have a extension that matches regex .txt
# Dr.Doom will run every Sunday at midnight and export the list of files it destroyed to /home/user/doom_victims.log
# The destroy process will use the OR logic between the rules
# Using the OR logic between the rules
```

## Example 2

```yaml
# Example 2
doom_path: /home/user
circle: * 14 * * 0
doom_export: /home/user/doom_victims.log
rule:
  age: 30d
  size: 10M
  name: "/^victim/"
rule_and: true

# Meaning: Dr.Doom will destroy all files in /home/user that are:
# - Live longer than 30 days
# - Have a size bigger than 10MB
# - Have a name that matches regex /^victim/
# - Have a extension that matches regex .txt
# Dr.Doom will run every Sunday at 2pm and export the list of files it destroyed to /home/user/doom_victims.log
# Using the AND logic between the rules
```

## Override default value

**WHAT OVERRIDE WHAT**

`Default Config` << `Environment Config` << `File Config`

### Using environment variable

#### Docker container

```bash
docker run -d --name dr-doom -e DOOM_PATH="/home_user" \
-e CIRCLE="0 0 * * *" \
-e DOOM_EXPORT="/home_user/doom_victims.log" \
-e RULE_AGE="30d" -e RULE_SIZE="100M" \
-e RULE_NAME=".*" -v /home/user:/home_user \
-e RULE_AND="true" \
--restart unless-stopped \
mrnim94/doctor-doom:latest \
./doctor-doom
```

#### Docker compose

```yaml
version: "3.7"
services:
  dr-doom:
    image: mrnim94/doctor-doom:latest
    container_name: dr-doom
    environment:
      - DOOM_PATH="/home_user"
      - CIRCLE="0 0 * * *"
      - DOOM_EXPORT="/var/log/doctor-doom/doom_victims.log"
      - RULE_AGE="30d"
      - RULE_SIZE="100M"
      - RULE_NAME="*"
      - RULE_AND="true"
    volumes:
      - /home/user:/home_user
      - /var/log:/var/log
    restart: unless-stopped
```

### Using config file

```yaml
# ./sample/config.yaml
doom_path: /home/user
circle: "* 14 * * *"
doom_export: /home/user/doom_victims.log
rule:
  age: 30d
  size: 10M
  name: "/^victim/"
rule_and: true

```

Usage

```bash
./doctor-doom --doom-config ./sample/config.yaml

# Use this command as ENTRYPOINT or CMD
```

## Dependencies

- [Uber Zap](https://github.com/uber-go/zap)
- [Lumberjack v2](https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2?utm_source=godoc)
- [Cron v3](https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0)
- [YAML v2](https://pkg.go.dev/gopkg.in/yaml.v2@v2.4.0)
- [CLI v2](https://pkg.go.dev/github.com/urfave/cli/v2@v2.23.6)
