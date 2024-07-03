const tempStyle = document.createElement('style')
tempStyle.innerHTML = 'video {visibility: hidden;}'
document.head.appendChild(tempStyle)
const init = () => {
    const keyList = [
        " ", "ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown", ",", ".", ">", "<", "Home", "End",
        "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "/", "f", "c", "k", "m", "j", "l"
    ]
    const mask = document.getElementById("mask")
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
        if (window.document.fullscreenElement) {
            document.exitFullscreen()
        } else {
            document.body.requestFullscreen()
        }
    }
    document.body.onkeydown = (e) => {
        if (!keyList.includes(e.key) || e.target.id == 'search' || e.ctrlKey) return;
        e.preventDefault();
        switch(e.key) {
            case " ":
            case "k":
                document.getElementById('center').style.transform = 'scale(0.875)'
                playOrPause()
                setTimeout(()=>{document.getElementById('center').style.transform = null}, 100)
                break
            case "j":
                video.currentTime -= 5
            case "ArrowLeft":
                video.currentTime -= 5
                break
            case "l":
                video.currentTime += 5
            case "ArrowRight":
                video.currentTime += 5
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
                break
            case "Home":
                video.currentTime = 0
                break
            case "End":
                video.currentTime = video.duration
                break
            case "<":
                video.playbackRate = Math.round(video.playbackRate * 4 - 1) / 4
                break
            case ">":
                video.playbackRate = Math.round(video.playbackRate * 4 + 1) / 4
                break
            case "m":
                video.muted = !video.muted
                updateVolume()
                break
            case "ArrowUp":
                if (video.volume < 0.95) video.volume = video.volume + 1 / 20; else video.volume = 1;
                updateVolume()
                break
            case "ArrowDown":
                if (video.volume > 0.05) video.volume = video.volume - 1 / 20; else video.volume = 0;
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
        if (e.target.id != 'mask') return
        if (video.currentTime == 0) {
            video.play()
            return
        }
        if (e.pointerType == "touch") {
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
    var pointerWaitTime = 0
    mask.onpointermove = (e) => {
        if (e.pointerType == "touch") return
        updateSeeker()
        mask.classList.add('shown')
        pointerWaitTime = 3000
        setTimeout(()=>{if (pointerWaitTime <= 0) mask.classList.remove('shown')}, 5000)
    }
    mask.ondblclick = (e) => {
        if (e.target.id != 'mask') return
        e.preventDefault()
        if (e.layerX / e.target.clientWidth < 2 / 5) {
            video.currentTime -= 10
            updateSeeker()
        } else if (e.layerX / e.target.clientWidth > 3 / 5) {
            video.currentTime += 10
            updateSeeker()
        }
    }
    document.getElementById('center').onpointerdown = () => {
        playOrPause()
        document.getElementById('center').style.transform = 'scale(0.875)'
    }
    document.getElementById('center').onpointerup = () => {
        document.getElementById('center').style.transform = null
    }
    document.getElementById('center').onpointerleave = () => {
        document.getElementById('center').style.transform = null
    }
    document.getElementById('play-pause').onmousedown = playOrPause
    document.getElementById('seeker').oninput = () => {
        video.currentTime = document.getElementById('seeker').value / 1000
    }

    document.getElementById('volume-icon').onmousedown = () => {
        video.muted = !video.muted
        updateVolume()
    }
    document.getElementById('volume').oninput = () => {
        video.volume = document.getElementById('volume').value
        updateVolume()
    }
    document.getElementById('fullscreen').onmousedown = getFullscreen
    document.body.onfullscreenchange = (e) => {
        if (window.document.fullscreenElement) {
            player.classList.add('fullscreen')
        } else {
            player.classList.remove('fullscreen')
        }
    }
    if (document.getElementById('search')) document.getElementById('search').onkeyup = () => {
        const key = document.getElementById('search').value.toLocaleLowerCase()
        document.querySelectorAll('#playlist > a').forEach((v)=>{
            if (v.title.toLocaleLowerCase().includes(key) || v.name.toLocaleLowerCase().includes(key)) v.className = ''; else v.className = 'hidden';
        })
    }
    const updateSeeker = () => {
        document.getElementById("seeker").max = video.duration * 1000
        document.getElementById("seeker").value = video.currentTime * 1000
        const ct = video.currentTime / video.duration
        var bt
        for (var i=0; i<video.buffered.length; i++) {
            if (video.buffered.start(i) < video.currentTime && video.currentTime < video.buffered.end(i)) {
                bt = video.buffered.end(i) / video.duration
                break
            }
        }
        if (!bt) bt = ct
        document.getElementById("seeker").style.backgroundImage = `linear-gradient(to right, rgba(0,0,0,0) 6px, var(--main-color) 6px, var(--main-color) calc(6px + ${ct * 100}% - 12px * ${ct}), rgba(223,223,223,0.75) calc(6px + ${ct * 100}% - 12px * ${ct}), rgba(223,223,223,0.75) calc(6px + ${bt * 100}% - 12px * ${bt}), rgba(191,191,191,0.5) calc(6px + ${bt * 100}% - 12px * ${bt}), rgba(191,191,191,0.5) calc(100% - 6px), rgba(191,191,191,0.5) calc(100% - 6px), rgba(0,0,0,0) calc(100% - 6px))`
        document.getElementById("current").innerText = (Math.floor(video.currentTime / 60) < 10 ? '0' : '') + Math.floor(video.currentTime / 60) + ':' + (Math.floor(video.currentTime % 60) < 10 ? '0' : '') + Math.floor(video.currentTime % 60)
        + ' / ' + (Math.floor(video.duration / 60) < 10 ? '0' : '') + Math.floor(video.duration / 60) + ':' + (Math.floor(video.duration % 60) < 10 ? '0' : '') + Math.floor(video.duration % 60)
        if (!!(video.currentTime > 0 && !video.paused && !video.ended && video.readyState > 2)) {
            document.getElementById("center").style.backgroundImage = 'url(/static/pause.png)'
            document.getElementById("play-pause").style.backgroundImage = 'url(/static/pause.png)'
        } else {
            document.getElementById("center").style.backgroundImage = 'url(/static/play.png)'
            document.getElementById("play-pause").style.backgroundImage = 'url(/static/play.png)'
        }
    }
    const updateVolume = () => {
        document.getElementById("volume").style.backgroundImage = `linear-gradient(to right, rgba(0,0,0,0) 6px, var(--main-color) 6px, var(--main-color) calc(6px + ${video.volume * 100}% - 12px * ${video.volume}), rgba(191,191,191,0.5) calc(6px + ${video.volume * 100}% - 12px * ${video.volume}), rgba(191,191,191,0.5) calc(100% - 6px), rgba(0,0,0,0) calc(100% - 6px))`
        document.getElementById("volume").value = video.volume
        if (video.muted) {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/volume-off.png)'
        } else if (video.volume == 0) {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/volume-low.png)'
        } else if (video.volume < 1/2) {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/volume-medium.png)'
        } else {
            document.getElementById("volume-icon").style.backgroundImage = 'url(/static/volume-high.png)'
        }
    }
    const repeat = () => setTimeout(()=>{updateSeeker(); pointerWaitTime -= 20; repeat()}, 20)
    updateVolume()
    tempStyle.remove()
    const current = document.querySelector('a.playing')
    if (current) current.scrollIntoView()
    repeat()
}