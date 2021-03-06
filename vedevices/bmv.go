package vedevices

var RegisterListBmv700Essential = Registers{
	"MainVoltage": Register{
		Address:       0xED8D,
		Factor:        0.01,
		Unit:          "V",
		Signed:        false,
		RoundDecimals: 2,
	},
	"Current": Register{
		Address:       0xED8F,
		Factor:        0.1,
		Unit:          "A",
		Signed:        true,
		RoundDecimals: 1,
	},
	"Power": Register{
		Address:       0xED8E,
		Factor:        1,
		Unit:          "W",
		Signed:        true,
		RoundDecimals: 0,
	},
}

var RegisterListBmv700 = mergeRegisters(
	RegisterListBmv700Essential,
	Registers{
		"Consumed": Register{
			Address:       0xEEFF,
			Factor:        0.1,
			Unit:          "Ah",
			Signed:        true,
			RoundDecimals: 1,
		},
		"StateOfCharge": Register{
			Address:       0x0FFF,
			Factor:        0.01,
			Unit:          "%",
			Signed:        false,
			RoundDecimals: 0,
		},
		"TimeToGo": Register{
			Address:       0x0FFE,
			Factor:        1,
			Unit:          "min",
			Signed:        false,
			RoundDecimals: 0,
		},
		"Temperature": Register{
			Address:       0xEDEC,
			Factor:        0.01,
			Unit:          "K",
			Signed:        false,
			RoundDecimals: 1,
		},
		"DepthOfTheDeepestDischarge": Register{
			Address:       0x0300,
			Factor:        0.1,
			Unit:          "Ah",
			Signed:        true,
			RoundDecimals: 0,
		},
		"DepthOfTheLastDischarge": Register{
			Address:       0x0301,
			Factor:        0.1,
			Unit:          "Ah",
			Signed:        true,
			RoundDecimals: 0,
		},
		"DepthOfTheAverageDischarge": Register{
			Address:       0x0302,
			Factor:        0.1,
			Unit:          "Ah",
			Signed:        true,
			RoundDecimals: 0,
		},
		"NumberOfCycles": Register{
			Address:       0x0303,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		"NumberOfFullDischarges": Register{
			Address:       0x0304,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		"CumulativeAmpHours": Register{
			Address:       0x0305,
			Factor:        0.1,
			Unit:          "Ah",
			Signed:        true,
			RoundDecimals: 0,
		},
		"MainVoltageMinimum": Register{
			Address:       0x0306,
			Factor:        0.01,
			Unit:          "V",
			Signed:        false,
			RoundDecimals: 2,
		},
		"MainVoltageMaximum": Register{
			Address:       0x0307,
			Factor:        0.01,
			Unit:          "V",
			Signed:        false,
			RoundDecimals: 2,
		},
		"HoursSinceFullCharge": Register{
			Address:       0x0308,
			Factor:        float64(24) / float64(86400),
			Unit:          "h",
			Signed:        false,
			RoundDecimals: 1,
		},
		"NumberOfAutomaticSynchronizations": Register{
			Address:       0x0309,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		"NumberOfLowMainVoltageAlarms": Register{
			Address:       0x030A,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		"NumberOfHighMainVoltageAlarms": Register{
			Address:       0x030B,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		"AmountOfDischargedEnergy": Register{
			Address:       0x0310,
			Factor:        0.01,
			Unit:          "kWh",
			Signed:        false,
			RoundDecimals: 1,
		},
		"AmountOfChargedEnergy": Register{
			Address:       0x0311,
			Factor:        0.01,
			Unit:          "kWh",
			Signed:        false,
			RoundDecimals: 1,
		},
	},
)

var RegisterListBmv702 = mergeRegisters(
	RegisterListBmv700,
	Registers{
		"AuxVoltage": Register{
			Address:       0xED7D,
			Factor:        0.01,
			Unit:          "V",
			Signed:        false,
			RoundDecimals: 2,
		},
		/*
		"Synchronized": Register{
			Address:       0xEEB6,
			Factor:        1,
			Unit:          "1",
			Signed:        false,
			RoundDecimals: 0,
		},
		*/
		"MidPointVoltage": Register{
			Address:       0x0382,
			Factor:        0.01,
			Unit:          "V",
			Signed:        false,
			RoundDecimals: 2,
		},
		"MidPointVoltageDeviation": Register{
			Address:       0x0383,
			Factor:        0.1,
			Unit:          "%",
			Signed:        true,
			RoundDecimals: 1,
		},
		/*
		"NumberOfLowAuxVoltageAlarms": Register{
			Address:       0x030C,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		"NumberOfHighAuxVoltageAlarms": Register{
			Address:       0x030D,
			Factor:        1,
			Unit:          "",
			Signed:        false,
			RoundDecimals: 0,
		},
		*/
		"AuxVoltageMinimum": Register{
			Address:       0x030E,
			Factor:        0.01,
			Unit:          "V",
			Signed:        true,
			RoundDecimals: 2,
		},
		"AuxVoltageMaximum": Register{
			Address:       0x030F,
			Factor:        0.01,
			Unit:          "V",
			Signed:        true,
			RoundDecimals: 2,
		},
	},
)
