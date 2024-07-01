const init = () => {
    document.body.style.overflowY = 'scroll'
    const input = document.createElement('input')
    input.id = 'search'
    input.type = 'text'
    input.placeholder = "Search..."
    input.onkeyup = () => {
        const key = input.value
        document.querySelectorAll('#homeview > a').forEach((v)=>{
            if (v.title.includes(key) || v.name.includes(key)) v.className = ''; else v.className = 'hidden';
        })
    }
    document.querySelector('nav').appendChild(input)
}