package medtronic

import (
	"fmt"
	"time"
)

type HistoryRecordType byte

//go:generate stringer -type HistoryRecordType

const (
	Bolus                   HistoryRecordType = 0x01
	Prime                   HistoryRecordType = 0x03
	Alarm                   HistoryRecordType = 0x06
	DailyTotal              HistoryRecordType = 0x07
	BasalProfileBefore      HistoryRecordType = 0x08
	BasalProfileAfter       HistoryRecordType = 0x09
	BgCapture               HistoryRecordType = 0x0A
	SensorAlarm             HistoryRecordType = 0x0B
	ClearAlarm              HistoryRecordType = 0x0C
	ChangeBasalPattern      HistoryRecordType = 0x14
	TempBasalDuration       HistoryRecordType = 0x16
	ChangeTime              HistoryRecordType = 0x17
	NewTime                 HistoryRecordType = 0x18
	LowBattery              HistoryRecordType = 0x19
	BatteryChange           HistoryRecordType = 0x1A
	SetAutoOff              HistoryRecordType = 0x1B
	SuspendPump             HistoryRecordType = 0x1E
	ResumePump              HistoryRecordType = 0x1F
	SelfTest                HistoryRecordType = 0x20
	Rewind                  HistoryRecordType = 0x21
	EnableChildBlock        HistoryRecordType = 0x23
	MaxBolus                HistoryRecordType = 0x24
	EnableRemote            HistoryRecordType = 0x26
	MaxBasal                HistoryRecordType = 0x2C
	EnableBolusWizard       HistoryRecordType = 0x2D
	SetAlarmClockTime       HistoryRecordType = 0x32
	TempBasalRate           HistoryRecordType = 0x33
	LowReservoir            HistoryRecordType = 0x34
	SensorStatus            HistoryRecordType = 0x3B
	EnableMeter             HistoryRecordType = 0x3C
	ChangeBolusWizardSetup  HistoryRecordType = 0x4F
	BolusWizardSetup        HistoryRecordType = 0x5A
	BolusWizard             HistoryRecordType = 0x5B
	UnabsorbedInsulin       HistoryRecordType = 0x5C
	EnableVariableBolus     HistoryRecordType = 0x5E
	ChangeEasyBolus         HistoryRecordType = 0x5F
	EnableBgReminder        HistoryRecordType = 0x60
	EnableAlarmClock        HistoryRecordType = 0x61
	ChangeTempBasalType     HistoryRecordType = 0x62
	ChangeAlarmType         HistoryRecordType = 0x63
	ChangeTimeFormat        HistoryRecordType = 0x64
	ChangeReservoirWarning  HistoryRecordType = 0x65
	EnableBolusReminder     HistoryRecordType = 0x66
	SetBolusReminderTime    HistoryRecordType = 0x67
	DeleteBolusReminderTime HistoryRecordType = 0x68
	DeleteAlarmClockTime    HistoryRecordType = 0x6A
	DailyTotal522           HistoryRecordType = 0x6D
	DailyTotal523           HistoryRecordType = 0x6E
	BasalProfileStart       HistoryRecordType = 0x7B
	ConnectOtherDevices     HistoryRecordType = 0x7C
	ChangeOtherDevice       HistoryRecordType = 0x81
	DeleteOtherDevice       HistoryRecordType = 0x82
)

type (
	HistoryRecord struct {
		Data              []byte                   `json:",omitempty"`
		Time              time.Time                `json:",omitempty"`
		Value             *int                     `json:",omitempty"`
		Enabled           *bool                    `json:",omitempty"`
		Insulin           *MilliUnits              `json:",omitempty"`
		Duration          *time.Duration           `json:",omitempty"`
		TempBasalType     *TempBasalType           `json:",omitempty"`
		BasalProfile      BasalRateSchedule        `json:",omitempty"`
		BasalProfileStart *BasalProfileStartRecord `json:",omitempty"`
		Prime             *PrimeRecord             `json:",omitempty"`
		Bolus             *BolusRecord             `json:",omitempty"`
		BolusWizard       *BolusWizardRecord       `json:",omitempty"`
		BolusWizardSetup  *BolusWizardSetupRecord  `json:",omitempty"`
		UnabsorbedInsulin []UnabsorbedBolus        `json:",omitempty"`
	}

	BasalProfileStartRecord struct {
		ProfileIndex int
		BasalRate    BasalRate
	}

	PrimeRecord struct {
		Fixed  MilliUnits
		Manual MilliUnits
	}

	BolusRecord struct {
		Programmed MilliUnits
		Amount     MilliUnits
		Unabsorbed MilliUnits
		Duration   time.Duration // non-zero for square wave bolus
	}

	BolusWizardRecord struct {
		GlucoseInput int // mg/dL or μmol/L
		TargetLow    int // mg/dL or μmol/L
		TargetHigh   int // mg/dL or μmol/L
		Sensitivity  int // mg/dL or μmol/L reduction per insulin unit
		CarbInput    int // grams or exchanges
		CarbRatio    int // grams or exchanges covered by TEN insulin units
		Unabsorbed   MilliUnits
		Correction   MilliUnits
		Food         MilliUnits
		Bolus        MilliUnits
	}

	BolusWizardConfig struct {
		Ratios        CarbRatioSchedule
		Sensitivities InsulinSensitivitySchedule
		Targets       GlucoseTargetSchedule
	}

	BolusWizardSetupRecord struct {
		Before BolusWizardConfig
		After  BolusWizardConfig
	}

	UnabsorbedBolus struct {
		Bolus MilliUnits
		Age   time.Duration
	}

	UnknownRecordTypeError struct {
		Data []byte
	}
)

