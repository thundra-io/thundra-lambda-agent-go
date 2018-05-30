package plugin

var ApplicationName string
var ApplicationId string
var ApplicationVersion string
var ApplicationProfile string
var TransactionId string
var Region string
var MemorySize int

func init() {
	ApplicationName = getApplicationName()
	ApplicationId = getAppId()
	ApplicationVersion = getApplicationVersion()
	ApplicationProfile = getApplicationProfile()
	Region = getRegion()
	MemorySize = getMemorySize()
}
