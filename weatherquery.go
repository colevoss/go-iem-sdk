package iem

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type WeatherDataQueryFormat string
type WeatherDataQueryYesNo string
type WeatherDataQueryMissing string
type WeatherDataQueryTrace string

type WeatherDataData string

const (
	All                 WeatherDataData = "all"                // All Available
	TempF               WeatherDataData = "tmpf"               // Air Temperature [F]
	TempC               WeatherDataData = "tmpc"               // Air Temperature [C]
	DewPointF           WeatherDataData = "dwpf"               // Dew Point [F]
	DewPointC           WeatherDataData = "dwpc"               // Dew Point [C]
	RelativeHumidity    WeatherDataData = "relh"               // Relative Humidity [%]
	Feel                WeatherDataData = "feel"               // Heat Index/Wind Chill [F]
	WindDirection       WeatherDataData = "drct"               // Wind Direction
	WindSpeedKnots      WeatherDataData = "sknt"               // Wind Speed [knots]
	WindSpeedMPH        WeatherDataData = "sped"               // Wind Speed [mph]
	Altimeter           WeatherDataData = "alti"               // Altimeter [inches]
	SeaLevelPressure    WeatherDataData = "mslp"               // Sea Level Pressure [mb]
	PrecipMM            WeatherDataData = "p01m"               // 1 hour Precipitation [mm]
	PrecipInch          WeatherDataData = "p01i"               // 1 hour Precipitation [inch]
	Visibility          WeatherDataData = "vsby"               // Visibility [miles]
	WindGustKnots       WeatherDataData = "gust"               // Wind Gust [knots]
	WindGustMPH         WeatherDataData = "gust_mph"           // Wind Gust [mph]
	CloudCoverageL1     WeatherDataData = "skyc1"              // Cloud Coverage Level 1
	CloudCoverageL2     WeatherDataData = "skyc2"              // Cloud Coverage Level 2
	CloudCoverageL3     WeatherDataData = "skyc3"              // Cloud Coverage Level 3
	CloudHeightL1       WeatherDataData = "skyl1"              // Cloud Height Level 1 [ft]
	CloudHeightL2       WeatherDataData = "skyl2"              // Cloud Height Level 2 [ft]
	CloudHeightL3       WeatherDataData = "skyl3"              // Cloud Height Level 3 [ft]
	PresentWeatherCodes WeatherDataData = "wxcodes"            // Present Weather Code(s)
	IceAccretion1HR     WeatherDataData = "ice_accretion_1hr"  // Ice Accretion 1 Hour
	IceAccretion3HR     WeatherDataData = "ice_accretion_3hr"  // Ice Accretion 3 Hour
	IceAccretion6HR     WeatherDataData = "ice_accretion_6hr"  // Ice Accretion 6 Hour
	PeakWindGustKnots   WeatherDataData = "peak_wind_gust"     // Peak Wind Gust [knots]
	PeakWindGustMPH     WeatherDataData = "peak_wind_gust_mph" // Peak Wind Gust [MPH]
	PeakWindDirection   WeatherDataData = "peak_wind_drct"     // Peak Wind Direction [deg]
	PeakWindTime        WeatherDataData = "peak_wind_time"     // Peak Wind Time
	SnowDepth           WeatherDataData = "snowdepth"          // Snow Depth (4-group) [inch]
	METAR               WeatherDataData = "metar"              // Raw METAR
)

const (
	OnlyComma WeatherDataQueryFormat = "onlycomma"
	OnlyTDF   WeatherDataQueryFormat = "onlytdf"
	Comma     WeatherDataQueryFormat = "comma"
	TDF       WeatherDataQueryFormat = "tdf"

	Yes WeatherDataQueryYesNo = "yes"
	No  WeatherDataQueryYesNo = "no"

	MissingM     WeatherDataQueryMissing = "M"
	MissingNull  WeatherDataQueryMissing = "null"
	MissingEmpty WeatherDataQueryMissing = "empty"

	TraceT     WeatherDataQueryTrace = "T"
	TraceNull  WeatherDataQueryTrace = "null"
	TraceEmpty WeatherDataQueryTrace = "empty"
	TraceFloat WeatherDataQueryTrace = "0.0001"
)

const defaultTimeZone = "Etc/UTC"
const defaultLatLon = false

// const defaultLatLon = No
const defaultFormat = OnlyComma

// const defaultElevation = No
const defaultElevation = false
const defaultMissing = MissingM
const defaultTrace = TraceT

// const defaultDirect = No
const defaultDirect = false

var defaultReportType = []int{3, 4}

type WeatherDataQueryBuilderError struct {
	msg string
}

func (err WeatherDataQueryBuilderError) Error() string {
	return err.msg
}

func requiredError(prop string) WeatherDataQueryBuilderError {
	return WeatherDataQueryBuilderError{
		msg: fmt.Sprintf("WeatherDataQueryBuilder: %s property is required", prop),
	}
}

// WeatherDataQueryBuilder represents the url query sent to the IEM API
// when requesting weather data for stations
type WeatherDataQueryBuilder struct {
	// Stations to get weather data from
	stations []string

	// List of data points requested from IEM API
	data []WeatherDataData

	// Start date to query for
	start time.Time

	// End date to query for
	end time.Time

	// Etc/UTC
	// America/New_York
	// America/Chicago
	// America/Denver
	// America/Los_Angeles
	// America/Anchorage
	// Timezone used for dates
	tz string

	// onlycomma
	// onlytdf
	// comma
	// tdf
	// Format of returned weather data
	format WeatherDataQueryFormat

	// no
	// yes
	// latlon WeatherDataQueryYesNo
	// Include lat and lon for weather data in response
	latlon bool

	// no
	// yes
	// elev WeatherDataQueryYesNo
	// Include elev for weather data in response
	elev bool

	// M
	// null (as string)
	// empty
	// How missing data is represented
	missing WeatherDataQueryMissing

	// T
	// null (as string)
	// empty
	// 0.0001 (as string)
	// How trace data is represnted
	trace WeatherDataQueryTrace

	// direct WeatherDataQueryYesNo // should just be "no"
	direct bool // should just be "no"

	reportType []int // probably just 3 and 4
}