func (r HistoryRecord) Type() HistoryRecordType {
	return HistoryRecordType(r.Data[0])
}

func decodeBase(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Time: decodeTimestamp(data[2:7]),
		Data: data[:7],
	}
}

func decodeValue(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	v := int(data[1])
	r.Value = &v
	return r
}

func decodeBolus(data []byte, newerPump bool) HistoryRecord {
	if newerPump {
		return HistoryRecord{
			Bolus: &BolusRecord{
				Programmed: twoByteMilliUnits(data[1:3], true),
				Amount:     twoByteMilliUnits(data[3:5], true),
				Unabsorbed: twoByteMilliUnits(data[5:7], true),
				Duration:   scheduleToDuration(data[7]),
			},
			Time: decodeTimestamp(data[8:13]),
			Data: data[:13],
		}
	} else {
		return HistoryRecord{
			Bolus: &BolusRecord{
				Programmed: byteToMilliUnits(data[1], false),
				Amount:     byteToMilliUnits(data[2], false),
				Duration:   scheduleToDuration(data[3]),
			},
			Time: decodeTimestamp(data[4:9]),
			Data: data[:9],
		}
	}
}

func decodePrime(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Prime: &PrimeRecord{
			Fixed:  byteToMilliUnits(data[2], false),
			Manual: byteToMilliUnits(data[4], false),
		},
		Time: decodeTimestamp(data[5:10]),
		Data: data[:10],
	}
}

func decodeAlarm(data []byte, newerPump bool) HistoryRecord {
	r := HistoryRecord{
		Time: decodeTimestamp(data[4:9]),
		Data: data[:9],
	}
	alarm := int(data[1])
	r.Value = &alarm
	return r
}

func decodeSensorAlarm(data []byte, newerPump bool) HistoryRecord {
	r := HistoryRecord{
		Time: decodeTimestamp(data[3:8]),
		Data: data[:8],
	}
	alarm := int(data[1])
	r.Value = &alarm
	return r
}

func decodeDailyTotal(data []byte, newerPump bool) HistoryRecord {
	t := decodeDate(data[5:7])
	total := twoByteMilliUnits(data[3:5], true)
	if newerPump {
		return HistoryRecord{
			Time:    t,
			Insulin: &total,
			Data:    data[:10],
		}
	} else {
		return HistoryRecord{
			Time:    t,
			Insulin: &total,
			Data:    data[:7],
		}
	}
}

// Note that this is a different format than the response to BasalRates.
func decodeBasalRate(data []byte) BasalRate {
	return BasalRate{
		Start: scheduleToDuration(data[0]),
		Rate:  byteToMilliUnits(data[1], true),
		// data[2] unused
	}
}

func decodeBasalProfile(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	body := data[7:]
	sched := []BasalRate{}
	for i := 0; i < 144; i += 3 {
		b := decodeBasalRate(body[i : i+3])
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 0 && b.Start == 0 && b.Rate == 0 {
			break
		}
		sched = append(sched, b)
	}
	r.BasalProfile = sched
	r.Data = data[:152]
	return r
}

func decodeBgCapture(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	bg := int(data[1]) | int(data[6]>>7)<<8
	r.Value = &bg
	return r
}

func decodeClearAlarm(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	alarm := int(data[1])
	r.Value = &alarm
	return r
}

func decodeTempBasalDuration(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	d := scheduleToDuration(data[1])
	r.Duration = &d
	return r
}

func decodeSetAutoOff(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	d := time.Duration(data[1]) * time.Hour
	r.Duration = &d
	return r
}

func decodeEnable(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	enabled := false
	if data[1] != 0 {
		enabled = true
	}
	r.Enabled = &enabled
	return r
}

func decodeMax(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	max := byteToMilliUnits(data[1], true)
	r.Insulin = &max
	return r
}

