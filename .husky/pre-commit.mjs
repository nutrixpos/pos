import { spawnSync } from 'node:child_process'

function run(cmd, args, opts = {}) {
  const res = spawnSync(cmd, args, { stdio: 'inherit', ...opts })
  if (res.error) {
    console.error(res.error.message)
    process.exit(1)
  }
  if (res.status !== 0) process.exit(res.status)
}

const staged = spawnSync('git', ['diff', '--cached', '--name-only', '--diff-filter=ACMR', '-z', '--', '*.go'])
if (staged.status !== 0) {
  console.error(staged.stderr.toString())
  process.exit(staged.status ?? 1)
}

const files = staged.stdout.toString().split('\0').filter(Boolean)
if (files.length === 0) process.exit(0)

run('gofmt', ['-w', ...files])
run('git', ['add', ...files])

run('npm', ['run', 'vet'], { shell: process.platform === 'win32' })
run('npm', ['run', 'test'], { shell: process.platform === 'win32' })