// Creates a new WeatherDataQueryBuilder with defaults set to optional fields
func NewWeatherDataQuery() *WeatherDataQueryBuilder {
	return &WeatherDataQueryBuilder{
		tz:         defaultTimeZone,
		latlon:     defaultLatLon,
		elev:       defaultElevation,
		missing:    defaultMissing,
		trace:      defaultTrace,
		direct:     defaultDirect,
		reportType: defaultReportType,
		format:     defaultFormat,

		start: time.Now(),
		end:   time.Now(),
	}
}

func (b *WeatherDataQueryBuilder) isMissingOrTrace(value string) bool {
	if value == string(b.missing) {
		return true
	}

	if value == string(b.trace) {
		return true
	}

	return false
}

// Appends stations to builder.station
func (b *WeatherDataQueryBuilder) Stations(stations ...string) *WeatherDataQueryBuilder {
	for _, s := range stations {
		b.stations = append(b.stations, s)
	}

	return b
}

// Appends data to query builder.data
func (b *WeatherDataQueryBuilder) Data(data ...WeatherDataData) *WeatherDataQueryBuilder {
	for _, d := range data {
		b.data = append(b.data, d)
	}

	return b
}

// Sets query builder start date (defaults to today)
func (b *WeatherDataQueryBuilder) Start(t time.Time) *WeatherDataQueryBuilder {
	b.start = t
	return b
}

// Sets query builder end date (defaults to today)
func (b *WeatherDataQueryBuilder) End(t time.Time) *WeatherDataQueryBuilder {
	b.end = t
	return b
}

// Sets query builder end date (defaults to Etc/UTC)
func (b *WeatherDataQueryBuilder) Timezone(tz string) *WeatherDataQueryBuilder {
	b.tz = tz
	return b
}

// Sets query builder format (defaults to Etc/UTC)
func (b *WeatherDataQueryBuilder) Format(format WeatherDataQueryFormat) *WeatherDataQueryBuilder {
	b.format = format
	return b
}

// Sets query builder latlon property (defaults to false)
func (b *WeatherDataQueryBuilder) LatLon(latlon bool) *WeatherDataQueryBuilder {
	b.latlon = latlon

	return b
}

// Sets query builder elev property (defaults to false)
func (b *WeatherDataQueryBuilder) Elevation(elev bool) *WeatherDataQueryBuilder {
	b.elev = elev

	return b
}

// Sets query builder missing property (defaults to M)
func (b *WeatherDataQueryBuilder) Missing(missing WeatherDataQueryMissing) *WeatherDataQueryBuilder {
	b.missing = missing

	return b
}

// Sets query builder trace property (defaults to T)
func (b *WeatherDataQueryBuilder) Trace(trace WeatherDataQueryTrace) *WeatherDataQueryBuilder {
	b.trace = trace

	return b
}

// Sets query builder trace property (defaults to T)
// func (b *WeatherDataQueryBuilder) Direct(direct bool) *WeatherDataQueryBuilder {
// 	b.direct = direct
//
// 	return b
// }

// Sets query builder report_type property (defaults to {3, 4})
func (b *WeatherDataQueryBuilder) ReportType(types ...int) *WeatherDataQueryBuilder {
	b.reportType = types

	return b
}

// Creates url.Values with validated data from query builder
func (b *WeatherDataQueryBuilder) BuildUrl() (url.Values, error) {
	v := url.Values{}

	// stations
	if b.stations == nil {
		return nil, requiredError("stations")
	}

	for _, s := range b.stations {
		v.Add("station", s)
	}

	// data
	if b.data == nil {
		return nil, requiredError("data")
	}

	for _, d := range b.data {
		v.Add("data", string(d))
	}

	// start
	v.Add("year1", strconv.Itoa(b.start.Year()))
	v.Add("month1", strconv.Itoa(int(b.start.Month())))
	v.Add("day1", strconv.Itoa(b.start.Day()))

	// end
	v.Add("year2", strconv.Itoa(b.end.Year()))
	v.Add("month2", strconv.Itoa(int(b.end.Month())))
	v.Add("day2", strconv.Itoa(b.end.Day()))

	// tz
	v.Add("tz", b.tz)

	// format
	v.Add("format", string(b.format))

	// latlon
	v.Add("latlon", boolToYesNo(b.latlon))

	// elev
	v.Add("elev", boolToYesNo(b.elev))

	// missing
	v.Add("missing", string(b.missing))

	// trace
	v.Add("trace", string(b.trace))

	// direct
	v.Add("direct", boolToYesNo(b.direct))

	// report_type

	for _, rt := range b.reportType {
		v.Add("report_type", strconv.Itoa(rt))
	}

	return v, nil
}

func boolToYesNo(b bool) string {
	if b {
		return string(Yes)
	} else {
		return string(No)
	}
}

// https://mesonet.agron.iastate.edu/cgi-bin/request/asos.py
// station=1HW
// ?data=tmpf multiple = (data=tmpf&data=tmpc)
// &year1=2023
// &month1=10
// &day1=3
// &year2=2023
// &month2=10
// &day2=4
// &tz=Etc%2FUTC
// &format=onlycomma
// &latlon=no
// &elev=no
// &missing=M
// &trace=T
// &direct=no
// &report_type=3
// &report_type=4
