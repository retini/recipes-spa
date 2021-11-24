window.onload = () => {
    btn = document.querySelector(".ajax-login");
    if (btn) {
        btn.addEventListener("click", ajaxAuth, false);
    }
    setListeners();
}

// callServer calls the given server method with the given parameters.
// It returns the response or throw an exception if an error occurs.
async function callServer(method, params) {
    let body = "";
    try {
        for (let i = 0; i < params.length; i++) {
            body += JSON.stringify(params[i]) + "\n";
        }
    } catch (e) {
        throw new Error("Cannot serialize arguments: " + e.message);
    }
    let res = await fetch(method + ".json", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        redirect: "error",
        body: body
    })
    if (res.status !== 200) {
        throw new Error("Unexpected HTTP status " + e.status);
    }
    let data;
    try {
        data = await res.text();
    } catch (e) {
        throw new Error("Cannot read HTTP response: " + e.message);
    }
    if ( data.length === 0 ) {
        return null;
    }
    try {
        return JSON.parse(data.substr(9));
    } catch (e) {
        throw new Error("Cannot parse JSON response: " + e.message);
    }
}

async function ajaxAuth(e) {
    e.preventDefault();
    var u = document.querySelector("input#user").value;
    var p = document.querySelector("input#password").value;
    var logged = await callServer("login", [u, p]);
    if (logged) {
        window.location = "index.html";
    }
    var errorBox = document.querySelector(".error-msg");
    errorBox.classList.add("opened");
    setTimeout(() => {
        errorBox.classList.remove("opened");
    }, 1000);
    errorBox.innerHTML = "Username o password errati.";
    return
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