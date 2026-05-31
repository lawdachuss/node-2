package channel

// Global upload semaphore — limits total concurrent uploads across ALL channels
// to prevent bandwidth saturation.  Default: 10 concurrent uploads.
var UploadSem = make(chan struct{}, 10)
