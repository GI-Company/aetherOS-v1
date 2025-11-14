package hostapi

import (
    "context"
    "log"
)

func RegisterHostAPIs(ctx context.Context, r interface{}) {
    // In wazero you'd use NewHostModuleBuilder to register functions
    // Example host methods:
    // host.kv_get(ptr, len) -> writes to module memory, etc.
    log.Println("host APIs would be registered here")
}