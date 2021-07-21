package models

import "time"

type Meter_grid_stat struct {
	Current_Unbalance  float64   `json:"Current unbalance"`
	Frequency          float64   `json:"Frequency"`
	I1_Current         float64   `json:"I1 Current"`
	I1_Current_THD     float64   `json:"I1 Current THD"`
	I2_Current         float64   `json:"I2 Current"`
	I2_Current_THD     float64   `json:"I2 Current THD"`
	I3_Current         float64   `json:"I3 Current"`
	I3_Current_THD     float64   `json:"I3 Current THD"`
	In_Neutral_Current float64   `json:"In (neutral) Current"`
	Power_Factor_L1    float64   `json:"Power factor L1"`
	Power_Factor_L2    float64   `json:"Power factor L2"`
	Power_Factor_L3    float64   `json:"Power factor L3"`
	Timestamp          time.Time `json:"Timestamp" gorm:"type:time"`
	V1_Voltage         float64   `json:"V1 Voltage"`
	V1_Voltage_THD     float64   `json:"V1 Voltage THD"`
	V2_Voltage         float64   `json:"V2 Voltage"`
	V2_Voltage_THD     float64   `json:"V2 Voltage THD"`
	V3_Voltage         float64   `json:"V3 Voltage"`
	V3_Voltage_THD     float64   `json:"V3 Voltage THD"`
	Voltage_Unbalance  float64   `json:"Voltage unbalance"`
	KVA_L1             float64   `json:"kVA L1"`
	KVA_L2             float64   `json:"kVA L2"`
	KVA_L3             float64   `json:"kVA L3"`
	KW_L1              float64   `json:"kW L1"`
	KW_L2              float64   `json:"kW L2"`
	KW_L3              float64   `json:"kW L3"`
	Kvar_L1            float64   `json:"kvar L1"`
	Kvar_L2            float64   `json:"kvar L2"`
	Kvar_L3            float64   `json:"kvar L3"`
}
