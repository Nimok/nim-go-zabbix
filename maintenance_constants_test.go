package zabbix

import "testing"

func TestMaintenanceConstants(t *testing.T) {
	if MaintenanceTypeWithData != MaintenanceWithData {
		t.Fatalf("MaintenanceTypeWithData alias mismatch: got %d want %d", MaintenanceTypeWithData, MaintenanceWithData)
	}
	if MaintenanceTypeWithData != 0 {
		t.Fatalf("MaintenanceTypeWithData value mismatch: got %d", MaintenanceTypeWithData)
	}

	if MaintenanceTypeNoData != MaintenanceNoData {
		t.Fatalf("MaintenanceTypeNoData alias mismatch: got %d want %d", MaintenanceTypeNoData, MaintenanceNoData)
	}
	if MaintenanceTypeNoData != 1 {
		t.Fatalf("MaintenanceTypeNoData value mismatch: got %d", MaintenanceTypeNoData)
	}

	if MaintenanceTagsEvalTypeAndOr != 0 {
		t.Fatalf("MaintenanceTagsEvalTypeAndOr mismatch: got %d", MaintenanceTagsEvalTypeAndOr)
	}
	if MaintenanceTagsEvalTypeOr != 2 {
		t.Fatalf("MaintenanceTagsEvalTypeOr mismatch: got %d", MaintenanceTagsEvalTypeOr)
	}

	if MaintenanceTimePeriodOneTime != 0 ||
		MaintenanceTimePeriodDaily != 2 ||
		MaintenanceTimePeriodWeekly != 3 ||
		MaintenanceTimePeriodMonthly != 4 {
		t.Fatalf(
			"MaintenanceTimePeriod constants mismatch: got one=%d daily=%d weekly=%d monthly=%d",
			MaintenanceTimePeriodOneTime,
			MaintenanceTimePeriodDaily,
			MaintenanceTimePeriodWeekly,
			MaintenanceTimePeriodMonthly,
		)
	}

	if MaintenanceTagOperatorEquals != 0 || MaintenanceTagOperatorContains != 2 {
		t.Fatalf(
			"MaintenanceTagOperator constants mismatch: got equals=%d contains=%d",
			MaintenanceTagOperatorEquals,
			MaintenanceTagOperatorContains,
		)
	}

	if MaintenanceDayMonday != 1 || MaintenanceDaySunday != 64 {
		t.Fatalf("Maintenance day constants mismatch: got monday=%d sunday=%d", MaintenanceDayMonday, MaintenanceDaySunday)
	}

	if MaintenanceMonthJanuary != 1 || MaintenanceMonthDecember != 2048 || MaintenanceMonthAll != 4095 {
		t.Fatalf(
			"Maintenance month constants mismatch: got january=%d december=%d all=%d",
			MaintenanceMonthJanuary,
			MaintenanceMonthDecember,
			MaintenanceMonthAll,
		)
	}

	if MaintenanceWeekFirst != 1 || MaintenanceWeekLast != 5 {
		t.Fatalf("Maintenance week constants mismatch: got first=%d last=%d", MaintenanceWeekFirst, MaintenanceWeekLast)
	}

	if MaintenanceDefaultPeriodSeconds != 3600 {
		t.Fatalf("MaintenanceDefaultPeriodSeconds mismatch: got %d", MaintenanceDefaultPeriodSeconds)
	}
}
