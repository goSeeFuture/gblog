function backToTop() {
    document.body.scrollTop = document.documentElement.scrollTop = 0;
}

function showTOC() {
    var c = document.querySelector('.toc-content');
    if (c.hasAttribute('hidden')) {
        c.removeAttribute('hidden');
    } else {
        c.setAttribute('hidden', 'hidden');
    }
}