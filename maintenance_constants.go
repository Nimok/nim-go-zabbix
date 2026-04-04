package zabbix

// Maintenance type values (`maintenance_type`).
const (
	// Aliases for existing host maintenance type constants.
	MaintenanceTypeWithData = MaintenanceWithData
	MaintenanceTypeNoData   = MaintenanceNoData
)

// Maintenance tag evaluation type values (`tags_evaltype`).
const (
	MaintenanceTagsEvalTypeAndOr = 0
	MaintenanceTagsEvalTypeOr    = 2
)

// Maintenance time period type values (`timeperiod_type`).
const (
	MaintenanceTimePeriodOneTime = 0
	MaintenanceTimePeriodDaily   = 2
	MaintenanceTimePeriodWeekly  = 3
	MaintenanceTimePeriodMonthly = 4
)

// Maintenance problem tag operator values (`operator`).
const (
	MaintenanceTagOperatorEquals   = 0
	MaintenanceTagOperatorContains = 2
)

// Maintenance weekly/monthly day-of-week bitmask values (`dayofweek`).
const (
	MaintenanceDayMonday    = 1
	MaintenanceDayTuesday   = 2
	MaintenanceDayWednesday = 4
	MaintenanceDayThursday  = 8
	MaintenanceDayFriday    = 16
	MaintenanceDaySaturday  = 32
	MaintenanceDaySunday    = 64
)

// Maintenance monthly month bitmask values (`month`).
const (
	MaintenanceMonthJanuary   = 1
	MaintenanceMonthFebruary  = 2
	MaintenanceMonthMarch     = 4
	MaintenanceMonthApril     = 8
	MaintenanceMonthMay       = 16
	MaintenanceMonthJune      = 32
	MaintenanceMonthJuly      = 64
	MaintenanceMonthAugust    = 128
	MaintenanceMonthSeptember = 256
	MaintenanceMonthOctober   = 512
	MaintenanceMonthNovember  = 1024
	MaintenanceMonthDecember  = 2048
	MaintenanceMonthAll       = 4095
)

// Maintenance monthly week selector values (`every` for monthly schedules).
const (
	MaintenanceWeekFirst  = 1
	MaintenanceWeekSecond = 2
	MaintenanceWeekThird  = 3
	MaintenanceWeekFourth = 4
	MaintenanceWeekLast   = 5
)

// Default period in seconds used by Zabbix when `period` is omitted.
const (
	MaintenanceDefaultPeriodSeconds = 3600
)
