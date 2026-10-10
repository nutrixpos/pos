// Enforces a minimum statement-coverage threshold over the packages that are
// in scope for the project coverage goal.
//
// Usage: node .husky/coverage-threshold.mjs [threshold] [profile]
//
// The coverage profile is produced with `go test -covermode=atomic
// -coverpkg=<in-scope packages> -coverprofile=coverage.out ...`. When using
// -coverpkg the same block can appear once per test binary, so blocks are
// de-duplicated by position keeping the highest hit count.
//
// Out-of-scope packages and the hard-to-test files below are ignored so they do
// not affect the metric (see codecov.yml for the matching report exclusions).

import { readFileSync } from 'node:fs'

const threshold = Number(process.argv[2] ?? '80')
const profilePath = process.argv[3] ?? 'coverage.out'

const includePackages = [
  '/common/config/',
  '/common/helpers/',
  '/modules/auth/middlewares/',
  '/modules/core/handlers/',
  '/modules/core/services/',
]

const excludeFiles = [
  '/modules/core/services/receipt.go',
  '/modules/core/services/notification_singleton.go',
  '/modules/core/handlers/notifications.go',
]

const inScope = (file) =>
  includePackages.some((p) => file.includes(p)) && !excludeFiles.some((f) => file.includes(f))

let raw
try {
  raw = readFileSync(profilePath, 'utf8')
} catch (err) {
  console.error(`coverage-threshold: cannot read profile ${profilePath}: ${err.message}`)
  process.exit(1)
}

// key = "file:start,end" -> { numStmts, count } keeping the highest count.
const blocks = new Map()
for (const line of raw.split('\n')) {
  const trimmed = line.trim()
  if (trimmed === '' || trimmed.startsWith('mode:')) continue

  const match = trimmed.match(/^(\S+):(\d+\.\d+),(\d+\.\d+) (\d+) (\d+)$/)
  if (!match) continue

  const [, file, start, end, numStmtsStr, countStr] = match
  if (!inScope(file)) continue

  const key = `${file}:${start},${end}`
  const numStmts = Number(numStmtsStr)
  const count = Number(countStr)
  const existing = blocks.get(key)
  if (!existing || count > existing.count) {
    blocks.set(key, { numStmts, count })
  }
}

let covered = 0
let total = 0
for (const { numStmts, count } of blocks.values()) {
  total += numStmts
  if (count > 0) covered += numStmts
}

if (total === 0) {
  console.error('coverage-threshold: no in-scope statements found in the profile')
  process.exit(1)
}

const pct = (covered / total) * 100
console.log(`coverage (in scope): ${pct.toFixed(1)}% (${covered}/${total} statements)`)

if (pct + 1e-9 < threshold) {
  console.error(`coverage-threshold: ${pct.toFixed(1)}% is below the required ${threshold}%`)
  process.exit(1)
}