func decodeEnableRemote(data []byte, newerPump bool) HistoryRecord {
	r := decodeEnable(data, newerPump)
	r.Data = data[:21]
	return r
}

func decodeTempBasalRate(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	rate := int(data[1])
	tempBasalType := TempBasalType(data[7] >> 3)
	if tempBasalType == Absolute {
		rate *= 25
	}
	r.Value = &rate
	r.TempBasalType = &tempBasalType
	r.Data = data[:8]
	return r
}

func decodeLowReservoir(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	amount := int(data[1]) * 100
	r.Value = &amount
	return r
}

func decodeEnableMeter(data []byte, newerPump bool) HistoryRecord {
	r := decodeEnable(data, newerPump)
	r.Data = data[:21]
	return r
}

func decodeChangeBolusWizardSetup(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.Data = data[:39]
	return r
}

func decodeBolusWizardConfig(data []byte, newerPump bool) BolusWizardConfig {
	const numEntries = 8
	conf := BolusWizardConfig{}
	carbUnits := Grams // FIXME
	step := carbRatioStep(newerPump)
	conf.Ratios = decodeCarbRatioSchedule(data[2:2+numEntries*step], carbUnits, newerPump)
	data = data[2+numEntries*step:]
	bgUnits := MgPerDeciLiter // FIXME
	conf.Sensitivities = decodeInsulinSensitivitySchedule(data[:numEntries*2], bgUnits)
	if newerPump {
		data = data[numEntries*2+2:]
	} else {
		data = data[numEntries*2:]
	}
	conf.Targets = decodeGlucoseTargetSchedule(data[:numEntries*3], bgUnits)
	return conf
}

func decodeBolusWizardSetup(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	if newerPump {
		r.Data = data[:144]
	} else {
		r.Data = data[:124]
	}
	n := len(r.Data) - 1
	body := data[7:n]
	half := (n - 7) / 2
	r.BolusWizardSetup = &BolusWizardSetupRecord{
		Before: decodeBolusWizardConfig(body[:half], newerPump),
		After:  decodeBolusWizardConfig(body[half:], newerPump),
	}
	return r
}

func decodeBolusWizard(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	bg := int(data[1])
	body := data[7:]
	if newerPump {
		r.BolusWizard = &BolusWizardRecord{
			GlucoseInput: bg | int(body[1]&0x3)<<8,
			CarbInput:    int(body[1]&0xC)<<6 + int(body[0]),
			CarbRatio:    (int(body[2]&0x7)<<8 | int(body[3])),
			Sensitivity:  int(body[4]),
			TargetLow:    int(body[5]),
			Correction:   intToMilliUnits(int(body[9]&0x38)<<5+int(body[6]), true),
			Food:         twoByteMilliUnits(body[7:9], true),
			Unabsorbed:   twoByteMilliUnits(body[10:12], true),
			Bolus:        twoByteMilliUnits(body[12:14], true),
			TargetHigh:   int(body[14]),
		}
		r.Data = data[:22]
	} else {
		r.BolusWizard = &BolusWizardRecord{
			GlucoseInput: bg | int(body[1]&0xF)<<8,
			CarbInput:    int(body[0]),
			CarbRatio:    int(body[2]),
			Sensitivity:  int(body[3]),
			TargetLow:    int(body[4]),
			Correction:   intToMilliUnits(int(body[7])+int(body[5]&0xF), false),
			Food:         byteToMilliUnits(body[6], false),
			Unabsorbed:   byteToMilliUnits(body[9], false),
			Bolus:        byteToMilliUnits(body[11], false),
			TargetHigh:   int(body[12]),
		}
		r.Data = data[:20]
	}
	return r
}

func decodeUnabsorbedInsulin(data []byte, newerPump bool) HistoryRecord {
	n := int(data[1]) - 2
	body := data[2:]
	unabsorbed := []UnabsorbedBolus{}
	for i := 0; i < n; i += 3 {
		amount := byteToMilliUnits(body[i], true)
		curve := body[i+2]
		age := time.Duration(body[i+1]+(curve&0x30)<<4) * time.Minute
		unabsorbed = append(unabsorbed, UnabsorbedBolus{
			Bolus: amount,
			Age:   age,
		})
	}
	return HistoryRecord{
		Data:              data[:n+2],
		UnabsorbedInsulin: unabsorbed,
	}
}

func decodeChangeTempBasalType(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	tempBasalType := TempBasalType(data[1])
	r.TempBasalType = &tempBasalType
	return r
}

func decodeChangeReservoirWarning(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	v := data[1]
	if v&0x1 == 0 {
		amount := MilliUnits(1000 * int(v>>2))
		r.Insulin = &amount
	} else {
		d := scheduleToDuration(v >> 2)
		r.Duration = &d
	}
	return r
}

