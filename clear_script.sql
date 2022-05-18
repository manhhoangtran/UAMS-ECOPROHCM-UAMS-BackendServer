
IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[schedulers]') AND type in (N'U'))
DELETE FROM [dbo].[schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[doorlocks]') AND type in (N'U'))
DELETE FROM [dbo].[doorlocks]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateways]') AND type in (N'U'))
DELETE FROM [dbo].[gateways]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[areas]') AND type in (N'U'))
DELETE FROM [dbo].[areas]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[employees]') AND type in (N'U'))
DELETE FROM [dbo].[employees]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[students]') AND type in (N'U'))
DELETE FROM [dbo].[students]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[customers]') AND type in (N'U'))
DELETE FROM [dbo].[customers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[passwords]') AND type in (N'U'))
DELETE FROM [dbo].[passwords]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateway_logs]') AND type in (N'U'))
DELETE FROM [dbo].[gateway_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[secret_keys]') AND type in (N'U'))
DELETE FROM [dbo].[secret_keys]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[doorlock_status_logs]') AND type in (N'U'))
DELETE FROM [dbo].[doorlock_status_logs]

GO
