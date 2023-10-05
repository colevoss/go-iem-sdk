package iem

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

// IEMWeatherData represents a weather reading at a station at a time
// Properties are present depending on the query used to fetch weather data
type IEMWeatherData struct {
	Station string `json:"station"` // Station recorded at (station)

	Time *time.Time `json:"time"` // Time recorded at (valid)

	Lon float64 `json:"lon,omitempty"` // Longitude recorded at (lon)
	Lat float64 `json:"lat,omitempty"` // Latitude recorded at (lat)

	Elevation string `json:"elevation,omitempty"` // Elevation recorded at (elev)

	TemperatureF float64 `json:"tmpf,omitempty"` // Air Temperature [F] (tempf)
	TemperatureC float64 `json:"tmpc,omitempty"` // Air Temperature [C] (tempc)

	DewPointF float64 `json:"dwpf,omitempty"` // Dew Point [F] (dwpf)
	DewPointC float64 `json:"dwpc,omitempty"` // Dew Point [C] (dwpc)

	RelativeHumidity float64 `json:"relh,omitempty"` // Relative humidity [%] (relh)

	Feel float64 `json:"feel,omitempty"` // Heat Index/Wind Chil [F] (feel)

	WindDirection  float64 `json:"drct,omitempty"` // Wind Direction [deg] (drct)
	WindSpeedKnots float64 `json:"sknt,omitempty"` // Wind Speed [knots] (sknt)
	WindSpeedMPH   float64 `json:"sped,omitempty"` // Wind Speed [mph] (sped)

	WindGustKnots float64 `json:"gust,omitempty"`     // Wind Gust [knots] (gust)
	WindGustMPH   float64 `json:"gust_mph,omitempty"` // Wind Gust [mph] (gust_mph)

	PeakWindGustKnots float64    `json:"peak_wind_gust,omitempty"` // Peak Wind Gust [knots] (peak_wind_gust)
	PeakWindGustMPH   float64    `json:"peak_wind_mph,omitempty"`  // Peak Wind Gust [mph] (peak_wind_mph)
	PeakWindDirection float64    `json:"peak_wind_drct,omitempty"` // Peak Wind Direction [deg] (peak_wind_drct)
	PeakWindTime      *time.Time `json:"peak_wind_time,omitempty"` // Peak Wind Time (peak_wind_time)

	Altimeter float64 `json:"alti,omitempty"` // Altimeter [inches] (alti)

	SeaLevelPressure float64 `json:"mslp,omitempty"` // Sea Level Pressure [mb] (mslp)

	PrecipMM   float64 `json:"p01m,omitempty"` // 1 hour Precipitation [mm] (p01m)
	PrecipInch float64 `json:"p01i,omitempty"` // 1 hour Precipitation [inch] (po1i)

	Visibility float64 `json:"vsby,omitempty"` // Visibility [miles] (vsby)

	CloudCoverageL1 string `json:"skyc1,omitempty"` // Cloud Coverage Level 1 (skyc1)
	CloudCoverageL2 string `json:"skyc2,omitempty"` // Cloud Coverage Level 2 (skyc2)
	CloudCoverageL3 string `json:"skyc3,omitempty"` // Cloud Coverage Level 3 (skyc3)

	CloudHeightL1 float64 `json:"skyl1,omitempty"` // Cloud Height Level 1 [ft] (skyl1)
	CloudHeightL2 float64 `json:"skyl2,omitempty"` // Cloud Height Level 2 [ft] (skyl2)
	CloudHeightL3 float64 `json:"skyl3,omitempty"` // Cloud Height Level 3 [ft] (skyl3)

	PresentWeatherCodes string `json:"wxcodes,omitempty"` // Present Weather Code(s)

	IceAccretion1HR float64 `json:"ice_accretion_1hr,omitempty"` // Ice Accretion 1 Hour (ice_accretion_1hr)
	IceAccretion3HR float64 `json:"ice_accretion_3hr,omitempty"` // Ice Accretion 3 Hour (ice_accretion_3hr)
	IceAccretion6HR float64 `json:"ice_accretion_6hr,omitempty"` // Ice Accretion 6 Hour (ice_accretion_6hr)

	SnowDepth float64 `json:"snowdepth,omitempty"` // Snow Depth (4-group) [inch] (snowdepth)

	METAR string `json:"metar,omitempty"` // Raw METAR (metar)
}

// Parse weather data from a io.Reader that reads CSV data based on a WeatherDataQueryBuilder
func ParseWeatherData(reader io.Reader, query *WeatherDataQueryBuilder) ([]*IEMWeatherData, error) {
	csvReader := csv.NewReader(reader)

	csvReader.ReuseRecord = true

	data := []*IEMWeatherData{}
	keyIndecies := weatherDataIndecies(make(map[string]int))

	r := 0

	for {
		csvRecord, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if r == 0 {
			r++

			for i, header := range csvRecord {
				keyIndecies[header] = i
			}

			continue
		}

		weatherData, err := keyIndecies.csvRecordToWeatherData(&csvRecord, query)

		if err != nil {
			return nil, err
		}

		data = append(data, weatherData)
	}

	return data, nil
}

