// Code generated by "stringer -type HistoryRecordType"; DO NOT EDIT

package medtronic

import "fmt"

const _HistoryRecordType_name = "BolusPrimeAlarmDailyTotalBasalProfileBeforeBasalProfileAfterBGCaptureSensorAlarmClearAlarmChangeBasalPatternTempBasalDurationChangeTimeNewTimeLowBatteryBatteryChangeSetAutoOffSuspendPumpResumePumpSelfTestRewindEnableChildBlockMaxBolusEnableRemoteMaxBasalEnableBolusWizardChangeBGReminderSetAlarmClockTimeTempBasalRateLowReservoirAlarmClockSensorStatusEnableMeterMealMarkerExerciseMarkerInsulinMarkerOtherMarkerChangeBolusWizardSetupChangeGlucoseUnitsBolusWizardSetupBolusWizardUnabsorbedInsulinEnableVariableBolusChangeEasyBolusEnableBGReminderEnableAlarmClockChangeTempBasalTypeChangeAlarmTypeChangeTimeFormatChangeReservoirWarningEnableBolusReminderSetBolusReminderTimeDeleteBolusReminderTimeDeleteAlarmClockTimeDailyTotal522DailyTotal523ChangeCarbUnitsBasalProfileStartConnectOtherDevicesChangeOtherDeviceChangeMarriageDeleteOtherDeviceEnableCaptureEvent"

var _HistoryRecordType_map = map[HistoryRecordType]string{
	1:   _HistoryRecordType_name[0:5],
	3:   _HistoryRecordType_name[5:10],
	6:   _HistoryRecordType_name[10:15],
	7:   _HistoryRecordType_name[15:25],
	8:   _HistoryRecordType_name[25:43],
	9:   _HistoryRecordType_name[43:60],
	10:  _HistoryRecordType_name[60:69],
	11:  _HistoryRecordType_name[69:80],
	12:  _HistoryRecordType_name[80:90],
	20:  _HistoryRecordType_name[90:108],
	22:  _HistoryRecordType_name[108:125],
	23:  _HistoryRecordType_name[125:135],
	24:  _HistoryRecordType_name[135:142],
	25:  _HistoryRecordType_name[142:152],
	26:  _HistoryRecordType_name[152:165],
	27:  _HistoryRecordType_name[165:175],
	30:  _HistoryRecordType_name[175:186],
	31:  _HistoryRecordType_name[186:196],
	32:  _HistoryRecordType_name[196:204],
	33:  _HistoryRecordType_name[204:210],
	35:  _HistoryRecordType_name[210:226],
	36:  _HistoryRecordType_name[226:234],
	38:  _HistoryRecordType_name[234:246],
	44:  _HistoryRecordType_name[246:254],
	45:  _HistoryRecordType_name[254:271],
	49:  _HistoryRecordType_name[271:287],
	50:  _HistoryRecordType_name[287:304],
	51:  _HistoryRecordType_name[304:317],
	52:  _HistoryRecordType_name[317:329],
	53:  _HistoryRecordType_name[329:339],
	59:  _HistoryRecordType_name[339:351],
	60:  _HistoryRecordType_name[351:362],
	64:  _HistoryRecordType_name[362:372],
	65:  _HistoryRecordType_name[372:386],
	66:  _HistoryRecordType_name[386:399],
	67:  _HistoryRecordType_name[399:410],
	79:  _HistoryRecordType_name[410:432],
	86:  _HistoryRecordType_name[432:450],
	90:  _HistoryRecordType_name[450:466],
	91:  _HistoryRecordType_name[466:477],
	92:  _HistoryRecordType_name[477:494],
	94:  _HistoryRecordType_name[494:513],
	95:  _HistoryRecordType_name[513:528],
	96:  _HistoryRecordType_name[528:544],
	97:  _HistoryRecordType_name[544:560],
	98:  _HistoryRecordType_name[560:579],
	99:  _HistoryRecordType_name[579:594],
	100: _HistoryRecordType_name[594:610],
	101: _HistoryRecordType_name[610:632],
	102: _HistoryRecordType_name[632:651],
	103: _HistoryRecordType_name[651:671],
	104: _HistoryRecordType_name[671:694],
	106: _HistoryRecordType_name[694:714],
	109: _HistoryRecordType_name[714:727],
	110: _HistoryRecordType_name[727:740],
	111: _HistoryRecordType_name[740:755],
	123: _HistoryRecordType_name[755:772],
	124: _HistoryRecordType_name[772:791],
	125: _HistoryRecordType_name[791:808],
	129: _HistoryRecordType_name[808:822],
	130: _HistoryRecordType_name[822:839],
	131: _HistoryRecordType_name[839:857],
}

func (i HistoryRecordType) String() string {
	if str, ok := _HistoryRecordType_map[i]; ok {
		return str
	}
	return fmt.Sprintf("HistoryRecordType(%d)", i)
}
