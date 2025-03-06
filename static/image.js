const init = () => {
    if (location.pathname.endsWith('/')) window.history.replaceState('', null, window.location.origin + window.location.pathname.slice(0, -1) + window.location.search + window.location.hash);
    const keyList = [
        'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'PageUp', 'PageDown', 'Home', 'End'
    ]
    const list = document.querySelectorAll('#viewer img'), len = list.length
    let cp = 0, direction = 'ltr'
    if (location.hash) cp = parseInt(location.hash.slice(1))-1
    if (cp < 0 || cp >= len) cp = 0;
    const range = document.createElement('input')
    range.id = 'range'
    range.type = 'range'
    range.min = 0
    range.max = len - 1
    range.step = 1
    range.value = cp
    range.style.direction = direction
    document.querySelector('nav').appendChild(range)
    window.history.replaceState(null, null, window.location.origin + window.location.pathname + window.location.search + '#' + (cp+1));
    const load = (smooth) => {
        toPage(cp, smooth)
    }
    const toPage = (p, smooth) => {
        if (list[p]) {
            // document.getElementById('number').innerHTML = p + 1 + '<br>' + len
            range.value = p;
            // if (document.getElementById("img_"+p).complete || navigator.userAgent.indexOf("Firefox") != -1) {
            if (list[p].complete) {
                list[p].scrollIntoView(smooth ? {behavior: 'smooth'} : undefined);
            } else {
                // document.getElementById("loading").style.display = 'block'
                list[p].onload = function() {
                    list[p].scrollIntoView();
                    // document.getElementById("loading").style.display = 'none'
                    list[p].onload = null;
                }
            }
            window.history.replaceState(null, null, window.location.origin + window.location.pathname + window.location.search + '#' + (p+1));
        }
    }
    const prev = () => {
        if (direction == 'rtl' && cp > 0) cp-=1; else if (direction != 'rtl' && cp < len-1) cp+=1;
        load();
    }
    const next = () => {
        if (direction == 'rtl' && cp < len-1) cp+=1; else if (direction != 'rtl' && cp > 0) cp-=1;
        load();
    }
    document.onkeydown = function(e) {
        if (!keyList.includes(e.key) || e.ctrlKey) return
        e.preventDefault();
        if (e.key=="ArrowLeft") next();
        if (e.key=="ArrowRight") prev();
        if (e.key=="ArrowUp" || e.key=="PageUp") {
            if (cp > 0) cp -= 1
            load();
        }
        if (e.key=="ArrowDown" || e.key=="PageDown") {
            if (cp < len-1) cp += 1
            load();
        }
        if (e.key=="Home") {
            cp = 0
            load()
        }
        if (e.key=="End") {
            cp = len-1
            load()
        }
    }
    var posX = 0, posY = 0, relX = 0, relY = 0, time;
    document.onpointerdown = function (e) {
        if (e.button || ['INPUT', 'DIV', 'NAV', 'SPAN', 'A'].includes(e.target.tagName.toUpperCase())) return;
        e.preventDefault();
        posX = e.clientX;
        posY = e.clientY;
        time = Date.now()
        document.onpointerup = document.onpointercancel = document.onpointerout = document.onpointerleave = function (e) {
            relX = posX - e.clientX;
            relY = posY - e.clientY;
            if (Math.abs(relX+relY) >= 16) {
                if (relX*relY == 0 && Math.abs(relX) || relX*relY != 0 && Math.abs(relX/relY) >= 1) {
                    if (relX > 0) prev();
                    if (relX < 0) next();
                } else if (relX*relY == 0 && Math.abs(relY) || relX*relY != 0 && Math.abs(relY/relX) >= 1) {
                    if (relY > 0) document.documentElement.requestFullscreen();
                    if (relY < 0) document.exitFullscreen();
                }
            } else if (Date.now() - time < 400) prev();
            document.onpointerup = document.onpointercancel = document.onpointerout = document.onpointerleave = null;
        }
    }
    document.onwheel = (e) => {
        if (e.ctrlKey) return
        e.preventDefault();
        if (e.deltaY > 0 && cp < len-1) cp+=1;
        if (e.deltaY < 0 && cp > 0) cp-=1;
        load(true);
    }
    document.body.onhashchange = () => {
        if (cp != parseInt(location.hash.slice(1))-1) {
            cp = parseInt(location.hash.slice(1))-1
            load();
        }
    }
    range.oninput = () => {
        cp = parseInt(range.value)
        load();
    }
    document.onfullscreenchange = () => {
        if (document.fullscreenElement) {
            document.querySelector('nav').style.display = 'none';
            document.getElementById('viewer').style.height = '100vh';
        } else {
            document.querySelector('nav').style.display = '';
            document.getElementById('viewer').style.height = 'calc(100vh - 50px)';
        }
        load()
    }
    document.getElementById('viewer').style.overflowY = 'hidden';
    document.getElementById('paths').style.display = 'none';
}
