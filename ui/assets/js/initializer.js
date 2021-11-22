buttons = document.querySelectorAll(".link");
// Aggiungi ai bottoni l'event listener.
buttons.forEach((btn) => {
    btn.addEventListener("click", getComponent, false);
})

window.addEventListener("popstate", getComponent, false);

async function getComponent(e) {
    // previeni il ricaricamento della pagina.
    e.preventDefault();
    contentDiv = document.querySelector("div.content");
    // se l'evento è un popstate allora effettua la richiesta con l'href
    // dell'ultimo componente.
    if (e.type == "popstate") {
        fetch(window.location.hostname + e.state.url + "?is-from-js=true", {credentials: "include"})
        .then(response => response.text())
        .then(html => contentDiv.innerHTML = html)
        .catch(error => console.log(error))
        return
    }
    // prendi il link del bottone.
    href = e.target.getAttribute("href");
    if (href == window.location.pathname) {
        return;
    }
    // aggiungi il link alla url della pagina, senza ricaricarla.
    history.pushState({"url": href}, "", href);
    currentUrl = window.location.hostname + href;
    // effettua la richiesta con la query is-from-js, la quale è usata dal
    // server per capire se la richiesta arriva da javascript.
    res = await fetch(currentUrl+"?is-from-js=true", {credentials: "include"});
    html = await res.text();
    contentDiv.innerHTML = html;
    
    // registra i metodi di nuovo, in modo da coprire anche gli eventuali nuovi
    // elementi ottenuti da componenti.
    buttons = document.querySelectorAll(".link");
    buttons.forEach((btn) => {
        btn.addEventListener("click", getComponent, false);
    })
}