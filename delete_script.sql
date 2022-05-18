
IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[schedulers]') AND type in (N'U'))
DROP TABLE [dbo].[schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[doorlocks]') AND type in (N'U'))
DROP TABLE [dbo].[doorlocks]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gw_networks]') AND type in (N'U'))
DROP TABLE [dbo].[gw_networks]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateways]') AND type in (N'U'))
DROP TABLE [dbo].[gateways]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[areas]') AND type in (N'U'))
DROP TABLE [dbo].[areas]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[employees]') AND type in (N'U'))
DROP TABLE [dbo].[employees]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[students]') AND type in (N'U'))
DROP TABLE [dbo].[students]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[customers]') AND type in (N'U'))
DROP TABLE [dbo].[customers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[passwords]') AND type in (N'U'))
DROP TABLE [dbo].[passwords]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateway_logs]') AND type in (N'U'))
DROP TABLE [dbo].[gateway_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[secret_keys]') AND type in (N'U'))
DROP TABLE [dbo].[secret_keys]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[doorlock_status_logs]') AND type in (N'U'))
DROP TABLE [dbo].[doorlock_status_logs]

GO
