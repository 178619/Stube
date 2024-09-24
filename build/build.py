# Build Go project for different platforms and create zip archive of project
# directory structure.

import os
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
    
    e = ''
    if os.name == 'nt':
        for y in env.split(' '):
            e += 'set '+y+'&&'
    else:
        e = env
        
    x = os.system('{env} go build -o bin/{bin} {src}'.format(
        env=e,
        bin=bin,
        src=src,
    ))
    
    if x != 0:
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
    print('Done.')
