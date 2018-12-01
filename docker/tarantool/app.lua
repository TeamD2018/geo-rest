local os = require('os')

box.cfg {
    memtx_max_tuple_size = os.getenv(MEMTX_MAX_TUPLE_SIZE) or 10 * 1024 * 1024,
    memtx_memory = os.getenv(MEMTX_MEMORY) or 1024 * 1024 * 1024,
    listen = 3301,
    checkpoint_interval = 10 * 60,
    feedback_enabled = false,
}