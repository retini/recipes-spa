window.onload = () => {
    submenu = document.querySelector(".submenu");
    recipeMenu = document.querySelector("#menu-2");
    recipeMenu.onclick = () => {
        submenu.classList.toggle("active");
    }


}