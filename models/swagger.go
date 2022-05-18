package models

type SwagCreateGateway struct {
	AreaID    uint   `json:"areaId"`
	GatewayID string `json:"gatewayId"`
	Name      string `json:"name"`
}

type SwagUpateGateway struct {
	GormModel
	SwagCreateGateway
}

type SwagCreateArea struct {
	Gateway Gateway `json:"gateway"`
	Name    string  `json:"name"`
	Manager string  `json:"manager"`
}

type SwagUpdateArea struct {
	GormModel
	SwagCreateArea
}

type SwagCreateDoorlock struct {
	AreaID          uint   `json:"areaId"`
	GatewayID       uint   `json:"gatewayId"`
	SchedulerID     uint   `json:"schedulerId"`
	Description     string `json:"description"`
	Location        string `json:"location"`
	DoorlockAddress string `json:"doorlockAddress"`
}

type SwagUpdateDoorlock struct {
	GormModel
	ActiveState     string `json:"activeState"`
	BlockId         string `json:"blockId"`
	Description     string `json:"description"`
	DoorSerialID    string `json:"doorSerialId"`
	DoorlockAddress string `json:"doorlockAddress"`
	GatewayID       string `json:"gatewayId"`
	FloorId         string `json:"floorId"`
	RoomId          string `json:"roomId"`
	Location        string `json:"location"`
}

type SwagUpdatePassword struct {
	GormModel
	GatewayID    string `json:"gatewayId"`
	PasswordType string `json:"passwordType"`
	PasswordHash string `json:"passwordHash"`
}

type SwagCreateCustomer struct {
	CCCD  string `json:"cccd"  binding:"required"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	UserPass
}

type SwagCreateStudent struct {
	MSSV  string `json:"mssv"  binding:"required"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Major string `json:"major"`
	UserPass
}

type SwagCreateEmployee struct {
	MSNV            string `json:"msnv" binding:"required"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Department      string `json:"department"`
	Role            string `json:"role"`
	HighestPriority bool   `json:"highestPriority"`
	UserPass
}
type SwagCreateScheduler struct {
	Base           string `json:"base"`
	RoomRow        string `json:"roomRow"`
	RoomID         string `json:"roomId"`
	RoomName       string `json:"roomName"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	ClassID        string `json:"classId"`
	ClassName      string `json:"className"`
	LecturerID     string `json:"lecturerId"`
	LecturerName   string `json:"lecturerName"`
	Capacity       uint   `json:"capacity"`
	WeekDay        uint   `json:"weekDay"`
	StartClassTime uint   `json:"startClassTime"`
	EndClassTime   uint   `json:"endClassTime"`
	Amount         uint   `json:"amount"`
	Status         string `json:"status"`
}

type SwagUpdateScheduler struct {
	ID             uint   `json:"id"`
	Base           string `json:"base"`
	RoomRow        string `json:"roomRow"`
	RoomID         string `json:"roomId"`
	RoomName       string `json:"roomName"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	ClassID        string `json:"classId"`
	ClassName      string `json:"className"`
	LecturerID     string `json:"lecturerId"`
	LecturerName   string `json:"lecturerName"`
	Capacity       uint   `json:"capacity"`
	WeekDay        uint   `json:"weekDay"`
	StartClassTime uint   `json:"startClassTime"`
	EndClassTime   uint   `json:"endClassTime"`
	Amount         uint   `json:"amount"`
	Status         string `json:"status"`
}

type SwagCreateSecretKey struct {
	SecretKey string `json:"secret"`
}

type SwaggerDoorlockCmd struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	Duration string `json:"duration"`
}
