(function () {    
    const dialog = document.querySelector("dialog");
    const showButton = document.getElementById("deleteButton")
    const cancelButton = document.getElementById("cancelButton")

    // "Show the dialog" button opens the dialog modally
    showButton.addEventListener("click", () => {
        dialog.showModal();
    });

    cancelButton.addEventListener("click", () => {
        dialog.close();
    });
})();