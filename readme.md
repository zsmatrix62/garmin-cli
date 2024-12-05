## Install

```bash
go install github.com/zsmatrix62/garmin-cli@latest
```

## Basic Usage

### Arguments

```bash
garmin-cli \
    -u <username> \
    -p <password> \
    -w <flow> \
    --persist_state=[true, true|false] \
    --save_state="./garmin-cli-states" \
    [flow args]
```

### Examples

Upload Fit file to Garmin Connect

```bash
garmin-cli -u whoami@gmail.com -p iamapassword -w upload-fit 2024-11-27-19-06-00.fit
```
