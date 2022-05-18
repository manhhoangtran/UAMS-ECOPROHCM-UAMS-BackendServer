# DMS_BackendServer

## How to run app
```bash
    go install github.com/google/wire/cmd/wire@latest #require
    wire ./initializers
    go run .
```

## How to access MSSQL from VSCode's SQL Server extension

1. Server name: `server host`, `mssql port`
2. DB name: DevDB
3. SQL Login
4. User name: sa
5. Pass: `env file`

## How to deploy

1. Install `Docker` on your local
2. Install `sshpass`
3. Check connection host env variables ( `SERVER_HOST` is for development only )
3. make deploy SERVER_PW= `contact server admin`
4. SSH to server: 
 - Run cmd to clear *old* ssh key: `ssh-keygen -R iot.hcmue.space`
 - `ssh -p 2223 sviot@iot.hcmue.space`
5. `cd ./iot && docker load -i ./dms-be.tar && docker-compose down && docker-compose up -d`

## How to deploy swagger
### NOTES: Use [swag](github.com/swaggo/swag/cmd/swag@v1.7.8) to do `make swagger` & avoid weird error => It'll break CI

1. Update new `gin-swagger` comments for new APIs
2. `make swagger` ( Remember to add $GOPATH to your $PATH env variable )
3. SSH to server
4. Copy content of `docs/swagger.yaml` into `iot/swagger-ui/doc/api.yaml`

## How to use Logger
### About logger
Logger is upper layer based on [logrus](https://github.com/sirupsen/logrus) framework. Although `logrus` is a powerful logging framework but its default supported formatter was not match with logging format (JSONFormatter, TextFormatter) for our project, so defined our own Logger APIs based on it with customized third party formatter will be more flexible and easy to manage.
 - **Logger level**: Panic, Fatal, Error, Warn, Info, Debug
 - **Logger component**: Server's components are those framework using in project such as Gin, MQTT, SQL, GORM,...

### Logger formatter
 - **With fields options**: `<time> - LEVEL - Msg=Message here     {Param: key1:value1 - key2:value2 - ...}`
 - **Without fields options**: `<time> - LEVEL - Msg=Message here`
### Using Logger APIs
#### **NOTES**: Using level Panic and Fatal will end process
1. Like `logrus` Logger include APIs for default log by level, and its will be without field format
 - `func Fatal(args ...interface{})`
 - `func Debug(args ...interface{})`
 - `func Info(args ...interface{}`
 - `func Error(args ...interface{})`
 - `func Warn(args ...interface{})`
 - `func Panic(args ...interface{})`
2. Logger provides APIs for log with or without fields
 - `func LogWithFields(comp ServerComponent, level log.Level, fields LoggerFields, message ...interface{})`
```
logger.LogWithFields(logger.DMSSERVER, logger.InfoLevel, logger.LoggerFields{
		"Name":         "Gateway",
		"Descriptions": "Manage Doorlock",
	}, "Message for gateway")
```
Output
```
<2022-03-31 20:19:47.5015587 +07:00> - INFO - Msg= Message for gateway   {Params: Descriptions=Manage Doorlock - Component=DMS_SERVER - Name=Gateway}
```
 - `func LogfWithFields(comp ServerComponent, level log.Level, fields LoggerFields, messageFormat string, args ...interface{})`
```
logger.LogfWithFields(logger.DMSSERVER, logger.InfoLevel, logger.LoggerFields{
		"Name":         "Gateway",
		"Descriptions": "Manage Doorlock",
	}, "Message for gateway ID %s", "112312")
```
Output
```
<2022-03-31 20:19:47.5025684 +07:00> - INFO - Msg= Message for gateway ID 112312   {Params: Component=DMS_SERVER - Name=Gateway - Descriptions=Manage Doorlock}
```
 - `func LogWithoutFields(comp ServerComponent, level log.Level, message ...interface{})`
```
logger.LogWithoutFields(logger.GINROUTER, logger.DebugLevel, "Failed to GET")
```
Output
```
<2022-04-01 00:27:09.2521683 +07:00> - DEBUG - Msg= Failed to GET   {Params: Component=GIN_ROUTER}
```
 - `func LogfWithoutFields(comp ServerComponent, level log.Level, messageFormat string, args ...interface{})`
```
logger.LogfWithoutFields(logger.GINROUTER, logger.InfoLevel, "Failed to GET gateway %s", "123123")
```
Output
```
<2022-04-01 00:27:09.2532097 +07:00> - INFO - Msg= Failed to GET gateway 123123   {Params: Component=GIN_ROUTER}
```
### Modify Logger
1. Can change GinLogger middleware in `logs/log.go` with different parameters
2. Can add caller function when at debug level
3. Add or modify third party Formatter in `logs/formatter.go`
4. Add rotate logs function if needed
5. Change or add more logger APIs

