package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

const (
	DEFAULT_TIME_FORMAT       string        = "2006-01-02 15:04:05.999999999 -07:00" // Sync with SQL format
	DEFAULT_CLEAN_LOGS_PERIOD time.Duration = time.Hour * 24 * 7                     // 1 week
)

type GatewayLogTime struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
type GatewayLog struct {
	ID        uint      `gorm:"primarykey;"`
	GatewayID string    `json:"gatewayId"`
	LogType   string    `json:"-"` // info, warn, debug, error, fatal
	Content   string    `json:"content"`
	LogTime   time.Time `json:"logTime"`
	CreatedAt time.Time `swaggerignore:"true"`
}

type logsCleanTicker struct {
	ticker            *time.Ticker
	done              chan bool
	db                *gorm.DB
	period            time.Duration
	scheduleTimer     *time.Timer
	isScheduleRunning bool
}
type LogSvc struct {
	db          *gorm.DB
	cleanTicker *logsCleanTicker
}

func NewLogSvc(db *gorm.DB) *LogSvc {
	logSvc := &LogSvc{
		db:          db,
		cleanTicker: newCleanTicker(DEFAULT_CLEAN_LOGS_PERIOD, db),
	}
	logSvc.cleanTicker.start()
	return logSvc
}

func (ls *LogSvc) FindAllGatewayLog(ctx context.Context) (glList []GatewayLog, err error) {
	result := ls.db.Find(&glList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return glList, nil
}

func (ls *LogSvc) FindGatewayLogByID(ctx context.Context, id string) (gl *GatewayLog, err error) {
	result := ls.db.Preload("Doorlocks").First(&gl, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gl, nil
}

func (ls *LogSvc) CreateGatewayLog(ctx context.Context, gl *GatewayLog) (*GatewayLog, error) {
	if err := ls.db.Create(&gl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gl, nil
}

func (ls *LogSvc) UpdateGatewayLogCleanPeriod(p *GatewayLogTime) error {
	fmt.Println("PERIOD ", p.toTimeDuration())
	ls.cleanTicker.setPeriod(p.toTimeDuration()).restart()
	return nil
}

func (ls *LogSvc) FindGatewayLogsByTime(gatewayId string, from string, to string) (glList *[]GatewayLog, err error) {
	result := ls.db.Where("gateway_id = ? AND log_time >= ? AND log_time <= ?", gatewayId, from, to).Find(&glList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return glList, nil
}

func (ls *LogSvc) FindGatewayLogsTypeByTime(gatewayId string, logType string, from string, to string) (glList *[]GatewayLog, err error) {
	result := ls.db.Where("gateway_id = ? AND log_type = ? AND log_time >= ? AND log_time <= ?", gatewayId, strings.ToUpper(logType), from, to).Find(&glList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return glList, nil
}

func newCleanTicker(period time.Duration, db *gorm.DB) *logsCleanTicker {
	return &logsCleanTicker{
		db:                db,
		period:            period,
		isScheduleRunning: false,
	}
}

func (lc *logsCleanTicker) runBackground() {
	for {
		select {
		case <-lc.done:
			return
		case tick := <-lc.ticker.C:
			fmt.Printf("At %s Gateway logs cleaner start...\n", tick)
			currentTime := time.Now()
			beginOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
				4, 59, 59, 0, currentTime.Location())
			endOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
				23, 59, 59, 0, currentTime.Location())
			if currentTime.Before(beginOfDay) {
				lc.cleanGatewayLogs(&currentTime)
			} else {
				waitTime := endOfDay.Sub(currentTime)
				fmt.Println("Wait", waitTime, "(s)", "to clean up gateway logs")
				lc.setSchedule(waitTime, &currentTime).stop()
			}
		}
	}
}

func (lc *logsCleanTicker) start() {
	lc.ticker = time.NewTicker(lc.period)
	lc.done = make(chan bool)
	go lc.runBackground()
}

func (lc *logsCleanTicker) stop() {
	lc.ticker.Stop()
	lc.done <- true
	close(lc.done)
	fmt.Println("Gateway Logs clean Ticker Stopped!")
}

func (lc *logsCleanTicker) cleanGatewayLogs(current *time.Time) {
	timeOffset := current.Add(time.Duration(-lc.period))
	formatedTime := timeOffset.Format(DEFAULT_TIME_FORMAT)
	fmt.Printf("Remove Gateway Logs before %s\n", formatedTime)
	result := lc.db.Where("log_time <= ?", formatedTime).Delete(&GatewayLog{})
	if err := result.Error; err != nil {
		fmt.Println("Gateway Logs clean up failed")
	}
}

func (lc *logsCleanTicker) setPeriod(period time.Duration) *logsCleanTicker {
	lc.period = period
	return lc
}

func (lc *logsCleanTicker) restart() {
	if lc.isScheduleRunning {
		return
	}
	lc.stop()
	lc.start()
}

func (lc *logsCleanTicker) setSchedule(scheduleTime time.Duration, currentTime *time.Time) *logsCleanTicker {
	lc.newTimer(scheduleTime)
	go func() {
		<-lc.scheduleTimer.C
		lc.cleanGatewayLogs(currentTime)
		lc.setPeriod(lc.period).start()
		lc.isScheduleRunning = false
	}()
	return lc
}

func (lc *logsCleanTicker) newTimer(scheduleTime time.Duration) *logsCleanTicker {
	lc.scheduleTimer = time.NewTimer(scheduleTime)
	lc.isScheduleRunning = true
	return lc
}

func (glt *GatewayLogTime) toTimeDuration() time.Duration {
	timeOffset := time.Now()
	endTime := timeOffset.Add(time.Hour * 24 * time.Duration(glt.Day)).
		Add(time.Hour * time.Duration(glt.Hour)).
		Add(time.Minute * time.Duration(glt.Minute)).
		Add(time.Second * time.Duration(glt.Second))
	period := endTime.Sub(timeOffset)
	return period
}
