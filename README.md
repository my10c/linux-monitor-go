
## See check for their release status
## Current 	: working on next check
## Next    	: check-process with stats
## Released ; check-mysql, check-kafka (broker only) : all minus stats
## Released ; check-disks, check-memory, check-load

## linux-monitor-go : Linux monitor programs written in GO

### Background
I been writting nagios plugins for years (who remember netsaint?), most in
bash, perl and python and now I'm planning to re-write all these in Go and
make them public

The reason of my choice for Go, is simple, I wanted a single binary and able
to use a single configuration file. A side effect is you could use this code
as a nagios-plugin framework without to have re-investing the wheel :)

### Single configuration
An other thing was the check's flags, sometime there are a lot of them, so now instead
of having these given on the command-line, they are now defined in the configuration file,
Example `warning` instead of `-w value`, just set `warning: value` in configuration file.

My choice for the file format is yaml. The reason is very simple, is easy to read
and to create. Given the flag `-setup`, the check will show you what are the available
configurations such as threshlold name. Some flags can be query with the keyword `help`,
example `-mode help` to show all the valid modes and the required configuration keys name.

For each check create a section `check-name:` and under it add the configuration value,
such warning- and critical-thresholds.
example

```
	check-momo:
	  user: momo
	  password: momo
	  warning: 10%
	  critical: 5%
```

#### Trick
With the flag `-config` you can use the same check for different needs. An other neat trick
is that could copy the same binary and now you can have single configuration for the same check.
Example:
	ln -s check-a check-b
	ln -s check-a check-c
the configuration would then look like this
```
check-a:
  username: momo
check-b:
  username: mimi
check-c:
  username: mumu
```

configuration that applies to all checks:
```
	# When enable stats (-stats) then the keys below are required, below are the default values.
	  	statsdir: /var/log/nagios-plugins-go.stat
	  	statsfile: {check-name}.stat
	# How many iteration to perform and how much to wait between (seconds) them before its an issue.
	  	iter: 3
	  	iterwait: 10
```

Here are the shared configurarion:
seems like a lot, but you should only need to configure these once, or disable the one
you do not care about or use the default.

`Values shown are the default values. Any section can be ommited, it will then use the default values.`
```
if the value a the key is shown empty and the key is used to disable the section, then the section
is by default disabled
```
```
	common:
	  nolog: false
	  debug: false
	  noalert: false
	  stats: false
	# to disable set an empty `logfile`, if shown empty, then its disable by default.
	log:
	  logmaxbackups: 14
	  logmaxage: 14
	  logdir: /var/log/nagios-plugins-go
	  logfile: check-mysql.log
	  logmaxsize: 128
	# to disable set an empty `emailto`, if shown empty, then its disable by default.
	email:
	  emailfrom:
	  emailfromname:
	  emailsubjecttag: [MONITOR]
	  emailpass:
	  emailhostport: 25
	  emailto:
	  emailtoname:
	  emailuser:
	  emailhost: localhost
	# to disable set an empty `tagfile`, if shown empty, then its disable by default.
	tag:
	  tagfile:
	  tagkeyname:
	# to disable set `syslogtag: off`, if shown empty, then its disable by default.
	syslog:
	  syslogtag: [{check-name}]
	  syslogpriority: LOG_INFO
	  syslogfacility: LOG_SYSLOG
	# to disable set an empty `pdservicekey`, if shown empty, then its disable by default.
	pagerduty:
	  pdservicename:
	  pdvalidunit: hour
	  pdevent: MONITOR ALERT
	  pdservicekey:
	# to disable set an empty `slackservicekey`, if shown empty, then its disable by default.
	slack:
	  slackservicekey:
	  slackchannel:
	  slackuser: MONITOR
```

NOTE
```
	* The key must be all lowercase!
	* Any key value that contains any of these charaters: ':#[]()*' must be double quoted!
	* The key `logmaxsize` value unit is megabytes.
	* The `tagfile` and `tagkeyname` keys are use to get a tag; useful in AWS, info by looking for
	  the keyword `tagkeyname` in the configured file `tagfile`, line format: 'keyname value', nothing fancy!
	* The pagerduty `pdvalidunit` is the unit used to create an event-id so no duplicate is created.
	  Valid choices are hour or minute. If an event was create at hour X (or minute X) then pagerduty
	  will not create a new event until the next hour, it sees it as an update to an existing event,.
	  because it has the same event-id, but do realize there always the possiblity that it could
	  overlap, certainly if it set to minute, you could get alert every minute!.
	  If the `pdvalidunit` value is invalid then it defaults to hour, valid options are `hour` or `minute`.
	* The key `emailsubjecttag` is use for email filtering.
	* Syslog Valid `syslogpriority`: LOG_INFO LOG_EMERG LOG_ALERT LOG_CRIT LOG_ERR LOG_NOTICE LOG_WARNING LOG_DEBUG
	* Syslog Valid `syslogfacility`: LOG_AUTH LOG_NEWS LOG_LOCAL7 LOG_MAIL LOG_DAEMON LOG_SYSLOG
		LOG_UUCP LOG_LOCAL5 LOG_LPR LOG_AUTHPRIV LOG_LOCAL3 LOG_LOCAL4 LOG_LOCAL6
		LOG_CRON LOG_FTP LOG_LOCAL0 LOG_LOCAL1 LOG_LOCAL2
```


### Checks
This is the list of check I plan to build:

check-cert status `not started yet` 	: check cert expiration

check-fd status `not started yet` 		: check file descriptors

check-http status `not started yet`		: check http port reply

check-disk status `first release`		: check disk space

check-kafka status `broker mode only`	: check kakfa

check-load status `first release`		: check system load

check-memory status `first release`		: check available memory

check-mysql status `first release`		: check mysql health include slave/replication

check-network status `back in queue`	: check network status such as TX, RX and error

check-nginx status `not started yet`	: check nginx status

check-process status `starting`			: check if given process(es) is running, /proc basesd

check-psql status `not started yet`		: check mysql health include slave/replication

Any other that you would like to see? shoot me an email


*note*: check-kafka checks to be finished : topic and pubsub

### How to build

create a work directory then set GOPATH : export GOPATH=full-path-work-directory

```
go get github.com/my10c/linux-monitor-go/xxxx
with xxxx the name of the check/script
```

that's it


### Feedback
Feedback and bug report welcome...

Enjoy, Momo
