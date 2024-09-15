const init = () => {
    document.body.style.overflowY = 'scroll'
    const input = document.createElement('input')
    input.id = 'search'
    input.type = 'text'
    input.placeholder = "Search..."
    input.oninput = () => {
        const key = input.value.toLocaleLowerCase()
        document.querySelectorAll('#homeview > a').forEach((v)=>{
            if (v.getAttribute('key').toLocaleLowerCase().includes(key)) v.className = ''; else v.className = 'hidden';
        })
    }
    document.querySelector('nav').appendChild(input)
}