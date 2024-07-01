# Build Go project for different platforms and create zip archive of project
# directory structure.

import os
import zipfile

def build(pkg, bin, env):
    src = 'github.com/178619/tube'
    if os.name == 'nt':
        e = ''
        for y in env.split(' '):
            e += 'set '+y+'&&'
        x = os.system('{env} go build -o bin/{bin} {src}'.format(
            env=e,
            bin=bin,
            src=src,
        ))
    else:
        x = os.system('{env} go build -o bin/{bin} {src}'.format(
            env=env,
            bin=bin,
            src=src,
        ))
    if x != 0:
        print('Error building ' + pkg)
        return
    z = zipfile.ZipFile('bin/' + pkg, mode='w')
    z.write('bin/' + bin, arcname=bin)
    z.write('../config.json', arcname='config.json')
    z.write('../README.md', 'README.md')
    z.write('../videos/README.md', 'videos/README.md')
    for filename in os.listdir('../static'):
        z.write('../static/' + filename, 'static/' + filename)
    for filename in os.listdir('../templates'):
        z.write('../templates/' + filename, 'templates/' + filename)
    # cleanup executable
    os.remove('bin/' + bin)
    print('Built ' + pkg)

if __name__ == "__main__":
    build(
        pkg='tube_linux.zip',
        bin='tube',
        env='GOOS=linux GOARCH=amd64',
    )
    build(
        pkg='tube_windows.zip',
        bin='tube.exe',
        env='GOOS=windows GOARCH=amd64',
    )
    build(
        pkg='tube_osx.zip',
        bin='tube',
        env='GOOS=darwin GOARCH=amd64',
    )
    build(
        pkg='tube_arm6.zip',
        bin='tube',
        env='GOOS=linux GOARCH=arm GOARM=6',
    )
    input('Done.')
