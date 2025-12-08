window.addEventListener('load', () => {
    document.body.style.overflowY = 'scroll'
    const input = document.createElement('input')
    input.id = 'search'
    input.type = 'text'
    input.placeholder = "Search..."
    input.oninput = () => {
        const key = input.value.toLocaleLowerCase()
        document.querySelectorAll('#homeview > a').forEach((v)=>{
            v.classList.toggle('hidden', !v.getAttribute('key').toLocaleLowerCase().includes(key))
        })
    }
    document.querySelector('nav').appendChild(input)
})