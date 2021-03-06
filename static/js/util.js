function backToTop() {
    document.body.scrollTop = document.documentElement.scrollTop = 0
}

function isMobile() {
    return window.getComputedStyle(document.getElementById("float-mobile-ctl"), null).display != "none"
}

function toggleTOC() {
    if (isMobile()) {
        var c = document.querySelector('.modal')
        if (c.getAttribute('class').lastIndexOf('active') == -1) {
            c.setAttribute('class', 'modal active')
        }
    } else {
        console.log("toggleTOC pc")
        var c = document.querySelector('.toc')
        if (c.hasAttribute('hidden')) {
            c.removeAttribute('hidden')
        } else {
            c.setAttribute('hidden', 'hidden')
        }
    }
}

function closeMobileTOC() {
    var c = document.getElementById('modal-id')
    console.log("closeMobileTOC", c.getAttribute('class'))
    if (c.getAttribute('class').lastIndexOf('active') == -1) {
        c.setAttribute('class', 'modal active')
    } else {
        c.setAttribute('class', 'modal')
    }
}
