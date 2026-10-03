// Runs the Go test suite with the race detector against the external MongoDB
// backend. TEST_MONGO_URI defaults to mongodb://127.0.0.1:27017 when unset.
import { spawnSync } from 'node:child_process'

const pkgs = ['.', './cmd/...', './common/...', './internal/...', './modules/...']
const uri = process.env.TEST_MONGO_URI || 'mongodb://127.0.0.1:27017'

const result = spawnSync('go', ['test', '-race', ...pkgs], {
  stdio: 'inherit',
  env: { ...process.env, TEST_MONGO_URI: uri },
})

process.exit(result.status ?? 1)
