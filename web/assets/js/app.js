const app = (function () {    
    const bindingDialogButtons = () => {
        const dialog = document.querySelector("dialog");
        const showButton = document.getElementById("deleteButton")
        const cancelButton = document.getElementById("cancelButton")
        
        if (!!showButton) {
            showButton.addEventListener("click", () => {
                dialog.showModal();
            });
        }

        cancelButton.addEventListener("click", () => {
            dialog.close();
        });
    }

    htmx.onLoad(function(elt) {
        if (htmx.findAll(elt, "#contentRoot").length > 0) {
            bindingDialogButtons();
        }
    })

    return {
        closeDialog: () => {dialog.close();}
    }
})();