type weatherDataIndecies map[string]int

func (w *weatherDataIndecies) getIndex(key string) (int, bool) {
	i, ok := (*w)[key]

	return i, ok
}

// TODO: Better error handling
func (w *weatherDataIndecies) csvRecordToWeatherData(record *[]string, query *WeatherDataQueryBuilder) (*IEMWeatherData, error) {
	data := &IEMWeatherData{}

	err := w.setString("station", record, &data.Station, query)
	err = w.setTime("valid", record, &data.Time, query)
	err = w.setFloat("lon", record, &data.Lon, query)
	err = w.setFloat("lat", record, &data.Lat, query)
	err = w.setString("elevation", record, &data.Elevation, query)
	err = w.setFloat("tmpf", record, &data.TemperatureF, query)
	err = w.setFloat("tmpc", record, &data.TemperatureC, query)
	err = w.setFloat("dwpf", record, &data.DewPointF, query)
	err = w.setFloat("dwpc", record, &data.DewPointC, query)
	err = w.setFloat("relh", record, &data.RelativeHumidity, query)
	err = w.setFloat("feel", record, &data.Feel, query)
	err = w.setFloat("drct", record, &data.WindDirection, query)
	err = w.setFloat("sknt", record, &data.WindSpeedKnots, query)
	err = w.setFloat("sped", record, &data.WindSpeedMPH, query)
	err = w.setFloat("gust", record, &data.WindGustKnots, query)
	err = w.setFloat("gust_mph", record, &data.WindGustMPH, query)
	err = w.setFloat("peak_wind_gust", record, &data.PeakWindGustKnots, query)
	err = w.setFloat("peak_wind_gust_mph", record, &data.PeakWindGustMPH, query)
	err = w.setFloat("peak_wind_drct", record, &data.PeakWindDirection, query)
	err = w.setTime("peak_wind_time", record, &data.PeakWindTime, query)
	err = w.setFloat("alti", record, &data.Altimeter, query)
	err = w.setFloat("mslp", record, &data.SeaLevelPressure, query)
	err = w.setFloat("p01m", record, &data.PrecipMM, query)
	err = w.setFloat("p01i", record, &data.PrecipInch, query)
	err = w.setFloat("vsby", record, &data.Visibility, query)
	err = w.setString("skyc1", record, &data.CloudCoverageL1, query)
	err = w.setString("skyc2", record, &data.CloudCoverageL2, query)
	err = w.setString("skyc3", record, &data.CloudCoverageL3, query)
	err = w.setFloat("skyl1", record, &data.CloudHeightL1, query)
	err = w.setFloat("skyl2", record, &data.CloudHeightL2, query)
	err = w.setFloat("skyl3", record, &data.CloudHeightL3, query)
	err = w.setString("wxcodes", record, &data.PresentWeatherCodes, query)
	err = w.setFloat("ice_accretion_1hr", record, &data.IceAccretion1HR, query)
	err = w.setFloat("ice_accretion_3hr", record, &data.IceAccretion3HR, query)
	err = w.setFloat("ice_accretion_6hr", record, &data.IceAccretion6HR, query)
	err = w.setFloat("snowdepth", record, &data.SnowDepth, query)
	err = w.setString("metar", record, &data.METAR, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (w *weatherDataIndecies) setString(key string, record *[]string, data *string, query *WeatherDataQueryBuilder) error {
	idx, ok := w.getIndex(key)
	if !ok {
		return nil
	}

	v := (*record)[idx]

	if query.isMissingOrTrace(v) {
		return nil
	}

	*data = (*record)[idx]
	return nil
}

func (w *weatherDataIndecies) setFloat(key string, record *[]string, data *float64, query *WeatherDataQueryBuilder) error {
	idx, ok := w.getIndex(key)
	if !ok {
		return nil
	}

	v := (*record)[idx]

	if query.isMissingOrTrace(v) {
		return nil
	}

	f, err := strconv.ParseFloat(v, 64)

	if err != nil {
		return fmt.Errorf("error parsing %s: [%w]", key, err)
	}

	*data = f

	return nil
}

func (w *weatherDataIndecies) setTime(key string, record *[]string, data **time.Time, query *WeatherDataQueryBuilder) error {
	idx, ok := w.getIndex(key)
	if !ok {
		return nil
	}

	v := (*record)[idx]

	if query.isMissingOrTrace(v) {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04", v)

	if err != nil {
		return fmt.Errorf("error parsing %s [%w]", key, err)
	}

	*data = &t

	return nil
}
