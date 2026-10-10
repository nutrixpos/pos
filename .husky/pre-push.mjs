import { spawnSync } from 'node:child_process'

function run(script) {
  const res = spawnSync('npm', ['run', script], { stdio: 'inherit', shell: process.platform === 'win32' })
  if (res.error) {
    console.error(res.error.message)
    process.exit(1)
  }
  if (res.status !== 0) process.exit(res.status)
}

run('test-race')
run('vuln')