func decodeBolusReminder(data []byte, newerPump bool) HistoryRecord {
	r := decodeEnable(data, newerPump)
	r.Data = data[:9]
	return r
}

func decodeDailyTotal522(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Time: decodeDate(data[1:3]),
		Data: data[:44],
	}
}

func decodeDailyTotal523(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Time: decodeDate(data[1:3]),
		Data: data[:52],
	}
}

func decodeBasalProfileStart(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.BasalProfileStart = &BasalProfileStartRecord{
		ProfileIndex: int(data[1]),
		BasalRate:    decodeBasalRate(data[7:10]),
	}
	r.Data = data[:10]
	return r
}

func decodeOtherDevice(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.Data = data[:12]
	return r
}

var decode = map[HistoryRecordType]func([]byte, bool) HistoryRecord{
	Bolus:                   decodeBolus,
	Prime:                   decodePrime,
	Alarm:                   decodeAlarm,
	DailyTotal:              decodeDailyTotal,
	BasalProfileBefore:      decodeBasalProfile,
	BasalProfileAfter:       decodeBasalProfile,
	BgCapture:               decodeBgCapture,
	SensorAlarm:             decodeSensorAlarm,
	ClearAlarm:              decodeValue,
	ChangeBasalPattern:      decodeValue,
	TempBasalDuration:       decodeTempBasalDuration,
	ChangeTime:              decodeBase,
	NewTime:                 decodeBase,
	LowBattery:              decodeBase,
	BatteryChange:           decodeBase,
	SetAutoOff:              decodeSetAutoOff,
	SuspendPump:             decodeBase,
	ResumePump:              decodeBase,
	SelfTest:                decodeBase,
	Rewind:                  decodeBase,
	EnableChildBlock:        decodeEnable,
	MaxBolus:                decodeMax,
	EnableRemote:            decodeEnableRemote,
	MaxBasal:                decodeMax,
	EnableBolusWizard:       decodeEnable,
	SetAlarmClockTime:       decodeBase,
	TempBasalRate:           decodeTempBasalRate,
	LowReservoir:            decodeLowReservoir,
	SensorStatus:            decodeEnable,
	EnableMeter:             decodeEnableMeter,
	ChangeBolusWizardSetup:  decodeChangeBolusWizardSetup,
	BolusWizardSetup:        decodeBolusWizardSetup,
	BolusWizard:             decodeBolusWizard,
	UnabsorbedInsulin:       decodeUnabsorbedInsulin,
	EnableVariableBolus:     decodeEnable,
	ChangeEasyBolus:         decodeBase,
	EnableBgReminder:        decodeEnable,
	EnableAlarmClock:        decodeEnable,
	ChangeTempBasalType:     decodeChangeTempBasalType,
	ChangeAlarmType:         decodeValue,
	ChangeTimeFormat:        decodeValue,
	ChangeReservoirWarning:  decodeChangeReservoirWarning,
	EnableBolusReminder:     decodeEnable,
	SetBolusReminderTime:    decodeBolusReminder,
	DeleteBolusReminderTime: decodeBolusReminder,
	DeleteAlarmClockTime:    decodeBase,
	DailyTotal522:           decodeDailyTotal522,
	DailyTotal523:           decodeDailyTotal523,
	BasalProfileStart:       decodeBasalProfileStart,
	ConnectOtherDevices:     decodeEnable,
	ChangeOtherDevice:       decodeOtherDevice,
	DeleteOtherDevice:       decodeOtherDevice,
}

func DecodeHistoryRecord(data []byte, newerPump bool) (HistoryRecord, error) {
	decoder := decode[HistoryRecordType(data[0])]
	if decoder != nil {
		return decoder(data, newerPump), nil
	} else {
		return HistoryRecord{}, unknownRecord(data)
	}
}

func (e UnknownRecordTypeError) Error() string {
	return fmt.Sprintf("unknown record type here: % X", e.Data)
}

func unknownRecord(data []byte) error {
	return UnknownRecordTypeError{
		Data: data,
	}
}

// Decode records in a page of data and return them in reverse order,
// to match the order of the history pages themselves.
func DecodeHistoryRecords(data []byte, newerPump bool) ([]HistoryRecord, error) {
	results := []HistoryRecord{}
	r := HistoryRecord{}
	err := error(nil)
	for !allZero(data) {
		r, err = DecodeHistoryRecord(data, newerPump)
		if err != nil {
			break
		}
		results = append(results, r)
		data = data[len(r.Data):]
	}
	reverseHistoryRecords(results)
	return results, err
}

// Partially filled history pages are padded to the end with zero bytes.
func allZero(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

func reverseHistoryRecords(a []HistoryRecord) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
