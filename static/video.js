const tempStyle = document.createElement('style')
tempStyle.innerHTML = 'video {visibility: hidden;}'
document.head.appendChild(tempStyle)
const init = () => {
    const miniAlert = (v) => {
        const d = document.createElement('div')
        d.classList.add('minialert')
        d.innerHTML = v
        document.getElementById("mask").appendChild(d)
        setTimeout(()=>{if (d) d.style.opacity = 0}, 2250)
        setTimeout(()=>{if (d) d.remove()}, 3000)
        return d
    }
    const oneAlertHandler = []
    const oneAlert = (v) => {
        if (!oneAlertHandler.length) return
        oneAlertHandler.push(miniAlert(v))
        if (oneAlertHandler[0]) oneAlertHandler[0].remove()
        oneAlertHandler.shift()
    }
    const openMenu = (t, v) => {
        const d = document.createElement('div')
        d.classList.add('menu')
        const head = document.createElement('div')
        head.classList.add('head')
        head.innerText = t
        const closer = document.createElement('button')
        closer.classList.add('closer')
        closer.onclick = (e) => {
            if (e.pointerType == 'mouse' && e.button != 0) return
            removeMenus()
        }
        head.appendChild(closer)
        d.appendChild(head)
        d.appendChild(v)
        document.getElementById("mask").appendChild(d)
        return d
    }
    const setMenuHandler = []
    const setMenu = (t, v) => {
        if (setMenuHandler.length) {
            setMenuHandler[0].remove()
            setMenuHandler.shift()
        }
        setMenuHandler.push(openMenu(t, v))
    }
    const removeMenus = () => {
        while (setMenuHandler.length) {
            setMenuHandler[0].remove()
            setMenuHandler.shift()
        }
    }
    let pointerWaitTime = 0
    const keyList = [
        " ", "ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown", ",", ".", ">", "<", "Home", "End",
        "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "/", "f", "c", "k", "m", "j", "l"
    ]
    const getTimeString = (v) => {
        return (v >= 3600 ? Math.floor(v / 3600) + ':' : '')
        + (Math.floor(v % 3600 / 60) < 10 ? '0' : '') + Math.floor(v % 3600 / 60) + ':'
        + (Math.floor(v % 60) < 10 ? '0' : '') + Math.floor(v % 60)
    }
    const languageNames = new Intl.DisplayNames(navigator.languages, { type: 'language' });
    if (document.getElementById('search')) document.getElementById('search').onkeyup = () => {
        const key = document.getElementById('search').value.toLocaleLowerCase()
        document.querySelectorAll('#playlist > a').forEach((v)=>{
            if (v.title.toLocaleLowerCase().includes(key) || v.name.toLocaleLowerCase().includes(key)) v.className = ''; else v.className = 'hidden';
        })
    }
    const mask = document.getElementById("mask")
    if (!mask) {
        tempStyle.remove()
        return
    }
    mask.style.display = 'block'
    const video = document.getElementById("video")
    if (video.hasAttribute('controls')) video.removeAttribute('controls')
    const player = document.getElementById("player")
    const playOrPause = () => {
        if (!!(video.currentTime > 0 && !video.paused && !video.ended && video.readyState > 2)) {
            video.pause()
        } else {
            video.play()
        }
    }
    const getFullscreen = () => {
        if (document.fullscreenElement) {
            document.exitFullscreen()
        } else {
            document.body.requestFullscreen()
        }
    }
    const addTime = (v) => {
        let q = video.currentTime + v
        if (video.loop) {
            while (q > video.duration) q -= video.duration
            while (q < 0) q += video.duration
        }
        video.currentTime = q
        oneAlert('To: '+getTimeString(video.currentTime))
    }

    const locationOptions = Object.fromEntries(location.search.slice(1).split('&').filter((v)=>{return v.length}).map((v)=>{return [v.split('=')[0], v.split('=').slice(1).join('=')]}))

    document.body.onkeydown = (e) => {
        if (!keyList.includes(e.key) || e.target.id == 'search' || e.ctrlKey) return
        e.preventDefault()
        switch(e.key) {
            case " ":
            case "k":
                document.getElementById('center').style.transform = 'scale(0.875)'
                playOrPause()
                setTimeout(()=>{document.getElementById('center').style.transform = null}, 100)
                break
            case "j":
                addTime(-5)
            case "ArrowLeft":
                addTime(-5)
                break
            case "l":
                addTime(5)
            case "ArrowRight":
                addTime(5)
                break
            case "0":
            case "1":
            case "2":
            case "3":
            case "4":
            case "5":
            case "6":
            case "7":
            case "8":
            case "9":
                video.currentTime = video.duration / 10 * parseInt(e.key)
                oneAlert('To: '+getTimeString(video.currentTime))
                break
            case "Home":
                video.currentTime = 0
                oneAlert('To: '+getTimeString(video.currentTime))
                break
            case "End":
                video.currentTime = video.duration
                oneAlert('To: '+getTimeString(video.duration))
                break
            case "<":
                if (video.playbackRate > 0.5) video.playbackRate = Math.round(video.playbackRate * 4 - 1) / 4; else video.playbackRate = 0.25;
                oneAlert('Playspeed: ' + video.playbackRate + 'x')
                break
            case ">":
                if (video.playbackRate < 3.75) video.playbackRate = Math.round(video.playbackRate * 4 + 1) / 4;  else video.playbackRate = 4;
                oneAlert('Playspeed: ' + video.playbackRate + 'x')
                break
            case "m":
                video.muted = !video.muted
                updateVolume()
                break
            case "ArrowUp":
                if (video.volume < 0.95) video.volume = Math.round(video.volume * 20 + 1) / 20; else video.volume = 1;
                updateVolume()
                break
            case "ArrowDown":
                if (video.volume > 0.05) video.volume = Math.round(video.volume * 20 - 1) / 20; else video.volume = 0;
                updateVolume()
                break
            case "f":
                getFullscreen()
                break
        }
        updateSeeker()
    }
    mask.onmouseenter = (e) => {
        if (!mask.classList.toString().includes('shown')) {
            updateSeeker()
            mask.classList.add('shown')
        }
    }
    mask.onmouseleave = (e) => {
        if (mask.classList.toString().includes('shown')) {
            document.getElementById('center').style.transform = null
            mask.classList.remove('shown')
        }
    }
    mask.onpointerdown = (e) => {
        if (e.target.id != 'mask' || e.pointerType == 'mouse' && e.button != 0) return
        if (e.pointerType == 'touch') {
            if (mask.classList.toString().includes('shown')) {
                mask.classList.remove('shown')
            } else {
                updateSeeker()
                mask.classList.add('shown')
            }
        } else {
            playOrPause()
            document.getElementById('center').style.transform = 'scale(0.875)'
        }
    }
    mask.onpointerup = (e) => {
        if (e.pointerType != "touch") {
            document.getElementById('center').style.transform = null
        }
    }
    const hideMask = () => {setTimeout(()=>{if (pointerWaitTime <= 0) mask.classList.remove('shown'); else hideMask()}, 100)}
    mask.onpointermove = (e) => {
        if (e.pointerType == "touch") return
        updateSeeker()
        mask.classList.add('shown')
        pointerWaitTime = 3000
        hideMask()
    }
    mask.ondblclick = (e) => {
        if (e.target.id != 'mask') return
        e.preventDefault()
        if (e.layerX / e.target.clientWidth < 2 / 5) {
            addTime(-10)
            updateSeeker()
        } else if (e.layerX / e.target.clientWidth > 3 / 5) {
            addTime(10)
            updateSeeker()
        }
    }
    document.querySelector('#mask h1').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        if (!document.querySelector('.embed')) return
        open(location.origin + location.pathname.replace('/e/', '/v/'))
    }
    document.getElementById('center').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        playOrPause()
        document.getElementById('center').style.transform = 'scale(0.875)'
    }
    document.getElementById('center').onpointerup = () => {
        document.getElementById('center').style.transform = null
    }
    document.getElementById('center').onpointerleave = () => {
        document.getElementById('center').style.transform = null
    }
    document.getElementById('seeker').onwheel = (e) => {
        if (e.ctrlKey) return
        e.preventDefault()
        addTime(e.wheelDeltaY / 100)
    }
    document.getElementById('seeker').oninput = () => {
        video.currentTime = document.getElementById('seeker').value / 1000
        oneAlert('To: '+getTimeString(video.currentTime))
    }
    document.getElementById('play-pause').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        playOrPause()
    }
    document.getElementById('loop').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        video.loop = !video.loop
        if (!video.loop) {
            document.getElementById("loop").style.backgroundImage = 'url(/static/icons/repeat-off.svg)'
            oneAlert('Loop Turned Off')
        } else if (true) {
            document.getElementById("loop").style.backgroundImage = 'url(/static/icons/repeat-once.svg)'
            oneAlert('Loop Turned On')
        }
    }
    document.getElementById('volume-icon').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        video.muted = !video.muted
        updateVolume()
    }
    document.getElementById('volume-icon').onwheel = (e) => {
        if (e.ctrlKey) return
        e.preventDefault()
        if (e.deltaY < 0) {
            if (video.volume < 0.95) video.volume = Math.round(video.volume * 20 + 1) / 20; else video.volume = 1;
        } else if (e.deltaY > 0) {
            if (video.volume > 0.05) video.volume = Math.round(video.volume * 20 - 1) / 20; else video.volume = 0;
        }
        document.getElementById('volume-icon').style.transform = 'scale(0.875)'
        setTimeout(()=>{document.getElementById('volume-icon').style.transform = null}, 100)
        updateVolume()
    }
    document.getElementById('volume').oninput = () => {
        video.volume = document.getElementById('volume').value
        updateVolume()
    }
    document.getElementById('volume').onwheel = (e) => {
        if (e.ctrlKey) return
        e.preventDefault()
        if (e.deltaY < 0) {
            if (video.volume < 0.95) video.volume = Math.round(video.volume * 20 + 1) / 20; else video.volume = 1;
        } else if (e.deltaY > 0) {
            if (video.volume > 0.05) video.volume = Math.round(video.volume * 20 - 1) / 20; else video.volume = 0;
        }
        updateVolume()
    }
    if (!video.textTracks.length) document.getElementById('captions').style.display = 'none'
    document.getElementById('captions').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        const tracks = Array.from(video.textTracks)
        const ul = document.createElement('ul')
        tracks.forEach((v)=>{
            const li = document.createElement('li')
            try {
                li.innerText = languageNames.of(v.language)
            } catch {
                li.innerText = v.language
            }
            if (v.mode == 'showing') li.classList.add('enabled')
            li.onpointerdown = (e) => {
                if (e.pointerType == 'mouse' && e.button != 0) return
                if (v.mode == 'showing') {
                    v.mode = 'disabled'
                    li.classList.remove('enabled')
                } else {
                    v.mode = 'showing'
                    li.classList.add('enabled')
                }
            }
            ul.appendChild(li)
        })
        setMenu('Captions', ul)
    }
    document.getElementById('screenshot').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        if (!video.videoWidth || !video.videoHeight) {
            const a = document.createElement('a')
            a.href = '/t/' + location.pathname.slice(3)
            a.download = decodeURIComponent(location.pathname.split('/').pop())+'.svg'
            document.body.appendChild(a)
            a.click()
            a.remove()
            oneAlert('Screenshot Taken')
            return
        }
        const canvas = document.createElement('canvas')
        canvas.width = video.videoWidth
        canvas.height = video.videoHeight
        const ctx = canvas.getContext('2d')
        ctx.drawImage(video, 0, 0, video.videoWidth, video.videoHeight)
        const a = document.createElement('a')
        a.href = canvas.toDataURL('image/png')
        a.download = decodeURIComponent(location.pathname.split('/').pop())+'-'+getTimeString(video.currentTime)+(video.currentTime%1+'').replace('.',':').slice(1,5)+'.svg'
        document.body.appendChild(a)
        a.click()
        a.remove()
        oneAlert('Screenshot Taken')
    }
    const speedList = [0.25, 0.5, 0.75, 1, 1.25, 1.5, 2, 2.5, 3, 4]
    document.getElementById('playspeed').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        const ul = document.createElement('ul')
        speedList.forEach((v)=>{
            const li = document.createElement('li')
            li.innerText = v + 'x'
            li.onpointerdown = (e) => {
                if (e.pointerType == 'mouse' && e.button != 0) return
                video.playbackRate = v
                oneAlert('Playspeed: ' + video.playbackRate + 'x')
            }
            ul.appendChild(li)
        })
        setMenu('Playspeed', ul)
    }
    document.getElementById('playspeed').onwheel = (e) => {
        if (e.ctrlKey) return
        e.preventDefault()
        if (e.deltaY < 0) {
            if (video.playbackRate < 3.75) video.playbackRate = Math.round(video.playbackRate * 4 + 1) / 4; else video.playbackRate = 4;
        } else if (e.deltaY > 0) {
            if (video.playbackRate > 0.5) video.playbackRate = Math.round(video.playbackRate * 4 - 1) / 4; else video.playbackRate = 0.25;
        }
        oneAlert('Playspeed: ' + video.playbackRate + 'x')
        document.getElementById('playspeed').style.transform = 'scale(0.875)'
        setTimeout(()=>{document.getElementById('playspeed').style.transform = null}, 100)
    }
    document.getElementById('embedlink').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button == 1) {
            open(location.origin + location.pathname.replace('/v/', '/e/'))
            return
        }
        if (e.pointerType == 'mouse' && e.button != 0) return
        navigator.clipboard.writeText('<iframe src="'+location.origin+location.pathname.replace('/v/', '/e/')+'" width=320 height=180 allowfullscreen></iframe>')
        oneAlert('HTML code Copied to Clipboard.')
    }
    document.getElementById('filelink').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        const a = document.createElement('a')
        a.href = document.getElementById('filelink').getAttribute('href')
        document.body.appendChild(a)
        a.click()
        a.remove()
    }
    document.getElementById('current').style.display = 'block';
    document.getElementById('collapse').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        if (document.getElementById('current').style.display != 'block') {
            document.getElementById('current').style.display = 'block';
            document.getElementById('captions').style.display = '';
            document.getElementById('screenshot').style.display = '';
            document.getElementById('playspeed').style.display = '';
            document.getElementById('embedlink').style.display = '';
            document.getElementById('filelink').style.display = '';
            document.getElementById("collapse").style.backgroundImage = 'url(/static/icons/tune.svg)'
        } else {
            document.getElementById('current').style.display = '';
            document.getElementById('captions').style.display = 'block';
            document.getElementById('screenshot').style.display = 'block';
            document.getElementById('playspeed').style.display = 'block';
            document.getElementById('embedlink').style.display = 'block';
            document.getElementById('filelink').style.display = 'block';
            document.getElementById("collapse").style.backgroundImage = 'url(/static/icons/tune-vertical.svg)'
        }
    }
    document.getElementById('fullscreen').onpointerdown = (e) => {
        if (e.pointerType == 'mouse' && e.button != 0) return
        getFullscreen()
    }
    document.body.onfullscreenchange = (e) => {
        if (document.fullscreenElement) {
            player.classList.add('fullscreen')
        } else {
            player.classList.remove('fullscreen')
        }
    }
    const updateSeeker = () => {
        document.getElementById("seeker").max = video.duration * 1000
        document.getElementById("seeker").value = video.currentTime * 1000
        const ct = video.currentTime / video.duration
        let bt
        for (let i=0; i<video.buffered.length; i++) {
            if (video.buffered.start(i) < video.currentTime && video.currentTime < video.buffered.end(i)) {
                bt = video.buffered.end(i) / video.duration
                break
            }
        }
        if (!bt) bt = ct
        document.getElementById("seeker").style.backgroundImage = `linear-gradient(to right, rgba(0,0,0,0) 6px, var(--main-color) 6px, var(--main-color) calc(6px + ${ct * 100}% - 12px * ${ct}), rgba(223,223,223,0.75) calc(6px + ${ct * 100}% - 12px * ${ct}), rgba(223,223,223,0.75) calc(6px + ${bt * 100}% - 12px * ${bt}), rgba(191,191,191,0.5) calc(6px + ${bt * 100}% - 12px * ${bt}), rgba(191,191,191,0.5) calc(100% - 6px), rgba(191,191,191,0.5) calc(100% - 6px), rgba(0,0,0,0) calc(100% - 6px))`
        document.getElementById("current").innerText = getTimeString(video.currentTime) + ' / ' + getTimeString(video.duration)
        if (!!(video.currentTime > 0 && !video.paused && !video.ended && video.readyState > 2)) {
            document.getElementById("center").style.backgroundImage = 'url(/static/icons/pause.svg)'
            document.getElementById("play-pause").style.backgroundImage = 'url(/static/icons/pause.svg)'
        } else if (video.currentTime == video.duration) {
            document.getElementById("center").style.backgroundImage = 'url(/static/icons/replay.svg)'
            document.getElementById("play-pause").style.backgroundImage = 'url(/static/icons/replay.svg)'
        } else {
            document.getElementById("center").style.backgroundImage = 'url(/static/icons/play.svg)'
            document.getElementById("play-pause").style.backgroundImage = 'url(/static/icons/play.svg)'
        }
    }
    const updateVolume = () => {
        document.getElementById("volume").style.backgroundImage = `linear-gradient(to right, rgba(0,0,0,0) 6px, var(--main-color) 6px, var(--main-color) calc(6px + ${video.volume * 100}% - 12px * ${video.volume}), rgba(191,191,191,0.5) calc(6px + ${video.volume * 100}% - 12px * ${video.volume}), rgba(191,191,191,0.5) calc(100% - 6px), rgba(0,0,0,0) calc(100% - 6px))`
        document.getElementById("volume").value = video.volume
        if (video.muted) {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/icons/volume-off.svg)'
        } else if (video.volume == 0) {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/icons/volume-low.svg)'
        } else if (video.volume < 1/2) {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/icons/volume-medium.svg)'
        } else {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/icons/volume-high.svg)'
        }
        oneAlert('Volume: '+(video.muted ? 'Muted' : Math.floor(video.volume * 100)+'%'))
    }
    const repeat = () => setTimeout(()=>{updateSeeker(); pointerWaitTime -= 20; repeat()}, 20)
    updateVolume()
    tempStyle.remove()
    const current = document.querySelector('a.playing')
    if (current) current.scrollIntoView()
    document.body.scrollIntoView()
    document.body.click()
    oneAlertHandler.push(null)
    repeat()
}