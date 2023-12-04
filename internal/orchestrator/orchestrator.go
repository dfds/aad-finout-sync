package orchestrator

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/handler"
	"go.dfds.cloud/aad-finout-sync/internal/util"

	configUtils "go.dfds.cloud/utils/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

const AzureAdToFinoutName = "aadToFinout"
const CostCentreToFinoutName = "costCentreToFinout"

var currentJobsGauge prometheus.Gauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name:      "jobs_running",
	Help:      "Current jobs that are running",
	Namespace: "aad_finout_sync",
})

var currentJobStatus *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "job_is_running",
	Help:      "Is {job_name} running. 1 = in progress, 0 = not running",
	Namespace: "aad_finout_sync",
}, []string{"name"})

var jobFailedCount *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "job_failed_count",
	Help:      "How many times has {job_name} failed.",
	Namespace: "aad_finout_sync",
}, []string{"name"})

var jobSuccessfulCount *prometheus.GaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "job_success_count",
	Help:      "How many times has {job_name} successfully completed.",
	Namespace: "aad_finout_sync",
}, []string{"name"})

// Orchestrator
// Used for managing long-lived fully fledged sync jobs
type Orchestrator struct {
	status     map[string]*SyncStatus
	scheduling map[string]*Schedule
	Jobs       map[string]*Job
	ctx        context.Context
	wg         *sync.WaitGroup
}

type SyncStatus struct {
	mu     sync.Mutex
	active bool
}

func (s *SyncStatus) InProgress() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}

func (s *SyncStatus) SetStatus(status bool) {
	s.mu.Lock()
	s.active = status
	s.mu.Unlock()
}

type Schedule struct {
	name         string
	enabled      bool
	interval     time.Duration
	lastExecuted time.Time
}

func (s *Schedule) Enabled() bool {
	return s.enabled
}

func (s *Schedule) Interval() time.Duration {
	return s.interval
}

func (s *Schedule) LoadConfig() {
	prefix := fmt.Sprintf("AFS_SCHEDULER_JOB_%s", strings.ToUpper(s.name))
	configPath := fmt.Sprintf("%s_ENABLE", prefix)
	s.enabled = configUtils.GetEnvBool(configPath, false)

	configPath = fmt.Sprintf("%s_INTERVAL", prefix)
	s.interval = parseDuration(configUtils.GetEnvValue(configPath, "1h"))

	util.Logger.Info(fmt.Sprintf("Job schedule %s loaded with the following configuration: Enabled: %t, Interval(in seconds): %d", s.name, s.Enabled(), int64(s.Interval().Seconds())))
}

func (s *Schedule) TimeToRun() bool {
	now := time.Now().Unix()
	nextExecution := s.lastExecuted.Add(s.interval).Unix()
	diff := nextExecution - now
	if diff < 0 {
		return true
	} else {
		return false
	}
}

func NewOrchestrator(ctx context.Context, wg *sync.WaitGroup) *Orchestrator {
	return &Orchestrator{
		Jobs:       map[string]*Job{},
		scheduling: map[string]*Schedule{},
		status:     map[string]*SyncStatus{},
		ctx:        ctx,
		wg:         wg,
	}
}

func (o *Orchestrator) AddJob(job *Job, schedule *Schedule) {
	o.status[job.Name] = &SyncStatus{active: false}
	o.scheduling[job.Name] = schedule
	job.context = o.ctx
	job.wg = o.wg
	job.Status = o.status[job.Name]
	job.Schedule = schedule

	schedule.name = job.Name
	schedule.LoadConfig()
	schedule.lastExecuted = time.Now().Add(-schedule.interval)

	o.Jobs[job.Name] = job
}

func (o *Orchestrator) JobStatus(name string) *SyncStatus {
	return o.status[name]
}

func (o *Orchestrator) JobStatusProgress(name string) bool {
	if status, exists := o.status[name]; !exists {
		return false // ideally this should never happen
	} else {
		return status.InProgress()
	}
}

func (o *Orchestrator) Init(conf config.Config) {
	o.AddJob(&Job{
		Name:    AzureAdToFinoutName,
		handler: handler.Azure2FinoutHandler,
	}, &Schedule{})

	o.AddJob(&Job{
		Name:    CostCentreToFinoutName,
		handler: handler.CostCentre2FinoutHandler,
	}, &Schedule{})
}

type Job struct {
	Name     string
	Status   *SyncStatus
	context  context.Context
	handler  func(ctx context.Context) error
	wg       *sync.WaitGroup
	Schedule *Schedule
}

func (j *Job) Run() {
	if j.Status.InProgress() {
		util.Logger.Warn("Can't start Job because Job is already in progress.", zap.String("jobName", j.Name))
		return
	}
	j.Schedule.lastExecuted = time.Now()
	j.Status.SetStatus(true)
	j.wg.Add(1)
	currentJobsGauge.Inc()
	currentJobStatus.WithLabelValues(j.Name).Set(1)
	util.Logger.Warn("Job started", zap.String("jobName", j.Name))

	go func() {
		defer j.wg.Done()
		err := j.handler(j.context)
		if err != nil {
			jobFailedCount.WithLabelValues(j.Name).Inc()
			util.Logger.Error("Job failed", zap.String("jobName", j.Name), zap.Error(err))
		} else {
			jobSuccessfulCount.WithLabelValues(j.Name).Inc()
		}
		currentJobsGauge.Dec()
		currentJobStatus.WithLabelValues(j.Name).Set(0)
		j.Status.SetStatus(false)
		util.Logger.Warn("Job ended", zap.String("jobName", j.Name))

	}()
}

func fakeJobGen(ctx context.Context, msg string, jobName string) {
	entityCount := rand.Intn(100-20) + 20
	sleepInSecs := rand.Intn(7-2) + 2

	for i := 0; i < entityCount; i++ {
		select {
		case <-ctx.Done():
			util.Logger.Info("Job cancelled", zap.String("jobName", jobName))
			return
		default:
			util.Logger.Info(fmt.Sprintf("%s %d", msg, i), zap.String("jobName", jobName))
			time.Sleep(time.Second * time.Duration(sleepInSecs))
		}
	}
}

func parseDuration(val string) time.Duration {
	duration, err := time.ParseDuration(val)
	if err != nil {
		panic(err) // ideally this should never happen
	}
	return duration
}
