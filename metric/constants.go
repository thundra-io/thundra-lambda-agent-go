package metric

const metricType = "Metric"
const cpuMetric = "CPUMetric"
const gcMetric = "GcMetric"
const goroutineMetric = "GoroutineMetric"
const heapMetric = "HeapMetric"
const memoryMetric = "MemoryMetric"
const diskMetric = "DiskMetric"
const netMetric = "NetMetric"

// CPU Metrics
const appCpuLoad = "app.cpuLoad"
const sysCpuLoad = "sys.cpuLoad"

// GC Metrics
const pauseTotalNs = "pauseTotalNs"
const pauseNs = "pauseNs"
const numGc = "numGC"
const nextGc = "nextGC"
const gcCpuFraction = "gcCPUFraction"
const deltaNumGc = "deltaNumGC"
const deltaPauseTotalNs = "deltaPauseTotalNs"

// Disk Metrics
const readBytes = "readBytes"
const writeBytes = "writeBytes"
const readCount = "readCount"
const writeCount = "writeCount"

// GC Metrics
const numGoroutine = "numGoroutine"

// Heap Metrics
const heapAlloc = "heapAlloc"
const heapSys = "heapSys"
const heapInuse = "heapInuse"
const heapObjects = "heapObjects"
const memoryPercent = "memoryPercent"

// Net Metrics
const bytesRecv = "bytesRecv"
const bytesSent = "bytesSent"
const packetsRecv = "packetsRecv"
const packetsSent = "packetsSent"
const errIn = "errIn"
const errOut = "errOut"

// Memory Metrics
const appUsedMemory = "app.usedMemory"
const appMaxMemory = "app.maxMemory"
const sysUsedMemory = "sys.usedMemory"
const sysMaxMemory = "sys.maxMemory"
