package types

type BuildingInfo struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

var Buildings = map[string]BuildingInfo{
	// === FAIRFAX CAMPUS ===
	"AB":     {Name: "Art and Design Building", Lat: 38.8285, Lng: -77.3094},
	"ACGC":   {Name: "Angel Cabrera Global Center", Lat: 38.8355, Lng: -77.3005},
	"AFC":    {Name: "Aquatic and Fitness Center", Lat: 38.8263, Lng: -77.3115},
	"AQ":     {Name: "Aquia Building", Lat: 38.8311, Lng: -77.3072},
	"BL":     {Name: "Blue Ridge Hall", Lat: 38.8332, Lng: -77.3055},
	"BUCHAN": {Name: "Buchanan Hall", Lat: 38.8298, Lng: -77.3088},
	"CAROW":  {Name: "Carow Hall", Lat: 38.8268, Lng: -77.3082},
	"CFA":    {Name: "Center for the Arts", Lat: 38.8290, Lng: -77.3089},
	"CH":     {Name: "College Hall", Lat: 38.8324, Lng: -77.3085},
	"DK":     {Name: "David J. King Hall", Lat: 38.8306, Lng: -77.3060},
	"E":      {Name: "East Building", Lat: 38.8327, Lng: -77.3088},
	"ENGR":   {Name: "Nguyen Engineering Building", Lat: 38.8270, Lng: -77.3054},
	"ENT":    {Name: "Enterprise Hall", Lat: 38.8276, Lng: -77.3059},
	"ESNHWR": {Name: "Eisenhower", Lat: 38.8340, Lng: -77.3060},
	"ESTSHR": {Name: "Eastern Shore", Lat: 38.8325, Lng: -77.3045},
	"EXPL":   {Name: "Exploratory Hall", Lat: 38.8288, Lng: -77.3059},
	"FENWCK": {Name: "Fenwick Library", Lat: 38.8306, Lng: -77.3076},
	"FH":     {Name: "Field House", Lat: 38.8250, Lng: -77.3150},
	"FIELD":  {Name: "Athletic Fields (West Campus)", Lat: 38.8320, Lng: -77.3250},
	"FINLEY": {Name: "Finley Building", Lat: 38.8320, Lng: -77.3080},
	"HNOVR":  {Name: "Hanover Hall", Lat: 38.8350, Lng: -77.3065},
	"HORIZN": {Name: "Horizon Hall", Lat: 38.8294, Lng: -77.3065},
	"HR":     {Name: "Hampton Roads", Lat: 38.8335, Lng: -77.3040},
	"HT":     {Name: "Harris Theater", Lat: 38.8300, Lng: -77.3085},
	"HUB":    {Name: "The Hub", Lat: 38.8299, Lng: -77.3051},
	"IN":     {Name: "Innovation Hall", Lat: 38.8282, Lng: -77.3073},
	"JC":     {Name: "Johnson Center", Lat: 38.8299, Lng: -77.3074},
	"KB":     {Name: "Krasnow Building", Lat: 38.8260, Lng: -77.3020},
	"KH":     {Name: "Krug Hall", Lat: 38.8319, Lng: -77.3064},
	"LH":     {Name: "Lecture Hall", Lat: 38.8315, Lng: -77.3085},
	"MAINST": {Name: "Main Street (9900 Main Street)", Lat: 38.8360, Lng: -77.3000},
	"MERTEN": {Name: "Merten Hall", Lat: 38.8348, Lng: -77.3087},
	"MTB":    {Name: "Music Theater Building", Lat: 38.8289, Lng: -77.3089},
	"PAB":    {Name: "de Laski Performing Arts Building", Lat: 38.8289, Lng: -77.3089},
	"PETRSN": {Name: "Peterson Hall", Lat: 38.8335, Lng: -77.3081},
	"PIEDMT": {Name: "Piedmont Hall", Lat: 38.8328, Lng: -77.3035},
	"PLANET": {Name: "Planetary Hall", Lat: 38.8283, Lng: -77.3053},
	"RAC":    {Name: "Recreation Athletic Complex", Lat: 38.8305, Lng: -77.3130},
	"ROGER":  {Name: "Roger Hall", Lat: 38.8265, Lng: -77.3065},
	"RSCH":   {Name: "Research Hall", Lat: 38.8276, Lng: -77.3049},
	"SNDBGE": {Name: "Sandbridge Hall", Lat: 38.8330, Lng: -77.3038},
	"SUBI":   {Name: "Student Union I", Lat: 38.8315, Lng: -77.3068},
	"T":      {Name: "Thompson Hall", Lat: 38.8327, Lng: -77.3099},
	"W":      {Name: "West Building", Lat: 38.8325, Lng: -77.3091},

	// === ARLINGTON CAMPUS ===
	"ARL1":   {Name: "Hazel Hall", Lat: 38.8850, Lng: -77.1030},
	"ARLVM":  {Name: "Van Metre Hall", Lat: 38.8845, Lng: -77.1025},
	"ARLVSH": {Name: "Vernon Smith Hall", Lat: 38.8842, Lng: -77.1020},
	"ARFUSE": {Name: "Fuse at Mason Square", Lat: 38.8848, Lng: -77.1035},

	// === SCITECH CAMPUS ===
	"PW-ABR":  {Name: "Advanced Biomedical Research", Lat: 38.7580, Lng: -77.5220},
	"PW-CH":   {Name: "Colgan Hall", Lat: 38.7590, Lng: -77.5230},
	"PW-DH":   {Name: "Discovery Hall", Lat: 38.7585, Lng: -77.5215},
	"PW-FC":   {Name: "Freedom Aquatic Center", Lat: 38.7550, Lng: -77.5180},
	"PW-KJH":  {Name: "Katherine Johnson Hall", Lat: 38.7595, Lng: -77.5225},
	"PW-LSEB": {Name: "Life Sciences and Engineering", Lat: 38.7588, Lng: -77.5210},

	// === OTHER ===
	"C":          {Name: "Commerce Building", Lat: 38.8500, Lng: -77.3000}, // Approx
	"OFF_CAMPUS": {Name: "Off Campus", Lat: 0.0, Lng: 0.0},
	"ON_LINE":    {Name: "Online", Lat: 0.0, Lng: 0.0},
}
