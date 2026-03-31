#! /usr/bin/python3

# Build Go project for different platforms and create zip archive of project
# directory structure.

import os
import subprocess
import zipfile

def listdir(a):
    t = []
    with os.scandir(a) as s:
        for x in s:
            if x.is_dir():
                t += listdir(a+'/'+x.name)
            elif x.is_file():
                t.append(a+'/'+x.name)
    return t

def build(pkg, bin, env):
    src = 'github.com/178619/stube'
    try:
        subprocess.run(
            ['go', 'build', '-o', f'bin/{bin}', src],
            env=env,
            check=True
        )
    except subprocess.CalledProcessError:
        print('Error building ' + pkg)
        return
    
    with zipfile.ZipFile('bin/patch_' + pkg, mode='w') as f:
        f.write('bin/' + bin, arcname=bin)
        for filename in listdir('../static'):
            f.write('../static/' + filename, 'static/' + filename)
        for filename in listdir('../templates'):
            f.write('../templates/' + filename, 'templates/' + filename)
            
    with zipfile.ZipFile('bin/' + pkg, mode='w') as z:
        z.write('bin/' + bin, bin)
        z.write('../config.json', 'config.json')
        z.write('../README.md', 'README.md')
        z.write('../videos/README.md', 'videos/README.md')
        for filename in listdir('../static'):
            z.write('../static/' + filename, 'static/' + filename)
        for filename in listdir('../templates'):
            z.write('../templates/' + filename, 'templates/' + filename)

    # cleanup executable
    os.remove('bin/' + bin)
    print('Built ' + pkg)

if __name__ == "__main__":
    current_env = os.environ.copy()

    targets = [
        ('freebsd', '386'),
        ('freebsd', 'amd64'),
        ('freebsd', 'arm'),
        ('freebsd', 'arm64'),
        ('netbsd', '386'),
        ('netbsd', 'amd64'),
        ('netbsd', 'arm'),
        ('netbsd', 'arm64'),
        ('openbsd', '386'),
        ('openbsd', 'amd64'),
        ('openbsd', 'arm'),
        ('openbsd', 'arm64'),
        ('linux', '386'),
        ('linux', 'amd64'),
        ('linux', 'arm'),
        ('linux', 'arm64'),
        ('linux', 'riscv64'),
        ('darwin', 'amd64'),
        ('darwin', 'arm64'),
        ('windows', '386'),
        ('windows', 'amd64'),
        ('windows', 'arm64')
    ]

    for OS, ARCH in targets:
        env = current_env | {'CGO_ENABLED': '0', 'GOOS': OS, 'GOARCH': ARCH}
        if ARCH == 'arm':
            ARCH = 'armv6'
            env['GOARM'] = '6'
        elif ARCH == '386': ARCH = 'i386'
        if OS == 'darwin': OS = 'macos'
        build(
            pkg=f'tube_{OS}_{ARCH}.zip',
            bin='tube.exe' if OS == 'windows' else 'tube',
            env=env
        )
        
    print('Done.')
