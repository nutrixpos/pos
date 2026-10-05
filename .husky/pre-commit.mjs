import { spawnSync } from 'node:child_process'

function run(cmd, args, opts = {}) {
  const res = spawnSync(cmd, args, { stdio: 'inherit', ...opts })
  if (res.error) {
    console.error(res.error.message)
    process.exit(1)
  }
  if (res.status !== 0) process.exit(res.status)
}

function listFiles(gitArgs) {
  const res = spawnSync('git', gitArgs)
  if (res.status !== 0) {
    console.error(res.stderr.toString())
    process.exit(res.status ?? 1)
  }
  return res.stdout.toString().split('\0').filter(Boolean)
}

const staged = listFiles(['diff', '--cached', '--name-only', '--diff-filter=ACMR', '-z', '--', '*.go'])
if (staged.length === 0) process.exit(0)

// Refuse to run when a staged file also has unstaged edits. gofmt rewrites the
// whole working-tree file and `git add` would then stage those unrelated
// changes, breaking partial staging. Stop with a clear error instead.
const unstaged = new Set(listFiles(['diff', '--name-only', '-z', '--', '*.go']))
const partiallyStaged = staged.filter((file) => unstaged.has(file))
if (partiallyStaged.length > 0) {
  console.error('pre-commit: these Go files have both staged and unstaged changes:')
  for (const file of partiallyStaged) console.error('  ' + file)
  console.error('Stage the remaining changes or stash them, then commit again.')
  console.error('(gofmt would otherwise reformat and stage the whole working-tree file.)')
  process.exit(1)
}

run('gofmt', ['-w', ...staged])
run('git', ['add', ...staged])

run('npm', ['run', 'vet'], { shell: process.platform === 'win32' })
run('npm', ['run', 'test'], { shell: process.platform === 'win32' })
