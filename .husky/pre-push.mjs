import { spawnSync } from 'node:child_process'

const res = spawnSync('npm', ['run', 'test-race'], { stdio: 'inherit', shell: process.platform === 'win32' })
if (res.error) {
  console.error(res.error.message)
  process.exit(1)
}
process.exit(res.status ?? 1)
