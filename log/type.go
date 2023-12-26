package log

import "time"

const (
	LabelStartup  string = "STARTUP   "
	LabelMonitor  string = "MONITOR   "
	LabelShutdown string = "SHUTDOWN  "
	LabelQueue    string = "QUEUE	 "
	LabelMPC      string = "MPC	 "

	levelInfo  string = "[INFO ]"
	levelError string = "[ERROR]"
	levelFatal string = "[FATAL]"
)

var timezone *time.Location = time.Local
