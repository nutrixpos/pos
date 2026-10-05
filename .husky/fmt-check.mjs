// Fails when any tracked Go file is not gofmt-compliant.
import { spawnSync } from 'node:child_process'

const listed = spawnSync('git', ['ls-files', '-z', '--', '*.go'])
if (listed.status !== 0) {
  console.error(listed.stderr.toString())
  process.exit(listed.status ?? 1)
}

const files = listed.stdout.toString().split('\0').filter(Boolean)
if (files.length === 0) {
  process.exit(0)
}

const check = spawnSync('gofmt', ['-l', ...files])
const unformatted = check.stdout.toString().trim()
if (unformatted !== '') {
  console.error('gofmt needed for:\n' + unformatted)
  process.exit(1)
}

process.exit(check.status ?? 0)
