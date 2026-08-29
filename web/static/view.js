let wrapper = document.getElementsByClassName("top-bar__wrapper")[0],
    auth = document.getElementsByClassName("top-bar__section_auth")[0],
    highlights = document.getElementsByClassName("top-bar__section_highlights")[0]

const toggleElement = el => el.classList.toggle("top-bar__section_hidden-on-mobile")
toggleElement(auth)
toggleElement(highlights)

let hamburger = document.createElement("button")
hamburger.classList.add("top-bar__hamburger")
hamburger.onclick = _ => {
    toggleElement(auth)
    toggleElement(highlights)
}
hamburger.innerText = "Menu"

let hamburgerWrapper = document.createElement("div")
hamburgerWrapper.classList.add("top-bar__hamburger-wrapper")

let hamburgerSection = document.createElement("li")
hamburgerSection.classList.add("top-bar__section", "top-bar__section_hamburger")

hamburgerWrapper.appendChild(hamburger)
hamburgerSection.appendChild(hamburgerWrapper)
wrapper.appendChild(hamburgerSection);

(async () => {
    const input = document.querySelector('.js-add-cat-name'),
        datalist = document.querySelector('.js-add-cat-list')
    if (!input || !datalist) return;

    const categories = await fetch('/category')
        .then(resp => resp.text())
        .then(html => {
            return Array
                .from(new DOMParser()
                    .parseFromString(html, 'text/html')
                    .querySelectorAll('.mv-tag .p-name'))
                .map(a => a.innerText);
        });

    for (let cat of categories) {
        let optionElement = document.createElement('option')
        optionElement.value = cat
        datalist.appendChild(optionElement)
    }
    input.setAttribute('list', 'cat-name-options')
})();

(() => {
    const copyIconSVG = `<svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.3" aria-hidden="true"><rect x="5" y="5" width="9" height="9" rx="1.5"/><path d="M3.5 10.5h-1a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1h7a1 1 0 0 1 1 1v1"/></svg>`
    const checkIconSVG = `<svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" aria-hidden="true"><path d="M3 8.5l3 3 7-7"/></svg>`

    for (const pre of document.querySelectorAll('pre.codeblock')) {
        const codeEl = pre.querySelector('code')
        if (!codeEl) continue

        const button = document.createElement('button')
        button.type = 'button'
        button.classList.add('codeblock__copy-btn')
        button.setAttribute('aria-label', rrh.l10n('Copy to clipboard'))
        button.title = rrh.l10n('Copy to clipboard')
        button.innerHTML = copyIconSVG

        let resetTimer = null
        button.addEventListener('click', () => {
            const copying = navigator.clipboard
                ? navigator.clipboard.writeText(codeEl.textContent)
                : Promise.reject(new Error('Clipboard API unavailable'))

            copying.then(() => {
                button.innerHTML = checkIconSVG
                button.classList.add('codeblock__copy-btn_success')
                button.classList.remove('codeblock__copy-btn_error')
                button.setAttribute('aria-label', rrh.l10n('Copied!'))
            }, () => {
                button.classList.add('codeblock__copy-btn_error')
                button.classList.remove('codeblock__copy-btn_success')
                button.setAttribute('aria-label', rrh.l10n('Copy failed'))
            })

            clearTimeout(resetTimer)
            resetTimer = setTimeout(() => {
                button.innerHTML = copyIconSVG
                button.classList.remove('codeblock__copy-btn_success', 'codeblock__copy-btn_error')
                button.setAttribute('aria-label', rrh.l10n('Copy to clipboard'))
            }, 2000)
        })

        pre.appendChild(button)
    }
})();

