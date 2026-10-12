package zabbix

import "testing"

func TestItemConstants(t *testing.T) {
	if ItemTypeZabbixAgent != 0 || ItemTypeTrapper != 2 || ItemTypeCalculated != 15 ||
		ItemTypeDependent != 18 || ItemTypeHTTPAgent != 19 || ItemTypeSNMPAgent != 20 ||
		ItemTypeScript != 21 || ItemTypeBrowser != 22 {
		t.Fatalf(
			"ItemType constants mismatch: got agent=%d trapper=%d calculated=%d dependent=%d http=%d snmp=%d script=%d browser=%d",
			ItemTypeZabbixAgent,
			ItemTypeTrapper,
			ItemTypeCalculated,
			ItemTypeDependent,
			ItemTypeHTTPAgent,
			ItemTypeSNMPAgent,
			ItemTypeScript,
			ItemTypeBrowser,
		)
	}

	if ItemValueTypeFloat != 0 || ItemValueTypeChar != 1 || ItemValueTypeLog != 2 ||
		ItemValueTypeUnsigned != 3 || ItemValueTypeText != 4 || ItemValueTypeBinary != 5 {
		t.Fatalf(
			"ItemValueType constants mismatch: got float=%d char=%d log=%d unsigned=%d text=%d binary=%d",
			ItemValueTypeFloat,
			ItemValueTypeChar,
			ItemValueTypeLog,
			ItemValueTypeUnsigned,
			ItemValueTypeText,
			ItemValueTypeBinary,
		)
	}

	if ItemStatusEnabled != 0 || ItemStatusDisabled != 1 {
		t.Fatalf("ItemStatus constants mismatch: got enabled=%d disabled=%d", ItemStatusEnabled, ItemStatusDisabled)
	}

	if ItemStateNormal != 0 || ItemStateNotSupported != 1 {
		t.Fatalf("ItemState constants mismatch: got normal=%d notsupported=%d", ItemStateNormal, ItemStateNotSupported)
	}

	if ItemRequestMethodGET != 0 || ItemRequestMethodHEAD != 3 ||
		ItemPostTypeJSON != 2 || ItemPostTypeXML != 3 ||
		ItemRetrieveModeBoth != 2 || ItemHTTPAuthDigest != 4 {
		t.Fatalf(
			"Item HTTP constants mismatch: got get=%d head=%d postjson=%d postxml=%d retrieveboth=%d digest=%d",
			ItemRequestMethodGET,
			ItemRequestMethodHEAD,
			ItemPostTypeJSON,
			ItemPostTypeXML,
			ItemRetrieveModeBoth,
			ItemHTTPAuthDigest,
		)
	}

	if ItemPreprocessingCustomMultiplier != 1 || ItemPreprocessingJSONPath != 12 ||
		ItemPreprocessingDiscardUnchanged != 19 || ItemPreprocessingJavaScript != 21 ||
		ItemPreprocessingSNMPGetValue != 30 {
		t.Fatalf(
			"ItemPreprocessingType constants mismatch: got multiplier=%d jsonpath=%d discard=%d js=%d snmpget=%d",
			ItemPreprocessingCustomMultiplier,
			ItemPreprocessingJSONPath,
			ItemPreprocessingDiscardUnchanged,
			ItemPreprocessingJavaScript,
			ItemPreprocessingSNMPGetValue,
		)
	}

	if ItemPreprocessingErrorDefault != 0 || ItemPreprocessingErrorSetError != 3 {
		t.Fatalf(
			"ItemPreprocessingErrorHandler constants mismatch: got default=%d seterror=%d",
			ItemPreprocessingErrorDefault,
			ItemPreprocessingErrorSetError,
		)
	}
}
