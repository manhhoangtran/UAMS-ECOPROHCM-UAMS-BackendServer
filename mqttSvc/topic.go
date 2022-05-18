package mqttSvc

const (
	TOPIC_GW_LOG_C           string = "gateway/log/create"
	TOPIC_GW_DOORLOCK_STATUS string = "gateway/doorlock/status"
	TOPIC_GW_DOORLOCK_C      string = "gateway/doorlock/create"
	TOPIC_GW_DOORLOCK_U      string = "gateway/doorlock/update"
	TOPIC_GW_DOORLOCK_D      string = "gateway/doorlock/delete"

	TOPIC_GW_BOOTUP   string = "gateway/bootup"
	TOPIC_GW_SHUTDOWN string = "gateway/shutdown"
	TOPIC_GW_LASTWILL string = "gateway/lastwill"

	TOPIC_SV_DOORLOCK_C      string = "server/doorlock/create"
	TOPIC_SV_DOORLOCK_U      string = "server/doorlock/update"
	TOPIC_SV_DOORLOCK_D      string = "server/doorlock/delete"
	TOPIC_SV_DOORLOCK_CMD    string = "server/doorlock/command"
	TOPIC_SV_DOORLOCK_BOOTUP string = "server/doorlock/bootup"

	TOPIC_SV_GATEWAY_U string = "server/gateway/update"
	TOPIC_SV_GATEWAY_D string = "server/gateway/delete"

	TOPIC_SV_SCHEDULER_C      string = "server/register/create"
	TOPIC_SV_SCHEDULER_U      string = "server/register/update"
	TOPIC_SV_SCHEDULER_D      string = "server/register/delete"
	TOPIC_SV_SCHEDULER_BOOTUP string = "server/register/bootup"

	TOPIC_SV_HP_BOOTUP string = "server/hp/bootup"
	TOPIC_SV_HP_C      string = "server/hp/create"
	TOPIC_SV_HP_U      string = "server/hp/update"
	TOPIC_SV_HP_D      string = "server/hp/delete"

	TOPIC_SV_USER_U string = "server/user/update"
	TOPIC_SV_USER_D string = "server/user/delete"

	TOPIC_SV_SYSTEM_U      string = "server/system/update"
	TOPIC_SV_LASTWILL      string = "server/lastwill"
	TOPIC_SV_SYSTEM_BOOTUP string = "server/system/bootup"
)
