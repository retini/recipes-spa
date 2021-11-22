window.onload = () => {
    setListeners()
}

function setListeners() {
    buttons = document.querySelectorAll(".link");
    buttons.forEach(btn => {
        btn.addEventListener("click", ajaxRedirect, false);
    })
    window.addEventListener("popstate", lastComponent, false);
    formBtn = document.querySelector(".ajax-search");
    if (formBtn) {
        formBtn.addEventListener("click", ajaxSubmit, false);
    }
    searchBtn = document.querySelector(".search-btn");
    if (searchBtn) {
        searchBtn.addEventListener("click", (e) => {
            overlay = document.querySelector(".overlay");
            overlay.classList.toggle("opened");
        })
    }
}

function ajaxSubmit(e) {
    e.preventDefault();
    href = e.target.getAttribute("href");
    history.pushState({"url": href}, "", href);
    i = document.querySelector("input#ingredient").value;
    q = document.querySelector("input#quantity").value;
    url = `${href}?ingredient=${i}&quantity=${q}&is-from-js=true`;
    console.log(url);
    placeComponent(url);
}

function ajaxRedirect(e) {
    e.preventDefault();
    href = e.target.getAttribute("href");
    if (href == window.location.pathname) {
        return
    }
    history.pushState({"url": href}, "", href);
    url = href + "?is-from-js=true";
    placeComponent(url);
}

async function lastComponent(e) {
    e.preventDefault();
    contentDiv = document.querySelector("div.content");
    try {
        response = await fetch(window.location.hostname + e.state.url + "?is-from-js=true");
        html = await response.text();
    } catch(err) {
        console.log(err);
    }
    contentDiv.innerHTML = html;
    setListeners();
    return
}

async function placeComponent(url) {
    contentDiv = document.querySelector("div.content");
    absUrl = window.location.hostname + url;
    try {
        response = await fetch(absUrl);
        html = await response.text();
    } catch(err) {
        console.log(err);
    }
    contentDiv.innerHTML = html;
    setListeners();
    return
}